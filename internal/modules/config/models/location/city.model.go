package config_location_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceItem struct {
	Value        float64             `bson:"value" json:"value"`
	CurrencyCode string              `bson:"currencyCode" json:"currencyCode"`
	UnitId       *primitive.ObjectID `bson:"unitId,omitempty" json:"unitId,omitempty"`
}
type CityBranchZone struct {
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

type CityBranch struct {
	Name                       types.Localizable[string]             `bson:"name" json:"name"`
	Image                      types.Localizable[types.FileFieldDto] `bson:"image" json:"image"`
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
	TimeBeforeOrder            *uint64                               `bson:"timeBeforeOrder,omitempty" json:"timeBeforeOrder,omitempty"`
	Zones                      *[]CityBranchZone                     `bson:"zones,omitempty" json:"zones,omitempty"`
	OpeningDays                *OpeningDays                          `bson:"openingDays,omitempty" json:"openingDays,omitempty"`
	VendorId                   *primitive.ObjectID                   `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	models.MetaDataBaseModel   `bson:",inline"`
}

type BaseCity struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     string `bson:"name" json:"name"`
	IsActive                 *bool  `bson:"isActive,omitempty" json:"isActive,omitempty"`

	CountryId    primitive.ObjectID `bson:"countryId,omitempty" json:"countryId,omitempty"`
	CountryName  *string            `bson:"countryName,omitempty" json:"countryName,omitempty"`
	MultiPolygon [][][]float64      `bson:"multiPolygon" json:"multiPolygon"`
}

type City struct {
	BaseCity `bson:",inline"`
}

type CityResponse struct {
	City     `bson:",inline"`
	Country  *Country     `bson:"country,omitempty" json:"country,omitempty"`
	Branches []CityBranch `bson:"branches" json:"branches"`
}

type CitiesWithPagination = types.DataWithPagination[CityResponse]
