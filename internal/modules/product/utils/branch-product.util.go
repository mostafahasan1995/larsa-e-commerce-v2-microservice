package product_utils

import (
	"fmt"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	"time"
)

type MillisecondsToHoursAndMinutesResult struct {
	Hours        uint64
	Minutes      uint64
	Seconds      uint64
	Milliseconds uint64
	Full         string
}

// converts a base 10 number in milliseconds to a base 60 number in hours and minutes
func MillisecondsToHoursAndMinutes(milliseconds uint64) MillisecondsToHoursAndMinutesResult {
	hours := milliseconds / 3600000
	minutes := (milliseconds % 3600000) / 60000
	seconds := (milliseconds % 60000) / 1000
	return MillisecondsToHoursAndMinutesResult{
		Hours:        hours,
		Minutes:      minutes,
		Seconds:      seconds,
		Milliseconds: milliseconds,
		Full:         fmt.Sprintf("%d:%d:%d", hours, minutes, seconds),
	}
}

func ExtractHoursAndMinutesFromDate(date time.Time) MillisecondsToHoursAndMinutesResult {
	hours := date.Hour()
	minutes := date.Minute()
	seconds := date.Second()
	milliseconds := product_shared.ToMilliseconds(hours, minutes, seconds)
	return MillisecondsToHoursAndMinutesResult{
		Hours:        uint64(hours),
		Minutes:      uint64(minutes),
		Seconds:      uint64(seconds),
		Milliseconds: uint64(milliseconds),
		Full:         fmt.Sprintf("%d:%d:%d", hours, minutes, seconds),
	}
}
