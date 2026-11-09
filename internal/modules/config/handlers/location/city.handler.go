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

type CityHandler struct {
	cityService        config_location_services.CityService
	validationInstance *validator.Validate
}

func NewCityHandler(i *do.Injector, r *chi.Mux) {
	h := CityHandler{
		cityService:        do.MustInvoke[config_location_services.CityService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.CITY.GetAdminRoute(), func(r chi.Router) {
		//Get All Seed Cities
		r.Get("/seeds", helpers.Make(h.GetAllSeeds))

		//Create City
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.Put)),
			middleware.ValidateDTO(dtos.CreateCityDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get City By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.CityFilters{}, middleware.QUERY, i),
		).Get("/{id}", helpers.Make(h.GetById))

		//Get All Cities
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.Get)),
			middleware.ValidateDTO(config_filters.CityFilters{}, middleware.QUERY, i),
		).Get("/all", helpers.Make(h.GetAll))

		//Get Cities With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.Get)),
			middleware.ValidateDTO(config_filters.CityFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetWithPagination))

		//Update City By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.Put)),
			middleware.ValidateDTO(dtos.CreateCityDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		//Delete City By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.City, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteById))

	})
}

func (h *CityHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateCityDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.cityService.Create(ctx, dto)
	if createErr != nil {
		return createErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CityHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.cityService.GetById(ctx, objectId, query)
	if getErr != nil {
		return getErr
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CityHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.cityService.GetAll(ctx, query)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CityHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.cityService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CityHandler) GetAllSeeds(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	countryName := r.URL.Query().Get("countryName")
	if countryName == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "countryName is required",
		}
	}

	result, err := h.cityService.GetAllSeedsCities(ctx, countryName)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CityHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateCityDto](r)
	if err != nil {
		return err
	}

	result, updateErr := h.cityService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		return updateErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CityHandler) DeleteById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.cityService.DeleteById(ctx, objectId); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
