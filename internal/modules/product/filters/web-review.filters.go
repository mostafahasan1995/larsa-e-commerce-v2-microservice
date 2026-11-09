package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebReviewFilters struct {
	base_filters.BaseFilters
	ProductId *primitive.ObjectID `json:"productId,omitempty"`
}

func NewWebReviewFilters(query string) (*WebReviewFilters, error) {
	reviewFilters := WebReviewFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &reviewFilters); err != nil {
			return nil, err
		}
	}

	return &reviewFilters, nil
}

func (reviewFilters *WebReviewFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if reviewFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	if reviewFilters.ProductId != nil {
		matchStage["productId"] = reviewFilters.ProductId
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &reviewFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.REVIEW)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(reviewFilters.PopFields, enums.REVIEW, pipeline)

		pipeline = base_filters.BuildSortingPipeline(reviewFilters.Sorts, pipeline)
	}

	return pipeline
}
