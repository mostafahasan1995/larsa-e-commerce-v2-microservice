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

// AddonHandler handles HTTP requests for the addon module
type AddonHandler struct {
	addonService       services.AddonService
	validationInstance *validator.Validate
}

// NewAddonHandler creates a new AddonHandler and sets up the routes
func NewAddonHandler(i *do.Injector, r *chi.Mux) {
	h := AddonHandler{
		addonService:       do.MustInvoke[services.AddonService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.ADDON.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.Post.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.Post.String(),
			}),
			middleware.ValidateDTO(dtos.CreateAddonDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.GetOne.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.GetOne.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", h.GetById)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.Get.String(),
			}),
			middleware.ValidateDTO(product_filters.AddonFilters{}, middleware.QUERY, i),
		).Get("/", h.GetWithPagination)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.Get.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.Get.String(),
			}),
		).Get("/all", h.GetAll)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.Put.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.Put.String(),
			}),
			middleware.ValidateDTO(dtos.UpdateAddonDto{}, middleware.BODY, i),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}", helpers.Make(h.UpdateById))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.Delete.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", h.DeleteById)

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Addon.SuperAdmin.DeleteMany.String(),
				OwnerCapability:      enums.CapabilitiesRegistry.Addon.Vendor.DeleteMany.String(),
			}),
			middleware.ValidateDTO(types.BulkDeleteDto{}, middleware.BODY, i),
		).Delete("/", helpers.Make(h.BulkDelete))
	})
}

// Create handles the creation of a new addon
func (h *AddonHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.CreateAddonDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, createErr := h.addonService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

// GetById handles the retrieval of an addon by ID
func (h *AddonHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	result, getErr := h.addonService.GetById(ctx, objectId, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetWithPagination handles the retrieval of paginated addons
func (h *AddonHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	ctx, _ := util.AddCtxAppCfg(r)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.addonService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetAll handles the retrieval of all addons without pagination
func (h *AddonHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	result, getErr := h.addonService.GetAll(ctx, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// UpdateById handles the update of an addon by ID
func (h *AddonHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	data, err := middleware.GetValidatedDTO[dtos.UpdateAddonDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	result, updateErr := h.addonService.UpdateById(ctx, objectId, data)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

// DeleteById handles the deletion of an addon by ID
func (h *AddonHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	if deleteErr := h.addonService.DeleteById(ctx, objectId); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BulkDelete handles the deletion of multiple addons
func (h *AddonHandler) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	data, err := middleware.GetValidatedDTO[types.BulkDeleteDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	if bulkDeleteErr := h.addonService.BulkDelete(ctx, data.Ids); bulkDeleteErr != nil {
		helpers.HandleApiError(w, bulkDeleteErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
