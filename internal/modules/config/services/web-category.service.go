package config_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebCategoryRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError)
}

type WebCategoryService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError)
}

type webCategory struct {
	repo WebCategoryRepo
}

func NewWebCategoryService(i *do.Injector) (WebCategoryService, error) {
	return &webCategory{
		repo: do.MustInvoke[WebCategoryRepo](i),
	}, nil
}

func (p *webCategory) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *webCategory) GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *webCategory) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}
