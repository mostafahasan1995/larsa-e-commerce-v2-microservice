package settings_sales_models

import settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"

type OrderStatus struct {
	Title string `bson:"title" json:"title"`
	Color string `bson:"color" json:"color"`
	Order int    `bson:"order" json:"order"`

	IsSystemStatus   bool                                   `bson:"isSystemStatus" json:"isSystemStatus"`
	SystemStatusType *settings_sales_enums.SystemStatusType `bson:"systemStatusType" json:"systemStatusType,omitempty"`

	// AppliesTo determines which delivery types this status is valid for.
	// If empty or nil, it applies to both B2B and B2C for backward compatibility.
	AppliesTo []settings_sales_enums.FlowAppliesTo `bson:"appliesTo,omitempty" json:"appliesTo,omitempty"`

	MoveForward                  bool `bson:"moveForward" json:"moveForward"`
	MoveForwardOnlyIfPaymentDone bool `bson:"moveForwardOnlyIfPaymentDone" json:"moveForwardOnlyIfPaymentDone"`

	CancelOrder     bool `bson:"cancelOrder" json:"cancelOrder"`
	RefundOrder     bool `bson:"refundOrder" json:"refundOrder"`
	AssignToDriver  bool `bson:"assignToDriver" json:"assignToDriver"`
	MoveBackward    bool `bson:"moveBackward" json:"moveBackward"`
	ReopenOrder     bool `bson:"reopenOrder" json:"reopenOrder"`
	ConfirmDelivery bool `bson:"confirmDelivery" json:"confirmDelivery"`

	ConfirmDeliveryTemplate *settings_sales_enums.ConfirmDeliveryTemplate `bson:"confirmDeliveryTemplate" json:"confirmDeliveryTemplate"`
	IsActive                bool                                          `bson:"isActive" json:"isActive"`
}

// Logistic Settings
type LogisticSettings struct {
	EnablePickup   bool `bson:"enablePickup" json:"enablePickup"`
	EnableDelivery bool `bson:"enableDelivery" json:"enableDelivery"`

	Pickup   *LogisticPickup   `bson:"pickup" json:"pickup"`
	Delivery *LogisticDelivery `bson:"delivery" json:"delivery"`
}

type LogisticPickup struct {
	BranchName bool `bson:"branchName" json:"branchName"`
	PickUpDate bool `bson:"pickUpDate" json:"pickUpDate"`
	PickUpTime bool `bson:"pickUpTime" json:"pickUpTime"`
}

type LogisticDelivery struct {
	ShipmentMethod   bool `bson:"shipmentMethod" json:"shipmentMethod"`
	TrackingNumber   bool `bson:"trackingNumber" json:"trackingNumber"`
	TransportCompany bool `bson:"transportCompany" json:"transportCompany"`
	ShipmentDate     bool `bson:"shipmentDate" json:"shipmentDate"`
	DeliveryFee      bool `bson:"deliveryFee" json:"deliveryFee"`
	CustomAddress    bool `bson:"customAddress" json:"customAddress"`
}

// Recent Activity
type RecentActivity struct {
	EnableRecentActivity bool `bson:"enableRecentActivity" json:"enableRecentActivity"`

	NumberOfProcessedOrders  bool `bson:"numberOfProcessedOrders" json:"numberOfProcessedOrders"`
	NumberOfIncompleteOrders bool `bson:"numberOfIncompleteOrders" json:"numberOfIncompleteOrders"`
	NumberOfCustomers        bool `bson:"numberOfCustomers" json:"numberOfCustomers"`
	UnpaidAmount             bool `bson:"unpaidAmount" json:"unpaidAmount"`
	OverdueAmount            bool `bson:"overdueAmount" json:"overdueAmount"`
}

// Invoice Terms
type InvoiceTerms struct {
	DefaultInvoiceTerms     int                                  `bson:"defaultInvoiceTerms" json:"defaultInvoiceTerms"`
	DefaultInvoiceTermsUnit settings_sales_enums.InvoiceTermUnit `bson:"defaultInvoiceTermsUnit" json:"defaultInvoiceTermsUnit"`
}
