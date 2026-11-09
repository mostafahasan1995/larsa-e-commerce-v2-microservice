package config_location_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO: add location web repo
type LocationWebService interface {
	GetAllCountries(ctx context.Context) ([]config_location_models.CountryResponse, *helpers.ApiError)
	GetCitiesByCountry(ctx context.Context, query string) ([]config_location_models.CityResponse, *helpers.ApiError)
	GetAreasByCity(ctx context.Context, query string) ([]config_location_models.AreaResponse, *helpers.ApiError)
}

type locationWeb struct {
	countryRepo CountryRepo
	cityRepo    CityRepo
	areaRepo    AreaRepo
}

func NewLocationWebService(i *do.Injector) (LocationWebService, error) {
	return &locationWeb{
		countryRepo: do.MustInvoke[CountryRepo](i),
		cityRepo:    do.MustInvoke[CityRepo](i),
		areaRepo:    do.MustInvoke[AreaRepo](i),
	}, nil
}

func (s *locationWeb) GetAllCountries(ctx context.Context) ([]config_location_models.CountryResponse, *helpers.ApiError) {
	isActive := true
	filter := config_filters.CountryFilters{
		IsActive: &isActive,
	}
	pipeline := filter.BuildPipeline(bson.M{}, ctx, false)

	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	pipeline = append(pipeline, projectionStage)

	return s.countryRepo.GetAll(ctx, pipeline)
}

func (s *locationWeb) GetCitiesByCountry(ctx context.Context, query string) ([]config_location_models.CityResponse, *helpers.ApiError) {
	isActive := true
	cityFilters, _ := config_filters.NewCityFilters(query)
	cityFilters.IsActive = &isActive
	// cityFilters.BaseFilters.PopFields = &[]string{"branches"}
	// filter := config_filters.CityFilters{
	// 	CountryName: &countryName,
	// 	IsActive:    &isActive,
	// 	BaseFilters: filters.BaseFilters{
	// 		PopFields: &[]string{"branches"},
	// 	},
	// }
	pipeline := cityFilters.BuildPipeline(bson.M{}, ctx, true)

	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	pipeline = append(pipeline, projectionStage)

	return s.cityRepo.GetAll(ctx, pipeline)
}

func (s *locationWeb) GetAreasByCity(ctx context.Context, query string) ([]config_location_models.AreaResponse, *helpers.ApiError) {
	fmt.Println("query", query)
	isActive := true
	areaFilters, _ := config_filters.NewAreaFilters(query)
	areaFilters.IsActive = &isActive
	pipeline := areaFilters.BuildPipeline(bson.M{}, ctx, false)

	// Add projection to exclude multiPolygon
	projectionStage := bson.M{"$project": bson.M{"multiPolygon": 0}}
	pipeline = append(pipeline, projectionStage)

	return s.areaRepo.GetAll(ctx, pipeline)
}
