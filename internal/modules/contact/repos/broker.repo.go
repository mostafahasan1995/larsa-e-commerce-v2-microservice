package contact_repos

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
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos"
	filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/contact/services"
	utils "larsa-e-commerce-microservice-v2/internal/modules/contact/utils"
)

type BrokerRepo struct {
	*base_repos.BaseRepo[models.Broker, models.BrokerResponse]
}

func NewBrokerRepo(i *do.Injector) (services.BrokerRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &BrokerRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Broker, models.BrokerResponse](db, enums.BROKER),
	}, nil
}

func (repo *BrokerRepo) Create(ctx context.Context, createBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	broker := utils.ConvertToBrokerModel(createBrokerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *broker)
}

func (repo *BrokerRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BrokerResponse, *helpers.ApiError) {
	// Create filters from query parameters
	brokerFilters, _ := filters.NewBrokerFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}
	pipeline := brokerFilters.BuildPipeline(matchStage, ctx, true)

	broker, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return broker, nil
}

func (repo *BrokerRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BrokersWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	brokerFilters, _ := filters.NewBrokerFilters(query)

	// Build pipeline with filters
	pipeline := brokerFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &brokerFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *BrokerRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	broker := utils.ConvertToBrokerModel(updateBrokerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *broker, matchStage)
}

func (repo *BrokerRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *BrokerRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
