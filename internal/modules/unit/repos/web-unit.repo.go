package unit_repos

import (
	"context"
	filters "larsa-e-commerce-microservice-v2/internal/modules/unit/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/unit/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"
)

type WebUnitRepo struct {
	*base_repos.BaseRepo[models.Unit, models.UnitResponse]
}

func NewWebUnitRepo(i *do.Injector) (services.WebUnitRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebUnitRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Unit, models.UnitResponse](db, enums.UNIT),
	}, nil
}

func (repo *WebUnitRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError) {

	// Create filters from query parameters
	unitFilters, _ := filters.NewUnitFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := unitFilters.BuildPipeline(matchStage, ctx, true)

	unit, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return unit, nil
}

func (repo *WebUnitRepo) GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	unitFilters, _ := filters.NewUnitFilters(query)

	// Build pipeline with filters
	pipeline := unitFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WebUnitRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	unitFilters, _ := filters.NewUnitFilters(query)

	// Build pipeline with filters
	pipeline := unitFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &unitFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
