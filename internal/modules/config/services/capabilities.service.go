package config_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"reflect"
	"sync"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
)

type CapabilitiesService interface {
	SetupCapabilities(ctx context.Context, headerParams *common.HeaderParams) error
}

type capabilitiesService struct{}

func NewCapabilitiesService(i *do.Injector) (CapabilitiesService, error) {
	return &capabilitiesService{}, nil
}

func (s *capabilitiesService) SetupCapabilities(ctx context.Context, headerParams *common.HeaderParams) error {
	registry := enums.CapabilitiesRegistry

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	registryValue := reflect.ValueOf(registry)
	registryType := registryValue.Type()

	for i := 0; i < registryValue.NumField(); i++ {
		fieldName := registryType.Field(i).Name
		fieldValue := registryValue.Field(i)

		if fieldName == "General" {
			continue
		}

		if fieldValue.Type() == reflect.TypeOf(enums.ModuleCapability{}) {
			moduleCapability := fieldValue.Interface().(enums.ModuleCapability)

			wg.Add(1)
			go s.processUserTypeOps(moduleCapability.SuperAdmin, &wg, headerParams, errChan)

			wg.Add(1)
			go s.processUserTypeOps(moduleCapability.Vendor, &wg, headerParams, errChan)

			wg.Add(1)
			go s.processUserTypeOps(moduleCapability.Member, &wg, headerParams, errChan)

			wg.Add(1)
			go s.processUserTypeOps(moduleCapability.Agent, &wg, headerParams, errChan)

			wg.Add(1)
			go s.processUserTypeOps(moduleCapability.ShippingCompany, &wg, headerParams, errChan)

		}
	}

	if registry.General != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.registerCapability(registry.General, headerParams, errChan)
		}()
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func (s *capabilitiesService) processUserTypeOps(ops enums.UserTypeOps, wg *sync.WaitGroup, headerParams *common.HeaderParams, errChan chan<- error) {
	defer wg.Done()

	opsValue := reflect.ValueOf(ops)
	for i := 0; i < opsValue.NumField(); i++ {
		field := opsValue.Field(i)
		if field.IsValid() {
			operation := field.Interface().(enums.Operation)
			capabilityString := operation.String()
			if capabilityString != "" {
				s.registerCapability(capabilityString, headerParams, errChan)
			}
		}
	}
}

func (s *capabilitiesService) registerCapability(capability string, headerParams *common.HeaderParams, errChan chan<- error) {
	if capability == "" {
		return
	}

	reqOp := &common.RequestParams{
		Service: "users",
		Path:    "capabilities",
		Method:  "POST",
		Data: map[string]interface{}{
			"label": capability,
			"name":  capability,
		},
		Header: headerParams,
	}

	_, err := common.CallService(reqOp)
	if err != nil {
		errChan <- fmt.Errorf("failed to create capability %s: %v", capability, err)
	}
}
