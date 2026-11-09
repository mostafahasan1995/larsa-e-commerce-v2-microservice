package cart_shared

import (
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculateCartItemsPricesResponse struct {
	//Discount value
	//Coupon Value
	//Loyalty Value

	DiscountValue float64 `json:"discountValue"`
	CouponValue   float64 `json:"couponValue"`
	LoyaltyValue  float64 `json:"loyaltyValue"`

	//SubTotalPrice
	//Total discount value
	//Total Price
	SubTotalPrice      float64 `json:"subTotalPrice"`
	TotalDiscountValue float64 `json:"totalDiscountValue"`
	TotalPrice         float64 `json:"totalPrice"`

	DeliveryFee               float64               `json:"deliveryFee"`
	LoyaltyDiscountPercentage float64               `json:"loyaltyDiscountPercentage"`
	LoyaltyOrdersRemaining    int64                 `json:"loyaltyOrdersRemaining"`
	EstimatedTime             uint64                `json:"estimatedTime"`
	CartId                    primitive.ObjectID    `json:"cartId"`
	ZoneId                    *primitive.ObjectID   `json:"zoneId,omitempty"`
	BranchId                  *primitive.ObjectID   `json:"branchId,omitempty"`
	CurrencyCode              *string               `json:"currencyCode,omitempty"`
	AddressId                 *primitive.ObjectID   `json:"addressId,omitempty"`
	CouponCode                *string               `json:"couponCode,omitempty"`
	OrderType                 order_enums.OrderType `json:"orderType"`
	// CartItems                 []CartItemResponse             `json:"cartItems"`
	OrderItems []order_dto.CreateOrderItemDto `json:"-"`
	Note       *string                        `json:"note,omitempty"`
}

type GetMineArgs struct {
	Query         string
	CurrencyCode  string
	BranchId      string
	AddressId     string
	CouponCode    string
	CalculateList bool
	OrderType     order_enums.OrderType
	ZoneId        string
}
