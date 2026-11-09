package plan_service

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	plan_models "larsa-e-commerce-microservice-v2/internal/modules/plan/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebPlanRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*plan_models.PlanResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]plan_models.PlanResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*plan_models.PlansWithPagination, *helpers.ApiError)
}

type WebPlanService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*plan_models.PlanResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]plan_models.PlanResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*plan_models.PlansWithPagination, *helpers.ApiError)
}

type webPlan struct {
	repo WebPlanRepo
}

func NewWebPlanService(i *do.Injector) (WebPlanService, error) {
	return &webPlan{
		repo: do.MustInvoke[WebPlanRepo](i),
	}, nil
}

func (p *webPlan) GetById(ctx context.Context, id primitive.ObjectID, query string) (*plan_models.PlanResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *webPlan) GetAll(ctx context.Context, query string) ([]plan_models.PlanResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *webPlan) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*plan_models.PlansWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}
