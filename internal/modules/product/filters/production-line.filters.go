package product_filters

import (
	"context"
	"encoding/json"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductionLineFilters struct {
	base_filters.BaseFilters
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	VendorId    *string `json:"vendorId,omitempty"`
}

func NewProductionLineFilters(query string) (*ProductionLineFilters, error) {
	ProductionLineFilters := ProductionLineFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &ProductionLineFilters); err != nil {
			return nil, err
		}
	}

	return &ProductionLineFilters, nil
}

func (productionLineFilters *ProductionLineFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if productionLineFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &productionLineFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PRODUCTION_LINE)

	// Apply vendorId filter
	if productionLineFilters.VendorId != nil {
		vendorId, err := primitive.ObjectIDFromHex(*productionLineFilters.VendorId)
		if err == nil {
			matchStage["vendorId"] = vendorId
		}
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})
	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(productionLineFilters.PopFields, enums.PRODUCTION_LINE, pipeline)

		pipeline = base_filters.BuildSortingPipeline(productionLineFilters.Sorts, pipeline)
	}

	return pipeline
}
