package payment_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceRepo interface {
	Create(ctx context.Context, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.InvoiceResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.InvoicesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type InvoiceService interface {
	Create(ctx context.Context, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.InvoiceResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.InvoicesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type invoice struct {
	repo InvoiceRepo
}

func NewInvoiceService(i *do.Injector) (InvoiceService, error) {
	return &invoice{
		repo: do.MustInvoke[InvoiceRepo](i),
	}, nil
}

func (i *invoice) Create(ctx context.Context, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError) {
	return i.repo.Create(ctx, invoice)
}

func (i *invoice) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.InvoiceResponse, *helpers.ApiError) {
	return i.repo.GetById(ctx, id, query)
}

func (i *invoice) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.InvoicesWithPagination, *helpers.ApiError) {
	return i.repo.GetWithPagination(ctx, skip, limit, query)
}

func (i *invoice) UpdateById(ctx context.Context, id primitive.ObjectID, invoice *dtos.CreateInvoiceDto) (*models.InvoiceResponse, *helpers.ApiError) {
	return i.repo.UpdateById(ctx, id, invoice)
}

func (i *invoice) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return i.repo.DeleteById(ctx, id)
}

func (i *invoice) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return i.repo.BulkDelete(ctx, ids)
}
