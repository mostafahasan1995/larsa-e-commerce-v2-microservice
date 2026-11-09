package loyalty_handlers

import (
	"encoding/json"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_dtos "larsa-e-commerce-microservice-v2/internal/modules/loyalty/dtos"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type LoyaltyOrderHandler struct {
	loyaltyOrderService loyalty_services.LoyaltyOrderService
	validationInstance  *validator.Validate
}

func NewLoyaltyOrderHandler(i *do.Injector, r *chi.Mux) {
	h := LoyaltyOrderHandler{
		loyaltyOrderService: do.MustInvoke[loyalty_services.LoyaltyOrderService](i),
		validationInstance:  do.MustInvoke[*validator.Validate](i),
	}

	r.With(
		middleware.Auth([]string{}...),
	).Post(enums.LOYALTY_POINTS.GetWebRoute()+"/assign-order", h.AssignOrderForPoints)

	// r.With(
	// 	middleware.Auth([]string{}...),
	// ).Post(enums.LOYALTY_POINTS.GetWebRoute()+"/process-discount-usage", h.ProcessDiscountUsage)
}

func (h *LoyaltyOrderHandler) AssignOrderForPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var dto loyalty_dtos.AssignOrderForPointsDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid request body",
		})
		return
	}

	if err := h.validationInstance.Struct(dto); err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Validation failed: " + err.Error(),
		})
		return
	}

	// Create a minimal order structure for loyalty processing
	// In a real scenario, you might want to validate the order exists first
	cfg, _ := util.GetReqAppCfg(ctx)
	order := &order_models.Order{
		BaseOrder: order_models.BaseOrder{
			BuyerInfo: order_models.OrderParticipantInfo{
				Id: cfg.User.Id, // Use the authenticated user as buyer
			},
			PaymentInfo: &order_models.PaymentInfo{
				Status: order_enums.PAID, // Assume order is paid for loyalty points
			},
		},
	}
	order.Id = dto.OrderId

	// Process order for loyalty points
	result, loyaltyErr := h.loyaltyOrderService.ProcessOrderForLoyaltyPoints(ctx, order)
	if loyaltyErr != nil {
		helpers.HandleApiError(w, loyaltyErr)
		return
	}

	response, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// func (h *LoyaltyOrderHandler) ProcessDiscountUsage(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	ctx, _ := util.AddCtxAppCfg(r)

// 	var dto loyalty_dtos.ProcessDiscountUsageDTO
// 	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
// 		helpers.HandleApiError(w, &helpers.ApiError{
// 			StatusCode: 400,
// 			Message:    "Invalid request body",
// 		})
// 		return
// 	}

// 	if err := h.validationInstance.Struct(dto); err != nil {
// 		helpers.HandleApiError(w, &helpers.ApiError{
// 			StatusCode: 400,
// 			Message:    "Validation failed: " + err.Error(),
// 		})
// 		return
// 	}

// 	cfg, _ := util.GetReqAppCfg(ctx)
// 	userId := cfg.User.Id

// 	// Process loyalty discount usage
// 	result, loyaltyErr := h.loyaltyOrderService.ProcessLoyaltyDiscountUsage(ctx, userId, dto.OrderId, dto.DiscountAmount)
// 	if loyaltyErr != nil {
// 		helpers.HandleApiError(w, loyaltyErr)
// 		return
// 	}

// 	response, _ := json.Marshal(result)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }
