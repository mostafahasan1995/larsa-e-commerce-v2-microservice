package config_location_dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: add validator to check city id
type CreateAreaDto struct {
	Name     string             `json:"name" validate:"required"`
	IsActive *bool              `json:"isActive,omitempty"`
	CityId   primitive.ObjectID `json:"cityId" validate:"required"`
}
