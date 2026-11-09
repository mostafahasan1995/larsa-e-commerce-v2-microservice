package coupon_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type CouponFilters struct {
	base_filters.BaseFilters
	Code  *string                    `json:"code,omitempty"`
	Types *[]coupon_enums.CouponType `json:"types,omitempty"`
}

func NewCouponFilters(query string) (*CouponFilters, error) {
	couponFilters := CouponFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &couponFilters); err != nil {
			return nil, err
		}
	}

	return &couponFilters, nil
}

func (couponFilters *CouponFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if couponFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &couponFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.COUPON)

	if couponFilters.Code != nil {
		matchStage["code"] = *couponFilters.Code
	}

	if couponFilters.Types != nil && len(*couponFilters.Types) > 0 {
		matchStage["type"] = bson.M{"$in": *couponFilters.Types}
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		cfg, _ := util.GetReqAppCfg(ctx)
		if cfg != nil && cfg.User != nil {
			pipeline = append(pipeline, bson.M{
				"$addFields": bson.M{
					"userId": cfg.User.Id,
				},
			})
		}
		pipeline = base_filters.BuildPopulatePipeline(couponFilters.PopFields, enums.COUPON, pipeline)

		pipeline = base_filters.BuildSortingPipeline(couponFilters.Sorts, pipeline)

	}

	return pipeline
}
