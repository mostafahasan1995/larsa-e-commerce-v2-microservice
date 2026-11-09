package cart_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var PopulateCartItemsProductsPipeline = []bson.M{

	{
		"$lookup": bson.M{
			"from":         enums.CART_ITEM.GetCollectionName(),
			"localField":   "cartItemIds",
			"foreignField": "_id",
			"as":           "items",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$items",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$lookup": bson.M{
			"from":         enums.PRODUCT.GetCollectionName(),
			"localField":   "items.itemId",
			"foreignField": "_id",
			"as":           "items.product",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$items.product",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$group": bson.M{
			"_id":             "$_id",
			"cartItemIds":     bson.M{"$first": "$cartItemIds"},
			"userId":          bson.M{"$first": "$userId"},
			"totalQuantities": bson.M{"$first": "$totalQuantities"},
			"items":           bson.M{"$push": "$items"},
		},
	},
}
