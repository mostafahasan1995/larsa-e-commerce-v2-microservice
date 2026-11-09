package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type CurrencyFilters struct {
	base_filters.BaseFilters
	Code *string `json:"code,omitempty"`
}

func NewCurrencyFilters(query string) (*CurrencyFilters, error) {
	currencyFilters := CurrencyFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &currencyFilters); err != nil {
			return nil, err
		}
	}

	return &currencyFilters, nil
}

func (currencyFilters *CurrencyFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if currencyFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &currencyFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CURRENCY)

	if currencyFilters.Code != nil {
		matchStage["code"] = *currencyFilters.Code
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(currencyFilters.PopFields, enums.CURRENCY, pipeline)

		pipeline = base_filters.BuildSortingPipeline(currencyFilters.Sorts, pipeline)

	}

	return pipeline
}
