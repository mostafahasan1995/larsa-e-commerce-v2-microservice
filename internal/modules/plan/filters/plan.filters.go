package plan_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	plan_enums "larsa-e-commerce-microservice-v2/internal/modules/plan/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type PlanFilters struct {
	base_filters.BaseFilters
	Type plan_enums.PlanType `json:"type"`
}

func NewPlanFilters(query string) (*PlanFilters, error) {
	planFilters := PlanFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &planFilters); err != nil {
			return nil, err
		}
	}
	return &planFilters, nil
}

func (planFilters *PlanFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{"deletedAt": bson.M{"$exists": false}}

	// Apply base filters
	baseFilters := &planFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PLAN)

	// Apply plan-specific filters
	if planFilters.Type != "" {
		matchStage["type"] = planFilters.Type
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchStage,
		})
	}

	if addExtraOperations {
		// pipeline = base_filters.BuildPopulatePipeline(planFilters.PopFields, enums.PLAN, pipeline)

		pipeline = base_filters.BuildSortingPipeline(planFilters.Sorts, pipeline)

	}
	return pipeline
}
