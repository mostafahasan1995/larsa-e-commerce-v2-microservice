package quotation_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_filters "larsa-e-commerce-microservice-v2/internal/modules/quotation/filters"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

func (h *QuotationAdminHandler) InitializeVendorRoutes(r chi.Router, i *do.Injector) {
	r.With(
		middleware.AdminAuth(middleware.Restriction{
			SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			VendorCapability:     enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
		}),
		middleware.ValidateDTO(quotation_dto.CreateQuotationDTO{}, middleware.BODY, i),
	).Post("/{requestQuotationId}/response", helpers.Make(h.CreateResponseQuotation))

	r.With(
		middleware.AdminAuth(middleware.Restriction{
			SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			VendorCapability:     enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
		}),
		middleware.ValidateDTO(quotation_dto.CreateQuotationDTO{}, middleware.BODY, i),
	).Post("/general", helpers.Make(h.CreateGeneralQuotation))

	r.With(
		middleware.AdminAuth(middleware.Restriction{
			SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			VendorCapability:     enums.CapabilitiesRegistry.Quotation.Vendor.Get.String(),
		}),
		middleware.ValidateDTO(quotation_filters.QuotationFilters{}, middleware.QUERY, i),
	).Get("/request", helpers.Make(h.GetRequestQuotations))

}

func (h *QuotationAdminHandler) CreateGeneralQuotation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	dto, _ := middleware.GetValidatedDTO[quotation_dto.CreateQuotationDTO](r)

	result, apiErr := h.quotationVendorService.CreateGeneralQuotation(ctx, dto)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) CreateResponseQuotation(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	requestQuotationId := chi.URLParam(r, "requestQuotationId")
	if requestQuotationId == "" {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Request quotation ID is required",
		}
	}

	dto, dtoError := middleware.GetValidatedDTO[quotation_dto.CreateQuotationDTO](r)
	if dtoError != nil {
		return dtoError
	}

	result, apiErr := h.quotationVendorService.CreateResponseQuotation(ctx, requestQuotationId, dto)
	if apiErr != nil {
		return apiErr
	}

	response, marshalError := json.MarshalContext(ctx, result)
	if marshalError != nil {
		return marshalError
	}

	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) GetRequestQuotations(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		return errGetPaginate
	}
	result, getErr := h.quotationVendorService.GetRequestQuotations(ctx, skip, limit, query)
	if getErr != nil {
		return getErr
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
	return nil
}
