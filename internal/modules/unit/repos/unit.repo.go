package unit_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/unit/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/unit/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/unit/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/unit/utils"
)

type UnitRepo struct {
	*base_repos.BaseRepo[models.Unit, models.UnitResponse]
}

func NewUnitRepo(i *do.Injector) (services.UnitRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &UnitRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Unit, models.UnitResponse](db, enums.UNIT),
	}, nil
}

func (repo *UnitRepo) Create(ctx context.Context, createUnitDto *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	unit := utils.ConvertToUnitModel(createUnitDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *unit)
}

func (repo *UnitRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError) {
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

func (repo *UnitRepo) GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError) {
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

func (repo *UnitRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError) {
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

func (repo *UnitRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateUnitDto *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	unit := utils.ConvertToUnitModel(updateUnitDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *unit, matchStage)
}

func (repo *UnitRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *UnitRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
