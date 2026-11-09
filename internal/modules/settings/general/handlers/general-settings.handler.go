package settings_general_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_general_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/general/dtos"
	settings_general_services "larsa-e-commerce-microservice-v2/internal/modules/settings/general/services"

	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type GeneralSettingsHandler struct {
	service settings_general_services.GeneralSettingsService
}

func NewGeneralSettingsHandler(i *do.Injector, r chi.Router) {
	h := &GeneralSettingsHandler{
		service: do.MustInvoke[settings_general_services.GeneralSettingsService](i),
	}

	r.Route(enums.GENERAL_SETTINGS.GetAdminRoute(), func(sr chi.Router) {
		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.GeneralSettings, enums.Get)),
		).Get("/", helpers.Make(h.getGeneralSettings))

		sr.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.GeneralSettings, enums.Put)),
			middleware.ValidateDTO(&settings_general_dtos.UpdateGeneralSettingsDTO{}, middleware.BODY, i),
		).Put("/", helpers.Make(h.updateGeneralSettings))
	})
}

func (h *GeneralSettingsHandler) getGeneralSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	settings, err := h.service.GetGeneralSettings(ctx)
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

func (h *GeneralSettingsHandler) updateGeneralSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, apiErr := middleware.GetValidatedDTO[settings_general_dtos.UpdateGeneralSettingsDTO](r)
	if apiErr != nil {
		return apiErr
	}

	updatedSettings, serviceErr := h.service.UpdateGeneralSettings(ctx, *dto)
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
