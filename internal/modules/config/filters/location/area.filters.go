package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AreaFilters struct {
	base_filters.BaseFilters
	CityId   *primitive.ObjectID `json:"cityId,omitempty"`
	CityName *string             `json:"cityName,omitempty"`
	IsActive *bool               `json:"isActive,omitempty"`
}

func NewAreaFilters(query string) (*AreaFilters, error) {
	areaFilters := AreaFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &areaFilters); err != nil {
			return nil, err
		}
	}

	return &areaFilters, nil
}

func (areaFilter *AreaFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if areaFilter == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := areaFilter.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.AREA)

	if areaFilter.CityId != nil {
		matchStage["cityId"] = *areaFilter.CityId
	}

	if areaFilter.CityName != nil {
		matchStage["cityName"] = *areaFilter.CityName
	}

	if areaFilter.IsActive != nil {
		matchStage["isActive"] = *areaFilter.IsActive
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})
	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(areaFilter.PopFields, enums.AREA, pipeline)
		pipeline = base_filters.BuildSortingPipeline(areaFilter.Sorts, pipeline)
	}

	return pipeline
}
