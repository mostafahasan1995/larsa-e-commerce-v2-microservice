package product_utils

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertAddonValueDtoToModel converts an AddonValueDto to an AddonValue model
func ConvertAddonValueDtoToModel(dto dtos.AddonValueDto) models.AddonValue {
	return models.AddonValue{
		Id:           primitive.NewObjectID(),
		Label:        dto.Label,
		Image:        dto.Image,
		Price:        dto.Price,
		MinSelection: dto.MinSelection,
		MaxSelection: dto.MaxSelection,
	}
}

// ConvertUpdateAddonValueDtoToModel converts an UpdateAddonValueDto to an AddonValue model
func ConvertUpdateAddonValueDtoToModel(dto dtos.UpdateAddonValueDto) models.AddonValue {
	addonValue := ConvertAddonValueDtoToModel(dto.AddonValueDto)

	if dto.Id != nil {
		addonValue.Id = *dto.Id
	}

	return addonValue
}

// ConvertToAddonModel converts a CreateAddonDto to an Addon model
func ConvertToAddonModel(dto *dtos.CreateAddonDto, options types.AdditionalOptions) *models.Addon {
	now := time.Now()
	addonValues := make([]models.AddonValue, len(dto.AddonValues))

	for i, valueDto := range dto.AddonValues {
		addonValues[i] = ConvertAddonValueDtoToModel(valueDto)
	}

	addon := &models.Addon{
		BaseAddon: models.BaseAddon{
			Label:        dto.Label,
			MinSelection: dto.MinSelection,
			MaxSelection: dto.MaxSelection,
			AddonValues:  addonValues,
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id: primitive.NewObjectID(),
			},
		},
	}

	if options.Mode == types.CREATE {
		addon.CreatedAt = now
		addon.UpdatedAt = now
		addon.CreatedBy = &options.UserId
		addon.UpdatedBy = &options.UserId
	} else if options.Mode == types.UPDATE {
		addon.UpdatedAt = now
		addon.UpdatedBy = &options.UserId
	}

	return addon
}
