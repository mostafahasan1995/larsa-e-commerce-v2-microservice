package cart_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type CartFilters struct {
	base_filters.BaseFilters
	CartItem *CartItemFilters `json:"cartItem,omitempty"`
}

func NewCartFilters(query string) (*CartFilters, error) {
	cartFilters := CartFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &cartFilters); err != nil {
			return nil, err
		}
	}
	return &cartFilters, nil
}

func (cartFilters *CartFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if cartFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &cartFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CART)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {

		if cartFilters.CartItem != nil {
			pipeline = append(pipeline, bson.M{"$addFields": bson.M{"vendorId": cartFilters.CartItem.VendorId}})
			pipeline = append(pipeline, bson.M{"$addFields": bson.M{"branchId": cartFilters.CartItem.BranchId}})
		}

		pipeline = base_filters.BuildPopulatePipeline(cartFilters.PopFields, enums.CART, pipeline)

		pipeline = base_filters.BuildSortingPipeline(cartFilters.Sorts, pipeline)

	}

	return pipeline
}
