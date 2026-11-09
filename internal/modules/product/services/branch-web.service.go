package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/services/geocoding"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	unserved_delivery_request_services "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/services"

	"math"
	"net/http"
	"slices"
	"sort"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchWebService interface {
	GetByLatLng(ctx context.Context, lat, lng float64,country, city, area *string) (*models.BranchResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BranchesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError)
	GetById(ctx context.Context, id string) (*models.BranchResponse, *helpers.ApiError)
}

type BranchWebRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, branchFilters filters.BranchFilters) (*models.BranchResponse, *helpers.ApiError)
	GetAll(ctx context.Context, branchFilters filters.BranchFilters) ([]models.BranchResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, branchFilters filters.BranchFilters) (*models.BranchesWithPagination, *helpers.ApiError)
}

type branchWeb struct {
	branchWebRepo BranchWebRepo
	zoneService   config_services.ZoneService
	unservedDeliveryRequestSvc unserved_delivery_request_services.UnservedDeliveryRequestService
}

func NewBranchWebService(i *do.Injector) (BranchWebService, error) {
	unservedSvc := do.MustInvoke[unserved_delivery_request_services.UnservedDeliveryRequestService](i)
	return &branchWeb{
		branchWebRepo: do.MustInvoke[BranchWebRepo](i),
		zoneService:   do.MustInvoke[config_services.ZoneService](i),
		unservedDeliveryRequestSvc: unservedSvc,
	}, nil
}

// calculateDistance calculates the distance between two points using the Haversine formula
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// findNearestBranch finds the branch closest to the given coordinates
func findNearestBranch(branches []models.BranchResponse, targetLat, targetLng float64) models.BranchResponse {
	if len(branches) == 0 {
		return models.BranchResponse{}
	}

	nearestBranch := branches[0]
	minDistance := math.MaxFloat64
	validBranches := 0

	for _, branch := range branches {
		if len(branch.Coordinates) >= 2 {
			branchLat := branch.Coordinates[1] // Latitude is typically at index 1
			branchLng := branch.Coordinates[0] // Longitude is typically at index 0

			// Validate coordinates are within valid ranges
			if branchLat >= -90 && branchLat <= 90 && branchLng >= -180 && branchLng <= 180 {
				distance := calculateDistance(targetLat, targetLng, branchLat, branchLng)
				validBranches++

				if distance < minDistance {
					minDistance = distance
					nearestBranch = branch
				}
			}
		}
	}

	// If no valid coordinates found, return the first branch as fallback
	if validBranches == 0 {
		return branches[0]
	}

	return nearestBranch
}
func (p *branchWeb) recordUnservedAndReturn404(ctx context.Context, lat, lng float64, country, city, area *string) (*models.BranchResponse, *helpers.ApiError) {
	var userId *primitive.ObjectID
	if cfg, _ := util.GetReqAppCfg(ctx); cfg != nil && cfg.User != nil {
		userId = &cfg.User.Id
	}

	// Fill missing location from reverse geocoding when country, city, or area are not provided
	if country == nil || city == nil || area == nil {
		if gcCountry, gcCity, gcArea, err := geocoding.ReverseGeocode(lat, lng); err == nil {
			if country == nil && gcCountry != "" {
				country = &gcCountry
			}
			if city == nil && gcCity != "" {
				city = &gcCity
			}
			if area == nil && gcArea != "" {
				area = &gcArea
			}
		}
	}

	p.unservedDeliveryRequestSvc.RecordRequest(ctx, lat, lng, userId, country, city, area, nil)
	return nil, &helpers.ApiError{
		StatusCode: http.StatusNotFound,
		Message:    "No branch found for the given coordinates",
	}
}

func (p *branchWeb) GetByLatLng(ctx context.Context, lat, lng float64,country, city, area *string) (*models.BranchResponse, *helpers.ApiError) {
	zones, err := p.zoneService.FindZonesContainingPoint(ctx, lat, lng)
	if err != nil {
		return nil, err
	}
	if len(zones) == 0 {
		return p.recordUnservedAndReturn404(ctx, lat, lng, country, city, area)
	
	}

	zoneIds := make([]primitive.ObjectID, len(zones))
	for i, zone := range zones {
		zoneIds[i] = zone.Id
	}

	branchFilters := filters.BranchFilters{
		ZoneIds: zoneIds,
	}

	branches, err := p.branchWebRepo.GetAll(ctx, branchFilters)
	if err != nil {
		return nil, err
	}

	if len(branches) == 0 {
		return p.recordUnservedAndReturn404(ctx, lat, lng, country, city, area)
	}

	// Find the nearest branch when multiple branches are found
	branch := findNearestBranch(branches, lat, lng)
	mode := branch.DeliveryPricingMode
	if mode == "" {
		mode = string(product_enums.DeliveryPricingModeZone)
	}

	if mode == string(product_enums.DeliveryPricingModeZone) && branch.Zones != nil {
		activeZones := []models.BranchZone{}
		for _, zone := range *branch.Zones {
			if slices.Contains(zoneIds, zone.ZoneId) {
				activeZones = append(activeZones, zone)
			}
		}

		sort.Slice(activeZones, func(i, j int) bool {
			// If MinOrderAmount is the same, compare by DeliveryFee.Value (lower is better)
			return activeZones[i].DeliveryFee.Value < activeZones[j].DeliveryFee.Value
		})

		branch.Zones = &activeZones
	}

	return &branch, nil
}

func (p *branchWeb) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BranchesWithPagination, *helpers.ApiError) {
	branchFilters, err := filters.NewBranchFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}
	return p.branchWebRepo.GetWithPagination(ctx, skip, limit, *branchFilters)
}

func (p *branchWeb) GetAll(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError) {
	branchFilters, err := filters.NewBranchFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}
	return p.branchWebRepo.GetAll(ctx, *branchFilters)
}

func (p *branchWeb) GetById(ctx context.Context, id string) (*models.BranchResponse, *helpers.ApiError) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch id format",
		}
	}

	return p.branchWebRepo.GetById(ctx, objectId, filters.BranchFilters{})
}
