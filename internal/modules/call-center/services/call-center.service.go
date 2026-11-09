package call_center_services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"

	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/services/geocoding"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	call_center_models "larsa-e-commerce-microservice-v2/internal/modules/call-center/models"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CallCenterService interface {
	Geocode(ctx context.Context, address string) (*call_center_models.GeocodeResponse, *helpers.ApiError)
	GetDeliveryFeasibility(ctx context.Context, lat, lng float64, currencyCode *string) (*call_center_models.DeliveryFeasibilityResponse, *helpers.ApiError)
	SearchBranchMenu(ctx context.Context, branchId string, query *string, skip, limit int64, sorts *string, currencyCode *string) (*call_center_models.BranchMenuSearchResponse, *helpers.ApiError)
	Lookup(ctx context.Context, address string, currencyCode *string, menuQuery *string, skip, limit int64, sorts *string) (*call_center_models.LookupResponse, *helpers.ApiError)
}

type callCenterService struct {
	branchWebService  product_services.BranchWebService
	zoneService       config_services.ZoneService
	branchService     product_services.BranchService
	webProductService product_services.WebProductService
	productService    product_services.ProductService
}

func NewCallCenterService(i *do.Injector) (CallCenterService, error) {
	return &callCenterService{
		branchWebService:  do.MustInvoke[product_services.BranchWebService](i),
		zoneService:       do.MustInvoke[config_services.ZoneService](i),
		branchService:     do.MustInvoke[product_services.BranchService](i),
		webProductService: do.MustInvoke[product_services.WebProductService](i),
		productService:    do.MustInvoke[product_services.ProductService](i),
	}, nil
}

// Geocode converts an address string to coordinates and location details
func (s *callCenterService) Geocode(ctx context.Context, address string) (*call_center_models.GeocodeResponse, *helpers.ApiError) {
	lat, lng, country, city, area, err := geocoding.Geocode(address)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("Geocoding failed: %v", err),
		}
	}

	return &call_center_models.GeocodeResponse{
		Latitude:  lat,
		Longitude: lng,
		Country:   country,
		City:      city,
		Area:      area,
	}, nil
}

// GetDeliveryFeasibility checks if delivery is possible and calculates fee and ETA
func (s *callCenterService) GetDeliveryFeasibility(ctx context.Context, lat, lng float64, currencyCode *string) (*call_center_models.DeliveryFeasibilityResponse, *helpers.ApiError) {
	// Get nearest branch by coordinates
	branch, apiErr := s.branchWebService.GetByLatLng(ctx, lat, lng, nil, nil, nil)
	if apiErr != nil {
		// No branch found - return not feasible
		return &call_center_models.DeliveryFeasibilityResponse{
			Feasible:           false,
			Branch:             nil,
			DeliveryFee:        0,
			EstimatedMinutes:   0,
			ZoneId:             nil,
			Coordinates:        nil,
			DeliveryPricingMod: "",
		}, nil
	}

	// Build address from coordinates for delivery fee calculation
	addressForFee := &address_models.AddressResponse{
		AddressModel: address_models.AddressModel{
			Location: &address_models.GeoLocation{
				Type:        "Point",
				Coordinates: []float64{lng, lat}, // [lng, lat] format
			},
		},
	}

	// Determine zone ID if zone-based pricing
	var zoneId *primitive.ObjectID
	if branch.Zones != nil && len(*branch.Zones) > 0 {
		// Use the first zone (already sorted by delivery fee in BranchWebService)
		zoneId = &(*branch.Zones)[0].ZoneId
	}

	// Use default currency if not provided
	defaultCurrency := "USD"
	if currencyCode == nil {
		currencyCode = &defaultCurrency
	}

	// Calculate delivery fee and estimated time
	fee, estimatedMinutes, feeErr := product_shared.CalculateDeliveryFee(product_shared.DeliveryFeeArgs{
		Branch:       branch,
		ZoneId:       zoneId,
		Address:      addressForFee,
		OrderType:    order_enums.DELIVERY,
		CurrencyCode: currencyCode,
	})

	if feeErr != nil {
		return nil, feeErr
	}

	// Distance mode returns 0 for estimated time; use first zone when branch has zones
	if estimatedMinutes == 0 && branch.Zones != nil && len(*branch.Zones) > 0 {
		estimatedMinutes = (*branch.Zones)[0].EstimatedTime
	}

	// Prepare coordinates info
	coordinates := &call_center_models.CoordinatesInfo{
		Customer: []float64{lng, lat},
		Branch:   branch.Coordinates,
	}
	var everyFeeKm float64
	var extraCostPerKm product_models.PriceItem
	var distanceKm *float64
	deliveryPricingMod := "zone"
	if branch.DeliveryPricingMode == string(product_enums.DeliveryPricingModeDistance) {
		deliveryPricingMod = "distance"
		everyFeeKm = branch.DistancePricingConfig.ExtraCostPerKm.Value // cost per block (e.g. 2 AED per everyKmFee km)
		extraCostPerKm = branch.DistancePricingConfig.ExtraCostPerKm
		d := product_shared.DeliveryDistanceKm(branch, addressForFee)
		distanceKm = &d
	}

	return &call_center_models.DeliveryFeasibilityResponse{
		Feasible:           true,
		Branch:             branch,
		DeliveryFee:        fee,
		EstimatedMinutes:   estimatedMinutes,
		ZoneId:             zoneId,
		Coordinates:        coordinates,
		DistanceKm:         distanceKm,
		DeliveryPricingMod:  deliveryPricingMod,
		EveryFeeKm:         everyFeeKm,
		ExtraCostPerKm:     extraCostPerKm,
	}, nil
}

