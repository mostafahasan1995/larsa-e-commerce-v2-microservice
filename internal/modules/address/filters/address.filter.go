package address_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	address_enums "larsa-e-commerce-microservice-v2/internal/modules/address/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressFilters struct {
	filters.BaseFilters
	Name      *string                    `json:"name,omitempty"`
	Type      *address_enums.AddressType `json:"type,omitempty"`
	Country   *string                    `json:"country,omitempty"`
	City      *string                    `json:"city,omitempty"`
	IsDefault *bool                      `json:"isDefault,omitempty"`
	UserId    *primitive.ObjectID        `json:"userId,omitempty"`
}

func NewAddressFilters(query string) (*AddressFilters, error) {
	addressFilters := AddressFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &addressFilters); err != nil {
			return nil, err
		}
	}

	return &addressFilters, nil
}

func (addressFilters *AddressFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if addressFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &addressFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.ADDRESS)

	// Apply specific address filters
	if addressFilters.Name != nil {
		matchStage["name"] = addressFilters.Name
	}

	if addressFilters.Type != nil {
		matchStage["type"] = addressFilters.Type
	}

	if addressFilters.Country != nil {
		matchStage["country"] = addressFilters.Country
	}

	if addressFilters.City != nil {
		matchStage["city"] = addressFilters.City
	}

	if addressFilters.IsDefault != nil {
		matchStage["isDefault"] = addressFilters.IsDefault
	}

	if addressFilters.UserId != nil {
		matchStage["userId"] = addressFilters.UserId
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = filters.BuildPopulatePipeline(addressFilters.PopFields, enums.ADDRESS, pipeline)
		pipeline = filters.BuildSortingPipeline(addressFilters.Sorts, pipeline)
	}

	return pipeline
}
