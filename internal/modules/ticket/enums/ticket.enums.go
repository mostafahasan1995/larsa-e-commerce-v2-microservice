package ticket_enums

type TicketStatus string

const (
	CLOSED      TicketStatus = "closed"
	PENDING     TicketStatus = "pending"
	DRAFT       TicketStatus = "draft"
	IN_PROGRESS TicketStatus = "in_progress"
)

type TicketPriority string

const (
	LOW    TicketPriority = "low"
	MEDIUM TicketPriority = "medium"
	HIGH   TicketPriority = "high"
	URGENT TicketPriority = "urgent"
)
