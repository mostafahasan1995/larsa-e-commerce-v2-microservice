package enums

import (
	"strings"
)

type OperationType string

const (
	Get        OperationType = "Get"
	Post       OperationType = "Post"
	Put        OperationType = "Put"
	Delete     OperationType = "Delete"
	GetOne     OperationType = "GetOne"
	DeleteMany OperationType = "DeleteMany"
)

type ModuleType string

const (
	ProductModule           ModuleType = "product"
	CategoryModule          ModuleType = "category"
	CountryModule           ModuleType = "country"
	GroupModule             ModuleType = "group"
	UnitModule              ModuleType = "unit"
	BranchModule            ModuleType = "branch"
	BannerModule            ModuleType = "banner"
	AttributeModule         ModuleType = "attribute"
	ProductionLineModule    ModuleType = "productionLine"
	BrokerModule            ModuleType = "brokerConsumer"
	BrokerAgentModule       ModuleType = "BrokerAgent"
	BrokerCustomerModule    ModuleType = "brokerCustomer"
	BrokerShippingModule    ModuleType = "brokerShipping"
	ProfileModule           ModuleType = "profile"
	FeatureModule           ModuleType = "feature"
	FeatureGroupModule      ModuleType = "featureGroup"
	PlanModule              ModuleType = "plan"
	CurrencyModule          ModuleType = "currency"
	PaymentModule           ModuleType = "payment"
	InvoiceModule           ModuleType = "invoice"
	PayoutModule            ModuleType = "payout"
	TicketModule            ModuleType = "ticket"
	TicketTypeModule        ModuleType = "ticketType"
	DocumentModule          ModuleType = "document"
	QuotationModule         ModuleType = "quotation"
	ShippingModule          ModuleType = "Shipping"
	QuotationSettingsModule ModuleType = "quotationSettings"
	QuotationRequestModule  ModuleType = "quotationRequest"
	TransactionModule       ModuleType = "transaction"
	AppointmentModule       ModuleType = "appointment"
	OrderModule             ModuleType = "order"
	AddonModule             ModuleType = "addon"
	ZoneModule              ModuleType = "zone"
	AreaModule              ModuleType = "area"
	CityModule              ModuleType = "city"
	CartItemModule          ModuleType = "cartItem"
	CheckoutModule          ModuleType = "checkout"
	BranchProductModule     ModuleType = "branchProduct"
	WatchListModule         ModuleType = "watchList"
	SalesSettingsModule     ModuleType = "salesSettings"
	InventorySettingsModule ModuleType = "inventorySettings"
	GeneralSettingsModule   ModuleType = "generalSettings"
	AddressModule           ModuleType = "address"
	CouponModule            ModuleType = "coupons"
	QuestionModule          ModuleType = "questions"
	WaiterModule            ModuleType = "waiters"
	FeedbackModule          ModuleType = "feedbacks"
	RecipeModule            ModuleType = "recipe"
	DiscountModule          ModuleType = "discount"
	PromotionModule         ModuleType = "promotion"
	LoyaltyPointsModule     ModuleType = "loyaltyPoints"
	LoyaltySettingsModule   ModuleType = "loyaltySettings"
)

type Operation struct {
	name string
}

func (o Operation) String() string {
	return o.name
}

