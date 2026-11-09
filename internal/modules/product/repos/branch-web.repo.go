package product_repos

import (
	"context"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type BranchWebRepo struct {
	*base_repos.BaseRepo[models.Branch, models.BranchResponse]
}

func NewBranchWebRepo(i *do.Injector) (services.BranchWebRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &BranchWebRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Branch, models.BranchResponse](db, enums.BRANCH),
	}, nil
}

func (repo *BranchWebRepo) GetById(ctx context.Context, id primitive.ObjectID, branchFilters filters.BranchFilters) (*models.BranchResponse, *helpers.ApiError) {
	matchStage := bson.M{"isActive": true}
	branchFilters.ProductBranchFilters = &product_filters.ProductBranchFilters{
		IsActive: &[]bool{true}[0],
	}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, true)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *BranchWebRepo) GetAll(ctx context.Context, branchFilters filters.BranchFilters) ([]models.BranchResponse, *helpers.ApiError) {
	matchStage := bson.M{"isActive": true}
	branchFilters.ProductBranchFilters = &product_filters.ProductBranchFilters{
		IsActive: &[]bool{true}[0],
	}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, true)
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *BranchWebRepo) GetWithPagination(ctx context.Context, skip, limit int64, branchFilters filters.BranchFilters) (*models.BranchesWithPagination, *helpers.ApiError) {
	matchStage := bson.M{"isActive": true}
	branchFilters.ProductBranchFilters = &product_filters.ProductBranchFilters{
		IsActive: &[]bool{true}[0],
	}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, false)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &branchFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}
