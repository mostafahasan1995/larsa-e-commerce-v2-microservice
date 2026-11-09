package product_dtos

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBranchZoneDto struct {
	MinOrderAmount float64            `json:"minOrderAmount" validate:"required,min=0"`
	DeliveryFee    PriceItemDto       `json:"deliveryFee" validate:"required"`
	EstimatedTime  uint64             `json:"estimatedTime" validate:"required"`
	ZoneId         primitive.ObjectID `json:"zoneId" validate:"required"`
}

type DateRangeDto struct {
	From uint64 `json:"from" validate:"required,min=0"`
	To   uint64 `json:"to" validate:"required,min=0,gtfield=From,max=86400000"`
}

type DayScheduleDto struct {
	Label      *string        `json:"label,omitempty"`
	TimeRanges []DateRangeDto `json:"timeRanges"`
}

type OpeningDaysDto struct {
	Days map[time.Weekday]DayScheduleDto `json:"days"`
}

type TimeZoneDto struct {
	Location  string `json:"location" validate:"required"`
	UtcOffset int64  `json:"utcOffset"`
}

type DistancePricingConfigDto struct {
	BaseFee PriceItemDto `json:"baseFee" validate:"required"`
	// BaseDistanceKm is the distance (in kilometers) included in the BaseFee.
	BaseDistanceKm float64 `json:"baseDistanceKm" validate:"required,min=0"`
	// FreeDistanceKm is the distance (in kilometers) that is free before any fee calculation.
	FreeDistanceKm float64 `json:"freeDistanceKm" validate:"omitempty,gte=0"`
	// EveryKmFee is the number of kilometers per block; each block beyond BaseDistanceKm adds ExtraCostPerKm.Value.
	EveryKmFee float64 `json:"everyKmFee" validate:"required,gt=0"`
	// ExtraCostPerKm is the fee (per block of EveryKmFee km) for distance beyond BaseDistanceKm.
	ExtraCostPerKm PriceItemDto `json:"extraCostPerKm" validate:"required"`
	// EstimatedMinutesPerKm optional; estimated time for distance mode comes from zone.
	EstimatedMinutesPerKm float64 `json:"estimatedMinutesPerKm,omitempty" validate:"omitempty,gt=0"`
}

type CreateBranchDto struct {
	Name                          types.Localizable[string]             `json:"name" validate:"required"`
	Image                         types.Localizable[types.FileFieldDto] `json:"image" validate:"required"`
	MobileImage                   types.Localizable[types.FileFieldDto] `json:"mobileImage" validate:"omitempty"`
	CountryIso3                   string                                `json:"countryIso3" validate:"required"`
	City                          string                                `json:"city" validate:"required"`
	Area                          string                                `json:"area" validate:"required"`
	Coordinates                   []float64                             `json:"coordinates" validate:"required,min=2,max=2"`
	IsActive                      bool                                  `json:"isActive" validate:"required"`
	TimeZone                      TimeZoneDto                           `json:"timeZone" validate:"required"`
	PhoneNumber                   string                                `json:"phoneNumber" validate:"required"`
	MobileNumber                  string                                `json:"mobileNumber" validate:"omitempty"`
	BaseCurrencyCode              string                                `json:"baseCurrencyCode" validate:"required"`
	SupportedCurrencyCodes        *[]string                             `json:"supportedCurrencyCodes" validate:"omitempty,min=1,dive"`
	RealTimeCurrencyConversion    bool                                  `json:"realTimeCurrencyConversion"`
	CanSetDeliveryOnPickup        bool                                  `json:"canSetDeliveryOnPickup"`
	MinTimeBeforeOrder            *uint64                               `json:"minTimeBeforeOrder" validate:"omitempty,min=0"`
	MaxTimeBeforeOrder            *uint64                               `json:"maxTimeBeforeOrder" validate:"omitempty,min=0"`
	Products                      *[]LinkBranchToProductDto             `json:"products" validate:"omitempty,min=0,dive"`
	Zones                         *[]CreateBranchZoneDto                `json:"zones" validate:"min=0"`
	OpeningDays                   *OpeningDaysDto                       `json:"openingDays" validate:"required"`
	VendorId                      *primitive.ObjectID                   `bson:"vendorId,omitempty" json:"vendorId,omitempty" validate:"omitempty"`
	ManagersIds                   *[]primitive.ObjectID                 `json:"managersIds,omitempty"`
	base_models.MetaDataBaseModel `bson:",inline"`
	DeliveryPricingMode           string                                `json:"deliveryPricingMode" validate:"omitempty,oneof=zone distance"`
	DistancePricingConfig         *DistancePricingConfigDto             `json:"distancePricingConfig,omitempty" validate:"omitempty"`

}

