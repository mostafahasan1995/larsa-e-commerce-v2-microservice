package feature_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	feature_filters "larsa-e-commerce-microservice-v2/internal/modules/feature/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	feature_service "larsa-e-commerce-microservice-v2/internal/modules/feature/service"
	feature_utils "larsa-e-commerce-microservice-v2/internal/modules/feature/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeatureGroupRepository struct {
	*repo.BaseRepo[models.FeatureGroupModel, models.FeatureGroupResponse]
}

func NewFeatureGroupRepository(i *do.Injector) (feature_service.FeatureGroupRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &FeatureGroupRepository{
		BaseRepo: repo.NewBaseRepo[models.FeatureGroupModel, models.FeatureGroupResponse](db, enums.FEATURE_GROUP),
	}, nil
}

func (repo *FeatureGroupRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter feature_filters.FeatureGroupFilters) (*types.DataWithPagination[models.FeatureGroupResponse], *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *FeatureGroupRepository) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*models.FeatureGroupResponse, *helpers.ApiError) {
	pipeline = append(pipeline, feature_utils.PopulateFeaturesCount...)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}
