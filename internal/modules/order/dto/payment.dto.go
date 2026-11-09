package order_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// PaymentInitiationResponseDTO defines the response structure for payment initiation
type PaymentInitiationResponseDTO struct {
	OrderId     primitive.ObjectID `json:"orderId"`
	IsRedirect  bool               `json:"isRedirect"`
	RedirectURL *string            `json:"redirectUrl,omitempty"`
}
