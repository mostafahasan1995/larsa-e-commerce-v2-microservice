package feature_dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateFeatureGroupDTO struct {
	Name     string         `json:"name" validate:"required"`
	Features []FeatureGroup `json:"features"`
}

type UpdateFeatureGroupDTO struct {
	Name     string         `json:"name" validate:"required"`
	Features []FeatureGroup `json:"features"`
}

type FeatureGroup struct {
	Id     primitive.ObjectID `json:"id"`
	Active bool               `json:"active"`
}
