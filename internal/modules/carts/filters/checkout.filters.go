package cart_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type CheckoutFilters struct {
	base_filters.BaseFilters
	Code *string `json:"code,omitempty"`
}

func NewCheckoutFilters(query string) (*CheckoutFilters, error) {
	checkoutFilters := CheckoutFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &checkoutFilters); err != nil {
			return nil, err
		}
	}

	return &checkoutFilters, nil
}

func (checkoutFilters *CheckoutFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if checkoutFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &checkoutFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CHECKOUT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(checkoutFilters.PopFields, enums.CHECKOUT, pipeline)

		pipeline = base_filters.BuildSortingPipeline(checkoutFilters.Sorts, pipeline)

	}

	return pipeline
}
