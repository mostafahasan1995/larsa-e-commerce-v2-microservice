package filters

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DateRange represents a date range with greater than or equal and less than or equal
type DateRange struct {
	Gte *time.Time `json:"gte,omitempty"`
	Lte *time.Time `json:"lte,omitempty"`
}

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

// Sort represents a sort configuration for a field
type Sort struct {
	Field string    `json:"field" validate:"required"`
	Order SortOrder `json:"order" validate:"required,oneof=	"` // "asc" or "desc"
}

// SearchField represents a field to search in
type SearchField struct {
	Name        string `json:"name"`
	IsLocalized bool   `json:"isLocalized"`
}

// BaseFilters represents common filter fields that can be reused across different filters
type BaseFilters struct {
	// ID filters
	Ids        *[]primitive.ObjectID `json:"ids,omitempty" validate:"omitempty"`
	ExcludeIds *[]primitive.ObjectID `json:"excludeIds,omitempty" validate:"omitempty"`

	// Date range filters
	CreatedAt *DateRange `json:"createdAt,omitempty" validate:"omitempty"`
	UpdatedAt *DateRange `json:"updatedAt,omitempty" validate:"omitempty"`

	// Populate fields
	PopFields      *[]string `json:"popFields,omitempty" validate:"omitempty"`
	IncludeDeleted *bool     `json:"includeDeleted,omitempty" validate:"omitempty"`

	// Sort fields
	Sorts []Sort `json:"sorts,omitempty" validate:"omitempty,dive"`

	// Search query
	Q           *string `json:"q,omitempty" validate:"omitempty"`
	SearchQuery *string `json:"searchQuery,omitempty" validate:"omitempty"`
}

// ApplyBaseFilters applies the common base filters to the MongoDB query
// By default, it excludes soft-deleted documents (where deletedAt doesn't exist).
// To include soft-deleted documents, set IncludeDeleted to true.
func (filters *BaseFilters) ApplyBaseFilters(ctx context.Context, m bson.M, collectionName enums.CollectionName) {
	cfg, _ := util.GetReqAppCfg(ctx)
	acceptLanguage := cfg.Lang

	if acceptLanguage == "" {
		acceptLanguage = "en"
	}

	searchFields := base_utils.GetSearchQueryField(collectionName)
	if _, exists := m["$and"]; !exists {
		m["$and"] = []bson.M{}
	}
	// Search query filter
	if filters.Q != nil && *filters.Q != "" && len(searchFields) > 0 {
		searchConditions := make([]bson.M, 0)
		for _, field := range searchFields {
			if field.IsLocalized {
				if collectionName == enums.PRODUCT {
					// For product, search both name.en and name.ar
					searchConditions = append(searchConditions, bson.M{
						field.FieldName + ".en": bson.M{"$regex": filters.Q, "$options": "i"},
					})
					searchConditions = append(searchConditions, bson.M{
						field.FieldName + ".ar": bson.M{"$regex": filters.Q, "$options": "i"},
					})
				} else {
					fieldName := field.FieldName + "." + acceptLanguage
					searchConditions = append(searchConditions, bson.M{
						fieldName: bson.M{"$regex": filters.Q, "$options": "i"},
					})
				}
			} else {
				searchConditions = append(searchConditions, bson.M{
					field.FieldName: bson.M{"$regex": filters.Q, "$options": "i"},
				})
			}
		}
		if len(searchConditions) > 0 {
			m["$and"] = append(m["$and"].([]bson.M), bson.M{"$or": searchConditions})
		}
	}

	if filters.SearchQuery != nil && *filters.SearchQuery != "" && len(searchFields) > 0 {
		searchConditions := make([]bson.M, 0)
		for _, field := range searchFields {
			if field.IsLocalized {
				if collectionName == enums.PRODUCT {
					// For product, search both name.en and name.ar
					searchConditions = append(searchConditions, bson.M{
						field.FieldName + ".en": bson.M{"$regex": filters.SearchQuery, "$options": "i"},
					})
					searchConditions = append(searchConditions, bson.M{
						field.FieldName + ".ar": bson.M{"$regex": filters.SearchQuery, "$options": "i"},
					})
				} else {
					fieldName := field.FieldName + "." + acceptLanguage
					searchConditions = append(searchConditions, bson.M{
						fieldName: bson.M{"$regex": filters.SearchQuery, "$options": "i"},
					})
				}
			} else {
				searchConditions = append(searchConditions, bson.M{
					field.FieldName: bson.M{"$regex": filters.SearchQuery, "$options": "i"},
				})
			}
		}
		if len(searchConditions) > 0 {
			m["$and"] = append(m["$and"].([]bson.M), bson.M{"$or": searchConditions})
		}
	}

	// ID filters
	if filters.Ids != nil && len(*filters.Ids) > 0 {
		m["_id"] = bson.M{"$in": *filters.Ids}
	}
	if filters.ExcludeIds != nil && len(*filters.ExcludeIds) > 0 {
		m["_id"] = bson.M{"$nin": *filters.ExcludeIds}
	}

	// Date range filters
	if filters.CreatedAt != nil {
		if filters.CreatedAt.Gte != nil {
			m["createdAt"] = bson.M{"$gte": *filters.CreatedAt.Gte}
		}
		if filters.CreatedAt.Lte != nil {
			if _, exists := m["createdAt"]; exists {
				m["createdAt"].(bson.M)["$lte"] = *filters.CreatedAt.Lte
			} else {
				m["createdAt"] = bson.M{"$lte": *filters.CreatedAt.Lte}
			}
		}
	}
	if filters.UpdatedAt != nil {
		if filters.UpdatedAt.Gte != nil {
			m["updatedAt"] = bson.M{"$gte": *filters.UpdatedAt.Gte}
		}
		if filters.UpdatedAt.Lte != nil {
			if _, exists := m["updatedAt"]; exists {
				m["updatedAt"].(bson.M)["$lte"] = *filters.UpdatedAt.Lte
			} else {
				m["updatedAt"] = bson.M{"$lte": *filters.UpdatedAt.Lte}
			}
		}
	}

	// Apply soft-delete filter unless explicitly included
	if filters.IncludeDeleted == nil || !*filters.IncludeDeleted {
		m = base_utils.AddSoftDeletedFilter(m)
	}
	//delete and if empty array
	if _, exists := m["$and"]; exists {
		if len(m["$and"].([]bson.M)) == 0 {
			delete(m, "$and")
		}
	}

}

