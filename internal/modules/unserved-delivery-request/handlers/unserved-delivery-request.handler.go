package unserved_delivery_request_handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	unserved_delivery_request_filters "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/filters"
	unserved_delivery_request_services "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/services"
)

type UnservedDeliveryRequestHandler struct {
	service unserved_delivery_request_services.UnservedDeliveryRequestService
	excel   excel_services.ExcelService
}

func NewUnservedDeliveryRequestHandler(i *do.Injector, r *chi.Mux) {
	h := UnservedDeliveryRequestHandler{
		service: do.MustInvoke[unserved_delivery_request_services.UnservedDeliveryRequestService](i),
		excel:   do.MustInvoke[excel_services.ExcelService](i),
	}

	r.Route(enums.UNSERVED_DELIVERY_REQUEST.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
		).Get("/", helpers.Make(h.GetWithPagination))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
		).Get("/aggregate", helpers.Make(h.GetAggregated))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
		).Get("/export", helpers.Make(h.Export))
	})
}

func (h *UnservedDeliveryRequestHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, err := util.Paginate(r)
	if err != nil {
		return err
	}

	q := r.URL.Query()
	filter := h.buildFilterFromQuery(q)

	result, apiErr := h.service.GetWithPagination(ctx, skip, limit, filter)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return nil
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *UnservedDeliveryRequestHandler) GetAggregated(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	q := r.URL.Query()
	groupBy := q.Get("groupBy")
	if groupBy == "" {
		groupBy = "date"
	}

	filter := h.buildFilterFromQuery(q)

	result, apiErr := h.service.GetAggregated(ctx, groupBy, filter)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return nil
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *UnservedDeliveryRequestHandler) Export(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := util.AddCtxAppCfg(r)

	q := r.URL.Query()
	filter := h.buildFilterFromQuery(q)

	result, apiErr := h.service.GetAll(ctx, filter)
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return nil
	}

	// Columns match list response data without id fields (no id, no userId)
	columns := []excel_services.ExcelColumn{
		{Key: "latitude", Header: "Latitude", Type: excel_services.ExcelTypeFloat},
		{Key: "longitude", Header: "Longitude", Type: excel_services.ExcelTypeFloat},
		{Key: "country", Header: "Country", Type: excel_services.ExcelTypeString},
		{Key: "city", Header: "City", Type: excel_services.ExcelTypeString},
		{Key: "area", Header: "Area", Type: excel_services.ExcelTypeString},
		{Key: "source", Header: "Source", Type: excel_services.ExcelTypeString},
		{Key: "createdAt", Header: "Created At", Type: excel_services.ExcelTypeDateTime, DateFormat: "2006-01-02 15:04:05"},
	}

	data := make([]map[string]any, 0, len(result))
	for _, row := range result {
		country := ""
		if row.Country != nil {
			country = *row.Country
		}
		city := ""
		if row.City != nil {
			city = *row.City
		}
		area := ""
		if row.Area != nil {
			area = *row.Area
		}
		source := ""
		if row.Source != nil {
			source = *row.Source
		}
		createdAt := ""
		if !row.CreatedAt.IsZero() {
			createdAt = row.CreatedAt.Format("2006-01-02 15:04:05")
		}
		data = append(data, map[string]any{
			"latitude":  row.Latitude,
			"longitude": row.Longitude,
			"country":   country,
			"city":      city,
			"area":      area,
			"source":    source,
			"createdAt": createdAt,
		})
	}

	bytes, apiErr := h.excel.Export(ctx, data, columns, "Unserved Delivery Requests")
	if apiErr != nil {
		helpers.HandleApiError(w, apiErr)
		return nil
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=unserved-delivery-requests.xlsx")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
	return nil
}

func (h *UnservedDeliveryRequestHandler) buildFilterFromQuery(q interface{ Get(string) string }) *unserved_delivery_request_filters.UnservedDeliveryRequestFilters {
	filter := &unserved_delivery_request_filters.UnservedDeliveryRequestFilters{}

	if v := q.Get("country"); v != "" {
		filter.Country = &v
	}
	if v := q.Get("city"); v != "" {
		filter.City = &v
	}
	if v := q.Get("area"); v != "" {
		filter.Area = &v
	}

	startDate := util.GetDateOrDefault(q.Get("startDate"), true)
	endDate := util.GetDateOrDefault(q.Get("endDate"), false)
	var gte, lte *time.Time
	if startDate != "" {
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			gte = &t
		}
	}
	if endDate != "" {
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			lte = &t
		}
	}
	if gte != nil || lte != nil {
		filter.CreatedAt = &filters.DateRange{Gte: gte, Lte: lte}
	}

	return filter
}
