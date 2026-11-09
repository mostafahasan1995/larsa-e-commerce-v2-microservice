package quotation_enums

type ShippingType string

const (
	CUSTOMS_BROKERS ShippingType = "customs_brokers"
)

type ShippingMethod string

const (
	SEA               ShippingMethod = "sea"
	LAND              ShippingMethod = "land"
	INTERNATIONAL_AIR ShippingMethod = "international_air"
	LOCAL_AIR         ShippingMethod = "local_air"
)

type ShippingStatus string

const (
	SHIPPING_PENDING   ShippingStatus = "pending"
	SHIPPING_APPROVED  ShippingStatus = "approved"
	SHIPPING_COMPLETED ShippingStatus = "completed"
)

type PlatformCostType string

const (
	FIXED      PlatformCostType = "fixed"
	PERCENTAGE PlatformCostType = "percentage"
)

type ShippingContextKey string

const (
	SHIPPING_CONTEXT_KEY ShippingContextKey = "shipping"
)
