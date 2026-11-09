package promotion_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasePromotion struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     string               `bson:"name" json:"name"`
	Description              string               `bson:"description" json:"description"`
	DiscountPercentage       float64              `bson:"discountPercentage" json:"discountPercentage"`
	StartDate                time.Time            `bson:"startDate" json:"startDate"`
	EndDate                  time.Time            `bson:"endDate" json:"endDate"`
	ProductIds               []primitive.ObjectID `bson:"productIds" json:"productIds"`
	AllProducts              bool                 `bson:"allProducts" json:"allProducts"`
	BranchIds                []primitive.ObjectID `bson:"branchIds" json:"branchIds"`
	IsActive                 bool                 `bson:"isActive" json:"isActive"`
}

type Promotion struct {
	BasePromotion `bson:",inline"`
}

type PromotionResponse struct {
	BasePromotion `bson:",inline"`
	Branches      []*product_models.BranchResponse  `bson:"branches" json:"branches"`
	Products      []*product_models.ProductResponse `bson:"products" json:"products"`
}

type PromotionsWithPagination = types.DataWithPagination[PromotionResponse]
