package config_handlers

import (
	"fmt"
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

type CountryHandler struct {
	countryService     config_location_services.CountryService
	validationInstance *validator.Validate
}

func NewCountryHandler(i *do.Injector, r *chi.Mux) {
	h := CountryHandler{
		countryService:     do.MustInvoke[config_location_services.CountryService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.COUNTRY.GetAdminRoute(), func(r chi.Router) {
		r.Get("/seed/all", helpers.Make(h.GetAllSeeds))

		//Create Country
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Put)),
			middleware.ValidateDTO(dtos.CreateCountryDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Country By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.CountryFilters{}, middleware.QUERY, i),
		).Get("/{id}", helpers.Make(h.GetById))

		//Get All Countries
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Get)),
			middleware.ValidateDTO(config_filters.CountryFilters{}, middleware.QUERY, i),
		).Get("/all", helpers.Make(h.GetAll))

		//Get Countries With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Get)),
			middleware.ValidateDTO(config_filters.CountryFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetWithPagination))

		//Update Country By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Put)),
			middleware.ValidateDTO(dtos.CreateCountryDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		//Delete Country By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteById))

		// Get all countries with locations
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Country, enums.Get)),
		).Get("/all-with-locations", helpers.Make(h.GetAllWithLocations))
	})
}

func (h *CountryHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateCountryDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.countryService.Create(ctx, dto)
	if createErr != nil {
		return createErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CountryHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.countryService.GetById(ctx, objectId, query)
	if getErr != nil {
		return getErr
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CountryHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.countryService.GetAll(ctx, query)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CountryHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.countryService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

	return nil
}

func (h *CountryHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateCountryDto](r)
	if err != nil {
		return err
	}

	result, updateErr := h.countryService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		return updateErr
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CountryHandler) DeleteById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.countryService.DeleteById(ctx, objectId); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *CountryHandler) GetAllSeeds(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.countryService.GetAllSeedsCountries(ctx)
	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *CountryHandler) GetAllWithLocations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.countryService.GetAllWithLocations(ctx)
	if err != nil {
		return err
	}

	// Log the structure before marshalling for debugging
	fmt.Printf("Data before marshalling: %+v\n", result)

	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		fmt.Printf("Error marshalling response: %v\n", marshalErr)
		return helpers.DefaultError(marshalErr)
	}

	w.Write(response)
	return nil
}
