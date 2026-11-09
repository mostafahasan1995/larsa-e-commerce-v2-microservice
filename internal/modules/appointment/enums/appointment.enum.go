package appointment_enums

type AppointmentStatus string
type AppointmentLocation string
type TranslateType string

const (
	PENDING   AppointmentStatus = "PENDING"
	APPROVED  AppointmentStatus = "APPROVED"
	REJECTED  AppointmentStatus = "REJECTED"
	COMPLETED AppointmentStatus = "COMPLETED"
	CANCELED  AppointmentStatus = "CANCELED"
)

const (
	ONLINE  AppointmentLocation = "ONLINE"
	OFFLINE AppointmentLocation = "OFFLINE"
)

const (
	AI    TranslateType = "AI"
	HUMAN TranslateType = "HUMAN"
)
