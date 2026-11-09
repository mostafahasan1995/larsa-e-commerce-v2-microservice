package coupon_utils

import (
	"fmt"
	"net/http"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
)

type CheckCouponRedemptionArgs struct {
	Coupon    *coupon_models.CouponResponse
	SubTotal  *float64
	OrderType *order_enums.OrderType
}

func CheckCouponRedemption(args CheckCouponRedemptionArgs) *helpers.ApiError {
	coupon := args.Coupon

	if coupon.UserCoupon != nil {
		if coupon.TotalUsages == coupon.UserCouponTotalUsages {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Coupon has reached the total usage limit",
			}
		}

		if coupon.UserCoupon.Usages != nil && coupon.Usage != nil && coupon.Usage.MaxUsagePerCustomer != nil {
			if int(coupon.UserCoupon.TotalUsages) == int(*coupon.Usage.MaxUsagePerCustomer) {
				return &helpers.ApiError{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Coupon max usage per customer has reached its limit",
				}
			}
		}

		if coupon.Usage != nil && coupon.Usage.Interval != nil && coupon.UserCoupon.Usages != nil {

			/*
			* Group usages by the interval type
			* This section tracks how many times a coupon has been used within specific time intervals
			* (daily, weekly, monthly, or yearly) to enforce usage limits.
			*
			* Example scenarios:
			* 1. Daily limit of 2 uses:
			*    - User uses coupon at 10:00 AM -> usageCount["2024-03-20"] = 1
			*    - User uses coupon at 2:00 PM -> usageCount["2024-03-20"] = 2
			*    - User tries to use at 6:00 PM -> Error: "Coupon usage limit reached for the current period"
			*
			* 2. Weekly limit of 3 uses:
			*    - User uses coupon on Monday -> usageCount["2024-W12"] = 1
			*    - User uses coupon on Wednesday -> usageCount["2024-W12"] = 2
			*    - User uses coupon on Friday -> usageCount["2024-W12"] = 3
			*    - User tries to use on Sunday -> Error: "Coupon usage limit reached for the current period"
			*
			* 3. Monthly limit of 5 uses:
			*    - User uses coupon 5 times in March -> usageCount["2024-03"] = 5
			*    - User tries to use again in March -> Error: "Coupon usage limit reached for the current period"
			*    - User can use again in April -> usageCount["2024-04"] = 1
			*
			* 4. Yearly limit of 10 uses:
			*    - User uses coupon 10 times in 2024 -> usageCount["2024"] = 10
			*    - User tries to use again in 2024 -> Error: "Coupon usage limit reached for the current period"
			*    - User can use again in 2025 -> usageCount["2025"] = 1
			 */
			usageCount := make(map[string]int)
			var nowKey = extractKey(coupon.Usage.Interval.Type, time.Now())

			for _, usage := range *coupon.UserCoupon.Usages {
				var key = extractKey(coupon.Usage.Interval.Type, usage.Timestamp)
				usageCount[key]++
			}

			// Check current period limit
			if usageCount[nowKey] >= int(coupon.Usage.Interval.Value) {
				return &helpers.ApiError{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Coupon usage limit reached for the current period",
					Details: &[]types.ValidationError{
						{
							Field:   "coupon_code",
							Message: fmt.Sprintf("Coupon usage limit reached for the current period: %v %v", coupon.Usage.Interval.Value, coupon.Usage.Interval.Type),
							Value:   coupon.Code,
						},
					},
				}
			}
		}
	}

	// Check minimum subtotal if provided
	if args.SubTotal != nil && coupon.Usage != nil && coupon.Usage.MinimumSubtotal != nil {
		if *args.SubTotal < *coupon.Usage.MinimumSubtotal {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Order subtotal does not meet the minimum requirement for this coupon",
			}
		}
	}

	// Check order type if provided
	if args.OrderType != nil && coupon.OrderType != nil {
		if *args.OrderType != *coupon.OrderType {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Coupon is not valid for this order type",
			}
		}
	}

	if coupon.Type == coupon_enums.AUTO_GENERATED {
		if coupon.UserCoupon == nil {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Coupon is not applied to you",
			}
		}

		return nil
	}

	if coupon.Type == coupon_enums.MANUAL {
		if coupon.Usage != nil && coupon.Usage.StartDate != nil && coupon.Usage.StartDate.After(time.Now()) {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Coupon is not active yet",
			}
		}

		if coupon.Usage != nil && coupon.Usage.EndDate != nil && coupon.Usage.EndDate.Before(time.Now()) {
			return &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Coupon has expired",
			}
		}

		return nil
	}

	return nil
}

func extractKey(usageType coupon_enums.CouponExpirationType, timestamp time.Time) string {
	var key string
	switch usageType {
	case coupon_enums.DAYS:
		// Format: "2006-01-02" (e.g., "2024-03-20")
		key = timestamp.Format("2006-01-02")
	case coupon_enums.WEEKS:
		// Format: "YYYY-WXX" (e.g., "2024-W12" for the 12th week of 2024)
		year, week := timestamp.ISOWeek()
		key = fmt.Sprintf("%d-W%02d", year, week)
	case coupon_enums.MONTHS:
		// Format: "2006-01" (e.g., "2024-03")
		key = timestamp.Format("2006-01")
	case coupon_enums.YEARS:
		// Format: "2006" (e.g., "2024")
		key = timestamp.Format("2006")
	}

	return key
}
