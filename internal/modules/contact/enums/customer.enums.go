package contact_enums

type CustomerStatus string

type NotificationMethod string

const (
	EMAIL     NotificationMethod = "email"
	SMS       NotificationMethod = "sms"
	WHATSAPP  NotificationMethod = "whatsapp"
	DASHBOARD NotificationMethod = "dashboard"
)
