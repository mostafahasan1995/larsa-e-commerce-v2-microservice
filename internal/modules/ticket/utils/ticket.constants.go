package ticket_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"

	"go.mongodb.org/mongo-driver/bson"
)

var TicketPopulateCommentsUsers = []bson.M{
	{
		"$lookup": bson.M{
			"from": enums.PROFILE.GetCollectionName(),
			"let": bson.M{
				"commentUserIds": bson.M{
					"$ifNull": []interface{}{
						"$comments.userId",
						[]interface{}{},
					},
				},
			},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$in": []interface{}{"$_id", "$$commentUserIds"},
						},
					},
				},
			},
			"as": "commentUsers",
		},
	},
	{
		"$addFields": bson.M{
			"comments": bson.M{
				"$map": bson.M{
					"input": bson.M{
						"$ifNull": []interface{}{"$comments", []interface{}{}},
					},
					"as": "comment",
					"in": bson.M{
						"text":       "$$comment.text",
						"attachment": "$$comment.attachment",
						"userId":     "$$comment.userId",
						"createdAt":  "$$comment.createdAt",
						"user": bson.M{
							"$arrayElemAt": []interface{}{
								bson.M{
									"$filter": bson.M{
										"input": "$commentUsers",
										"as":    "u",
										"cond":  bson.M{"$eq": []interface{}{"$$u._id", "$$comment.userId"}},
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
	{
		"$project": bson.M{
			"commentUsers": 0,
		},
	},
}