var capabilityPrefixes = map[ModuleType]string{
	ProductModule:           "ecommerceV2",
	DiscountModule:          "ecommerceV2",
	GroupModule:             "ecommerceV2",
	CategoryModule:          "ecommerceV2",
	CountryModule:           "ecommerceV2",
	UnitModule:              "ecommerceV2",
	BranchModule:            "ecommerceV2",
	BannerModule:            "ecommerceV2",
	BranchProductModule:     "ecommerceV2",
	AttributeModule:         "ecommerceV2",
	ProductionLineModule:    "ecommerceV2",
	BrokerModule:            "ecommerceV2",
	BrokerAgentModule:       "ecommerceV2",
	BrokerCustomerModule:    "ecommerceV2",
	BrokerShippingModule:    "ecommerceV2",
	ProfileModule:           "ecommerceV2",
	FeatureModule:           "ecommerceV2",
	FeatureGroupModule:      "ecommerceV2",
	PlanModule:              "ecommerceV2",
	CurrencyModule:          "ecommerceV2",
	PaymentModule:           "ecommerceV2",
	InvoiceModule:           "ecommerceV2",
	PayoutModule:            "ecommerceV2",
	TicketModule:            "ecommerceV2",
	DocumentModule:          "ecommerceV2",
	QuotationModule:         "ecommerceV2",
	QuotationSettingsModule: "ecommerceV2",
	QuotationRequestModule:  "ecommerceV2",
	TransactionModule:       "ecommerceV2",
	AppointmentModule:       "ecommerceV2",
	OrderModule:             "ecommerceV2",
	AddonModule:             "ecommerceV2",
	ZoneModule:              "ecommerceV2",
	AreaModule:              "ecommerceV2",
	CityModule:              "ecommerceV2",
	CartItemModule:          "ecommerceV2",
	CheckoutModule:          "ecommerceV2",
	WatchListModule:         "ecommerceV2",
	SalesSettingsModule:     "ecommerceV2",
	InventorySettingsModule: "ecommerceV2",
	GeneralSettingsModule:   "ecommerceV2",
	AddressModule:           "ecommerceV2",
	CouponModule:            "ecommerceV2",
	QuestionModule:          "ecommerceV2",
	WaiterModule:            "ecommerceV2",
	FeedbackModule:          "ecommerceV2",
	RecipeModule:            "ecommerceV2",
	TicketTypeModule:        "ecommerceV2",
	PromotionModule:         "ecommerceV2",
	LoyaltyPointsModule:     "ecommerceV2",
	LoyaltySettingsModule:   "ecommerceV2",
}

func generateCapability(module ModuleType, userType UserType, operation OperationType) string {
	prefix := capabilityPrefixes[module]
	if prefix != "" {
		moduleTitle := strings.Title(string(module))
		return prefix + string(userType) + moduleTitle + string(operation)
	}

	return string(module) + string(userType) + string(operation)
}

type UserTypeOps struct {
	Get        Operation
	Post       Operation
	Put        Operation
	Delete     Operation
	GetOne     Operation
	DeleteMany Operation
}

type ModuleCapability struct {
	SuperAdmin      UserTypeOps
	Vendor          UserTypeOps
	Member          UserTypeOps
	Agent           UserTypeOps
	ShippingCompany UserTypeOps
}

type CapabilityRegistry struct {
	Product           ModuleCapability
	BranchProduct     ModuleCapability
	Category          ModuleCapability
	Discount          ModuleCapability
	Country           ModuleCapability
	Currency          ModuleCapability
	Banner            ModuleCapability
	Payment           ModuleCapability
	Payout            ModuleCapability
	Branch            ModuleCapability
	Broker            ModuleCapability
	BrokerAgent       ModuleCapability
	BrokerCustomer    ModuleCapability
	BrokerShipping    ModuleCapability
	Unit              ModuleCapability
	Attribute         ModuleCapability
	ProductionLine    ModuleCapability
	Group             ModuleCapability
	Profile           ModuleCapability
	Invoice           ModuleCapability
	Feature           ModuleCapability
	FeatureGroup      ModuleCapability
	Plan              ModuleCapability
	Ticket            ModuleCapability
	TicketType        ModuleCapability
	General           string
	Quotation         ModuleCapability
	Shipping          ModuleCapability
	QuotationRequest  ModuleCapability
	Transaction       ModuleCapability
	Appointment       ModuleCapability
	Order             ModuleCapability
	Addon             ModuleCapability
	Zone              ModuleCapability
	Area              ModuleCapability
	City              ModuleCapability
	CartItem          ModuleCapability
	Checkout          ModuleCapability
	WatchList         ModuleCapability
	SalesSettings     ModuleCapability
	InventorySettings ModuleCapability
	GeneralSettings   ModuleCapability
	Address           ModuleCapability
	Coupon            ModuleCapability
	QuotationSettings ModuleCapability
	Question          ModuleCapability
	Waiter            ModuleCapability
	Feedback          ModuleCapability
	Recipe            ModuleCapability
	Promotion         ModuleCapability
	LoyaltyPoints     ModuleCapability
	LoyaltySettings   ModuleCapability
}

var CapabilitiesRegistry CapabilityRegistry

