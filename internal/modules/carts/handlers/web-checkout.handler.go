package cart_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckoutHandler struct {
	checkoutService    cart_services.CheckoutService
	validationInstance *validator.Validate
}

func NewCheckoutHandler(i *do.Injector, r *chi.Mux) {
	h := CheckoutHandler{
		checkoutService:    do.MustInvoke[cart_services.CheckoutService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.CHECKOUT.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(cart_dtos.CreateCheckoutDto{}, middleware.BODY, i),
		).Post("/try", helpers.Make(h.Try))

		r.With(
			middleware.Auth([]string{}...),
			middleware.ValidateDTO(cart_dtos.CreateCheckoutDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Checkout))

		// Debug endpoint to test working hours validation
		r.With(
			middleware.Auth([]string{}...),
		).Post("/debug-working-hours", helpers.Make(h.DebugWorkingHours))

	})
}

func (h *CheckoutHandler) Try(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[cart_dtos.CreateCheckoutDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	// Check working hours before processing
	if err := h.checkoutService.CheckWorkingHours(ctx, dto.BranchId); err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}

	result, createErr := h.checkoutService.CalculateCartItemsPrices(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[cart_dtos.CreateCheckoutDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	// Check working hours before processing
	if err := h.checkoutService.CheckWorkingHours(ctx, dto.BranchId); err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}

	result, createErr := h.checkoutService.ConvertCartToOrder(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)

	return nil
}

func (h *CheckoutHandler) DebugWorkingHours(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Parse request body to get branch ID
	var requestBody struct {
		BranchId string `json:"branchId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return helpers.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	branchId, err := primitive.ObjectIDFromHex(requestBody.BranchId)
	if err != nil {
		return helpers.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid branch ID format",
		})
	}

	// Check working hours using service
	if err := h.checkoutService.CheckWorkingHours(ctx, &branchId); err != nil {
		return helpers.WriteJson(w, http.StatusServiceUnavailable, map[string]interface{}{
			"error":  err.Message,
			"status": "closed",
		})
	}

	return helpers.WriteJson(w, http.StatusOK, map[string]interface{}{
		"status":  "open",
		"message": "Branch is currently open",
	})
}
