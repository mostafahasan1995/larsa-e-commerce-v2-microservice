package config_models

import (
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BaseSetting struct {
	models.MetaDataBaseModel `bson:",inline"`
	Collection               string `bson:"collection" json:"collection"`
	CurrentSerialNumber      uint64 `bson:"currentSerialNumber" json:"currentSerialNumber"`
}

type Setting struct {
	BaseSetting `bson:",inline"`
}

type SettingResponse struct {
	BaseSetting `bson:",inline"`
}
