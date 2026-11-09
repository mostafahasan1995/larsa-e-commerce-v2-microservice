package plan_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
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

type WebPlanRepo struct {
	*base_repos.BaseRepo[plan_models.Plan, plan_models.PlanResponse]
}

func NewWebPlanRepo(i *do.Injector) (plan_service.WebPlanRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebPlanRepo{
		BaseRepo: base_repos.NewBaseRepo[plan_models.Plan, plan_models.PlanResponse](db, enums.PLAN),
	}, nil
}

func (repo *WebPlanRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*plan_models.PlanResponse, *helpers.ApiError) {

	// Create filters from query parameters
	planFilters, _ := plan_filters.NewPlanFilters(query)

	// Build pipeline with filters
	matchStage := []bson.M{}

	pipeline := planFilters.BuildPipeline(matchStage, ctx, true)
	pipeline = append(pipeline, plan_utils.PopulatePlanFeatures...)

	plan, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return plan, nil
}

func (repo *WebPlanRepo) GetAll(ctx context.Context, query string) ([]plan_models.PlanResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := []bson.M{}

	// Create filters from query parameters
	planFilters, _ := plan_filters.NewPlanFilters(query)

	// Build pipeline with filters
	pipeline := planFilters.BuildPipeline(matchStage, ctx, true)
	pipeline = append(pipeline, plan_utils.PopulatePlanFeatures...)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WebPlanRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*plan_models.PlansWithPagination, *helpers.ApiError) {

	// Base match stage with deletedAt filter
	matchStage := []bson.M{}

	// Create filters from query parameters
	planFilters, _ := plan_filters.NewPlanFilters(query)

	// Build pipeline with filters
	pipeline := planFilters.BuildPipeline(matchStage, ctx, false)
	pipeline = append(pipeline, plan_utils.PopulatePlanFeatures...)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &planFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
