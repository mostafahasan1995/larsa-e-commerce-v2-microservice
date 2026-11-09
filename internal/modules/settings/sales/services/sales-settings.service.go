package settings_sales_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	settings_sales_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/dtos"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"
	settings_sales_utils "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/utils"

	"github.com/samber/do"
)

type SalesSettingsRepo interface {
	Get(ctx context.Context) (*settings_sales_models.SalesSettings, *helpers.ApiError)
	Update(ctx context.Context, model *settings_sales_models.SalesSettings) (*settings_sales_models.SalesSettings, *helpers.ApiError)
}

type SalesSettingsService interface {
	GetSalesSettings(ctx context.Context) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError)
	UpdateSalesSettings(ctx context.Context, dto settings_sales_dtos.UpdateSalesSettingsDTO) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError)
}

type salesSettingsService struct {
	repo SalesSettingsRepo
}

func NewSalesSettingsService(i *do.Injector) (SalesSettingsService, error) {
	repo := do.MustInvoke[SalesSettingsRepo](i)
	return &salesSettingsService{repo: repo}, nil
}

func (s *salesSettingsService) GetSalesSettings(ctx context.Context) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError) {
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

func (s *salesSettingsService) UpdateSalesSettings(ctx context.Context, dto settings_sales_dtos.UpdateSalesSettingsDTO) (*settings_sales_models.SalesSettingsResponse, *helpers.ApiError) {
	existingSettings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	settingsModel, convertErr := settings_sales_utils.ConvertUpdateDTOToSalesSettingsModel(ctx, dto, existingSettings)
	if convertErr != nil {
		return nil, convertErr
	}

	updatedSettings, updateErr := s.repo.Update(ctx, settingsModel)
	if updateErr != nil {
		return nil, updateErr
	}

	response := &settings_sales_models.SalesSettingsResponse{
		BaseSalesSettings: updatedSettings.BaseSalesSettings,
	}
	return response, nil
}
