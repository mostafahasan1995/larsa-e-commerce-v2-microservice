package order_service

import (
	"context"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	settings_sales_service "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/samber/do"
)

type OrderWebService interface {
	GetMyOrders(ctx context.Context, skip int64, limit int64, query string) (*order_models.OrdersWithPagination, *helpers.ApiError)
	PayOrder(ctx context.Context, dto *order_dto.PayOrderDTO) (*order_dto.PaymentInitiationResponseDTO, *helpers.ApiError)
}

type orderWebService struct {
	repo                 order_repos.OrderWebRepo
	settingsSalesService settings_sales_service.SalesSettingsService
	paymentGateway       *gateway.PaymentGw
}

func NewOrderWebService(i *do.Injector) (OrderWebService, error) {
	repo := do.MustInvoke[order_repos.OrderWebRepo](i)
	settingsSalesService := do.MustInvoke[settings_sales_service.SalesSettingsService](i)
	paymentGateway := do.MustInvoke[*gateway.PaymentGw](i)
	return &orderWebService{
		repo:                 repo,
		settingsSalesService: settingsSalesService,
		paymentGateway:       paymentGateway,
	}, nil
}

func (s *orderWebService) GetMyOrders(ctx context.Context, skip int64, limit int64, query string) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError) {
	config, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	buyerId := config.User.Id

	filters, err := order_filters.NewOrderFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	filters.BuyerId = &buyerId

	return s.repo.GetWithPagination(ctx, skip, limit, *filters)
}

func (s *orderWebService) PayOrder(ctx context.Context, dto *order_dto.PayOrderDTO) (*order_dto.PaymentInitiationResponseDTO, *helpers.ApiError) {
	config, _ := util.GetReqAppCfg(ctx)
	userId := config.User.Id

	order, apiErr := s.repo.GetById(ctx, dto.OrderId, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	if order.PaymentInfo != nil && order.PaymentInfo.Status == order_enums.PAID {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Order is already paid",
		}
	}
	successURL := ""
	cancelURL := ""
	username := order.BuyerInfo.FullName
	email := order.BuyerInfo.Email
	mobile := order.BuyerInfo.Phone
	productId := dto.OrderId.Hex()
	orderRequest := gateway.Order{
		ProductId:  productId,
		Amount:     order.TotalPrice,
		Currency:   order.CurrencyCode,
		SuccessURL: &successURL,
		CancelURL:  &cancelURL,
		Email:      email,
		Username:   username,
		Mobile:     mobile,
	}

	paymentResponse, err := s.paymentGateway.PayOrder(ctx, &gateway.OrderRequest{
		Provider: dto.PaymentGateway,
		Order:    orderRequest,
	})

	if err != nil {
		fmt.Printf("Failed to initiate payment: %v", err)
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Failed to initiate payment",
		}
	}

	updates := map[string]interface{}{
		"logHistory": append(order.LogHistory, order_models.LogHistory{
			Action:      order_enums.PAYMENT_ADDED,
			Description: "Payment initiated",
			Date:        time.Now(),
		}),
	}

	_, updateErr := s.repo.UpdateFieldsById(ctx, dto.OrderId, updates, userId)
	if updateErr != nil {
		return nil, updateErr
	}
	//
	return &order_dto.PaymentInitiationResponseDTO{
		OrderId:     dto.OrderId,
		IsRedirect:  paymentResponse.IsRedirect,
		RedirectURL: paymentResponse.RedirectURL,
	}, nil
}
