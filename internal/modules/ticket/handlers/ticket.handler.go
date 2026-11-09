package ticket_handlers

import (
	"encoding/json"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	ticket_dto "larsa-e-commerce-microservice-v2/internal/modules/ticket/dto"
	ticket_service "larsa-e-commerce-microservice-v2/internal/modules/ticket/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketHandler struct {
	service            ticket_service.TicketService
	validationInstance interface{}
}

func NewTicketHandler(i *do.Injector, r *chi.Mux) {
	h := TicketHandler{
		service: do.MustInvoke[ticket_service.TicketService](i),
	}

	r.Route(enums.TICKET.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Ticket, enums.Post)),
			middleware.ValidateDTO(ticket_dto.CreateTicketDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Ticket, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", helpers.Make(h.GetById))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Ticket, enums.Put)),
			middleware.ValidateDTO(ticket_dto.CreateTicketDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.Update))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Ticket, enums.Put)),
			middleware.ValidateDTO(ticket_dto.AddCommentDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/comment", helpers.Make(h.AddComment))

		r.With(
			middleware.Auth([]string{}...),
		).Get("/", helpers.Make(h.GetMyTickets))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Ticket, enums.GetOne)),
		).Get("/summary", helpers.Make(h.GetSummary))

	})

}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var dto ticket_dto.CreateTicketDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	ticket, apiErr := h.service.Create(ctx, dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid ticket ID")
		return &apiErr
	}

	ticket, apiErr := h.service.GetById(ctx, id)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, ticket)
}

func (h *TicketHandler) AddComment(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid ticket ID")
		return &apiErr
	}

	var dto ticket_dto.AddCommentDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	ticket, apiErr := h.service.AddComment(ctx, id, dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, ticket)
}

func (h *TicketHandler) GetMyTickets(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		apiErr := helpers.DefaultError(errGetPaginate)
		return apiErr
	}

	result, apiErr := h.service.GetMyTickets(ctx, skip, limit, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid ticket ID")
		return &apiErr
	}

	var dto ticket_dto.CreateTicketDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	ticket, apiErr := h.service.Update(ctx, id, dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, ticket)
}

func (h *TicketHandler) GetSummary(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	summary, apiErr := h.service.GetTicketSummary(ctx)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, summary)
}
