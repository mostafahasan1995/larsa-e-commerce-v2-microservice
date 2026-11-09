package config_location_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils/location"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CityRepo interface {
	CreateCity(ctx context.Context, city config_location_models.City) (*config_location_models.CityResponse, *helpers.ApiError)
	UpdateCity(ctx context.Context, updateCity config_location_models.City) (*config_location_models.CityResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.CityResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.CityFilters) (*config_location_models.CitiesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]config_location_models.CityResponse, *helpers.ApiError)
	GetAllSeedsCities(ctx context.Context, countryName string) ([]config_location_models.City, *helpers.ApiError)
}

type CityService interface {
	Create(ctx context.Context, city *dtos.CreateCityDto) (*config_location_models.CityResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.CityResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]config_location_models.CityResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.CitiesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, city *dtos.CreateCityDto) (*config_location_models.CityResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	//Seeds
	GetAllSeedsCities(ctx context.Context, countryName string) ([]config_location_models.City, *helpers.ApiError)
}

type city struct {
	repo        CityRepo
	areaService AreaService
}

func NewCityService(i *do.Injector) (CityService, error) {
	areaService, err := do.Invoke[AreaService](i)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return &city{
		repo:        do.MustInvoke[CityRepo](i),
		areaService: areaService,
	}, nil
}

func (p *city) Create(ctx context.Context, cityDto *dtos.CreateCityDto) (*config_location_models.CityResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// // Check if city with same name exists
	// existingCities, err := p.GetAll(ctx, `{"name":"`+cityDto.Name+`"}`)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(existingCities) > 0 {
	// 	return nil, &helpers.ApiError{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    "City with this name already exists",
	// 	}
	// }

	cityModel := config_utils.ConvertToCityModel(cityDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	// Create the city first
	_createdCity, apiErr := p.repo.CreateCity(ctx, cityModel)
	if apiErr != nil {
		return nil, apiErr
	}

	// Get city by id to include the country name
	createdCity, apiErr := p.GetById(ctx, _createdCity.Id, `{"popFields":["countryId"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	// If IncludeAreas is true, create all areas for this city
	if cityDto.IncludeAreas != nil && *cityDto.IncludeAreas {

		// Get seed areas for this city
		seedAreas, apiErr := p.areaService.GetAllSeedsAreas(ctx, cityDto.Name)
		if apiErr != nil {
			return createdCity, apiErr
		}

		// Create each area
		for _, seedArea := range seedAreas {
			areaDto := &dtos.CreateAreaDto{
				Name:     seedArea.Name,
				IsActive: cityDto.IsActive,
				CityId:   createdCity.Id,
			}
			_, apiErr := p.areaService.Create(ctx, areaDto)
			if apiErr != nil {
				return createdCity, apiErr
			}
		}
	}

	return createdCity, nil
}

func (p *city) GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.CityResponse, *helpers.ApiError) {
	cityFilters, err := config_filters.NewCityFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetById(ctx, id, cityFilters.BuildPipeline(bson.M{}, ctx, true))
}

func (p *city) GetAll(ctx context.Context, query string) ([]config_location_models.CityResponse, *helpers.ApiError) {
	cityFilters, err := config_filters.NewCityFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	pipeline := cityFilters.BuildPipeline(bson.M{}, ctx, false)
	return p.repo.GetAll(ctx, pipeline)
}

func (p *city) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.CitiesWithPagination, *helpers.ApiError) {
	cityFilters, err := config_filters.NewCityFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetWithPagination(ctx, skip, limit, cityFilters)
}

func (p *city) UpdateById(ctx context.Context, id primitive.ObjectID, city *dtos.CreateCityDto) (*config_location_models.CityResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	cityModel := config_utils.ConvertToCityModel(city, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	cityModel.Id = id

	// Update the city first
	updatedCity, apiErr := p.repo.UpdateCity(ctx, cityModel)
	if apiErr != nil {
		return nil, apiErr
	}

	// If IsActive status is being updated, cascade the change to areas
	fmt.Println("city.IsActive", city.IsActive)
	if city.IsActive != nil {
		// Get all areas for this city
		areaQuery := `{"cityId":"` + id.Hex() + `"}`
		areas, areaErr := p.areaService.GetAll(ctx, areaQuery)
		if areaErr != nil {
			return updatedCity, areaErr
		}
		fmt.Println("areas", len(areas))

		// Update each area's status
		for _, area := range areas {
			fmt.Println("area.IsActive", area)
			areaDto := &dtos.CreateAreaDto{
				Name:     area.Name,
				IsActive: city.IsActive,
				CityId:   id,
			}
			if _, err := p.areaService.UpdateById(ctx, area.Id, areaDto); err != nil {
				continue
			}
		}
	}

	return updatedCity, nil
}

func (p *city) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	// First get all areas for this city
	areaQuery := `{"cityId":"` + id.Hex() + `"}`
	areas, areaErr := p.areaService.GetAll(ctx, areaQuery)
	if areaErr != nil {
		return areaErr
	}

	// Delete all areas
	for _, area := range areas {
		if err := p.areaService.DeleteById(ctx, area.Id); err != nil {
			continue
		}
	}

	// Finally delete the city
	return p.repo.DeleteById(ctx, id, bson.M{})
}

func (p *city) GetAllSeedsCities(ctx context.Context, countryName string) ([]config_location_models.City, *helpers.ApiError) {
	cities, err := p.repo.GetAllSeedsCities(ctx, countryName)
	if err != nil {
		return nil, err
	}

	return cities, nil
}
