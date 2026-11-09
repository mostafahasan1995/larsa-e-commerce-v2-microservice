package config_location_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

// TODO: add geo indexing for multiPolygon
type BaseCountry struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     string `bson:"name" json:"name"`
	IsActive                 *bool  `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Iso3                     string `bson:"iso3" json:"iso3"`

	MultiPolygon [][][]float64 `bson:"multiPolygon" json:"multiPolygon"`
}

type Country struct {
	BaseCountry `bson:",inline"`
}
type CountryResponse struct {
	Country `bson:",inline"`
}

type CountriesWithPagination = types.DataWithPagination[CountryResponse]

type CityWithAreasResponse struct {
	City  `bson:",inline"`
	Areas []Area `bson:"areas" json:"areas"`
}

type CountryWithLocationsResponse struct {
	Country `bson:",inline"`
	Cities  []CityWithAreasResponse `bson:"cities" json:"cities"`
}
