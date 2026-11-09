package promotion_filters

import (
	"context"
	"encoding/json"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromotionFilters struct {
	base_filters.BaseFilters
	Name      *string              `json:"name,omitempty"`
	IsActive  *bool                `json:"isActive,omitempty"`
	StartDate *string              `json:"startDate,omitempty"`
	EndDate   *string              `json:"endDate,omitempty"`
	BranchIds []primitive.ObjectID `json:"branchIds,omitempty"`
}

func NewPromotionFilters(query string) (*PromotionFilters, error) {
	promotionFilters := PromotionFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &promotionFilters); err != nil {
			return nil, err
		}
	}

	return &promotionFilters, nil
}

func (promotionFilters *PromotionFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{"deletedAt": bson.M{"$exists": false}}

	// Apply base filters
	baseFilters := &promotionFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PROMOTION)

	// Apply promotion-specific filters
	if promotionFilters.Name != nil && *promotionFilters.Name != "" {
		matchStage["name"] = bson.M{"$regex": *promotionFilters.Name, "$options": "i"}
	}

	if promotionFilters.IsActive != nil {
		matchStage["isActive"] = *promotionFilters.IsActive
	}

	if promotionFilters.StartDate != nil && *promotionFilters.StartDate != "" {
		if startDate, err := time.Parse(time.RFC3339, *promotionFilters.StartDate); err == nil {
			matchStage["startDate"] = bson.M{"$gte": startDate}
		}
	}

	if promotionFilters.EndDate != nil && *promotionFilters.EndDate != "" {
		if endDate, err := time.Parse(time.RFC3339, *promotionFilters.EndDate); err == nil {
			matchStage["endDate"] = bson.M{"$lte": endDate}
		}
	}

	if len(promotionFilters.BranchIds) > 0 {
		matchStage["branchIds"] = bson.M{"$in": promotionFilters.BranchIds}
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchStage,
		})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(promotionFilters.PopFields, enums.PROMOTION, pipeline, ctx)

		pipeline = base_filters.BuildSortingPipeline(promotionFilters.Sorts, pipeline)
	}

	return pipeline
}
