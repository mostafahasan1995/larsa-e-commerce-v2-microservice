package feature_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	feature_dto "larsa-e-commerce-microservice-v2/internal/modules/feature/dto"
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateDTOToModel(dto *feature_dto.CreateFeatureDTO, userId primitive.ObjectID) (feature_models.Feature, *helpers.ApiError) {
	now := time.Now()

	return feature_models.Feature{
		BaseFeature: feature_models.BaseFeature{
			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        primitive.NewObjectID(),
				CreatedAt: now,
				CreatedBy: &userId,
				UpdatedAt: now,
				UpdatedBy: &userId,
			},
			Name:            dto.Name,
			Slug:            dto.Slug,
			GroupId:         dto.GroupId,
			Service:         dto.Service,
			Description:     dto.Description,
			HasQuantity:     dto.HasQuantity,
			IsActive:        dto.IsActive,
			Unit:            dto.Units,
			HasSubscription: dto.HasSubscription,
			Image:           dto.Image,
		},
		MonthlyPrice: feature_models.FeaturePrice{
			Active:     dto.MonthlyPrice.Active,
			Currencies: convertCurrencyDTOs(dto.MonthlyPrice.Currencies),
		},
		AnnualPrice: feature_models.FeaturePrice{
			Active:     dto.AnnualPrice.Active,
			Currencies: convertCurrencyDTOs(dto.AnnualPrice.Currencies),
		},
		PricePerUser: feature_models.FeaturePrice{
			Active:     dto.PricePerUser.Active,
			Currencies: convertCurrencyDTOs(dto.PricePerUser.Currencies),
		},
	}, nil
}

func ConvertUpdateDTOToModel(dto *feature_dto.UpdateFeatureDTO, id primitive.ObjectID, userId primitive.ObjectID) (feature_models.Feature, *helpers.ApiError) {
	now := time.Now()

	return feature_models.Feature{
		BaseFeature: feature_models.BaseFeature{
			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        id,
				UpdatedAt: now,
				UpdatedBy: &userId,
			},
			Name:            dto.Name,
			Slug:            dto.Slug,
			GroupId:         dto.GroupId,
			Service:         dto.Service,
			Description:     dto.Description,
			HasQuantity:     dto.HasQuantity,
			IsActive:        dto.IsActive,
			Unit:            dto.Units,
			HasSubscription: dto.HasSubscription,
			Image:           dto.Image,
		},
		MonthlyPrice: feature_models.FeaturePrice{
			Active:     dto.MonthlyPrice.Active,
			Currencies: convertCurrencyDTOs(dto.MonthlyPrice.Currencies),
		},
		AnnualPrice: feature_models.FeaturePrice{
			Active:     dto.AnnualPrice.Active,
			Currencies: convertCurrencyDTOs(dto.AnnualPrice.Currencies),
		},
		PricePerUser: feature_models.FeaturePrice{
			Active:     dto.PricePerUser.Active,
			Currencies: convertCurrencyDTOs(dto.PricePerUser.Currencies),
		},
	}, nil
}

func convertCurrencyDTOs(dtos []feature_dto.CurrencyDTO) []feature_models.FeaturePriceCurrency {
	models := make([]feature_models.FeaturePriceCurrency, len(dtos))
	for i, dto := range dtos {
		models[i] = feature_models.FeaturePriceCurrency{
			CurrencyId: dto.CurrencyId,
			Price:      dto.Price,
		}
	}
	return models
}
