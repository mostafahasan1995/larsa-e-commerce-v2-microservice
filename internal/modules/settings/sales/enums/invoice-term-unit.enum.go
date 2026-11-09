package settings_sales_enums

type InvoiceTermUnit string

const (
	INVOICE_TERM_UNIT_DAYS   InvoiceTermUnit = "days"
	INVOICE_TERM_UNIT_WEEKS  InvoiceTermUnit = "weeks"
	INVOICE_TERM_UNIT_MONTHS InvoiceTermUnit = "months"
)

type ConfirmDeliveryTemplate string

const (
	CONFIRM_DELIVERY_TEMPLATE_EMAIL ConfirmDeliveryTemplate = "email"
	CONFIRM_DELIVERY_TEMPLATE_SMS   ConfirmDeliveryTemplate = "sms"
)

type SystemStatusType string
type TransitionStatusType string

const (
	SYSTEM_STATUS_TYPE_PENDING   SystemStatusType = "pending"
	SYSTEM_STATUS_TYPE_CANCELED  SystemStatusType = "canceled"
	SYSTEM_STATUS_TYPE_FAILED    SystemStatusType = "failed"
	SYSTEM_STATUS_TYPE_REFUND    SystemStatusType = "refund"
	SYSTEM_STATUS_TYPE_COMPLETED SystemStatusType = "completed"
	SYSTEM_STATUS_TYPE_ACCEPTED  SystemStatusType = "accepted"

	SYSTEM_STATUS_TYPE_IN_PRODUCTION     SystemStatusType = "in_production"
	SYSTEM_STATUS_TYPE_READ_FOR_SHIPMENT SystemStatusType = "read_for_shipment"
	SYSTEM_STATUS_TYPE_IN_TRANSIT        SystemStatusType = "in_transit"
	SYSTEM_STATUS_TYPE_DELIVERED         SystemStatusType = "delivered"
)

const (
	TRANSITION_STATUS_TYPE_HANDED_OVER_TO_CARRIER    TransitionStatusType = "handed_over_to_carrier"
	TRANSITION_STATUS_TYPE_DELIVERED_TO_PORT         TransitionStatusType = "delivered_to_port"
	TRANSITION_STATUS_TYPE_ARRIVED_TO_PORT           TransitionStatusType = "arrived_to_port"
	TRANSITION_STATUS_TYPE_READY_FOR_DELIVERY        TransitionStatusType = "ready_for_delivery"
	TRANSITION_STATUS_TYPE_CONFIRM_AWAITING_DELIVERY TransitionStatusType = "confirm_awaiting_delivery"
	TRANSITION_STATUS_TYPE_RELEASED_FROM_PORT        TransitionStatusType = "released_from_port"
	TRANSITION_STATUS_TYPE_CONFIRM_DELIVERY          TransitionStatusType = "confirm_delivery"
)
