package stats_service

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	stats_dto "larsa-e-commerce-microservice-v2/internal/modules/stats/dto"
	stats_repos "larsa-e-commerce-microservice-v2/internal/modules/stats/repos"
	stats_utils "larsa-e-commerce-microservice-v2/internal/modules/stats/utils"
	"math"
	"strings"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatsOrderService interface {
	GetOrderStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.OrderStatsResp, error)
	GetSalesStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.OrderStatsResp, error)
	GetInvoiceStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.InvoiceStatsResp, error)

	GetTopCompanies(ctx context.Context, filter stats_dto.TopCompaniesFilterDto) (*stats_dto.TopCompaniesResp, error)
	GetCompanyOrders(ctx context.Context, filter stats_dto.CompanyOrdersFilterDto) (*stats_dto.CompanyOrdersWithPagination, error)
	GetCustomersActivity(ctx context.Context, filter stats_dto.CustomersActivityFilterDto) (*stats_dto.CustomersActivityResp, error)

	// Coupons usage
	GetTopCoupons(ctx context.Context, filter stats_dto.TopCouponsFilterDto) (*types.DataWithPagination[stats_dto.CouponUsageStat], error)
	GetCouponOrders(ctx context.Context, filter stats_dto.CouponOrdersFilterDto) (*stats_dto.CouponOrdersWithPagination, error)

	// Customer Profile
	GetCustomerProfile(ctx context.Context, filter stats_dto.CustomerProfileFilterDto) (*stats_dto.CustomerProfileResponseDto, error)
	GetCustomerOrders(ctx context.Context, filter stats_dto.CustomerOrdersFilterDto) (*stats_dto.CustomerOrdersTableDto, error)
	GetCustomerCouponUsageDetails(ctx context.Context, customerId string, startDate, endDate *time.Time, skip, limit int64) (*stats_dto.CustomerCouponUsageDetailsDto, error)

	// Notification Center
	GetGenericData(ctx context.Context, dto stats_dto.GenericFilter) ([]bson.M, error)

	// Charts
	GetCouponRedemptionsPerBranch(ctx context.Context, filter stats_dto.CouponRedemptionsPerBranchFilterDto) (*stats_dto.CouponRedemptionsPerBranchResponseDto, error)
	GetCouponRedemptionsByCustomerType(ctx context.Context, filter stats_dto.CouponRedemptionsByCustomerTypeFilterDto) (*stats_dto.CouponRedemptionsByCustomerTypeResponseDto, error)
	GetCustomersCouponRedeemsBar(ctx context.Context, filter stats_dto.CustomersCouponRedeemsBarFilterDto) (*stats_dto.CustomersCouponRedeemsBarResponseDto, error)
}

type statsOrderService struct {
	orderRepo stats_repos.StatsOrderRepo
	invRepo   stats_repos.StatsInvoiceRepo
	db        *mongo.Client
}

func NewStatsOrderService(i *do.Injector) (StatsOrderService, error) {
	return &statsOrderService{
		orderRepo: do.MustInvoke[stats_repos.StatsOrderRepo](i),
		invRepo:   do.MustInvoke[stats_repos.StatsInvoiceRepo](i),
		db:        do.MustInvoke[*mongo.Client](i),
	}, nil
}

