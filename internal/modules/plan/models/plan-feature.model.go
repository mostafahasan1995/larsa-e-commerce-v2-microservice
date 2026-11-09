package plan_models

import (
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasePlanFeature struct {
	FeatureId primitive.ObjectID `bson:"featureId" json:"featureId"`
	Active    bool               `bson:"active" json:"active"`
	Value     string             `bson:"value" json:"value"`
}

type PlanFeature struct {
	BasePlanFeature `bson:",inline"`
}

type PlanFeatureResponse struct {
	BasePlanFeature `bson:",inline"`
	FeatureObject   *feature_models.FeatureResponse `bson:"featureObject" json:"featureObject"`
}
