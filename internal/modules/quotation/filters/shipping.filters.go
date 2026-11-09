package quotation_filters

import (
	"context"
	"encoding/json"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type ShippingFilters struct {
	base_filters.BaseFilters
	Type *quotation_enums.ShippingType `json:"type,omitempty"`
}

func NewShippingFilters(query string) (*ShippingFilters, error) {
	shippingFilters := ShippingFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &shippingFilters); err != nil {
			return nil, err
		}
	}
	return &shippingFilters, nil
}

func (shippingFilters *ShippingFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if shippingFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &shippingFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.SHIPPING)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if shippingFilters.Type != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": shippingFilters.Type}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(shippingFilters.PopFields, enums.SHIPPING, pipeline)
		pipeline = base_filters.BuildSortingPipeline(shippingFilters.Sorts, pipeline)
	}

	return pipeline
}
