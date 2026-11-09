package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseLinkBranchProductDto struct {
	OverridePrice       *PriceItemDto              `json:"overridePrice" validate:"omitempty"`
	OverrideDescription *types.Localizable[string] `json:"overrideDescription" validate:"omitempty,localizable"`
	IsActive            bool                       `json:"isActive"`
	OutOfStock          bool                       `json:"outOfStock"`
}

type LinkProductToBranchDto struct {
	BaseLinkBranchProductDto
	ProductId *primitive.ObjectID `json:"-"`
	BranchId  primitive.ObjectID  `json:"branchId" validate:"required"`
}

type LinkBranchToProductDto struct {
	BaseLinkBranchProductDto
	ProductId primitive.ObjectID  `json:"productId" validate:"required"`
	BranchId  *primitive.ObjectID `json:"-"`
}

type UpsertProductBranchDto struct {
	BaseLinkBranchProductDto
	ProductId primitive.ObjectID `json:"productId" validate:"required"`
	BranchId  primitive.ObjectID `json:"branchId" validate:"required"`
}
