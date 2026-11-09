package survey_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_utils "larsa-e-commerce-microservice-v2/internal/modules/survey/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WaiterRepo interface {
	Create(ctx context.Context, waiter survey_models.WaiterModel) (*survey_models.WaiterResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, waiter survey_models.WaiterModel, filter bson.M) (*survey_models.WaiterResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.WaiterResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError)
}

type WaiterService interface {
	Create(ctx context.Context, dto *survey_dtos.CreateWaiterDTO) (*survey_models.WaiterResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateWaiterDTO) (*survey_models.WaiterResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.WaiterResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError)
}

type waiter struct {
	repo WaiterRepo
}

func NewWaiterService(i *do.Injector) (WaiterService, error) {
	return &waiter{
		repo: do.MustInvoke[WaiterRepo](i),
	}, nil
}

func (w *waiter) Create(ctx context.Context, dto *survey_dtos.CreateWaiterDTO) (*survey_models.WaiterResponse, *helpers.ApiError) {
	waiterModel, err := survey_utils.ConvertCreateWaiterDTOToModel(dto, ctx)
	if err != nil {
		return nil, err
	}

	return w.repo.Create(ctx, waiterModel)
}

func (w *waiter) Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateWaiterDTO) (*survey_models.WaiterResponse, *helpers.ApiError) {
	waiterModel, err := survey_utils.ConvertUpdateWaiterDTOToModel(dto, id, ctx)
	if err != nil {
		return nil, err
	}

	return w.repo.UpdateById(ctx, id, waiterModel, bson.M{})
}

func (w *waiter) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return w.repo.DeleteById(ctx, id, bson.M{})
}

func (w *waiter) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.WaiterResponse, *helpers.ApiError) {
	return w.repo.GetById(ctx, id, query)
}

func (w *waiter) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.WaitersWithPagination, *helpers.ApiError) {
	return w.repo.GetWithPagination(ctx, skip, limit, query)
}

func (w *waiter) GetAll(ctx context.Context, query string) ([]survey_models.WaiterResponse, *helpers.ApiError) {
	return w.repo.GetAll(ctx, query)
}
