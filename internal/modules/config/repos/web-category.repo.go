package config_repos

import (
	"context"
	filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

type WebCategoryRepo struct {
	*base_repos.BaseRepo[models.Category, models.CategoryResponse]
}

func NewWebCategoryRepo(i *do.Injector) (services.WebCategoryRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebCategoryRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Category, models.CategoryResponse](db, enums.CATEGORY),
	}, nil
}

func (repo *WebCategoryRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError) {

	// Create filters from query parameters
	categoryFilters, _ := filters.NewCategoryFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := categoryFilters.BuildPipeline(matchStage, ctx, true)

	category, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return category, nil
}

func (repo *WebCategoryRepo) GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	categoryFilters, _ := filters.NewCategoryFilters(query)

	// Build pipeline with filters
	pipeline := categoryFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WebCategoryRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	categoryFilters, _ := filters.NewCategoryFilters(query)

	// Build pipeline with filters
	pipeline := categoryFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &categoryFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
