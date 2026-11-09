package product_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

type DiscountHandler struct {
	discountService    services.DiscountService
	validationInstance *validator.Validate
}

func NewDiscountHandler(i *do.Injector, r *chi.Mux) {
	h := DiscountHandler{
		discountService:    do.MustInvoke[services.DiscountService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.DISCOUNT.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Addon, enums.GetOne)),
			middleware.ValidateDTO(dtos.HandleNewProductsDiscountDto{}, middleware.BODY, i),
		).Put("/", helpers.Make(h.HandleNewProductsDiscount))
	})
}

func (h *DiscountHandler) HandleNewProductsDiscount(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[dtos.HandleNewProductsDiscountDto](r)
	if err != nil {
		helpers.DefaultError(err)
		return err
	}

	updateErr := h.discountService.HandleNewProductsDiscount(ctx, *dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
