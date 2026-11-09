package payment_enums

type TransactionType string
type TransactionStatus string
type TransactionPriority string

const (
	COMPLAINTS TransactionType = "complaints"
	REQUEST    TransactionType = "request"
	INQUIRY    TransactionType = "inquiry"
)

const (
	CLOSED              TransactionStatus = "closed"
	TRANSACTION_PENDING TransactionStatus = "pending"
	DRAFT               TransactionStatus = "draft"
	IN_PROGRESS         TransactionStatus = "in_progress"
	OPEN                TransactionStatus = "open"
)

const (
	HIGH   TransactionPriority = "high"
	MEDIUM TransactionPriority = "medium"
	LOW    TransactionPriority = "low"
)
