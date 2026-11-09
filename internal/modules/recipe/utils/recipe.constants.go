package recipe_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

// PopulateRecipeProductsAndUnits builds a pipeline that:
// - looks up all products referenced by productItems.productId
// - looks up all units referenced by productItems.unitId
// - adds a computed "products" array that merges each item with its populated product and unit
// - removes temporary lookup arrays
var PopulateRecipeProductsAndUnits = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.PRODUCT.GetCollectionName(),
			"localField":   "productItems.productId",
			"foreignField": "_id",
			"as":           "populatedProducts",
		},
	},
	{
		"$lookup": bson.M{
			"from":         enums.UNIT.GetCollectionName(),
			"localField":   "productItems.unitId",
			"foreignField": "_id",
			"as":           "populatedUnits",
		},
	},
	{
		"$addFields": bson.M{
			"productItems": bson.M{
				"$map": bson.M{
					"input": bson.M{"$ifNull": []interface{}{"$productItems", []interface{}{}}},
					"as":    "item",
					"in": bson.M{
						"$mergeObjects": []interface{}{
							"$$item",
							bson.M{
								"product": bson.M{
									"$arrayElemAt": []interface{}{
										bson.M{
											"$filter": bson.M{
												"input": "$populatedProducts",
												"as":    "p",
												"cond":  bson.M{"$eq": []interface{}{"$$p._id", "$$item.productId"}},
											},
										},
										0,
									},
								},
								"unit": bson.M{
									"$cond": []interface{}{
										bson.M{"$ne": []interface{}{"$$item.unitId", nil}},
										bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$populatedUnits",
														"as":    "u",
														"cond":  bson.M{"$eq": []interface{}{"$$u._id", "$$item.unitId"}},
													},
												},
												0,
											},
										},
										nil,
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"$project": bson.M{
			"populatedProducts": 0,
			"populatedUnits":    0,
		},
	},
}
