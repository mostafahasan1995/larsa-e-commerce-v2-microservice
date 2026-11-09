package order_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var PopulateOrder = []bson.M{
	{
		"$lookup": bson.M{
			"from":         "quotations",
			"localField":   "quotationId",
			"foreignField": "_id",
			"as":           "quotation",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$quotation",
			"preserveNullAndEmptyArrays": true,
		},
	},
}

// OptimizedPopulateOrderItemsProducts - Optimized version that only fetches essential product fields
var OptimizedPopulateOrderItemsProducts = []bson.M{
	{
		"$lookup": bson.M{
			"from": enums.PRODUCT.GetCollectionName(),
			"let": bson.M{
				"productItemIds": bson.M{
					"$ifNull": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": bson.M{
									"$filter": bson.M{
										"input": bson.M{"$ifNull": []interface{}{"$orderItems", []interface{}{}}},
										"as":    "item",
										"cond": bson.M{
											"$eq": []interface{}{"$$item.type", "product"},
										},
									},
								},
								"as": "item",
								"in": "$$item.itemId",
							},
						},
						[]interface{}{},
					},
				},
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []interface{}{
								bson.M{"$gt": []interface{}{bson.M{"$size": "$$productItemIds"}, 0}},
								bson.M{"$in": []interface{}{"$_id", "$$productItemIds"}},
							},
						},
					},
				},
				// Project to exclude heavy fields only
				{
					"$project": bson.M{
						// Exclude heavy fields only - everything else included by default
						"createdAt":  0,
						"updatedAt":  0,
						"metadata":   0,
						"reviews":    0,
						"attributes": 0,
					},
				},
				// Limit images to first 3 in separate stage to avoid projection conflicts
				{
					"$addFields": bson.M{
						"images": bson.M{"$slice": []interface{}{"$images", 3}},
					},
				},
			},
			"as": "productLookup",
		},
	},
	{
		"$addFields": bson.M{
			"orderItems": bson.M{
				"$map": bson.M{
					"input": bson.M{"$ifNull": []interface{}{"$orderItems", []interface{}{}}},
					"as":    "item",
					"in": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": []interface{}{"$$item.type", "product"},
							},
							"then": bson.M{
								"$mergeObjects": []interface{}{
									"$$item",
									bson.M{
										"product": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$productLookup",
														"as":    "p",
														"cond": bson.M{
															"$eq": []interface{}{"$$p._id", "$$item.itemId"},
														},
													},
												},
												0,
											},
										},
									},
								},
							},
							"else": "$$item",
						},
					},
				},
			},
		},
	},
	{
		"$project": bson.M{
			"productLookup": 0,
		},
	},
}

var PopulateOrderItemsProducts = []bson.M{
	{
		"$lookup": bson.M{
			"from": enums.PRODUCT.GetCollectionName(),
			"let": bson.M{
				"productItemIds": bson.M{
					"$ifNull": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": bson.M{
									"$filter": bson.M{
										"input": bson.M{"$ifNull": []interface{}{"$orderItems", []interface{}{}}},
										"as":    "item",
										"cond": bson.M{
											"$eq": []interface{}{"$$item.type", "product"},
										},
									},
								},
								"as": "item",
								"in": "$$item.itemId",
							},
						},
						[]interface{}{},
					},
				},
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []interface{}{
								bson.M{"$gt": []interface{}{bson.M{"$size": "$$productItemIds"}, 0}},
								bson.M{"$in": []interface{}{"$_id", "$$productItemIds"}},
							},
						},
					},
				},
			},
			"as": "productLookup",
		},
	},
	{
		"$addFields": bson.M{
			"orderItems": bson.M{
				"$map": bson.M{
					"input": bson.M{"$ifNull": []interface{}{"$orderItems", []interface{}{}}},
					"as":    "item",
					"in": bson.M{
						"$cond": bson.M{
							"if": bson.M{
								"$eq": []interface{}{"$$item.type", "product"},
							},
							"then": bson.M{
								"$mergeObjects": []interface{}{
									"$$item",
									bson.M{
										"product": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$productLookup",
														"as":    "p",
														"cond": bson.M{
															"$eq": []interface{}{"$$p._id", "$$item.itemId"},
														},
													},
												},
												0,
											},
										},
									},
								},
							},
							"else": "$$item",
						},
					},
				},
			},
		},
	},
	{
		"$project": bson.M{
			"productLookup": 0,
		},
	},
}

var CalculateTotalPricePipeline = []bson.M{
	{
		"$addFields": bson.M{
			"totalPrice": bson.M{
				"$sum": bson.M{
					"$map": bson.M{
						"input": "$products",
						"as":    "product",
						"in": bson.M{
							"$multiply": []string{"$$product.quantity", "$$product.pricePerUnit"},
						},
					},
				},
			},
		},
	},
}
