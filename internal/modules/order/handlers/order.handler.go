package order_handlers

import (
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_service "larsa-e-commerce-microservice-v2/internal/modules/order/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	service                  order_service.OrderService
	orderSummaryService      order_service.OrderSummaryService
	orderReportsService      order_service.OrderReportsService
	orderUpdateService       order_service.OrderUpdateService
	orderCancelRefundService order_service.OrderCancelRefundService
}

func NewOrderHandler(i *do.Injector, r *chi.Mux) {
	h := OrderHandler{
		service:                  do.MustInvoke[order_service.OrderService](i),
		orderSummaryService:      do.MustInvoke[order_service.OrderSummaryService](i),
		orderReportsService:      do.MustInvoke[order_service.OrderReportsService](i),
		orderUpdateService:       do.MustInvoke[order_service.OrderUpdateService](i),
		orderCancelRefundService: do.MustInvoke[order_service.OrderCancelRefundService](i),
	}

	r.Route(enums.ORDER.GetAdminRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Post)),
		).Post("/", helpers.Make(h.Create))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.GetOne)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Get("/{id}", helpers.Make(h.GetById))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Get)),
		).Get("/", helpers.Make(h.GetWithPagination))

		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Put)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(order_dto.UpdateOrderDTO{}, middleware.BODY, i),
		).Put("/{id}", helpers.Make(h.UpdateOrder))

		//update order status
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Put)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/status/{status}", helpers.Make(h.UpdateOrderStatus))

		//update order transition status
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Put)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/transition-status/{status}", helpers.Make(h.UpdateTransitionStatus))

		//update admin seen status
		r.With(
			middleware.AdminAuth(middleware.CreateRestriction(enums.CapabilitiesRegistry.Order, enums.Put)),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
		).Put("/{id}/admin-seen", helpers.Make(h.UpdateAdminSeen))

		r.Delete("/all", helpers.Make(h.DeleteAllOrders))

		h.NewOrderReportHandler(i, r)
		h.NewOrderSummaryHandler(i, r)
	})
}

func (h *OrderHandler) DeleteAllOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	return h.service.DeleteAllOrders(ctx)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	var dto order_dto.CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	order, apiErr := h.service.Create(ctx, &dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetWithPagination(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		apiErr := helpers.DefaultError(errGetPaginate)
		return apiErr
	}

	result, apiErr := h.service.GetWithPagination(ctx, skip, limit, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, result)
}

func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	query := r.URL.Query().Get("query")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	order, apiErr := h.service.GetById(ctx, id, query)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	var dto order_dto.UpdateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	order, apiErr := h.orderUpdateService.UpdateOrder(ctx, id, &dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	newStatusName := chi.URLParam(r, "status")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	order, apiErr := h.service.UpdateOrderStatus(ctx, id, newStatusName)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateAdminSeen(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	var dto order_dto.UpdateAdminSeenDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	order, apiErr := h.service.UpdateAdminSeen(ctx, id, &dto)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateTransitionStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	status := chi.URLParam(r, "status")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	order, apiErr := h.service.UpdateTransitionStatus(ctx, id, status)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}
