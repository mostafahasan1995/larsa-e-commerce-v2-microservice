package config_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CurrencyRepo interface {
	Create(ctx context.Context, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CurrencyResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CurrencyResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CurrenciesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type CurrencyService interface {
	Create(ctx context.Context, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CurrencyResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CurrencyResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CurrenciesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type currency struct {
	repo CurrencyRepo
}

func NewCurrencyService(i *do.Injector) (CurrencyService, error) {
	return &currency{
		repo: do.MustInvoke[CurrencyRepo](i),
	}, nil
}

func (p *currency) Create(ctx context.Context, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, currency)
}

func (p *currency) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CurrencyResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *currency) GetAll(ctx context.Context, query string) ([]models.CurrencyResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *currency) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CurrenciesWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *currency) UpdateById(ctx context.Context, id primitive.ObjectID, currency *dtos.CreateCurrencyDto) (*models.CurrencyResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, currency)
}

func (p *currency) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *currency) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
