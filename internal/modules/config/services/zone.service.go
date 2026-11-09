package config_services

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	config_location_services "larsa-e-commerce-microservice-v2/internal/modules/config/services/location"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"

	"larsa-e-commerce-microservice-v2/internal/common/types"
)

type ZoneService interface {
	Create(ctx context.Context, dto *config_dtos.CreateZoneDto) (*config_models.ZoneResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_models.ZoneResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]config_models.ZoneResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[config_models.ZoneResponse], *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, dto *config_dtos.CreateZoneDto) (*config_models.ZoneResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	FindZonesContainingPoint(ctx context.Context, latitude, longitude float64) ([]config_models.ZoneResponse, *helpers.ApiError)
}

type ZoneRepo interface {
	Create(ctx context.Context, model config_models.Zone) (*config_models.ZoneResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_models.ZoneResponse, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]config_models.ZoneResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.ZoneFilters) (*types.DataWithPagination[config_models.ZoneResponse], *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, model *config_models.Zone) (*config_models.ZoneResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
}

type zoneService struct {
	zoneRepo       ZoneRepo
	countryService config_location_services.CountryService
	cityService    config_location_services.CityService
	areaService    config_location_services.AreaService
}

func NewZoneService(i *do.Injector) (ZoneService, error) {
	return &zoneService{
		zoneRepo:       do.MustInvoke[ZoneRepo](i),
		countryService: do.MustInvoke[config_location_services.CountryService](i),
		cityService:    do.MustInvoke[config_location_services.CityService](i),
		areaService:    do.MustInvoke[config_location_services.AreaService](i),
	}, nil
}

func (s *zoneService) Create(ctx context.Context, dto *config_dtos.CreateZoneDto) (*config_models.ZoneResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	zone := config_utils.ConvertToZoneModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return s.zoneRepo.Create(ctx, zone)
}

