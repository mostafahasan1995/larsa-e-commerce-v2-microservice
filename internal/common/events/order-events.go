package events

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
)

const SubjectProductPurchase = "automation.event"
const automationEvent = "automation.event"

// PublishProductPurchaseEvent publishes a product-purchase event to NATS when order is completed
func PublishProductPurchaseEvent(ctx context.Context, orderData order_models.OrderResponse) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return
	}

	// Extract product IDs from order items
	var productIds []string
	for _, item := range orderData.OrderItems {
		productIds = append(productIds, item.ItemId.Hex())
	}

	eventData := map[string]interface{}{
		"id":          primitive.NewObjectID().Hex(),
		"client":      cfg.Hp.Client,
		"type":        "product-purchase",
		"userId":      cfg.User.Id.Hex(),
		"serviceName": "ecommerce",
		"data": map[string]interface{}{
			"userId":             cfg.User.Id.Hex(),
			"products":           productIds,
			"orderId":            orderData.Id.Hex(),
			"orderCode":          orderData.Code,
			"totalPrice":         orderData.TotalPrice,
			"subTotalPrice":      orderData.SubTotalPrice,
			"totalDiscountValue": orderData.TotalDiscountValue,
			"currencyCode":       orderData.CurrencyCode,
			"status":             orderData.Status,
			"invoiceNumber":      orderData.InvoiceNumber,
			"email":              cfg.User.UserData.Email,
			"userName":           cfg.User.UserData.FirstName + " " + cfg.User.UserData.LastName,
			"orderType":          orderData.Type,
			"itemCount":          len(orderData.OrderItems),
		},
	}

	PublishNats(SubjectProductPurchase, eventData)
}

func PublishOrderStatusChangedEvent(ctx context.Context, oldStatus, newStatus string, orderData *order_models.OrderResponse) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return
	}

	userId := orderData.BuyerInfo.Id.Hex()

	eventData := map[string]interface{}{
		"client":      cfg.Hp.Client,
		"type":        "order-status-changed",
		"userId":      userId,
		"serviceName": "ecommerce",
		"data": map[string]any{
			"userId":                userId,
			"fullName":              orderData.BuyerInfo.FullName,
			"email":                 orderData.BuyerInfo.Email,
			"countryCode":           orderData.BuyerInfo.CountryCode,
			"phone":                 orderData.BuyerInfo.Phone,
			"alternativeMobile":     orderData.BuyerInfo.AlternativeMobile,
			"alternativeMobileCode": orderData.BuyerInfo.AlternativeMobileCode,
			"type":                  orderData.BuyerInfo.Type,
			"newStatus":             newStatus,
			"oldStatus":             oldStatus,
		},
	}

	PublishNats(automationEvent, eventData)
}
