package settings_sales_dtos

import settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"

type UpdateSalesSettingsDTO struct {
	OrderStatuses []OrderStatusDTO `json:"orderStatuses"`

	ShowVendorInfoInOrderDetails *bool                `json:"showVendorInfoInOrderDetails,omitempty"`
	LogisticSettings             *LogisticSettingsDTO `json:"logisticSettings,omitempty"`
	RecentActivity               *RecentActivityDTO   `json:"recentActivity,omitempty"`
	InvoiceTerms                 *InvoiceTermsDTO     `json:"invoiceTerms,omitempty"`

	EnableDiscountOnSales           *bool    `json:"enableDiscountOnSales,omitempty"`
	EnablePartialPayment            *bool    `json:"enablePartialPayment,omitempty"`
	SalesTaxesRate                  *float64 `json:"salesTaxesRate,omitempty"`
	EnableSalesNotes                *bool    `json:"enableSalesNotes,omitempty"`
	EnableAutoGenerateInvoiceNumber *bool    `json:"enableAutoGenerateInvoiceNumber,omitempty"`
	EnableInvoiceEmailNotification  *bool    `json:"enableInvoiceEmailNotification,omitempty"`

	// Platform order flow mode: b2c or b2b
	OrderFlowMode *settings_sales_enums.OrderFlowMode `json:"orderFlowMode,omitempty"`
}

type OrderStatusDTO struct {
	Title string `json:"title" validate:"required"`
	Color string `json:"color" validate:"required"`
	Order int    `json:"order"`

	IsSystemStatus   bool                                   `json:"isSystemStatus"`
	SystemStatusType *settings_sales_enums.SystemStatusType `json:"systemStatusType,omitempty"`

	// Which delivery types this status applies to
	AppliesTo []settings_sales_enums.FlowAppliesTo `json:"appliesTo,omitempty"`

	MoveForward                  bool `json:"moveForward"`
	MoveForwardOnlyIfPaymentDone bool `json:"moveForwardOnlyIfPaymentDone"`

	CancelOrder     bool `json:"cancelOrder"`
	RefundOrder     bool `json:"refundOrder"`
	AssignToDriver  bool `json:"assignToDriver"`
	MoveBackward    bool `json:"moveBackward"`
	ReopenOrder     bool `json:"reopenOrder"`
	ConfirmDelivery bool `json:"confirmDelivery"`

	ConfirmDeliveryTemplate *settings_sales_enums.ConfirmDeliveryTemplate `json:"confirmDeliveryTemplate,omitempty"`
	IsActive                bool                                          `json:"isActive"`
}

// Logistic Settings
type LogisticSettingsDTO struct {
	EnablePickup   *bool                `json:"enablePickup,omitempty"`
	EnableDelivery *bool                `json:"enableDelivery,omitempty"`
	Pickup         *LogisticPickupDTO   `json:"pickup,omitempty"`
	Delivery       *LogisticDeliveryDTO `json:"delivery,omitempty"`
}

type LogisticPickupDTO struct {
	BranchName *bool `json:"branchName,omitempty"`
	PickUpDate *bool `json:"pickUpDate,omitempty"`
	PickUpTime *bool `json:"pickUpTime,omitempty"`
}

type LogisticDeliveryDTO struct {
	ShipmentMethod   *bool `json:"shipmentMethod,omitempty"`
	TrackingNumber   *bool `json:"trackingNumber,omitempty"`
	TransportCompany *bool `json:"transportCompany,omitempty"`
	ShipmentDate     *bool `json:"shipmentDate,omitempty"`
	DeliveryFee      *bool `json:"deliveryFee,omitempty"`
	CustomAddress    *bool `json:"customAddress,omitempty"`
}

// Recent Activity
type RecentActivityDTO struct {
	EnableRecentActivity *bool `json:"enableRecentActivity,omitempty"`

	NumberOfProcessedOrders  *bool `json:"numberOfProcessedOrders,omitempty"`
	NumberOfIncompleteOrders *bool `json:"numberOfIncompleteOrders,omitempty"`
	NumberOfCustomers        *bool `json:"numberOfCustomers,omitempty"`
	UnpaidAmount             *bool `json:"unpaidAmount,omitempty"`
	OverdueAmount            *bool `json:"overdueAmount,omitempty"`
}

// Invoice Terms
type InvoiceTermsDTO struct {
	DefaultInvoiceTerms     *int                                  `json:"defaultInvoiceTerms,omitempty"`
	DefaultInvoiceTermsUnit *settings_sales_enums.InvoiceTermUnit `json:"defaultInvoiceTermsUnit,omitempty"`
}
