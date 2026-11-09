package coupon_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	filters "larsa-e-commerce-microservice-v2/internal/modules/coupon/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/coupon/utils"
)

type CouponRepo struct {
	*base_repos.BaseRepo[models.Coupon, models.CouponResponse]
}

func NewCouponRepo(i *do.Injector) (services.CouponRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CouponRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Coupon, models.CouponResponse](db, enums.COUPON),
	}, nil
}

func (repo *CouponRepo) Create(ctx context.Context, createCouponDto *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	coupon := utils.ConvertToCouponModel(*createCouponDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	if coupon.IsActive == nil && coupon.Type == coupon_enums.CouponType(coupon_enums.MANUAL) {
		coupon.IsActive = &[]bool{true}[0]
	}

	return repo.BaseRepo.Create(ctx, coupon)
}

func (repo *CouponRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CouponResponse, *helpers.ApiError) {

	// Create filters from query parameters
	couponFilters, _ := filters.NewCouponFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := couponFilters.BuildPipeline(matchStage, ctx, true)

	coupon, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return coupon, nil
}

func (repo *CouponRepo) GetAll(ctx context.Context, query string) ([]models.CouponResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	couponFilters, _ := filters.NewCouponFilters(query)

	// Build pipeline with filters
	pipeline := couponFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CouponRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CouponsWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	couponFilters, _ := filters.NewCouponFilters(query)

	// Build pipeline with filters
	pipeline := couponFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &couponFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *CouponRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateCouponDto *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	coupon := utils.ConvertToCouponModel(*updateCouponDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	// Marshal coupon to BSON to get update document
	couponBytes, err := bson.Marshal(coupon)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	var updateDoc bson.M
	if err := bson.Unmarshal(couponBytes, &updateDoc); err != nil {
		return nil, helpers.DefaultError(err)
	}

	metadataFieldsToExclude := []string{"_id", "createdAt", "createdBy", "deletedAt", "deletedBy"}
	for _, field := range metadataFieldsToExclude {
		delete(updateDoc, field)
	}

	// Handle OrderType based on Nullable behavior:
	// - If provided with value: use the value (already in updateDoc from marshal)
	// - If explicitly null: set to null in DB
	// - If not provided: exclude from update (don't change existing value)
	if updateCouponDto.OrderType.IsProvided() {
		if updateCouponDto.OrderType.IsExplicitlyNull() {
			updateDoc["orderType"] = nil
		}
		// If provided with value, it's already in updateDoc from marshal, so no action needed
	} else {
		// If not provided, remove from updateDoc to preserve existing value
		delete(updateDoc, "orderType")
	}

	// Convert bson.M to map[string]interface{} for UpdateFieldsById
	updateFields := make(map[string]interface{})
	for k, v := range updateDoc {
		updateFields[k] = v
	}

	// Use UpdateFieldsById for cleaner update logic
	return repo.BaseRepo.UpdateFieldsById(ctx, id, updateFields, cfg.User.Id)
}

func (repo *CouponRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *CouponRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
