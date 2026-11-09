package settings_sales_models

import (
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
)

type BaseSalesSettings struct {
	models.MetaDataBaseModel `bson:",inline"`

	OrderStatuses []OrderStatus `bson:"orderStatuses" json:"orderStatuses"`

	ShowVendorInfoInOrderDetails bool             `bson:"showVendorInfoInOrderDetails" json:"showVendorInfoInOrderDetails"`
	LogisticSettings             LogisticSettings `bson:"logisticSettings" json:"logisticSettings"`
	RecentActivity               RecentActivity   `bson:"recentActivity" json:"recentActivity"`
	InvoiceTerms                 InvoiceTerms     `bson:"invoiceTerms" json:"invoiceTerms"`

	// Platform order flow mode: b2c or b2b
	OrderFlowMode settings_sales_enums.OrderFlowMode `bson:"orderFlowMode" json:"orderFlowMode"`

	EnableDiscountOnSales bool    `bson:"enableDiscountOnSales" json:"enableDiscountOnSales"`
	EnablePartialPayment  bool    `bson:"enablePartialPayment" json:"enablePartialPayment"`
	SalesTaxesRate        float64 `bson:"salesTaxesRate" json:"salesTaxesRate"`

	//TODO:
	// EnableAutoAssignSalesperson     objectId         `bson:"enableAutoAssignSalesperson" json:"enableAutoAssignSalesperson"`
	EnableSalesNotes                bool `bson:"enableSalesNotes" json:"enableSalesNotes"`
	EnableAutoGenerateInvoiceNumber bool `bson:"enableAutoGenerateInvoiceNumber" json:"enableAutoGenerateInvoiceNumber"`
	EnableInvoiceEmailNotification  bool `bson:"enableInvoiceEmailNotification" json:"enableInvoiceEmailNotification"`
}

type SalesSettings struct {
	BaseSalesSettings `bson:",inline"`
}

type SalesSettingsResponse struct {
	BaseSalesSettings `bson:",inline"`
}
