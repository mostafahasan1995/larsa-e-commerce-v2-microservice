package config_location_dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: add validator to check country id
type CreateCityDto struct {
	Name         string             `json:"name" validate:"required"`
	IsActive     *bool              `json:"isActive,omitempty"`
	CountryId    primitive.ObjectID `json:"countryId,omitempty"`
	IncludeAreas *bool              `json:"includeAreas,omitempty"`
}
