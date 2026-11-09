package plan_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	plan_filters "larsa-e-commerce-microservice-v2/internal/modules/plan/filters"
	plan_models "larsa-e-commerce-microservice-v2/internal/modules/plan/models"
	plan_service "larsa-e-commerce-microservice-v2/internal/modules/plan/service"
	plan_utils "larsa-e-commerce-microservice-v2/internal/modules/plan/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type PlanRepository struct {
	*base_repos.BaseRepo[plan_models.Plan, plan_models.Plan]
}

func NewPlanRepository(i *do.Injector) (plan_service.PlanRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &PlanRepository{
		BaseRepo: base_repos.NewBaseRepo[plan_models.Plan, plan_models.Plan](db, enums.PLAN),
	}, nil
}

func (repo *PlanRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter plan_filters.PlanFilters) (*types.DataWithPagination[plan_models.Plan], *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, false)
	pipeline = append(pipeline, plan_utils.PopulatePlanFeatures...)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *PlanRepository) GetAll(ctx context.Context, pipeline []bson.M) ([]plan_models.Plan, *helpers.ApiError) {
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *PlanRepository) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*plan_models.Plan, *helpers.ApiError) {
	pipeline = append(pipeline, plan_utils.PopulatePlanFeatures...)
	plan, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}
	return plan, nil
}
