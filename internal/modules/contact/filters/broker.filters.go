package contact_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type BrokerFilters struct {
	base_filters.BaseFilters
	FullName                *string                    `json:"fullName,omitempty"`
	CompanyName             *types.Localizable[string] `json:"companyName,omitempty"`
	CustomsLicenseNumber    *string                    `json:"customsLicenseNumber,omitempty"`
	CommercialRegNumber     *string                    `json:"commercialRegNumber,omitempty"`
	TaxIdentificationNumber *string                    `json:"taxIdentificationNumber,omitempty"`
	AuthorizedCountry       *string                    `json:"authorizedCountry,omitempty"`
}

func NewBrokerFilters(query string) (*BrokerFilters, error) {
	brokerFilters := BrokerFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &brokerFilters); err != nil {
			return nil, err
		}
	}
	return &brokerFilters, nil
}

func (brokerFilters *BrokerFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if brokerFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &brokerFilters.BaseFilters

	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.BROKER)
	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// Apply fullName filter
	if brokerFilters.FullName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"fullName": *brokerFilters.FullName}})
	}

	if brokerFilters.AuthorizedCountry != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"authorizedCountry": *brokerFilters.AuthorizedCountry}})
	}

	if addExtraOperations {

		pipeline = base_filters.BuildPopulatePipeline(baseFilters.PopFields, enums.BROKER, pipeline)

		pipeline = base_filters.BuildSortingPipeline(baseFilters.Sorts, pipeline)
	}

	return pipeline
}
