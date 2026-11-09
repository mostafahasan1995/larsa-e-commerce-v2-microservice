package plan_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"
)

// TODO: replace base model with meta model
type BasePlan struct {
	Name           string              `bson:"name" json:"name"`
	Type           plan_enums.PlanType `bson:"type" json:"type"`
	TrailDays      int                 `bson:"trailDays" json:"trailDays"`
	IsDefault      bool                `bson:"isDefault" json:"isDefault"`
	IsActive       bool                `bson:"isActive" json:"isActive"`
	ExpirationDays int                 `bson:"expirationDays" json:"expirationDays"`
	SuspendDays    int                 `bson:"suspendDays" json:"suspendDays"`
	Description    string              `bson:"description" json:"description"`
	Image          string              `bson:"image" json:"image"`

	MonthlyPrices []PlanPrice `bson:"monthlyPrices" json:"monthlyPrices"`
	AnnualPrices  []PlanPrice `bson:"annualPrices" json:"annualPrices"`

	Customizations []PlanCustomization `bson:"customizations" json:"customizations"`

	models.BaseModel `bson:",inline"`
}

type Plan struct {
	BasePlan `bson:",inline"`
	Features []PlanFeature `bson:"features" json:"features"`
}

type PlanResponse struct {
	BasePlan `bson:",inline"`
	Features []PlanFeatureResponse `bson:"features" json:"features"`
}

type PlansWithPagination = types.DataWithPagination[PlanResponse]
