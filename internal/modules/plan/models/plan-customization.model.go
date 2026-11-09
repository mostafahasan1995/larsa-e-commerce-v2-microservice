package plan_models

import (
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasePlanCustomization struct {
	Type     plan_enums.PlanCustomizationType `bson:"type" json:"type"`
	UsersIds []primitive.ObjectID             `bson:"usersIds" json:"usersIds"`
}

type PlanCustomization struct {
	BasePlanCustomization `bson:",inline"`
}

type PlanCustomizationResponse struct {
	BasePlanCustomization `bson:",inline"`
	Users                 []*primitive.ObjectID `bson:"users" json:"users"`
}
