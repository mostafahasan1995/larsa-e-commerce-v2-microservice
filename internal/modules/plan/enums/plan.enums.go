package plan_enums

type PlanCustomizationType string
type PlanType string

const (
	PERSONAL PlanCustomizationType = "PERSONAL"
)

const (
	MARKETING   PlanType = "MARKETING"
	ECONOMIC    PlanType = "ECONOMIC"
	TRANSLATION PlanType = "TRANSLATION"
	INDUSTRY    PlanType = "INDUSTRY"
	COMMERCIAL  PlanType = "COMMERCIAL"
	LEGAL       PlanType = "LEGAL"
	TECHNOLOGY  PlanType = "TECHNOLOGY"
)
