package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BaseGroup struct {
	Name                          types.Localizable[string] `bson:"name" json:"name"`
	Code                          string                    `bson:"code" json:"code"`
	base_models.MetaDataBaseModel `bson:",inline"`
}

type Group struct {
	BaseGroup `bson:",inline"`
}

type GroupResponse struct {
	BaseGroup `bson:",inline"`
}

type GroupsWithPagination = types.DataWithPagination[GroupResponse]
