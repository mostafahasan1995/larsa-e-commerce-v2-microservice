package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddonValueDto represents the data needed to create or update an addon value
type AddonValueDto struct {
	Label        types.Localizable[string] `bson:"label" json:"label" validate:"required,localizable"`
	Image        *types.FileFieldDto       `bson:"image,omitempty" json:"image,omitempty"`
	Price        float64                   `bson:"price" json:"price" validate:"gte=0"`
	MinSelection *uint64                   `bson:"minSelection" json:"minSelection" validate:"omitempty,gte=0"`
	MaxSelection *uint64                   `bson:"maxSelection" json:"maxSelection" validate:"omitempty,gtefield=MinSelection"`
}

// UpdateAddonValueDto represents the data needed to update an addon value
type UpdateAddonValueDto struct {
	//for when the addon value is updated aleady one exits or add new addon value
	Id            *primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty" validate:"omitempty"`
	AddonValueDto `bson:",inline" json:",inline" validate:"required"`
}

// CreateAddonDto represents the data needed to create an addon
type CreateAddonDto struct {
	Label        types.Localizable[string] `bson:"label" json:"label" validate:"required,localizable"`
	MinSelection uint64                    `bson:"minSelection" json:"minSelection" validate:"gte=0"`
	MaxSelection uint64                    `bson:"maxSelection" json:"maxSelection" validate:"gtefield=MinSelection"`
	AddonValues  []AddonValueDto           `bson:"addonValues" json:"addonValues" validate:"required,min=1,dive"`
}

// UpdateAddonDto represents the data needed to update an addon
type UpdateAddonDto struct {
	Label        types.Localizable[string] `bson:"label" json:"label" validate:"required,localizable"`
	MinSelection uint64                    `bson:"minSelection" json:"minSelection" validate:"gte=0"`
	MaxSelection uint64                    `bson:"maxSelection" json:"maxSelection" validate:"gtefield=MinSelection"`
	AddonValues  []UpdateAddonValueDto     `bson:"addonValues" json:"addonValues" validate:"required,min=1,dive"`
}
