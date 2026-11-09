package cart_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseCart struct {
	CartItemIds     []primitive.ObjectID `bson:"cartItemIds" json:"cartItemIds"`
	UserId          primitive.ObjectID   `bson:"userId" json:"userId"`
	IsActive        bool                 `bson:"isActive" json:"isActive"`
	IsClone         bool                 `bson:"isClone" json:"isClone"`
	RefId           primitive.ObjectID   `bson:"refId,omitempty" json:"refId,omitempty"`
	TotalQuantities uint64               `bson:"totalQuantities" json:"totalQuantities"`
	// AddressId is set when the user gets cart with ?addressId=... (delivery); used by checkout when body omits addressId.
	AddressId *primitive.ObjectID `bson:"addressId,omitempty" json:"addressId,omitempty"`
}

type Cart struct {
	models.MetaDataBaseModel `bson:",inline"`
	BaseCart                 `bson:",inline"`
}

type CartResponse struct {
	models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseCart                         `bson:",inline"`
	CartItems                        []cart_shared.CartItemResponse                `bson:"cartItems" json:"cartItems"`
	Calculation                      *cart_shared.CalculateCartItemsPricesResponse `bson:"calculation,omitempty" json:"calculation,omitempty"`
}

type CartsWithPagination = types.DataWithPagination[CartResponse]
