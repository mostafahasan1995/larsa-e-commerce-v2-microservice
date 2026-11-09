package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"time"
)

func ConvertToBranchModel(createBranchDto dtos.CreateBranchDto, additionalOptions types.AdditionalOptions) models.Branch {
	var openingDays *product_models.OpeningDays
	if createBranchDto.OpeningDays != nil {
		// res := HandleOpeningDaysSpinningAcrossMultipleDays(createBranchDto.TimeZone.UtcOffset, *createBranchDto.OpeningDays)
		// openingDays = &res

		if createBranchDto.OpeningDays != nil {
			convertedOpeningDays := make(map[time.Weekday]product_models.DaySchedule)
			for day, daySchedule := range createBranchDto.OpeningDays.Days {
				convertedDateRanges := make([]product_models.DateRange, len(daySchedule.TimeRanges))
				for i, dateRange := range daySchedule.TimeRanges {
					convertedDateRanges[i] = product_models.DateRange{
						From: dateRange.From,
						To:   dateRange.To,
					}
				}
				convertedOpeningDays[day] = product_models.DaySchedule{
					Label:      daySchedule.Label,
					TimeRanges: convertedDateRanges,
				}
			}
			openingDays = &product_models.OpeningDays{
				Days: convertedOpeningDays,
			}
		}
	}

	branch := models.Branch{
		BaseBranch: models.BaseBranch{
			Name:                       createBranchDto.Name,
			Image:                      createBranchDto.Image,
			MobileImage:                createBranchDto.MobileImage,
			MetaDataBaseModel:          util.CreateMetaDataBaseModel(additionalOptions),
			CountryIso3:                createBranchDto.CountryIso3,
			City:                       createBranchDto.City,
			Area:                       createBranchDto.Area,
			Coordinates:                createBranchDto.Coordinates,
			IsActive:                   createBranchDto.IsActive,
			PhoneNumber:                createBranchDto.PhoneNumber,
			MobileNumber:               createBranchDto.MobileNumber,
			BaseCurrencyCode:           createBranchDto.BaseCurrencyCode,
			SupportedCurrencyCodes:     createBranchDto.SupportedCurrencyCodes,
			RealTimeCurrencyConversion: createBranchDto.RealTimeCurrencyConversion,
			CanSetDeliveryOnPickup:     createBranchDto.CanSetDeliveryOnPickup,
			MinTimeBeforeOrder:         createBranchDto.MinTimeBeforeOrder,
			MaxTimeBeforeOrder:         createBranchDto.MaxTimeBeforeOrder,
			OpeningDays:                openingDays,
			DeliveryPricingMode:        createBranchDto.DeliveryPricingMode,
			ManagersIds:                createBranchDto.ManagersIds,
			TimeZone: product_models.TimeZone{
				Location:  createBranchDto.TimeZone.Location,
				UtcOffset: createBranchDto.TimeZone.UtcOffset,
			},
		},
	}
	if createBranchDto.DistancePricingConfig != nil {
		branch.DistancePricingConfig = &product_models.DistancePricingConfig{
			BaseFee: product_models.PriceItem{
				Value:        createBranchDto.DistancePricingConfig.BaseFee.Value,
				CurrencyCode: createBranchDto.DistancePricingConfig.BaseFee.CurrencyCode,
				UnitId:       createBranchDto.DistancePricingConfig.BaseFee.UnitId,
			},
			BaseDistanceKm:   createBranchDto.DistancePricingConfig.BaseDistanceKm,
			FreeDistanceKm:   createBranchDto.DistancePricingConfig.FreeDistanceKm,
			EveryKmFee:       createBranchDto.DistancePricingConfig.EveryKmFee,
			ExtraCostPerKm: product_models.PriceItem{
				Value:        createBranchDto.DistancePricingConfig.ExtraCostPerKm.Value,
				CurrencyCode: createBranchDto.DistancePricingConfig.ExtraCostPerKm.CurrencyCode,
				UnitId:       createBranchDto.DistancePricingConfig.ExtraCostPerKm.UnitId,
			},
			EstimatedMinutesPerKm: createBranchDto.DistancePricingConfig.EstimatedMinutesPerKm,
		}
	}
	if createBranchDto.Zones != nil {
		branchZones := convertToBranchZoneModel(*createBranchDto.Zones)
		branch.Zones = &branchZones
	}
	return branch
}

func convertToBranchZoneModel(createBranchZoneDtos []dtos.CreateBranchZoneDto) []models.BranchZone {
	branchZones := make([]models.BranchZone, len(createBranchZoneDtos))
	for i, createBranchZoneDto := range createBranchZoneDtos {
		branchZones[i] = models.BranchZone{
			MinOrderAmount: createBranchZoneDto.MinOrderAmount,
			DeliveryFee: product_models.PriceItem{
				Value:        createBranchZoneDto.DeliveryFee.Value,
				CurrencyCode: createBranchZoneDto.DeliveryFee.CurrencyCode,
				UnitId:       createBranchZoneDto.DeliveryFee.UnitId,
			},
			EstimatedTime: createBranchZoneDto.EstimatedTime,
			ZoneId:        createBranchZoneDto.ZoneId,
		}
	}
	return branchZones
}
