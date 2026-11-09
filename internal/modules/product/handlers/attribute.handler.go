package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttributeHandler struct {
	attributeService   services.AttributeService
	validationInstance *validator.Validate
}

func NewAttributeHandler(i *do.Injector, r *chi.Mux) {
	h := AttributeHandler{
		attributeService:   do.MustInvoke[services.AttributeService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.ATTRIBUTE.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: "authenticate",
				OwnerCapability:      "authenticate",
				MemberCapability:     "authenticate",
			}),
			middleware.ValidateDTO(dtos.CreateAttributeDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Attribute.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Attribute.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Attribute.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(product_filters.AttributeFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.General,
				OwnerCapability:      enums.CapabilitiesRegistry.General,
				MemberCapability:     enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(product_filters.AttributeFilters{}, middleware.QUERY, i)).Get("/all", h.GetAll)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Attribute.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Attribute.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Attribute.Member.Get.String(),
			}),
			middleware.ValidateDTO(product_filters.AttributeFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Attribute.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Attribute.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Attribute.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateAttributeDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Attribute.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Attribute.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Attribute.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Attribute.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Attribute.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.Attribute.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
	})
}

func (h *AttributeHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.attributeService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AttributeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.attributeService.GetAll(ctx, query)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AttributeHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.attributeService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *AttributeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateAttributeDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.attributeService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AttributeHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	data, err := middleware.GetValidatedDTO[dtos.CreateAttributeDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.attributeService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *AttributeHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.attributeService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AttributeHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.attributeService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
