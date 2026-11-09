package stats_handlers

import (
	"encoding/json"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	stats_dto "larsa-e-commerce-microservice-v2/internal/modules/stats/dto"
	stats_order_service "larsa-e-commerce-microservice-v2/internal/modules/stats/stats_service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

type StatsHandler struct {
	orderSvc stats_order_service.StatsOrderService
	excel    excel_services.ExcelService
}

func NewStatsHandler(i *do.Injector, r *chi.Mux) {
	h := StatsHandler{
		orderSvc: do.MustInvoke[stats_order_service.StatsOrderService](i),
		excel:    do.MustInvoke[excel_services.ExcelService](i),
	}

	r.Route(enums.STATS.GetAdminRoute(), func(r chi.Router) {
		r.Get("/orders", helpers.Make(h.OrderStats))
		r.Get("/sales", helpers.Make(h.SalesStats))
		r.Get("/invoices", helpers.Make(h.InvoicesStats))

		r.Get("/top-companies", helpers.Make(h.TopCompanies))
		r.Get("/top-companies/export", helpers.Make(h.ExportTopCompanies))
		r.Get("/company-orders", helpers.Make(h.CompanyOrders))
		r.Get("/customers-activity", helpers.Make(h.CustomersActivity))
		r.Get("/customers-activity/export", helpers.Make(h.ExportCustomersActivity))
		r.Get("/top-coupons", helpers.Make(h.TopCoupons))
		r.Get("/top-coupons/export", helpers.Make(h.ExportTopCoupons))
		r.Get("/coupon-orders", helpers.Make(h.CouponOrders))
		r.Get("/coupon-orders/export", helpers.Make(h.ExportCouponOrders))

		// Customer Profile endpoints
		r.Get("/customer-profile", helpers.Make(h.GetCustomerProfile))
		r.Get("/customer-orders", helpers.Make(h.GetCustomerOrders))
		r.Get("/customer-coupon-usage", helpers.Make(h.GetCustomerCouponUsageDetails))

		// Charts
		r.Get("/coupon-redemptions-per-branch", helpers.Make(h.GetCouponRedemptionsPerBranch))
		r.Get("/coupon-redemptions-by-customer-type", helpers.Make(h.GetCouponRedemptionsByCustomerType))
		r.Get("/customers-coupon-redeems-bar", helpers.Make(h.GetCustomersCouponRedeemsBar))
		// Notification Center
		r.Post("/generic", helpers.Make(h.GetGenericData))
	})
}

func (h *StatsHandler) GetGenericData(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var dto stats_dto.GenericFilter
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		apiErr := helpers.BadRequest(fmt.Sprintf("Invalid request body: %v", err))
		return &apiErr
	}

	data, err := h.orderSvc.GetGenericData(ctx, dto)
	if err != nil {
		return err

	}

	return helpers.WriteJson(w, http.StatusOK, data)
}

func (h *StatsHandler) OrderStats(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get startDate and endDate from query params, default to current month if not provided
	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)

	filterDto := stats_dto.StatsFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
	}

	stats, err := h.orderSvc.GetOrderStats(ctx, filterDto)
	if err != nil {
		return err
	}

	if stats == nil {
		apiE := helpers.BadRequest("Stats service returned nil response")
		return &apiE
	}

	return helpers.WriteJson(w, http.StatusOK, stats)
}
func (h *StatsHandler) SalesStats(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get startDate and endDate from query params, default to current month if not provided
	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)

	filterDto := stats_dto.StatsFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
	}
	stats, err := h.orderSvc.GetSalesStats(ctx, filterDto)
	if err != nil {
		return err
	}

	if stats == nil {
		apiE := helpers.BadRequest("Stats service returned nil response")
		return &apiE
	}

	return helpers.WriteJson(w, http.StatusOK, stats)
}

func (h *StatsHandler) InvoicesStats(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	// Get startDate and endDate from query params, default to current month if not provided
	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)

	filterDto := stats_dto.StatsFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
	}
	stats, err := h.orderSvc.GetInvoiceStats(ctx, filterDto)
	if err != nil {
		return err
	}

	if stats == nil {
		apiE := helpers.BadRequest("Stats service returned nil response")
		return &apiE
	}

	return helpers.WriteJson(w, http.StatusOK, stats)
}

