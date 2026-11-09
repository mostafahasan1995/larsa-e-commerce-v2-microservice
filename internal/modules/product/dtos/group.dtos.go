package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type CreateGroupDto struct {
	Name                          types.Localizable[string] `json:"name" validate:"required,localizable"`
	Code                          *string                   `json:"-"`
	base_models.MetaDataBaseModel `bson:",inline"`
}
