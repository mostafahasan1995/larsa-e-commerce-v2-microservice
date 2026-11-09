package product_repos

import (
	"context"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type WebProductRepo struct {
	*base_repos.BaseRepo[models.Product, models.ProductResponse]
}

func NewWebProductRepo(i *do.Injector) (services.WebProductRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebProductRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Product, models.ProductResponse](db, enums.PRODUCT),
	}, nil
}

func (repo *WebProductRepo) GetById(ctx context.Context, args product_utils.GetByIdArgs) (*models.ProductResponse, *helpers.ApiError) {
	// Create filters from query parameters
	productFilters, _ := filters.NewProductFilters(args.Query)
	productFilters.IncludeAllProductTypes = &[]bool{true}[0]

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := productFilters.BuildPipeline(matchStage, ctx, filters.BuildPipelineArgs{
		AddExtraOperations: true,
		BranchId:           args.BranchId,
		CurrencyCode:       args.CurrencyCode,
		UnitId:             args.UnitId,
	})

	product, findErr := repo.BaseRepo.GetById(ctx, args.Id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return product, nil
}

func (repo *WebProductRepo) GetWithPagination(ctx context.Context, args product_utils.GetWithPaginationArgs) (*models.ProductsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	productFilters, _ := filters.NewProductFilters(args.Query)

	// Build pipeline with filters
	pipeline := productFilters.BuildPipeline(matchStage, ctx, filters.BuildPipelineArgs{
		AddExtraOperations: true,
		BranchId:           args.BranchId,
		UnitId:             args.UnitId,
		CurrencyCode:       args.CurrencyCode,
	})

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &productFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           args.Skip,
		Limit:          args.Limit,
	})

	if err != nil {
		return nil, err
	}

	// Remove attributes to reduce the response size
	for i, _ := range results.Data {
		if productFilters.IncludeAttributes != nil && *productFilters.IncludeAttributes == false {
			results.Data[i].Attributes = []models.ProductAttribute{}
		}
	}

	return results, nil
}
