package feature_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	feature_dto "larsa-e-commerce-microservice-v2/internal/modules/feature/dto"
	feature_filters "larsa-e-commerce-microservice-v2/internal/modules/feature/filters"
	feature_service "larsa-e-commerce-microservice-v2/internal/modules/feature/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeatureGroupHandler struct {
	featureGroupService feature_service.FeatureGroupService
	validationInstance  *validator.Validate
}

func NewFeatureGroupHandler(i *do.Injector, r *chi.Mux) {
	h := FeatureGroupHandler{
		featureGroupService: do.MustInvoke[feature_service.FeatureGroupService](i),
		validationInstance:  do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.FEATURE_GROUP.GetAdminRoute(), func(r chi.Router) {
		//Create Feature Group
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.Post)),
			middleware.ValidateDTO(feature_dto.CreateFeatureGroupDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Feature Group By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		//Get Feature Groups With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.Get)),
			middleware.ValidateDTO(feature_filters.FeatureGroupFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Feature Groups
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.Get)),
			middleware.ValidateDTO(feature_filters.FeatureGroupFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Update Feature Group
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.Put)),
			middleware.ValidateDTO(feature_dto.UpdateFeatureGroupDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		//Delete Feature Group
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.FeatureGroup, enums.Delete)),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}

func (h *FeatureGroupHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[feature_dto.CreateFeatureGroupDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.featureGroupService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *FeatureGroupHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.featureGroupService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeatureGroupHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	query := r.URL.Query().Get("query")

	result, getErr := h.featureGroupService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeatureGroupHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[feature_dto.UpdateFeatureGroupDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.featureGroupService.Update(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *FeatureGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.featureGroupService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FeatureGroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, getErr := h.featureGroupService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
