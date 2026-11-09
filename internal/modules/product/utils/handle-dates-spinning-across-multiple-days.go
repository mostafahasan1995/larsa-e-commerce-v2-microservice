package product_utils

import (
	"time"
)

// DateRange represents a start and end date pair
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// HandleDatesSpinningAcrossMultipleDays takes a start and end date and breaks them down into
// individual day segments if they span across multiple days.
func HandleDatesSpinningAcrossMultipleDays(startDate time.Time, endDate time.Time) []DateRange {
	/*
	 * E.x 2024-09-23T01:05:10.000+00:00 is gonna be 2024-09-23T00:00:00.000+00:00
	 */
	startDateBase := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)

	/*
	 * E.x 2024-09-24T01:05:10.000+00:00 is gonna be 2024-09-24T00:00:00.000+00:00
	 */
	endDateBase := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.UTC)

	differenceInDays := int(endDateBase.Sub(startDateBase).Hours() / 24)

	datesAcrossMultipleDays := make([]DateRange, 0)

	/*
	 * If it's in the same day we don't change anything.
	 */
	if differenceInDays == 0 {
		datesAcrossMultipleDays = append(datesAcrossMultipleDays, DateRange{
			StartDate: startDate,
			EndDate:   endDate,
		})
	} else {
		/*
		 * NOTE: Remember the starting point of the loop is 0 not 1.
		 * This important for later.
		 */

		/*
		 * We need to add one day to calculate the difference between all dates,
		 * cause the difference may be one day but that means we have two dates,
		 * and we need to calculate the dates of all days.
		 *
		 * E.x
		 * startPoint.checkedAt: 2024-09-23T15:30:00.000+00:00
		 * endPoint.checkedAt: 2024-09-24T15:30:00.000+00:00
		 * differenceInDays is 1, but we have to dates (two days).
		 * In order to solve this issue we add one
		 */
		differenceInDaysPlusOne := differenceInDays + 1
		for day := 0; day < differenceInDaysPlusOne; day++ {
			var newStartDate time.Time
			var newEndDate time.Time

			/*
			 * If day is 0 that mean it is the first day, so we have to keep it as it is.
			 *
			 *
			 * E.x
			 * startDate: 2024-09-23T15:30:00.000+00:00
			 * endDate: 2024-09-25T15:30:00.000+00:00
			 *
			 * The Result is
			 * startDate: 2024-09-23T15:30:00.000+00:00
			 */
			if day == 0 {
				newStartDate = startDate
			} else {
				/*
				 * We need to add the start of the next date.
				 *
				 * E.x
				 * startDate: 2024-09-23T15:30:00.000+00:00
				 * startDateBase: 2024-09-23T00:00:00.000+00:00
				 * endDate: 2024-09-25T15:30:00.000+00:00
				 *
				 * The Result is
				 * lets say it is the last day which is 2 so $$day is 2,
				 * and multiply $$day with 24 hours
				 * and is to be added to the start date base $startDateBase
				 * startDate: 2024-09-25T00:00:00.000+00:00
				 */
				newStartDate = startDateBase.AddDate(0, 0, day)
			}

			/*
			 * If day is the last day it means we have to take as it is.
			 *
			 * E.x
			 * startDate: 2024-09-23T15:30:00.000+00:00
			 * endDate: 2024-09-25T15:30:00.000+00:00
			 *
			 * the Result is
			 * endDate: 2024-09-25T15:30:00.000+00:00
			 */
			if day == differenceInDays {
				newEndDate = endDate
			} else {
				/*
				 * There are two cases here either it's the first day,
				 * or in between day.
				 */
				/*
				 * We have to add the end of the day,
				 * and we need to add n days to
				 * reach the current date of the
				 * starting point date,
				 * by using $$day.
				 *
				 * E.x
				 * startDate: 2024-09-23T15:30:00.000+00:00
				 * endDate: 2024-09-25T15:30:00.000+00:00
				 *
				 * The result is if the day is not 0
				 * startDate: 2024-09-24T00:00:00.000+00:00
				 * endDate: 2024-09-24T23:59:59.999+00:00
				 */
				dayEnd := startDateBase.AddDate(0, 0, day)
				newEndDate = time.Date(dayEnd.Year(), dayEnd.Month(), dayEnd.Day(), 23, 59, 59, 999999999, time.UTC)
			}

			datesAcrossMultipleDays = append(datesAcrossMultipleDays, DateRange{
				StartDate: newStartDate,
				EndDate:   newEndDate,
			})
		}
	}

	return datesAcrossMultipleDays
}
