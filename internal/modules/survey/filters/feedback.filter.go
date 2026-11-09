package survey_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackFilters struct {
	base_filters.BaseFilters
	WaiterId   *primitive.ObjectID `json:"waiterId,omitempty" schema:"waiterId"`
	BranchId   *primitive.ObjectID `json:"branchId,omitempty" schema:"branchId"`
	VendorId   *primitive.ObjectID `json:"vendorId,omitempty" schema:"vendorId"`
	ClientName *string             `json:"clientName,omitempty" schema:"clientName"`
	Phone      *string             `json:"phone,omitempty" schema:"phone"`
	Email      *string             `json:"email,omitempty" schema:"email"`
}

func NewFeedbackFilters(query string) (*FeedbackFilters, error) {
	feedbackFilters := FeedbackFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &feedbackFilters); err != nil {
			return nil, err
		}
	}

	return &feedbackFilters, nil
}

func (feedbackFilters *FeedbackFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if feedbackFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &feedbackFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.FEEDBACK)

	// Apply waiterId filter
	if feedbackFilters.WaiterId != nil {
		matchStage["waiterId"] = *feedbackFilters.WaiterId
	}

	// Apply branchId filter
	if feedbackFilters.BranchId != nil {
		matchStage["branchId"] = *feedbackFilters.BranchId
	}

	// Apply vendorId filter
	if feedbackFilters.VendorId != nil {
		matchStage["vendorId"] = *feedbackFilters.VendorId
	}

	// Apply clientName filter
	if feedbackFilters.ClientName != nil && *feedbackFilters.ClientName != "" {
		matchStage["clientName"] = bson.M{"$regex": *feedbackFilters.ClientName, "$options": "i"}
	}

	// Apply phone filter
	if feedbackFilters.Phone != nil && *feedbackFilters.Phone != "" {
		matchStage["phone"] = *feedbackFilters.Phone
	}

	// Apply email filter
	if feedbackFilters.Email != nil && *feedbackFilters.Email != "" {
		matchStage["email"] = *feedbackFilters.Email
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(feedbackFilters.PopFields, enums.FEEDBACK, pipeline, ctx)
		pipeline = base_filters.BuildSortingPipeline(feedbackFilters.Sorts, pipeline)
	}

	return pipeline
}
