package cart_shared

import (
	"context"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculateLoyaltyDiscountArgs struct {
	TotalPrice   float64
	CurrencyCode string
	UserId       primitive.ObjectID
}

type CalculateLoyaltyDiscountResult struct {
	LoyaltyDiscountValue float64
	IsEligible           bool
	DiscountPercentage   float64
	OrdersRemaining      int64
}

func CalculateLoyaltyDiscount(ctx context.Context, calculateLoyaltyDiscountArgs CalculateLoyaltyDiscountArgs, loyaltyDiscountService loyalty_services.LoyaltyDiscountService) CalculateLoyaltyDiscountResult {
	totalPrice := calculateLoyaltyDiscountArgs.TotalPrice
	loyaltyDiscountValue := 0.0
	isEligible := false
	discountPercentage := 0.0
	ordersRemaining := int64(0)

	// Check eligibility for loyalty discount
	eligibility, err := loyaltyDiscountService.CheckDiscountEligibility(ctx, calculateLoyaltyDiscountArgs.UserId)
	if err != nil {
		// If there's an error checking eligibility, continue without loyalty discount
		return CalculateLoyaltyDiscountResult{
			LoyaltyDiscountValue: 0,
			IsEligible:           false,
			DiscountPercentage:   0,
			OrdersRemaining:      0,
		}
	}

	if eligibility.IsEligible {
		isEligible = true
		discountPercentage = eligibility.DiscountPercentage
		ordersRemaining = eligibility.OrdersRemaining

		loyaltyDiscountValue = totalPrice * discountPercentage / 100
	}

	return CalculateLoyaltyDiscountResult{
		LoyaltyDiscountValue: loyaltyDiscountValue,
		IsEligible:           isEligible,
		DiscountPercentage:   discountPercentage,
		OrdersRemaining:      ordersRemaining,
	}
}
