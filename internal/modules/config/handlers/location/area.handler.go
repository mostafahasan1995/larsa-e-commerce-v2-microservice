package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_services "larsa-e-commerce-microservice-v2/internal/modules/config/services/location"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AreaHandler struct {
	areaService        config_location_services.AreaService
	validationInstance *validator.Validate
}

func NewAreaHandler(i *do.Injector, r *chi.Mux) {
	h := AreaHandler{
		areaService:        do.MustInvoke[config_location_services.AreaService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.AREA.GetAdminRoute(), func(r chi.Router) {
		//Get All Seed Areas
		r.Get("/seeds", helpers.Make(h.GetAllSeeds))

		//Create  Areas
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.Post)),
			middleware.ValidateDTO(dtos.CreateAreaDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Area By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.AreaFilters{}, middleware.QUERY, i),
		).Get("/{id}", helpers.Make(h.GetById))

		//Get All Areas
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.Get)),
			middleware.ValidateDTO(config_filters.AreaFilters{}, middleware.QUERY, i),
		).Get("/all", helpers.Make(h.GetAll))

		//Get Areas With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.Get)),
			middleware.ValidateDTO(config_filters.AreaFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetWithPagination))

		//Update Area By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.Put)),
			middleware.ValidateDTO(dtos.CreateAreaDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		//Delete Area By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Area, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteById))

	})
}

func (h *AreaHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateAreaDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.areaService.Create(ctx, dto)
	if createErr != nil {
		return createErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.areaService.GetById(ctx, objectId, query)
	if getErr != nil {
		return getErr
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.areaService.GetAll(ctx, query)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.areaService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) GetAllSeeds(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	cityName := r.URL.Query().Get("cityName")
	if cityName == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "cityName is required",
		}
	}

	result, err := h.areaService.GetAllSeedsAreas(ctx, cityName)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateAreaDto](r)
	if err != nil {
		return err
	}

	result, updateErr := h.areaService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		return updateErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AreaHandler) DeleteById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.areaService.DeleteById(ctx, objectId); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
