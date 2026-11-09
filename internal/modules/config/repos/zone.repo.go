package config_repos

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
)

type ZoneRepo struct {
	*base_repos.BaseRepo[config_models.Zone, config_models.ZoneResponse]
}

func NewZoneRepo(i *do.Injector) (config_services.ZoneRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ZoneRepo{
		BaseRepo: base_repos.NewBaseRepo[config_models.Zone, config_models.ZoneResponse](db, enums.ZONE),
	}, nil
}

func (repo *ZoneRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.ZoneFilters) (*types.DataWithPagination[config_models.ZoneResponse], *helpers.ApiError) {
	pipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		Skip:           skip,
		Limit:          limit,
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
	})
}

func (repo *ZoneRepo) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_models.ZoneResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *ZoneRepo) UpdateById(ctx context.Context, id primitive.ObjectID, model *config_models.Zone) (*config_models.ZoneResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, *model, bson.M{})
}

func (repo *ZoneRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.DeleteById(ctx, id, bson.M{})
}
