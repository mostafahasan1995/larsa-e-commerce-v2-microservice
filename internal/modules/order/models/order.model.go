package order_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseOrder struct {
	models.MetaDataBaseModel `bson:",inline"`

	//Type: pickup, delivery
	Type order_enums.OrderType `bson:"type" json:"type"`

	//PickUp: (pickup info)
	PickupInfo *PickupInfo `bson:"pickupInfo,omitempty" json:"pickupInfo,omitempty"`

	//Delivery: (delivery info)
	DeliveryInfo *DeliveryInfo `bson:"deliveryInfo,omitempty" json:"deliveryInfo,omitempty"`
	//DeliveryInfo:
	/*
		Type: B2B, B2C
			B2B
				shippen method
				tracking number
				transport company
				shipment date
				delivery fee
				shipment status
			B2C
				customer address
				delivery fee
	*/

	TotalDiscountValue *float64 `bson:"totalDiscountValue,omitempty" json:"totalDiscountValue,omitempty"`

	//Payment
	/*
		status
		payments
			payment method
			payment status
			payment date
			payment amount
	*/
	PaymentInfo *PaymentInfo `bson:"paymentInfo,omitempty" json:"paymentInfo,omitempty"`

	InvoiceNumber *string `bson:"invoiceNumber,omitempty" json:"invoiceNumber,omitempty"`

	//Note: it is not an enum, because it is dynamic based on sales settings
	Status string `bson:"status" json:"status"`

	//TotalPriceBeforeDiscount
	SubTotalPrice float64 `bson:"subTotalPrice" json:"subTotalPrice"`

	//CurrencyCode
	CurrencyCode string `bson:"currencyCode" json:"currencyCode"`

	CouponCode *string `bson:"couponCode,omitempty" json:"couponCode,omitempty"`

	CouponDiscountValue float64 `bson:"couponDiscountValue" json:"couponDiscountValue"`

	LoyaltyDiscountValue float64 `bson:"loyaltyDiscountValue" json:"loyaltyDiscountValue"`

	TotalPriceAfterProductDiscount float64 `bson:"totalPriceAfterProductDiscount" json:"totalPriceAfterProductDiscount"`

	CartId *primitive.ObjectID `bson:"cartId,omitempty" json:"cartId,omitempty"`

	Note *string `bson:"note,omitempty" json:"note,omitempty"`

	//TotalPrice
	TotalPrice float64 `bson:"totalPrice" json:"totalPrice"`

	//BuyerInfo
	BuyerInfo OrderParticipantInfo `bson:"buyerInfo" json:"buyerInfo"`

	//VendorInfo
	VendorInfo OrderParticipantInfo `bson:"vendorInfo" json:"vendorInfo"`

	//OrderItems
	OrderItems []OrderItem `bson:"orderItems" json:"orderItems"`

	//Unique code for order
	Code string `bson:"code" json:"code"` // Auto generated with specific format

	Documents  []Document   `bson:"documents" json:"documents"`
	Comments   []Comment    `bson:"comments" json:"comments"`
	LogHistory []LogHistory `bson:"logHistory" json:"logHistory"`

	QuotationId *primitive.ObjectID `bson:"quotationId,omitempty" json:"quotationId,omitempty"`

	// Admin tracking
	AdminSeen bool `bson:"adminSeen" json:"adminSeen"`

	// We use this field to track if the order is migrated from old system
	MigratedFrom *string `bson:"migratedFrom,omitempty" json:"migratedFrom,omitempty"`
}

type Order struct {
	BaseOrder `bson:",inline"`
}

type OrderResponse struct {
	BaseOrder `bson:",inline"`
	Quotation *quotation_models.Quotation `bson:"quotation,omitempty" json:"quotation,omitempty"`
}

type OrderStats struct {
	TotalSales       float64 `bson:"totalSales"`
	NumberOfOrders   int     `bson:"numberOfOrders"`
	AvgOrderValue    float64 `bson:"avgOrderValue"`
	TopSellingItemId string  `bson:"topSellingItemId"`
	TopSellingName   string  `bson:"topSellingName"`
	TopQuantity      int     `bson:"topQuantity"`
}

