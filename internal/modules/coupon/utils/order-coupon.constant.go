package coupon_utils

import "go.mongodb.org/mongo-driver/bson"

var GroupByBuyerIdPipeline = bson.M{
	"$facet": bson.M{
		"numberOfOrders": []bson.M{
			{
				"$group": bson.M{
					"_id":         "$buyerInfo.id",
					"totalOrders": bson.M{"$sum": 1},
					"userId":      bson.M{"$first": "$buyerInfo.id"},
				},
			},
		},
	},
}

var UnwindTotalOrders = []bson.M{
	{
		"$unwind": "$numberOfOrders",
	},
	{
		"$replaceRoot": bson.M{
			"newRoot": "$numberOfOrders",
		},
	},
	{
		"$project": bson.M{
			"userId":      1,
			"totalOrders": 1,
			"_id":         0,
		},
	},
}

func NumberOfOrdersConditionPipeline(operator string, conditionValue interface{}) []bson.M {
	pipeline := []bson.M{
		GroupByBuyerIdPipeline,
		{
			"$match": bson.M{
				"numberOfOrders.totalOrders": bson.M{operator: conditionValue},
			},
		},
	}
	return append(pipeline, UnwindTotalOrders...)
}
