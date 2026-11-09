package recipe_service

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	recipe_models "larsa-e-commerce-microservice-v2/internal/modules/recipe/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebRecipeRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError)
}

type WebRecipeService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError)
}

type webRecipe struct {
	repo WebRecipeRepo
}

func NewWebRecipeService(i *do.Injector) (WebRecipeService, error) {
	return &webRecipe{
		repo: do.MustInvoke[WebRecipeRepo](i),
	}, nil
}

func (s *webRecipe) GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, query)
}

func (s *webRecipe) GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError) {
	return s.repo.GetAll(ctx, query)
}

func (s *webRecipe) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError) {
	return s.repo.GetWithPagination(ctx, skip, limit, query)
}
