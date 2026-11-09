package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
)

type BranchProductRepo struct {
	*base_repos.BaseRepo[models.BranchProduct, models.BranchProductResponse]
}

func NewBranchProductRepo(i *do.Injector) (services.BranchProductRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &BranchProductRepo{
		BaseRepo: base_repos.NewBaseRepo[models.BranchProduct, models.BranchProductResponse](db, enums.BRANCH_PRODUCT),
	}, nil
}

func (repo *BranchProductRepo) Create(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) (*models.BranchProductResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	branchProduct := utils.ConvertToBranchProductModel(upsertProductBranchDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	return repo.BaseRepo.Create(ctx, branchProduct)
}

func (repo *BranchProductRepo) BulkUpsert(ctx context.Context, productBranches []models.BranchProduct) *helpers.ApiError {
	bulkWriteArgs := make([]base_repos.BulkWriteArgs[models.BranchProduct, models.BranchProductResponse], len(productBranches))
	for i, productBranch := range productBranches {

		filter := base_utils.AddSoftDeletedFilter(bson.M{
			"branchId":  productBranch.BranchId,
			"productId": productBranch.ProductId,
		})
		bulkWriteArgs[i] = base_repos.BulkWriteArgs[models.BranchProduct, models.BranchProductResponse]{
			Filter: filter,
			Upsert: productBranch,
		}

	}

	return repo.BaseRepo.BulkUpsert(ctx, bulkWriteArgs)
}

func (repo *BranchProductRepo) GetByBranchIdAndProductId(ctx context.Context, args product_utils.GetByBranchIdAndProductIdArgs, pipeline []bson.M) (*models.BranchProductResponse, *helpers.ApiError) {
	filter := base_utils.AddSoftDeletedFilter(bson.M{
		"branchId":  args.BranchId,
		"productId": args.ProductId,
	})

	pipeline = append(pipeline, bson.M{
		"$match": filter,
	})

	branchProduct, err := repo.BaseRepo.FindOne(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if branchProduct == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "No branch product found for the specified branch and product",
		}
	}
	return branchProduct, nil
}

func (repo *BranchProductRepo) GetAll(ctx context.Context, pipeline []bson.M) ([]models.BranchProductResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *BranchProductRepo) UpdateMany(ctx context.Context, filter bson.M, updateFields map[string]interface{}, userId primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.UpdateMany(ctx, filter, updateFields, userId)
}
