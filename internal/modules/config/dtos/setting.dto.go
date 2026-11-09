package config_dtos

import "go.mongodb.org/mongo-driver/bson"

type TruncateCollectionDto struct {
	CollectionName string  `json:"collectionName" validate:"required"`
	Filter         *bson.M `json:"filter" validate:"omitempty"`
}

type TruncateCollectionDtoList struct {
	Items []TruncateCollectionDto `json:"items" validate:"required,min=1,dive"`
}
