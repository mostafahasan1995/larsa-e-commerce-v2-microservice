package feature_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseFeature struct {
	models.MetaDataBaseModel `bson:",inline"`

	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	Slug        string              `bson:"slug" json:"slug"`
	GroupId     *primitive.ObjectID `bson:"groupId" json:"groupId"`

	Service plan_enums.PlanType `bson:"service" json:"service"`

	HasQuantity     bool   `bson:"hasQuantity" json:"hasQuantity"`
	IsActive        bool   `bson:"isActive" json:"isActive"`
	Unit            string `bson:"unit" json:"unit"`
	HasSubscription bool   `bson:"hasSubscription" json:"hasSubscription"`
	Image           string `bson:"image" json:"image"`

	FeatureDependencies []primitive.ObjectID `bson:"featureDependencies" json:"featureDependencies"`
}

type Feature struct {
	BaseFeature  `bson:",inline"`
	MonthlyPrice FeaturePrice `bson:"monthlyPrice" json:"monthlyPrice"`
	AnnualPrice  FeaturePrice `bson:"annualPrice" json:"annualPrice"`
	PricePerUser FeaturePrice `bson:"pricePerUser" json:"pricePerUser"`
}

type FeatureResponse struct {
	BaseFeature         `bson:",inline"`
	Group               *FeatureGroupResponse `bson:"group" json:"group,omitempty"`
	FeatureDependencies []*FeatureResponse    `bson:"featureDependencies" json:"featureDependencies"`
	MonthlyPrice        FeaturePriceResponse  `bson:"monthlyPrice" json:"monthlyPrice"`
	AnnualPrice         FeaturePriceResponse  `bson:"annualPrice" json:"annualPrice"`
	PricePerUser        FeaturePriceResponse  `bson:"pricePerUser" json:"pricePerUser"`
}

type FeaturesWithPagination = types.DataWithPagination[FeatureResponse]
