package order_service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	product_service "larsa-e-commerce-microservice-v2/internal/modules/product/services"

	"github.com/samber/do"
)

// This service is created by Mohemmad Zarzar
type OrderReportsService interface {
	GetOrdersByFilter(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReport, *helpers.ApiError)
	GetOrdersByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReport, *helpers.ApiError)
	GetProductSalesReport(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError)
	GetOrdersByItemFilter(
		ctx context.Context,
		filter *order_dto.OrderItemFilterDTO,
	) ([]order_models.OrderResponse, *helpers.ApiError)
	GetOrdersStat(ctx context.Context, orders []order_models.OrderResponse) *order_models.OrderStats
	GetBranchSales(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError)
	GetBranchSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError)
	GetBranchSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError)
	GetProductSales(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError)
	GetProductSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError)
	GetProductSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError)

	GetFilteredOrdersStatsV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReportSales, *helpers.ApiError)
	// GetFilteredOrdersStats(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReportSales, *helpers.ApiError)
	GetOrdersStatSales(orders []order_models.OrderResponse) order_models.OrderStat
	GetProductSalesReportV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError)
	GetReportDependProductV2(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.AggregatedProductSales, *helpers.ApiError)
	GetReportDependProductV3(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.AggregatedProductSales, *helpers.ApiError)
}

type orderReportsService struct {
	repo          order_repos.OrderReportsRepo
	branchService product_service.BranchService
}

func NewOrderReportsService(i *do.Injector) (OrderReportsService, error) {
	repo := do.MustInvoke[order_repos.OrderReportsRepo](i)
	branchService := do.MustInvoke[product_service.BranchService](i)
	return &orderReportsService{
		repo:          repo,
		branchService: branchService,
	}, nil
}

func isRevenueEligible(order order_models.OrderResponse) bool {
	switch order.Status {
	case "completed", "accepted":
		return true
	case "pending":
		return order.PaymentInfo.Status == order_enums.PAID
	default:
		return false
	}
}

func (s *orderReportsService) GetOrdersStat(ctx context.Context, orders []order_models.OrderResponse) *order_models.OrderStats {

	var totalSales float64
	var totalOrders int
	itemSales := make(map[string]struct {
		name     string
		quantity uint64
	})

	for _, order := range orders {
		totalOrders++
		if isRevenueEligible(order) {
			totalSales += order.TotalPrice
		}

		for _, item := range order.OrderItems {
			itemIDStr := item.ItemId.Hex() // convert ObjectID to string
			entry := itemSales[itemIDStr]
			entry.name = item.Name["en"] // You can use .Ar if Arabic is preferred
			entry.quantity += item.Quantity
			itemSales[itemIDStr] = entry
		}
	}

	// Find top-selling item
	var topItemID, topItemName string
	var topQuantity uint64
	for id, entry := range itemSales {
		if entry.quantity > topQuantity {
			topItemID = id
			topItemName = entry.name
			topQuantity = entry.quantity
		}
	}

	avgOrderValue := 0.0
	if totalOrders > 0 {
		avgOrderValue = totalSales / float64(totalOrders)
	}

	stats := &order_models.OrderStats{
		TotalSales:       totalSales,
		NumberOfOrders:   totalOrders,
		AvgOrderValue:    avgOrderValue,
		TopSellingItemId: topItemID,
		TopSellingName:   topItemName,
		TopQuantity:      int(topQuantity),
	}
	return stats
}

func (s *orderReportsService) GetOrdersStatSales(orders []order_models.OrderResponse) order_models.OrderStat {
	customerSet := make(map[string]struct{})
	totalPurchases := len(orders)
	var totalQuantity uint64
	totalQuantity = 0.0
	totalAmount := 0.0

	for _, order := range orders {
		customerSet[order.BuyerInfo.Id.Hex()] = struct{}{}
		for _, item := range order.OrderItems {
			totalQuantity += item.Quantity
		}
		if isRevenueEligible(order) {
			totalAmount += order.TotalPrice
		}
	}

	return order_models.OrderStat{
		CustomersPurchasing: len(customerSet),
		TotalPurchases:      totalPurchases,
		TotalQuantity:       totalQuantity,
		TotalAmount:         totalAmount,
	}
}

