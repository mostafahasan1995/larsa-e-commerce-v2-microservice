package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketTypeHandler struct {
	ticketTypeService  config_services.TicketTypeService
	validationInstance *validator.Validate
}

func NewTicketTypeHandler(i *do.Injector, r *chi.Mux) {
	h := TicketTypeHandler{
		ticketTypeService:  do.MustInvoke[config_services.TicketTypeService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.TICKET_TYPE.GetAdminRoute(), func(r chi.Router) {
		//Create Ticket Type
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.Put)),
			middleware.ValidateDTO(config_dtos.CreateTicketTypeDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		//Get Ticket Type By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(config_filters.TicketTypeFilter{}, middleware.QUERY, i),
		).Get("/{id}", helpers.Make(h.GetById))

		//Get All Ticket Types
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.Get)),
			middleware.ValidateDTO(config_filters.TicketTypeFilter{}, middleware.QUERY, i),
		).Get("/all", helpers.Make(h.GetAll))

		//Get Ticket Types With Pagination
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.Get)),
			middleware.ValidateDTO(config_filters.TicketTypeFilter{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetWithPagination))

		//Update Ticket Type By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.Put)),
			middleware.ValidateDTO(config_dtos.CreateTicketTypeDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		//Delete Ticket Type By ID
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.TicketType, enums.Delete)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteById))
	})
}

func (h *TicketTypeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[config_dtos.CreateTicketTypeDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.ticketTypeService.Create(ctx, dto)
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

func (h *TicketTypeHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.ticketTypeService.GetById(ctx, objectId, query)
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

func (h *TicketTypeHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.ticketTypeService.GetAll(ctx, query)
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

func (h *TicketTypeHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	result, getErr := h.ticketTypeService.GetWithPagination(ctx, skip, limit, query)
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

func (h *TicketTypeHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[config_dtos.CreateTicketTypeDto](r)
	if err != nil {
		return err
	}

	result, updateErr := h.ticketTypeService.UpdateById(ctx, objectId, data)
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

func (h *TicketTypeHandler) DeleteById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.ticketTypeService.DeleteById(ctx, objectId); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