// TopCompanies returns ranked companies for the current month
func (h *StatsHandler) TopCompanies(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	filter := stats_dto.TopCompaniesFilterDto{
		Limit:     limit,
		Skip:      skip,
		StartDate: util.GetDateOrDefault(r.URL.Query().Get("startDate"), true),
		EndDate:   util.GetDateOrDefault(r.URL.Query().Get("endDate"), false),
		BranchId:  r.URL.Query().Get("branchId"),
	}
	resp, err := h.orderSvc.GetTopCompanies(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

func (h *StatsHandler) ExportTopCompanies(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	filter := stats_dto.TopCompaniesFilterDto{
		StartDate: util.GetDateOrDefault(q.Get("startDate"), true),
		EndDate:   util.GetDateOrDefault(q.Get("endDate"), false),
		BranchId:  q.Get("branchId"),
		Limit:     limit,
		Skip:      skip,
	}

	// fetch all rows (make sure this calls GetTopCompaniesAll, not the paged one)
	rows, err := h.orderSvc.GetTopCompanies(ctx, filter)
	if err != nil {
		return err
	}

	columns := []excel_services.ExcelColumn{
		{Key: "companyId", Header: "Company ID", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "companyName", Header: "Company Name", Type: excel_services.ExcelTypeString},
		{Key: "ordersThisMonth", Header: "Orders This Month", Type: excel_services.ExcelTypeInt},
		{Key: "totalOrders", Header: "Total Orders", Type: excel_services.ExcelTypeInt},
		{Key: "revenueThisMonth", Header: "Revenue This Month", Type: excel_services.ExcelTypeFloat},
		{Key: "totalRevenue", Header: "Total Revenue", Type: excel_services.ExcelTypeFloat},
	}

	data := make([]map[string]any, 0, len(rows.Data))
	for _, rrow := range rows.Data {
		data = append(data, map[string]any{
			"companyId": rrow.CompanyId,
			"companyName": func() string {
				if rrow.CompanyName != nil {
					return *rrow.CompanyName
				}
				return ""
			}(),
			"ordersThisMonth":  rrow.OrdersThisMonth,
			"totalOrders":      rrow.TotalOrders,
			"revenueThisMonth": rrow.RevenueThisMonth,
			"totalRevenue":     rrow.TotalRevenue,
		})
	}

	bytes, apiErr := h.excel.Export(ctx, data, columns, "Top Companies")
	if apiErr != nil {
		return apiErr
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=top-companies.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
	return nil
}

// CompanyOrders returns paginated orders for a given company with optional date range
func (h *StatsHandler) CompanyOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	q := r.URL.Query()
	filter := stats_dto.CompanyOrdersFilterDto{
		CompanyId: q.Get("companyId"),
		StartDate: util.GetDateOrDefault(r.URL.Query().Get("startDate"), true),
		EndDate:   util.GetDateOrDefault(r.URL.Query().Get("endDate"), false),
		Page:      int(skip/int64(limit)) + 1, // Convert skip/limit to page number
		Limit:     int(limit),
	}

	resp, err := h.orderSvc.GetCompanyOrders(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

// CustomersActivity returns customers with orders >= threshold this month
func (h *StatsHandler) CustomersActivity(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	fmt.Println("skip, limit :=", skip, limit)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	// Parse query parameters
	threshold := 5
	if t := r.URL.Query().Get("minOrders"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n > 0 {
			threshold = n
		}
	}

	var maxOrders *int
	if t := r.URL.Query().Get("maxOrders"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n > 0 {
			maxOrders = &n
		}
	}

	var minTotalQuantity *int
	if t := r.URL.Query().Get("minTotalQuantity"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n >= 0 {
			minTotalQuantity = &n
		}
	}

	var maxTotalQuantity *int
	if t := r.URL.Query().Get("maxTotalQuantity"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n >= 0 {
			maxTotalQuantity = &n
		}
	}

	var minTotalSalesAmount *float64
	if t := r.URL.Query().Get("minTotalSalesAmount"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil && n >= 0 {
			minTotalSalesAmount = &n
		}
	}

	var maxTotalSalesAmount *float64
	if t := r.URL.Query().Get("maxTotalSalesAmount"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil && n >= 0 {
			maxTotalSalesAmount = &n
		}
	}

	var minAverageOrderValue *float64
	if t := r.URL.Query().Get("minAverageOrderValue"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil && n >= 0 {
			minAverageOrderValue = &n
		}
	}

	var maxAverageOrderValue *float64
	if t := r.URL.Query().Get("maxAverageOrderValue"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil && n >= 0 {
			maxAverageOrderValue = &n
		}
	}

	var minUsedCouponsCount *int
	if t := r.URL.Query().Get("minUsedCouponsCount"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n >= 0 {
			minUsedCouponsCount = &n
		}
	}

	var maxUsedCouponsCount *int
	if t := r.URL.Query().Get("maxUsedCouponsCount"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n >= 0 {
			maxUsedCouponsCount = &n
		}
	}

	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)
	branchId := r.URL.Query().Get("branchId")

	customerName := r.URL.Query().Get("customerName")
	email := r.URL.Query().Get("email")
	phoneNumber := r.URL.Query().Get("phoneNumber")
	customerType := r.URL.Query().Get("customerType")
	lastBranchOrderedFrom := r.URL.Query().Get("lastBranchOrderedFrom")
	minLastOrderDate := r.URL.Query().Get("minLastOrderDate")
	maxLastOrderDate := r.URL.Query().Get("maxLastOrderDate")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	if sortDir != "" && sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	filter := stats_dto.CustomersActivityFilterDto{
		StartDate:             startDate,
		EndDate:               endDate,
		MinOrders:             threshold,
		MaxOrders:             maxOrders,
		BranchId:              branchId,
		CustomerName:          customerName,
		Email:                 email,
		PhoneNumber:           phoneNumber,
		MinTotalQuantity:      minTotalQuantity,
		MaxTotalQuantity:      maxTotalQuantity,
		MinTotalSalesAmount:   minTotalSalesAmount,
		MaxTotalSalesAmount:   maxTotalSalesAmount,
		MinAverageOrderValue:  minAverageOrderValue,
		MaxAverageOrderValue:  maxAverageOrderValue,
		MinUsedCouponsCount:   minUsedCouponsCount,
		MaxUsedCouponsCount:   maxUsedCouponsCount,
		CustomerType:          customerType,
		LastBranchOrderedFrom: lastBranchOrderedFrom,
		MinLastOrderDate:      minLastOrderDate,
		MaxLastOrderDate:      maxLastOrderDate,
		SortBy:                sortBy,
		SortDir:               sortDir,
		Limit:                 limit,
		Skip:                  skip,
	}
	resp, err := h.orderSvc.GetCustomersActivity(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

// ExportCustomersActivity exports the same report as Excel
func (h *StatsHandler) ExportCustomersActivity(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	threshold := 5
	if t := r.URL.Query().Get("minOrders"); t != "" {
		if n, err := strconv.Atoi(t); err == nil && n > 0 {
			threshold = n
		}
	}
	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)
	branchId := r.URL.Query().Get("branchId")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}
	filter := stats_dto.CustomersActivityFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
		MinOrders: threshold,
		BranchId:  branchId,
		SortBy:    sortBy,
		SortDir:   sortDir,
		Skip:      skip,
		Limit:     limit,
	}
	resp, err := h.orderSvc.GetCustomersActivity(ctx, filter)
	if err != nil {
		return err
	}

	// Prepare export columns
	columns := []excel_services.ExcelColumn{
		{Key: "customerName", Header: "Customer Name", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "companyName", Header: "Company Name", Type: excel_services.ExcelTypeString},
		{Key: "email", Header: "Email", Type: excel_services.ExcelTypeEmail, IsRequired: true},
		{Key: "phoneNumber", Header: "Phone Number", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "totalOrdersCount", Header: "Total Orders Count", Type: excel_services.ExcelTypeInt},
		{Key: "totalQuantityOrdered", Header: "Total Quantity Ordered", Type: excel_services.ExcelTypeInt},
		{Key: "totalMonthSalesAmount", Header: "Total Month Sales Amount", Type: excel_services.ExcelTypeFloat},
		{Key: "lastOrderDate", Header: "Last Order Date", Type: excel_services.ExcelTypeDateTime, DateFormat: "2006-01-02 15:04:05"},
		{Key: "usedCouponsCount", Header: "Used Coupons Count", Type: excel_services.ExcelTypeInt},
	}

	// Flatten data for export
	data := make([]map[string]any, 0, len(resp.Data))
	for _, row := range resp.Data {
		data = append(data, map[string]any{
			"customerName":          row.CustomerName,
			"companyName":           row.CompanyName,
			"email":                 *row.Email,
			"phoneNumber":           *row.PhoneNumber,
			"totalOrdersCount":      row.TotalOrdersCount,
			"totalQuantityOrdered":  row.TotalQuantityOrdered,
			"totalMonthSalesAmount": row.TotalMonthSalesAmount,
			"lastOrderDate":         row.LastOrderDate.Format("2006-01-02 15:04:05"),
			"usedCouponsCount":      row.UsedCouponsCount,
		})
	}

	bytes, apiErr := h.excel.Export(ctx, data, columns, "Customers Monthly Activity")
	if apiErr != nil {
		return apiErr
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=customers-monthly-activity.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
	return nil
}

// TopCoupons returns ranked coupons by usage within date range
func (h *StatsHandler) TopCoupons(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	q := r.URL.Query()

	startDate := util.GetDateOrDefault(q.Get("startDate"), true)
	endDate := util.GetDateOrDefault(q.Get("endDate"), false)
	sortBy := q.Get("sortBy")
	sortDir := q.Get("sortDir")
	if sortDir != "" && sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	filter := stats_dto.TopCouponsFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
		BranchId:  q.Get("branchId"),
		SortBy:    sortBy,
		SortDir:   sortDir,
		Limit:     limit,
		Skip:      skip,
	}

	// Parse additional filters
	if couponCode := q.Get("couponCode"); couponCode != "" {
		filter.CouponCode = couponCode
	}

	if couponType := q.Get("couponType"); couponType != "" {
		filter.CouponType = couponType
	}

	if status := q.Get("status"); status != "" {
		filter.Status = status
	}

	if minRedeemCount := q.Get("minRedeemCount"); minRedeemCount != "" {
		if n, err := strconv.Atoi(minRedeemCount); err == nil && n >= 0 {
			filter.MinRedeemCount = &n
		}
	}

	if maxRedeemCount := q.Get("maxRedeemCount"); maxRedeemCount != "" {
		if n, err := strconv.Atoi(maxRedeemCount); err == nil && n >= 0 {
			filter.MaxRedeemCount = &n
		}
	}

	if minUniqueCustomers := q.Get("minUniqueCustomers"); minUniqueCustomers != "" {
		if n, err := strconv.Atoi(minUniqueCustomers); err == nil && n >= 0 {
			filter.MinUniqueCustomers = &n
		}
	}

	if maxUniqueCustomers := q.Get("maxUniqueCustomers"); maxUniqueCustomers != "" {
		if n, err := strconv.Atoi(maxUniqueCustomers); err == nil && n >= 0 {
			filter.MaxUniqueCustomers = &n
		}
	}

	if minTotalDiscountValue := q.Get("minTotalDiscountValue"); minTotalDiscountValue != "" {
		if f, err := strconv.ParseFloat(minTotalDiscountValue, 64); err == nil && f >= 0 {
			filter.MinTotalDiscountValue = &f
		}
	}

	if maxTotalDiscountValue := q.Get("maxTotalDiscountValue"); maxTotalDiscountValue != "" {
		if f, err := strconv.ParseFloat(maxTotalDiscountValue, 64); err == nil && f >= 0 {
			filter.MaxTotalDiscountValue = &f
		}
	}

	if minCouponValuePerUse := q.Get("minCouponValuePerUse"); minCouponValuePerUse != "" {
		if f, err := strconv.ParseFloat(minCouponValuePerUse, 64); err == nil && f >= 0 {
			filter.MinCouponValuePerUse = &f
		}
	}

	if maxCouponValuePerUse := q.Get("maxCouponValuePerUse"); maxCouponValuePerUse != "" {
		if f, err := strconv.ParseFloat(maxCouponValuePerUse, 64); err == nil && f >= 0 {
			filter.MaxCouponValuePerUse = &f
		}
	}

	resp, err := h.orderSvc.GetTopCoupons(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

// ExportTopCoupons exports Top Coupons to Excel
func (h *StatsHandler) ExportTopCoupons(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	filter := stats_dto.TopCouponsFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
		BranchId:  r.URL.Query().Get("branchId"),
		SortBy:    sortBy,
		SortDir:   sortDir,
		Limit:     limit,
		Skip:      skip,
	}

	resp, err := h.orderSvc.GetTopCoupons(ctx, filter)
	if err != nil {
		return err
	}

	columns := []excel_services.ExcelColumn{
		{Key: "couponCode", Header: "Coupon Code", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "redeemCount", Header: "Redeem Count", Type: excel_services.ExcelTypeInt},
		{Key: "uniqueCustomers", Header: "Customers Count", Type: excel_services.ExcelTypeInt},
		{Key: "totalDiscountValue", Header: "Total Discount Value", Type: excel_services.ExcelTypeFloat},
		{Key: "couponValuePerUse", Header: "Coupon Value Per Use", Type: excel_services.ExcelTypeFloat},
	}

	data := make([]map[string]any, 0, len(resp.Data))
	for _, row := range resp.Data {
		data = append(data, map[string]any{
			"couponCode":         row.CouponCode,
			"redeemCount":        row.RedeemCount,
			"uniqueCustomers":    row.UniqueCustomers,
			"totalDiscountValue": row.TotalDiscountValue,
			"couponValuePerUse":  row.CouponValuePerUse,
		})
	}

	bytes, apiErr := h.excel.Export(ctx, data, columns, "Top Coupons")
	if apiErr != nil {
		return apiErr
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=top-coupons.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)

	return nil
}

// CouponOrders returns paginated orders that used a specific coupon with filters
func (h *StatsHandler) CouponOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	q := r.URL.Query()
	filter := stats_dto.CouponOrdersFilterDto{
		CouponCode:  q.Get("couponCode"),
		StartDate:   util.GetDateOrDefault(r.URL.Query().Get("startDate"), true),
		EndDate:     util.GetDateOrDefault(r.URL.Query().Get("endDate"), false),
		OrderSource: q.Get("orderSource"),
		Status:      q.Get("status"),
		Page:        int(skip/int64(limit)) + 1, // Convert skip/limit to page number
		Limit:       int(limit),
	}
	if strings.TrimSpace(filter.CouponCode) == "" {
		apiE := helpers.BadRequest("couponCode is required")
		return &apiE
	}
	resp, err := h.orderSvc.GetCouponOrders(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

// ExportCouponOrders exports coupon orders to Excel
func (h *StatsHandler) ExportCouponOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	q := r.URL.Query()
	filter := stats_dto.CouponOrdersFilterDto{
		CouponCode:  q.Get("couponCode"),
		StartDate:   util.GetDateOrDefault(r.URL.Query().Get("startDate"), true),
		EndDate:     util.GetDateOrDefault(r.URL.Query().Get("endDate"), false),
		OrderSource: q.Get("orderSource"),
		Status:      q.Get("status"),
	}
	if strings.TrimSpace(filter.CouponCode) == "" {
		apiE := helpers.BadRequest("couponCode is required")
		return &apiE
	}

	rows, err := h.orderSvc.GetCouponOrders(ctx, filter)
	if err != nil {
		return err
	}

	columns := []excel_services.ExcelColumn{
		{Key: "orderId", Header: "Order ID", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "orderDate", Header: "Order Date", Type: excel_services.ExcelTypeDateTime, DateFormat: "2006-01-02 15:04:05"},
		{Key: "orderSource", Header: "Order Source", Type: excel_services.ExcelTypeString},
		{Key: "customerName", Header: "Customer Name", Type: excel_services.ExcelTypeString, IsRequired: true},
		{Key: "originalOrderAmount", Header: "Original Amount", Type: excel_services.ExcelTypeFloat},
		{Key: "discountAmount", Header: "Discount Amount", Type: excel_services.ExcelTypeFloat},
		{Key: "finalOrderAmount", Header: "Final Amount", Type: excel_services.ExcelTypeFloat},
		{Key: "orderStatus", Header: "Status", Type: excel_services.ExcelTypeString},
	}

	data := make([]map[string]any, 0, len(rows.Data))
	for _, rrow := range rows.Data {
		data = append(data, map[string]any{
			"orderId":   rrow.OrderId,
			"orderDate": rrow.OrderDate.Format("2006-01-02 15:04:05"),
			"orderSource": func() string {
				if rrow.OrderSource != nil {
					return *rrow.OrderSource
				}
				return ""
			}(),
			"customerName":        rrow.CustomerName,
			"originalOrderAmount": rrow.OriginalAmount,
			"discountAmount":      rrow.DiscountAmount,
			"finalOrderAmount":    rrow.FinalAmount,
			"orderStatus":         rrow.OrderStatus,
		})
	}

	bytes, apiErr := h.excel.Export(ctx, data, columns, "Coupon Orders")
	if apiErr != nil {
		return apiErr
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=coupon-orders.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
	return nil
}

// Customer Profile Handlers
func (h *StatsHandler) GetCustomerProfile(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get customerId from query params
	customerId := r.URL.Query().Get("customerId")
	if customerId == "" {
		apiE := helpers.BadRequest("customerId query parameter is required")
		return &apiE
	}

	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)

	filterDto := stats_dto.CustomerProfileFilterDto{
		CustomerId: customerId,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	profile, err := h.orderSvc.GetCustomerProfile(ctx, filterDto)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, profile)
}

func (h *StatsHandler) GetCustomerOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}

	// Get customerId from query params
	customerId := r.URL.Query().Get("customerId")
	if customerId == "" {
		apiE := helpers.BadRequest("customerId query parameter is required")
		return &apiE
	}

	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)
	orderStatus := r.URL.Query().Get("orderStatus")
	branchId := r.URL.Query().Get("branchId")
	channel := r.URL.Query().Get("channel")
	paymentMethod := r.URL.Query().Get("paymentMethod")

	// Calculate page from skip and limit
	page := int(skip/int64(limit)) + 1

	filterDto := stats_dto.CustomerOrdersFilterDto{
		CustomerId:    customerId,
		StartDate:     &startDate,
		EndDate:       &endDate,
		OrderStatus:   &orderStatus,
		BranchId:      &branchId,
		Channel:       &channel,
		PaymentMethod: &paymentMethod,
		Page:          page,
		Limit:         int(limit),
	}

	result, err := h.orderSvc.GetCustomerOrders(ctx, filterDto)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *StatsHandler) GetCustomerCouponUsageDetails(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get customerId from query params
	customerId := r.URL.Query().Get("customerId")
	if customerId == "" {
		apiE := helpers.BadRequest("customerId query parameter is required")
		return &apiE
	}

	startDate := util.GetDateOrDefault(r.URL.Query().Get("startDate"), true)
	endDate := util.GetDateOrDefault(r.URL.Query().Get("endDate"), false)

	var startPtr, endPtr *time.Time

	if sd, err := time.Parse(time.RFC3339, startDate); err == nil {
		startPtr = &sd
	}

	if ed, err := time.Parse(time.RFC3339, endDate); err == nil {
		endPtr = &ed
	}

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return helpers.DefaultError(errGetPaginate)
	}
	couponUsage, err := h.orderSvc.GetCustomerCouponUsageDetails(ctx, customerId, startPtr, endPtr, skip, limit)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, couponUsage)
}

