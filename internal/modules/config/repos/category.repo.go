package config_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"
)

type CategoryRepo struct {
	*base_repos.BaseRepo[models.Category, models.CategoryResponse]
}

func NewCategoryRepo(i *do.Injector) (services.CategoryRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CategoryRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Category, models.CategoryResponse](db, enums.CATEGORY),
	}, nil
}

func (repo *CategoryRepo) Create(ctx context.Context, createCategoryDto *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	category := utils.ConvertToCategoryModel(createCategoryDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *category)
}

func (repo *CategoryRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError) {
	// Build pipeline with filters
	matchStage := bson.M{}

	// Create filters from query parameters
	categoryFilters, _ := filters.NewCategoryFilters(query)

	pipeline := categoryFilters.BuildPipeline(matchStage, ctx, true)

	category, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return category, nil
}

func (repo *CategoryRepo) GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError) {
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

func (repo *CategoryRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError) {

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
		PostPipeline:   config_utils.CategoryCalculateProductAndChildrenCount,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CategoryRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateCategoryDto *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	matchStage := bson.M{}

	category := utils.ConvertToCategoryModel(updateCategoryDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *category, matchStage)
}

func (repo *CategoryRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
