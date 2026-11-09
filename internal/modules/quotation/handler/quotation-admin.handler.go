package quotation_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_service "larsa-e-commerce-microservice-v2/internal/modules/quotation/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type QuotationAdminHandler struct {
	quotationAdminService  quotation_service.QuotationAdminService
	validationInstance     *validator.Validate
	quotationVendorService quotation_service.QuotationVendorService
}

func NewQuotationAdminHandler(i *do.Injector, r *chi.Mux) {
	h := QuotationAdminHandler{
		quotationAdminService:  do.MustInvoke[quotation_service.QuotationAdminService](i),
		validationInstance:     do.MustInvoke[*validator.Validate](i),
		quotationVendorService: do.MustInvoke[quotation_service.QuotationVendorService](i),
	}

	r.Route(enums.QUOTATION.GetAdminRoute(), func(r chi.Router) {
		// ----------------- //
		// Admin & Vendor Shared Endpoints
		// ----------------- //
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
				VendorCapability:     enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
			}),
		).Post("/{id}/status", helpers.Make(h.UpdateQuotationStatus))

		// this endpoint will used by admin, vendor, agent and shipping company
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability:      enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
				VendorCapability:          enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
				AgentCapability:           enums.CapabilitiesRegistry.Quotation.Agent.Get.String(),
				ShippingCompanyCapability: enums.CapabilitiesRegistry.Quotation.ShippingCompany.Get.String(),
			}),
			middleware.ValidateDTO(quotation_filters.QuotationFilters{}, middleware.QUERY, i),
		).Get("/", helpers.Make(h.GetQuotations))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability:      enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
				VendorCapability:          enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
				AgentCapability:           enums.CapabilitiesRegistry.Quotation.Agent.Get.String(),
				ShippingCompanyCapability: enums.CapabilitiesRegistry.Quotation.ShippingCompany.Get.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", helpers.Make(h.GetQuotationById))

		// ----------------- //
		// Admin Endpoints
		// ----------------- //
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			}),
			middleware.ValidateDTO(quotation_dto.CostDto{}, middleware.BODY, i),
		).Post("/{id}/admin-fee", helpers.Make(h.SetSuperAdminFee))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Delete.String(),
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Delete("/{id}", helpers.Make(h.DeleteQuotationById))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.DeleteMany.String(),
			}),
		).Delete("/", helpers.Make(h.DeleteAllQuotations))

		// ----------------- //
		// Vendor Endpoints
		// ----------------- //
		h.InitializeVendorRoutes(r, i)

		// ----------------- //
		// Agent & Shipping Company Endpoints
		// ----------------- //
		h.InitializeAgentShippingRoutes(r, i)
	})
}

func (h *QuotationAdminHandler) UpdateQuotationStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	status := r.URL.Query().Get("status")

	if status == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Status is required",
		}
	}

	statusEnum := quotation_enums.QuotationStatus(status)

	if updateErr := h.quotationAdminService.UpdateQuotationStatus(ctx, id, statusEnum); updateErr != nil {
		return updateErr
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *QuotationAdminHandler) GetQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}

	cfg, _ := util.GetReqAppCfg(ctx)
	result := &quotation_models.QuotationsWithPagination{}
	var err *helpers.ApiError
	if cfg.User.Type == enums.Vendor {
		result, err = h.quotationVendorService.GetWithPagination(ctx, skip, limit, query)
	} else {
		// For SuperAdmin, Agent, ShippingCompany, etc., use admin service with auto-filtering
		result, err = h.quotationAdminService.GetWithPagination(ctx, skip, limit, query)
	}

	if err != nil {
		return err
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) GetQuotationById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	query := r.URL.Query().Get("query")

	result, getErr := h.quotationAdminService.GetById(ctx, id, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) SetSuperAdminFee(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")
	dto, err := middleware.GetValidatedDTO[quotation_dto.CostDto](r)
	if err != nil {
		return err
	}
	updated, apiErr := h.quotationAdminService.SetAdminFee(ctx, id, *dto)
	if apiErr != nil {
		return apiErr
	}
	response, _ := json.MarshalContext(ctx, updated)
	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) DeleteQuotationById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	if deleteErr := h.quotationAdminService.DeleteById(ctx, id); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *QuotationAdminHandler) DeleteAllQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	if deleteErr := h.quotationAdminService.DeleteAll(ctx); deleteErr != nil {
		return deleteErr
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
