package settings_general_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_general_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/general/dtos"
	settings_general_models "larsa-e-commerce-microservice-v2/internal/modules/settings/general/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertUpdateDTOToGeneralSettingsModel(
	ctx context.Context,
	dto settings_general_dtos.UpdateGeneralSettingsDTO,
	existingModel *settings_general_models.GeneralSettings,
) (*settings_general_models.GeneralSettings, *helpers.ApiError) {
	now := time.Now()
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return &settings_general_models.GeneralSettings{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	userId := cfg.User.Id

	modelToUpdate := settings_general_models.GeneralSettings{}
	if existingModel != nil {
		modelToUpdate = *existingModel
		modelToUpdate.UpdatedBy = &userId
		modelToUpdate.UpdatedAt = now
	} else {
		modelToUpdate.Id = primitive.NewObjectID()
		modelToUpdate.CreatedBy = &userId
		modelToUpdate.CreatedAt = now
		modelToUpdate.UpdatedBy = &userId
		modelToUpdate.UpdatedAt = now
	}

	modelToUpdate.SupportedLanguages = dto.SupportedLanguages
	modelToUpdate.FeedbackCoupon = dto.FeedbackCoupon

	return &modelToUpdate, nil
}
