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

type CountryRepo struct {
	*base_repos.BaseRepo[config_location_models.Country, config_location_models.CountryResponse]
	db *mongo.Client
}

func NewCountryRepo(i *do.Injector) (config_location_services.CountryRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CountryRepo{
		BaseRepo: base_repos.NewBaseRepo[config_location_models.Country, config_location_models.CountryResponse](db, enums.COUNTRY),
		db:       db,
	}, nil
}

func (repo *CountryRepo) CreateCountry(ctx context.Context, country config_location_models.Country) (*config_location_models.CountryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedCountryCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_COUNTRY.GetCollectionName())
	seedCountry := config_location_models.Country{}

	seedCountryCollection.FindOne(ctx, bson.M{"name": country.Name}).Decode(&seedCountry)

	if seedCountry.Name == "" {
		return nil, helpers.DefaultError(errors.New("country not found"))
	}
	country.MultiPolygon = seedCountry.MultiPolygon
	country.Iso3 = seedCountry.Iso3

	return repo.Create(ctx, country)
}

func (repo *CountryRepo) UpdateCountry(ctx context.Context, updateCountry config_location_models.Country) (*config_location_models.CountryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	seedCountryCollection := repo.db.Database(cfg.Db).Collection(enums.SEED_COUNTRY.GetCollectionName())
	seedCountry := config_location_models.Country{}

	seedCountryCollection.FindOne(ctx, bson.M{"name": updateCountry.Name}).Decode(&seedCountry)

	if seedCountry.Name == "" {
		return nil, helpers.DefaultError(errors.New("country not found"))
	}

	updateCountry.MultiPolygon = seedCountry.MultiPolygon
	updateCountry.Iso3 = seedCountry.Iso3

	return repo.BaseRepo.UpdateById(ctx, updateCountry.Id, updateCountry, bson.M{})
}

func (repo *CountryRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.CountryFilters) (*types.DataWithPagination[config_location_models.CountryResponse], *helpers.ApiError) {
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

func (repo *CountryRepo) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.CountryResponse, *helpers.ApiError) {
	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	// Ensure pipeline is not nil before appending
	if pipeline == nil {
		pipeline = []bson.M{}
	}
	pipeline = append(pipeline, projectionStage)

	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *CountryRepo) GetAllSeedsCountries(ctx context.Context) ([]config_location_models.Country, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	collection := repo.db.Database(cfg.Db).Collection(enums.SEED_COUNTRY.GetCollectionName())

	opts := options.Find().SetProjection(bson.M{"multiPolygon": 0})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	var countries []config_location_models.Country
	if err := cursor.All(ctx, &countries); err != nil {
		return nil, helpers.DefaultError(err)
	}

	return countries, nil
}
