package payment_enums

type InvoiceStatus string
type InvoicePaymentStatus string

const (
	INVOICE_DRAFT     InvoiceStatus = "draft"
	INVOICE_PENDING   InvoiceStatus = "pending"
	INVOICE_APPROVED  InvoiceStatus = "approved"
	INVOICE_REJECTED  InvoiceStatus = "rejected"
	INVOICE_COMPLETED InvoiceStatus = "completed"
)

const (
	PAYMENT_PENDING   InvoicePaymentStatus = "pending"
	PAYMENT_PAID      InvoicePaymentStatus = "paid"
	PAYMENT_PARTIALLY InvoicePaymentStatus = "partially"
	PAYMENT_OVERDUE   InvoicePaymentStatus = "overdue"
)
