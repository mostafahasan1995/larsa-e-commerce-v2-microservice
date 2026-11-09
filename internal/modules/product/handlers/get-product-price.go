package product_handlers

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPriceArgs struct {
	Product       product_models.ProductResponse
	BranchId      string
	UnitId        string
	CurrencyCode  string
	CalculateList bool
}

func GetPrice(ctx context.Context, args GetPriceArgs) *product_models.ProductResponse {
	branchId := args.BranchId
	currencyCode := args.CurrencyCode
	var unitId *primitive.ObjectID

	/*
	 * NOTE: If the unit id is not empty, we will convert it to a primitive object id
	 * Timestamp: 2025-06-25 17:42:00
	 */
	if args.UnitId != "" {
		objectId, err := primitive.ObjectIDFromHex(args.UnitId)
		if err != nil {
			return &args.Product
		}
		unitId = &objectId
	}

	if currencyCode == "" {
		currencyCode = configuration.EnvironmentVariables.DefaultCurrencyCode
	}

	var branchIdObject *primitive.ObjectID
	if branchId != "" {
		objectId, err := primitive.ObjectIDFromHex(branchId)
		if err != nil {
			return &args.Product
		}
		branchIdObject = &objectId
	}

	productPrice, err := product_shared.CalculateProductPrice(ctx, &product_shared.CalculateCartItemsPricesArgs{
		Product:       args.Product,
		Quantity:      1,
		BranchId:      branchIdObject,
		UnitId:        unitId,
		CurrencyCode:  currencyCode,
		CalculateList: args.CalculateList,
	})

	if err != nil {
		return &args.Product
	}

	args.Product.Price = productPrice

	return &args.Product
}
