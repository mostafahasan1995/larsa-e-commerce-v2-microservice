package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
* NOTE: Discount does not take into consideration the currency code
* the discount only have one currency code
 */
type Discount struct {
	Value        float64                    `bson:"value" json:"value"`
	CurrencyCode *string                    `bson:"currencyCode,omitempty" json:"currencyCode,omitempty"`
	Type         product_enums.DiscountType `bson:"type" json:"type"`
	UpTo         *[]PriceItem               `bson:"upTo,omitempty" json:"upTo,omitempty"`
}

type PriceItem struct {
	Value        float64             `bson:"value" json:"value"`
	CurrencyCode string              `bson:"currencyCode" json:"currencyCode"`
	UnitId       *primitive.ObjectID `bson:"unitId,omitempty" json:"unitId,omitempty"`
}

type BasePrice struct {
	Name           *types.Localizable[string]     `bson:"name,omitempty" json:"name,omitempty"`
	BranchIds      *[]primitive.ObjectID          `bson:"branchIds,omitempty" json:"branchIds,omitempty"`
	Items          []PriceItem                    `bson:"items" json:"items"`
	CustomerGroups *[]product_enums.CustomerGroup `bson:"customerGroups,omitempty" json:"customerGroups,omitempty"`
	Platforms      *[]product_enums.PricePlatform `bson:"platforms,omitempty" json:"platforms,omitempty"`
	Discount       *Discount                      `bson:"discount,omitempty" json:"discount,omitempty"`
	IsDefault      *bool                          `bson:"isDefault,omitempty" json:"isDefault,omitempty"`
}

type Price struct {
	BasePrice `bson:",inline"`
}

type PriceResponse struct {
	BasePrice `bson:",inline"`
}

type PricesWithPagination = types.DataWithPagination[PriceResponse]
