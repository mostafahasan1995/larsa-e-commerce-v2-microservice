package events

import (
	"encoding/json"
	"fmt"

	"git.larsa.io/mahdawi/microservices-commons.git/nats"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
)

// PublishAddToCartEvent publishes an add-to-cart event to NATS
func PublishAddToCartEvent(clientId, userId, email string, dto *cart_dtos.UpsertCartItemDto) {
	eventData := map[string]interface{}{
		"id":          primitive.NewObjectID().Hex(),
		"client":      clientId,
		"type":        "add-to-cart",
		"userId":      userId,
		"serviceName": "ecommerce",
		"data": map[string]interface{}{
			"itemId":   dto.ItemId,
			"userId":   userId,
			"quantity": dto.Quantity,
			"type":     dto.Type,
			"branchId": dto.BranchId,
			"unitId":   dto.UnitId,
			"email":    email,
		},
	}

	eventDataJSON, _ := json.MarshalIndent(eventData, "", "  ")
	fmt.Printf("publishAddToCartEvent eventData %s\n", string(eventDataJSON))

	PublishNats(SubjectAddToCart, eventData)
}

// PublishAccessCheckoutEvent publishes an access-checkout-page event to NATS
func PublishAccessCheckoutEvent(clientId string, checkoutData order_models.OrderResponse) {
	eventData := map[string]interface{}{
		"id":          primitive.NewObjectID().Hex(),
		"client":      clientId,
		"type":        "access-checkout-page",
		"userId":      checkoutData.BuyerInfo.Id,
		"serviceName": "ecommerce",
		"data": map[string]interface{}{
			"user_name":          checkoutData.BuyerInfo.FullName,
			"order_id":           checkoutData.Id,
			"coupon":             checkoutData.CouponCode,
			"email":              checkoutData.BuyerInfo.Email,
			"userId":             checkoutData.BuyerInfo.Id,
			"totalPrice":         checkoutData.TotalPrice,
			"totalDiscountValue": checkoutData.TotalDiscountValue,
			"subTotalPrice":      checkoutData.SubTotalPrice,
			"currencyCode":       checkoutData.CurrencyCode,
			"status":             checkoutData.Status,
			"paymentInfo":        checkoutData.PaymentInfo,
			"invoiceNumber":      checkoutData.InvoiceNumber,
		},
	}

	eventDataJSON, _ := json.MarshalIndent(eventData, "", "  ")
	fmt.Printf("publishAccessCheckoutEvent eventData %s\n", string(eventDataJSON))

	PublishNats(SubjectAccessCheckout, eventData)
}

// TODO: make this function global
func PublishNats(topic string, data any) {
	env := configuration.EnvironmentVariables.NatsJetStream
	if env == "false" {
		nats.Publish(topic, data)
	} else {
		nats.PublishJet(topic, data)
	}
}
