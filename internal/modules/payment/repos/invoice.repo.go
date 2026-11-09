package payment_repos

import (
	"context"
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

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepo struct {
	*base_repos.BaseRepo[models.Invoice, models.InvoiceResponse]
}

func NewInvoiceRepo(i *do.Injector) (services.InvoiceRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &InvoiceRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Invoice, models.InvoiceResponse](db, enums.INVOICE),
	}, nil
}

func (repo *InvoiceRepo) Create(ctx context.Context, createInvoiceDto *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	invoice := utils.ConvertToInvoiceModel(createInvoiceDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *invoice)
}

func (repo *InvoiceRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.InvoiceResponse, *helpers.ApiError) {
	// Create filters from query parameters
	invoiceFilters, _ := filters.NewInvoiceFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := invoiceFilters.BuildPipeline(matchStage, ctx, true)

	invoice, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return invoice, nil
}

func (repo *InvoiceRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.InvoicesWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	invoiceFilters, _ := filters.NewInvoiceFilters(query)

	// Build pipeline with filters
	pipeline := invoiceFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &invoiceFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *InvoiceRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateInvoiceDto *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	invoice := utils.ConvertToInvoiceModel(updateInvoiceDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *invoice, matchStage)
}

func (repo *InvoiceRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *InvoiceRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
