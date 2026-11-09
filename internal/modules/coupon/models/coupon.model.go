package coupon_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"time"
)

type Usage struct {
	Type  coupon_enums.CouponExpirationType `bson:"type" json:"type"`
	Value uint16                            `bson:"value" json:"value"`
}

type CouponUsage struct {
	Interval            *Usage                        `bson:"interval,omitempty" json:"interval,omitempty"`
	MaxUsagePerCustomer *uint16                       `bson:"maxUsagePerCustomer,omitempty" json:"maxUsagePerCustomer,omitempty"`
	MinimumSubtotal     *float64                      `bson:"minimumSubtotal,omitempty" json:"minimumSubtotal,omitempty"`
	StartDate           *time.Time                    `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate             *time.Time                    `bson:"endDate,omitempty" json:"endDate,omitempty"`
	RelatedTo           *coupon_enums.CouponRelatedTo `bson:"relatedTo,omitempty" json:"relatedTo,omitempty"`
}

type CouponRule struct {
	IsActive   bool                      `bson:"isActive" json:"isActive"`
	Title      types.Localizable[string] `bson:"title" json:"title"`
	StartDate  time.Time                 `bson:"startDate" json:"startDate"`
	EndDate    time.Time                 `bson:"endDate" json:"endDate"`
	Priority   uint16                    `bson:"priority" json:"priority"`
	Conditions CouponRuleConditions      `bson:"conditions" json:"conditions"`
}

type CouponRuleCondition struct {
	Type     coupon_enums.CouponRuleConditionType `bson:"type" json:"type"`
	Operator enums.Operator                       `bson:"operator" json:"operator"`
	Value    interface{}                          `bson:"value" json:"value"`
}

type CouponRuleConditions struct {
	AllowedIf coupon_enums.CouponRuleConditionAllowedIf `bson:"allowedIf" json:"allowedIf"`
	When      bool                                      `bson:"when" json:"when"`
	Values    []CouponRuleCondition                     `bson:"values" json:"values"`
}

type BaseCoupon struct {
	Title                types.Localizable[string]    `bson:"title" json:"title"`
	IsActive             *bool                        `bson:"isActive,omitempty" json:"isActive,omitempty"`
	Code                 string                       `bson:"code" json:"code"`
	Description          *types.Localizable[string]   `bson:"description,omitempty" json:"description,omitempty"`
	Type                 coupon_enums.CouponType      `bson:"type" json:"type"`
	OrderType            *order_enums.OrderType       `bson:"orderType" json:"orderType"`
	AppliedOn            coupon_enums.CouponAppliedOn `bson:"appliedOn" json:"appliedOn"`
	Expiration           *Usage                       `bson:"expiration,omitempty" json:"expiration,omitempty"`
	Discount             product_models.Discount      `bson:"discount" json:"discount"`
	TotalUsages          *uint64                      `bson:"totalUsages,omitempty" json:"totalUsages,omitempty"`
	Usage                *CouponUsage                 `bson:"usage,omitempty" json:"usage,omitempty"`
	AppliedAutomatically bool                         `bson:"appliedAutomatically" json:"appliedAutomatically"`
	Rule                 *CouponRule                  `bson:"rule,omitempty" json:"rule,omitempty"`
}

type Coupon struct {
	BaseCoupon               `bson:",inline"`
	models.MetaDataBaseModel `bson:",inline"`
}

type CouponResponse struct {
	BaseCoupon                       `bson:",inline"`
	models.OptionalMetaDataBaseModel `bson:",inline"`
	UserCoupon                       *UserCouponResponse `bson:"userCoupon,omitempty" json:"userCoupon,omitempty"`
	UserCouponTotalUsages            *uint64             `bson:"userCouponTotalUsages,omitempty" json:"userCouponTotalUsages,omitempty"`
	TotalOrders                      *uint64             `bson:"totalOrders,omitempty" json:"totalOrders,omitempty"`
	NumberOfOrders                   *uint64             `bson:"numberOfOrders,omitempty" json:"numberOfOrders,omitempty"`
}

type CouponsWithPagination = types.DataWithPagination[CouponResponse]
