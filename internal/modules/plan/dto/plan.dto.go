package plan_dto

import (
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePlanDTO struct {
	Name           string                    `json:"name" validate:"required"`
	Type           plan_enums.PlanType       `json:"type" validate:"required"`
	TrailDays      int                       `json:"trailDays"`
	IsDefault      bool                      `json:"isDefault"`
	IsActive       bool                      `json:"isActive"`
	ExpirationDays int                       `json:"expirationDays"`
	SuspendDays    int                       `json:"suspendDays"`
	Description    string                    `json:"description"`
	Image          string                    `json:"image"`
	MonthlyPrices  []CreatePlanPrice         `json:"monthlyPrices"`
	AnnualPrices   []CreatePlanPrice         `json:"annualPrices"`
	Features       []CreatePlanFeature       `json:"features"`
	Customizations []CreatePlanCustomization `json:"customizations"`
}

type UpdatePlanDTO struct {
	Name           *string                   `json:"name"`
	Type           *plan_enums.PlanType      `json:"type"`
	TrailDays      *int                      `json:"trailDays"`
	IsDefault      *bool                     `json:"isDefault"`
	IsActive       *bool                     `json:"isActive"`
	ExpirationDays *int                      `json:"expirationDays"`
	SuspendDays    *int                      `json:"suspendDays"`
	Description    *string                   `json:"description"`
	Image          *string                   `json:"image"`
	MonthlyPrices  []UpdatePlanPrice         `json:"monthlyPrices"`
	AnnualPrices   []UpdatePlanPrice         `json:"annualPrices"`
	Features       []UpdatePlanFeature       `json:"features"`
	Customizations []UpdatePlanCustomization `json:"customizations"`
}

type CreatePlanPrice struct {
	CurrencyId primitive.ObjectID `json:"currencyId" validate:"required"`
	Price      float64            `json:"price" validate:"required"`
}

type UpdatePlanPrice struct {
	CurrencyId *primitive.ObjectID `json:"currencyId,omitempty"`
	Price      *float64            `json:"price,omitempty"`
}

type CreatePlanFeature struct {
	FeatureId primitive.ObjectID `json:"featureId" validate:"required"`
	Active    bool               `json:"active"`
	Value     string             `json:"value"`
}

type UpdatePlanFeature struct {
	FeatureId *primitive.ObjectID `json:"featureId,omitempty"`
	Active    *bool               `json:"active,omitempty"`
	Value     *string             `json:"value,omitempty"`
}

type CreatePlanCustomization struct {
	Type     plan_enums.PlanCustomizationType `json:"type"`
	UsersIds []primitive.ObjectID             `json:"usersIds"`
}

type UpdatePlanCustomization struct {
	Type     *plan_enums.PlanCustomizationType `json:"type,omitempty"`
	UsersIds []primitive.ObjectID              `json:"usersIds,omitempty"`
}
