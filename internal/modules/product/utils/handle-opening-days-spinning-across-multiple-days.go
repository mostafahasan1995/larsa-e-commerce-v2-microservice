// Package product_utils provides utility functions for handling product-related operations,
// particularly focusing on managing opening hours and time ranges.
package product_utils

import (
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"time"
)

var refFieldFrom26MayToWeekday = map[time.Weekday]int{
	time.Sunday:    25,
	time.Monday:    26,
	time.Tuesday:   27,
	time.Wednesday: 28,
	time.Thursday:  29,
	time.Friday:    30,
	time.Saturday:  31,
}

// HandleOpeningDaysSpinningAcrossMultipleDays processes opening days data and handles cases where
// business hours span across multiple days. It takes into account UTC offset and converts
// time ranges appropriately.
//
// The function:
// - Processes each day's opening hours
// - Converts time ranges from milliseconds to hours and minutes
// - Handles cases where business hours extend past midnight
// - Applies UTC offset adjustments
// - Returns a map of days to their corresponding time ranges
//
// Example:
//
//	openingDays := product_dtos.CreateOpeningDaysDto{
//	    Days: map[product_models.Day][]product_dtos.OpeningHours{
//	        product_models.Monday: {
//	            {
//	                From: 86400000,  // 24:00 (midnight)
//	                To:   90000000,  // 25:00 (1 AM next day)
//	            },
//	        },
//	    },
//	    UtcOffset: 2,
//	}
//	result := HandleOpeningDaysSpinningAcrossMultipleDays(openingDays)
//	Result will contain the time ranges split across days if they span midnight
//
// Parameters:
//   - utcOffset: The UTC offset for the branch
//   - openingDays: A DTO containing the opening hours for each day
//
// Returns:
//   - product_models.OpeningDays: A map of days to their processed time ranges
func HandleOpeningDaysSpinningAcrossMultipleDays(utcOffsetInMs int64, openingDays product_dtos.OpeningDaysDto) product_models.OpeningDays {
	openingDaysModel := product_models.OpeningDays{
		Days: make(map[time.Weekday]product_models.DaySchedule),
	}
	utcOffsetInSeconds := int64(utcOffsetInMs) / 1000
	for day, daySchedule := range openingDays.Days {
		// Initialize the day schedule if it doesn't exist
		if _, exists := openingDaysModel.Days[day]; !exists {
			openingDaysModel.Days[day] = product_models.DaySchedule{
				Label:      daySchedule.Label,
				TimeRanges: []product_models.DateRange{},
			}
		}

		for _, openingHour := range daySchedule.TimeRanges {
			from := MillisecondsToHoursAndMinutes(openingHour.From)
			to := MillisecondsToHoursAndMinutes(openingHour.To)
			constantDate := time.Date(2025, 5, refFieldFrom26MayToWeekday[day], 0, 0, 0, 0, time.UTC)
			startDate := time.Date(constantDate.Year(), constantDate.Month(), constantDate.Day(), int(from.Hours), int(from.Minutes), int(from.Seconds)+int(-1*utcOffsetInSeconds), 0, time.UTC)
			endDate := time.Date(constantDate.Year(), constantDate.Month(), constantDate.Day(), int(to.Hours), int(to.Minutes), int(to.Seconds)+int(-1*utcOffsetInSeconds), 0, time.UTC)

			dates := HandleDatesSpinningAcrossMultipleDays(startDate, endDate)

			for _, date := range dates {
				dayNumber := date.StartDate.Weekday()

				from := ExtractHoursAndMinutesFromDate(date.StartDate)
				to := ExtractHoursAndMinutesFromDate(date.EndDate)

				// Initialize the day schedule if it doesn't exist
				if _, exists := openingDaysModel.Days[dayNumber]; !exists {
					openingDaysModel.Days[dayNumber] = product_models.DaySchedule{
						Label:      daySchedule.Label,
						TimeRanges: []product_models.DateRange{},
					}
				}

				// Update the day schedule
				dayScheduleModel := openingDaysModel.Days[dayNumber]
				dayScheduleModel.TimeRanges = append(dayScheduleModel.TimeRanges, product_models.DateRange{
					From: from.Milliseconds,
					To:   to.Milliseconds,
				})
				openingDaysModel.Days[dayNumber] = dayScheduleModel
			}
		}
	}

	return openingDaysModel
}
