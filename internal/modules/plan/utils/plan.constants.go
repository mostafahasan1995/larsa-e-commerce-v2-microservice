package plan_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var PopulatePlanFeatures = []bson.M{
	{
		"$unwind": bson.M{
			"path":                       "$features",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$lookup": bson.M{
			"from":         enums.FEATURE.GetCollectionName(),
			"localField":   "features.featureId",
			"foreignField": "_id",
			"as":           "features.featureObject",
		},
	},
	{
		"$addFields": bson.M{
			"features.featureObject": bson.M{
				"$arrayElemAt": []interface{}{"$features.featureObject", 0},
			},
		},
	},
	{
		"$group": bson.M{
			"_id": "$_id",
			"name": bson.M{
				"$first": "$name",
			},
			"type": bson.M{
				"$first": "$type",
			},
			"trailDays": bson.M{
				"$first": "$trailDays",
			},
			"isDefault": bson.M{
				"$first": "$isDefault",
			},
			"isActive": bson.M{
				"$first": "$isActive",
			},
			"expirationDays": bson.M{
				"$first": "$expirationDays",
			},
			"suspendDays": bson.M{
				"$first": "$suspendDays",
			},
			"description": bson.M{
				"$first": "$description",
			},
			"image": bson.M{
				"$first": "$image",
			},
			"monthlyPrices": bson.M{
				"$first": "$monthlyPrices",
			},
			"annualPrices": bson.M{
				"$first": "$annualPrices",
			},
			"features": bson.M{
				"$push": "$features",
			},
			"customizations": bson.M{
				"$first": "$customizations",
			},
			"createdAt": bson.M{
				"$first": "$createdAt",
			},
			"updatedAt": bson.M{
				"$first": "$updatedAt",
			},
			"createdBy": bson.M{
				"$first": "$createdBy",
			},
			"updatedBy": bson.M{
				"$first": "$updatedBy",
			},
		},
	},
}
