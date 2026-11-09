package product_repos

import (
	"context"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
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
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type WebReviewRepo struct {
	*base_repos.BaseRepo[models.Review, models.ReviewResponse]
}

func NewWebReviewRepo(i *do.Injector) (services.WebReviewRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebReviewRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Review, models.ReviewResponse](db, enums.REVIEW),
	}, nil
}

func (webReviewRepo *WebReviewRepo) Create(ctx context.Context, productId primitive.ObjectID, createReviewDto product_dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	review := product_utils.ConvertToReviewModel(productId, createReviewDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return webReviewRepo.BaseRepo.Create(ctx, review)
}

func (webReviewRepo *WebReviewRepo) UpdateByProductId(ctx context.Context, productId primitive.ObjectID, updateReviewDto product_dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	review := product_utils.ConvertToReviewModel(productId, updateReviewDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return webReviewRepo.UpdateOne(ctx, bson.M{"productId": productId, "userId": cfg.User.Id}, bson.M{"$set": review}, cfg.User.Id)
}

func (webReviewRepo *WebReviewRepo) DeleteByProductId(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError {
	return webReviewRepo.BaseRepo.DeleteOne(ctx, base_utils.AddSoftDeletedFilter(bson.M{"productId": productId}))
}

func (webReviewRepo *WebReviewRepo) GetByProductIdWithPagination(ctx context.Context, productId primitive.ObjectID, skip int64, limit int64, query string) (*models.ReviewsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	reviewFilters, _ := product_filters.NewWebReviewFilters(query)
	reviewFilters.ProductId = &productId

	// Build pipeline with filters
	pipeline := reviewFilters.BuildPipeline(matchStage, ctx, true)

	results, err := webReviewRepo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &reviewFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (webReviewRepo *WebReviewRepo) ValidateIfUserCanReview(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	productService := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](webReviewRepo.GetDB(), enums.PRODUCT)
	orderRepo := base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](webReviewRepo.GetDB(), enums.ORDER)

	product, err := productService.FindOne(ctx, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"_id": productId,
			}),
		},
	})
	if err != nil {
		return err
	}
	if product == nil {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Product not found",
		}
	}

	order, apiErr := orderRepo.FindOne(ctx, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"buyerInfo.id":       cfg.User.Id,
				"orderItems.itemId":  productId,
				"paymentInfo.status": order_enums.PAID,
			}),
		},
	})
	if apiErr != nil {
		return apiErr
	}

	if order == nil {
		return &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "You have not purchased this product",
		}
	}

	return nil
}
