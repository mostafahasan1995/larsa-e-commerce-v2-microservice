package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddonValue represents a selectable option for an addon
type AddonValue struct {
	Id           primitive.ObjectID        `bson:"_id" json:"id"`
	Label        types.Localizable[string] `bson:"label" json:"label"`
	Image        *types.FileFieldDto       `bson:"image,omitempty" json:"image,omitempty"`
	Price        float64                   `bson:"price" json:"price"`
	MinSelection *uint64                   `bson:"minSelection" json:"minSelection"`
	MaxSelection *uint64                   `bson:"maxSelection" json:"maxSelection"`
}

// BaseAddon contains the common fields for Addon
type BaseAddon struct {
	Label        types.Localizable[string] `bson:"label" json:"label"`
	MinSelection uint64                    `bson:"minSelection" json:"minSelection"`
	MaxSelection uint64                    `bson:"maxSelection" json:"maxSelection"`
	AddonValues  []AddonValue              `bson:"addonValues" json:"addonValues"`

	base_models.MetaDataBaseModel `bson:",inline"`
}

// Addon represents an addon that can be attached to a product
type Addon struct {
	BaseAddon `bson:",inline"`
}

// AddonResponse is the response model for Addon
type AddonResponse struct {
	BaseAddon `bson:",inline"`
}

// AddonsWithPagination represents a paginated list of addons
type AddonsWithPagination = types.DataWithPagination[AddonResponse]
