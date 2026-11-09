package coupon_utils

import (
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserCouponsFromOrdersArgs struct {
	Coupon        coupon_models.CouponResponse
	OrderPipeline []bson.M
}

// BuyerOrderCount represents the result of grouping orders by buyer
type BuyerOrderCount struct {
	BuyerID primitive.ObjectID `bson:"_id" json:"buyerId"`
	Count   int                `bson:"count" json:"count"`
}

type NumberOfOrdersConditionPipelineResult struct {
	UserId      string `bson:"userId" json:"userId"`
	TotalOrders uint64 `bson:"totalOrders" json:"totalOrders"`
}
