package stats_dto

type CouponRedemptionsPerBranchFilterDto struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	CouponCode    string `json:"couponCode"`
	PaymentMethod string `json:"paymentMethod"`
}

// Response DTO for coupon redemptions per branch chart
type CouponRedemptionsPerBranchResponseDto struct {
	Data []CouponRedemptionPerBranchData `json:"data"`
}

type CouponRedemptionPerBranchData struct {
	BranchId               string  `json:"branchId" bson:"branchId"`
	BranchName             string  `json:"branchName" bson:"branchName"`
	TotalCouponRedemptions int     `json:"totalCouponRedemptions" bson:"totalCouponRedemptions"`
	TotalDiscountValue     float64 `json:"totalDiscountValue" bson:"totalDiscountValue"`
}

type CouponRedemptionsByCustomerTypeFilterDto struct {
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	CouponCode string `json:"couponCode"`
	// PaymentMethod string `json:"paymentMethod"`
}

type CouponRedemptionsByCustomerTypeData struct {
	CustomerType   string `json:"customerType" bson:"customerType"`
	CustomersCount int    `json:"customersCount" bson:"customersCount"`
	OrdersCount    int    `json:"ordersCount" bson:"ordersCount"`
}

type CouponRedemptionsByCustomerTypeResponseDto struct {
	Data []CouponRedemptionsByCustomerTypeData `json:"data"`
}

type CustomersCouponRedeemsBarFilterDto struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	BranchId  string `json:"branchId"`
	Limit     int64  `json:"limit"`
	Skip      int64  `json:"skip"`
}

type CustomersCouponRedeemsBarData struct {
	CustomerId   string  `json:"customerId" bson:"customerId"`
	CustomerName string  `json:"customerName" bson:"customerName"`
	TotalSales   float64 `json:"totalSales" bson:"totalSales"`
	DidRedeem    bool    `json:"didRedeem" bson:"didRedeem"`
}

type CustomersCouponRedeemsBarResponseDto struct {
	Data []CustomersCouponRedeemsBarData `json:"data"`
}
