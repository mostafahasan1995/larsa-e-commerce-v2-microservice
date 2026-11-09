package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseProductionLine struct {
	Name        types.Localizable[string] `bson:"name,omitempty" json:"name,omitempty"`
	Description types.Localizable[string] `bson:"description,omitempty" json:"description,omitempty"`
	Attributes  []ProductAttribute        `bson:"attributes,omitempty" json:"attributes,omitempty"`
	CategoryIds []primitive.ObjectID      `bson:"categoryIds,omitempty" json:"categoryIds,omitempty"`
	ProductIds  []primitive.ObjectID      `bson:"productIds,omitempty" json:"productIds,omitempty"`

	VendorId                      *primitive.ObjectID `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	base_models.MetaDataBaseModel `bson:",inline"`
}

// ProductionLine represents a production line in the system
type ProductionLine struct {
	BaseProductionLine `bson:",inline"`
}

type ProductionLineResponse struct {
	BaseProductionLine `bson:",inline"`
	Products           []ProductResponse                `bson:"products,omitempty" json:"products,omitempty"`
	Categories         []config_models.CategoryResponse `bson:"categories,omitempty" json:"categories,omitempty"`
}

type ProductionLinesWithPagination = types.DataWithPagination[ProductionLineResponse]
