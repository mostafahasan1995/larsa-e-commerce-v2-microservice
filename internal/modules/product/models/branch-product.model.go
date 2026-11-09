package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseBranchProduct struct {
	OverridePrice       *PriceItem                 `bson:"overridePrice" json:"overridePrice"`
	OverrideDescription *types.Localizable[string] `bson:"overrideDescription" json:"overrideDescription"`
	IsActive            bool                       `bson:"isActive" json:"isActive"`
	OutOfStock          bool                       `bson:"outOfStock" json:"outOfStock"`
	AutoResetOutOfStock bool                       `bson:"autoResetOutOfStock" json:"autoResetOutOfStock"`
}

type BranchProduct struct {
	BranchId                 primitive.ObjectID `bson:"branchId" json:"branchId"`
	ProductId                primitive.ObjectID `bson:"productId" json:"productId"`
	BaseBranchProduct        `bson:",inline"`
	models.MetaDataBaseModel `bson:",inline"`
}

type BranchProductResponse struct {
	BranchId                         *primitive.ObjectID `bson:"branchId,omitempty" json:"branchId,omitempty"`
	ProductId                        *primitive.ObjectID `bson:"productId,omitempty" json:"productId,omitempty"`
	BaseBranchProduct                `bson:",inline"`
	models.OptionalMetaDataBaseModel `bson:",inline"`
}
