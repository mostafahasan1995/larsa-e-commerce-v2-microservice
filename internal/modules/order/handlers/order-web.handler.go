package order_handlers

import (
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_service "larsa-e-commerce-microservice-v2/internal/modules/order/service"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderWebHandler struct {
	service            order_service.OrderWebService
	validationInstance interface{}
}

// s
func NewOrderWebHandler(i *do.Injector, r *chi.Mux) {
	h := OrderWebHandler{
		service: do.MustInvoke[order_service.OrderWebService](i),
	}

	r.Route(enums.ORDER.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
		).Get("/sent", helpers.Make(h.GetMyOrders))

		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
		).Post("/pay", helpers.Make(h.PayOrder))
	})
}

func (h *OrderWebHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		apiErr := helpers.DefaultError(errGetPaginate)
		return apiErr
	}

	upcomingOrders, apiErr := h.service.GetMyOrders(ctx, skip, limit, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, upcomingOrders)
}

func (h *OrderWebHandler) PayOrder(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	orderIdParam := r.URL.Query().Get("orderId")
	paymentGateway := r.URL.Query().Get("paymentGateway")
	orderId, err := primitive.ObjectIDFromHex(orderIdParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	paymentResponse, apiErr := h.service.PayOrder(ctx, &order_dto.PayOrderDTO{
		OrderId:        orderId,
		PaymentGateway: order_enums.PaymentGateway(paymentGateway),
	})

	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, paymentResponse)
}
