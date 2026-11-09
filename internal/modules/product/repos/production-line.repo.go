package product_repos

import (
	"context"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
)

type ProductionLineRepo struct {
	*base_repos.BaseRepo[models.ProductionLine, models.ProductionLineResponse]
}

func NewProductionLineRepo(i *do.Injector) (services.ProductionLineRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ProductionLineRepo{
		BaseRepo: base_repos.NewBaseRepo[models.ProductionLine, models.ProductionLineResponse](db, enums.PRODUCTION_LINE),
	}, nil
}

func (repo *ProductionLineRepo) Create(ctx context.Context, createProductionLineDto *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	productionLine := utils.ConvertToProductionLineModel(createProductionLineDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *productionLine)
}

func (repo *ProductionLineRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductionLineResponse, *helpers.ApiError) {

	// Create filters from query parameters
	productionLineFilters, _ := filters.NewProductionLineFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := productionLineFilters.BuildPipeline(matchStage, ctx, true)

	productionLine, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return productionLine, nil
}

func (repo *ProductionLineRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductionLinesWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	productFilters, _ := filters.NewProductionLineFilters(query)

	// Build pipeline with filters
	pipeline := productFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &productFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ProductionLineRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateProductionLineDto *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	// Build pipeline with filters
	matchStage := bson.M{}

	productionLine := utils.ConvertToProductionLineModel(updateProductionLineDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *productionLine, matchStage)
}

func (repo *ProductionLineRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
		if err != nil {
			return nil, err
		}

		unlinkProductErr := repo.unlinkProductionLineFromProduct(ctx, sessCtx, id)
		if unlinkProductErr != nil {
			return nil, unlinkProductErr
		}
		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *ProductionLineRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
		if bulkDeleteErr != nil {
			return nil, bulkDeleteErr
		}

		for _, id := range ids {

			unlinkProductErr := repo.unlinkProductionLineFromProduct(ctx, sessCtx, id)
			if unlinkProductErr != nil {
				return nil, unlinkProductErr
			}
		}
		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *ProductionLineRepo) unlinkProductionLineFromProduct(ctx context.Context, sessCtx context.Context, productionLineId primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	productCollection := repo.BaseRepo.GetDB().Database(cfg.Db).Collection(enums.PRODUCT.GetCollectionName())

	// Remove Old One
	productUpdateManyFilter := bson.M{
		"productionLineId": productionLineId,
		"deletedAt":        bson.M{"$exists": false},
	}

	productUpdateMany := bson.M{
		"$set": bson.M{
			"productionLineId": nil,
			"updatedAt":        time.Now(),
			"updatedBy":        cfg.User.Id,
		},
	}
	_, updateManyErr := productCollection.UpdateMany(sessCtx, productUpdateManyFilter, productUpdateMany)
	if updateManyErr != nil {
		return helpers.DefaultError(updateManyErr)
	}

	return nil
}
