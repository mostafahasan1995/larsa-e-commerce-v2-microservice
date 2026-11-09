package quotation_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	agent_services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/agent"
	shipping_services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/shipping"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_repos "larsa-e-commerce-microservice-v2/internal/modules/quotation/repo"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationVendorService interface {
	CreateResponseQuotation(ctx context.Context, requestQuotationID string, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
	CreateGeneralQuotation(ctx context.Context, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError)
	UpdateQuotationStatus(ctx context.Context, quotationID string, status quotation_enums.QuotationStatus) *helpers.ApiError
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
	GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError)
	GetRequestQuotations(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError)
}

type quotationVendor struct {
	repo            quotation_repos.QuotationRepo
	profileService  profile_services.ProfileService
	agentService    agent_services.AgentService
	shippingService shipping_services.ShippingService
}

func NewQuotationVendorService(i *do.Injector) (QuotationVendorService, error) {
	return &quotationVendor{
		repo:            do.MustInvoke[quotation_repos.QuotationRepo](i),
		profileService:  do.MustInvoke[profile_services.ProfileService](i),
		agentService:    do.MustInvoke[agent_services.AgentService](i),
		shippingService: do.MustInvoke[shipping_services.ShippingService](i),
	}, nil
}

func (s *quotationVendor) CreateResponseQuotation(ctx context.Context, requestQuotationID string, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	requestQuotation, err := s.GetById(ctx, requestQuotationID, "")
	if err != nil {
		return nil, err
	}

	quotation, apiErr := quotation_utils.ConvertCreateQuotationDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	config, _ := util.GetReqAppCfg(ctx)
	quotation.Code = quotation.GenerateCode()
	quotation.Status = quotation_enums.PENDING_ADMIN
	quotation.Type = quotation_enums.RESPONSE
	quotation.VendorProfileId = &config.User.Id
	quotation.RequestQuotationId = &requestQuotation.Id

	fromIso := quotation.Delivery.FromLocation.CountryIso3
	toIso := quotation.Delivery.ToLocation.CountryIso3

	buyerInfo := requestQuotation.QuotationBuyerInfo
	quotation.QuotationBuyerInfo = &quotation_models.QuotationBuyerInfo{
		BaseQuotationBuyerInfo: buyerInfo.BaseQuotationBuyerInfo,
	}

	toLocation := requestQuotation.Delivery.ToLocation
	quotation.Delivery.ToLocation = toLocation

	s.assignAgentsAndShipping(ctx, &quotation, fromIso, toIso)

	result, apiErr := s.repo.Create(ctx, quotation)
	if apiErr != nil {
		return nil, apiErr
	}

	s.notifyAgentsAndShipping(ctx, result)

	return result, nil
}

func (s *quotationVendor) CreateGeneralQuotation(ctx context.Context, dto *quotation_dto.CreateQuotationDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotation, apiErr := quotation_utils.ConvertCreateQuotationDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	config, _ := util.GetReqAppCfg(ctx)
	quotation.Code = quotation.GenerateCode()
	quotation.Status = quotation_enums.PENDING_ADMIN
	quotation.Type = quotation_enums.GENERAL
	quotation.VendorProfileId = &config.User.Id

	vendorProfile, _ := s.profileService.GetById(ctx, *quotation.VendorProfileId, "")
	if vendorProfile != nil {
		quotation.FromCountryIso3 = vendorProfile.Country
	}

	if quotation.Delivery.FromLocation == (quotation_models.QuotationDeliveryLocation{}) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "From location is required",
		}
	}
	if quotation.Delivery.ToLocation == (quotation_models.QuotationDeliveryLocation{}) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "To location is required",
		}
	}

	fromIso := quotation.Delivery.FromLocation.CountryIso3
	toIso := quotation.Delivery.ToLocation.CountryIso3

	s.assignAgentsAndShipping(ctx, &quotation, fromIso, toIso)

	result, apiErr := s.repo.Create(ctx, quotation)
	if apiErr != nil {
		return nil, apiErr
	}

	s.notifyAgentsAndShipping(ctx, result)

	return result, nil
}

