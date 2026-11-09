package config_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	cart_services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
)

type WebCouponHandler struct {
	webCouponService   services.WebCouponService
	cartService        cart_services.CartService
	validationInstance *validator.Validate
}

func NewWebCouponHandler(i *do.Injector, r *chi.Mux) {
	h := WebCouponHandler{
		webCouponService:   do.MustInvoke[services.WebCouponService](i),
		cartService:        do.MustInvoke[cart_services.CartService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.COUPON.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.Auth([]string{}...),
		).Get("/mine", h.GetMine)
		r.With(
			middleware.Auth([]string{}...),
		).Get("/{code}", h.CheckCouponRedemption)

	})
}

func (h *WebCouponHandler) GetMine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	result, getErr := h.webCouponService.GetMine(ctx)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}
	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

func (h *WebCouponHandler) CheckCouponRedemption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	couponCode := chi.URLParam(r, "code")

	// Get optional query parameters for cart calculation
	currencyCode := r.URL.Query().Get("currencyCode")
	if currencyCode == "" {
		currencyCode = configuration.EnvironmentVariables.DefaultCurrencyCode
	}
	branchId := r.URL.Query().Get("branchId")
	zoneId := r.URL.Query().Get("zoneId")
	orderTypeStr := r.URL.Query().Get("orderType")

	var orderTypeEnum order_enums.OrderType = order_enums.PICKUP
	if orderTypeStr != "" {
		orderTypeEnum = order_enums.OrderType(orderTypeStr)
	}

	var subTotal *float64
	var orderType *order_enums.OrderType

	// Get user's cart with calculation to extract subTotal and orderType
	cart, cartErr := h.cartService.GetMine(ctx, &cart_shared.GetMineArgs{
		Query:         "",
		CurrencyCode:  currencyCode,
		BranchId:      branchId,
		CouponCode:    "",
		CalculateList: true,
		ZoneId:        zoneId,
		OrderType:     orderTypeEnum,
	})

	if cartErr == nil && cart != nil && cart.Calculation != nil {
		subTotal = &cart.Calculation.SubTotalPrice
		orderType = &cart.Calculation.OrderType
	} else if orderTypeStr != "" {
		// If cart calculation failed but orderType was provided, use it
		orderType = &orderTypeEnum
	}

	result, getErr := h.webCouponService.CheckCouponRedemption(ctx, couponCode, subTotal, orderType)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
