package product_shared

const MAX_HOURS_IN_A_DAY = 23
const MAX_MINUTES_IN_A_HOUR = 59
const MAX_SECONDS_IN_A_MINUTE = 59
const MAX_MILLISECONDS_IN_A_SECOND = 1000

func ToMilliseconds(hours int, minutes int, seconds int) int64 {
	return int64(hours*3600000 + minutes*60000 + seconds*1000)
}

var MAX_TIME_SPAN_IN_A_DAY = uint64(ToMilliseconds(MAX_HOURS_IN_A_DAY, MAX_MINUTES_IN_A_HOUR, MAX_SECONDS_IN_A_MINUTE))

var MINIMUM_UTC_OFFSET = ToMilliseconds(-12, 0, 0)
var MAXIMUM_UTC_OFFSET = ToMilliseconds(14, 0, 0)
