package loyalty_dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLoyaltyPointsDTO struct {
	UserId primitive.ObjectID `json:"userId" validate:"required"`
}
