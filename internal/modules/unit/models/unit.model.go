package unit_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	unit_enums "larsa-e-commerce-microservice-v2/internal/modules/unit/enums"
)

type BaseUnit struct {
	Name     types.Localizable[string] `bson:"name" json:"name"`
	Type     unit_enums.UnitType       `bson:"type" json:"type"`
	Category unit_enums.UnitCategory   `bson:"category" json:"category"`
	Rate     *float32                  `bson:"rate,omitempty" json:"rate,omitempty"`
	Active   bool                      `bson:"active" json:"active"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Unit struct {
	BaseUnit `bson:",inline"`
}

type UnitResponse struct {
	BaseUnit `bson:",inline"`
}

type UnitsWithPagination = types.DataWithPagination[UnitResponse]

//
