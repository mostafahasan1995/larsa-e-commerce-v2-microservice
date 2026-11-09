package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type GroupFilters struct {
	base_filters.BaseFilters
}

func NewGroupFilters(query string) (*GroupFilters, error) {
	groupFilters := GroupFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &groupFilters); err != nil {
			return nil, err
		}
	}

	return &groupFilters, nil
}

func (groupFilters *GroupFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if groupFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &groupFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.GROUP)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(groupFilters.PopFields, enums.GROUP, pipeline)

		pipeline = base_filters.BuildSortingPipeline(groupFilters.Sorts, pipeline)
	}

	return pipeline
}
