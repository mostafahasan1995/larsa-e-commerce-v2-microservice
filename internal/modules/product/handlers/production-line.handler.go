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

	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type ProductionLineHandler struct {
	productionLineService services.ProductionLineService
	validationInstance    *validator.Validate
}

func NewProductionLineHandler(i *do.Injector, r *chi.Mux) {
	h := ProductionLineHandler{
		productionLineService: do.MustInvoke[services.ProductionLineService](i),
		validationInstance:    do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.PRODUCTION_LINE.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.Post.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateProductionLineDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.GetOne.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(product_filters.ProductionLineFilters{}, middleware.QUERY, i)).Get("/{id}", h.GetById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.Get.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.Get.String(),
			}), middleware.ValidateDTO(product_filters.ProductionLineFilters{}, middleware.QUERY, i)).Get("/", h.GetWithPagination)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.Put.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.Put.String(),
			}),
			middleware.ValidateDTO(dtos.CreateProductionLineDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.Delete.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.ProductionLine.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.ProductionLine.Vendor.DeleteMany.String(),
				MemberCapability:     enums.CapabilitiesRegistry.ProductionLine.Member.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
	})
}

func (h *ProductionLineHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateProductionLineDto](r)

	if err != nil {
		helpers.DefaultError(err)
		return nil
	}

	result, createErr := h.productionLineService.Create(ctx, dto)
	if createErr != nil {

		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *ProductionLineHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.productionLineService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

}

func (h *ProductionLineHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.productionLineService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)

}

func (h *ProductionLineHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.CreateProductionLineDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return nil
	}

	result, updateErr := h.productionLineService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *ProductionLineHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.productionLineService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductionLineHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return nil
	}

	if bulkDeleteErr := h.productionLineService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
