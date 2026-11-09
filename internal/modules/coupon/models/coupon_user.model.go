package coupon_models

import (
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCouponUsage struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type BaseUserCoupon struct {
	UserId      primitive.ObjectID `bson:"userId" json:"userId"`
	CouponCode  string             `bson:"couponCode" json:"couponCode"`
	Usages      *[]UserCouponUsage `bson:"usages,omitempty" json:"usages,omitempty"`
	TotalUsages uint64             `bson:"totalUsages" json:"totalUsages"`
}

type UserCoupon struct {
	models.MetaDataBaseModel `bson:",inline"`
	BaseUserCoupon           `bson:",inline"`
}

type UserCouponResponse struct {
	models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseUserCoupon                   `bson:",inline"`
}
