package contact_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type ShippingFilters struct {
	base_filters.BaseFilters
	Name    *types.Localizable[string] `json:"name,omitempty"`
	Code    *string                    `json:"code,omitempty"`
	Country *string                    `json:"country,omitempty"`
}

func NewShippingFilters(query string) (*ShippingFilters, error) {
	ShippingFilters := ShippingFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &ShippingFilters); err != nil {
			return nil, err
		}
	}
	return &ShippingFilters, nil
}

func (filters *ShippingFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if filters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &filters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.SHIPPING_COMPANY)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// Apply name filter
	if filters.Name != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"name": filters.Name}})
	}

	// Apply code filter
	if filters.Code != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"code": *filters.Code}})
	}

	// Apply country filter
	if filters.Country != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"country": *filters.Country}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(filters.PopFields, enums.SHIPPING_COMPANY, pipeline)

		pipeline = base_filters.BuildSortingPipeline(baseFilters.Sorts, pipeline)
	}

	return pipeline
}
