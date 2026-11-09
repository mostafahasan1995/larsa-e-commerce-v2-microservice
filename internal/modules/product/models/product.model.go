package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	unit_models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAddonValue struct {
	Id         primitive.ObjectID `bson:"id" json:"id"`
	TotalPrice float64            `bson:"totalPrice" json:"totalPrice"`
	Quantity   uint64             `bson:"quantity" json:"quantity"`
}

type ProductPromotionInfo struct {
	PromotionId            primitive.ObjectID `json:"promotionId"`
	Name                   string             `json:"name"`
	Description            string             `json:"description"`
	DiscountPercentage     float64            `json:"discountPercentage"`
	StartDate              string             `json:"startDate"`
	EndDate                string             `json:"endDate"`
	PromotionDiscountValue float64            `json:"promotionDiscountValue"`
}

type CalculateProductPriceResult struct {
	List            *[]CalculateProductPriceResult `json:"list,omitempty"`
	SubTotalPrice   float64                        `json:"subTotalPrice"`
	TotalPrice      float64                        `json:"totalPrice"`
	DiscountValue   float64                        `json:"discountValue"`
	PriceItem       PriceItem                      `json:"priceItem"`
	Discount        *Discount                      `json:"discount"`
	AddonValues     []ProductAddonValue            `json:"addonValues,omitempty"`
	AddonTotalPrice float64                        `json:"addonTotalPrice,omitempty"`
	Addons          *[]AddonResponse               `json:"addons,omitempty"`
	ActivePromotion *ProductPromotionInfo          `json:"activePromotion,omitempty"`
}

type ProductAttribute struct {
	Code  string                             `bson:"code" json:"code"`
	Value interface{}                        `bson:"value" json:"value"`
	Type  product_enums.ProductAttributeType `bson:"type" json:"type"`
}
type BaseProduct struct {
	Name        types.Localizable[string] `bson:"name,omitempty" json:"name,omitempty"`
	Description types.Localizable[string] `bson:"description,omitempty" json:"description,omitempty"`

	CategoryIds      []primitive.ObjectID `bson:"categoryIds,omitempty" json:"categoryIds,omitempty"`
	UnitIds          []primitive.ObjectID `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
	Attributes       []ProductAttribute   `bson:"attributes,omitempty" json:"attributes,omitempty"`
	Images           []types.FileFieldDto `bson:"images,omitempty" json:"images,omitempty"`
	ProductionLineId *primitive.ObjectID  `bson:"productionLineId,omitempty" json:"productionLineId,omitempty"`
	Prices           []Price              `bson:"prices,omitempty" json:"prices,omitempty"`

	//NOTE: Flag is unused field, we need to remove it
	Flag                          *product_enums.ProductFlag `bson:"flag,omitempty" json:"flag,omitempty"`
	ProductType                   product_enums.ProductType  `bson:"productType,omitempty" json:"productType,omitempty"`
	VendorId                      *primitive.ObjectID        `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	InventoryProductIds           *[]primitive.ObjectID      `bson:"inventoryProductIds,omitempty" json:"inventoryProductIds,omitempty"`
	base_models.MetaDataBaseModel `bson:",inline"`
}
type Product struct {
	BaseProduct `bson:",inline"`
}
type ProductResponse struct {
	BaseProduct     `bson:",inline"`
	IsInWatchListed *bool                             `bson:"isInWatchListed,omitempty" json:"isInWatchListed,omitempty"`
	ProductionLine  *ProductionLineResponse           `bson:"productionLine,omitempty" json:"productionLine,omitempty"`
	Units           *[]unit_models.UnitResponse       `bson:"units,omitempty" json:"units,omitempty"`
	Categories      *[]config_models.CategoryResponse `bson:"categories,omitempty" json:"categories,omitempty"`
	Branches        *[]BranchResponse                 `bson:"branches,omitempty" json:"branches,omitempty"`
	Price           *CalculateProductPriceResult      `bson:"price,omitempty" json:"price,omitempty"`
	Reviews         *[]ReviewResponse                 `bson:"reviews,omitempty" json:"reviews,omitempty"`
	AverageRating   *float64                          `bson:"averageRating,omitempty" json:"averageRating,omitempty"`
	TotalReviews    *int64                            `bson:"totalReviews,omitempty" json:"totalReviews,omitempty"`
	OutOfStock      *bool                             `bson:"outOfStock,omitempty" json:"outOfStock,omitempty"`
	//NOTE: This is a workaround to avoid circular reference
	InventoryProducts []InventoryProductResponse `bson:"inventoryProducts" json:"inventoryProducts,omitempty"`
}

type InventoryProductResponse struct {
	BaseProduct     `bson:",inline"`
	IsInWatchListed *bool                             `bson:"isInWatchListed,omitempty" json:"isInWatchListed,omitempty"`
	ProductionLine  *ProductionLineResponse           `bson:"productionLine,omitempty" json:"productionLine,omitempty"`
	Units           *[]unit_models.UnitResponse       `bson:"units,omitempty" json:"units,omitempty"`
	Categories      *[]config_models.CategoryResponse `bson:"categories,omitempty" json:"categories,omitempty"`
	Branches        *[]BranchResponse                 `bson:"branches,omitempty" json:"branches,omitempty"`
	Price           *CalculateProductPriceResult      `bson:"price,omitempty" json:"price,omitempty"`
	Reviews         *[]ReviewResponse                 `bson:"reviews,omitempty" json:"reviews,omitempty"`
	AverageRating   *float64                          `bson:"averageRating,omitempty" json:"averageRating,omitempty"`
	TotalReviews    *int64                            `bson:"totalReviews,omitempty" json:"totalReviews,omitempty"`
	OutOfStock      *bool                             `bson:"outOfStock,omitempty" json:"outOfStock,omitempty"`
}

type ProductsWithPagination = types.DataWithPagination[ProductResponse]

type ProductPrice struct {
	ProductResponse
}
