package feature_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var PopulateFeaturesCount = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.FEATURE.GetCollectionName(),
			"localField":   "_id",
			"foreignField": "groupId",
			"as":           "features",
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"deletedAt": bson.M{"$exists": false},
					},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"featuresCount": bson.M{"$size": "$features"},
		},
	},
}
