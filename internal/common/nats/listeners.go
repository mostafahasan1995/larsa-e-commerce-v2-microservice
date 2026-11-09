package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	"time"

	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"git.larsa.io/mahdawi/microservices-commons.git/nats"
	natss "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentOrder struct {
	Id          primitive.ObjectID `json:"_id"`
	ProductId   primitive.ObjectID `json:"productId"`
	ProductType string             `json:"productType"`
}

type WebhookMsg struct {
	Headers struct {
		Client       string `json:"xClient"`
		Service      string `json:"xService"`
		ServiceToken string `json:"xServiceToken"`
	} `json:"headers"`
	Data struct {
		Order   PaymentOrder `json:"order"`
		Payment struct {
			Id primitive.ObjectID `json:"_id"`
		} `json:"payment"`
	}
	ReceiptURL string `json:"receiptURL"`
	InvoicePDF string `json:"invoicePDF"`
	Result     string `json:"result"`
}

func InitListeners(injector *do.Injector) {
	env := configuration.EnvironmentVariables.NatsJetStream
	if env == "false" {
		// fmt.Println("starting nats listeners to: ", env)
		InitNatsListeners(injector)
	} else {
		// fmt.Println("starting nats jet listeners to: ", env)
		InitNatsJet(injector)
	}

}

func InitNatsListeners(injector *do.Injector) {
	nats.Listen(SubjectUpdatePayment, "ecommerce", "ecommerce", func(m *stan.Msg) {
		m.Ack()
		var msg WebhookMsg
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			fmt.Printf("error unmarshalling message: %+v\n", err)
			return
		}
		// fmt.Printf("listening to ecommerce msg: %+v\n", msg)

		switch msg.Data.Order.ProductType {
		case "ecommerceOrder":
			handleOrder(injector, msg)
		}

	})
}

func InitNatsJet(injector *do.Injector) {
	natsJetLinsterError := nats.ListenJet(SubjectUpdatePayment, "ecommerce", "ecommerce", func(m *natss.Msg) {
		m.Ack()
		var msg WebhookMsg
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			fmt.Printf("error unmarshalling message: %+v\n", err)
			return
		}

		log.Info().Msgf("listening to ecommerce msg: %+v", msg)

		switch msg.Data.Order.ProductType {
		case "ecommerceOrder":
			handleOrder(injector, msg)
		}
	})
	if natsJetLinsterError != nil {
		fmt.Printf("error listening to nats jet stream: %+v\n", natsJetLinsterError)
	}
}

func handleOrder(injector *do.Injector, msg WebhookMsg) {

	log.Info().Msgf("handling order: %+v", msg)
	// Create a context with application configuration
	cfg := &types.AppCfg{
		Db:   msg.Headers.Client,
		Lang: "en",
		Hp: &common.HeaderParams{
			Client: msg.Headers.Client,
		},
	}
	ctx := context.WithValue(context.TODO(), util.ReqAppCfg, cfg)

	orderRepo := do.MustInvoke[order_repos.OrderRepo](injector)
	logHistory := []order_models.LogHistory{
		{
			Action:      order_enums.PAYMENT_ADDED,
			Description: "Payment confirmed and order activated",
			Date:        time.Now(),
		},
	}

	// Create filter that includes deleted orders since we need to update soft-deleted orders
	includeDeleted := true
	filter := &order_filters.OrderFilters{
		IncludeOnlineUnPaidOrders: true,
		BaseFilters: base_filters.BaseFilters{
			IncludeDeleted: &includeDeleted,
		},
	}

	// fmt.Printf("getting order by id in handleOrder listener: %+v\n", msg.Data.Order.ProductId)
	order, err := orderRepo.GetById(ctx, msg.Data.Order.ProductId, filter)
	if err != nil {
		// fmt.Printf("error getting order: %+v\n", err)
		return
	}

	// Check if payment already exists to prevent duplication
	paymentExists := false
	paymentNumber := msg.Data.Payment.Id.Hex()
	for i, p := range order.PaymentInfo.Payments {
		if p.PaymentNumber == paymentNumber {
			order.PaymentInfo.Payments[i].PaymentStatus = order_enums.PAID
			order.PaymentInfo.Payments[i].PaymentDate = time.Now()
			paymentExists = true
			break
		}
	}

	if !paymentExists {
		payToVendor := order_enums.PAY_TO_VENDOR
		order.PaymentInfo.Payments = append(order.PaymentInfo.Payments, order_models.Payments{
			PaymentMethod: order_enums.ONLINE_PAYMENT,
			PaymentStatus: order_enums.PAID,
			PaymentDate:   time.Now(),
			PaymentAmount: order.TotalPrice,
			PaymentNumber: paymentNumber,
			PaidTo:        &payToVendor,
			Attachment:    &types.FileFieldDto{},
		})
	}
	order.PaymentInfo.Status = order_enums.PAID

	updates := map[string]interface{}{
		"logHistory":  logHistory,
		"paymentInfo": order.PaymentInfo,
		"deletedAt":   nil,
	}
	updatedOrder, err := orderRepo.UpdateFieldsById(ctx, msg.Data.Order.ProductId, updates, primitive.NilObjectID)
	if err != nil {
		fmt.Printf("error updating order: %+v\n", err)
		return
	}
	// fmt.Printf("updatedOrder: %+v\n", updatedOrder)

	// Process loyalty points after payment confirmation
	loyaltyOrderService := do.MustInvoke[loyalty_services.LoyaltyOrderService](injector)
	orderForLoyalty := order_models.Order{
		BaseOrder: updatedOrder.BaseOrder,
	}
	_, loyaltyErr := loyaltyOrderService.ProcessOrderForLoyaltyPoints(ctx, &orderForLoyalty)
	if loyaltyErr != nil {
		fmt.Printf("error processing loyalty points: %+v\n", loyaltyErr)
		// Don't fail the order update if loyalty processing fails
	}

	if err := EmptyClonedCart(ctx, injector, order); err != nil {
		fmt.Printf("error handling cart after payment: %+v\n", err)
		// Don't fail the order update if cart handling fails
	}
}

func EmptyClonedCart(ctx context.Context, i *do.Injector, order *order_models.OrderResponse) error {
	db := do.MustInvoke[*mongo.Client](i)

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("error emptying cloned cart: %+v\n", err)
		}
	}()

	cartRepo := *base_repos.NewBaseRepo[cart_models.Cart, cart_models.CartResponse](db, enums.CART)
	cartItemRepo := *base_repos.NewBaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse](db, enums.CART_ITEM)

	clonedCart, err := cartRepo.FindOne(ctx, []bson.M{
		{
			"$match": bson.M{
				"refId":   *order.CartId,
				"isClone": true,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error getting cart: %+v", err)
	}

	if _, err := cartRepo.BulkWrite(ctx, []mongo.WriteModel{
		mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": *clonedCart.Id}),
	}); err != nil {
		return fmt.Errorf("error deleting cloned carts: %+v", err)
	}

	if _, err := cartItemRepo.BulkWrite(ctx, []mongo.WriteModel{
		mongo.NewDeleteManyModel().SetFilter(bson.M{"cartId": *clonedCart.Id}),
	}); err != nil {
		return fmt.Errorf("error deleting cloned cart items: %+v", err)
	}

	return nil
}
