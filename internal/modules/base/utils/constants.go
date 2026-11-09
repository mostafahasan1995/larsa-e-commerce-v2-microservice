package base_utils

import "go.mongodb.org/mongo-driver/bson"

func AddSoftDeletedFilter(matchStage bson.M) bson.M {
	softDeleteFilter := bson.M{
		"$or": []bson.M{
			{
				"deletedAt": bson.M{"$exists": false},
			},
			{
				"deletedAt": bson.M{"$exists": true, "$eq": nil},
			},
		},
	}

	if _, exists := matchStage["$and"]; !exists {
		matchStage["$and"] = []bson.M{}
	}
	matchStage["$and"] = append(matchStage["$and"].([]bson.M), softDeleteFilter)
	return matchStage
}
