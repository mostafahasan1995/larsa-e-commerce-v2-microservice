package loyalty_dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type AssignOrderForPointsDTO struct {
	OrderId primitive.ObjectID `json:"orderId" validate:"required"`
}

type ProcessDiscountUsageDTO struct {
	OrderId        primitive.ObjectID `json:"orderId" validate:"required"`
	DiscountAmount float64            `json:"discountAmount" validate:"required,min=0"`
}
