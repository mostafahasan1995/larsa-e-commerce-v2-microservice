package config_location_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	util "larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils/location"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CountryRepo interface {
	CreateCountry(ctx context.Context, country config_location_models.Country) (*config_location_models.CountryResponse, *helpers.ApiError)
	UpdateCountry(ctx context.Context, updateCountry config_location_models.Country) (*config_location_models.CountryResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.CountryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.CountryFilters) (*config_location_models.CountriesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]config_location_models.CountryResponse, *helpers.ApiError)
	GetAllSeedsCountries(ctx context.Context) ([]config_location_models.Country, *helpers.ApiError)
}

type CountryService interface {
	Create(ctx context.Context, country *dtos.CreateCountryDto) (*config_location_models.CountryResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.CountryResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]config_location_models.CountryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.CountriesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, country *dtos.CreateCountryDto) (*config_location_models.CountryResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	//Seeds
	GetAllSeedsCountries(ctx context.Context) ([]config_location_models.Country, *helpers.ApiError)
	// Get all countries with their locations
	GetAllWithLocations(ctx context.Context) ([]config_location_models.CountryWithLocationsResponse, *helpers.ApiError)
}

type country struct {
	repo        CountryRepo
	areaService AreaService
	cityService CityService
}

func NewCountryService(i *do.Injector) (CountryService, error) {
	areaService, err := do.Invoke[AreaService](i)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	cityService, err := do.Invoke[CityService](i)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return &country{
		repo:        do.MustInvoke[CountryRepo](i),
		areaService: areaService,
		cityService: cityService,
	}, nil
}

func (p *country) Create(ctx context.Context, countryDto *dtos.CreateCountryDto) (*config_location_models.CountryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// // Check if country with same name exists
	// existingCountries, err := p.GetAll(ctx, `{"name":"`+countryDto.Name+`"}`)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(existingCountries) > 0 {
	// 	return nil, &helpers.ApiError{
	// 		StatusCode: http.StatusBadRequest,
	// 		Message:    "Country with this name already exists",
	// 	}
	// }

	countryModel := config_utils.ConvertToCountryModel(countryDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	// Create the country first
	createdCountry, apiErr := p.repo.CreateCountry(ctx, countryModel)
	if apiErr != nil {
		return nil, apiErr
	}

	var cities []config_location_models.CityResponse
	// If IncludeCities is true, create all cities for this country
	if countryDto.IncludeCities != nil && *countryDto.IncludeCities {
		seedCities, apiErr := p.cityService.GetAllSeedsCities(ctx, countryDto.Name)
		if apiErr != nil {
			return createdCountry, apiErr
		}
		for _, seedCity := range seedCities {
			cityDto := &dtos.CreateCityDto{
				Name:      seedCity.Name,
				IsActive:  countryDto.IsActive,
				CountryId: createdCountry.Id,
			}
			city, apiErr := p.cityService.Create(ctx, cityDto)
			if apiErr != nil {
				return createdCountry, apiErr
			}
			cities = append(cities, *city)
		}
	}

	// If IncludeAreas is true, create all areas for this country
	if countryDto.IncludeAreas != nil && *countryDto.IncludeAreas {
		for _, city := range cities {
			seedAreas, apiErr := p.areaService.GetAllSeedsAreas(ctx, city.Name)
			if apiErr != nil {
				return createdCountry, apiErr
			}

			for _, seedArea := range seedAreas {
				areaDto := &dtos.CreateAreaDto{
					Name:     seedArea.Name,
					IsActive: countryDto.IsActive,
					CityId:   city.Id,
				}
				_, apiErr := p.areaService.Create(ctx, areaDto)
				if apiErr != nil {
					return createdCountry, apiErr
				}
			}
		}
	}

	return createdCountry, nil
}

func (p *country) GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.CountryResponse, *helpers.ApiError) {
	countryFilters, err := config_filters.NewCountryFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetById(ctx, id, countryFilters.BuildPipeline(bson.M{}, ctx, true))
}

func (p *country) GetAll(ctx context.Context, query string) ([]config_location_models.CountryResponse, *helpers.ApiError) {
	countryFilters, err := config_filters.NewCountryFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetAll(ctx, countryFilters.BuildPipeline(bson.M{}, ctx, false))
}

func (p *country) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.CountriesWithPagination, *helpers.ApiError) {
	countryFilters, err := config_filters.NewCountryFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetWithPagination(ctx, skip, limit, countryFilters)
}

