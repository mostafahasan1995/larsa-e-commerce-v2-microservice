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

type PayoutRepo struct {
	*base_repos.BaseRepo[models.Payout, models.PayoutResponse]
}

func NewPayoutRepo(i *do.Injector) (services.PayoutRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &PayoutRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Payout, models.PayoutResponse](db, enums.PAYOUT),
	}, nil
}

func (repo *PayoutRepo) Create(ctx context.Context, createPayoutDto *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	payout := utils.ConvertToPayoutModel(createPayoutDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *payout)
}

func (repo *PayoutRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PayoutResponse, *helpers.ApiError) {
	// Create filters from query parameters
	payoutFilters, _ := filters.NewPayoutFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}
	pipeline := payoutFilters.BuildPipeline(matchStage, ctx, true)

	payout, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return payout, nil
}

func (repo *PayoutRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PayoutsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}
	// Create filters from query parameters
	payoutFilters, _ := filters.NewPayoutFilters(query)

	// Build pipeline with filters
	pipeline := payoutFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &payoutFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PayoutRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updatePayoutDto *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	payout := utils.ConvertToPayoutModel(updatePayoutDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *payout, matchStage)
}

func (repo *PayoutRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *PayoutRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
