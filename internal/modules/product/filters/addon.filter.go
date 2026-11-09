package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

// AddonFilters stores filter parameters for addons
type AddonFilters struct {
	base_filters.BaseFilters
	Label        *helpers.StringLocalizableFilter `json:"label,omitempty"`
	MinSelection *helpers.NumericFilter           `json:"minSelection,omitempty"`
	MaxSelection *helpers.NumericFilter           `json:"maxSelection,omitempty"`
	Pricing      *helpers.NumericFilter           `json:"pricing,omitempty"`
}

// NewAddonFilters creates and initializes a new AddonFilters instance
func NewAddonFilters(queryString string) (AddonFilters, error) {
	var filters AddonFilters
	if queryString != "" {
		if err := json.Unmarshal([]byte(queryString), &filters); err != nil {
			return AddonFilters{}, err
		}
	}
	return filters, nil
}

// BuildPipeline builds the aggregation pipeline for addon queries
func (addonFilters *AddonFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if addonFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &addonFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.ADDON)

	if addonFilters.Label != nil && *addonFilters.Label != "" {
		labelFilter := addonFilters.Label.BuildFilter("label")
		if labelFilter != nil {
			pipeline = append(pipeline, bson.M{"$match": labelFilter})
		}
	}

	// Apply minSelection filter
	if addonFilters.MinSelection != nil {
		minSelectionFilter := addonFilters.MinSelection.BuildFilter()
		if minSelectionFilter != nil {
			matchStage["minSelection"] = minSelectionFilter
		}
	}

	// Apply maxSelection filter
	if addonFilters.MaxSelection != nil {
		maxSelectionFilter := addonFilters.MaxSelection.BuildFilter()
		if maxSelectionFilter != nil {
			matchStage["maxSelection"] = maxSelectionFilter
		}
	}

	// Apply pricing filter
	if addonFilters.Pricing != nil {
		pricingFilter := addonFilters.Pricing.BuildFilter()
		if pricingFilter != nil {
			matchStage["pricing"] = pricingFilter
		}
	}

	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {

		pipeline = base_filters.BuildPopulatePipeline(addonFilters.PopFields, enums.ADDON, pipeline, ctx)
		pipeline = base_filters.BuildSortingPipeline(addonFilters.Sorts, pipeline)
	}

	return pipeline
}
