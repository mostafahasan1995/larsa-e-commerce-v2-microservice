package product_shared

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExtractPricePipeline returns MongoDB aggregation pipeline stages to filter products by price
// based on the ExtractPrice function logic
func ExtractPricePipeline(candidateBranchId *primitive.ObjectID, currencyCode string, unitId *primitive.ObjectID) []bson.M {
	pipeline := []bson.M{
		// Add a field to find the appropriate price based on ExtractPrice logic
		{
			"$addFields": bson.M{
				"price": bson.M{
					"$let": bson.M{
						"vars": bson.M{
							// First try to find price matching the branch ID
							"branchSpecificPrice": bson.M{
								"$cond": []interface{}{
									bson.M{"$ne": []interface{}{candidateBranchId, nil}},
									bson.M{
										"$arrayElemAt": []interface{}{
											bson.M{
												"$filter": bson.M{
													"input": "$prices",
													"cond": bson.M{
														"$and": []bson.M{
															// Check if price has branchIds and it's not null/empty
															{"$ne": []interface{}{"$$this.branchIds", nil}},
															{"$gt": []interface{}{bson.M{"$size": bson.M{"$ifNull": []interface{}{"$$this.branchIds", []interface{}{}}}}, 0}},
															// Check if candidateBranchId exists in price.branchIds
															{
																"$in": []interface{}{
																	candidateBranchId,
																	"$$this.branchIds",
																},
															},
															bson.M{
																"$anyElementTrue": bson.M{
																	"$map": bson.M{
																		"input": bson.M{"$ifNull": []interface{}{"$$this.items", []interface{}{}}},
																		"as":    "item",
																		"in":    filterPriceItemPipeline("item", currencyCode, unitId),
																	},
																},
															},
														},
													},
												},
											},
											0,
										},
									},
									nil,
								},
							},
							// If no branch-specific price, find default price
							"defaultPrice": getDefaultPricePipeline(currencyCode, unitId),
						},
						"in": bson.M{
							"$cond": []interface{}{
								bson.M{"$ne": []interface{}{
									bson.M{"$ifNull": []interface{}{"$$branchSpecificPrice", nil}},
									nil,
								}},
								"$$branchSpecificPrice",
								"$$defaultPrice",
							},
						},
					},
				},
			},
		},

		// Filter out products that don't have a valid selected price
		// {
		// 	"$match": bson.M{
		// 		"price": bson.M{"$ne": nil},
		// 	},
		// },
	}

	return pipeline
}

// getDefaultPricePipeline returns the default price pipeline with currency and unitId filtering
func getDefaultPricePipeline(currencyCode string, unitId *primitive.ObjectID) bson.M {
	return bson.M{
		"$arrayElemAt": []interface{}{
			bson.M{
				"$filter": bson.M{
					"input": "$prices",
					"cond": bson.M{
						"$and": []bson.M{
							bson.M{
								"$ne": []interface{}{
									"$$this.isDefault",
									nil,
								},
							},
							bson.M{
								"$eq": []interface{}{
									"$$this.isDefault",
									true,
								},
							},
							bson.M{
								"$anyElementTrue": bson.M{
									"$map": bson.M{
										"input": bson.M{"$ifNull": []interface{}{"$$this.items", []interface{}{}}},
										"as":    "item",
										"in":    filterPriceItemPipeline("item", currencyCode, unitId),
									},
								},
							},
						},
					},
				},
			},
			0,
		},
	}
}

// ExtractPriceWithCurrencyPipeline returns MongoDB aggregation pipeline stages to filter products by price
// and specific currency code based on the ExtractPrice function logic
func ExtractPriceWithCurrencyPipeline(candidateBranchId *primitive.ObjectID, currencyCode string, unitId *primitive.ObjectID) []bson.M {
	pipeline := ExtractPricePipeline(candidateBranchId, currencyCode, unitId)

	// Add currency filtering stage
	currencyFilter := bson.M{
		"$addFields": bson.M{
			"priceItem": bson.M{
				"$let": bson.M{
					"vars": bson.M{
						"priceWithCurrency": bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$filter": bson.M{
										"input": bson.M{"$ifNull": []interface{}{"$price.items", []interface{}{}}},
										"cond":  filterPriceItemPipeline("this", currencyCode, unitId),
									},
								},
								0,
							},
						},
					},
					"in": "$$priceWithCurrency",
				},
			},
		},
	}

	// Insert currency filter before the final match and project stages
	// newPipeline := make([]bson.M, 0, len(pipeline)+1)
	// newPipeline = append(newPipeline, pipeline[:len(pipeline)-2]...)
	pipeline = append(pipeline, currencyFilter)
	// newPipeline = append(newPipeline, pipeline[len(pipeline)-2:]...)

	// Update the final match to include currency check
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				{"price": bson.M{"$ne": nil}},
				{"priceItem": bson.M{"$ne": nil}},
			},
		},
	})

	return pipeline
}

// ExtractPriceWithCurrencyAndDiscountConversionPipeline returns a complete pipeline that:
// 1. Extracts price with currency
// 2. Converts fixed discounts to percentage
// 3. Can be used for consistent discount filtering
func ExtractPriceWithCurrencyAndDiscountConversionPipeline(candidateBranchId *primitive.ObjectID, currencyCode string, unitId *primitive.ObjectID) []bson.M {
	pipeline := ExtractPriceWithCurrencyPipeline(candidateBranchId, currencyCode, unitId)
	pipeline = append(pipeline, ConvertFixedDiscountToPercentagePipeline(currencyCode)...)
	return pipeline
}

func filterPriceItemPipeline(key string, currencyCode string, unitId *primitive.ObjectID) bson.M {
	return bson.M{
		"$and": []bson.M{
			{
				"$eq": []interface{}{
					"$$" + key + ".currencyCode",
					currencyCode,
				},
			},
			{
				"$cond": []interface{}{
					bson.M{"$ne": []interface{}{unitId, nil}},
					// When unitId is provided, filter by both currency and unitId
					bson.M{
						"$eq": []interface{}{
							"$$" + key + ".unitId",
							unitId,
						},
					},
					// When unitId is not provided, don't filter by unitId
					bson.M{"$eq": []interface{}{true, true}},
				},
			},
		},
	}
}
