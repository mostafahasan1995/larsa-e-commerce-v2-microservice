package product_enums

type PricePlatform string

const (
	WEB    PricePlatform = "web"
	MOBILE PricePlatform = "mobile"
	POS    PricePlatform = "pos"
)

type CustomerGroup string

const (
	LOYALTY CustomerGroup = "loyalty"
	ENGAGED CustomerGroup = "engaged"
)

type DiscountType string

const (
	PERCENTAGE DiscountType = "percentage"
	FIXED      DiscountType = "fixed"
)