// filterBranchProductsByName filters branch products by product name using regex (case-insensitive)
func (s *callCenterService) filterBranchProductsByName(ctx context.Context, branchProducts []product_models.BranchProductResponse, query string) []product_models.BranchProductResponse {
	if len(branchProducts) == 0 {
		return branchProducts
	}

	productIds := make([]primitive.ObjectID, 0, len(branchProducts))
	for _, bp := range branchProducts {
		if bp.ProductId != nil {
			productIds = append(productIds, *bp.ProductId)
		}
	}
	if len(productIds) == 0 {
		return branchProducts
	}

	// Fetch products (only _id and name) for regex match
	products, apiErr := s.productService.GetAll(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": productIds}}},
		{"$project": bson.M{"_id": 1, "name": 1}},
	})
	if apiErr != nil || len(products) == 0 {
		return branchProducts
	}

	// Build productId -> all name translations (en, ar, etc.) for matching in any language
	productNamesByLang := make(map[primitive.ObjectID][]string)
	for _, p := range products {
		if len(p.Name) > 0 {
			names := make([]string, 0)
			for _, v := range p.Name {
				if v != "" {
					names = append(names, v)
				}
			}
			productNamesByLang[p.Id] = names
		}
	}

	// Compile regex - treat query as literal substring, case-insensitive (works for both en and ar)
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(query))
	if err != nil {
		return branchProducts
	}

	filtered := make([]product_models.BranchProductResponse, 0, len(branchProducts))
	for _, bp := range branchProducts {
		if bp.ProductId == nil {
			continue
		}
		names := productNamesByLang[*bp.ProductId]
		matches := false
		for _, name := range names {
			if re.MatchString(name) {
				matches = true
				break
			}
		}
		if matches {
			filtered = append(filtered, bp)
		}
	}
	return filtered
}

