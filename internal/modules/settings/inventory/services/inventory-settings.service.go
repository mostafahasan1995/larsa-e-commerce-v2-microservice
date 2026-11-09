package settings_inventory_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	settings_inventory_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/dtos"
	settings_inventory_models "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/models"
	settings_inventory_utils "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/utils"

	"github.com/samber/do"
)

type InventorySettingsRepo interface {
	Get(ctx context.Context) (*settings_inventory_models.InventorySettings, *helpers.ApiError)
	Update(ctx context.Context, model *settings_inventory_models.InventorySettings) (*settings_inventory_models.InventorySettings, *helpers.ApiError)
}

type InventorySettingsService interface {
	GetInventorySettings(ctx context.Context) (*settings_inventory_models.InventorySettingsResponse, *helpers.ApiError)
	UpdateInventorySettings(ctx context.Context, dto settings_inventory_dtos.UpdateInventorySettingsDTO) (*settings_inventory_models.InventorySettingsResponse, *helpers.ApiError)
}

type inventorySettingsService struct {
	repo InventorySettingsRepo
}

func NewInventorySettingsService(i *do.Injector) (InventorySettingsService, error) {
	repo := do.MustInvoke[InventorySettingsRepo](i)
	return &inventorySettingsService{repo: repo}, nil
}

func (s *inventorySettingsService) GetInventorySettings(ctx context.Context) (*settings_inventory_models.InventorySettingsResponse, *helpers.ApiError) {
	settings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	if settings == nil {
		return &settings_inventory_models.InventorySettingsResponse{
			BaseInventorySettings: settings_inventory_models.BaseInventorySettings{
				UnitsActive: false,
			},
		}, nil
	}

	response := &settings_inventory_models.InventorySettingsResponse{
		BaseInventorySettings: settings.BaseInventorySettings,
	}
	return response, nil
}

func (s *inventorySettingsService) UpdateInventorySettings(ctx context.Context, dto settings_inventory_dtos.UpdateInventorySettingsDTO) (*settings_inventory_models.InventorySettingsResponse, *helpers.ApiError) {
	existingSettings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	settingsModel, convertErr := settings_inventory_utils.ConvertUpdateDTOToInventorySettingsModel(ctx, dto, existingSettings)
	if convertErr != nil {
		return nil, convertErr
	}

	updatedSettings, updateErr := s.repo.Update(ctx, settingsModel)
	if updateErr != nil {
		return nil, updateErr
	}

	response := &settings_inventory_models.InventorySettingsResponse{
		BaseInventorySettings: updatedSettings.BaseInventorySettings,
	}
	return response, nil
}
