package contact_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/agent"
	agent_filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters/agent"
	services "larsa-e-commerce-microservice-v2/internal/modules/contact/services/agent"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentHandler struct {
	agentService       services.AgentService
	validationInstance *validator.Validate
}

func NewAgentHandler(i *do.Injector, r *chi.Mux) {
	h := AgentHandler{
		agentService:       do.MustInvoke[services.AgentService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.AGENT.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.BrokerAgent.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.BrokerAgent.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.BrokerAgent.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateAgentDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.BrokerAgent.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.BrokerAgent.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.BrokerAgent.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(agent_filters.AgentFilters{}, middleware.QUERY, i),
		).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(agent_filters.AgentFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.BrokerAgent.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.BrokerAgent.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.BrokerAgent.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateAgentDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.BrokerAgent.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.BrokerAgent.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.BrokerAgent.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.BrokerAgent.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.BrokerAgent.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.BrokerAgent.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
	})
}

func (h *AgentHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateAgentDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.agentService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AgentHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	result, getErr := h.agentService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AgentHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.agentService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AgentHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	data, err := middleware.GetValidatedDTO[dtos.CreateAgentDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.agentService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AgentHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.agentService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AgentHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.agentService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
