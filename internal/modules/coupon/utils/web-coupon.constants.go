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

type GetMineCouponPipelineArgs struct {
	Code                  *string
	CheckCouponRedemption bool
	CalculateTotalUsages  bool
	AppliedAutomatically  *bool
}

func GetMineCouponPipeline(ctx context.Context, args GetMineCouponPipelineArgs) []bson.M {
	cfg, _ := util.GetReqAppCfg(ctx)

	couponFilter := coupon_filters.CouponFilters{
		Code: args.Code,
		BaseFilters: filters.BaseFilters{
			PopFields: &[]string{"userCoupon"},
		},
	}

	pipeline := couponFilter.BuildPipeline(bson.M{}, ctx, true)

	if args.AppliedAutomatically != nil {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"appliedAutomatically": args.AppliedAutomatically,
			},
		})
	}

	// Add lookup for feedback when relatedTo is FEEDBACK
	pipeline = append(pipeline, []bson.M{
		{
			"$lookup": bson.M{
				"from": enums.FEEDBACK.GetCollectionName(),
				"let": bson.M{
					"userEmail": cfg.User.UserData.Email,
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{
										"$eq": []interface{}{
											"$email",
											"$$userEmail",
										},
									},
								},
							},
						},
					},
				},
				"as": "userFeedback",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$userFeedback",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}...)

	// Add match stage to ensure user has feedback when relatedTo is FEEDBACK
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"$or": []bson.M{
				{
					"usage.relatedTo": bson.M{
						"$ne": coupon_enums.FEEDBACK,
					},
				},
				{
					"$and": []bson.M{
						{
							"usage.relatedTo": coupon_enums.FEEDBACK,
						},
						{
							"userFeedback": bson.M{
								"$exists": true,
								"$ne":     nil,
							},
						},
					},
				},
			},
		},
	})

	if args.CalculateTotalUsages {
		pipeline = append(pipeline, calculateTotalUsagesPipeline(ctx)...)
	}

	if args.CheckCouponRedemption {
		pipeline = append(pipeline, []bson.M{
			{
				"$match": bson.M{
					"$and": []bson.M{
						{
							"$expr": bson.M{
								"$or": []bson.M{
									{
										"$eq": []interface{}{
											bson.M{
												"$ifNull": []interface{}{
													"$totalUsages",
													nil,
												},
											},
											nil,
										},
									},
									{
										"$gt": []interface{}{
											"$totalUsages",
											bson.M{
												"$ifNull": []interface{}{
													"$userCouponTotalUsages",
													0,
												},
											},
										},
									},
								},
							},
						},
						{
							"$expr": bson.M{
								"$or": []bson.M{
									{
										"$eq": []interface{}{
											bson.M{
												"$ifNull": []interface{}{
													"$usage.maxUsagePerCustomer",
													0,
												},
											},
											nil,
										},
									},
									{
										"$gt": []interface{}{
											"$usage.maxUsagePerCustomer",
											bson.M{
												"$ifNull": []interface{}{
													"$userCoupon.totalUsages",
													0,
												},
											},
										},
									},
								},
							},
						},
						{
							"$or": []bson.M{
								{
									"type": "auto_generated",
									"userCoupon": bson.M{
										"$ne":     nil,
										"$exists": true,
									},
								},
								{
									"type":     "manual",
									"isActive": true,
									"$expr": bson.M{
										"$and": []bson.M{
											{"$lt": []interface{}{"$usage.startDate", "$$NOW"}},
											{"$gt": []interface{}{"$usage.endDate", "$$NOW"}},
										},
									},
								},
							},
						},

						/*
						* TODO: Coupon Usage filtration must be implemented here
						* It must be implemented like in `CheckCouponRedemption` method
						 */
					},
				},
			},
		}...)
	}

	return pipeline
}

func calculateTotalUsagesPipeline(ctx context.Context) []bson.M {

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"let": bson.M{
					"code": "$code",
				},
				"from":         enums.USER_COUPON.GetCollectionName(),
				"localField":   "code",
				"foreignField": "couponCode",
				"as":           "userCouponTotalUsages",
				"pipeline": []bson.M{
					{
						"$group": bson.M{
							"_id": "$couponCode",
							"totalUsages": bson.M{
								"$sum": "$totalUsages",
							},
						},
					},
					{
						"$project": bson.M{
							"totalUsages": 1,
						},
					},
				},
			},
		},
		{
			"$unwind": bson.M{
				"preserveNullAndEmptyArrays": true,
				"path":                       "$userCouponTotalUsages",
			},
		},
		{
			"$addFields": bson.M{
				"userCouponTotalUsages": bson.M{
					"$ifNull": []interface{}{
						"$userCouponTotalUsages.totalUsages",
						0,
					},
				},
			},
		},
	}

	return pipeline
}
