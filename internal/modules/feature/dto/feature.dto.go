package feature_dto

import (
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrencyDTO struct {
	CurrencyId primitive.ObjectID `json:"currencyId" validate:"required"`
	Price      float64            `json:"price" validate:"required,min=0"`
}

type MonthlyPriceDTO struct {
	Active     bool          `json:"active"`
	Currencies []CurrencyDTO `json:"currencies" validate:"dive"`
}

type AnnualPriceDTO struct {
	Active     bool          `json:"active"`
	Currencies []CurrencyDTO `json:"currencies" validate:"dive"`
}

type PricePerUserDTO struct {
	Active     bool          `json:"active"`
	Currencies []CurrencyDTO `json:"currencies" validate:"dive"`
}

type CreateFeatureDTO struct {
	Name                string               `json:"name" validate:"required"`
	Slug                string               `json:"slug" validate:"required"`
	GroupId             *primitive.ObjectID  `json:"groupId" `
	Service             plan_enums.PlanType  `json:"service"`
	Description         string               `json:"description"`
	HasQuantity         bool                 `json:"hasQuantity"`
	IsActive            bool                 `json:"isActive"`
	Units               string               `json:"units"`
	HasSubscription     bool                 `json:"hasSubscription"`
	Image               string               `json:"image"`
	MonthlyPrice        MonthlyPriceDTO      `json:"monthlyPrice"`
	AnnualPrice         AnnualPriceDTO       `json:"annualPrice"`
	PricePerUser        PricePerUserDTO      `json:"pricePerUser"`
	FeatureDependencies []primitive.ObjectID `json:"featureDependencies"`
}

type UpdateFeatureDTO struct {
	Name                string               `json:"name" validate:"required"`
	Slug                string               `json:"slug" validate:"required"`
	GroupId             *primitive.ObjectID  `json:"groupId"`
	Service             plan_enums.PlanType  `json:"service"`
	Description         string               `json:"description"`
	HasQuantity         bool                 `json:"hasQuantity"`
	IsActive            bool                 `json:"isActive"`
	Units               string               `json:"units"`
	HasSubscription     bool                 `json:"hasSubscription"`
	Image               string               `json:"image"`
	MonthlyPrice        MonthlyPriceDTO      `json:"monthlyPrice"`
	AnnualPrice         AnnualPriceDTO       `json:"annualPrice"`
	PricePerUser        PricePerUserDTO      `json:"pricePerUser"`
	FeatureDependencies []primitive.ObjectID `json:"featureDependencies"`
}
