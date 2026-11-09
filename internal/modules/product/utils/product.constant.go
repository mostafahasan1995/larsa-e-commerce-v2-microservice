package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CalculateReviewPipeline = []bson.M{
	{
		"$group": bson.M{
			"_id": "$productId",
			"averageRating": bson.M{
				"$avg": "$rate",
			},
			"totalReviews": bson.M{
				"$sum": 1,
			},
		},
	},
	{
		"$project": bson.M{
			"_id":       0,
			"productId": "$_id",
			"averageRating": bson.M{
				"$round": bson.A{"$averageRating", 2},
			},
			"totalReviews": 1,
		},
	},
}

// ReviewStatsPipeline calculates review statistics for products
var ReviewStatsPipeline = []bson.M{
	{
		"$group": bson.M{
			"_id": "$productId",
			"averageRating": bson.M{
				"$avg": "$rate",
			},
			"totalReviews": bson.M{
				"$sum": 1,
			},
		},
	},
	{
		"$project": bson.M{
			"_id":       0,
			"productId": "$_id",
			"averageRating": bson.M{
				"$round": bson.A{"$averageRating", 2},
			},
			"totalReviews": 1,
		},
	},
}

// ActiveDiscountPipeline adds a boolean field `hasActiveDiscount` set to true when a product has a valid discount
// after price extraction. It normalizes fixed discounts using the currency-specific conversion.
var ActiveDiscountPipeline = func(currencyCode string) []bson.M {
	return []bson.M{
		// Normalize discount to percentage and ensure convertedDiscount exists
		{
			"$addFields": bson.M{
				"convertedDiscount": bson.M{
					"$cond": []interface{}{
						bson.M{"$and": []bson.M{
							{"$ne": []interface{}{"$price", nil}},
							{"$ne": []interface{}{"$price.discount", nil}},
							{"$eq": []interface{}{"$price.discount.type", product_enums.FIXED}},
							{"$eq": []interface{}{"$price.currencyCode", currencyCode}},
						}},
						bson.M{
							"$mergeObjects": []interface{}{
								"$price.discount",
								bson.M{
									"type": "percentage",
									"value": bson.M{
										"$cond": []interface{}{
											bson.M{"$and": []bson.M{
												{"$ne": []interface{}{"$priceItem", nil}},
												{"$gt": []interface{}{bson.M{"$ifNull": []interface{}{"$priceItem.value", 0}}, 0}},
											}},
											bson.M{
												"$multiply": []interface{}{
													bson.M{
														"$divide": []interface{}{
															"$price.discount.value",
															"$priceItem.value",
														},
													},
													100,
												},
											},
											"$price.discount.value",
										},
									},
								},
							},
						},
						bson.M{
							"$cond": []interface{}{
								bson.M{"$and": []bson.M{
									{"$ne": []interface{}{"$price", nil}},
									{"$ne": []interface{}{"$price.discount", nil}},
									{"$eq": []interface{}{"$price.discount.currencyCode", currencyCode}},
								}},
								"$price.discount",
								bson.M{"type": product_enums.PERCENTAGE, "value": 0},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"convertedDiscount": bson.M{"$ifNull": []interface{}{"$convertedDiscount", bson.M{"type": product_enums.PERCENTAGE, "value": 0}}},
			},
		},
		{
			"$addFields": bson.M{
				"hasActiveDiscount": bson.M{
					"$and": []bson.M{
						{"$ne": []interface{}{"$convertedDiscount", nil}},
						{"$gt": []interface{}{bson.M{"$ifNull": []interface{}{"$convertedDiscount.value", 0}}, 0}},
					},
				},
			},
		},
	}
}

// ActiveDiscountMatchStage is a reusable match stage to filter products with active discount (> 0)
var ActiveDiscountMatchStage = bson.M{
	"$match": bson.M{
		"convertedDiscount.value": bson.M{"$gt": 0},
	},
}

// ActivePromotionLookupPipeline looks up active promotions for the current product and flags hasActivePromotion
// Matches promotions where:
// - isActive == true
// - startDate <= now and endDate >= now
// - (allProducts == true) OR (productIds contains this product _id)
// - if branchId provided, promotion's branchIds includes branchId
var ActivePromotionLookupPipeline = func(branchId *primitive.ObjectID) []bson.M {
	innerPipeline := []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{}),
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$isActive", true}},
						bson.M{"$lte": []interface{}{"$startDate", "$$NOW"}},
						bson.M{"$gte": []interface{}{"$endDate", "$$NOW"}},
						bson.M{"$or": []interface{}{
							bson.M{"$eq": []interface{}{"$allProducts", true}},
							bson.M{"$in": []interface{}{"$$productId", "$productIds"}},
						}},
					},
				},
			},
		},
	}

	if branchId != nil {
		innerPipeline = append(innerPipeline, bson.M{
			"$match": bson.M{
				"branchIds": bson.M{"$in": []primitive.ObjectID{*branchId}},
			},
		})
	}

	return []bson.M{
		{
			"$lookup": bson.M{
				"from":     enums.PROMOTION.GetCollectionName(),
				"let":      bson.M{"productId": "$_id"},
				"pipeline": innerPipeline,
				"as":       "activePromotions",
			},
		},
		{
			"$addFields": bson.M{
				"hasActivePromotion": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": bson.M{"$ifNull": []interface{}{"$activePromotions", bson.A{}}}},
						0,
					},
				},
			},
		},
	}
}
