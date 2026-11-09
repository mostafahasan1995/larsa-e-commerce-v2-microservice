package payment_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type TransactionFilters struct {
	filters.BaseFilters
}

func NewTransactionFilters(query string) (*TransactionFilters, error) {
	transactionFilters := TransactionFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &transactionFilters); err != nil {
			return nil, err
		}
	}
	return &transactionFilters, nil
}

func (transactionFilters *TransactionFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if transactionFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &transactionFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.TRANSACTION)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(baseFilters.PopFields, enums.TRANSACTION, pipeline)

		pipeline = base_filters.BuildSortingPipeline(baseFilters.Sorts, pipeline)
	}
	return pipeline
}
