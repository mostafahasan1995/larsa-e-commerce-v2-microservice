package stats_dto

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomerProfileFilterDto represents the filter for customer profile data
type CustomerProfileFilterDto struct {
	CustomerId string `json:"customerId" validate:"required,objectid"`
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
}

// CustomerProfileDto represents the basic customer information
type CustomerProfileDto struct {
	Id          primitive.ObjectID `json:"id"`
	FullName    string             `json:"fullName"`
	CompanyName *string            `json:"companyName,omitempty"`
	Email       *string            `json:"email,omitempty"`
	PhoneNumber *string            `json:"phoneNumber,omitempty"`
}

// OrdersSummaryDto represents the summary cards for customer orders
type OrdersSummaryDto struct {
	TotalOrders            int                `json:"totalOrders"`
	TotalOrderValue        float64            `json:"totalOrderValue"`
	OrdersPerChannel       map[string]int     `json:"ordersPerChannel"`
	OrdersValuePerChannel  map[string]float64 `json:"ordersValuePerChannel"`
	TotalQuantityPurchased int                `json:"totalQuantityPurchased"`
	UncompletedOrdersCount int                `json:"uncompletedOrdersCount"`
}

// CustomerOrderDto represents a single order in the customer's order history
type CustomerOrderDto struct {
	OrderId        string    `json:"orderId"`
	OrderDate      time.Time `json:"orderDate"`
	OrderStatus    string    `json:"orderStatus"`
	SalesChannel   string    `json:"salesChannel"`
	BranchName     *string   `json:"branchName,omitempty"`
	InvoiceTotal   float64   `json:"invoiceTotal"`
	PaymentMethod  string    `json:"paymentMethod"`
	PaymentStatus  string    `json:"paymentStatus"`
	UsedCouponCode *string   `json:"usedCouponCode,omitempty"`
}

// CustomerOrdersTableDto represents the paginated orders table response
type CustomerOrdersTableDto = types.DataWithPagination[CustomerOrderDto]

// CustomerOrdersFilterDto represents the filter for customer orders table
type CustomerOrdersFilterDto struct {
	CustomerId    string  `json:"customerId" validate:"required,objectid"`
	StartDate     *string `json:"startDate,omitempty"`
	EndDate       *string `json:"endDate,omitempty"`
	OrderStatus   *string `json:"orderStatus,omitempty"`
	BranchId      *string `json:"branchId,omitempty"`
	Channel       *string `json:"channel,omitempty"`
	PaymentMethod *string `json:"paymentMethod,omitempty"`
	Page          int     `json:"page,omitempty"`
	Limit         int     `json:"limit,omitempty"`
}

// CouponUsageSummaryDto represents the summary card for coupon usage
type CouponUsageSummaryDto struct {
	TotalCoupons   int     `json:"totalCoupons"`
	TotalDiscount  float64 `json:"totalDiscount"`
	LastCouponUsed *string `json:"lastCouponUsed,omitempty"`
}

// CustomerCouponUsageDto represents detailed coupon usage information
type CustomerCouponUsageDto struct {
	CouponCode    string    `json:"couponCode"`
	DiscountValue float64   `json:"discountValue"`
	InvoiceTotal  float64   `json:"invoiceTotal"`
	DateUsed      time.Time `json:"dateUsed"`
	OrderId       string    `json:"orderId"`
	Channel       string    `json:"channel"`
	Branch        *string   `json:"branch,omitempty"`
}

// CustomerCouponUsageDetailsDto represents the detailed coupon usage table response
type CustomerCouponUsageDetailsDto = types.DataWithPagination[CustomerCouponUsageDto]

// CustomerProfileResponseDto represents the complete customer profile response
type CustomerProfileResponseDto struct {
	Profile            CustomerProfileDto    `json:"profile"`
	OrdersSummary      OrdersSummaryDto      `json:"ordersSummary"`
	CouponUsageSummary CouponUsageSummaryDto `json:"couponUsageSummary"`
}
