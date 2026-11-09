package call_center_handlers

import (
	"net/http"
	"strconv"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	call_center_services "larsa-e-commerce-microservice-v2/internal/modules/call-center/services"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type CallCenterHandler struct {
	callCenterService call_center_services.CallCenterService
}

func NewCallCenterHandler(i *do.Injector, r *chi.Mux) {
	h := CallCenterHandler{
		callCenterService: do.MustInvoke[call_center_services.CallCenterService](i),
	}

	// Admin routes with General capability
	r.Route("/admin/call-center", func(r chi.Router) {
		r.Get("/geocode", helpers.Make(h.Geocode))
		r.Get("/delivery-feasibility", helpers.Make(h.GetDeliveryFeasibility))
		r.Get("/branch-menu", helpers.Make(h.SearchBranchMenu))
		r.Get("/lookup", helpers.Make(h.Lookup))
	})
}

// Geocode converts an address to coordinates and location details
// GET /admin/call-center/geocode?address=<address>
func (h *CallCenterHandler) Geocode(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	address := r.URL.Query().Get("address")
	if address == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "address query parameter is required",
		}
	}

	result, apiErr := h.callCenterService.Geocode(ctx, address)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

// GetDeliveryFeasibility checks delivery feasibility for given coordinates
// GET /admin/call-center/delivery-feasibility?lat=<lat>&lng=<lng>&currencyCode=<code>
func (h *CallCenterHandler) GetDeliveryFeasibility(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "lat and lng query parameters are required",
		}
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid latitude value",
		}
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid longitude value",
		}
	}

	var currencyCode *string
	if code := r.URL.Query().Get("currencyCode"); code != "" {
		currencyCode = &code
	}

	result, apiErr := h.callCenterService.GetDeliveryFeasibility(ctx, lat, lng, currencyCode)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

// SearchBranchMenu searches products at a specific branch
// GET /admin/call-center/branch-menu?branchId=<id>&currencyCode=<code>&q=<query>&skip=<skip>&limit=<limit>&sorts=<sorts>
func (h *CallCenterHandler) SearchBranchMenu(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	branchId := r.URL.Query().Get("branchId")
	if branchId == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "branchId query parameter is required",
		}
	}

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	var query *string
	if q := r.URL.Query().Get("q"); q != "" {
		query = &q
	}

	var sorts *string
	if s := r.URL.Query().Get("sorts"); s != "" {
		sorts = &s
	}

	var currencyCode *string
	if code := r.URL.Query().Get("currencyCode"); code != "" {
		currencyCode = &code
	}

	result, apiErr := h.callCenterService.SearchBranchMenu(ctx, branchId, query, skip, limit, sorts, currencyCode)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

// Lookup performs combined geocoding, delivery feasibility, and optional menu search
// GET /admin/call-center/lookup?address=<address>&currencyCode=<code>&menuQuery=<query>&skip=<skip>&limit=<limit>&sorts=<sorts>
func (h *CallCenterHandler) Lookup(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	address := r.URL.Query().Get("address")
	if address == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "address query parameter is required",
		}
	}

	var currencyCode *string
	if code := r.URL.Query().Get("currencyCode"); code != "" {
		currencyCode = &code
	}

	var menuQuery *string
	if q := r.URL.Query().Get("menuQuery"); q != "" {
		menuQuery = &q
	}

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	var sorts *string
	if s := r.URL.Query().Get("sorts"); s != "" {
		sorts = &s
	}

	result, apiErr := h.callCenterService.Lookup(ctx, address, currencyCode, menuQuery, skip, limit, sorts)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}
