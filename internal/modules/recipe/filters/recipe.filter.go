package recipe_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	recipe_enums "larsa-e-commerce-microservice-v2/internal/modules/recipe/enums"
	recipe_utils "larsa-e-commerce-microservice-v2/internal/modules/recipe/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeFilters struct {
	base_filters.BaseFilters

	Name       string                         `json:"name"`
	Category   string                         `json:"category"`
	Difficulty []recipe_enums.DifficultyLevel `json:"difficulty"`
	RecipeType []recipe_enums.RecipeType      `json:"recipeType"`
	IsActive   *bool                          `json:"isActive"`
	Tags       []string                       `json:"tags"`
	ProductId  *primitive.ObjectID            `json:"productId"`

	// Cooking time filters
	CookingTimeMin *int `json:"cookingTimeMin"`
	CookingTimeMax *int `json:"cookingTimeMax"`
}

func NewRecipeFilters(query string) (*RecipeFilters, error) {
	recipeFilters := RecipeFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &recipeFilters); err != nil {
			return nil, err
		}
	}
	return &recipeFilters, nil
}

func (f *RecipeFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{}

	// Apply base filters
	baseFilters := &f.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.RECIPE)

	// Name filter
	if f.Name != "" {
		matchStage["$or"] = []bson.M{
			{"name.en": bson.M{"$regex": f.Name, "$options": "i"}},
			{"name.ar": bson.M{"$regex": f.Name, "$options": "i"}},
		}
	}

	// Category filter
	if f.Category != "" {
		matchStage["category"] = bson.M{"$regex": f.Category, "$options": "i"}
	}

	// Difficulty filter (multiple values)
	if len(f.Difficulty) > 0 {
		matchStage["difficultyLevel"] = bson.M{"$in": f.Difficulty}
	}

	// Recipe type filter (multiple values)
	if len(f.RecipeType) > 0 {
		matchStage["recipeType"] = bson.M{"$in": f.RecipeType}
	}

	// IsActive filter
	if f.IsActive != nil {
		matchStage["isActive"] = *f.IsActive
	}

	// Tags filter
	if len(f.Tags) > 0 {
		matchStage["tags"] = bson.M{"$in": f.Tags}
	}

	// Product filter (recipes that contain a specific product)
	if f.ProductId != nil {
		matchStage["productItems.productId"] = *f.ProductId
	}

	// Cooking time range filter (min/max)
	if f.CookingTimeMin != nil || f.CookingTimeMax != nil {
		cookingTimeFilter := bson.M{}
		if f.CookingTimeMin != nil {
			cookingTimeFilter["$gte"] = *f.CookingTimeMin
		}
		if f.CookingTimeMax != nil {
			cookingTimeFilter["$lte"] = *f.CookingTimeMax
		}
		matchStage["cookingTime"] = cookingTimeFilter
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchStage,
		})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(f.PopFields, enums.RECIPE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(f.Sorts, pipeline)
		pipeline = append(pipeline, recipe_utils.PopulateRecipeProductsAndUnits...)
	}

	return pipeline
}