func (s *orderReportsService) GetOrdersByItemFilter(
	ctx context.Context,
	filter *order_dto.OrderItemFilterDTO,
) ([]order_models.OrderResponse, *helpers.ApiError) {
	return s.repo.GetOrdersByItemFilter(ctx, filter)
}

func (s *orderReportsService) GetReportDependProductV2(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.AggregatedProductSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}

	aggregatedData, apiErr := s.repo.GetAggregatedProductSalesV2(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	branches, _ := s.branchService.GetAll(ctx, "")
	_ = branches // kept for parity; can be used later

	weekdayMap := map[string]string{"0": "Saturday", "1": "Sunday", "2": "Monday", "3": "Tuesday", "4": "Wednesday", "5": "Thursday", "6": "Friday"}
	result := make(map[string][]order_models.AggregatedProductSales)
	for _, data := range aggregatedData {
		dayKey := data.GroupKey
		if filter.Period == "This Week" {
			dayKey = weekdayMap[data.GroupKey]
		}
		result[dayKey] = append(result[dayKey], data)
	}
	return result, nil
}

func (s *orderReportsService) GetReportDependProductV3(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.AggregatedProductSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}
	// Validate date inputs
	if filter.StartDate == "" || filter.EndDate == "" {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "StartDate and EndDate are required"}
	}

	// Validate period type
	if filter.PeriodType == "" {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "PeriodType is required"}
	}

	// Parse dates
	startDate, err := time.Parse(time.RFC3339, filter.StartDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid StartDate format"}
	}

	endDate, err := time.Parse(time.RFC3339, filter.EndDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid EndDate format"}
	}

	// Generate all possible period keys for the date range
	periodKeys := generatePeriodKeys(startDate, endDate, filter.PeriodType)

	// Initialize result with all period keys and empty arrays
	result := make(map[string][]order_models.AggregatedProductSales)
	for _, key := range periodKeys {
		result[key] = []order_models.AggregatedProductSales{}
	}

	// Get aggregated data from repository
	aggregatedData, apiErr := s.repo.GetAggregatedProductSalesV3(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	// Fill in the actual data where available
	for _, data := range aggregatedData {
		groupKey := data.GroupKey

		// Format weekly key if needed
		if filter.PeriodType == "weekly" {
			parts := strings.Split(groupKey, "-")
			if len(parts) == 2 {
				groupKey = fmt.Sprintf("%s-W%s", parts[0], parts[1])
			}
		}

		result[groupKey] = append(result[groupKey], data)
	}

	return result, nil
}

func (s *orderReportsService) adjustFilterDates(filter *order_dto.OrderReportDTO) *helpers.ApiError {
	start, err1 := time.Parse(time.RFC3339, filter.StartDate)
	end, err2 := time.Parse(time.RFC3339, filter.EndDate)
	if err1 != nil || err2 != nil {
		return nil // Or handle as error if required. Returning nil for parity with existing logic.
	}

	startHour := "00:00"
	if filter.StartHour != "" {
		startHour = filter.StartHour
	}
	var sh, sm int
	fmt.Sscanf(startHour, "%d:%d", &sh, &sm)
	start = time.Date(start.Year(), start.Month(), start.Day(), sh, sm, 0, 0, start.Location())

	endHour := "23:59"
	if filter.EndHour != "" {
		endHour = filter.EndHour
	}
	var eh, em int
	fmt.Sscanf(endHour, "%d:%d", &eh, &em)
	end = time.Date(end.Year(), end.Month(), end.Day(), eh, em, 59, 999999999, end.Location())

	filter.StartDate = start.Format(time.RFC3339)
	filter.EndDate = end.Format(time.RFC3339)
	return nil
}

