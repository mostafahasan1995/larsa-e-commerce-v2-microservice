package settings_sales_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_sales_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/dtos"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type SalesSettingsHandler struct {
	service settings_sales_services.SalesSettingsService
}

func NewSalesSettingsHandler(i *do.Injector, r chi.Router) {
	h := &SalesSettingsHandler{
		service: do.MustInvoke[settings_sales_services.SalesSettingsService](i),
	}

	r.Route(enums.SALES_SETTINGS.GetAdminRoute(), func(sr chi.Router) {
		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.SalesSettings, enums.Get)),
		).Get("/", helpers.Make(h.getSalesSettings))

		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.SalesSettings, enums.Put)),
			middleware.ValidateDTO(&settings_sales_dtos.UpdateSalesSettingsDTO{}, middleware.BODY, i),
		).Put("/", helpers.Make(h.updateSalesSettings))
	})
}

func (h *SalesSettingsHandler) getSalesSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	settings, err := h.service.GetSalesSettings(ctx)
	if err != nil {
		return err
	}
	response, marshalErr := json.MarshalContext(ctx, settings)
	if marshalErr != nil {
		return marshalErr
	}

	w.Write(response)
	return nil
}

func (h *SalesSettingsHandler) updateSalesSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, apiErr := middleware.GetValidatedDTO[settings_sales_dtos.UpdateSalesSettingsDTO](r)
	if apiErr != nil {
		return apiErr
	}

	updatedSettings, serviceErr := h.service.UpdateSalesSettings(ctx, *dto)
	if serviceErr != nil {
		return serviceErr
	}

	response, marshalErr := json.MarshalContext(ctx, updatedSettings)
	if marshalErr != nil {
		return marshalErr
	}

	w.Write(response)
	return nil
}
