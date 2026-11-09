package address_enums

type AddressType string

const (
	WarehouseAddress AddressType = "warehouse"
	BuildingAddress  AddressType = "building"
	OfficeAddress    AddressType = "office"
)
