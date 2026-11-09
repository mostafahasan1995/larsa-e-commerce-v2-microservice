package promotion_services

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	promotion_dtos "larsa-e-commerce-microservice-v2/internal/modules/promotion/dtos"
	promotion_filters "larsa-e-commerce-microservice-v2/internal/modules/promotion/filters"
	promotion_models "larsa-e-commerce-microservice-v2/internal/modules/promotion/models"
	promotion_utils "larsa-e-commerce-microservice-v2/internal/modules/promotion/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromotionRepo interface {
	Create(ctx context.Context, promotion promotion_models.Promotion) (*promotion_models.PromotionResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, filter *promotion_filters.PromotionFilters) (*promotion_models.PromotionResponse, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]promotion_models.PromotionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter promotion_filters.PromotionFilters) (*promotion_models.PromotionsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, promotion promotion_models.Promotion) (*promotion_models.PromotionResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	GetActivePromotionsForProductAndBranch(ctx context.Context, productId primitive.ObjectID, branchId *primitive.ObjectID, timezone *time.Location) (promotion_models.PromotionResponse, *helpers.ApiError)
}

type PromotionService interface {
	Create(ctx context.Context, dto *promotion_dtos.CreatePromotionDTO) (*promotion_models.PromotionResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*promotion_models.PromotionResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]promotion_models.PromotionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*promotion_models.PromotionsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, dto *promotion_dtos.UpdatePromotionDTO) (*promotion_models.PromotionResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	GetActivePromotionsForProductAndBranch(ctx context.Context, productId primitive.ObjectID, branchId *primitive.ObjectID, timezone *time.Location) (promotion_models.PromotionResponse, *helpers.ApiError)
}

type promotion struct {
	repo PromotionRepo
}

func NewPromotionService(i *do.Injector) (PromotionService, error) {
	return &promotion{
		repo: do.MustInvoke[PromotionRepo](i),
	}, nil
}

func (s *promotion) Create(ctx context.Context, dto *promotion_dtos.CreatePromotionDTO) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	model := promotion_utils.ConvertCreateDTOToModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	return s.repo.Create(ctx, model)
}

func (s *promotion) GetById(ctx context.Context, id primitive.ObjectID, query string) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	filters, err := promotion_filters.NewPromotionFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	return s.repo.GetById(ctx, id, filters)
}

func (s *promotion) GetAll(ctx context.Context, query string) ([]promotion_models.PromotionResponse, *helpers.ApiError) {
	filters, err := promotion_filters.NewPromotionFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := filters.BuildPipeline([]bson.M{}, ctx, false)
	return s.repo.GetAll(ctx, pipeline)
}

func (s *promotion) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*promotion_models.PromotionsWithPagination, *helpers.ApiError) {
	filters, err := promotion_filters.NewPromotionFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	return s.repo.GetWithPagination(ctx, skip, limit, *filters)
}

func (s *promotion) UpdateById(ctx context.Context, id primitive.ObjectID, dto *promotion_dtos.UpdatePromotionDTO) (*promotion_models.PromotionResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	model := promotion_utils.ConvertUpdateDTOToModel(dto, id, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return s.repo.UpdateById(ctx, id, model)
}

func (s *promotion) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id)
}

func (s *promotion) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return s.repo.BulkDelete(ctx, ids)
}

func (s *promotion) GetActivePromotionsForProductAndBranch(ctx context.Context, productId primitive.ObjectID, branchId *primitive.ObjectID, timezone *time.Location) (promotion_models.PromotionResponse, *helpers.ApiError) {
	return s.repo.GetActivePromotionsForProductAndBranch(ctx, productId, branchId, timezone)
}