func (s *zoneService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_models.ZoneResponse, *helpers.ApiError) {
	filters, err := config_filters.NewZoneFilters(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	pipeline := filters.BuildPipeline(bson.M{}, ctx, true)
	return s.zoneRepo.GetById(ctx, id, pipeline)
}

func (s *zoneService) GetAll(ctx context.Context, query string) ([]config_models.ZoneResponse, *helpers.ApiError) {
	filters, err := config_filters.NewZoneFilters(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	pipeline := filters.BuildPipeline(bson.M{}, ctx, true)
	result, apiErr := s.zoneRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

func (s *zoneService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[config_models.ZoneResponse], *helpers.ApiError) {
	filters, err := config_filters.NewZoneFilters(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	return s.zoneRepo.GetWithPagination(ctx, skip, limit, filters)
}

func (s *zoneService) UpdateById(ctx context.Context, id primitive.ObjectID, dto *config_dtos.CreateZoneDto) (*config_models.ZoneResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	zone := config_utils.ConvertToZoneModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return s.zoneRepo.UpdateById(ctx, id, &zone)
}

func (s *zoneService) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.zoneRepo.DeleteById(ctx, id)
}

//We remove this function because we will only link zone with polygon
/**
func (s *zoneService) FindZonesContainingLocation(ctx context.Context, locationType config_enums.LocationType, locationName string) ([]config_models.ZoneResponse, *helpers.ApiError) {
	// Get all zones
	zones, apiErr := s.GetAll(ctx, `{"popFields":["countryIds","cityIds","areaIds"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	zonesResponse := make([]config_models.ZoneResponse, 0)

	for _, zone := range zones {
		switch locationType {
		case config_enums.AREA:
			for _, area := range zone.Areas {
				if area.Name == locationName {
					zonesResponse = append(zonesResponse, zone)
					break
				}
			}
		case config_enums.CITY:
			for _, city := range zone.Cities {
				if city.Name == locationName {
					zonesResponse = append(zonesResponse, zone)
					break
				}
			}
		case config_enums.COUNTRY:
			for _, country := range zone.Countries {
				if country.Name == locationName {
					zonesResponse = append(zonesResponse, zone)
					break
				}
			}
		default:
			apiErr := helpers.NewApiError(400, fmt.Errorf("invalid location type: %s", locationType))
			return nil, &apiErr
		}
	}

	return zonesResponse, nil
}
**/

//We remove this function because we will only link zone with polygon
/**
func (s *zoneService) findZonesByLocation(ctx context.Context, latitude, longitude float64) ([]config_models.ZoneResponse, *helpers.ApiError) {
	// First check if the point is inside our countries
	countries, apiErr := s.countryService.GetAll(ctx, ``)
	if apiErr != nil {
		return nil, apiErr
	}

	countriesResponse := config_location_models.CountryResponse{}
	for _, country := range countries {
		if !config_utils.IsPointInPolygon(latitude, longitude, country.MultiPolygon) {
			return nil, nil
		} else {
			countriesResponse = country
		}
	}

	//then check if the point is inside our cities
	cities, apiErr := s.cityService.GetAll(ctx, `{"countryId": "`+countriesResponse.Id.Hex()+`"}`)
	if apiErr != nil {
		return nil, apiErr
	}

	citiesResponse := config_location_models.CityResponse{}
	for _, city := range cities {
		if !config_utils.IsPointInPolygon(latitude, longitude, city.MultiPolygon) {
			return nil, nil
		} else {
			citiesResponse = city
		}
	}

	//then check if the point is inside our areas
	areas, apiErr := s.areaService.GetAll(ctx, `{"cityId": "`+citiesResponse.Id.Hex()+`"}`)
	if apiErr != nil {
		return nil, apiErr
	}

	areasResponse := config_location_models.AreaResponse{}
	for _, area := range areas {
		if !config_utils.IsPointInPolygon(latitude, longitude, area.MultiPolygon) {
			return nil, nil
		} else {
			areasResponse = area
		}
	}

	// Get zones that match the country, city, and area IDs
	return s.GetAll(ctx, `{"countryIds":["`+countriesResponse.Id.Hex()+`"], "cityIds":["`+citiesResponse.Id.Hex()+`"], "areaIds":["`+areasResponse.Id.Hex()+`"]}`)
}
**/

func (s *zoneService) findZonesByCoordinates(ctx context.Context, latitude, longitude float64) ([]config_models.ZoneResponse, *helpers.ApiError) {
	allZones, apiErr := s.GetAll(ctx, ``)
	if apiErr != nil {
		return nil, apiErr
	}

	// Filter zones to only include those where the point is within their polygons
	matchingZones := make([]config_models.ZoneResponse, 0)
	for _, zone := range allZones {
		// Check if point is within any of the zone's polygons
		for _, polygon := range zone.Polygons {
			// Convert zone polygon points to the format expected by IsPointInPolygon
			points := make([][]float64, len(polygon.Points))
			for i, point := range polygon.Points {
				points[i] = []float64{point.Longitude, point.Latitude}
			}
			if config_utils.IsPointInPolygon(latitude, longitude, [][][]float64{points}) {
				matchingZones = append(matchingZones, zone)
				break // Found a matching polygon, no need to check others for this zone
			}
		}
	}

	return matchingZones, nil
}

// FindZonesContainingPoint finds all zones that contain the given point, combining results from both
// location-based matching and coordinate-based polygon matching
func (s *zoneService) FindZonesContainingPoint(ctx context.Context, latitude, longitude float64) ([]config_models.ZoneResponse, *helpers.ApiError) {
	// We will remove this function because we will only link zone with polygon
	// Get zones by location hierarchy (country -> city -> area)
	// locationZones, apiErr := s.findZonesByLocation(ctx, latitude, longitude)
	// if apiErr != nil {
	// 	return nil, apiErr
	// }

	// Get zones by direct coordinate matching
	coordinateZones, apiErr := s.findZonesByCoordinates(ctx, latitude, longitude)
	if apiErr != nil {
		return nil, apiErr
	}

	// Combine both results
	// Use a map to deduplicate zones
	uniqueZones := make(map[primitive.ObjectID]config_models.ZoneResponse)

	// for _, zone := range locationZones {
	// 	uniqueZones[zone.Id] = zone
	// }

	for _, zone := range coordinateZones {
		uniqueZones[zone.Id] = zone
	}

	// Convert map back to slice
	result := make([]config_models.ZoneResponse, 0, len(uniqueZones))
	for _, zone := range uniqueZones {
		result = append(result, zone)
	}

	return result, nil
}
