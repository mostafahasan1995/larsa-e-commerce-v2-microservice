package config_dtos

type DataMigrationDto struct {
	FromDatabase   string `json:"fromDatabase" validate:"required"`
	ToDatabase     string `json:"toDatabase" validate:"required"`
	CollectionName string `json:"collectionName" validate:"required"`
}
