package cart_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItemFilters struct {
	base_filters.BaseFilters
	Code     *string             `json:"code,omitempty"`
	BranchId *primitive.ObjectID `json:"branchId,omitempty"`
	VendorId *primitive.ObjectID `json:"vendorId,omitempty"`
	CartId   *primitive.ObjectID `json:"cartId,omitempty"`
}

func NewCartItemFilters(query string) (*CartItemFilters, error) {
	cartItemFilters := CartItemFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &cartItemFilters); err != nil {
			return nil, err
		}
	}

	return &cartItemFilters, nil
}

func (cartItemFilters *CartItemFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if cartItemFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &cartItemFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CART_ITEM)

	if cartItemFilters.BranchId != nil {
		matchStage["branchId"] = cartItemFilters.BranchId
	}

	if cartItemFilters.VendorId != nil {
		matchStage["vendorId"] = cartItemFilters.VendorId
	}

	if cartItemFilters.CartId != nil {
		matchStage["cartId"] = cartItemFilters.CartId
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(cartItemFilters.PopFields, enums.CART_ITEM, pipeline)

		pipeline = base_filters.BuildSortingPipeline(cartItemFilters.Sorts, pipeline)

	}

	return pipeline
}