type OrderReport struct {
	Orders   []OrderResponse `bson:"orders,omitempty" json:"orders,omitempty"`
	Stat     *OrderStats     `bson:"stat,omitempty" json:"stat,omitempty"`
	LastStat *OrderStats     `bson:"lastStat,omitempty" json:"lastStat,omitempty"`
	// LastStat OrderStats `bson:"orders,omitempty" json:"orders,omitempty"`
}

type OrdersWithPagination = types.DataWithPagination[OrderResponse]

type BranchSales struct {
	BranchId       string  `json:"branchId"`
	BranchName     string  `json:"branchName"`
	SalesValue     float64 `json:"salesValue"`
	NumberOfOrders int     `json:"numberOfOrders"`
}
type AggregatedBranchSales struct {
	GroupKey       string  `bson:"groupKey"`
	BranchId       string  `bson:"branchId"`
	SalesValue     float64 `bson:"salesValue"`
	NumberOfOrders int     `bson:"numberOfOrders"`
}
type ProductSales struct {
	Name       string  `json:"name"`
	Quantity   uint64  `json:"quantity"`
	SalesValue float64 `json:"salesValue"`
}
type OrderStat struct {
	CustomersPurchasing int     `json:"customersPurchasing"`
	TotalPurchases      int     `json:"totalPurchases"`
	TotalQuantity       uint64  `json:"totalQuantity"`
	TotalAmount         float64 `json:"totalAmount"`
}
type PaginatedOrdersWithStats struct {
	Orders []OrderResponse `bson:"orders"`
	Stats  []OrderStat     `bson:"stats"`
}
type CustomerResponse struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	CustomerName  string             `bson:"customerName" json:"customerName"`
	CustomerEmail string             `bson:"customerEmail" json:"customerEmail"`
	CustomerPhone string             `bson:"customerPhone" json:"customerPhone"`
	LocationCity  string             `bson:"locationCity" json:"locationCity"`
	TotalOrder    int                `bson:"totalOrder" json:"totalOrder"`
}

type PaginatedCustomersWithStats struct {
	Customers []CustomerResponse `bson:"customers" json:"customers"`
	Stats     []OrderStat        `bson:"stats" json:"stats"`
}

type OrderReportSales struct {
	Customers []CustomerResponse `json:"orders"`
	Stat      OrderStat          `json:"stat"`
	LastStat  OrderStat          `json:"lastStat"`
}
type ProductSalesReportResult struct {
	Report   []ProductSalesReport `json:"report"`
	Stat     ProductSalesStat     `json:"stat"`
	LastStat ProductSalesStat     `json:"lastStat"`
}

type ProductSalesReport struct {
	ProductName     string                    `json:"productName"`
	BranchId        *primitive.ObjectID       `json:"branchId"`
	OrderTotal      float64                   `json:"orderTotal"`
	Quantity        uint64                    `json:"quantity"`
	CustomerName    string                    `json:"customerName"`
	CustomerEmail   *string                   `json:"customerEmail,omitempty"`
	CustomerAddress *string                   `json:"customerAddress"`
	PaymentMethod   order_enums.PaymentMethod `json:"paymentMethod"`
}

type ProductSalesStat struct {
	TotalProductsSold  int     `json:"totalProductsSold"`
	TotalQuantity      uint64  `json:"totalQuantity"`
	TotalSales         float64 `json:"totalSales"`
	BestSellingProduct string  `json:"bestSellingProduct"`
}
type AggregatedProductSales struct {
	GroupKey    string  `bson:"groupKey" json:"groupKey"`
	BranchId    string  `bson:"branchId" json:"branchId"`
	ItemId      string  `bson:"itemId" json:"itemId"`
	ProductName string  `bson:"productName" json:"productName"`
	SalesValue  any `bson:"salesValue" json:"salesValue"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	Price       any     `bson:"price" json:"price"`
}