func (dto *CreateBranchDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	zoneRepo := base_repos.NewBaseRepo[config_models.Zone, config_models.ZoneResponse](args.DB, enums.ZONE)
	productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)

	// Validate MinTimeBeforeOrder and MaxTimeBeforeOrder
	if dto.MinTimeBeforeOrder != nil && dto.MaxTimeBeforeOrder != nil {
		if *dto.MinTimeBeforeOrder > *dto.MaxTimeBeforeOrder {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "MinTimeBeforeOrder must be less than or equal to MaxTimeBeforeOrder",
			}
		}
	}

	if dto.TimeZone.UtcOffset < product_shared.MINIMUM_UTC_OFFSET || dto.TimeZone.UtcOffset > product_shared.MAXIMUM_UTC_OFFSET {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    fmt.Sprintf("UTC offset must be between %v and %v", product_shared.MINIMUM_UTC_OFFSET, product_shared.MAXIMUM_UTC_OFFSET),
		}

	}
	if dto.DeliveryPricingMode == "" {
		dto.DeliveryPricingMode = string(product_enums.DeliveryPricingModeZone)
	}

	// Validate distance-based configuration when applicable
	if dto.DeliveryPricingMode == string(product_enums.DeliveryPricingModeDistance) {
		if dto.DistancePricingConfig == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "distancePricingConfig is required when deliveryPricingMode is distance",
			}
		}

		// Basic non-negative checks are already enforced by struct tags; here we enforce currency alignment.
		baseCurrency := dto.BaseCurrencyCode
		supported := map[string]struct{}{baseCurrency: {}}
		if dto.SupportedCurrencyCodes != nil {
			for _, code := range *dto.SupportedCurrencyCodes {
				supported[code] = struct{}{}
			}
		}

		if _, ok := supported[dto.DistancePricingConfig.BaseFee.CurrencyCode]; !ok {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "distance pricing baseFee currencyCode must match branch base or supported currencies",
			}
		}
		if _, ok := supported[dto.DistancePricingConfig.ExtraCostPerKm.CurrencyCode]; !ok {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "distance pricing extraCostPerKm currencyCode must match branch base or supported currencies",
			}
		}
	}

	if dto.OpeningDays != nil {
		for day, dateRanges := range dto.OpeningDays.Days {
			if day < 0 || day > 6 {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "day must be between 0 and 6",
				}
			}
			for _, dateRange := range dateRanges.TimeRanges {

				if dateRange.From > product_shared.MAX_TIME_SPAN_IN_A_DAY || dateRange.To > product_shared.MAX_TIME_SPAN_IN_A_DAY {
					return nil, &helpers.ApiError{
						StatusCode: http.StatusBadRequest,
						Message:    fmt.Sprintf("Time must be less than %v", product_shared.MAX_TIME_SPAN_IN_A_DAY),
					}
				}

				if dateRange.From > dateRange.To {
					return nil, &helpers.ApiError{
						StatusCode: http.StatusBadRequest,
						Message:    "From must be less than to",
					}
				}
			}
		}
	}

	if dto.Products != nil {
		productIds := make([]primitive.ObjectID, len(*dto.Products))
		for i, product := range *dto.Products {
			productIds[i] = product.ProductId
		}

		if util.HasDuplicateValues(productIds) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "products must be unique",
			}
		}

		_, err := productRepo.GetByIds(args.Context, productIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
		if err != nil {
			return nil, err
		}

	}

	if dto.Zones != nil {
		zoneIds := make([]primitive.ObjectID, len(*dto.Zones))
		for i, zone := range *dto.Zones {
			zoneIds[i] = zone.ZoneId
		}

		if util.HasDuplicateValues(zoneIds) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "zones must be unique",
			}
		}

		_, err := zoneRepo.GetByIds(args.Context, zoneIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
		if err != nil {
			return nil, err
		}

	}
	currencyRepo := base_repos.NewBaseRepo[config_models.Currency, config_models.CurrencyResponse](args.DB, enums.CURRENCY)

	uniqueCurrencyCodes := make([]string, 0)
	allCurrencyCodes := make([]string, 0)
	allCurrencyCodes = append(allCurrencyCodes, dto.BaseCurrencyCode)
	if dto.SupportedCurrencyCodes != nil {
		allCurrencyCodes = append(allCurrencyCodes, *dto.SupportedCurrencyCodes...)
		if util.HasDuplicateValues(allCurrencyCodes) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "currencies Codes must be unique",
			}
		}

		uniqueCurrencyCodes = util.ExtractUniqueValues(allCurrencyCodes)

	}

	currencies, err := currencyRepo.GetAll(args.Context, []bson.M{{"$match": bson.M{"code": bson.M{"$in": uniqueCurrencyCodes}}}})
	if err != nil {
		return nil, err
	}

	if len(currencies) != len(uniqueCurrencyCodes) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Some currencies codes are not found",
		}
	}
	return nil, nil
}
