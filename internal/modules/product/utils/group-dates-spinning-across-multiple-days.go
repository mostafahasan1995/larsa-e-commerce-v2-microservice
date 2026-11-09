package product_utils

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	"math"
	"time"
)

func GroupDatesSpinningAcrossMultipleDays(utcOffsetInMs int64, openingDays *product_models.OpeningDays) *product_models.OpeningDays {
	if openingDays == nil {
		return nil
	}

	if utcOffsetInMs == 0 {
		return openingDays
	}

	newOpeningHours := product_models.OpeningDays{
		Days: make(map[time.Weekday]product_models.DaySchedule),
	}
	toBeExcluded := product_models.OpeningDays{
		Days: make(map[time.Weekday]product_models.DaySchedule),
	}

	for day, daySchedule := range openingDays.Days {
		for _, dateRange := range daySchedule.TimeRanges {
			newDay := day
			newFrom := int64(dateRange.From)
			newTo := int64(dateRange.To)

			previousDay := getPreviousDay(day)
			isPositiveUtcOffset := utcOffsetInMs > 0

			// Initialize maps if they don't exist
			if _, exists := newOpeningHours.Days[previousDay]; !exists {
				newOpeningHours.Days[previousDay] = product_models.DaySchedule{
					Label:      openingDays.Days[previousDay].Label,
					TimeRanges: []product_models.DateRange{},
				}
			}
			if _, exists := toBeExcluded.Days[previousDay]; !exists {
				toBeExcluded.Days[previousDay] = product_models.DaySchedule{
					Label:      openingDays.Days[previousDay].Label,
					TimeRanges: []product_models.DateRange{},
				}
			}

			previousDateRange := util.Find(openingDays.Days[previousDay].TimeRanges, func(el product_models.DateRange) bool {
				return el.To == product_shared.MAX_TIME_SPAN_IN_A_DAY
			})

			if dateRange.From == 0 && previousDateRange != nil {
				previousDay := getPreviousDay(day)
				newFrom = 0

				utcOffsetAbs := uint64(math.Abs(float64(utcOffsetInMs)))

				if previousDateRange.From+utcOffsetAbs < product_shared.MAX_TIME_SPAN_IN_A_DAY {
					// Update toBeExcluded for day
					toBeExcludedDay := toBeExcluded.Days[day]
					toBeExcludedDay.TimeRanges = append(toBeExcludedDay.TimeRanges, dateRange)
					toBeExcluded.Days[day] = toBeExcludedDay

					// Update toBeExcluded for previousDay
					toBeExcludedPrevDay := toBeExcluded.Days[previousDay]
					toBeExcludedPrevDay.TimeRanges = append(toBeExcludedPrevDay.TimeRanges, *previousDateRange)
					toBeExcluded.Days[previousDay] = toBeExcludedPrevDay

					newDay = previousDay
					newFrom = int64(previousDateRange.From) - int64(utcOffsetAbs)
					newTo = int64(previousDateRange.To) - (int64(utcOffsetAbs) - int64(dateRange.To))

					// Update newOpeningHours for newDay
					newOpeningDay := newOpeningHours.Days[newDay]
					newOpeningDay.TimeRanges = append(newOpeningDay.TimeRanges, product_models.DateRange{
						From: uint64(newFrom),
						To:   uint64(newTo) + 1000,
					})
					newOpeningHours.Days[newDay] = newOpeningDay

				} else if previousDateRange.From+utcOffsetAbs > product_shared.MAX_TIME_SPAN_IN_A_DAY {
					// Update toBeExcluded for day
					toBeExcludedDay := toBeExcluded.Days[day]
					toBeExcludedDay.TimeRanges = append(toBeExcludedDay.TimeRanges, dateRange)
					toBeExcluded.Days[day] = toBeExcludedDay

					// Update toBeExcluded for previousDay
					toBeExcludedPrevDay := toBeExcluded.Days[previousDay]
					toBeExcludedPrevDay.TimeRanges = append(toBeExcludedPrevDay.TimeRanges, *previousDateRange)
					toBeExcluded.Days[previousDay] = toBeExcludedPrevDay

					newDay = day
					newFrom = int64(math.Abs(float64((int64(product_shared.MAX_TIME_SPAN_IN_A_DAY) - int64(previousDateRange.From)) - (int64(utcOffsetAbs) - 1000))))
					newTo = int64(dateRange.To) + int64(utcOffsetInMs)

					// Update newOpeningHours for newDay
					newOpeningDay := newOpeningHours.Days[newDay]
					newOpeningDay.TimeRanges = append(newOpeningDay.TimeRanges, product_models.DateRange{
						From: uint64(newFrom),
						To:   uint64(newTo),
					})
					newOpeningHours.Days[newDay] = newOpeningDay
				}

			} else {
				newDay = day
				newFrom = int64(dateRange.From) + int64(utcOffsetInMs)
				newTo = int64(dateRange.To) + int64(utcOffsetInMs)

				// Initialize maps if they don't exist
				if _, exists := newOpeningHours.Days[getNextDay(day)]; !exists {
					newOpeningHours.Days[getNextDay(day)] = product_models.DaySchedule{
						Label:      openingDays.Days[getNextDay(day)].Label,
						TimeRanges: []product_models.DateRange{},
					}
				}

				nextDateRange := util.Find(openingDays.Days[getNextDay(day)].TimeRanges, func(el product_models.DateRange) bool {
					return el.From == 0
				})

				isThisDateRangeSpanningAcrossMultipleDays := false

				if nextDateRange != nil && dateRange.To == product_shared.MAX_TIME_SPAN_IN_A_DAY {
					isThisDateRangeSpanningAcrossMultipleDays = true
				}

				if previousDateRange != nil && dateRange.From == 0 {
					isThisDateRangeSpanningAcrossMultipleDays = true
				}

				if !isThisDateRangeSpanningAcrossMultipleDays {
					if isPositiveUtcOffset {

						newFrom = int64(dateRange.From) + int64(utcOffsetInMs)
						newTo = int64(dateRange.To) + int64(utcOffsetInMs)

						if newFrom > int64(product_shared.MAX_TIME_SPAN_IN_A_DAY) && newTo > int64(product_shared.MAX_TIME_SPAN_IN_A_DAY) {

							newFrom = newFrom - int64(product_shared.MAX_TIME_SPAN_IN_A_DAY) - 1000
							newTo = newTo - int64(product_shared.MAX_TIME_SPAN_IN_A_DAY) - 1000
							newDay = getNextDay(newDay)
						}

						// Initialize map if it doesn't exist
						if _, exists := newOpeningHours.Days[newDay]; !exists {
							newOpeningHours.Days[newDay] = product_models.DaySchedule{
								Label:      daySchedule.Label,
								TimeRanges: []product_models.DateRange{},
							}
						}

						// Update newOpeningHours for newDay
						newOpeningDay := newOpeningHours.Days[newDay]
						newOpeningDay.TimeRanges = append(newOpeningDay.TimeRanges, product_models.DateRange{
							From: uint64(newFrom),
							To:   uint64(newTo),
						})
						newOpeningHours.Days[newDay] = newOpeningDay
					} else {

						newFrom = int64(dateRange.From) + int64(utcOffsetInMs)
						newTo = int64(dateRange.To) + int64(utcOffsetInMs)

						// Initialize map if it doesn't exist
						if _, exists := newOpeningHours.Days[newDay]; !exists {
							newOpeningHours.Days[newDay] = product_models.DaySchedule{
								Label:      daySchedule.Label,
								TimeRanges: []product_models.DateRange{},
							}
						}

						// Update newOpeningHours for newDay
						newOpeningDay := newOpeningHours.Days[newDay]
						newOpeningDay.TimeRanges = append(newOpeningDay.TimeRanges, product_models.DateRange{
							From: uint64(newFrom),
							To:   uint64(newTo),
						})
						newOpeningHours.Days[newDay] = newOpeningDay
					}
				}
			}
		}
	}

	fmt.Println("newOpeningHours", newOpeningHours)
	return &newOpeningHours
}

func getPreviousDay(day time.Weekday) time.Weekday {
	if day >= 1 && day <= 6 {
		return day - 1
	}

	if day == 0 {
		return 6
	}

	return 0

}

func getNextDay(day time.Weekday) time.Weekday {
	if day >= 0 && day <= 5 {
		return day + 1
	}

	if day == 6 {
		return 1
	}

	return 0
}
