package config_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseZone struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     string `bson:"name" json:"name"`
	IsActive                 *bool  `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Code                     string `bson:"code" json:"code"`

	CountryIds []primitive.ObjectID `bson:"countryIds,omitempty" json:"countryIds,omitempty"`
	CityIds    []primitive.ObjectID `bson:"cityIds,omitempty" json:"cityIds,omitempty"`
	AreaIds    []primitive.ObjectID `bson:"areaIds,omitempty" json:"areaIds,omitempty"`

	Polygons []Polygon `bson:"polygons" json:"polygons"`

	//TODO: add vendor ids like product vendor ids.
}

type Polygon struct {
	Points []Point `bson:"points" json:"points"`
}

type Point struct {
	Longitude float64 `bson:"longitude,omitempty" json:"longitude,omitempty"`
	Latitude  float64 `bson:"latitude,omitempty" json:"latitude,omitempty"`
}

type Zone struct {
	BaseZone `bson:",inline"`
}

type ZoneResponse struct {
	BaseZone  `bson:",inline"`
	Countries []config_location_models.Country `bson:"countries,omitempty" json:"countries,omitempty"`
	Cities    []config_location_models.City    `bson:"cities,omitempty" json:"cities,omitempty"`
	Areas     []config_location_models.Area    `bson:"areas,omitempty" json:"areas,omitempty"`
}

type ZonesWithPagination = types.DataWithPagination[ZoneResponse]
