package payment_repos

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	filters "larsa-e-commerce-microservice-v2/internal/modules/payment/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/payment/services"
	utils "larsa-e-commerce-microservice-v2/internal/modules/payment/utils"
)

type TransactionRepo struct {
	*base_repos.BaseRepo[models.Transaction, models.TransactionResponse]
}

func NewTransactionRepo(i *do.Injector) (services.TransactionRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &TransactionRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Transaction, models.TransactionResponse](db, enums.TRANSACTION),
	}, nil
}

func (repo *TransactionRepo) Create(ctx context.Context, createTransactionDto *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	transaction := utils.ConvertToTransactionModel(createTransactionDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *transaction)
}

func (repo *TransactionRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.TransactionResponse, *helpers.ApiError) {
	// Create filters from query parameters
	transactionFilters, _ := filters.NewTransactionFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := transactionFilters.BuildPipeline(matchStage, ctx, true)

	transaction, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return transaction, nil
}

func (repo *TransactionRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.TransactionsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	transactionFilters, _ := filters.NewTransactionFilters(query)

	// Build pipeline with filters
	pipeline := transactionFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &transactionFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *TransactionRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateTransactionDto *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	transaction := utils.ConvertToTransactionModel(updateTransactionDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *transaction, matchStage)
}

func (repo *TransactionRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
