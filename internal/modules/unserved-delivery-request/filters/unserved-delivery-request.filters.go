package unserved_delivery_request_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type UnservedDeliveryRequestFilters struct {
	filters.BaseFilters
	Country *string `json:"country,omitempty"`
	City    *string `json:"city,omitempty"`
	Area    *string `json:"area,omitempty"`
}

func NewUnservedDeliveryRequestFilters(query string) (*UnservedDeliveryRequestFilters, error) {
	f := UnservedDeliveryRequestFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &f); err != nil {
			return nil, err
		}
	}
	return &f, nil
}

func (f *UnservedDeliveryRequestFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if f == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)
	baseFilters := &f.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.UNSERVED_DELIVERY_REQUEST)

	if f.Country != nil && *f.Country != "" {
		matchStage["country"] = bson.M{"$regex": *f.Country, "$options": "i"}
	}
	if f.City != nil && *f.City != "" {
		matchStage["city"] = bson.M{"$regex": *f.City, "$options": "i"}
	}
	if f.Area != nil && *f.Area != "" {
		matchStage["area"] = bson.M{"$regex": *f.Area, "$options": "i"}
	}

	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = filters.BuildPopulatePipeline(f.PopFields, enums.UNSERVED_DELIVERY_REQUEST, pipeline)
		pipeline = filters.BuildSortingPipeline(f.Sorts, pipeline)
	}

	return pipeline
}