func (s *statsOrderService) GetGenericData(ctx context.Context, dto stats_dto.GenericFilter) ([]bson.M, error) {
	cfg, _ := util.GetReqAppCfg(ctx)

	pipeLine := stats_dto.BuildGenericMongoPipeline(dto)
	coll := s.db.Database(cfg.Db).Collection(dto.Collection)

	cursor, err := coll.Aggregate(ctx, pipeLine)
	if err != nil {
		return nil, fmt.Errorf("aggregate failed: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("cursor decode failed: %w", err)
	}

	return results, nil
}

func (s *statsOrderService) GetOrderStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.OrderStatsResp, error) {
	resp := &stats_dto.OrderStatsResp{}
	fmt.Println("start: ", filter.StartDate)
	fmt.Println("end: ", filter.EndDate)

	// Parse dates
	startDate, err1 := time.Parse(time.RFC3339, filter.StartDate)
	endDate, err2 := time.Parse(time.RFC3339, filter.EndDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format: %v, %v", err1, err2)
	}

	active, err := s.orderRepo.GetActiveCustomersCount(ctx, startDate, endDate, true, &filter.BranchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active customers count: %v", err)
	}
	resp.ActiveCount = active
	inActive, err := s.orderRepo.GetActiveCustomersCount(ctx, startDate, endDate, false, &filter.BranchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get active customers count: %v", err)
	}
	resp.InActiveCount = inActive

	// Build tiny line chart series for Active/Inactive counts over the selected range (max 7 buckets)
	// Determine number of buckets based on total days, capped at 7
	totalDuration := endDate.Sub(startDate)
	if totalDuration < 0 {
		return nil, fmt.Errorf("endDate must be after startDate")
	}
	// Calculate number of days (ceil) in range
	dayCount := int(math.Ceil(totalDuration.Hours() / 24.0))
	numBuckets := 7
	if dayCount > 0 && dayCount < numBuckets {
		numBuckets = dayCount
	}
	if numBuckets < 1 {
		numBuckets = 1
	}
	step := time.Duration(0)
	if numBuckets > 0 {
		step = totalDuration / time.Duration(numBuckets)
	}

	activeSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)
	inactiveSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)
	for i := 0; i < numBuckets; i++ {
		bucketStart := startDate.Add(step * time.Duration(i))
		var bucketEnd time.Time
		if i == numBuckets-1 {
			bucketEnd = endDate
		} else {
			bucketEnd = startDate.Add(step * time.Duration(i+1)).Add(-time.Nanosecond)
		}

		ac, err := s.orderRepo.GetActiveCustomersCount(ctx, bucketStart, bucketEnd, true, &filter.BranchId)
		if err != nil {
			return nil, fmt.Errorf("failed to get active customers count for bucket %d: %v", i, err)
		}
		inactive, err := s.orderRepo.GetActiveCustomersCount(ctx, bucketStart, bucketEnd, false, &filter.BranchId)
		if err != nil {
			return nil, fmt.Errorf("failed to get inactive customers count for bucket %d: %v", i, err)
		}
		activeSeries = append(activeSeries, stats_dto.TinySeriesPoint{Value: ac, Date: bucketStart.Format(time.RFC3339)})
		inactiveSeries = append(inactiveSeries, stats_dto.TinySeriesPoint{Value: inactive, Date: bucketStart.Format(time.RFC3339)})
	}
	resp.ActiveSeries = activeSeries
	resp.InActiveSeries = inactiveSeries

	// Total orders per bucket for TinyLineChart of totalCount
	totalSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)
	for i := 0; i < numBuckets; i++ {
		bucketStart := startDate.Add(step * time.Duration(i))
		var bucketEnd time.Time
		if i == numBuckets-1 {
			bucketEnd = endDate
		} else {
			bucketEnd = startDate.Add(step * time.Duration(i+1)).Add(-time.Nanosecond)
		}
		cnt, err := s.orderRepo.CountOrders(ctx, bucketStart, bucketEnd, &filter.BranchId)
		if err != nil {
			return nil, fmt.Errorf("failed to count orders for bucket %d: %v", i, err)
		}
		totalSeries = append(totalSeries, stats_dto.TinySeriesPoint{Value: cnt, Date: bucketStart.Format(time.RFC3339)})
	}
	resp.TotalCountSeries = totalSeries

	ordersCount, err := s.orderRepo.CountOrders(ctx, startDate, endDate, &filter.BranchId)
	if err != nil {
		return nil, fmt.Errorf("failed to count orders: %v", err)
	}
	resp.TotalCount = float64(ordersCount)

	currentCount, previousCount, percentage, err := s.orderRepo.GetOrdersComparisonByDateRange(ctx, startDate, endDate, &filter.BranchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders comparison by date range: %v", err)
	}
	resp.OrdersComparision = stats_dto.Comparision{
		Current:    float64(currentCount),
		Prev:       float64(previousCount),
		Percentage: percentage,
	}
	currentGrowth, prevGrowth, percentBetween, err := s.orderRepo.GetGrowthRateComparison(ctx, startDate, endDate, &filter.BranchId)
	if err != nil {
		return nil, fmt.Errorf("failed to get growth rate %v", err)
	}
	resp.GrowthRateComp = stats_dto.Comparision{
		Current:    currentGrowth,
		Prev:       prevGrowth,
		Percentage: percentBetween,
	}

	// Round float values to 2 decimal places
	stats_utils.RoundComparison(&resp.OrdersComparision)
	stats_utils.RoundComparison(&resp.GrowthRateComp)
	resp.TotalCount = stats_utils.RoundFloatToTwoDecimals(resp.TotalCount)

	return resp, nil
}

