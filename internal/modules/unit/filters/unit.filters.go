package unit_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type UnitFilters struct {
	base_filters.BaseFilters
	Name       *types.Localizable[string] `json:"name,omitempty"`
	Type       *string                    `json:"type,omitempty"`
	CategoryId *string                    `json:"categoryId,omitempty"`
	Rate       *int                       `json:"rate,omitempty"`
	Active     *bool                      `json:"active,omitempty"`
}

func NewUnitFilters(query string) (*UnitFilters, error) {
	unitFilters := UnitFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &unitFilters); err != nil {
			return nil, err
		}
	}

	return &unitFilters, nil
}

func (unitFilters *UnitFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if unitFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &unitFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.UNIT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// Apply type filter
	if unitFilters.Type != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": *unitFilters.Type}})
	}

	// Apply categoryId filter
	if unitFilters.CategoryId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"categoryId": *unitFilters.CategoryId}})
	}

	// Apply rate filter
	if unitFilters.Rate != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"rate": *unitFilters.Rate}})
	}

	// Apply active filter
	if unitFilters.Active != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"active": *unitFilters.Active}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(unitFilters.PopFields, enums.UNIT, pipeline)

		pipeline = base_filters.BuildSortingPipeline(unitFilters.Sorts, pipeline)
	}

	return pipeline
}
