package config_location_repos

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	config_location_services "larsa-e-commerce-microservice-v2/internal/modules/config/services/location"
)

type AreaRepo struct {
	*base_repos.BaseRepo[config_location_models.Area, config_location_models.AreaResponse]
	db *mongo.Client
}

func NewAreaRepo(i *do.Injector) (config_location_services.AreaRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &AreaRepo{
		BaseRepo: base_repos.NewBaseRepo[config_location_models.Area, config_location_models.AreaResponse](db, enums.AREA),
		db:       db,
	}, nil
}

func (repo *AreaRepo) CreateArea(ctx context.Context, area config_location_models.Area) (*config_location_models.AreaResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedAreaCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_AREA.GetCollectionName())
	seedArea := config_location_models.Area{}

	seedAreaCollection.FindOne(ctx, bson.M{"name": area.Name}).Decode(&seedArea)

	// If area exists in seed data, use the seed data values
	if seedArea.Name != "" {
		area.MultiPolygon = seedArea.MultiPolygon
		area.CityName = seedArea.CityName
	}
	// If area doesn't exist in seed data, proceed with creation using provided values
	// MultiPolygon and CityName should be provided in the input area

	return repo.Create(ctx, area)
}

func (repo *AreaRepo) UpdateArea(ctx context.Context, updateArea config_location_models.Area) (*config_location_models.AreaResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedAreaCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_AREA.GetCollectionName())
	seedArea := config_location_models.Area{}

	seedAreaCollection.FindOne(ctx, bson.M{"name": updateArea.Name}).Decode(&seedArea)

	// If area exists in seed data, use the seed data values
	if seedArea.Name != "" {
		updateArea.MultiPolygon = seedArea.MultiPolygon
		updateArea.CityName = seedArea.CityName
	}
	// If area doesn't exist in seed data, proceed with update using provided values
	// MultiPolygon and CityName should be provided in the input updateArea

	return repo.BaseRepo.UpdateById(ctx, updateArea.Id, updateArea, bson.M{})
}

func (repo *AreaRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.AreaFilters) (*types.DataWithPagination[config_location_models.AreaResponse], *helpers.ApiError) {
	pipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	pipeline = append(pipeline, projectionStage)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		Skip:           skip,
		Limit:          limit,
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
	})
}

func (repo *AreaRepo) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.AreaResponse, *helpers.ApiError) {
	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	// Ensure pipeline is not nil before appending
	if pipeline == nil {
		pipeline = []bson.M{}
	}
	pipeline = append(pipeline, projectionStage)

	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *AreaRepo) GetAllSeedsAreas(ctx context.Context, cityName string) ([]config_location_models.Area, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection := repo.db.Database(cfg.Db).Collection(enums.SEED_AREA.GetCollectionName())

	opts := options.Find().SetProjection(bson.M{"multiPolygon": 0})
	cursor, err := collection.Find(ctx, bson.M{"cityName": cityName}, opts)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	var areas []config_location_models.Area
	if err := cursor.All(ctx, &areas); err != nil {
		return nil, helpers.DefaultError(err)
	}

	return areas, nil
}