func (s *statsOrderService) GetSalesStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.OrderStatsResp, error) {
	resp := &stats_dto.OrderStatsResp{}

	// Parse dates
	startDate, err1 := time.Parse(time.RFC3339, filter.StartDate)
	endDate, err2 := time.Parse(time.RFC3339, filter.EndDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format: %v, %v", err1, err2)
	}

	currentGrowth, prevGrowth, percentBetween, err := s.orderRepo.GetSalesGrowthRateComparison(ctx, startDate, endDate, &filter.BranchId)
	if err != nil {
		return nil, err
	}
	resp.GrowthRateComp = stats_dto.Comparision{
		Current:    currentGrowth,
		Prev:       prevGrowth,
		Percentage: percentBetween,
	}

	currentSales, previousSales, percentage, err := s.orderRepo.GetSalesComparisonByDateRange(ctx, startDate, endDate, &filter.BranchId)
	if err != nil {
		return nil, err
	}
	resp.OrdersComparision = stats_dto.Comparision{
		Current:    currentSales,
		Prev:       previousSales,
		Percentage: float64(percentage),
	}

	total, err := s.orderRepo.GetTotalSalesSinceLaunch(ctx, &filter.BranchId)
	if err != nil {
		return nil, err
	}
	resp.TotalCount = float64(total)
	resp.ActiveCount = 0

	// Round float values to 2 decimal places
	stats_utils.RoundComparison(&resp.GrowthRateComp)
	stats_utils.RoundComparison(&resp.OrdersComparision)
	resp.TotalCount = stats_utils.RoundFloatToTwoDecimals(resp.TotalCount)

	return resp, nil
}

