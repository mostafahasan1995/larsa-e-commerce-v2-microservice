package payment_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepo interface {
	Create(ctx context.Context, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PaymentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PaymentsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type PaymentService interface {
	Create(ctx context.Context, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PaymentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PaymentsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type payment struct {
	repo PaymentRepo
}

func NewPaymentService(i *do.Injector) (PaymentService, error) {
	return &payment{
		repo: do.MustInvoke[PaymentRepo](i),
	}, nil
}

func (p *payment) Create(ctx context.Context, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, payment)
}

func (p *payment) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PaymentResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *payment) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PaymentsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *payment) UpdateById(ctx context.Context, id primitive.ObjectID, payment *dtos.CreatePaymentDto) (*models.PaymentResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, payment)
}

func (p *payment) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *payment) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
