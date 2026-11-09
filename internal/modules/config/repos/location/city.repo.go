package config_location_repos

import (
	"context"
	"errors"

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

type CityRepo struct {
	*base_repos.BaseRepo[config_location_models.City, config_location_models.CityResponse]
	db *mongo.Client
}

func NewCityRepo(i *do.Injector) (config_location_services.CityRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CityRepo{
		BaseRepo: base_repos.NewBaseRepo[config_location_models.City, config_location_models.CityResponse](db, enums.CITY),
		db:       db,
	}, nil
}

func (repo *CityRepo) CreateCity(ctx context.Context, city config_location_models.City) (*config_location_models.CityResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedCityCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_CITY.GetCollectionName())
	seedCity := config_location_models.City{}

	seedCityCollection.FindOne(ctx, bson.M{"name": city.Name}).Decode(&seedCity)

	if seedCity.Name == "" {
		return nil, helpers.DefaultError(errors.New("city not found in seed data"))
	}

	city.MultiPolygon = seedCity.MultiPolygon
	city.CountryName = seedCity.CountryName

	return repo.Create(ctx, city)
}

func (repo *CityRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.CityFilters) (*types.DataWithPagination[config_location_models.CityResponse], *helpers.ApiError) {
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

func (repo *CityRepo) UpdateCity(ctx context.Context, updateCity config_location_models.City) (*config_location_models.CityResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedCityCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_CITY.GetCollectionName())
	seedCity := config_location_models.City{}

	seedCityCollection.FindOne(ctx, bson.M{"name": updateCity.Name}).Decode(&seedCity)

	if seedCity.Name == "" {
		return nil, helpers.DefaultError(errors.New("city not found"))
	}

	updateCity.MultiPolygon = seedCity.MultiPolygon
	updateCity.CountryName = seedCity.CountryName

	return repo.BaseRepo.UpdateById(ctx, updateCity.Id, updateCity, bson.M{})
}

func (repo *CityRepo) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.CityResponse, *helpers.ApiError) {
	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	// Ensure pipeline is not nil before appending
	if pipeline == nil {
		pipeline = []bson.M{}
	}
	pipeline = append(pipeline, projectionStage)

	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *CityRepo) GetAllSeedsCities(ctx context.Context, countryName string) ([]config_location_models.City, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection := repo.db.Database(cfg.Db).Collection(enums.SEED_CITY.GetCollectionName())

	opts := options.Find().SetProjection(bson.M{"multiPolygon": 0})
	cursor, err := collection.Find(ctx, bson.M{"countryName": countryName}, opts)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	var cities []config_location_models.City
	if err := cursor.All(ctx, &cities); err != nil {
		return nil, helpers.DefaultError(err)
	}

	return cities, nil
}
