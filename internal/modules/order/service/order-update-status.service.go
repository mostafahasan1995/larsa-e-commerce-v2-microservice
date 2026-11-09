package order_service

import (
	"context"
	"fmt"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/events"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *orderService) UpdateOrderStatus(ctx context.Context, id primitive.ObjectID, newStatusName string) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	currentOrder, currentStatusConfig, newStatusConfig, platformFlow, apiErr := s.prepareStatusTransition(ctx, id, newStatusName)
	if apiErr != nil {
		return nil, apiErr
	}

	// B2B-specific guards and hooks (refactored)
	isB2BOrder := platformFlow == settings_sales_enums.ORDER_FLOW_MODE_B2B
	if isB2BOrder {
		if err := s.handleB2BStatusTransition(ctx, currentOrder, newStatusConfig); err != nil {
			return nil, err
		}
	} else {
		if err := s.handleB2CStatusTransition(ctx, cfg, newStatusConfig, currentStatusConfig, currentOrder); err != nil {
			return nil, err
		}
	}

	updates := map[string]interface{}{"status": newStatusName}

	updatedOrder, apiErr := s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}

	// Stubbed notification for any status change
	s.notifyStatusChange(ctx, currentOrder, newStatusName)

	// Process loyalty points if the order is being completed and payment is confirmed
	if newStatusConfig.IsSystemStatus && newStatusConfig.SystemStatusType != nil &&
		*newStatusConfig.SystemStatusType == settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED {
		orderForLoyalty := order_models.Order{
			BaseOrder: updatedOrder.BaseOrder,
		}
		go func() {
			ctx := context.Background()
			ctx = context.WithValue(ctx, util.ReqAppCfg, cfg)
			_, loyaltyErr := s.loyaltyOrderService.ProcessOrderForLoyaltyPoints(ctx, &orderForLoyalty)
			if loyaltyErr != nil {
				fmt.Printf("Error processing loyalty points for order %s: %+v\n", id.Hex(), loyaltyErr)
			} else {
				// fmt.Printf("Successfully processed loyalty points for order %s\n", id.Hex())
			}
			events.PublishProductPurchaseEvent(ctx, *updatedOrder)
		}()

	}

	return updatedOrder, nil
}

// prepareStatusTransition encapsulates fetching order and sales settings,
// resolving current/new status configs, and validating applicability.
func (s *orderService) prepareStatusTransition(
	ctx context.Context,
	id primitive.ObjectID,
	newStatusName string,
) (
	currentOrder *order_models.OrderResponse,
	currentStatusConfig *settings_sales_models.OrderStatus,
	newStatusConfig *settings_sales_models.OrderStatus,
	platformFlow settings_sales_enums.OrderFlowMode,
	apiErr *helpers.ApiError,
) {
	currentOrder, apiErr = s.repo.GetById(ctx, id, nil)
	if apiErr != nil {
		return nil, nil, nil, "", apiErr
	}

	apiErr = s.ensureOnlyRole(ctx, currentOrder.Status, newStatusName)
	if apiErr != nil {
		return nil, nil, nil, "", apiErr
	}

	salesSettings, apiErr := s.settingsSalesService.GetSalesSettings(ctx)
	if apiErr != nil {
		return nil, nil, nil, "", apiErr
	}
	if salesSettings == nil || salesSettings.OrderStatuses == nil || len(salesSettings.OrderStatuses) == 0 {
		return nil, nil, nil, "", &helpers.ApiError{
			StatusCode: 400,
			Message:    "Order statuses are not configured in sales settings.",
		}
	}

	for i, os := range salesSettings.OrderStatuses {
		if os.Title == currentOrder.Status {
			currentStatusConfig = &salesSettings.OrderStatuses[i]
		}
		if os.Title == newStatusName {
			newStatusConfig = &salesSettings.OrderStatuses[i]
		}
	}

	if newStatusConfig == nil {
		return nil, nil, nil, "", &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid target order status provided: " + newStatusName + " (not found in sales settings)",
		}
	}

	platformFlow = salesSettings.OrderFlowMode
	if platformFlow == "" {
		platformFlow = settings_sales_enums.ORDER_FLOW_MODE_B2C
	}

	if len(newStatusConfig.AppliesTo) > 0 {
		valid := false
		for _, a := range newStatusConfig.AppliesTo {
			if string(a) == string(platformFlow) {
				valid = true
				break
			}
		}
		if !valid {
			return nil, nil, nil, "", &helpers.ApiError{
				StatusCode: 400,
				Message:    "Status '" + newStatusName + "' is not applicable for current order flow",
			}
		}
	}

	return currentOrder, currentStatusConfig, newStatusConfig, platformFlow, nil
}

// --- B2B transition handler ---

