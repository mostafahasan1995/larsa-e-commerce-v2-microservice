package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductFilters represents the filters for a product
type BranchProductFilters struct {
	base_filters.BaseFilters
	ProductId *primitive.ObjectID `json:"productId,omitempty"`
}

func NewBranchProductFilters(query string) (*BranchProductFilters, error) {
	branchProductFilters := BranchProductFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &branchProductFilters); err != nil {
			return nil, err
		}
	}

	return &branchProductFilters, nil
}

// BuildPipeline builds a MongoDB aggregation pipeline based on the provided filters
func (branchProductFilters *BranchProductFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if branchProductFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &branchProductFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PRODUCT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(branchProductFilters.PopFields, enums.PRODUCT, pipeline)

		pipeline = base_filters.BuildSortingPipeline(branchProductFilters.Sorts, pipeline)
	}

	return pipeline
}
