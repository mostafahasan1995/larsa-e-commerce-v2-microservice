package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type BannerFilters struct {
	base_filters.BaseFilters
	IsActive *bool `json:"isActive,omitempty" bson:"isActive,omitempty"`
}

func NewBannerFilters(query string) (*BannerFilters, error) {
	bannerFilters := BannerFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &bannerFilters); err != nil {
			return nil, err
		}
	}

	return &bannerFilters, nil
}

func (bannerFilters *BannerFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if bannerFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	if bannerFilters.IsActive != nil {
		matchStage["isActive"] = *bannerFilters.IsActive
	}

	// Apply base filters
	baseFilters := &bannerFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.BANNER)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(bannerFilters.PopFields, enums.BANNER, pipeline)

		pipeline = base_filters.BuildSortingPipeline(bannerFilters.Sorts, pipeline)
	}

	return pipeline
}
