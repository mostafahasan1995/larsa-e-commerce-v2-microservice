package plan_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	plan_filters "larsa-e-commerce-microservice-v2/internal/modules/plan/filters"
	plan_service "larsa-e-commerce-microservice-v2/internal/modules/plan/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebPlanHandler struct {
	webPlanService     plan_service.WebPlanService
	validationInstance *validator.Validate
}

func NewWebPlanHandler(i *do.Injector, r *chi.Mux) {
	h := WebPlanHandler{
		webPlanService:     do.MustInvoke[plan_service.WebPlanService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.PLAN.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(plan_filters.PlanFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.ValidateDTO(plan_filters.PlanFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.ValidateDTO(plan_filters.PlanFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
	})
}

func (h *WebPlanHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.webPlanService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebPlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.webPlanService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebPlanHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webPlanService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
