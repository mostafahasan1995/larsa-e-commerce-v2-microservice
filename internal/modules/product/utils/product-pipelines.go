package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProductRelatedBranchPipeline(branchId primitive.ObjectID, outOfStock *product_enums.OutOfStock) []bson.M {
	outOfStockFilter := bson.M{
		"$match": bson.M{
			"$or": []bson.M{
				{
					"outOfStock": false,
				},
				{
					"outOfStock": nil,
				},
				{
					"outOfStock": bson.M{"$exists": false},
				},
			},
		},
	}

	if outOfStock != nil {
		switch *outOfStock {
		case product_enums.OUT_OF_STOCK_IN_STOCK:
			break
		case product_enums.OUT_OF_STOCK_OUT_OF_STOCK:
			outOfStockFilter = bson.M{
				"$match": bson.M{"outOfStock": true},
			}
		case product_enums.OUT_OF_STOCK_UNKNOWN:
			outOfStockFilter = bson.M{"$match": bson.M{}}
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"attributes.code": product_enums.SYSTEM_BRANCH,
				/*
				 * NOTE:
				 * This's a workaround to support both ObjectID and string values for branchId
				 * because the branchId is stored as a string in the database
				 * and we need right now to support both types of values in the filter
				 * to avoid wrong filter results.
				 *
				 * timestamp: 2025-06-09:15:30:00
				 */
				"$or": []bson.M{
					{
						"attributes.value": branchId,
					},
					{
						"attributes.value": branchId.Hex(),
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         enums.BRANCH_PRODUCT.GetCollectionName(),
				"localField":   "_id",
				"foreignField": "productId",
				"as":           "branchProducts",
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"branchId": branchId,
							"isActive": true,
						},
					},
					{
						"$match": base_utils.AddSoftDeletedFilter(bson.M{}),
					},
					outOfStockFilter,
				},
			},
		},
		{
			"$addFields": bson.M{
				"branchProductsSize": bson.M{
					"$size": "$branchProducts",
				},
			},
		},
		{
			"$match": bson.M{
				"branchProductsSize": bson.M{
					"$gte": 1,
				},
			},
		},
	}

	return append(pipeline, productOutOfStockPipeline...)
}

var productOutOfStockPipeline = []bson.M{
	{
		"$addFields": bson.M{
			"branchProduct": bson.M{
				"$ifNull": bson.A{
					bson.M{
						"$arrayElemAt": bson.A{
							"$branchProducts",
							0,
						},
					},
					bson.M{},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"outOfStock": bson.M{
				"$ifNull": bson.A{
					"$branchProduct.outOfStock",
					true,
				},
			},
		},
	},
}
