package config_location_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseArea struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     string `bson:"name" json:"name"`
	IsActive                 *bool  `bson:"isActive,omitempty" json:"isActive,omitempty"`

	CityId       primitive.ObjectID `bson:"cityId,omitempty" json:"cityId,omitempty"`
	CityName     *string            `bson:"cityName,omitempty" json:"cityName,omitempty"`
	MultiPolygon [][][]float64      `bson:"multiPolygon" json:"multiPolygon"`
}

type Area struct {
	BaseArea `bson:",inline"`
}

type AreaResponse struct {
	Area `bson:",inline"`
	City *City `bson:"city,omitempty" json:"city,omitempty"`
}

type AreasWithPagination = types.DataWithPagination[AreaResponse]
