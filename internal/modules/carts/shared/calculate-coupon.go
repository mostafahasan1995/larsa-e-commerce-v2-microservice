package cart_shared

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type CalculateCouponArgs struct {
	TotalPrice   float64
	CurrencyCode string
	Coupon       coupon_models.CouponResponse
}

type CalculateCouponResult struct {
	CouponDiscountValue float64
}

func CalculateCoupon(ctx context.Context, calculateCouponArgs CalculateCouponArgs) CalculateCouponResult {
	couponDiscountValue := 0.0

	if calculateCouponArgs.Coupon.AppliedOn == coupon_enums.GENERAL || calculateCouponArgs.Coupon.AppliedOn == coupon_enums.PRODUCT {
		var couponUpTo *product_models.PriceItem
		if calculateCouponArgs.Coupon.Discount.UpTo != nil {
			couponUpTo = util.Find(*calculateCouponArgs.Coupon.Discount.UpTo, func(priceItem product_models.PriceItem) bool {
				return priceItem.CurrencyCode == calculateCouponArgs.CurrencyCode
			})
		}

		var discountModel = product_models.Discount{
			Value:        calculateCouponArgs.Coupon.Discount.Value,
			Type:         calculateCouponArgs.Coupon.Discount.Type,
			CurrencyCode: calculateCouponArgs.Coupon.Discount.CurrencyCode,
		}
		discountValue := CalculateDiscount(discountModel, AddDiscountArgs{
			CurrencyCode: calculateCouponArgs.CurrencyCode,
			Amount:       calculateCouponArgs.TotalPrice,
		})

		if couponUpTo != nil && couponUpTo.Value < discountValue {
			discountValue = couponUpTo.Value
		}

		couponDiscountValue = discountValue
	}

	return CalculateCouponResult{
		CouponDiscountValue: couponDiscountValue,
	}

}
