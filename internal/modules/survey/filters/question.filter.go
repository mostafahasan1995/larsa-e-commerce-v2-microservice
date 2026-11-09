package survey_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	survey_enums "larsa-e-commerce-microservice-v2/internal/modules/survey/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionFilters struct {
	base_filters.BaseFilters
	Type     *survey_enums.QuestionType `json:"type,omitempty" schema:"type"`
	BranchId *primitive.ObjectID        `json:"branchId,omitempty" schema:"branchId"`
	VendorId *primitive.ObjectID        `json:"vendorId,omitempty" schema:"vendorId"`
}

func NewQuestionFilters(query string) (*QuestionFilters, error) {
	questionFilters := QuestionFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &questionFilters); err != nil {
			return nil, err
		}
	}

	return &questionFilters, nil
}

func (questionFilters *QuestionFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if questionFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &questionFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.QUESTION)

	// Apply type filter
	if questionFilters.Type != nil {
		matchStage["type"] = *questionFilters.Type
	}

	// Apply branchId filter - include questions targeting all branches or specific branch
	if questionFilters.BranchId != nil {
		matchStage["$or"] = []bson.M{
			{"allBranches": true},
			{"branchIds": *questionFilters.BranchId},
		}
	}

	// Apply vendorId filter
	if questionFilters.VendorId != nil {
		matchStage["vendorId"] = *questionFilters.VendorId
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(questionFilters.PopFields, enums.QUESTION, pipeline, ctx)
		pipeline = base_filters.BuildSortingPipeline(questionFilters.Sorts, pipeline)
	}

	return pipeline
}
