package product_shared

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"math"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryFeeArgs struct {
	Branch       *product_models.BranchResponse
	ZoneId       *primitive.ObjectID
	Address      *address_models.AddressResponse
	OrderType    order_enums.OrderType
	CurrencyCode *string
}

// calculateDistanceKm calculates the distance between two points using the Haversine formula.
// lat/lng are expected in degrees.
func calculateDistanceKm(lat1, lng1, lat2, lng2 float64) float64 {
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

// DeliveryDistanceKm returns the distance in km between the branch and the address (same formula as fee calculation).
// Branch and address coordinates must be [lng, lat]. Returns 0 if either is missing coordinates.
func DeliveryDistanceKm(branch *product_models.BranchResponse, address *address_models.AddressResponse) float64 {
	if branch == nil || len(branch.Coordinates) < 2 ||
		address == nil || address.Location == nil || len(address.Location.Coordinates) < 2 {
		return 0
	}
	branchLng := branch.Coordinates[0]
	branchLat := branch.Coordinates[1]
	addressLng := address.Location.Coordinates[0]
	addressLat := address.Location.Coordinates[1]
	d := calculateDistanceKm(branchLat, branchLng, addressLat, addressLng)
	if d < 0 {
		return 0
	}
	return d
}

// CalculateDeliveryFee returns the delivery fee and estimated time (in minutes) for the given arguments.
// It is intentionally forgiving for zone-based pricing (to keep backward compatibility) and stricter
// for distance-based pricing where configuration issues should surface as errors.
func CalculateDeliveryFee(args DeliveryFeeArgs) (float64, uint64, *helpers.ApiError) {
	if args.Branch == nil || args.CurrencyCode == nil {
		// Not enough information to calculate a delivery fee; keep previous behavior by returning zeros.
		return 0, 0, nil
	}

	mode := args.Branch.DeliveryPricingMode
	if mode == "" {
		mode = string(product_enums.DeliveryPricingModeZone)
	}

	switch mode {
	case string(product_enums.DeliveryPricingModeDistance):
		return calculateDistanceBasedFee(args)
	default: // zone-based (default / backward compatible)
		return calculateZoneBasedFee(args)
	}
}

func calculateZoneBasedFee(args DeliveryFeeArgs) (float64, uint64, *helpers.ApiError) {
	// This logic intentionally mirrors the existing checkout behavior:
	// - If anything is missing, it logs (in the caller) and returns zero fee and time.
	if args.Branch.Zones == nil || args.ZoneId == nil || args.OrderType != order_enums.DELIVERY {
		return 0, 0, nil
	}

	branchZone := findBranchZoneById(args.Branch.Zones, *args.ZoneId)
	if branchZone == nil {
		return 0, 0, nil
	}

	if branchZone.DeliveryFee.CurrencyCode != *args.CurrencyCode {
		return 0, 0, nil
	}

	return branchZone.DeliveryFee.Value, branchZone.EstimatedTime, nil
}

func findBranchZoneById(zones *[]product_models.BranchZone, id primitive.ObjectID) *product_models.BranchZone {
	if zones == nil {
		return nil
	}
	for i := range *zones {
		if (*zones)[i].ZoneId == id {
			return &(*zones)[i]
		}
	}
	return nil
}

func calculateDistanceBasedFee(args DeliveryFeeArgs) (float64, uint64, *helpers.ApiError) {
	if args.OrderType != order_enums.DELIVERY {
		// No delivery fee for non-delivery orders.
		return 0, 0, nil
	}

	cfg := args.Branch.DistancePricingConfig
	if cfg == nil {
		return 0, 0, &helpers.ApiError{
			StatusCode: 500,
			Message:    "distance pricing configuration is missing for branch",
		}
	}

	if args.Address == nil || args.Address.Location == nil || len(args.Address.Location.Coordinates) < 2 {
		return 0, 0, &helpers.ApiError{
			StatusCode: 400,
			Message:    "address location is required for distance-based delivery pricing",
		}
	}

	// Branch coordinates: [lng, lat]
	if len(args.Branch.Coordinates) < 2 {
		return 0, 0, &helpers.ApiError{
			StatusCode: 500,
			Message:    "branch coordinates are missing for distance-based delivery pricing",
		}
	}

	branchLng := args.Branch.Coordinates[0]
	branchLat := args.Branch.Coordinates[1]

	addressLng := args.Address.Location.Coordinates[0]
	addressLat := args.Address.Location.Coordinates[1]

	distanceKm := calculateDistanceKm(branchLat, branchLng, addressLat, addressLng)

	// Ensure non-negative values
	if distanceKm < 0 {
		distanceKm = 0
	}

	// If total distance is within the free-distance threshold, delivery is completely free
	if cfg.FreeDistanceKm > 0 && distanceKm <= cfg.FreeDistanceKm {
		return 0, 0, nil
	}

	// Apply free distance: first FreeDistanceKm are free, only the rest is chargeable
	freeKm := 0.0
	if cfg.FreeDistanceKm > 0 {
		freeKm = cfg.FreeDistanceKm
	}
	chargeableKm := distanceKm - freeKm
	if chargeableKm < 0 {
		chargeableKm = 0
	}

	baseFee := cfg.BaseFee.Value
	if baseFee < 0 {
		baseFee = 0
	}

	extraCostPerKm := cfg.ExtraCostPerKm.Value
	if extraCostPerKm < 0 {
		extraCostPerKm = 0
	}

	extraKm := chargeableKm - cfg.BaseDistanceKm
	if extraKm < 0 {
		extraKm = 0
	}

	everyKmFee := cfg.EveryKmFee
	if everyKmFee <= 0 {
		everyKmFee = 1
	}
	// Round extra fee: if remainder >= half of everyKmFee round up, else round down
	var steps float64
	if extraKm <= 0 {
		steps = 0
	} else {
		fullBlocks := math.Floor(extraKm / everyKmFee)
		remainder := extraKm - fullBlocks*everyKmFee
		if remainder >= everyKmFee/2 {
			steps = fullBlocks + 1
		} else {
			steps = fullBlocks
		}
	}
	if steps < 0 {
		steps = 0
	}
	extraFee := steps * extraCostPerKm
	fee := baseFee + extraFee
	if fee < 0 {
		fee = 0
	}

	// Distance mode: estimated time comes from zone (caller sets it); return 0 here.
	return fee, 0, nil
}

