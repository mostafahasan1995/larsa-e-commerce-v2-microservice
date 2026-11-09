package recipe_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	recipe_dto "larsa-e-commerce-microservice-v2/internal/modules/recipe/dtos"
	recipe_models "larsa-e-commerce-microservice-v2/internal/modules/recipe/models"
	recipe_utils "larsa-e-commerce-microservice-v2/internal/modules/recipe/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeRepo interface {
	Create(ctx context.Context, recipe recipe_models.Recipe) (*recipe_models.RecipeResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, recipe recipe_models.Recipe, filter bson.M) (*recipe_models.RecipeResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}, userId primitive.ObjectID) (*recipe_models.RecipeResponse, *helpers.ApiError)
}

type RecipeService interface {
	Create(ctx context.Context, dto *recipe_dto.CreateRecipeDTO) (*recipe_models.RecipeResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *recipe_dto.UpdateRecipeDTO) (*recipe_models.RecipeResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*recipe_models.RecipeResponse, *helpers.ApiError)
}

type recipe struct {
	repo RecipeRepo
}

func NewRecipeService(i *do.Injector) (RecipeService, error) {
	return &recipe{
		repo: do.MustInvoke[RecipeRepo](i),
	}, nil
}

func (s *recipe) Create(ctx context.Context, dto *recipe_dto.CreateRecipeDTO) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	config, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}
	recipe, apiErr := recipe_utils.ConvertCreateRecipeDTOToModel(dto, config.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}
	result, apiErr := s.repo.Create(ctx, recipe)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *recipe) Update(ctx context.Context, id primitive.ObjectID, dto *recipe_dto.UpdateRecipeDTO) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	config, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}
	recipe, apiErr := recipe_utils.ConvertUpdateRecipeDTOToModel(dto, id, config.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}

	result, apiErr := s.repo.UpdateById(ctx, id, recipe, bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *recipe) GetById(ctx context.Context, id primitive.ObjectID, query string) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, query)
}

func (s *recipe) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id, bson.M{})
}

func (s *recipe) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*recipe_models.RecipesWithPagination, *helpers.ApiError) {
	return s.repo.GetWithPagination(ctx, skip, limit, query)
}

func (s *recipe) GetAll(ctx context.Context, query string) ([]recipe_models.RecipeResponse, *helpers.ApiError) {
	return s.repo.GetAll(ctx, query)
}

func (s *recipe) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*recipe_models.RecipeResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}

	return s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)
}
