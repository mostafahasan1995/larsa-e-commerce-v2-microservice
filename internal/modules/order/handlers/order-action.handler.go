package order_handlers

import (
	"encoding/json"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *OrderHandler) UpdateShippingStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		apiErr := helpers.BadRequest("Invalid order ID")
		return &apiErr
	}

	var dto order_dto.UpdateOrderStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return err
	}

	order, apiErr := h.service.UpdateOrderStatus(ctx, id, dto.Status)
	if apiErr != nil {
		return apiErr
	}

	return helpers.WriteJson(w, http.StatusOK, order)
}
