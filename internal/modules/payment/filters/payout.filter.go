package payment_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type PayoutFilters struct {
	base_filters.BaseFilters
	Amount *float64 `json:"amount,omitempty"`
}

func NewPayoutFilters(query string) (*PayoutFilters, error) {
	payoutFilters := PayoutFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &payoutFilters); err != nil {
			return nil, err
		}
	}

	return &payoutFilters, nil
}

func (payoutFilters *PayoutFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if payoutFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &payoutFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PAYOUT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(payoutFilters.PopFields, enums.PAYMENT, pipeline)
		pipeline = base_filters.BuildSortingPipeline(payoutFilters.Sorts, pipeline)
	}

	return pipeline
}
