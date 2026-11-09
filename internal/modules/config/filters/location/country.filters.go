package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type CountryFilters struct {
	base_filters.BaseFilters
	Iso3     *string `json:"iso3,omitempty"`
	IsActive *bool   `json:"isActiveisActive,omitempty"`
}

func NewCountryFilters(query string) (*CountryFilters, error) {
	countryFilters := CountryFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &countryFilters); err != nil {
			return nil, err
		}
	}

	return &countryFilters, nil
}

func (countryFilters *CountryFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if countryFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &countryFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.COUNTRY)

	if countryFilters.Iso3 != nil {
		matchStage["iso3"] = *countryFilters.Iso3
	}
	if countryFilters.IsActive != nil {
		matchStage["isActive"] = *countryFilters.IsActive
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(countryFilters.PopFields, enums.COUNTRY, pipeline)

		pipeline = base_filters.BuildSortingPipeline(countryFilters.Sorts, pipeline)

	}

	return pipeline
}