func (s *orderService) handleB2BStatusTransition(ctx context.Context, order *order_models.OrderResponse, newStatus *settings_sales_models.OrderStatus) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	if newStatus == nil || !newStatus.IsSystemStatus || newStatus.SystemStatusType == nil {
		return nil
	}

	switch *newStatus.SystemStatusType {
	case settings_sales_enums.SYSTEM_STATUS_TYPE_IN_PRODUCTION:
		if order.VendorInfo.Id != cfg.User.Id {
			return &helpers.ApiError{StatusCode: 400, Message: "only vendor can mark order as in production"}
		}
	case settings_sales_enums.SYSTEM_STATUS_TYPE_READ_FOR_SHIPMENT:
		if order.DeliveryInfo == nil || order.DeliveryInfo.B2BInfo == nil || order.DeliveryInfo.B2BInfo.ShipmentDetails == nil || order.DeliveryInfo.B2BInfo.ShipmentDetails.ShipmentDate == nil {
			return &helpers.ApiError{StatusCode: 400, Message: "Shipment date is required before marking as ready for shipment"}
		}
		if order.PaymentInfo.Status != order_enums.PAID {
			return &helpers.ApiError{StatusCode: 400, Message: "Payment must be paid before ready for shipment"}
		}
		// Stubs
		s.notifyShipmentWindowApproaching(ctx, order)
		s.createBrokerTodo(ctx, order)
	case settings_sales_enums.SYSTEM_STATUS_TYPE_DELIVERED:
		s.activateBuyerConfirmTags(ctx, order)
	default:
		// No special handling
	}
	return nil
}

func (s *orderService) handleB2CStatusTransition(ctx context.Context, cfg *types.AppCfg, newStatusConfig, currentStatusConfig *settings_sales_models.OrderStatus, currentOrder *order_models.OrderResponse) *helpers.ApiError {
	if currentStatusConfig != nil && currentStatusConfig.Title != newStatusConfig.Title {
		if newStatusConfig.IsSystemStatus && newStatusConfig.SystemStatusType != nil {
			if *newStatusConfig.SystemStatusType == settings_sales_enums.SYSTEM_STATUS_TYPE_CANCELED {
				if !currentStatusConfig.CancelOrder {
					return &helpers.ApiError{
						StatusCode: 400,
						Message:    "Cannot cancel order from status '" + currentStatusConfig.Title + "'.",
					}
				}
				// Continue with normal status update instead of calling CancelOrder
			}
			if *newStatusConfig.SystemStatusType == settings_sales_enums.SYSTEM_STATUS_TYPE_REFUND {
				if !currentStatusConfig.RefundOrder {
					return &helpers.ApiError{
						StatusCode: 400,
						Message:    "Cannot refund order from status '" + currentStatusConfig.Title + "'.",
					}
				}
				_, err := s.orderCancelRefundService.RefundOrder(ctx, currentOrder.Id)
				return err
			}
			if *newStatusConfig.SystemStatusType == settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED {
				if currentOrder.PaymentInfo.PaymentMethod == order_enums.CASH_ON_DELIVERY || currentOrder.PaymentInfo.PaymentMethod == order_enums.CASH_TERMINAL {
					// Update payment status to paid
					updates := map[string]interface{}{"paymentInfo.status": order_enums.PAID}
					s.repo.UpdateFieldsById(ctx, currentOrder.Id, updates, cfg.User.Id)
				}
			}
		}

		if newStatusConfig.Order > currentStatusConfig.Order { // Moving Forward
			if !currentStatusConfig.MoveForward {
				return &helpers.ApiError{
					StatusCode: 400,
					Message:    "Cannot move order forward from status '" + currentStatusConfig.Title + "' to '" + newStatusConfig.Title + "'.",
				}
			}
			if currentStatusConfig.MoveForwardOnlyIfPaymentDone && currentOrder.PaymentInfo.Status != order_enums.PAID { // ASSUMPTION: currentOrder.PaymentStatus
				return &helpers.ApiError{
					StatusCode: 400,
					Message:    "Cannot move order forward from status '" + currentStatusConfig.Title + "': Order payment is not complete.",
				}
			}
		} else if newStatusConfig.Order < currentStatusConfig.Order { // Moving Backward
			if !currentStatusConfig.MoveBackward {
				return &helpers.ApiError{
					StatusCode: 400,
					Message:    "Cannot move order backward from status '" + currentStatusConfig.Title + "' to '" + newStatusConfig.Title + "'.",
				}
			}
		}
	}
	return nil
}

// --- Notification and B2B helper stubs ---

func (s *orderService) notifyStatusChange(ctx context.Context, order *order_models.OrderResponse, newStatusName string) {
	// Guard against nil pointers
	if order == nil {
		return
	}

	events.PublishOrderStatusChangedEvent(ctx, order.Status, newStatusName, order)
}

func (s *orderService) notifyShipmentWindowApproaching(ctx context.Context, order *order_models.OrderResponse) {
	// TODO: implement daily reminder if shipment date - now <= 7 days
}

func (s *orderService) createBrokerTodo(ctx context.Context, order *order_models.OrderResponse) {
	// TODO: create a to-do in the vendor country's customs broker dashboard
}

func (s *orderService) activateBuyerConfirmTags(ctx context.Context, order *order_models.OrderResponse) {
	// TODO: activate "Awaiting Delivery" and later "Confirm Delivery" UI states
}

// TODO : check other cases  or maybe refactor
func (s *orderService) ensureOnlyRole(ctx context.Context, from, to string) *helpers.ApiError {
	if from == "canceled" && to == "completed" {
		isAdminOrSuper, err := s.rolesGw.IsUserHaveAnyRole(ctx, "administrator", "super")
		if err != nil {
			return &helpers.ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to check user roles",
			}
		}

		if !isAdminOrSuper {
			return &helpers.ApiError{
				StatusCode: http.StatusForbidden,
				Message:    "User is not authorized to update order status",
			}
		}

	}

	return nil

}
