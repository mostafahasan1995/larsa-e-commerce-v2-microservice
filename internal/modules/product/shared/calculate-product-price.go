package product_shared

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	util "larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CalculateCartItemsPricesArgs struct {
	Product       product_models.ProductResponse
	CurrencyCode  string
	BranchId      *primitive.ObjectID
	Branch        *product_models.BranchResponse
	UnitId        *primitive.ObjectID
	Quantity      uint64
	Addons        *[]cart_shared.CartAddon
	CalculateList bool
}

func CalculateProductPrice(ctx context.Context, args *CalculateCartItemsPricesArgs) (*product_models.CalculateProductPriceResult, *helpers.ApiError) {
	db := do.MustInvoke[*mongo.Client](util.GetInjector(ctx))
	branchProductRepo := base_repos.NewBaseRepo[product_models.BranchProduct, product_models.BranchProductResponse](db, enums.BRANCH_PRODUCT)
	addonRepo := base_repos.NewBaseRepo[product_models.Addon, product_models.AddonResponse](db, enums.ADDON)

	if args.BranchId != nil {
		/*
		 * NOTE: Get the branch product even if it's not active or soft deleted.
		 * because we've checked the availability in the add to cart step.
		 * So if the product branch is not active or soft deleted, it means the product can't be added to the cart,
		 * but can be checked out.
		 * Timestamp: 2025-05-27 17:42:00
		 */
		branchProduct, err := branchProductRepo.FindOne(ctx, []bson.M{
			{
				"$match": bson.M{
					"productId": args.Product.Id,
					"branchId":  args.BranchId,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		if branchProduct == nil {
			branchIdMessage := "this branch"
			if args.BranchId != nil {
				branchIdMessage = "this branch with id " + args.BranchId.Hex()
			}

			return nil, &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "product is not available in this branch",
				Details: &[]types.ValidationError{
					{
						Message: "product with id " + args.Product.Id.Hex() + " is not available in " + branchIdMessage,
					},
				},
			}
		}
	}

	addonValues := []product_models.ProductAddonValue{}
	addonTotalPrice := 0.0
	addons := []product_models.AddonResponse{}
	if args.Addons != nil {
		addonIds := []primitive.ObjectID{}
		for _, addon := range *args.Addons {
			addonIds = append(addonIds, addon.Id)
		}

		storedAddons, err := addonRepo.GetByIds(ctx, addonIds, []bson.M{}, base_utils.AdditionalOptions{
			ReturnError: true,
		})

		if err != nil {
			return nil, err
		}

		for _, addon := range storedAddons {
			for _, value := range addon.AddonValues {
				argsAddon := util.Find(*args.Addons, func(argsAddon cart_shared.CartAddon) bool {
					return addon.Id == argsAddon.Id
				})

				argsAddonValue := util.Find(argsAddon.Values, func(argsAddonValue cart_shared.CartAddonValue) bool {
					return value.Id == argsAddonValue.Id
				})

				/*
				 * NOTE: If the addon value is not nil, we will calculate the price.
				 * because the addon value is not nil, it means the addon value is not hard deleted
				 * Timestamp: 2025-05-26 17:42:00
				 */
				if argsAddonValue != nil {
					price := value.Price * float64(argsAddonValue.Quantity)
					addonTotalPrice = addonTotalPrice + price
					addonValues = append(addonValues, product_models.ProductAddonValue{
						Id:         value.Id,
						TotalPrice: price,
						Quantity:   argsAddonValue.Quantity,
					})
				}
			}
		}

		addons = storedAddons
	}

	extractedPriceResult, priceErr := ExtractPrice(ExtractPriceArgs{
		ProductId:    args.Product.Id,
		Prices:       args.Product.Prices,
		BranchId:     args.BranchId,
		CurrencyCode: args.CurrencyCode,
		UnitId:       args.UnitId,
	})

	if priceErr != nil {
		return nil, priceErr
	}

	priceItemList := []product_models.CalculateProductPriceResult{}
	priceItem := extractedPriceResult.PriceItem

	if args.CalculateList {
		for _, item := range extractedPriceResult.Price.Items {
			result, err := CalculatePriceItem(ctx, CalculatePriceItemArgs{
				Price:        *extractedPriceResult.Price,
				CurrencyCode: args.CurrencyCode,
				UnitId:       args.UnitId,
				ProductId:    args.Product.Id,
				BranchId:     args.BranchId,
				Branch:       args.Branch,
				Quantity:     args.Quantity,
				Addons:       addons,
				AddonValues:  addonValues,
				priceItem:    &item,
			})
			if err != nil {
				return nil, err
			}
			priceItemList = append(priceItemList, *result)
		}
	}

	calculatePriceItemResult, err := CalculatePriceItem(ctx, CalculatePriceItemArgs{
		Price:        *extractedPriceResult.Price,
		CurrencyCode: args.CurrencyCode,
		UnitId:       args.UnitId,
		ProductId:    args.Product.Id,
		BranchId:     args.BranchId,
		Branch:       args.Branch,
		Quantity:     args.Quantity,
		Addons:       addons,
		AddonValues:  addonValues,
		priceItem:    priceItem,
	})

	if err != nil {
		return nil, err
	}

	calculatePriceItemResult.List = &priceItemList

	return calculatePriceItemResult, nil
}