func (p *country) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	// First get all cities for this country
	cityQuery := `{"countryId":"` + id.Hex() + `"}`
	cities, cityErr := p.cityService.GetAll(ctx, cityQuery)
	if cityErr != nil {
		return cityErr
	}

	// Delete all areas for each city
	for _, city := range cities {
		areaQuery := `{"cityId":"` + city.Id.Hex() + `"}`
		areas, areaErr := p.areaService.GetAll(ctx, areaQuery)
		if areaErr != nil {
			continue
		}

		// Delete each area
		for _, area := range areas {
			if err := p.areaService.DeleteById(ctx, area.Id); err != nil {
				continue
			}
		}

		// Delete the city
		if err := p.cityService.DeleteById(ctx, city.Id); err != nil {
			continue
		}
	}

	// Finally delete the country
	return p.repo.DeleteById(ctx, id, bson.M{})
}

func (p *country) UpdateById(ctx context.Context, id primitive.ObjectID, country *dtos.CreateCountryDto) (*config_location_models.CountryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	countryModel := config_utils.ConvertToCountryModel(country, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	countryModel.Id = id

	// Update the country first
	updatedCountry, apiErr := p.repo.UpdateCountry(ctx, countryModel)
	if apiErr != nil {
		return nil, apiErr
	}

	// If IsActive status is being updated, cascade the change to cities and areas
	if country.IsActive != nil {
		// Get all cities for this country
		cityQuery := `{"countryId":"` + id.Hex() + `"}`
		cities, cityErr := p.cityService.GetAll(ctx, cityQuery)
		if cityErr != nil {
			return updatedCountry, cityErr
		}

		// Update each city and its areas
		for _, city := range cities {
			// Update city status
			cityDto := &dtos.CreateCityDto{
				Name:      city.Name,
				IsActive:  country.IsActive,
				CountryId: id,
			}
			if _, err := p.cityService.UpdateById(ctx, city.Id, cityDto); err != nil {
				continue
			}

			// Get and update all areas for this city
			areaQuery := `{"cityId":"` + city.Id.Hex() + `"}`
			areas, areaErr := p.areaService.GetAll(ctx, areaQuery)
			if areaErr != nil {
				continue
			}

			// Update each area's status
			for _, area := range areas {
				areaDto := &dtos.CreateAreaDto{
					Name:     area.Name,
					IsActive: country.IsActive,
					CityId:   city.Id,
				}
				if _, err := p.areaService.UpdateById(ctx, area.Id, areaDto); err != nil {
					continue
				}
			}
		}
	}

	return updatedCountry, nil
}

func (p *country) GetAllSeedsCountries(ctx context.Context) ([]config_location_models.Country, *helpers.ApiError) {
	countries, err := p.repo.GetAllSeedsCountries(ctx)
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (p *country) GetAllWithLocations(ctx context.Context) ([]config_location_models.CountryWithLocationsResponse, *helpers.ApiError) {
	countries, apiErr := p.GetAll(ctx, `{"isActive":true}`)
	if apiErr != nil {
		return nil, apiErr
	}
	result := make([]config_location_models.CountryWithLocationsResponse, 0, len(countries))
	for _, countryRes := range countries {

		cityQuery := `{"countryId":"` + countryRes.Id.Hex() + `","isActive":true}`
		cities, cityErr := p.cityService.GetAll(ctx, cityQuery)
		if cityErr != nil {
			continue
		}

		citiesWithAreas := make([]config_location_models.CityWithAreasResponse, 0, len(cities))

		for _, cityRes := range cities {
			areaQuery := `{"cityId":"` + cityRes.Id.Hex() + `","isActive":true}`
			areas, areaErr := p.areaService.GetAll(ctx, areaQuery)
			if areaErr != nil {
				continue
			}
			areasWithAreas := make([]config_location_models.Area, 0, len(areas))
			for _, areaRes := range areas {
				areasWithAreas = append(areasWithAreas, config_location_models.Area{
					BaseArea: config_location_models.BaseArea{
						MetaDataBaseModel: areaRes.MetaDataBaseModel,
						Name:              areaRes.Name,
						CityId:            cityRes.Id,
					},
				})
			}
			citiesWithAreas = append(citiesWithAreas, config_location_models.CityWithAreasResponse{
				City: config_location_models.City{
					BaseCity: config_location_models.BaseCity{
						MetaDataBaseModel: cityRes.MetaDataBaseModel,
						Name:              cityRes.Name,
						CountryId:         countryRes.Id,
					},
				},
				Areas: areasWithAreas,
			})
		}

		result = append(result, config_location_models.CountryWithLocationsResponse{
			Country: config_location_models.Country{
				BaseCountry: config_location_models.BaseCountry{
					MetaDataBaseModel: countryRes.MetaDataBaseModel,
					Name:              countryRes.Name,
					Iso3:              countryRes.Iso3,
				},
			},
			Cities: citiesWithAreas,
		})
	}

	return result, nil
}
