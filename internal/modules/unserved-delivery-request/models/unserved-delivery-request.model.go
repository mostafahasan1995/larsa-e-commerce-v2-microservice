package unserved_delivery_request_models

import (
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnservedDeliveryRequest struct {
	base_models.MetaDataBaseModel `bson:",inline"`

	Latitude  float64             `bson:"latitude" json:"latitude"`
	Longitude float64             `bson:"longitude" json:"longitude"`
	UserId    *primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	Country   *string             `bson:"country,omitempty" json:"country,omitempty"`
	City      *string             `bson:"city,omitempty" json:"city,omitempty"`
	Area      *string             `bson:"area,omitempty" json:"area,omitempty"`
	Source    *string             `bson:"source,omitempty" json:"source,omitempty"`
}

type UnservedDeliveryRequestResponse struct {
	UnservedDeliveryRequest `bson:",inline" json:",inline"`
}
