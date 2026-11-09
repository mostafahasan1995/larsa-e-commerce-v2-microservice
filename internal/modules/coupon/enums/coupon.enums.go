package coupon_enums

type CouponAppliedOn string

const (
	GENERAL CouponAppliedOn = "general"
	PRODUCT CouponAppliedOn = "product"
)

type CouponType string

const (
	AUTO_GENERATED     CouponType = "auto_generated"
	MANUAL             CouponType = "manual"
	FEEDBACK_GENERATED CouponType = "feedback_generated"
)

type CouponExpirationType string

const (
	DAYS   CouponExpirationType = "days"
	WEEKS  CouponExpirationType = "weeks"
	MONTHS CouponExpirationType = "months"
	YEARS  CouponExpirationType = "years"
)

type CouponRuleConditionType string

const (
	NUMBER_OF_ORDERS     CouponRuleConditionType = "number_of_orders"
	ORDER_PAYMENT_STATUS CouponRuleConditionType = "order_payment_status"
	DELIVERY_METHOD      CouponRuleConditionType = "delivery_method"
	BRANCH               CouponRuleConditionType = "branch"
	SUB_TOTAL_AMOUNT     CouponRuleConditionType = "sub_total_amount"
	ORDER_TIMEFRAME      CouponRuleConditionType = "order_timeframe"
)

type CouponRuleConditionAllowedIf string

const (
	ALL  CouponRuleConditionAllowedIf = "all"
	ANY  CouponRuleConditionAllowedIf = "any"
	NONE CouponRuleConditionAllowedIf = "none"
)

type CouponRelatedTo string

const (
	FEEDBACK CouponRelatedTo = "feedback"
)
