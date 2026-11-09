package survey_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaiterFilters struct {
	base_filters.BaseFilters
	Name     *string             `json:"name,omitempty" schema:"name"`
	BranchId *primitive.ObjectID `json:"branchId,omitempty" schema:"branchId"`
	VendorId *primitive.ObjectID `json:"vendorId,omitempty" schema:"vendorId"`
}

func NewWaiterFilters(query string) (*WaiterFilters, error) {
	waiterFilters := WaiterFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &waiterFilters); err != nil {
			return nil, err
		}
	}

	return &waiterFilters, nil
}

// Helper function to check if bson.M is empty
func isEmptyBSON(m bson.M) bool {
	return len(m) == 0
}

// Helper function to check if a string is empty
func isEmptyString(s string) bool {
	return s == ""
}

func (f *WaiterFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if f == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &f.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.WAITER)

	// Apply name filter
	if f.Name != nil && *f.Name != "" {
		matchStage["name"] = bson.M{"$regex": *f.Name, "$options": "i"}
	}

	// Apply branchId filter
	if f.BranchId != nil {
		matchStage["branchId"] = *f.BranchId
	}

	// Apply vendorId filter
	if f.VendorId != nil {
		matchStage["vendorId"] = *f.VendorId
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(f.PopFields, enums.WAITER, pipeline, ctx)
		pipeline = base_filters.BuildSortingPipeline(f.Sorts, pipeline)
	}

	return pipeline
}
