package stats_dto

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatsFilterDto struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	PaymentMethod string `json:"paymentMethod"`
	BranchId      string `json:"branchId"`
	ProductId     string `json:"productId"`
	DeliveryType  string `json:"deliveryType"` // e.g. "standard", "express"
}

type OrderStatsResp struct {
	ActiveCount       int         `json:"activeCount"`
	InActiveCount     int         `json:"inActiveCount"`
	TotalCount        float64     `json:"totalCount"`
	OrdersComparision Comparision `json:"ordersComp"`
	GrowthRateComp    Comparision `json:"growthRateComp"`
	// Tiny line chart series for active/inactive customers within the selected range (max 7 buckets)
	ActiveSeries     []TinySeriesPoint `json:"activeSeries,omitempty"`
	InActiveSeries   []TinySeriesPoint `json:"inActiveSeries,omitempty"`
	TotalCountSeries []TinySeriesPoint `json:"totalCountSeries,omitempty"`
}

type Comparision struct {
	Current    float64 `json:"current"`
	Prev       float64 `json:"prev"`
	Percentage float64 `json:"percentage"`
}

type InvoiceStatsResp struct {
	PaidInvCount                int               `json:"paidInvoicesCount"`
	UnpaidInvCount              int               `json:"unpaidInvoicesCount"`
	PartiallyPaid               int               `json:"partiallyPaidInvoicesCount"`
	PaidInvoicesSeries          []TinySeriesPoint `json:"paidInvoicesSeries,omitempty"`
	UnpaidInvoicesSeries        []TinySeriesPoint `json:"unpaidInvoicesSeries,omitempty"`
	PartiallyPaidInvoicesSeries []TinySeriesPoint `json:"partiallyPaidInvoicesSeries,omitempty"`
}

// TinySeriesPoint represents one point in a tiny line chart series
type TinySeriesPoint struct {
	Value int    `json:"value"`
	Date  string `json:"date"`
}

// Top-performing companies (by buyer profile) for the current month
type TopCompanyStat struct {
	CompanyId        string  `json:"companyId"`
	CompanyName      *string `json:"companyName,omitempty"`
	OrdersThisMonth  int     `json:"ordersThisMonth"`
	TotalOrders      int     `json:"totalOrders"`
	RevenueThisMonth float64 `json:"revenueThisMonth"`
	TotalRevenue     float64 `json:"totalRevenue"`
}

type TopCompaniesResp = types.DataWithPagination[TopCompanyStat]

// Filter for viewing company orders details
type CompanyOrdersFilterDto struct {
	CompanyId string `json:"companyId" validate:"required,objectid"`
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	BranchId  string `json:"branchId,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type TopCompaniesFilterDto struct {
	Limit     int64
	Skip      int64
	StartDate string
	EndDate   string
	BranchId  string
}

// A compact order details record for company-level analysis
type CompanyOrderDetail struct {
	OrderDate     time.Time `json:"orderDate"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	CouponApplied bool      `json:"couponApplied"`
	CustomerName  string    `json:"customerName"`
	OrderCode     string    `json:"orderCode"`
}

// Paginated response for company orders details
type CompanyOrdersWithPagination = types.DataWithPagination[CompanyOrderDetail]

// CouponUsageDetail holds coupon usage info per order
type CouponUsageDetail struct {
	Code          string    `json:"code"`
	UsageDate     time.Time `json:"usageDate"`
	DiscountValue float64   `json:"discountValue"`
}

type CustomersActivityFilterDto struct {
	StartDate             string   `json:"startDate"`
	EndDate               string   `json:"endDate"`
	MinOrders             int      `json:"minOrders"`
	MaxOrders             *int     `json:"maxOrders,omitempty"`
	BranchId              string   `json:"branchId,omitempty"`
	CustomerName          string   `json:"customerName,omitempty"`
	Email                 string   `json:"email,omitempty"`
	PhoneNumber           string   `json:"phoneNumber,omitempty"`
	MinTotalQuantity      *int     `json:"minTotalQuantity,omitempty"`
	MaxTotalQuantity      *int     `json:"maxTotalQuantity,omitempty"`
	MinTotalSalesAmount   *float64 `json:"minTotalSalesAmount,omitempty"`
	MaxTotalSalesAmount   *float64 `json:"maxTotalSalesAmount,omitempty"`
	MinAverageOrderValue  *float64 `json:"minAverageOrderValue,omitempty"`
	MaxAverageOrderValue  *float64 `json:"maxAverageOrderValue,omitempty"`
	MinUsedCouponsCount   *int     `json:"minUsedCouponsCount,omitempty"`
	MaxUsedCouponsCount   *int     `json:"maxUsedCouponsCount,omitempty"`
	CustomerType          string   `json:"customerType,omitempty"` // "new" or "returning"
	LastBranchOrderedFrom string   `json:"lastBranchOrderedFrom,omitempty"`
	MinLastOrderDate      string   `json:"minLastOrderDate,omitempty"`
	MaxLastOrderDate      string   `json:"maxLastOrderDate,omitempty"`
	SortBy                string   `json:"sortBy,omitempty"`
	SortDir               string   `json:"sortDir,omitempty"`
	Limit                 int64
	Skip                  int64
}