func (s *orderReportsService) GetProductSalesReportV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return order_models.ProductSalesReportResult{}, apiErr
	}

	start, _ := time.Parse(time.RFC3339, filter.StartDate)
	end, _ := time.Parse(time.RFC3339, filter.EndDate)

	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := end.Add(-duration)

	lastFilter := *filter
	lastFilter.StartDate = lastStart.Format(time.RFC3339)
	lastFilter.EndDate = lastEnd.Format(time.RFC3339)

	reportResult, err := s.repo.GetProductSalesAggregatedReportV2(ctx, filter)
	if err != nil {
		return order_models.ProductSalesReportResult{}, err
	}
	lastReportResult, err := s.repo.GetProductSalesAggregatedReportV2(ctx, &lastFilter)
	if err != nil {
		return order_models.ProductSalesReportResult{}, err
	}

	result := order_models.ProductSalesReportResult{
		Report:   reportResult.Report,
		Stat:     reportResult.Stat,
		LastStat: lastReportResult.Stat,
	}
	return result, nil
}

func (s *orderReportsService) GetProductSalesReport(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return order_models.ProductSalesReportResult{}, apiErr
	}

	start, _ := time.Parse(time.RFC3339, filter.StartDate)
	end, _ := time.Parse(time.RFC3339, filter.EndDate)

	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := end.Add(-duration)

	lastFilter := *filter
	lastFilter.StartDate = lastStart.Format(time.RFC3339)
	lastFilter.EndDate = lastEnd.Format(time.RFC3339)

	orders, err := s.repo.GetOrdersByFilter(ctx, filter)
	if err != nil {
		return order_models.ProductSalesReportResult{}, err
	}
	lastOrders, err := s.repo.GetOrdersByFilter(ctx, &lastFilter)
	if err != nil {
		return order_models.ProductSalesReportResult{}, err
	}

	var report []order_models.ProductSalesReport
	currentStats := make(map[string]struct {
		Qty  uint64
		Name string
	})
	var totalQuantity uint64
	var totalSales float64
	for _, order := range orders {
		isEligible := isRevenueEligible(order)
		for _, item := range order.OrderItems {
			entry := order_models.ProductSalesReport{ProductName: item.Name["en"], BranchId: item.BranchId, OrderTotal: item.TotalPrice, Quantity: item.Quantity, CustomerName: order.BuyerInfo.FullName, CustomerEmail: order.BuyerInfo.Email, CustomerAddress: order.BuyerInfo.Location.Address, PaymentMethod: order.PaymentInfo.PaymentMethod}
			report = append(report, entry)
			itemID := item.ItemId.Hex()
			stat := currentStats[itemID]
			stat.Qty += item.Quantity
			stat.Name = item.Name["en"]
			currentStats[itemID] = stat
			totalQuantity += item.Quantity
			if isEligible {
				totalSales += item.TotalPrice
			}
		}
	}
	stat := buildProductSalesStat(currentStats, totalQuantity, totalSales)

	lastStats := make(map[string]struct {
		Qty  uint64
		Name string
	})
	var lastTotalQuantity uint64
	var lastTotalSales float64
	for _, order := range lastOrders {
		isEligible := isRevenueEligible(order)
		for _, item := range order.OrderItems {
			itemID := item.ItemId.Hex()
			st := lastStats[itemID]
			st.Qty += item.Quantity
			st.Name = item.Name["en"]
			lastStats[itemID] = st
			lastTotalQuantity += item.Quantity
			if isEligible {
				lastTotalSales += item.TotalPrice
			}
		}
	}
	lastStat := buildProductSalesStat(lastStats, lastTotalQuantity, lastTotalSales)

	return order_models.ProductSalesReportResult{Report: report, Stat: stat, LastStat: lastStat}, nil
}

