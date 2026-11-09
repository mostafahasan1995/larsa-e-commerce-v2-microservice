package call_center_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GeocodeResponse represents the result of forward geocoding
type GeocodeResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Area      string  `json:"area"`
}

// CoordinatesInfo represents coordinate information for customer and branch
type CoordinatesInfo struct {
	Customer []float64 `json:"customer"` // [lng, lat]
	Branch   []float64 `json:"branch"`   // [lng, lat]
}

// DeliveryFeasibilityResponse represents delivery feasibility information
type DeliveryFeasibilityResponse struct {
	Feasible           bool                           `json:"feasible"`
	Branch             *product_models.BranchResponse `json:"branch,omitempty"`
	DeliveryFee        float64                        `json:"deliveryFee"`
	EstimatedMinutes   uint64                         `json:"estimatedMinutes"`
	ZoneId             *primitive.ObjectID            `json:"zoneId,omitempty"`
	Coordinates        *CoordinatesInfo               `json:"coordinates,omitempty"`
	DistanceKm         *float64                       `json:"distanceKm,omitempty"`         // distance used for fee (distance mode only)
	DeliveryPricingMod string                         `json:"deliveryPricingMod"`             // "zone" or "distance" — how the fee was computed
	EveryFeeKm         float64                        `json:"everyFeeKm,omitempty"`          // extra cost per block of km (distance mode only) = ExtraCostPerKm.Value
	ExtraCostPerKm     product_models.PriceItem       `json:"extraCostPerKm,omitempty"`      // full price item for extra per km (distance mode only)
}

// BranchMenuSearchResponse represents paginated products with full details (name, description, etc.)
type BranchMenuSearchResponse struct {
	*types.DataWithPagination[product_models.ProductResponse]
}

// LookupResponse combines geocoding, delivery feasibility, and optional menu search
type LookupResponse struct {
	Geocode             *GeocodeResponse             `json:"geocode"`
	DeliveryFeasibility *DeliveryFeasibilityResponse `json:"deliveryFeasibility"`
	BranchMenu          *BranchMenuSearchResponse    `json:"branchMenu,omitempty"`
}
