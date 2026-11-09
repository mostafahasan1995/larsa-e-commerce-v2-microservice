package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ZoneHandler struct {
	zoneService        config_services.ZoneService
	validationInstance *validator.Validate
}

func NewZoneHandler(i *do.Injector, r *chi.Mux) {
	h := ZoneHandler{
		zoneService:        do.MustInvoke[config_services.ZoneService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.ZONE.GetAdminRoute(), func(r chi.Router) {
		//Create Zone
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Put)),
			middleware.ValidateDTO(dtos.CreateZoneDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Zone By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.ZoneFilters{}, middleware.QUERY, i),
		).Get("/{id}", helpers.Make(h.GetById))

		//Get All Zones
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Get)),
			middleware.ValidateDTO(config_filters.ZoneFilters{}, middleware.QUERY, i),
		).Get("/all", helpers.Make(h.GetAll))

		//Get Zones With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Get)),
			middleware.ValidateDTO(config_filters.ZoneFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetWithPagination))

		//Update Zone By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Put)),
			middleware.ValidateDTO(dtos.CreateZoneDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		//Delete Zone By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteById))

		//Find Zones Containing Point
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Zone, enums.Get)),
			middleware.ValidateDTO(dtos.FindZonesContainingPointDto{}, middleware.BODY, i),
		).Post("/containing-point", helpers.Make(h.FindZonesContainingPoint))
	})
}

func (h *ZoneHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateZoneDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.zoneService.Create(ctx, dto)
	if createErr != nil {
		return createErr
	}

	response, marshalErr := json.MarshalContext(ctx, *result)
	if marshalErr != nil {
		return marshalErr
	}
	w.Write(response)

	return nil
}

func (h *ZoneHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.zoneService.GetById(ctx, objectId, query)
	if getErr != nil {
		return getErr
	}

	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		return marshalErr
	}
	w.Write(response)

	return nil
}

func (h *ZoneHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.zoneService.GetAll(ctx, query)
	if err != nil {
		return err
	}

	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		return helpers.DefaultError(marshalErr)
	}

	w.Write(response)

	return nil
}

func (h *ZoneHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.zoneService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		return marshalErr
	}
	w.Write(response)

	return nil
}

func (h *ZoneHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateZoneDto](r)
	if err != nil {
		return err
	}

	result, updateErr := h.zoneService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		return updateErr
	}

	response, marshalErr := json.MarshalContext(ctx, *result)
	if marshalErr != nil {
		return marshalErr
	}
	w.Write(response)

	return nil
}

func (h *ZoneHandler) DeleteById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.zoneService.DeleteById(ctx, objectId); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ZoneHandler) FindZonesContainingPoint(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.FindZonesContainingPointDto](r)
	if err != nil {
		return err
	}

	result, apiErr := h.zoneService.FindZonesContainingPoint(ctx, dto.Latitude, dto.Longitude)
	if apiErr != nil {
		return apiErr
	}

	response, marshalErr := json.MarshalContext(ctx, result)
	if marshalErr != nil {
		return marshalErr
	}
	w.Write(response)

	return nil
}