// CustomerActivityRecord represents per-customer stats for current month
type CustomerActivityRecord struct {
	CustomerId            primitive.ObjectID  `json:"customerId"`
	CustomerName          string              `json:"customerName"`
	CompanyName           *string             `json:"companyName,omitempty"`
	Email                 *string             `json:"email,omitempty"`
	PhoneNumber           *string             `json:"phoneNumber,omitempty"`
	TotalOrdersCount      int                 `json:"totalOrdersCount"`
	TotalQuantityOrdered  int                 `json:"totalQuantityOrdered"`
	TotalMonthSalesAmount float64             `json:"totalMonthSalesAmount"`
	LastOrderDate         time.Time           `json:"lastOrderDate"`
	UsedCouponsCount      int                 `json:"usedCouponsCount"`
	CouponDetails         []CouponUsageDetail `json:"couponDetails,omitempty"`

	AverageOrderValue     float64 `json:"averageOrderValue"` // total amount of orders / total number of orders
	CustomerType          string  `json:"customerType"`      // new (first time ordering) OR returning  (has placed at least two orders last 90 days)]
	LastBranchOrderedFrom string  `json:"lastBranchOrderedFrom"`
	IsTopCustomer         bool    `json:"isTopCustomer"` // flag indicating if this is the highest customer among all
}

type CustomersActivityResp = types.DataWithPagination[CustomerActivityRecord]

// Top used coupons within a date range
type CouponUsageStat struct {
	CouponCode         string     `json:"couponCode"`
	RedeemCount        int        `json:"redeemCount"`
	UniqueCustomers    int        `json:"uniqueCustomers"`
	TotalDiscountValue float64    `json:"totalDiscountValue"`
	CouponValuePerUse  float64    `json:"couponValuePerUse"`
	ValidityStartDate  *time.Time `json:"validityStartDate,omitempty"`
	ValidityEndDate    *time.Time `json:"validityEndDate,omitempty"`
	CouponType         string     `json:"couponType"`  // e.g., "percentage", "fixed", "free_delivery"
	Status             string     `json:"status"`      // "active", "expired", "scheduled"
	IsTopCoupon        bool       `json:"isTopCoupon"` // flag indicating if this is the highest coupon among all
}

type TopCouponsFilterDto struct {
	StartDate             string   `json:"startDate"`
	EndDate               string   `json:"endDate"`
	BranchId              string   `json:"branchId,omitempty"`
	CouponCode            string   `json:"couponCode,omitempty"`
	CouponType            string   `json:"couponType,omitempty"`
	Status                string   `json:"status,omitempty"`
	MinRedeemCount        *int     `json:"minRedeemCount,omitempty"`
	MaxRedeemCount        *int     `json:"maxRedeemCount,omitempty"`
	MinUniqueCustomers    *int     `json:"minUniqueCustomers,omitempty"`
	MaxUniqueCustomers    *int     `json:"maxUniqueCustomers,omitempty"`
	MinTotalDiscountValue *float64 `json:"minTotalDiscountValue,omitempty"`
	MaxTotalDiscountValue *float64 `json:"maxTotalDiscountValue,omitempty"`
	MinCouponValuePerUse  *float64 `json:"minCouponValuePerUse,omitempty"`
	MaxCouponValuePerUse  *float64 `json:"maxCouponValuePerUse,omitempty"`
	SortBy                string   `json:"sortBy,omitempty"`
	SortDir               string   `json:"sortDir,omitempty"`
	Limit                 int64    `json:"limit"`
	Skip                  int64    `json:"skip"`
}

// Coupon orders details filtering
type CouponOrdersFilterDto struct {
	CouponCode  string `json:"couponCode" validate:"required"`
	StartDate   string `json:"startDate,omitempty"`
	EndDate     string `json:"endDate,omitempty"`
	BranchId    string `json:"branchId,omitempty"`
	OrderSource string `json:"orderSource,omitempty"`
	Status      string `json:"status,omitempty"`
	Page        int    `json:"page,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

// Detailed order record for a specific coupon usage
type CouponOrderDetail struct {
	OrderId        string    `json:"orderId"`
	OrderDate      time.Time `json:"orderDate"`
	OrderSource    *string   `json:"orderSource,omitempty"`
	CustomerName   string    `json:"customerName"`
	OriginalAmount float64   `json:"originalOrderAmount"`
	DiscountAmount float64   `json:"discountAmount"`
	FinalAmount    float64   `json:"finalOrderAmount"`
	OrderStatus    string    `json:"orderStatus"`
}

type CouponOrdersWithPagination = types.DataWithPagination[CouponOrderDetail]
