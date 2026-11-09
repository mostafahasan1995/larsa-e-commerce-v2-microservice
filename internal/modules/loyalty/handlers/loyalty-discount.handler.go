package loyalty_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type LoyaltyDiscountHandler struct {
	loyaltyDiscountService loyalty_services.LoyaltyDiscountService
	validationInstance     *validator.Validate
}

func NewLoyaltyDiscountHandler(i *do.Injector, r *chi.Mux) {
	h := LoyaltyDiscountHandler{
		loyaltyDiscountService: do.MustInvoke[loyalty_services.LoyaltyDiscountService](i),
		validationInstance:     do.MustInvoke[*validator.Validate](i),
	}

	r.With(
		middleware.Auth([]string{}...),
	).Get(enums.LOYALTY_POINTS.GetWebRoute()+"/discount/eligibility", h.CheckDiscountEligibility)
}

func (h *LoyaltyDiscountHandler) CheckDiscountEligibility(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	cfg, _ := util.GetReqAppCfg(ctx)
	userId := cfg.User.Id

	result, err := h.loyaltyDiscountService.CheckDiscountEligibility(ctx, userId)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to check discount eligibility",
		})
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