func (s *orderReportsService) GetOrdersByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReport, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return order_models.OrderReport{}, apiErr
	}

	start, _ := time.Parse(time.RFC3339, filter.StartDate)
	end, _ := time.Parse(time.RFC3339, filter.EndDate)

	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := end.Add(-duration)

	lastFilter := *filter
	lastFilter.StartDate = lastStart.Format(time.RFC3339)
	lastFilter.EndDate = lastEnd.Format(time.RFC3339)

	orders, err := s.repo.GetOrdersByFilterV2(ctx, filter)
	if err != nil {
		return order_models.OrderReport{}, err
	}

	stats, err := s.repo.GetOrderStatsByFilterV2(ctx, filter)
	if err != nil {
		return order_models.OrderReport{}, err
	}
	lastStats, err := s.repo.GetOrderStatsByFilterV2(ctx, &lastFilter)
	if err != nil {
		return order_models.OrderReport{}, err
	}

	return order_models.OrderReport{Orders: orders, Stat: stats, LastStat: lastStats}, nil
}

func (s *orderReportsService) GetProductSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}
	// Validate date inputs
	if filter.StartDate == "" || filter.EndDate == "" {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "StartDate and EndDate are required"}
	}

	// Parse dates to validate format
	_, err := time.Parse(time.RFC3339, filter.StartDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid StartDate format"}
	}

	_, err = time.Parse(time.RFC3339, filter.EndDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid EndDate format"}
	}

	return s.repo.GetProductSalesAggregationV3(ctx, filter)
}

func (s *orderReportsService) GetProductSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}
	return s.repo.GetProductSalesAggregationV2(ctx, filter)
}

func (s *orderReportsService) GetBranchSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}

	aggregatedData, apiErr := s.repo.GetAggregatedBranchSalesV2(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	branches, _ := s.branchService.GetAll(ctx, "")
	branchNameMap := make(map[string]string)
	for _, br := range branches {
		branchNameMap[br.Id.Hex()] = br.Name["en"]
	}
	weekdayMap := map[string]string{"0": "Saturday", "1": "Sunday", "2": "Monday", "3": "Tuesday", "4": "Wednesday", "5": "Thursday", "6": "Friday"}

	result := make(map[string][]order_models.BranchSales)
	for _, data := range aggregatedData {
		dayKey := data.GroupKey
		if filter.Period == "This Week" {
			dayKey = weekdayMap[data.GroupKey]
		}
		result[dayKey] = append(result[dayKey], order_models.BranchSales{BranchId: data.BranchId, BranchName: branchNameMap[data.BranchId], SalesValue: data.SalesValue, NumberOfOrders: data.NumberOfOrders})
	}
	return result, nil
}

func (s *orderReportsService) GetBranchSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}
	// Validate date inputs
	if filter.StartDate == "" || filter.EndDate == "" {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "StartDate and EndDate are required"}
	}

	// Validate period type
	if filter.PeriodType == "" {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "PeriodType is required"}
	}

	// Parse dates
	startDate, err := time.Parse(time.RFC3339, filter.StartDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid StartDate format"}
	}

	endDate, err := time.Parse(time.RFC3339, filter.EndDate)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid EndDate format"}
	}

	// Get all branches
	branches, _ := s.branchService.GetAllV2(ctx, "")
	branchNameMap := make(map[string]string)
	for _, br := range branches {
		branchNameMap[br.Id.Hex()] = br.Name["en"]
	}

	// Generate all possible period keys for the date range
	periodKeys := generatePeriodKeys(startDate, endDate, filter.PeriodType)

	// Initialize result with all period keys and empty arrays
	result := make(map[string][]order_models.BranchSales)
	for _, key := range periodKeys {
		result[key] = []order_models.BranchSales{}
	}

	// Get aggregated data from repository
	aggregatedData, apiErr := s.repo.GetAggregatedBranchSalesV3(ctx, filter)
	if apiErr != nil {
		return nil, apiErr
	}

	// Fill in the actual data where available
	for _, data := range aggregatedData {
		groupKey := data.GroupKey

		// Format weekly key if needed
		if filter.PeriodType == "weekly" {
			parts := strings.Split(groupKey, "-")
			if len(parts) == 2 {
				groupKey = fmt.Sprintf("%s-W%s", parts[0], parts[1])
			}
		}

		// Add branch sales data to the result
		result[groupKey] = append(result[groupKey], order_models.BranchSales{
			BranchId:       data.BranchId,
			BranchName:     branchNameMap[data.BranchId],
			SalesValue:     data.SalesValue,
			NumberOfOrders: data.NumberOfOrders,
		})
	}

	return result, nil
}

