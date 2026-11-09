package promotion_dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePromotionDTO struct {
	Name               string               `json:"name" validate:"required"`
	Description        string               `json:"description"`
	DiscountPercentage float64              `json:"discountPercentage" validate:"required,gt=0,lte=100"`
	StartDate          time.Time            `json:"startDate" validate:"required"`
	EndDate            time.Time            `json:"endDate" validate:"required,gtfield=StartDate"`
	ProductIds         []primitive.ObjectID `json:"productIds"`
	AllProducts        bool                 `json:"allProducts"`
	BranchIds          []primitive.ObjectID `json:"branchIds"`
	IsActive           bool                 `json:"isActive"`
}
