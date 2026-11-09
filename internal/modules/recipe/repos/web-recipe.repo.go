package recipe_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	recipe_filters "larsa-e-commerce-microservice-v2/internal/modules/recipe/filters"
	recipe_models "larsa-e-commerce-microservice-v2/internal/modules/recipe/models"
	recipe_service "larsa-e-commerce-microservice-v2/internal/modules/recipe/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type WebRecipeRepository struct {
	*repo.BaseRepo[recipe_models.Recipe, recipe_models.RecipeResponse]
}

func NewWebRecipeRepository(i *do.Injector) (recipe_service.WebRecipeRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebRecipeRepository{
		BaseRepo: repo.NewBaseRepo[recipe_models.Recipe, recipe_models.RecipeResponse](db, enums.RECIPE),
	}, nil
}

func (repo *WebRecipeRepository) GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError) {
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

func (repo *WebRecipeRepository) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError) {
	filter, err := recipe_filters.NewRecipeFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *WebRecipeRepository) GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError) {
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
