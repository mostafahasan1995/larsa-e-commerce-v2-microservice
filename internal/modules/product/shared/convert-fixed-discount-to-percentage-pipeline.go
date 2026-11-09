package product_shared

import (
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson"
)

// ConvertFixedDiscountToPercentagePipeline returns MongoDB aggregation pipeline stages
// to convert fixed discount values to percentage values for filtering purposes
func ConvertFixedDiscountToPercentagePipeline(currencyCode string) []bson.M {
	return []bson.M{
		{
			"$addFields": bson.M{
				"convertedDiscount": bson.M{
					"$cond": []interface{}{
						bson.M{"$and": []bson.M{
							bson.M{"$ne": []interface{}{"$price", nil}},
							bson.M{"$ne": []interface{}{"$price.discount", nil}},
							bson.M{"$eq": []interface{}{"$price.discount.type", product_enums.FIXED}},
							bson.M{"$eq": []interface{}{"$price.currencyCode", currencyCode}},
						}},
						bson.M{
							"$mergeObjects": []interface{}{
								"$price.discount",
								bson.M{
									"type": "percentage",
									"value": bson.M{
										"$cond": []interface{}{
											bson.M{"$and": []bson.M{
												bson.M{"$ne": []interface{}{"$priceItem", nil}},
												bson.M{"$gt": []interface{}{bson.M{"$ifNull": []interface{}{"$priceItem.value", 0}}, 0}},
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
											"$price.discount.value", // Keep original value if priceItem is 0 or null
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
									{
										"$eq": []interface{}{
											"$price.discount.currencyCode",
											currencyCode,
										},
									},
								}},
								"$price.discount", // Keep original discount if no fixed discount
								bson.M{ // Set to 0 when there's no discount or no price
									"type":  "percentage",
									"value": 0,
								},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"convertedDiscount": bson.M{
					"$ifNull": []interface{}{
						"$convertedDiscount",
						bson.M{
							"type":  product_enums.PERCENTAGE,
							"value": 0,
						},
					},
				},
			},
		},
	}
}