// Helper function to generate all period keys for a date range
func generatePeriodKeys(startDate, endDate time.Time, periodType string) []string {
	var keys []string
	current := startDate

	for current.Before(endDate) || current.Equal(endDate) {
		var key string

		switch periodType {
		case "daily":
			key = current.Format("2006-01-02")
			current = current.AddDate(0, 0, 1) // Next day
		case "weekly":
			year, week := current.ISOWeek()
			key = fmt.Sprintf("%d-W%02d", year, week)
			current = current.AddDate(0, 0, 7) // Next week
		case "monthly":
			key = current.Format("2006-01")
			current = current.AddDate(0, 1, 0) // Next month
		case "yearly":
			key = current.Format("2006")
			current = current.AddDate(1, 0, 0) // Next year
		}

		keys = append(keys, key)
	}

	return keys
}
func (s *orderReportsService) GetBranchSales(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.BranchSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}

	branches, _ := s.branchService.GetAll(ctx, "")
	fmt.Println("branches")
	fmt.Println(branches)

	orders, _ := s.repo.GetOrdersByFilter(ctx, filter)
	groupMap := make(map[string]map[string]float64)
	orderCountMap := make(map[string]map[string]int)

	// We need the original start and end dates for boundary checks if needed
	startDate, _ := time.Parse(time.RFC3339, filter.StartDate)
	endDate, _ := time.Parse(time.RFC3339, filter.EndDate)
	for _, order := range orders {
		isEligible := isRevenueEligible(order)
		t := order.CreatedAt.UTC()
		var groupKey string
		switch filter.Period {
		case "This Week":
			if t.Before(startDate) || t.After(endDate) {
				continue
			}
			groupKey = strings.ToLower(t.Weekday().String())
		case "This Month":
			if t.Before(startDate) || t.After(endDate) {
				continue
			}
			groupKey = fmt.Sprintf("%d", t.Day())
		case "This Year":
			if t.Before(startDate) || t.After(endDate) {
				continue
			}
			groupKey = strings.ToLower(t.Month().String()[:3])
		default:
			continue
		}
		if _, ok := groupMap[groupKey]; !ok {
			groupMap[groupKey] = make(map[string]float64)
		}
		if _, ok := orderCountMap[groupKey]; !ok {
			orderCountMap[groupKey] = make(map[string]int)
		}
		countedBranches := make(map[string]bool)
		for _, item := range order.OrderItems {
			branchID := item.BranchId.Hex()
			if isEligible {
				groupMap[groupKey][branchID] += item.TotalPrice
			}
			if !countedBranches[branchID] {
				orderCountMap[groupKey][branchID]++
				countedBranches[branchID] = true
			}
		}
	}
	result := make(map[string][]order_models.BranchSales)
	for groupKey, branchSales := range groupMap {
		var sales []order_models.BranchSales
		for branchId, total := range branchSales {
			branchName := ""
			for _, br := range branches {
				if br.Id.Hex() == branchId {
					branchName = br.Name["en"]
				}
			}
			sales = append(sales, order_models.BranchSales{BranchId: branchId, BranchName: branchName, SalesValue: total, NumberOfOrders: orderCountMap[groupKey][branchId]})
		}
		result[groupKey] = sales
	}
	return result, nil
}

