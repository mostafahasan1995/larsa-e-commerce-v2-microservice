package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttributeFilters struct {
	base_filters.BaseFilters
	Name         *types.Localizable[string]     `json:"name,omitempty"`
	Tag          *product_enums.AttributeTag    `json:"tag,omitempty"`
	VendorId     *string                        `json:"vendorId,omitempty"`
	IsFiltrable  *bool                          `json:"isFiltrable,omitempty"`
	Required     *bool                          `json:"required,omitempty"`
	IsComparable *bool                          `json:"isComparable,omitempty"`
	Visibility   *bool                          `json:"visibility,omitempty"`
	Configurable *bool                          `json:"configurable,omitempty"`
	BelongsTo    *[]product_enums.BelongTo      `json:"belongsTo,omitempty"`
	Codes        *[]string                      `json:"codes,omitempty"`
	Types        *[]product_enums.AttributeType `json:"types,omitempty"`
}

func NewAttributeFilters(query string) (*AttributeFilters, error) {
	attributeFilters := AttributeFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &attributeFilters); err != nil {
			return nil, err
		}
	}

	return &attributeFilters, nil
}

func (attributeFilters *AttributeFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if attributeFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &attributeFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.ATTRIBUTE)

	// Apply types filter
	if attributeFilters.Types != nil && len(*attributeFilters.Types) > 0 {
		matchStage["type"] = bson.M{"$in": *attributeFilters.Types}
	}

	if attributeFilters.Codes != nil && len(*attributeFilters.Codes) > 0 {
		matchStage["code"] = bson.M{"$in": *attributeFilters.Codes}
	}

	// Apply tag filter
	if attributeFilters.Tag != nil {
		matchStage["tag"] = *attributeFilters.Tag
	}

	// Apply belongsTo filter
	if attributeFilters.BelongsTo != nil && len(*attributeFilters.BelongsTo) > 0 {
		matchStage["belongsTo"] = bson.M{"$in": *attributeFilters.BelongsTo}
	}

	// Apply vendorId filter
	if attributeFilters.VendorId != nil {
		vendorId, err := primitive.ObjectIDFromHex(*attributeFilters.VendorId)
		if err == nil {
			matchStage["vendorId"] = vendorId
		}
	}

	// Apply isFiltrable filter
	if attributeFilters.IsFiltrable != nil {
		matchStage["isFiltrable"] = *attributeFilters.IsFiltrable
	}

	// Apply required filter
	if attributeFilters.Required != nil {
		matchStage["required"] = *attributeFilters.Required
	}

	// Apply isComparable filter
	if attributeFilters.IsComparable != nil {
		matchStage["isComparable"] = *attributeFilters.IsComparable
	}

	// Apply visibility filter
	if attributeFilters.Visibility != nil {
		matchStage["visibility"] = *attributeFilters.Visibility
	}

	// Apply configurable filter
	if attributeFilters.Configurable != nil {
		matchStage["configurable"] = *attributeFilters.Configurable
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(attributeFilters.PopFields, enums.ATTRIBUTE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(attributeFilters.Sorts, pipeline)
	}

	return pipeline
}