func (s *statsOrderService) GetInvoiceStats(ctx context.Context, filter stats_dto.StatsFilterDto) (*stats_dto.InvoiceStatsResp, error) {
	resp := &stats_dto.InvoiceStatsResp{}
	startDate, err1 := time.Parse(time.RFC3339, filter.StartDate)
	endDate, err2 := time.Parse(time.RFC3339, filter.EndDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format: %v, %v", err1, err2)
	}

	pCount, err := s.orderRepo.CountPayment(ctx, order_enums.PAID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	resp.PaidInvCount = pCount

	upCount, err := s.orderRepo.CountPayment(ctx, order_enums.UNPAID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	resp.UnpaidInvCount = upCount

	partCount, err := s.orderRepo.CountPayment(ctx, order_enums.PARTIALLY_PAID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	resp.PartiallyPaid = partCount

	// Build tiny line chart series for invoice statuses over the selected range (max 7 buckets)
	totalDuration := endDate.Sub(startDate)
	if totalDuration < 0 {
		return nil, fmt.Errorf("endDate must be after startDate")
	}
	dayCount := int(math.Ceil(totalDuration.Hours() / 24.0))
	numBuckets := 7
	if dayCount > 0 && dayCount < numBuckets {
		numBuckets = dayCount
	}
	if numBuckets < 1 {
		numBuckets = 1
	}
	step := time.Duration(0)
	if numBuckets > 0 {
		step = totalDuration / time.Duration(numBuckets)
	}

	paidSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)
	unpaidSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)
	partiallyPaidSeries := make([]stats_dto.TinySeriesPoint, 0, numBuckets)

	for i := 0; i < numBuckets; i++ {
		bucketStart := startDate.Add(step * time.Duration(i))
		var bucketEnd time.Time
		if i == numBuckets-1 {
			bucketEnd = endDate
		} else {
			bucketEnd = startDate.Add(step * time.Duration(i+1)).Add(-time.Nanosecond)
		}

		paid, err := s.orderRepo.CountPayment(ctx, order_enums.PAID, bucketStart, bucketEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to count paid invoices for bucket %d: %v", i, err)
		}
		unpaid, err := s.orderRepo.CountPayment(ctx, order_enums.UNPAID, bucketStart, bucketEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to count unpaid invoices for bucket %d: %v", i, err)
		}
		partial, err := s.orderRepo.CountPayment(ctx, order_enums.PARTIALLY_PAID, bucketStart, bucketEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to count partially paid invoices for bucket %d: %v", i, err)
		}

		paidSeries = append(paidSeries, stats_dto.TinySeriesPoint{Value: paid, Date: bucketStart.Format(time.RFC3339)})
		unpaidSeries = append(unpaidSeries, stats_dto.TinySeriesPoint{Value: unpaid, Date: bucketStart.Format(time.RFC3339)})
		partiallyPaidSeries = append(partiallyPaidSeries, stats_dto.TinySeriesPoint{Value: partial, Date: bucketStart.Format(time.RFC3339)})
	}
	resp.PaidInvoicesSeries = paidSeries
	resp.UnpaidInvoicesSeries = unpaidSeries
	resp.PartiallyPaidInvoicesSeries = partiallyPaidSeries
	return resp, nil
}

func (s *statsOrderService) GetTopCompanies(ctx context.Context, filter stats_dto.TopCompaniesFilterDto) (*stats_dto.TopCompaniesResp, error) {
	startDate, err1 := time.Parse(time.RFC3339, filter.StartDate)
	endDate, err2 := time.Parse(time.RFC3339, filter.EndDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format: %v, %v", err1, err2)
	}
	stats, err := s.orderRepo.GetTopCompanies(ctx, startDate, endDate, &filter.BranchId, filter.Limit, filter.Skip)
	if err != nil {
		return nil, err
	}

	// Round float values to 2 decimal places
	for i := range stats.Data {
		stats_utils.RoundTopCompanyStat(&stats.Data[i])
	}

	return &stats, nil
}

func (s *statsOrderService) GetCompanyOrders(ctx context.Context, filter stats_dto.CompanyOrdersFilterDto) (*stats_dto.CompanyOrdersWithPagination, error) {
	// Parse companyId
	oid, err := primitive.ObjectIDFromHex(strings.TrimSpace(filter.CompanyId))
	if err != nil {
		return nil, fmt.Errorf("invalid companyId: %v", err)
	}
	// Parse optional dates
	var startPtr, endPtr *time.Time
	if filter.StartDate != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if filter.EndDate != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}

	result, err := s.orderRepo.GetCompanyOrders(ctx, oid, startPtr, endPtr, &filter.BranchId, filter.Page, filter.Limit)
	if err != nil {
		return nil, err
	}

	// Round float values to 2 decimal places
	for i := range result.Data {
		stats_utils.RoundCompanyOrderDetail(&result.Data[i])
	}

	return result, nil
}

func (s *statsOrderService) GetCustomersActivity(ctx context.Context, filter stats_dto.CustomersActivityFilterDto) (*stats_dto.CustomersActivityResp, error) {
	rows, totalCount, topCustomerId, err := s.orderRepo.GetCustomersActivity(ctx, filter)
	if err != nil {
		return nil, err
	}
	resp := stats_dto.CustomersActivityResp{
		Data: rows,
		Pagination: types.Pagination{
			TotalPages: func() float64 {
				if filter.Limit == 0 {
					return 0
				}
				return math.Ceil(float64(totalCount) / float64(filter.Limit))
			}(),
			PerPage:    filter.Limit,
			TotalCount: totalCount,
		},
	}

	// Round float values to 2 decimal places and identify top customer
	for i := range resp.Data {
		stats_utils.RoundCustomerActivityRecord(&resp.Data[i])
	}

	// Flag the top customer based on the ID returned from repository
	if topCustomerId != nil {
		for i := range resp.Data {
			if resp.Data[i].CustomerId == *topCustomerId {
				resp.Data[i].IsTopCustomer = true
				break
			}
		}
	}

	return &resp, nil
}

