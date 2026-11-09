package settings_general_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	settings_general_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/general/dtos"
	settings_general_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/general/enums"
	settings_general_models "larsa-e-commerce-microservice-v2/internal/modules/settings/general/models"
	settings_general_utils "larsa-e-commerce-microservice-v2/internal/modules/settings/general/utils"

	"github.com/samber/do"
)

type GeneralSettingsRepo interface {
	Get(ctx context.Context) (*settings_general_models.GeneralSettings, *helpers.ApiError)
	Update(ctx context.Context, model *settings_general_models.GeneralSettings) (*settings_general_models.GeneralSettings, *helpers.ApiError)
}

type GeneralSettingsService interface {
	GetGeneralSettings(ctx context.Context) (*settings_general_models.GeneralSettingsResponse, *helpers.ApiError)
	UpdateGeneralSettings(ctx context.Context, dto settings_general_dtos.UpdateGeneralSettingsDTO) (*settings_general_models.GeneralSettingsResponse, *helpers.ApiError)
}

type generalSettingsService struct {
	repo GeneralSettingsRepo
}

func NewGeneralSettingsService(i *do.Injector) (GeneralSettingsService, error) {
	repo := do.MustInvoke[GeneralSettingsRepo](i)
	return &generalSettingsService{repo: repo}, nil
}

func (s *generalSettingsService) GetGeneralSettings(ctx context.Context) (*settings_general_models.GeneralSettingsResponse, *helpers.ApiError) {
	settings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	if settings == nil {
		return &settings_general_models.GeneralSettingsResponse{
			BaseGeneralSettings: settings_general_models.BaseGeneralSettings{
				SupportedLanguages: []settings_general_enums.SupportedLanguage{},
			},
		}, nil
	}

	response := &settings_general_models.GeneralSettingsResponse{
		BaseGeneralSettings: settings.BaseGeneralSettings,
	}
	return response, nil
}

func (s *generalSettingsService) UpdateGeneralSettings(ctx context.Context, dto settings_general_dtos.UpdateGeneralSettingsDTO) (*settings_general_models.GeneralSettingsResponse, *helpers.ApiError) {
	existingSettings, apiErr := s.repo.Get(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	settingsModel, convertErr := settings_general_utils.ConvertUpdateDTOToGeneralSettingsModel(ctx, dto, existingSettings)
	if convertErr != nil {
		return nil, convertErr
	}

	updatedSettings, updateErr := s.repo.Update(ctx, settingsModel)
	if updateErr != nil {
		return nil, updateErr
	}

	response := &settings_general_models.GeneralSettingsResponse{
		BaseGeneralSettings: updatedSettings.BaseGeneralSettings,
	}
	return response, nil
}
