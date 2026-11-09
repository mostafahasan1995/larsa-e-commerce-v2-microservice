package survey_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebWaiterService interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*survey_models.WaiterResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError)
}

type webWaiter struct {
	repo WaiterRepo
}

func NewWebWaiterService(i *do.Injector) (WebWaiterService, error) {
	return &webWaiter{
		repo: do.MustInvoke[WaiterRepo](i),
	}, nil
}

func (w *webWaiter) GetById(ctx context.Context, id primitive.ObjectID) (*survey_models.WaiterResponse, *helpers.ApiError) {
	return w.repo.GetById(ctx, id, "")
}

func (w *webWaiter) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError) {
	return w.repo.GetWithPagination(ctx, skip, limit, query)
}

func (w *webWaiter) GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError) {
	return w.repo.GetAll(ctx, query)
}
