package product_shared

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculatePriceItemArgs struct {
	Price        product_models.Price
	CurrencyCode string
	UnitId       *primitive.ObjectID
	ProductId    primitive.ObjectID
	BranchId     *primitive.ObjectID
	Branch       *product_models.BranchResponse
	Quantity     uint64
	Addons       []product_models.AddonResponse
	AddonValues  []product_models.ProductAddonValue
	priceItem    *product_models.PriceItem
}

func CalculatePriceItem(ctx context.Context, args CalculatePriceItemArgs) (*product_models.CalculateProductPriceResult, *helpers.ApiError) {
	subTotalPrice := 0.0
	addonTotalPrice := 0.0

	if args.AddonValues != nil {
		for _, addonValue := range args.AddonValues {
			addonTotalPrice = addonTotalPrice + (addonValue.TotalPrice * float64(args.Quantity))
		}
	}

	if args.priceItem == nil {
		branchIdMessage := "default branch"
		if args.BranchId != nil {
			branchIdMessage = "the branch with id " + args.BranchId.Hex()
		}

		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "The Product price has not set a currency code for the branch",
			Details: &[]types.ValidationError{
				{
					Message: "The Product with id " + args.ProductId.Hex() + " has no price with the currency code " + args.CurrencyCode + " for " + branchIdMessage,
				},
			},
		}
	}

	subTotalPrice = subTotalPrice + args.priceItem.Value

	discountValue := 0.0

	if args.Price.Discount != nil && args.priceItem.CurrencyCode == args.CurrencyCode {
		discountValue = cart_shared.CalculateDiscount(*args.Price.Discount, cart_shared.AddDiscountArgs{
			CurrencyCode: args.CurrencyCode,
			Amount:       subTotalPrice,
		})
	}

	totalPrice := subTotalPrice - discountValue

	// Calculate promotion discount
	promotionResult, err := CalculatePromotionDiscount(ctx, CalculatePromotionDiscountArgs{
		ProductId:    args.ProductId,
		BranchId:     args.BranchId,
		Branch:       args.Branch, // Pass the Branch field
		TotalPrice:   totalPrice,
		CurrencyCode: args.CurrencyCode,
	})

	if err != nil {
		// If promotion calculation fails, continue without promotion
		promotionResult = &CalculatePromotionDiscountResult{
			ActivePromotion:          nil,
			TotalPriceAfterPromotion: totalPrice,
			PromotionDiscountValue:   0,
		}
	}

	totalPrice = totalPrice - promotionResult.PromotionDiscountValue
	discountValue = discountValue + promotionResult.PromotionDiscountValue

	return &product_models.CalculateProductPriceResult{
		SubTotalPrice:   (subTotalPrice * float64(args.Quantity)) + addonTotalPrice,
		TotalPrice:      (totalPrice * float64(args.Quantity)) + addonTotalPrice,
		DiscountValue:   discountValue * float64(args.Quantity),
		PriceItem:       *args.priceItem,
		Discount:        args.Price.Discount,
		AddonValues:     args.AddonValues,
		AddonTotalPrice: addonTotalPrice,
		Addons:          &args.Addons,
		ActivePromotion: promotionResult.ActivePromotion,
	}, nil
}
