package cart_shared

import (
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type AddDiscountArgs struct {
	CurrencyCode string
	Amount       float64
}

func CalculateDiscount(discount product_models.Discount, args AddDiscountArgs) float64 {
	discountValue := 0.0

	if discount.Type == product_enums.PERCENTAGE {
		discountValue = args.Amount * discount.Value / 100
	} else if discount.Type == product_enums.FIXED {
		if discount.CurrencyCode != nil && *discount.CurrencyCode == args.CurrencyCode {
			discountValue = discount.Value
		}
	}

	return discountValue
}
