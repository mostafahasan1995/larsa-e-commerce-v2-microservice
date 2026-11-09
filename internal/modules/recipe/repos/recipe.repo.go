package recipe_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	recipe_filters "larsa-e-commerce-microservice-v2/internal/modules/recipe/filters"
	recipe_models "larsa-e-commerce-microservice-v2/internal/modules/recipe/models"
	recipe_service "larsa-e-commerce-microservice-v2/internal/modules/recipe/services"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type RecipeRepository struct {
	*repo.BaseRepo[recipe_models.Recipe, recipe_models.RecipeResponse]
}

func NewRecipeRepository(i *do.Injector) (recipe_service.RecipeRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &RecipeRepository{
		BaseRepo: repo.NewBaseRepo[recipe_models.Recipe, recipe_models.RecipeResponse](db, enums.RECIPE),
	}, nil
}

func (repo *RecipeRepository) GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	filter, err := recipe_filters.NewRecipeFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := []bson.M{}
	pipeline = filter.BuildPipeline(pipeline, ctx, true)

	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *RecipeRepository) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError) {
	filter, err := recipe_filters.NewRecipeFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)
	return repo.BaseRepo.GetWithOptimizedPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *RecipeRepository) GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError) {
	filter, err := recipe_filters.NewRecipeFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := []bson.M{}
	pipeline = filter.BuildPipeline(pipeline, ctx, true)

	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *RecipeRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = userId

	result, updateErr := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if updateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update document: " + updateErr.Error(),
		}
	}

	if result.MatchedCount == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Document not found",
		}
	}

	return repo.BaseRepo.GetById(ctx, id, []bson.M{})
}
