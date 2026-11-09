package feature_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var PopulateFeatureGroup = []bson.M{
	{
		"$lookup": bson.M{
			"from":         enums.FEATURE_GROUP.GetCollectionName(),
			"localField":   "groupId",
			"foreignField": "_id",
			"as":           "group",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$group",
			"preserveNullAndEmptyArrays": true,
		},
	},
}
