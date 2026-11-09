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

type FeatureHandler struct {
	featureService     feature_service.FeatureService
	validationInstance *validator.Validate
}

func NewFeatureHandler(i *do.Injector, r *chi.Mux) {
	h := FeatureHandler{
		featureService:     do.MustInvoke[feature_service.FeatureService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.FEATURE.GetAdminRoute(), func(r chi.Router) {
		//Create Feature
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.Post)),
			middleware.ValidateDTO(feature_dto.CreateFeatureDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Feature By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		//Get Features With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.Get)),
			middleware.ValidateDTO(feature_filters.FeatureFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		//Get All Features
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.Get)),
			middleware.ValidateDTO(feature_filters.FeatureFilters{}, middleware.QUERY, i),
		).Get("/all", h.GetAll)

		//Update Feature
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.Put)),
			middleware.ValidateDTO(feature_dto.UpdateFeatureDTO{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		//Delete Feature
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Feature, enums.Delete)),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.Delete)
	})
}

func (h *FeatureHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, err := middleware.GetValidatedDTO[feature_dto.CreateFeatureDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}
	result, createErr := h.featureService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *FeatureHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.featureService.GetById(ctx, objectId)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeatureHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.featureService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *FeatureHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[feature_dto.UpdateFeatureDTO](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.featureService.Update(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *FeatureHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	if deleteErr := h.featureService.Delete(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FeatureHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	query := r.URL.Query().Get("query")

	result, getErr := h.featureService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
