package feature_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	feature_utils "larsa-e-commerce-microservice-v2/internal/modules/feature/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type FeatureGroupFilters struct {
	base_filters.BaseFilters
	Name string `json:"name"`
}

func NewFeatureGroupFilters(query string) (*FeatureGroupFilters, error) {
	featureGroupFilters := FeatureGroupFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &featureGroupFilters); err != nil {
			return nil, err
		}
	}
	return &featureGroupFilters, nil
}

func (featureGroupFilters *FeatureGroupFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{}

	// Apply base filters
	baseFilters := &featureGroupFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.FEATURE_GROUP)

	if featureGroupFilters.Name != "" {
		matchStage["name"] = bson.M{"$regex": featureGroupFilters.Name, "$options": "i"}
	}

	// Wrap the matchStage in a $match stage
	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{"$match": matchStage})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(featureGroupFilters.PopFields, enums.FEATURE_GROUP, pipeline)

		pipeline = base_filters.BuildSortingPipeline(featureGroupFilters.Sorts, pipeline)

		pipeline = append(pipeline, feature_utils.PopulateFeaturesCount...)
	}
	return pipeline
}
