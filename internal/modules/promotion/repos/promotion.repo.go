package promotion_repos

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	promotion_filters "larsa-e-commerce-microservice-v2/internal/modules/promotion/filters"
	promotion_models "larsa-e-commerce-microservice-v2/internal/modules/promotion/models"
	promotion_services "larsa-e-commerce-microservice-v2/internal/modules/promotion/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type promotionRepo struct {
	*base_repos.BaseRepo[promotion_models.Promotion, promotion_models.PromotionResponse]
}

func NewPromotionRepo(i *do.Injector) (promotion_services.PromotionRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &promotionRepo{
		BaseRepo: base_repos.NewBaseRepo[promotion_models.Promotion, promotion_models.PromotionResponse](db, enums.PROMOTION),
	}, nil
}

func (repo *promotionRepo) Create(ctx context.Context, promotion promotion_models.Promotion) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, promotion)
}

func (repo *promotionRepo) GetById(ctx context.Context, id primitive.ObjectID, filter *promotion_filters.PromotionFilters) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *promotionRepo) GetAll(ctx context.Context, pipeline []bson.M) ([]promotion_models.PromotionResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *promotionRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter promotion_filters.PromotionFilters) (*promotion_models.PromotionsWithPagination, *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)

	props := helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	}
	return repo.BaseRepo.GetWithPagination(ctx, props)
}

func (repo *promotionRepo) UpdateById(ctx context.Context, id primitive.ObjectID, promotion promotion_models.Promotion) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, promotion, bson.M{})
}

func (repo *promotionRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.DeleteById(ctx, id, bson.M{})
}

func (repo *promotionRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.BulkDelete(ctx, ids, bson.M{})
}

func (repo *promotionRepo) GetActivePromotionsForProductAndBranch(ctx context.Context, productId primitive.ObjectID, branchId *primitive.ObjectID, timezone *time.Location) (promotion_models.PromotionResponse, *helpers.ApiError) {
	var now time.Time
	if timezone != nil {
		now = time.Now().In(timezone)
	} else {
		now = time.Now().UTC()
	}

	// Build the product condition - either AllProducts is true OR the product is in ProductIds
	productCondition := bson.M{
		"$or": []bson.M{
			{"allProducts": true},
			{"productIds": productId},
		},
	}

	matchStage := bson.M{
		"deletedAt": bson.M{"$exists": false},
		"isActive":  true,
		"startDate": bson.M{"$lte": now},
		"endDate":   bson.M{"$gte": now},
		"$and":      []bson.M{productCondition},
	}

	if branchId != nil {
		matchStage["branchIds"] = bson.M{"$in": []primitive.ObjectID{*branchId}}
	}

	pipeline := []bson.M{
		{
			"$match": matchStage,
		},
		{
			"$sort": bson.M{
				"discountPercentage": -1,
			},
		},
		{
			"$limit": 1,
		},
	}

	promotions, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return promotion_models.PromotionResponse{}, err
	}

	if len(promotions) == 0 {
		return promotion_models.PromotionResponse{}, nil
	}

	return promotions[0], nil
}
