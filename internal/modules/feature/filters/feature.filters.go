package feature_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type FeatureFilters struct {
	base_filters.BaseFilters
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	GroupID   string `json:"groupId"`
	ServiceID string `json:"serviceId"`
	IsActive  bool   `json:"isActive"`
}

func NewFeatureFilters(query string) (*FeatureFilters, error) {
	featureFilters := FeatureFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &featureFilters); err != nil {
			return nil, err
		}
	}
	return &featureFilters, nil
}

func (featureFilters *FeatureFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{}

	// Apply base filters
	baseFilters := &featureFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.FEATURE)

	// Apply feature-specific filters
	if featureFilters.Name != "" {
		matchStage["name"] = bson.M{"$regex": featureFilters.Name, "$options": "i"}
	}
	if featureFilters.Slug != "" {
		matchStage["slug"] = bson.M{"$regex": featureFilters.Slug, "$options": "i"}
	}
	if featureFilters.GroupID != "" {
		matchStage["groupId"] = featureFilters.GroupID
	}
	if featureFilters.ServiceID != "" {
		matchStage["serviceId"] = featureFilters.ServiceID
	}
	if featureFilters.IsActive {
		matchStage["isActive"] = featureFilters.IsActive
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchStage,
		})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(featureFilters.PopFields, enums.FEATURE, pipeline)

		pipeline = base_filters.BuildSortingPipeline(featureFilters.Sorts, pipeline)
	}

	return pipeline
}
