package order_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *OrderHandler) NewOrderSummaryHandler(i *do.Injector, r chi.Router) {
	r.With(
		middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Get)),
	).Get("/upcoming-shipping", helpers.Make(h.GetUpcomingShipping))

	r.With(
		middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Get)),
	).Get("/recent-activity", helpers.Make(h.GetRecentActivity))

	r.With(
		middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Get)),
	).Get("/buyer/recent-activity", helpers.Make(h.GetBuyerActivity))

	r.With(
		middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Get)),
	).Get("/export", helpers.Make(h.ExportOrders))

}

func (h *OrderHandler) GetUpcomingShipping(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	limit := int64(3)

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.ParseInt(limitParam, 10, 64)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	upcomingOrders, apiErr := h.orderSummaryService.GetUpcomingShipments(ctx, limit)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, upcomingOrders)
}

func (h *OrderHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	branchIdParam := r.URL.Query().Get("branchId")
	branchId := primitive.NilObjectID
	if branchIdParam != "" {
		id, err := primitive.ObjectIDFromHex(branchIdParam)
		if err != nil {
			apiErr := helpers.BadRequest("Invalid branch ID")
			return &apiErr
		}
		branchId = id
	}

	query := r.URL.Query().Get("query")

	activityStats, apiErr := h.orderSummaryService.GetRecentActivity(ctx, &branchId, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, activityStats)
}

func (h *OrderHandler) GetBuyerActivity(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	branchIdParam := r.URL.Query().Get("branchId")
	branchId := primitive.NilObjectID
	if branchIdParam != "" {
		id, err := primitive.ObjectIDFromHex(branchIdParam)
		if err != nil {
			apiErr := helpers.BadRequest("Invalid branch ID")
			return &apiErr
		}
		branchId = id
	}

	query := r.URL.Query().Get("query")

	activityStats, apiErr := h.orderSummaryService.GetBuyerActivity(ctx, &branchId, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, activityStats)
}

func (h *OrderHandler) ExportOrders(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	limit := int64(0)
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		parsedLimit, err := strconv.ParseInt(limitParam, 10, 64)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	excelData, apiErr := h.orderSummaryService.ExportOrders(ctx, query, limit)
	if apiErr != nil {
		return apiErr
	}

	// Set headers for Excel file download
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=orders.xlsx")

	// Write Excel data directly to response
	_, err := w.Write(excelData)
	return err
}
