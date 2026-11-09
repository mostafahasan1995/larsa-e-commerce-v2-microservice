package feature_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BaseFeatureGroup struct {
	Name string `bson:"name" json:"name"`

	models.MetaDataBaseModel `bson:",inline"`
}

type FeatureGroupModel struct {
	BaseFeatureGroup `bson:",inline"`
}

type FeatureGroupResponse struct {
	BaseFeatureGroup `bson:",inline"`

	FeaturesCount int               `bson:"featuresCount" json:"featuresCount"`
	Features      []FeatureResponse `bson:"features" json:"features"`
}

type FeatureGroupsWithPagination = types.DataWithPagination[FeatureGroupResponse]
