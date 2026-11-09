package settings_sales_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"

	"github.com/samber/do"
)

type SalesSettingsWebRepo interface {
	Get(ctx context.Context) (*settings_sales_models.SalesSettings, *helpers.ApiError)
}

type SalesSettingsWebService interface {
	GetSalesSettings(ctx context.Context) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError)
}

type salesSettingsWebService struct {
	repo SalesSettingsWebRepo
}

func NewSalesSettingsWebService(i *do.Injector) (SalesSettingsWebService, error) {
	repo := do.MustInvoke[SalesSettingsWebRepo](i)
	return &salesSettingsWebService{repo: repo}, nil
}

func (s *salesSettingsWebService) GetSalesSettings(ctx context.Context) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError) {
	settings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	if settings == nil {
		return &settings_sales_models.SalesSettingsResponse{
			BaseSalesSettings: settings_sales_models.BaseSalesSettings{
				OrderStatuses: []settings_sales_models.OrderStatus{},
			},
		}, nil
	}

	response := &settings_sales_models.SalesSettingsResponse{
		BaseSalesSettings: settings.BaseSalesSettings,
	}
	return response, nil
}