func (s *statsOrderService) GetTopCoupons(ctx context.Context, filter stats_dto.TopCouponsFilterDto) (*types.DataWithPagination[stats_dto.CouponUsageStat], error) {
	stats, topCouponCode, err := s.orderRepo.GetTopCoupons(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Round float values to 2 decimal places
	for i := range stats.Data {
		stats_utils.RoundCouponUsageStat(&stats.Data[i])
	}

	// Flag the top coupon based on the code returned from repository
	if topCouponCode != nil {
		for i := range stats.Data {
			if stats.Data[i].CouponCode == *topCouponCode {
				stats.Data[i].IsTopCoupon = true
				break
			}
		}
	}

	return &stats, nil
}

func (s *statsOrderService) GetCouponOrders(ctx context.Context, filter stats_dto.CouponOrdersFilterDto) (*stats_dto.CouponOrdersWithPagination, error) {
	var startPtr, endPtr *time.Time
	if strings.TrimSpace(filter.StartDate) != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if strings.TrimSpace(filter.EndDate) != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}
	var sourcePtr, statusPtr *string
	if strings.TrimSpace(filter.OrderSource) != "" {
		v := filter.OrderSource
		sourcePtr = &v
	}
	if strings.TrimSpace(filter.Status) != "" {
		v := filter.Status
		statusPtr = &v
	}
	result, err := s.orderRepo.GetCouponOrders(ctx, filter.CouponCode, startPtr, endPtr, &filter.BranchId, sourcePtr, statusPtr, filter.Page, filter.Limit)
	if err != nil {
		return nil, err
	}

	// Round float values to 2 decimal places
	for i := range result.Data {
		stats_utils.RoundCouponOrderDetail(&result.Data[i])
	}

	return result, nil
}

// Customer Profile Methods

func (s *statsOrderService) GetCustomerProfile(ctx context.Context, filter stats_dto.CustomerProfileFilterDto) (*stats_dto.CustomerProfileResponseDto, error) {
	customerId, err := primitive.ObjectIDFromHex(filter.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %v", err)
	}

	var startPtr, endPtr *time.Time
	if filter.StartDate != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if filter.EndDate != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}

	// Get customer profile information
	profile, err := s.orderRepo.GetCustomerProfile(ctx, customerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer profile: %v", err)
	}

	// Get orders summary
	ordersSummary, err := s.orderRepo.GetCustomerOrdersSummary(ctx, customerId, startPtr, endPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer orders summary: %v", err)
	}

	// Get coupon usage summary
	couponUsageSummary, err := s.orderRepo.GetCustomerCouponUsageSummary(ctx, customerId, startPtr, endPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer coupon usage summary: %v", err)
	}

	// Round float values to 2 decimal places
	stats_utils.RoundOrdersSummaryDto(ordersSummary)
	stats_utils.RoundCouponUsageSummaryDto(couponUsageSummary)

	return &stats_dto.CustomerProfileResponseDto{
		Profile:            *profile,
		OrdersSummary:      *ordersSummary,
		CouponUsageSummary: *couponUsageSummary,
	}, nil
}