func (h *StatsHandler) GetCouponRedemptionsPerBranch(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Parse query parameters
	q := r.URL.Query()

	// Get startDate and endDate from query params, default to current month if not provided
	startDate := util.GetDateOrDefault(q.Get("startDate"), true)
	endDate := util.GetDateOrDefault(q.Get("endDate"), false)

	couponCode := q.Get("couponCode")
	paymentMethod := q.Get("paymentMethod")

	// Create filter DTO
	filterDto := stats_dto.CouponRedemptionsPerBranchFilterDto{
		StartDate:     startDate,
		EndDate:       endDate,
		CouponCode:    couponCode,
		PaymentMethod: paymentMethod,
	}

	// Call service method
	result, err := h.orderSvc.GetCouponRedemptionsPerBranch(ctx, filterDto)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *StatsHandler) GetCouponRedemptionsByCustomerType(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	q := r.URL.Query()
	startDate := util.GetDateOrDefault(q.Get("startDate"), true)
	endDate := util.GetDateOrDefault(q.Get("endDate"), false)
	couponCode := q.Get("couponCode")
	// paymentMethod := q.Get("paymentMethod")

	filter := stats_dto.CouponRedemptionsByCustomerTypeFilterDto{
		StartDate:  startDate,
		EndDate:    endDate,
		CouponCode: couponCode,
		// PaymentMethod: paymentMethod,
	}

	resp, err := h.orderSvc.GetCouponRedemptionsByCustomerType(ctx, filter)
	if err != nil {
		return err
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}

func (h *StatsHandler) GetCustomersCouponRedeemsBar(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	q := r.URL.Query()
	startDate := util.GetDateOrDefault(q.Get("startDate"), true)
	endDate := util.GetDateOrDefault(q.Get("endDate"), false)
	branchId := q.Get("branchId")

	// pagination
	skip, limit, err := util.Paginate(r)
	if err != nil {
		return helpers.DefaultError(err)
	}

	filter := stats_dto.CustomersCouponRedeemsBarFilterDto{
		StartDate: startDate,
		EndDate:   endDate,
		BranchId:  branchId,
		Skip:      skip,
		Limit:     limit,
	}
	resp, svcErr := h.orderSvc.GetCustomersCouponRedeemsBar(ctx, filter)
	if svcErr != nil {
		return svcErr
	}
	return helpers.WriteJson(w, http.StatusOK, resp)
}
