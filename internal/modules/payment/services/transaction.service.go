package payment_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionRepo interface {
	Create(ctx context.Context, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.TransactionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.TransactionsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type TransactionService interface {
	Create(ctx context.Context, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.TransactionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.TransactionsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type transaction struct {
	repo TransactionRepo
}

func NewTransactionService(i *do.Injector) (TransactionService, error) {
	return &transaction{
		repo: do.MustInvoke[TransactionRepo](i),
	}, nil
}

func (t *transaction) Create(ctx context.Context, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError) {
	return t.repo.Create(ctx, transaction)
}

func (t *transaction) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.TransactionResponse, *helpers.ApiError) {
	return t.repo.GetById(ctx, id, query)
}

func (t *transaction) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.TransactionsWithPagination, *helpers.ApiError) {
	return t.repo.GetWithPagination(ctx, skip, limit, query)
}

func (t *transaction) UpdateById(ctx context.Context, id primitive.ObjectID, transaction *dtos.CreateTransactionDto) (*models.TransactionResponse, *helpers.ApiError) {
	return t.repo.UpdateById(ctx, id, transaction)
}

func (t *transaction) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return t.repo.DeleteById(ctx, id)
}

func (t *transaction) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return t.repo.BulkDelete(ctx, ids)
}
