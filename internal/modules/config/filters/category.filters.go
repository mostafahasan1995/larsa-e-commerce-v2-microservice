package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryFilters struct {
	base_filters.BaseFilters
	Name         *types.Localizable[string]      `json:"name,omitempty"`
	Code         *types.LocalizableArray[string] `json:"code,omitempty"`
	ParentId     *primitive.ObjectID             `json:"parentId,omitempty"`
	Root         *bool                           `json:"root,omitempty"`
	Position     *int                            `json:"position,omitempty"`
	Visibilities []string                        `json:"visibilities,omitempty"`
}

func NewCategoryFilters(query string) (*CategoryFilters, error) {
	categoryFilters := CategoryFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &categoryFilters); err != nil {
			return nil, err
		}
	}

	return &categoryFilters, nil
}

func (categoryFilters *CategoryFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if categoryFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &categoryFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CATEGORY)

	// Apply parentId filter
	if categoryFilters.ParentId != nil {
		matchStage["parentId"] = *categoryFilters.ParentId

	}

	// Apply root filter
	if categoryFilters.Root != nil && *categoryFilters.Root {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"$or": []bson.M{
				{"parentId": bson.M{"$exists": false}},
				{"parentId": nil},
			},
		}})
	}

	// Apply position filter
	if categoryFilters.Position != nil {
		matchStage["position"] = *categoryFilters.Position
	}

	// Apply visibility filter
	if categoryFilters.Visibilities != nil {
		matchStage["visibilities"] = bson.M{"$in": categoryFilters.Visibilities}
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(categoryFilters.PopFields, enums.CATEGORY, pipeline)

		pipeline = base_filters.BuildSortingPipeline(categoryFilters.Sorts, pipeline)

		pipeline = append(pipeline, config_utils.CategoryCalculateProductAndChildrenCount...)
	}

	return pipeline
}
