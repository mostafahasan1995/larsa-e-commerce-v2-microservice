package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type BranchWebHandler struct {
	branchWebService   services.BranchWebService
	validationInstance *validator.Validate
}

func NewBranchWebHandler(i *do.Injector, r *chi.Mux) {
	h := BranchWebHandler{
		branchWebService:   do.MustInvoke[services.BranchWebService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.BRANCH.GetWebRoute(), func(r chi.Router) {
		r.Get("/by-location", helpers.Make(h.GetByLatLng))
		r.Get("/", helpers.Make(h.GetWithPagination))
		r.Get("/all", helpers.Make(h.GetAll))
		r.Get("/{id}", helpers.Make(h.GetById))
	})
}

func (h *BranchWebHandler) GetByLatLng(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Get lat and lng from query parameters
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		return &helpers.ApiError{StatusCode: http.StatusBadRequest, Message: "Latitude and longitude are required"}
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return &helpers.ApiError{StatusCode: http.StatusBadRequest, Message: "Invalid latitude value"}
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return &helpers.ApiError{StatusCode: http.StatusBadRequest, Message: "Invalid longitude value"}
	}
    var country, city, area *string
	if v := r.URL.Query().Get("country"); v != "" {
		country = &v
	}
	if v := r.URL.Query().Get("city"); v != "" {
		city = &v
	}
	if v := r.URL.Query().Get("area"); v != "" {
		area = &v
	}

	result, apiErr := h.branchWebService.GetByLatLng(ctx, lat, lng,country, city, area)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *BranchWebHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	query := r.URL.Query().Get("query")
	result, apiErr := h.branchWebService.GetWithPagination(ctx, skip, limit, query)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *BranchWebHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, apiErr := h.branchWebService.GetAll(ctx, query)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *BranchWebHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		return &helpers.ApiError{StatusCode: http.StatusBadRequest, Message: "Branch Id is required"}
	}

	result, apiErr := h.branchWebService.GetById(ctx, id)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}
