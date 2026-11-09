package config_dtos

import (
	config_enums "larsa-e-commerce-microservice-v2/internal/modules/config/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateZoneDto struct {
	Name     string `json:"name" validate:"required"`
	IsActive *bool  `json:"isActive,omitempty"`
	Code     string `json:"code" validate:"required"`

	CountryIds []primitive.ObjectID `json:"countryIds,omitempty"`
	CityIds    []primitive.ObjectID `json:"cityIds,omitempty"`
	AreaIds    []primitive.ObjectID `json:"areaIds,omitempty"`

	Polygons []PolygonDto `json:"polygons,omitempty"`
}

type PolygonDto struct {
	Points []PointDto `json:"points,omitempty"`
}

type PointDto struct {
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

type FindZonesContainingLocationDto struct {
	LocationType config_enums.LocationType `json:"locationType" validate:"required"`
	LocationName string                    `json:"locationName" validate:"required"`
}

type FindZonesContainingPointDto struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}
