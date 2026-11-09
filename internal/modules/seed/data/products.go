package seed_data

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	mostPopularFlag = product_enums.ProductFlag(product_enums.MOST_POPULAR)
)

var Products = []product_dtos.CreateProductDto{
	{
		Name: types.Localizable[string]{
			enums.EN: "iPhone 14 Pro",
			enums.AR: "آيفون 14 برو",
		},
		Description: types.Localizable[string]{
			enums.EN: "Latest iPhone model with advanced features",
			enums.AR: "أحدث موديل آيفون مع ميزات متقدمة",
		},
		Images: []types.FileFieldDto{
			{
				Id:   primitive.NewObjectID().Hex(),
				Path: "https://example.com/images/mens-accessories-2.jpg",
			},
			{
				Id:   primitive.NewObjectID().Hex(),
				Path: "https://example.com/images/mens-accessories-2.jpg",
			},
		},
		Attributes: []product_dtos.AttributeDto{
			{
				Code:  "color",
				Value: "Deep Purple",
				Type:  "string",
			},
			{
				Code:  "storage",
				Value: "256GB",
				Type:  "string",
			},
		},
		Prices: []product_dtos.CreatePriceDto{
			{
				Name: &types.Localizable[string]{
					enums.EN: "MacBook Pro 16",
					enums.AR: "ماك بوك برو 16",
				},
				BranchIds: &[]primitive.ObjectID{primitive.NewObjectID()},
				Items: []product_dtos.PriceItemDto{
					{
						CurrencyCode: "USD",
						Value:        999.99,
					},
					{
						CurrencyCode: "SAR",
						Value:        3749.99,
					},
				},
			},
		},
		Flag: &mostPopularFlag,
	},
}
