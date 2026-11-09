package cart_shared

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	cart_enums "larsa-e-commerce-microservice-v2/internal/modules/carts/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartAddon struct {
	Id     primitive.ObjectID `bson:"id" json:"id"`
	Values []CartAddonValue   `bson:"values" json:"values"`
}

type CartAddonValue struct {
	Id       primitive.ObjectID `bson:"id" json:"id"`
	Quantity uint64             `bson:"quantity" json:"quantity"`
}

type AdditionalInformation struct {
	Image *types.FileFieldDto `bson:"image,omitempty" json:"image,omitempty"`
	Logo  *types.FileFieldDto `bson:"logo,omitempty" json:"logo,omitempty"`
	Note  *string             `bson:"note,omitempty" json:"note,omitempty"`
}

type BaseCartItem struct {
	ItemId                primitive.ObjectID      `bson:"itemId" json:"itemId"`
	Type                  cart_enums.CartItemType `bson:"type" json:"type"`
	Quantity              uint64                  `bson:"quantity" json:"quantity"`
	UserId                primitive.ObjectID      `bson:"userId" json:"userId"`
	CartId                primitive.ObjectID      `bson:"cartId" json:"cartId"`
	UnitId                *primitive.ObjectID     `bson:"unitId" json:"unitId"`
	VendorId              *primitive.ObjectID     `bson:"vendorId" json:"vendorId"`
	BranchId              *primitive.ObjectID     `bson:"branchId" json:"branchId"`
	IsActive              bool                    `bson:"isActive" json:"isActive"`
	Addons                *[]CartAddon            `bson:"addons,omitempty" json:"addons,omitempty"`
	AdditionalInformation *AdditionalInformation  `bson:"additionalInformation,omitempty" json:"additionalInformation,omitempty"`
}

type CartItem struct {
	models.MetaDataBaseModel `bson:",inline"`
	BaseCartItem             `bson:",inline"`
}

type CartItemResponse struct {
	models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseCartItem                     `bson:",inline"`
	Product                          *product_models.ProductResponse `bson:"product,omitempty" json:"product,omitempty"`
	AddonList                        []product_models.Addon          `bson:"addonList" json:"addonList"`
}

type CartItemsWithPagination = types.DataWithPagination[CartItemResponse]
