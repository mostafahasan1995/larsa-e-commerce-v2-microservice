package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CityBranchFilters struct {
	Size *base_utils.NumericFilter `json:"size,omitempty"`
}

type CityFilters struct {
	base_filters.BaseFilters
	CountryId   *primitive.ObjectID `json:"countryId,omitempty"`
	CountryName *string             `json:"countryName,omitempty"`
	IsActive    *bool               `json:"isActive,omitempty"`
	Branches    *CityBranchFilters  `json:"branches,omitempty"`
}

func NewCityFilters(query string) (*CityFilters, error) {
	cityFilters := CityFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &cityFilters); err != nil {
			return nil, err
		}
	}

	return &cityFilters, nil
}

func (cityFilters *CityFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if cityFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &cityFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CITY)

	if cityFilters.CountryId != nil {
		matchStage["countryId"] = *cityFilters.CountryId
	}

	if cityFilters.CountryName != nil {
		matchStage["countryName"] = *cityFilters.CountryName
	}

	if cityFilters.IsActive != nil {
		matchStage["isActive"] = *cityFilters.IsActive
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(cityFilters.PopFields, enums.CITY, pipeline)

		if cityFilters.PopFields != nil && util.Contains(*cityFilters.PopFields, "branches") && cityFilters.Branches != nil && cityFilters.Branches.Size != nil {
			var operator string
			var value int

			if cityFilters.Branches.Size.Gte != nil {
				operator = "$gte"
				value = *cityFilters.Branches.Size.Gte
			}

			if cityFilters.Branches.Size.Lte != nil {
				operator = "$lte"
				value = *cityFilters.Branches.Size.Lte
			}

			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"$expr": bson.M{
						operator: []interface{}{bson.M{"$size": "$branches"}, value},
					},
				},
			})
		}

		pipeline = base_filters.BuildSortingPipeline(cityFilters.Sorts, pipeline)
	}

	return pipeline
}
