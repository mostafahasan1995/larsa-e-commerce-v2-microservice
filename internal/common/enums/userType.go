package enums

type UserType string

const (
	SuperAdmin      UserType = "SuperAdmin"
	Vendor          UserType = "Vendor"
	Member          UserType = "Member"
	Buyer           UserType = "Buyer"
	Normal          UserType = "Normal"
	Guest           UserType = "Guest"
	Agent           UserType = "Agent"
	ShippingCompany UserType = "ShippingCompany"
)

type RoleName string

const (
	VendorRoleName          RoleName = "vendor"
	AgentRoleName           RoleName = "agent"
	ShippingCompanyRoleName RoleName = "shippingCompany"
)
