package loyalty_handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"
)

type WebLoyaltyPointsHandler struct {
	loyaltyPointsService   loyalty_services.WebLoyaltyPointsService
	loyaltySettingsService loyalty_settings_services.LoyaltySettingsService
	validationInstance     *validator.Validate
}

func NewWebLoyaltyPointsHandler(i *do.Injector, r *chi.Mux) {
	h := WebLoyaltyPointsHandler{
		loyaltyPointsService:   do.MustInvoke[loyalty_services.WebLoyaltyPointsService](i),
		loyaltySettingsService: do.MustInvoke[loyalty_settings_services.LoyaltySettingsService](i),
		validationInstance:     do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.LOYALTY_POINTS.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
		).Get("/my-loyalty", h.GetMyLoyalty)

		r.With().Get("/settings", h.GetLoyaltySettings)
	})
}

func (h *WebLoyaltyPointsHandler) GetMyLoyalty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	cfg, _ := util.GetReqAppCfg(ctx)
	userId := cfg.User.Id

	result, err := h.loyaltyPointsService.GetByUserId(ctx, userId)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebLoyaltyPointsHandler) GetLoyaltySettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, err := h.loyaltySettingsService.GetActiveSettings(ctx)
	if err != nil {
		helpers.HandleApiError(w, err)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
