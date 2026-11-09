package util

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMetaDataBaseModel(additionalOptions types.AdditionalOptions) base_models.MetaDataBaseModel {
	metaDataBaseModel := base_models.MetaDataBaseModel{
		UpdatedAt: time.Now(),
		UpdatedBy: &additionalOptions.UserId,
	}

	if additionalOptions.Mode == types.CREATE {
		metaDataBaseModel.Id = primitive.NewObjectID()
		metaDataBaseModel.CreatedAt = time.Now()
		metaDataBaseModel.CreatedBy = &additionalOptions.UserId
	}

	return metaDataBaseModel
}
