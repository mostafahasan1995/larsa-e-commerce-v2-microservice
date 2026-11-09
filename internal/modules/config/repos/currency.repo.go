package config_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"

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

type CurrencyRepo struct {
	*base_repos.BaseRepo[models.Currency, models.CurrencyResponse]
}

func NewCurrencyRepo(i *do.Injector) (services.CurrencyRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CurrencyRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Currency, models.CurrencyResponse](db, enums.CURRENCY),
	}, nil
}

func (repo *CurrencyRepo) Create(ctx context.Context, createCurrencyDto *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	currency := utils.ConvertToCurrencyModel(createCurrencyDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, currency)
}

func (repo *CurrencyRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CurrencyResponse, *helpers.ApiError) {

	// Create filters from query parameters
	currencyFilters, _ := filters.NewCurrencyFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := currencyFilters.BuildPipeline(matchStage, ctx, true)

	currency, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return currency, nil
}

func (repo *CurrencyRepo) GetAll(ctx context.Context, query string) ([]models.CurrencyResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	currencyFilters, _ := filters.NewCurrencyFilters(query)

	// Build pipeline with filters
	pipeline := currencyFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CurrencyRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CurrenciesWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	currencyFilters, _ := filters.NewCurrencyFilters(query)

	// Build pipeline with filters
	pipeline := currencyFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &currencyFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CurrencyRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateCurrencyDto *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	matchStage := bson.M{}

	currency := utils.ConvertToCurrencyModel(updateCurrencyDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, currency, matchStage)
}

func (repo *CurrencyRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *CurrencyRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
