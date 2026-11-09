package order_enums

type OrderType string

const (
	DELIVERY OrderType = "delivery"
	PICKUP   OrderType = "pickup"
)

type DeliveryType string

const (
	B2B DeliveryType = "b2b"
	B2C DeliveryType = "b2c"
)

type ShipmentMethod string

const (
	CAR   ShipmentMethod = "car"
	PLAIN ShipmentMethod = "plain"
	SHIP  ShipmentMethod = "ship"
)

type PaymentStatus string

const (
	PAID           PaymentStatus = "paid"
	PARTIALLY_PAID PaymentStatus = "partially_paid"
	AGENT_PAID     PaymentStatus = "agent_paid"
	UNPAID         PaymentStatus = "unpaid"
)

type Unit string

const (
	KG    Unit = "kg"
	LB    Unit = "lb"
	CM    Unit = "cm"
	INCH  Unit = "inch"
	METER Unit = "meter"
)

type PaymentMethod string

const (
	CASH_TERMINAL    PaymentMethod = "cash_terminal"
	CASH_ON_DELIVERY PaymentMethod = "cash_on_delivery"
	ONLINE_PAYMENT   PaymentMethod = "online_payment"
)

type LogActionType string

const (
	CREATED          LogActionType = "created"
	STATUS_CHANGED   LogActionType = "status_changed"
	PAYMENT_ADDED    LogActionType = "payment_added"
	DOCUMENT_ADDED   LogActionType = "document_added"
	COMMENT_ADDED    LogActionType = "comment_added"
	SHIPMENT_UPDATED LogActionType = "shipment_updated"
	PAYMENT_UPDATED  LogActionType = "payment_updated"
)

type OrderParticipantType string

const (
	BUYER  OrderParticipantType = "buyer"
	VENDOR OrderParticipantType = "vendor"
)

type ShipmentStatus string

const (
	PENDING    ShipmentStatus = "pendin"
	PROCESSING ShipmentStatus = "processing"
	COMPLETED  ShipmentStatus = "completed"
)

type PayTo string

const (
	PAY_TO_BUYER  PayTo = "buyer"
	PAY_TO_VENDOR PayTo = "vendor"
	PAY_TO_ADMIN  PayTo = "admin"
)

type PaymentGateway string

const (
	GEIDEA PaymentGateway = "geidea"
	STRIPE PaymentGateway = "stripe"
	PAYMOB PaymentGateway = "paymob"
)
