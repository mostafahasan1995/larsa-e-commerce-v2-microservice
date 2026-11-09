package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateProductionLineDto represents the data needed to create a production line
type CreateProductionLineDto struct {
	Name        types.Localizable[string] `bson:"name" json:"name" validate:"required,localizable"`
	Description types.Localizable[string] `bson:"description" json:"description" validate:"required,localizable"`
	CategoryIds []primitive.ObjectID      `bson:"categoryIds,omitempty" json:"categoryIds,omitempty" validate:"omitempty,min=1,dive,required"`
	Attributes  []AttributeDto            `json:"attributes,omitempty" validate:"omitempty,min=1,dive,required"`
	VendorId    *primitive.ObjectID       `bson:"vendorId,omitempty" json:"vendorId,omitempty" validate:"omitempty"`
}
