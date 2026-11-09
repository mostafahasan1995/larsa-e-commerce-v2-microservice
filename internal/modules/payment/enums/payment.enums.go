package payment_enums

type PaymentMethod string
type PaymentStatus string

const (
	BANK_TRANSFER PaymentMethod = "bank_transfer"
)

const (
	PENDING   PaymentStatus = "pending"
	APPROVED  PaymentStatus = "approved"
	REJECTED  PaymentStatus = "rejected"
	COMPLETED PaymentStatus = "completed"
)
