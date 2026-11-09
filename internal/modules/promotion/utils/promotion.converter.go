package promotion_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	promotion_dtos "larsa-e-commerce-microservice-v2/internal/modules/promotion/dtos"
	promotion_models "larsa-e-commerce-microservice-v2/internal/modules/promotion/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateDTOToModel(dto *promotion_dtos.CreatePromotionDTO, additionalOptions types.AdditionalOptions) promotion_models.Promotion {
	return promotion_models.Promotion{
		BasePromotion: promotion_models.BasePromotion{
			MetaDataBaseModel:  util.CreateMetaDataBaseModel(additionalOptions),
			Name:               dto.Name,
			Description:        dto.Description,
			DiscountPercentage: dto.DiscountPercentage,
			StartDate:          dto.StartDate,
			EndDate:            dto.EndDate,
			ProductIds:         dto.ProductIds,
			AllProducts:        dto.AllProducts,
			BranchIds:          dto.BranchIds,
			IsActive:           dto.IsActive,
		},
	}
}

func ConvertUpdateDTOToModel(dto *promotion_dtos.UpdatePromotionDTO, id primitive.ObjectID, additionalOptions types.AdditionalOptions) promotion_models.Promotion {
	model := promotion_models.Promotion{
		BasePromotion: promotion_models.BasePromotion{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}

	// Set the ID for update
	model.Id = id

	// Update only provided fields
	if dto.Name != nil {
		model.Name = *dto.Name
	}
	if dto.Description != nil {
		model.Description = *dto.Description
	}
	if dto.DiscountPercentage != nil {
		model.DiscountPercentage = *dto.DiscountPercentage
	}
	if dto.StartDate != nil {
		model.StartDate = *dto.StartDate
	}
	if dto.EndDate != nil {
		model.EndDate = *dto.EndDate
	}
	if dto.ProductIds != nil {
		model.ProductIds = *dto.ProductIds
	}
	if dto.BranchIds != nil {
		model.BranchIds = *dto.BranchIds
	}
	if dto.IsActive != nil {
		model.IsActive = *dto.IsActive
	}
	if dto.AllProducts != nil {
		model.AllProducts = *dto.AllProducts
	}

	return model
}
