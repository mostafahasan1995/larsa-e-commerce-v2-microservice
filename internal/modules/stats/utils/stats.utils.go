package stats_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/util"
	stats_dto "larsa-e-commerce-microservice-v2/internal/modules/stats/dto"
)

// RoundFloatToTwoDecimals rounds a float64 value to 2 decimal places
func RoundFloatToTwoDecimals(value float64) float64 {
	return util.RoundToDecimalPlaces(value, 2)
}

// RoundComparison rounds all float values in a Comparision struct to 2 decimal places
func RoundComparison(comp *stats_dto.Comparision) {
	if comp != nil {
		comp.Current = RoundFloatToTwoDecimals(comp.Current)
		comp.Prev = RoundFloatToTwoDecimals(comp.Prev)
		comp.Percentage = RoundFloatToTwoDecimals(comp.Percentage)
	}
}

// RoundTopCompanyStat rounds all float values in a TopCompanyStat struct to 2 decimal places
func RoundTopCompanyStat(stat *stats_dto.TopCompanyStat) {
	if stat != nil {
		stat.RevenueThisMonth = RoundFloatToTwoDecimals(stat.RevenueThisMonth)
		stat.TotalRevenue = RoundFloatToTwoDecimals(stat.TotalRevenue)
	}
}

// RoundCompanyOrderDetail rounds all float values in a CompanyOrderDetail struct to 2 decimal places
func RoundCompanyOrderDetail(detail *stats_dto.CompanyOrderDetail) {
	if detail != nil {
		detail.Amount = RoundFloatToTwoDecimals(detail.Amount)
	}
}

// RoundCouponUsageDetail rounds all float values in a CouponUsageDetail struct to 2 decimal places
func RoundCouponUsageDetail(detail *stats_dto.CouponUsageDetail) {
	if detail != nil {
		detail.DiscountValue = RoundFloatToTwoDecimals(detail.DiscountValue)
	}
}

// RoundCustomerActivityRecord rounds all float values in a CustomerActivityRecord struct to 2 decimal places
func RoundCustomerActivityRecord(record *stats_dto.CustomerActivityRecord) {
	if record != nil {
		record.TotalMonthSalesAmount = RoundFloatToTwoDecimals(record.TotalMonthSalesAmount)

		// Round coupon details
		for i := range record.CouponDetails {
			RoundCouponUsageDetail(&record.CouponDetails[i])
		}
	}
}

// RoundCouponUsageStat rounds all float values in a CouponUsageStat struct to 2 decimal places
func RoundCouponUsageStat(stat *stats_dto.CouponUsageStat) {
	if stat != nil {
		stat.TotalDiscountValue = RoundFloatToTwoDecimals(stat.TotalDiscountValue)
		stat.CouponValuePerUse = RoundFloatToTwoDecimals(stat.CouponValuePerUse)
	}
}

// RoundCouponOrderDetail rounds all float values in a CouponOrderDetail struct to 2 decimal places
func RoundCouponOrderDetail(detail *stats_dto.CouponOrderDetail) {
	if detail != nil {
		detail.OriginalAmount = RoundFloatToTwoDecimals(detail.OriginalAmount)
		detail.DiscountAmount = RoundFloatToTwoDecimals(detail.DiscountAmount)
		detail.FinalAmount = RoundFloatToTwoDecimals(detail.FinalAmount)
	}
}

// RoundOrdersSummaryDto rounds all float values in an OrdersSummaryDto struct to 2 decimal places
func RoundOrdersSummaryDto(summary *stats_dto.OrdersSummaryDto) {
	if summary != nil {
		summary.TotalOrderValue = RoundFloatToTwoDecimals(summary.TotalOrderValue)

		// Round values in ordersValuePerChannel map
		for key, value := range summary.OrdersValuePerChannel {
			summary.OrdersValuePerChannel[key] = RoundFloatToTwoDecimals(value)
		}
	}
}

// RoundCustomerOrderDto rounds all float values in a CustomerOrderDto struct to 2 decimal places
func RoundCustomerOrderDto(order *stats_dto.CustomerOrderDto) {
	if order != nil {
		order.InvoiceTotal = RoundFloatToTwoDecimals(order.InvoiceTotal)
	}
}

// RoundCouponUsageSummaryDto rounds all float values in a CouponUsageSummaryDto struct to 2 decimal places
func RoundCouponUsageSummaryDto(summary *stats_dto.CouponUsageSummaryDto) {
	if summary != nil {
		summary.TotalDiscount = RoundFloatToTwoDecimals(summary.TotalDiscount)
	}
}

// RoundCustomerCouponUsageDto rounds all float values in a CustomerCouponUsageDto struct to 2 decimal places
func RoundCustomerCouponUsageDto(usage *stats_dto.CustomerCouponUsageDto) {
	if usage != nil {
		usage.DiscountValue = RoundFloatToTwoDecimals(usage.DiscountValue)
		usage.InvoiceTotal = RoundFloatToTwoDecimals(usage.InvoiceTotal)
	}
}
