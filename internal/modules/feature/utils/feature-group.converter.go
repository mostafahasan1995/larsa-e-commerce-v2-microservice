package feature_utils

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	feature_dto "larsa-e-commerce-microservice-v2/internal/modules/feature/dto"
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateGroupDTOToModel(dto *feature_dto.CreateFeatureGroupDTO, ctx context.Context) (feature_models.FeatureGroupModel, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return feature_models.FeatureGroupModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	now := time.Now()

	return feature_models.FeatureGroupModel{
		BaseFeatureGroup: feature_models.BaseFeatureGroup{
			Name: dto.Name,

			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        primitive.NewObjectID(),
				CreatedAt: now,
				CreatedBy: &cfg.User.Id,
				UpdatedAt: now,
				UpdatedBy: &cfg.User.Id,
			},
		},
	}, nil
}

func ConvertUpdateGroupDTOToModel(dto *feature_dto.UpdateFeatureGroupDTO, id primitive.ObjectID, ctx context.Context) (feature_models.FeatureGroupModel, *helpers.ApiError) {

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return feature_models.FeatureGroupModel{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	now := time.Now()

	return feature_models.FeatureGroupModel{
		BaseFeatureGroup: feature_models.BaseFeatureGroup{
			Name: dto.Name,

			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        id,
				UpdatedAt: now,
				UpdatedBy: &cfg.User.Id,
			},
		},
	}, nil
}
