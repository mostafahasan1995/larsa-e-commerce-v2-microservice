package contact_services

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/shipping"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/shipping"
)

type ShippingRepo interface {
	Create(ctx context.Context, createShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ShippingResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ShippingsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, updateShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	FindOneByCountryCode(ctx context.Context, iso3 string) (*models.ShippingResponse, *helpers.ApiError)
}

type ShippingService struct {
	repo ShippingRepo
}

func NewShippingService(i *do.Injector) (ShippingService, error) {
	repo := do.MustInvoke[ShippingRepo](i)
	return ShippingService{repo: repo}, nil
}

func (s ShippingService) Create(ctx context.Context, createShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError) {
	return s.repo.Create(ctx, createShippingDto)
}

func (s ShippingService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ShippingResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, query)
}

func (s ShippingService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ShippingsWithPagination, *helpers.ApiError) {
	return s.repo.GetWithPagination(ctx, skip, limit, query)
}

func (s ShippingService) UpdateById(ctx context.Context, id primitive.ObjectID, updateShippingDto *dtos.CreateShippingDto) (*models.ShippingResponse, *helpers.ApiError) {
	return s.repo.UpdateById(ctx, id, updateShippingDto)
}

func (s ShippingService) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id)
}

func (s ShippingService) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return s.repo.BulkDelete(ctx, ids)
}

func (s ShippingService) FindOneByCountryCode(ctx context.Context, iso3 string) (*models.ShippingResponse, *helpers.ApiError) {
	return s.repo.FindOneByCountryCode(ctx, iso3)
}