func createUserTypeOps(module ModuleType, userType UserType, operations ...OperationType) UserTypeOps {
	ops := UserTypeOps{}

	ops.Get = Operation{""}
	ops.Post = Operation{""}
	ops.Put = Operation{""}
	ops.Delete = Operation{""}
	ops.GetOne = Operation{""}
	ops.DeleteMany = Operation{""}

	for _, op := range operations {
		capabilityString := generateCapability(module, userType, op)
		switch op {
		case Get:
			ops.Get = Operation{capabilityString}
		case Post:
			ops.Post = Operation{capabilityString}
		case Put:
			ops.Put = Operation{capabilityString}
		case Delete:
			ops.Delete = Operation{capabilityString}
		case GetOne:
			ops.GetOne = Operation{capabilityString}
		case DeleteMany:
			ops.DeleteMany = Operation{capabilityString}
		}
	}

	return ops
}

func initModuleCapability(module ModuleType) ModuleCapability {
	return ModuleCapability{
		SuperAdmin:      createUserTypeOps(module, SuperAdmin, Get, Post, Put, Delete, GetOne, DeleteMany),
		Vendor:          createUserTypeOps(module, Vendor, Get, Post, Put, Delete, GetOne, DeleteMany),
		Member:          createUserTypeOps(module, Member, Get, GetOne),
		Agent:           createUserTypeOps(module, Agent, Get, Put, GetOne),
		ShippingCompany: createUserTypeOps(module, ShippingCompany, Get, Put, GetOne),
	}
}

func InitCapabilities() {
	CapabilitiesRegistry = CapabilityRegistry{
		General:           "authenticate",
		Product:           initModuleCapability(ProductModule),
		Discount:          initModuleCapability(DiscountModule),
		Category:          initModuleCapability(CategoryModule),
		Country:           initModuleCapability(CountryModule),
		Attribute:         initModuleCapability(AttributeModule),
		Group:             initModuleCapability(GroupModule),
		ProductionLine:    initModuleCapability(ProductionLineModule),
		Banner:            initModuleCapability(BannerModule),
		Unit:              initModuleCapability(UnitModule),
		Broker:            initModuleCapability(BrokerModule),
		BrokerAgent:       initModuleCapability(BrokerAgentModule),
		BrokerCustomer:    initModuleCapability(BrokerCustomerModule),
		BrokerShipping:    initModuleCapability(BrokerShippingModule),
		Payment:           initModuleCapability(PaymentModule),
		Invoice:           initModuleCapability(InvoiceModule),
		Profile:           initModuleCapability(ProfileModule),
		Feature:           initModuleCapability(FeatureModule),
		FeatureGroup:      initModuleCapability(FeatureGroupModule),
		Plan:              initModuleCapability(PlanModule),
		Currency:          initModuleCapability(CurrencyModule),
		Ticket:            initModuleCapability(TicketModule),
		TicketType:        initModuleCapability(TicketTypeModule),
		Quotation:         initModuleCapability(QuotationModule),
		QuotationSettings: initModuleCapability(QuotationSettingsModule),
		Shipping:          initModuleCapability(ShippingModule),
		QuotationRequest:  initModuleCapability(QuotationRequestModule),
		Transaction:       initModuleCapability(TransactionModule),
		Appointment:       initModuleCapability(AppointmentModule),
		Order:             initModuleCapability(OrderModule),
		Addon:             initModuleCapability(AddonModule),
		Branch:            initModuleCapability(BranchModule),
		Zone:              initModuleCapability(ZoneModule),
		Area:              initModuleCapability(AreaModule),
		City:              initModuleCapability(CityModule),
		BranchProduct:     initModuleCapability(BranchProductModule),
		CartItem:          initModuleCapability(CartItemModule),
		Checkout:          initModuleCapability(CheckoutModule),
		WatchList:         initModuleCapability(WatchListModule),
		SalesSettings:     initModuleCapability(SalesSettingsModule),
		InventorySettings: initModuleCapability(InventorySettingsModule),
		GeneralSettings:   initModuleCapability(GeneralSettingsModule),
		Address:           initModuleCapability(AddressModule),
		Coupon:            initModuleCapability(CouponModule),
		Question:          initModuleCapability(QuestionModule),
		Waiter:            initModuleCapability(WaiterModule),
		Feedback:          initModuleCapability(FeedbackModule),
		Recipe:            initModuleCapability(RecipeModule),
		Promotion:         initModuleCapability(PromotionModule),
		LoyaltyPoints:     initModuleCapability(LoyaltyPointsModule),
		LoyaltySettings:   initModuleCapability(LoyaltySettingsModule),
	}
}
