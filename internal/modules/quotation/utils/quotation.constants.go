package quotation_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var QuotationCalculateSubmittedQuotationCount = []bson.M{
	{
		"$lookup": bson.M{
			"from": enums.QUOTATION.GetCollectionName(),
			"let":  bson.M{"requestQuotationId": "$_id"},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$requestQuotationId", "$$requestQuotationId"}},
								{"$eq": []interface{}{"$status", quotation_enums.SENT_TO_BUYER}},
							},
						},
					},
				},
				{
					"$count": "count",
				},
			},
			"as": "submittedQuotations",
		},
	},
	{
		"$addFields": bson.M{
			"submittedQuotationCount": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$gt": []interface{}{bson.M{"$size": "$submittedQuotations"}, 0}},
					"then": bson.M{"$arrayElemAt": []interface{}{"$submittedQuotations.count", 0}},
					"else": 0,
				},
			},
		},
	},
	{
		"$project": bson.M{
			"submittedQuotations": 0,
		},
	},
}

var QuotationPopulateProducts = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.PRODUCT.GetCollectionName(),
			"localField":   "products.productId",
			"foreignField": "_id",
			"as":           "populatedProducts",
		},
	},
	{
		"$addFields": bson.M{
			"products": bson.M{
				"$map": bson.M{
					"input": "$products",
					"as":    "product",
					"in": bson.M{
						"$mergeObjects": []interface{}{
							"$$product",
							bson.M{
								"product": bson.M{
									"$arrayElemAt": []interface{}{
										bson.M{
											"$filter": bson.M{
												"input": "$populatedProducts",
												"as":    "p",
												"cond":  bson.M{"$eq": []interface{}{"$$p._id", "$$product.productId"}},
											},
										},
										0,
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
		},
	},
}
