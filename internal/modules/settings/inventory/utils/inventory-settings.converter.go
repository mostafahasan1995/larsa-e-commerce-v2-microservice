package settings_inventory_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_inventory_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/dtos"
	settings_inventory_models "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertUpdateDTOToInventorySettingsModel(
	ctx context.Context,
	dto settings_inventory_dtos.UpdateInventorySettingsDTO,
	existingModel *settings_inventory_models.InventorySettings,
) (*settings_inventory_models.InventorySettings, *helpers.ApiError) {
	now := time.Now()
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return &settings_inventory_models.InventorySettings{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	userId := cfg.User.Id

	modelToUpdate := settings_inventory_models.InventorySettings{}
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

	modelToUpdate.UnitsActive = dto.UnitsActive

	return &modelToUpdate, nil
}
