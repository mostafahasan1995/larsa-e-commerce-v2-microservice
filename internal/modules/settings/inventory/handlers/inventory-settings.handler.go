package settings_inventory_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_inventory_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/dtos"
	settings_inventory_services "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/services"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type InventorySettingsHandler struct {
	service settings_inventory_services.InventorySettingsService
}

func NewInventorySettingsHandler(i *do.Injector, r chi.Router) {
	h := &InventorySettingsHandler{
		service: do.MustInvoke[settings_inventory_services.InventorySettingsService](i),
	}

	r.Route(enums.INVENTORY_SETTINGS.GetAdminRoute(), func(sr chi.Router) {
		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.InventorySettings, enums.Get)),
		).Get("/", helpers.Make(h.getInventorySettings))

		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.InventorySettings, enums.Put)),
			middleware.ValidateDTO(&settings_inventory_dtos.UpdateInventorySettingsDTO{}, middleware.BODY, i),
		).Put("/", helpers.Make(h.updateInventorySettings))
	})
}

func (h *InventorySettingsHandler) getInventorySettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	settings, err := h.service.GetInventorySettings(ctx)
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

func (h *InventorySettingsHandler) updateInventorySettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, apiErr := middleware.GetValidatedDTO[settings_inventory_dtos.UpdateInventorySettingsDTO](r)
	if apiErr != nil {
		return apiErr
	}

	updatedSettings, serviceErr := h.service.UpdateInventorySettings(ctx, *dto)
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
