package payment_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/payment/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/payment/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/payment/utils"
)

type PaymentRepo struct {
	*base_repos.BaseRepo[models.Payment, models.PaymentResponse]
}

func NewPaymentRepo(i *do.Injector) (services.PaymentRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &PaymentRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Payment, models.PaymentResponse](db, enums.PAYMENT),
	}, nil
}

func (repo *PaymentRepo) Create(ctx context.Context, createPaymentDto *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	payment := utils.ConvertToPaymentModel(createPaymentDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *payment)
}

func (repo *PaymentRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PaymentResponse, *helpers.ApiError) {
	// Create filters from query parameters
	paymentFilters, _ := filters.NewPaymentFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := paymentFilters.BuildPipeline(matchStage, ctx, true)

	payment, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return payment, nil
}

func (repo *PaymentRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PaymentsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	paymentFilters, _ := filters.NewPaymentFilters(query)

	// Build pipeline with filters
	pipeline := paymentFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &paymentFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PaymentRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updatePaymentDto *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	payment := utils.ConvertToPaymentModel(updatePaymentDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *payment, matchStage)
}

func (repo *PaymentRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *PaymentRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
