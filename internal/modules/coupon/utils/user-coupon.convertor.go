package coupon_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConvertToUserCouponModelArgs struct {
	UserId     primitive.ObjectID `json:"userId"`
	CouponCode string             `json:"couponCode"`
}

func ConvertToUserCouponModel(createUserCouponDto ConvertToUserCouponModelArgs, additionalOptions types.AdditionalOptions) coupon_models.UserCoupon {
	coupon := coupon_models.UserCoupon{
		BaseUserCoupon: coupon_models.BaseUserCoupon{
			UserId:     createUserCouponDto.UserId,
			CouponCode: createUserCouponDto.CouponCode,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}

	return coupon

}
