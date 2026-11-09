package product_shared

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	util "larsa-e-commerce-microservice-v2/internal/common/util"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExtractPriceArgs struct {
	ProductId    primitive.ObjectID
	Prices       []product_models.Price
	BranchId     *primitive.ObjectID
	CurrencyCode string
	UnitId       *primitive.ObjectID
}

type ExtractPriceResult struct {
	Price     *product_models.Price
	PriceItem *product_models.PriceItem
}

func ExtractPrice(args ExtractPriceArgs) (*ExtractPriceResult, *helpers.ApiError) {
	var extractedPriceItem *product_models.PriceItem

	/*
	 * NOTE: If the override price is not found, we will get the price based on the branch id
	 * Timestamp: 2025-05-26 17:42:00
	 */
	price := util.Find(args.Prices, func(price product_models.Price) bool {

		/*
		 * NOTE: If the candidate branch id is nil or the price has no branch ids, we will return false
		 * since the price is not related to the branch id
		 * Timestamp: 2025-06-25 17:42:00
		 */
		if args.BranchId == nil || price.BranchIds == nil || len(*price.BranchIds) == 0 {
			return false
		}

		/*
		 * NOTE: We're only need to check if the branch id is in the list of branch ids
		 * because the branch id is not null
		 * the branch id may exist more than once in the list of branch ids
		 * we need to get the first branch id that is not null
		 * Timestamp: 2025-05-12 17:42:00
		 */
		branchId := util.Find(*price.BranchIds, func(branchId primitive.ObjectID) bool {
			return branchId.Hex() == args.BranchId.Hex()
		})

		if branchId == nil {
			return false
		}

		var doesBranchExist = branchId.Hex() == args.BranchId.Hex()

		if doesBranchExist {
			extractedPriceItem = ExtractPriceItem(price.Items, args.CurrencyCode, args.UnitId)
			if extractedPriceItem != nil {
				return true
			}
		}

		return extractedPriceItem != nil
	})

	if price == nil {
		price = util.Find(args.Prices, func(price product_models.Price) bool {
			return price.IsDefault != nil && *price.IsDefault
		})
	}

	if price != nil {
		extractedPriceItem = ExtractPriceItem(price.Items, args.CurrencyCode, args.UnitId)
	}

	/**
	 * NOTE: If the price is not found, we will get the price based only on the currency code
	 * Timestamp: 2025-06-30 17:42:00
	 */
	if price != nil && extractedPriceItem == nil {
		extractedPriceItem = ExtractPriceItem(price.Items, args.CurrencyCode, nil)
	}

	if price == nil {
		branchIdMessage := "default branch"
		if args.BranchId != nil {
			branchIdMessage = "the branch with id " + args.BranchId.Hex()
		}

		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "The Product price has not been set for the branch",
			Details: &[]types.ValidationError{
				{
					Message: "The Product with id " + args.ProductId.Hex() + " has no price for " + branchIdMessage,
				},
			},
		}
	}

	return &ExtractPriceResult{
		Price:     price,
		PriceItem: extractedPriceItem,
	}, nil
}

func ExtractPriceItem(priceItems []product_models.PriceItem, currencyCode string, unitId *primitive.ObjectID) *product_models.PriceItem {
	var resultPriceItem *product_models.PriceItem

	if len(priceItems) == 0 {
		return nil
	}

	for _, priceItem := range priceItems {
		if unitId == nil {
			if priceItem.CurrencyCode == currencyCode {
				resultPriceItem = &priceItem
				break
			}
		}

		if priceItem.UnitId == nil {
			continue
		}

		if priceItem.CurrencyCode == currencyCode && priceItem.UnitId.Hex() == unitId.Hex() {
			resultPriceItem = &priceItem
			break
		}

	}

	return resultPriceItem
}
