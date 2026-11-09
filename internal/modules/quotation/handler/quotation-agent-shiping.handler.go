package quotation_handler

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

func (h *QuotationAdminHandler) InitializeAgentShippingRoutes(r chi.Router, i *do.Injector) {
	r.With(
		middleware.AdminAuth(middleware.Restriction{
			SuperAdminCapability:      enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			AgentCapability:           enums.CapabilitiesRegistry.Quotation.Agent.Get.String(),
			ShippingCompanyCapability: enums.CapabilitiesRegistry.Quotation.ShippingCompany.Get.String(),
		}),
		middleware.ValidateDTO(quotation_dto.SetShippingCostDTO{}, middleware.BODY, i),
	).Post("/{id}/shipping-cost", helpers.Make(h.SetShippingCost))

	// Add agent fee by agent or shipping company
	r.With(
		middleware.AdminAuth(middleware.Restriction{
			SuperAdminCapability: enums.CapabilitiesRegistry.Quotation.SuperAdmin.Get.String(),
			AgentCapability:      enums.CapabilitiesRegistry.Quotation.Agent.Get.String(),
		}),
		middleware.ValidateDTO(quotation_dto.SetAgentFeeDTO{}, middleware.BODY, i),
	).Post("/{id}/agent-fee", helpers.Make(h.SetAgentFee))
}

func (h *QuotationAdminHandler) SetShippingCost(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	dto, dtoErr := middleware.GetValidatedDTO[quotation_dto.SetShippingCostDTO](r)
	if dtoErr != nil {
		return dtoErr
	}

	updated, apiErr := h.quotationAdminService.SetShippingCost(ctx, id, *dto)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, updated)
	w.Write(response)
	return nil
}

func (h *QuotationAdminHandler) SetAgentFee(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	id := chi.URLParam(r, "id")

	dto, dtoErr := middleware.GetValidatedDTO[quotation_dto.SetAgentFeeDTO](r)
	if dtoErr != nil {
		return dtoErr
	}
	updated, apiErr := h.quotationAdminService.SetAgentFee(ctx, id, *dto)
	if apiErr != nil {
		return apiErr
	}

	response, _ := json.MarshalContext(ctx, updated)
	w.Write(response)
	return nil
}
