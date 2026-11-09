package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchProductRepo interface {
	Create(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) (*models.BranchProductResponse, *helpers.ApiError)
	BulkUpsert(ctx context.Context, productBranches []models.BranchProduct) *helpers.ApiError
	GetByBranchIdAndProductId(ctx context.Context, args product_utils.GetByBranchIdAndProductIdArgs, pipeline []bson.M) (*models.BranchProductResponse, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]models.BranchProductResponse, *helpers.ApiError)
	UpdateMany(ctx context.Context, filter bson.M, updateFields map[string]interface{}, userId primitive.ObjectID) *helpers.ApiError
}

type BranchProductService interface {
	Create(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) (*models.BranchProductResponse, *helpers.ApiError)
	BulkUpsert(ctx context.Context, productBranches []models.BranchProduct) *helpers.ApiError
	GetByBranchIdAndProductId(ctx context.Context, args product_utils.GetByBranchIdAndProductIdArgs, pipeline []bson.M) (*models.BranchProductResponse, *helpers.ApiError)
}

type branchProduct struct {
	repo BranchProductRepo
}

func NewBranchProductService(i *do.Injector) (BranchProductService, error) {
	return &branchProduct{
		repo: do.MustInvoke[BranchProductRepo](i),
	}, nil
}

func (p *branchProduct) Create(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) (*models.BranchProductResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, upsertProductBranchDto)
}

func (p *branchProduct) BulkUpsert(ctx context.Context, productBranches []models.BranchProduct) *helpers.ApiError {
	return p.repo.BulkUpsert(ctx, productBranches)
}

func (p *branchProduct) GetByBranchIdAndProductId(ctx context.Context, args product_utils.GetByBranchIdAndProductIdArgs, pipeline []bson.M) (*models.BranchProductResponse, *helpers.ApiError) {
	return p.repo.GetByBranchIdAndProductId(ctx, args, pipeline)
}
