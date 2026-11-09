package order_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: things needed to be done
// 1. apply assign to driver feature with status
// 2. apply confirm delivery feature with status
type OrderCancelRefundService interface {
	RefundOrder(ctx context.Context, id primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError)
	CancelOrder(ctx context.Context, id primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError)
}

type orderCancelRefundService struct {
	repo order_repos.OrderRepo
}

func NewOrderCancelRefundService(i *do.Injector) (OrderCancelRefundService, error) {
	repo := do.MustInvoke[order_repos.OrderRepo](i)
	return &orderCancelRefundService{
		repo: repo,
	}, nil
}

func (s *orderCancelRefundService) RefundOrder(ctx context.Context, id primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError) {
	_, err := s.repo.GetById(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Implement refund logic
	// 1. Check if order is eligible for refund
	// 2. Process refund through payment provider
	// 3. Update order payment status
	// 4. Add log history entry using s.addLogHistory()
	// 5. Return updated order

	return nil, &helpers.ApiError{
		StatusCode: 501,
		Message:    "Refund functionality not yet implemented",
	}
}

func (s *orderCancelRefundService) CancelOrder(ctx context.Context, id primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError) {
	_, err := s.repo.GetById(ctx, id, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Implement cancel logic
	// 1. Check if order is eligible for cancellation
	// 2. Update order status
	// 3. Handle inventory updates
	// 4. Add log history entry using s.addLogHistory()
	// 5. Return updated order

	return nil, &helpers.ApiError{
		StatusCode: 501,
		Message:    "Cancel functionality not yet implemented",
	}
}
