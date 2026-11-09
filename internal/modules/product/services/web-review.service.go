package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebReviewRepo interface {
	Create(ctx context.Context, productId primitive.ObjectID, createReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError)
	UpdateByProductId(ctx context.Context, productId primitive.ObjectID, updateReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError)
	DeleteByProductId(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError
	GetByProductIdWithPagination(ctx context.Context, productId primitive.ObjectID, skip int64, limit int64, query string) (*models.ReviewsWithPagination, *helpers.ApiError)
	ValidateIfUserCanReview(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError
}

type WebReviewService interface {
	Create(ctx context.Context, productId primitive.ObjectID, createReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError)
	UpdateByProductId(ctx context.Context, productId primitive.ObjectID, updateReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError)
	DeleteByProductId(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError
	GetByProductIdWithPagination(ctx context.Context, productId primitive.ObjectID, skip int64, limit int64, query string) (*models.ReviewsWithPagination, *helpers.ApiError)
}

type webReview struct {
	repo WebReviewRepo
}

func NewWebReviewService(i *do.Injector) (WebReviewService, error) {
	return &webReview{
		repo: do.MustInvoke[WebReviewRepo](i),
	}, nil
}

func (webReview *webReview) Create(ctx context.Context, productId primitive.ObjectID, createReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError) {
	apiErr := webReview.repo.ValidateIfUserCanReview(ctx, productId)
	if apiErr != nil {
		return nil, apiErr
	}

	return webReview.repo.Create(ctx, productId, createReviewDto)
}

func (webReview *webReview) UpdateByProductId(ctx context.Context, productId primitive.ObjectID, updateReviewDto dtos.CreateReviewDto) (*models.ReviewResponse, *helpers.ApiError) {
	apiErr := webReview.repo.ValidateIfUserCanReview(ctx, productId)
	if apiErr != nil {
		return nil, apiErr
	}

	return webReview.repo.UpdateByProductId(ctx, productId, updateReviewDto)
}

func (webReview *webReview) DeleteByProductId(ctx context.Context, productId primitive.ObjectID) *helpers.ApiError {
	return webReview.repo.DeleteByProductId(ctx, productId)
}

func (webReview *webReview) GetByProductIdWithPagination(ctx context.Context, productId primitive.ObjectID, skip int64, limit int64, query string) (*models.ReviewsWithPagination, *helpers.ApiError) {
	return webReview.repo.GetByProductIdWithPagination(ctx, productId, skip, limit, query)
}
