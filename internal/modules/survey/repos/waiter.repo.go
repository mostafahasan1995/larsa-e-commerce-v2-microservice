package survey_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WaiterRepository struct {
	*repo.BaseRepo[survey_models.WaiterModel, survey_models.WaiterResponse]
}

func NewWaiterRepository(i *do.Injector) (survey_services.WaiterRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WaiterRepository{
		BaseRepo: repo.NewBaseRepo[survey_models.WaiterModel, survey_models.WaiterResponse](db, enums.WAITER),
	}, nil
}

func (repo *WaiterRepository) Create(ctx context.Context, waiterModel survey_models.WaiterModel) (*survey_models.WaiterResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, waiterModel)
}

func (repo *WaiterRepository) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.WaiterResponse, *helpers.ApiError) {
	// Create filters from query parameters
	waiterFilters, _ := survey_filters.NewWaiterFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := waiterFilters.BuildPipeline(matchStage, ctx, true)

	waiter, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return waiter, nil
}

func (repo *WaiterRepository) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	waiterFilters, _ := survey_filters.NewWaiterFilters(query)

	// Build pipeline with filters
	pipeline := waiterFilters.BuildPipeline(matchStage, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &waiterFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *WaiterRepository) GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	waiterFilters, _ := survey_filters.NewWaiterFilters(query)

	// Build pipeline with filters
	pipeline := waiterFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WaiterRepository) UpdateById(ctx context.Context, id primitive.ObjectID, waiter survey_models.WaiterModel, filter bson.M) (*survey_models.WaiterResponse, *helpers.ApiError) {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{}
	}

	return repo.BaseRepo.UpdateById(ctx, id, waiter, filter)
}

func (repo *WaiterRepository) DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{"deletedAt": bson.M{"$exists": false}}
	}

	return repo.BaseRepo.DeleteById(ctx, id, filter)
}
