package plan_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	plan_dto "larsa-e-commerce-microservice-v2/internal/modules/plan/dto"
	plan_filters "larsa-e-commerce-microservice-v2/internal/modules/plan/filters"
	plan_models "larsa-e-commerce-microservice-v2/internal/modules/plan/models"
	plan_utils "larsa-e-commerce-microservice-v2/internal/modules/plan/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanRepo interface {
	Create(ctx context.Context, plan plan_models.Plan) (*plan_models.Plan, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*plan_models.Plan, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter plan_filters.PlanFilters) (*types.DataWithPagination[plan_models.Plan], *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]plan_models.Plan, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, plan plan_models.Plan, filter bson.M) (*plan_models.Plan, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
}

type PlanService interface {
	Create(ctx context.Context, dto *plan_dto.CreatePlanDTO) (*plan_models.Plan, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*plan_models.Plan, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]plan_models.Plan, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[plan_models.Plan], *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *plan_dto.UpdatePlanDTO) (*plan_models.Plan, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
}

type plan struct {
	repo PlanRepo
}

func NewPlanService(i *do.Injector) (PlanService, error) {
	return &plan{
		repo: do.MustInvoke[PlanRepo](i),
	}, nil
}

func (s *plan) Create(ctx context.Context, dto *plan_dto.CreatePlanDTO) (*plan_models.Plan, *helpers.ApiError) {
	plan, apiErr := plan_utils.ConvertCreatePlanDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	result, apiErr := s.repo.Create(ctx, plan)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *plan) Update(ctx context.Context, id primitive.ObjectID, dto *plan_dto.UpdatePlanDTO) (*plan_models.Plan, *helpers.ApiError) {
	plan, apiErr := plan_utils.ConvertUpdateDTOToModel(dto, id, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	result, apiErr := s.repo.UpdateById(ctx, id, plan, bson.M{"deletedAt": bson.M{"$exists": false}})
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *plan) GetById(ctx context.Context, id primitive.ObjectID) (*plan_models.Plan, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, []bson.M{{"$match": bson.M{"deletedAt": bson.M{"$exists": false}}}})
}

func (s *plan) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
}

func (s *plan) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[plan_models.Plan], *helpers.ApiError) {
	filter, err := plan_filters.NewPlanFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *plan) GetAll(ctx context.Context, query string) ([]plan_models.Plan, *helpers.ApiError) {
	filter, err := plan_filters.NewPlanFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := []bson.M{}
	pipeline = filter.BuildPipeline(pipeline, ctx, true)

	return s.repo.GetAll(ctx, pipeline)
}