// BuildPopulatePipeline builds a pipeline for populating referenced fields
func BuildPopulatePipeline(popFields *[]string, collectionName enums.CollectionName, pipeline []bson.M, ctx ...context.Context) []bson.M {
	if popFields == nil || len(*popFields) == 0 {
		return pipeline
	}
	var user *types.User
	if len(ctx) > 0 {
		cfg, _ := util.GetReqAppCfg(ctx[0])
		user = cfg.User
	}
	// Get ref fields for module
	refFields := base_utils.GetRefField(collectionName)

	for _, field := range *popFields {
		ref, exists := refFields[field]
		// TODO: We should not have to check if the field exists, It should always exist.
		if !exists {
			continue
		}

		if ref.IsNeedAuth && user == nil {
			continue
		}

		lookupStage := bson.M{

			"from":         ref.Collection,
			"localField":   ref.LocalField,
			"foreignField": ref.ForeignField,
			"as":           ref.As,
		}

		if ref.Let != nil {
			lookupStage["let"] = ref.Let
		}

		if ref.Pipeline != nil {
			lookupStage["pipeline"] = ref.Pipeline
		}

		pipeline = append(pipeline, bson.M{"$lookup": lookupStage})

		if ref.ExternalPipeline != nil {
			pipeline = append(pipeline, ref.ExternalPipeline...)
		}

		// Only unwind if it's not an array field
		if !ref.IsArray {
			pipeline = append(pipeline, bson.M{
				"$unwind": bson.M{
					"path":                       "$" + ref.As,
					"preserveNullAndEmptyArrays": true,
				},
			})
		}
	}

	return pipeline
}

func BuildSortingPipeline(sorts []Sort, pipeline []bson.M) []bson.M {

	if len(sorts) == 0 {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
		return pipeline
	}

	sortStage := bson.M{}
	for _, sort := range sorts {
		// Convert order to 1 for ascending, -1 for descending
		var order int
		if sort.Order == Asc {
			order = 1
		} else {
			order = -1
		}
		sortStage[sort.Field] = order
	}

	pipeline = append(pipeline, bson.M{"$sort": sortStage})

	return pipeline
}