func (s *statsOrderService) GetCustomerOrders(ctx context.Context, filter stats_dto.CustomerOrdersFilterDto) (*stats_dto.CustomerOrdersTableDto, error) {
	customerId, err := primitive.ObjectIDFromHex(filter.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %v", err)
	}

	var startPtr, endPtr *time.Time
	if filter.StartDate != nil && *filter.StartDate != "" {
		sd, e := time.Parse(time.RFC3339, *filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if filter.EndDate != nil && *filter.EndDate != "" {
		ed, e := time.Parse(time.RFC3339, *filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}

	// Calculate skip from page and limit
	skip := int64((filter.Page - 1) * filter.Limit)
	limit := int64(filter.Limit)

	result, err := s.orderRepo.GetCustomerOrders(ctx, customerId, startPtr, endPtr, filter.OrderStatus, filter.BranchId, filter.Channel, filter.PaymentMethod, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer orders: %v", err)
	}

	return result, nil
}

func (s *statsOrderService) GetCustomerCouponUsageDetails(ctx context.Context, customerId string, startDate, endDate *time.Time, skip, limit int64) (*stats_dto.CustomerCouponUsageDetailsDto, error) {
	customerObjId, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %v", err)
	}

	coupons, total, err := s.orderRepo.GetCustomerCouponUsageDetails(ctx, customerObjId, startDate, endDate, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer coupon usage details: %v", err)
	}

	return &stats_dto.CustomerCouponUsageDetailsDto{
		Data: coupons,
		Pagination: types.Pagination{
			TotalCount: total,
			PerPage:    limit,
			TotalPages: math.Ceil(float64(total) / float64(limit)),
		},
	}, nil
}

// Charts Methods

func (s *statsOrderService) GetCouponRedemptionsPerBranch(ctx context.Context, filter stats_dto.CouponRedemptionsPerBranchFilterDto) (*stats_dto.CouponRedemptionsPerBranchResponseDto, error) {
	var startPtr, endPtr *time.Time

	if strings.TrimSpace(filter.StartDate) != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if strings.TrimSpace(filter.EndDate) != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}

	var couponCodePtr, paymentMethodPtr *string
	if strings.TrimSpace(filter.CouponCode) != "" {
		v := filter.CouponCode
		couponCodePtr = &v
	}
	if strings.TrimSpace(filter.PaymentMethod) != "" {
		v := filter.PaymentMethod
		paymentMethodPtr = &v
	}

	data, err := s.orderRepo.GetCouponRedemptionsPerBranch(ctx, startPtr, endPtr, couponCodePtr, paymentMethodPtr)
	if err != nil {
		return nil, err
	}

	return &stats_dto.CouponRedemptionsPerBranchResponseDto{
		Data: data,
	}, nil
}

func (s *statsOrderService) GetCouponRedemptionsByCustomerType(ctx context.Context, filter stats_dto.CouponRedemptionsByCustomerTypeFilterDto) (*stats_dto.CouponRedemptionsByCustomerTypeResponseDto, error) {
	var startPtr, endPtr *time.Time
	if strings.TrimSpace(filter.StartDate) != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if strings.TrimSpace(filter.EndDate) != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}
	var couponCodePtr *string
	if strings.TrimSpace(filter.CouponCode) != "" {
		v := filter.CouponCode
		couponCodePtr = &v
	}
	// if strings.TrimSpace(filter.PaymentMethod) != "" {
	// 	v := filter.PaymentMethod
	// 	paymentMethodPtr = &v
	// }

	data, err := s.orderRepo.GetCouponRedemptionsByCustomerType(ctx, startPtr, endPtr, couponCodePtr)
	if err != nil {
		return nil, err
	}
	return &stats_dto.CouponRedemptionsByCustomerTypeResponseDto{Data: data}, nil
}

func (s *statsOrderService) GetCustomersCouponRedeemsBar(ctx context.Context, filter stats_dto.CustomersCouponRedeemsBarFilterDto) (*stats_dto.CustomersCouponRedeemsBarResponseDto, error) {
	var startPtr, endPtr *time.Time
	if strings.TrimSpace(filter.StartDate) != "" {
		sd, e := time.Parse(time.RFC3339, filter.StartDate)
		if e != nil {
			return nil, fmt.Errorf("invalid startDate: %v", e)
		}
		startPtr = &sd
	}
	if strings.TrimSpace(filter.EndDate) != "" {
		ed, e := time.Parse(time.RFC3339, filter.EndDate)
		if e != nil {
			return nil, fmt.Errorf("invalid endDate: %v", e)
		}
		endPtr = &ed
	}
	rows, err := s.orderRepo.GetCustomersCouponRedeemsBar(ctx, startPtr, endPtr, &filter.BranchId, filter.Limit, filter.Skip)
	if err != nil {
		return nil, err
	}
	// Round numeric values
	for i := range rows {
		rows[i].TotalSales = stats_utils.RoundFloatToTwoDecimals(rows[i].TotalSales)
	}
	return &stats_dto.CustomersCouponRedeemsBarResponseDto{Data: rows}, nil
}
