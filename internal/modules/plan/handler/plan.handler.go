package plan_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	plan_dto "larsa-e-commerce-microservice-v2/internal/modules/plan/dto"
	plan_filters "larsa-e-commerce-microservice-v2/internal/modules/plan/filters"
	plan_service "larsa-e-commerce-microservice-v2/internal/modules/plan/service"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanHandler struct {
	planService        plan_service.PlanService
	validationInstance *validator.Validate
}

func NewPlanHandler(i *do.Injector, r *chi.Mux) {
	h := PlanHandler{
		planService:        do.MustInvoke[plan_service.PlanService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.PLAN.GetAdminRoute(), func(r chi.Router) {
		//Create Plan
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Post)),
			middleware.ValidateDTO(plan_dto.CreatePlanDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Plan By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Get)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		//Get Plans With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Get)),
			middleware.ValidateDTO(plan_filters.PlanFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Plans
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Get)),
			middleware.ValidateDTO(plan_filters.PlanFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Update Plan
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Put)),
			middleware.ValidateDTO(plan_dto.UpdatePlanDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		//Delete Plan
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Plan, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}
func (h *PlanHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[plan_dto.CreatePlanDTO](r)

	result, createErr := h.planService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *PlanHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.planService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *PlanHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.planService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *PlanHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[plan_dto.UpdatePlanDTO](r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, updateErr := h.planService.Update(ctx, objectId, dto)

	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}
	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *PlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.planService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, getErr := h.planService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
