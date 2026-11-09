package loyalty_dtos

import (
	"time"
)

// DiscountEligibility represents the eligibility status for loyalty discounts
type DiscountEligibility struct {
	IsEligible         bool       `json:"isEligible"`
	DiscountPercentage float64    `json:"discountPercentage"`
	OrdersRemaining    int64      `json:"ordersRemaining"`
	PeriodEndDate      *time.Time `json:"periodEndDate,omitempty"`
}

// DiscountResult represents the result of applying a loyalty discount
type DiscountResult struct {
	OriginalAmount     float64 `json:"originalAmount"`
	DiscountAmount     float64 `json:"discountAmount"`
	FinalAmount        float64 `json:"finalAmount"`
	DiscountPercentage float64 `json:"discountPercentage"`
	OrdersRemaining    int64   `json:"ordersRemaining"`
}
