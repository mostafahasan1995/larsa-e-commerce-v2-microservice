package profile_enums

// ProfileType defines the type of profile (Buyer or Vendor)
type ProfileType string

const (
	Vendor          ProfileType = "Vendor"
	SuperAdmin      ProfileType = "SuperAdmin"
	Agent           ProfileType = "Agent"
	ShippingCompany ProfileType = "ShippingCompany"
	Buyer           ProfileType = "Buyer"
	Guest           ProfileType = "Guest"
)
