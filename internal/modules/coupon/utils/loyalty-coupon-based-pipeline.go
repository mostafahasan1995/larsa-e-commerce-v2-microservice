package coupon_utils

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_filters "larsa-e-commerce-microservice-v2/internal/modules/coupon/filters"

	"go.mongodb.org/mongo-driver/bson"
)

func GetLoyaltyCouponBasedPipeline(ctx context.Context) []bson.M {
	cfg, _ := util.GetReqAppCfg(ctx)
	couponFilter := coupon_filters.CouponFilters{
		BaseFilters: filters.BaseFilters{},
	}

	pipeline := couponFilter.BuildPipeline(bson.M{}, ctx, true)

	orderPipeline := []bson.M{
		{
			"$match": bson.M{
				"buyerInfo.id": bson.M{"$eq": cfg.User.Id},
			},
		},
		GroupByBuyerIdPipeline,
	}
	orderPipeline = append(orderPipeline, UnwindTotalOrders...)

	facetPipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": []bson.M{
					{
						"appliedAutomatically": true,
					},
					{
						"type": "auto_generated",
					},
					{
						"$expr": bson.M{
							"$and": []bson.M{
								{
									"$lt": []interface{}{"$rule.startDate", "$$NOW"},
								},
								{
									"$gt": []interface{}{"$rule.endDate", "$$NOW"},
								},
							},
						},
					},
					{
						"rule.conditions.values": bson.M{
							"$elemMatch": bson.M{
								"type": coupon_enums.NUMBER_OF_ORDERS,
							},
						},
					},
				},
			},
		},

		{
			"$lookup": bson.M{
				"from": enums.ORDER.GetCollectionName(),
				// "localField":   "code",
				// "foreignField": "couponCode",
				"as":       "orders",
				"pipeline": orderPipeline,
			},
		},
		{
			"$addFields": bson.M{
				"totalOrders": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$isArray": "$orders.totalOrders"},
						"then": bson.M{
							"$arrayElemAt": []interface{}{"$orders.totalOrders", 0},
						},
						"else": 0,
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"numberOfOrders": bson.M{
					"$arrayElemAt": []interface{}{
						"$rule.conditions.values.value",
						bson.M{
							"$indexOfArray": []interface{}{
								"$rule.conditions.values.type",
								coupon_enums.NUMBER_OF_ORDERS,
							},
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{
							"$in": []interface{}{
								coupon_enums.NUMBER_OF_ORDERS,
								"$rule.conditions.values.type",
							},
						},
						{
							"$lt": []interface{}{
								"$totalOrders",
								"$numberOfOrders",
							},
						},
					},
				},
			},
		},
	}

	pipeline = append(pipeline, facetPipeline...)

	return pipeline
}
