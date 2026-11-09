package promotion_dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdatePromotionDTO struct {
	Name               *string               `json:"name,omitempty"`
	Description        *string               `json:"description,omitempty"`
	DiscountPercentage *float64              `json:"discountPercentage,omitempty" validate:"omitempty,gt=0,lte=100"`
	StartDate          *time.Time            `json:"startDate,omitempty"`
	EndDate            *time.Time            `json:"endDate,omitempty" validate:"omitempty,gtfield=StartDate"`
	ProductIds         *[]primitive.ObjectID `json:"productIds,omitempty"`
	AllProducts        *bool                 `json:"allProducts,omitempty"`
	BranchIds          *[]primitive.ObjectID `json:"branchIds,omitempty"`
	IsActive           *bool                 `json:"isActive,omitempty"`
}