// SearchBranchMenu returns products with full details (name, description, etc.) for a branch, with filter by product name regex.
// Prices match the website branch view when the same branchId and currencyCode are used.
func (s *callCenterService) SearchBranchMenu(ctx context.Context, branchId string, query *string, skip, limit int64, sorts *string, currencyCode *string) (*call_center_models.BranchMenuSearchResponse, *helpers.ApiError) {
	// Validate branch ID
	branchObjId, err := primitive.ObjectIDFromHex(branchId)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid branch ID format",
		}
	}

	// Load branch with products - same flow as GET /admin/branches/{id}?query={"popFields":["products"]}
	branch, apiErr := s.branchService.GetById(ctx, branchObjId, `{"popFields":["products"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	branchProducts := branch.Products
	if branchProducts == nil {
		branchProducts = []product_models.BranchProductResponse{}
	}

	// Apply search filter by product name (regex) when query is provided
	if query != nil && strings.TrimSpace(*query) != "" {
		branchProducts = s.filterBranchProductsByName(ctx, branchProducts, *query)
	}

	// Calculate pagination
	totalCount := int64(len(branchProducts))
	var totalPages float64
	if limit > 0 {
		totalPages = math.Ceil(float64(totalCount) / float64(limit))
	}

	// Paginate branch products
	start := skip
	end := skip + limit
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}
	paginatedProducts := branchProducts[start:end]

	// Convert to full ProductResponse via webProductService
	if len(paginatedProducts) == 0 {
		return &call_center_models.BranchMenuSearchResponse{
			DataWithPagination: &types.DataWithPagination[product_models.ProductResponse]{
				Data: []product_models.ProductResponse{},
				Pagination: types.Pagination{
					TotalPages: totalPages,
					PerPage:    limit,
					TotalCount: totalCount,
				},
			},
		}, nil
	}

	productIds := make([]primitive.ObjectID, 0, len(paginatedProducts))
	for _, bp := range paginatedProducts {
		if bp.ProductId != nil {
			productIds = append(productIds, *bp.ProductId)
		}
	}

	productQuery := struct {
		Ids []primitive.ObjectID `json:"ids"`
	}{Ids: productIds}
	productQueryJSON, _ := json.Marshal(productQuery)

	// Use default currency when not provided so prices match website branch view
	if currencyCode == nil {
		defaultCode := configuration.EnvironmentVariables.DefaultCurrencyCode
		currencyCode = &defaultCode
	}
	log.Printf("[CALL_CENTER_SEARCH_BRANCH_MENU] branchId=%s currencyCode=%s", branchId, *currencyCode)

	productsResult, apiErr := s.webProductService.GetWithPagination(ctx, product_utils.GetWithPaginationArgs{
		Skip:         0,
		Limit:        int64(len(productIds)),
		Query:        string(productQueryJSON),
		BranchId:     &branchObjId,
		CurrencyCode: currencyCode,
	})
	if apiErr != nil {
		return nil, apiErr
	}

	return &call_center_models.BranchMenuSearchResponse{
		DataWithPagination: &types.DataWithPagination[product_models.ProductResponse]{
			Data: productsResult.Data,
			Pagination: types.Pagination{
				TotalPages: totalPages,
				PerPage:    limit,
				TotalCount: totalCount,
			},
		},
	}, nil
}

// Lookup performs combined geocoding, delivery feasibility check, and optional menu search
func (s *callCenterService) Lookup(ctx context.Context, address string, currencyCode *string, menuQuery *string, skip, limit int64, sorts *string) (*call_center_models.LookupResponse, *helpers.ApiError) {
	// Step 1: Geocode the address
	geocodeResult, geocodeErr := s.Geocode(ctx, address)
	if geocodeErr != nil {
		return nil, geocodeErr
	}

	// Step 2: Get delivery feasibility
	deliveryFeasibility, deliveryErr := s.GetDeliveryFeasibility(ctx, geocodeResult.Latitude, geocodeResult.Longitude, currencyCode)
	if deliveryErr != nil {
		return nil, deliveryErr
	}

	response := &call_center_models.LookupResponse{
		Geocode:             geocodeResult,
		DeliveryFeasibility: deliveryFeasibility,
	}

	// Step 3: Optional menu search if menuQuery provided and branch exists
	if menuQuery != nil && *menuQuery != "" && deliveryFeasibility.Feasible && deliveryFeasibility.Branch != nil {
		branchIdHex := deliveryFeasibility.Branch.Id.Hex()
		menuResult, menuErr := s.SearchBranchMenu(ctx, branchIdHex, menuQuery, skip, limit, sorts, currencyCode)
		if menuErr != nil {
			// Don't fail the entire request if menu search fails, just omit it
			response.BranchMenu = nil
		} else {
			response.BranchMenu = menuResult
		}
	}

	return response, nil
}
