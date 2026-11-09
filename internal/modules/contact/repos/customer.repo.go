package contact_repos

import (
	"context"
	"net/http"

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

type CustomerRepo struct {
	*base_repos.BaseRepo[models.Customer, models.CustomerResponse]
}

func NewCustomerRepo(i *do.Injector) (services.CustomerRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CustomerRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Customer, models.CustomerResponse](db, enums.CUSTOMER),
	}, nil
}

func (repo *CustomerRepo) Create(ctx context.Context, createCustomerDto *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	customer := utils.ConvertToCustomerModel(createCustomerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	customer.Id = *createCustomerDto.UserId
	createdCustomer, err := repo.BaseRepo.Create(ctx, *customer)
	if err != nil {
		if err.StatusCode == http.StatusConflict {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusConflict,
				Message:    "Email already exists",
			}
		}
		return nil, err
	}
	return createdCustomer, nil
}

func (repo *CustomerRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CustomerResponse, *helpers.ApiError) {
	// Create filters from query parameters
	customerFilters, _ := filters.NewCustomerFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := customerFilters.BuildPipeline(matchStage, ctx, true)

	customer, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return customer, nil
}

func (repo *CustomerRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CustomersWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	customerFilters, _ := filters.NewCustomerFilters(query)

	// Build pipeline with filters
	pipeline := customerFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &customerFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CustomerRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateCustomerDto *dtos.CreateCustomerDto) (*models.CustomerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	customer := utils.ConvertToCustomerModel(updateCustomerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *customer, matchStage)
}

func (repo *CustomerRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *CustomerRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