func (s *quotationVendor) GetById(ctx context.Context, id string, query string) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	objectId, parseErr := primitive.ObjectIDFromHex(id)
	if parseErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid ID format",
		}
	}
	filter, err := quotation_filters.NewQuotationFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	return s.repo.GetById(ctx, objectId, *filter)
}

func (s *quotationVendor) UpdateQuotationStatus(ctx context.Context, quotationId string, status quotation_enums.QuotationStatus) *helpers.ApiError {
	quotation, apiErr := s.GetById(ctx, quotationId, `{"popFields":["quotationBuyerInfo.buyerProfileId","vendorProfileId"]}`)
	if apiErr != nil {
		return apiErr
	}
	config, _ := util.GetReqAppCfg(ctx)

	if !quotation_utils.IsUpdateQuotationStatusValid(quotation, status) {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid status",
		}
	}

	if !quotation_utils.IsUserAllowedToUpdateQuotationStatus(quotation, status, config.User) {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "User is not allowed to update quotation status",
		}
	}

	historyEntry := quotation_utils.GetHistoryEntry(quotation, status)

	newHistory := append(quotation.History, *historyEntry)

	_, apiErr = s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{
		"history": newHistory,
	}, config.User.Id)
	if apiErr != nil {
		return apiErr
	}

	_, apiErr = s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{
		"status": status,
	}, config.User.Id)
	if apiErr != nil {
		return apiErr
	}

	//TODO:
	// if status == quotation_enums.CONVERT_TO_ORDER {
	// 	quotation.Status = status
	// 	_, apiErr := s.ConvertToOrder(ctx, quotation)
	// 	if apiErr != nil {
	// 		return apiErr
	// 	}
	// 	_, apiErr = s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{
	// 		"status": quotation_enums.CONVERTED_TO_ORDER,
	// 	}, config.User.Id)
	// 	if apiErr != nil {
	// 		return apiErr
	// 	}
	// }
	return nil
}

func (s *quotationVendor) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	filter, err := quotation_filters.NewQuotationFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	filter.VendorProfileId = &config.User.Id

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *quotationVendor) GetRequestQuotations(ctx context.Context, skip, limit int64, query string) (*quotation_models.QuotationsWithPagination, *helpers.ApiError) {
	filter, err := quotation_filters.NewQuotationFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	//TODO: check if we show all request quotations for all vendors
	requestType := quotation_enums.REQUEST
	filter.Type = &requestType

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

// assignAgentsAndShipping sets AgentFromId, AgentToId, and ShippingCompanyId
// based on the provided ISO3 country codes and current services configuration.
func (s *quotationVendor) assignAgentsAndShipping(ctx context.Context, quotation *quotation_models.Quotation, fromIso string, toIso string) {
	// Find agent for From country (any agent authorized for that country)
	var fromAgentSupportsShipping bool
	if fromIso != "" {
		query := `{"authorizedCountryCodes":["` + fromIso + `"]}`
		agentsFrom, _ := s.agentService.GetWithPagination(ctx, 0, 1, query)
		if agentsFrom != nil && len(agentsFrom.Data) > 0 {
			quotation.AgentFromId = &agentsFrom.Data[0].Id
			fromAgentSupportsShipping = agentsFrom.Data[0].ResponsibleForLogistics
		}
	}

	// Find agent for To country
	if toIso != "" {
		query := `{"authorizedCountryCodes":["` + toIso + `"]}`
		agentsTo, _ := s.agentService.GetWithPagination(ctx, 0, 1, query)
		if agentsTo != nil && len(agentsTo.Data) > 0 {
			quotation.AgentToId = &agentsTo.Data[0].Id
		}
	}

	// Set shipping company: if agent from supports shipping => keep nil; otherwise find company by from country
	if fromIso != "" && !fromAgentSupportsShipping {
		shipping, _ := s.shippingService.FindOneByCountryCode(ctx, fromIso)
		if shipping != nil {
			quotation.ShippingCompanyId = &shipping.Id
		}
	}
}

// notifyAgentsAndShipping is a placeholder for future notifications to agents and shipping companies
func (s *quotationVendor) notifyAgentsAndShipping(ctx context.Context, quotation *quotation_models.QuotationResponse) {
	//TODO: implement notification
	// intentionally left blank for future implementation
}