func (s *orderReportsService) GetOrdersByFilter(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReport, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return order_models.OrderReport{}, apiErr
	}

	start, _ := time.Parse(time.RFC3339, filter.StartDate)
	end, _ := time.Parse(time.RFC3339, filter.EndDate)

	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := end.Add(-duration)
	lastFilter := *filter
	lastFilter.StartDate = lastStart.Format(time.RFC3339)
	lastFilter.EndDate = lastEnd.Format(time.RFC3339)
	orders, err := s.repo.GetOrdersByFilter(ctx, filter)
	if err != nil {
		return order_models.OrderReport{}, err
	}
	lastOrders, err := s.repo.GetOrdersByFilter(ctx, &lastFilter)
	if err != nil {
		return order_models.OrderReport{}, err
	}
	stats := s.GetOrdersStat(ctx, orders)
	lastStats := s.GetOrdersStat(ctx, lastOrders)
	return order_models.OrderReport{Orders: orders, Stat: stats, LastStat: lastStats}, nil
}

func (s *orderReportsService) GetProductSales(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return nil, apiErr
	}

	orders, err := s.repo.GetOrdersByFilter(ctx, filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Fail to fetch orders"}
	}
	productMap := make(map[string]*order_models.ProductSales)
	for _, order := range orders {
		isEligible := isRevenueEligible(order)
		for _, item := range order.OrderItems {
			name := item.Name["en"]
			if name == "" {
				continue
			}
			if _, exists := productMap[name]; !exists {
				productMap[name] = &order_models.ProductSales{Name: name}
			}
			productMap[name].Quantity += item.Quantity
			if isEligible {
				productMap[name].SalesValue += item.TotalPrice
			}
		}
	}
	var productSalesList []order_models.ProductSales
	for _, p := range productMap {
		productSalesList = append(productSalesList, *p)
	}
	sort.Slice(productSalesList, func(i, j int) bool { return productSalesList[i].SalesValue > productSalesList[j].SalesValue })
	return productSalesList, nil
}

func (s *orderReportsService) GetFilteredOrdersStatsV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.OrderReportSales, *helpers.ApiError) {
	if apiErr := s.adjustFilterDates(filter); apiErr != nil {
		return order_models.OrderReportSales{}, apiErr
	}

	start, _ := time.Parse(time.RFC3339, filter.StartDate)
	end, _ := time.Parse(time.RFC3339, filter.EndDate)

	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := end.Add(-duration)
	lastFilter := *filter
	lastFilter.StartDate = lastStart.Format(time.RFC3339)
	lastFilter.EndDate = lastEnd.Format(time.RFC3339)
	currentData, err := s.repo.GetOrdersWithStats(ctx, filter)
	if err != nil {
		return order_models.OrderReportSales{}, err
	}
	lastData, err := s.repo.GetOrdersWithStats(ctx, &lastFilter)
	if err != nil {
		return order_models.OrderReportSales{}, err
	}
	return order_models.OrderReportSales{Customers: currentData.Customers, Stat: currentData.Stats[0], LastStat: lastData.Stats[0]}, nil
}

func buildProductSalesStat(stats map[string]struct {
	Qty  uint64
	Name string
}, totalQuantity uint64, totalSales float64) order_models.ProductSalesStat {
	bestProductName := ""
	maxQty := uint64(0)

	for _, entry := range stats {
		if entry.Qty > maxQty {
			maxQty = entry.Qty
			bestProductName = entry.Name
		}
	}

	return order_models.ProductSalesStat{
		TotalProductsSold:  len(stats),
		TotalQuantity:      totalQuantity,
		TotalSales:         totalSales,
		BestSellingProduct: bestProductName,
	}
}

func isoWeekStart(year int, week int, loc *time.Location) time.Time {
	// ISO weeks: Monday is first day of the week
	// Week 1 is the week with the first Thursday of the year

	// January 4 is always in the first ISO week
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, loc)

	// Find the Monday of that week
	isoWeekDay := int(jan4.Weekday())
	if isoWeekDay == 0 {
		isoWeekDay = 7
	}
	monday := jan4.AddDate(0, 0, -isoWeekDay+1)

	// Add the number of weeks
	return monday.AddDate(0, 0, (week-1)*7)
}
