package payment_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayoutRepo interface {
	Create(ctx context.Context, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PayoutResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PayoutsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type PayoutService interface {
	Create(ctx context.Context, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PayoutResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PayoutsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type payout struct {
	repo PayoutRepo
}

func NewPayoutService(i *do.Injector) (PayoutService, error) {
	return &payout{
		repo: do.MustInvoke[PayoutRepo](i),
	}, nil
}

func (p *payout) Create(ctx context.Context, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, payout)
}

func (p *payout) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.PayoutResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *payout) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.PayoutsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *payout) UpdateById(ctx context.Context, id primitive.ObjectID, payout *dtos.CreatePayoutDto) (*models.PayoutResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, payout)
}

func (p *payout) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *payout) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
