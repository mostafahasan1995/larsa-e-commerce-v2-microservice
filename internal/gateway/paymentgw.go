package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentGw struct {
}

type OrderRequest struct {
	Provider order_enums.PaymentGateway `json:"provider"`
	Order    Order                      `json:"order"`
}

type Order struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	SuccessURL  *string `json:"successURL"`
	CancelURL   *string `json:"cancelURL"`
	Callback    string  `json:"callback,omitempty"` // to use with geidea
	ProductType *string `json:"productType"`
	Type        *string `json:"type"`
	//ProductId is the id of the order
	ProductId string       `json:"productId"`
	Email     *string      `json:"email"`
	Username  string       `json:"username"`
	Mobile    *string      `json:"mobile"`
	Items     *[]OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ProductName string  `json:"productName"`
	UnitPrice   float64 `json:"unitPrice"`
	Quantity    int     `json:"quantity"`
}

type GatewayInternalPaymentResponse struct {
	Result string `json:"result"`
	Flow   string `json:"flow"`
	Data   struct {
		RedirectURL string                 `json:"redirectURL"`
		Order       map[string]interface{} `json:"order"`
		Payment     struct {
			Id primitive.ObjectID `json:"id"`
		} `json:"payment"`
	} `json:"data"`
}

func NewPaymentGw(i *do.Injector) (*PaymentGw, error) {
	return &PaymentGw{}, nil
}

func (ugw *PaymentGw) PayOrder(ctx context.Context, requestData *OrderRequest) (*order_dto.PaymentInitiationResponseDTO, error) {
	cfg, _ := util.GetReqAppCfg(ctx)

	if requestData.Order.Email != nil && !util.IsValidEmail(*requestData.Order.Email) {
		return nil, errors.New("invalid email format")
	}

	if requestData.Order.Mobile != nil && !util.IsValidPhoneNumber(*requestData.Order.Mobile) {
		return nil, errors.New("invalid phone number format")
	}

	var successURL, cancelURL string

	frontendUrl, err := common.GetOptionValue("FRONTEND_URL", cfg.Hp)
	if err != nil {
		return nil, errors.New("frontend url not found in options")
	}
	if url, ok := frontendUrl.(string); ok {
		successURL = url + "/cart"
		cancelURL = url + "/cancel"
	}

	productType := "ecommerceOrder"
	orderType := "regular"
	apiRequest := &OrderRequest{
		Provider: requestData.Provider,
		Order: Order{
			Amount:      requestData.Order.Amount,
			Currency:    requestData.Order.Currency,
			SuccessURL:  &successURL,
			CancelURL:   &cancelURL,
			ProductType: &productType,
			Type:        &orderType,
			ProductId:   requestData.Order.ProductId,
			Email:       requestData.Order.Email,
			Username:    requestData.Order.Username,
			Mobile:      requestData.Order.Mobile,
			Items:       requestData.Order.Items,
		},
	}

	if requestData.Provider == order_enums.STRIPE {
		apiRequest.Order.Items = &[]OrderItem{
			{
				ProductName: "Total Order Price",
				UnitPrice:   requestData.Order.Amount,
				Quantity:    1,
			},
		}
	}

	if requestData.Provider == order_enums.GEIDEA {
		callbackURL, err := common.GetOptionValue("GEIDEA_CALLBACK_URL", cfg.Hp) // try get from options - useful when we deploy on different servers
		if err == nil {
			if url, ok := callbackURL.(string); ok {
				fmt.Println("geidea callbackURL", url)
				apiRequest.Order.Callback = url
			}
		} else {
			fmt.Println("geidea callbackURL not found in options", err)
		}
	}

	if requestData.Provider == order_enums.PAYMOB {
		callbackURL, err := common.GetOptionValue("PAYMOB_CALLBACK_URL", cfg.Hp) // try get from options - useful when we deploy on different servers
		if err == nil {
			if url, ok := callbackURL.(string); ok {
				fmt.Println("paymob callbackURL", url)
				apiRequest.Order.Callback = url
			}
		} else {
			fmt.Println("paymob callbackURL not found in options", err)
		}
	}

	var orderRequestMap map[string]interface{}
	jsonData, err := json.Marshal(apiRequest)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &orderRequestMap)
	if err != nil {
		return nil, err
	}
	fmt.Println("orderRequestMap", string(jsonData))
	reqOptions := &common.RequestParams{
		Service: "payment",
		Path:    "payments",
		Method:  "POST",
		Data:    orderRequestMap,
		Header:  cfg.Hp,
	}

	pRes, pErr := common.CallService(reqOptions)
	if pErr != nil {
		fmt.Printf("Failed to call payment service: %v", pErr)
		return nil, pErr
	}
	defer pRes.Body.Close()
	bodyBytes, readErr := io.ReadAll(pRes.Body)
	if readErr != nil {
		return nil, readErr
	}

	fmt.Println("payment request response", string(bodyBytes))
	var gatewayPaymentResp GatewayInternalPaymentResponse
	if err := json.Unmarshal(bodyBytes, &gatewayPaymentResp); err != nil {
		return nil, err
	}

	if gatewayPaymentResp.Result != "success" {
		return nil, errors.New("payment failed: " + string(bodyBytes))
	}
	response := &order_dto.PaymentInitiationResponseDTO{
		OrderId:     gatewayPaymentResp.Data.Payment.Id,
		IsRedirect:  gatewayPaymentResp.Flow == "redirect",
		RedirectURL: &gatewayPaymentResp.Data.RedirectURL,
	}

	return response, nil
}
