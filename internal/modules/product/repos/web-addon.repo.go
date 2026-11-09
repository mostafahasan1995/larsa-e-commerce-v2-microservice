package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
)

type WebAddonRepo struct {
	*base_repos.BaseRepo[models.Addon, models.AddonResponse]
}

func NewWebAddonRepo(i *do.Injector) (services.WebAddonRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebAddonRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Addon, models.AddonResponse](db, enums.ADDON),
	}, nil
}

func (repo *WebAddonRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError) {
	// Build pipeline for aggregation
	addonFilters, _ := filters.NewAddonFilters(query)

	// Build match stage
	matchStage := bson.M{}

	// Build pipeline using filters
	pipeline := addonFilters.BuildPipeline(matchStage, ctx, true)

	addon, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return addon, nil
}

func (repo *WebAddonRepo) GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError) {
	// Create addon filters from query string
	addonFilters, _ := filters.NewAddonFilters(query)

	// Build match stage
	matchStage := bson.M{}

	// Build pipeline using filters
	pipeline := addonFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WebAddonRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError) {
	// Create addon filters from query string
	addonFilters, _ := filters.NewAddonFilters(query)

	// Build match stage
	matchStage := bson.M{}

	// Build pipeline using filters
	pipeline := addonFilters.BuildPipeline(matchStage, ctx, true)

	results, apiErr := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &addonFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if apiErr != nil {
		return nil, apiErr
	}

	return results, nil
}
