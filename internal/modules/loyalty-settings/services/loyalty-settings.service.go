package loyalty_settings_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_settings_dtos "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/dtos"
	loyalty_settings_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/models"
	loyalty_settings_utils "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoyaltySettingsService interface {
	Upsert(ctx context.Context, dto *loyalty_settings_dtos.CreateLoyaltySettingsDTO) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
	GetActiveSettings(ctx context.Context) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
	UpdateSettings(ctx context.Context, dto *loyalty_settings_dtos.CreateLoyaltySettingsDTO, id primitive.ObjectID) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
}

type LoyaltySettingsRepo interface {
	Create(ctx context.Context, loyaltySettings loyalty_settings_models.LoyaltySettings) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
	GetActiveSettings(ctx context.Context) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, loyaltySettings loyalty_settings_models.LoyaltySettings, filter bson.M) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError)
}

type loyaltySettingsService struct {
	repo LoyaltySettingsRepo
}

func NewLoyaltySettingsService(i *do.Injector) (LoyaltySettingsService, error) {
	return &loyaltySettingsService{
		repo: do.MustInvoke[LoyaltySettingsRepo](i),
	}, nil
}

func (s *loyaltySettingsService) Upsert(ctx context.Context, dto *loyalty_settings_dtos.CreateLoyaltySettingsDTO) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError) {
	existing, err := s.repo.GetActiveSettings(ctx)

	if err != nil {
		if err.StatusCode == 404 {
			loyaltySettings := loyalty_settings_utils.ConvertCreateLoyaltySettingsDTOToModel(dto, types.AdditionalOptions{})
			return s.repo.Create(ctx, loyaltySettings)
		}
		return nil, err
	}

	return s.UpdateSettings(ctx, dto, existing.Id)
}

func (s *loyaltySettingsService) GetActiveSettings(ctx context.Context) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError) {
	return s.repo.GetActiveSettings(ctx)
}

func (s *loyaltySettingsService) UpdateSettings(ctx context.Context, dto *loyalty_settings_dtos.CreateLoyaltySettingsDTO, id primitive.ObjectID) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	loyaltySettings := loyalty_settings_utils.ConvertCreateLoyaltySettingsDTOToModel(dto, types.AdditionalOptions{
		Mode:   types.UPDATE,
		UserId: cfg.User.Id,
	})
	loyaltySettings.Id = id

	return s.repo.UpdateById(ctx, id, loyaltySettings, bson.M{})
}
