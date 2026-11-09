package loyalty_settings_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_settings_dtos "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/dtos"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type LoyaltySettingsHandler struct {
	loyaltySettingsService loyalty_settings_services.LoyaltySettingsService
	validationInstance     *validator.Validate
}

func NewLoyaltySettingsHandler(i *do.Injector, r *chi.Mux) {
	h := LoyaltySettingsHandler{
		loyaltySettingsService: do.MustInvoke[loyalty_settings_services.LoyaltySettingsService](i),
		validationInstance:     do.MustInvoke[*validator.Validate](i),
	}

	// Admin routes
	r.Route(enums.LOYALTY_SETTINGS.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.LoyaltySettings, enums.Post)),
			middleware.ValidateDTO(loyalty_settings_dtos.CreateLoyaltySettingsDTO{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Upsert))
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.LoyaltySettings, enums.Get)),
		).Get("/", helpers.Make(h.GetActiveSettings))
	})
}

func (h *LoyaltySettingsHandler) Upsert(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var dto loyalty_settings_dtos.CreateLoyaltySettingsDTO
	json.NewDecoder(r.Body).Decode(&dto)

	result, err := h.loyaltySettingsService.Upsert(ctx, &dto)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *LoyaltySettingsHandler) GetActiveSettings(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.loyaltySettingsService.GetActiveSettings(ctx)
	if err != nil {
		return err
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}
