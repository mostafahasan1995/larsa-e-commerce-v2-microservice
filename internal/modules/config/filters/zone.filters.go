package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type ZoneFilters struct {
	base_filters.BaseFilters
	Code       *string   `json:"3,omitempty"`
	CountryIds *[]string `json:"countryIds,omitempty"`
	CityIds    *[]string `json:"cityIds,omitempty"`
	AreaIds    *[]string `json:"areaIds,omitempty"`
}

func NewZoneFilters(query string) (*ZoneFilters, error) {
	zoneFilters := ZoneFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &zoneFilters); err != nil {
			return nil, err
		}
	}

	return &zoneFilters, nil
}

func (zoneFilters *ZoneFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if zoneFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &zoneFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.ZONE)

	if zoneFilters.Code != nil {
		matchStage["code"] = *zoneFilters.Code
	}

	if zoneFilters.CountryIds != nil && len(*zoneFilters.CountryIds) > 0 {
		matchStage["countryIds"] = bson.M{"$in": zoneFilters.CountryIds}
	}

	if zoneFilters.CityIds != nil && len(*zoneFilters.CityIds) > 0 {
		matchStage["cityIds"] = bson.M{"$in": zoneFilters.CityIds}
	}

	if zoneFilters.AreaIds != nil && len(*zoneFilters.AreaIds) > 0 {
		matchStage["areaIds"] = bson.M{"$in": zoneFilters.AreaIds}
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(zoneFilters.PopFields, enums.ZONE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(zoneFilters.Sorts, pipeline)
	}

	return pipeline
}
