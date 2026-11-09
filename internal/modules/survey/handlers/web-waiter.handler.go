package survey_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebWaiterHandler struct {
	waiterService      survey_services.WebWaiterService
	validationInstance *validator.Validate
}

func NewWebWaiterHandler(i *do.Injector, r *chi.Mux) {
	h := WebWaiterHandler{
		waiterService:      do.MustInvoke[survey_services.WebWaiterService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.WAITER.GetWebRoute(), func(r chi.Router) {
		// Get waiter by ID
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		// Get waiters with pagination
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(survey_filters.WaiterFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		// Get all waiters
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(survey_filters.WaiterFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)
	})
}

func (h *WebWaiterHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.waiterService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebWaiterHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.waiterService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebWaiterHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	result, getErr := h.waiterService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
