package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchZone struct {
	MinOrderAmount float64            `bson:"minOrderAmount" json:"minOrderAmount"`
	DeliveryFee    PriceItem          `bson:"deliveryFee" json:"deliveryFee"`
	EstimatedTime  uint64             `bson:"estimatedTime" json:"estimatedTime"`
	ZoneId         primitive.ObjectID `bson:"zoneId" json:"zoneId"`
}

type DateRange struct {
	From uint64 `bson:"from" json:"from"`
	To   uint64 `bson:"to" json:"to"`
}

type DaySchedule struct {
	Label      *string     `bson:"label" json:"label"`
	TimeRanges []DateRange `bson:"timeRanges" json:"timeRanges"`
}

type OpeningDays struct {
	Days map[time.Weekday]DaySchedule `bson:"days" json:"days"`
}

type TimeZone struct {
	Location  string `bson:"location" json:"location"`
	UtcOffset int64  `bson:"utcOffset" json:"utcOffset"`
}

type DistancePricingConfig struct {
	// BaseFee is the base delivery fee that covers the first BaseDistanceKm kilometers.
	BaseFee PriceItem `bson:"baseFee" json:"baseFee"`
	// BaseDistanceKm is the distance (in kilometers) included in the BaseFee.
	BaseDistanceKm float64 `bson:"baseDistanceKm" json:"baseDistanceKm"`
	// FreeDistanceKm is the distance (in kilometers) that is free before any fee calculation; when set, it is subtracted from total distance first.
	FreeDistanceKm float64 `bson:"freeDistanceKm" json:"freeDistanceKm"`
	// EveryKmFee is the number of kilometers per "block"; each block beyond BaseDistanceKm adds ExtraCostPerKm.Value.
	EveryKmFee float64 `bson:"everyKmFee" json:"everyKmFee"`
	// ExtraCostPerKm is the fee (per block of EveryKmFee km) for distance beyond BaseDistanceKm.
	ExtraCostPerKm PriceItem `bson:"extraCostPerKm" json:"extraCostPerKm"`
	// EstimatedMinutesPerKm unused for distance mode; estimated time comes from zone.
	EstimatedMinutesPerKm float64 `bson:"estimatedMinutesPerKm,omitempty" json:"estimatedMinutesPerKm,omitempty"`
}

type BaseBranch struct {
	Name                       types.Localizable[string]             `bson:"name" json:"name"`
	Image                      types.Localizable[types.FileFieldDto] `bson:"image" json:"image"`
	MobileImage                types.Localizable[types.FileFieldDto] `bson:"mobileImage" json:"mobileImage"`
	CountryIso3                string                                `bson:"countryIso3" json:"countryIso3"`
	City                       string                                `bson:"city" json:"city"`
	Area                       string                                `bson:"area" json:"area"`
	Coordinates                []float64                             `bson:"coordinates" json:"coordinates"`
	IsActive                   bool                                  `bson:"isActive" json:"isActive"`
	TimeZone                   TimeZone                              `bson:"timeZone" json:"timeZone"`
	PhoneNumber                string                                `bson:"phoneNumber" json:"phoneNumber"`
	MobileNumber               string                                `bson:"mobileNumber" json:"mobileNumber"`
	BaseCurrencyCode           string                                `bson:"baseCurrencyCode" json:"baseCurrencyCode"`
	SupportedCurrencyCodes     *[]string                             `bson:"supportedCurrencyCodes,omitempty" json:"supportedCurrencyCodes,omitempty"`
	RealTimeCurrencyConversion bool                                  `bson:"realTimeCurrencyConversion" json:"realTimeCurrencyConversion"`
	CanSetDeliveryOnPickup     bool                                  `bson:"canSetDeliveryOnPickup" json:"canSetDeliveryOnPickup"`
	MinTimeBeforeOrder         *uint64                               `bson:"minTimeBeforeOrder,omitempty" json:"minTimeBeforeOrder,omitempty"`
	MaxTimeBeforeOrder         *uint64                               `bson:"maxTimeBeforeOrder,omitempty" json:"maxTimeBeforeOrder,omitempty"`
	// DeliveryPricingMode controls how delivery fees are calculated for this branch.
	// When empty or \"zone\", zone-based pricing is used (backward compatible).
	// When \"distance\", DistancePricingConfig is used.
	DeliveryPricingMode           string                 `bson:"deliveryPricingMode,omitempty" json:"deliveryPricingMode,omitempty"`
	DistancePricingConfig         *DistancePricingConfig `bson:"distancePricingConfig,omitempty" json:"distancePricingConfig,omitempty"`
	Zones                         *[]BranchZone          `bson:"zones,omitempty" json:"zones,omitempty"`
	OpeningDays                   *OpeningDays           `bson:"openingDays,omitempty" json:"openingDays,omitempty"`
	VendorId                      *primitive.ObjectID    `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	ManagersIds                   *[]primitive.ObjectID  `bson:"managersIds,omitempty" json:"managersIds,omitempty"`
	base_models.MetaDataBaseModel `bson:",inline"`
}
type Branch struct {
	BaseBranch `bson:",inline"`
}

type BranchResponse struct {
	BaseBranch `bson:",inline"`
	Products   []BranchProductResponse `bson:"products" json:"products"`
}

type BranchesWithPagination = types.DataWithPagination[BranchResponse]
