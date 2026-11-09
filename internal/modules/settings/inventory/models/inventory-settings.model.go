package settings_inventory_models

import "larsa-e-commerce-microservice-v2/internal/modules/base/models"

type BaseInventorySettings struct {
	models.MetaDataBaseModel `bson:",inline"`

	UnitsActive bool `bson:"unitsActive" json:"unitsActive"`
}

type InventorySettings struct {
	BaseInventorySettings `bson:",inline"`
}

type InventorySettingsResponse struct {
	BaseInventorySettings `bson:",inline"`
}
