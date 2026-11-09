package config_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepo interface {
	Create(ctx context.Context, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type CategoryService interface {
	Create(ctx context.Context, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type category struct {
	repo CategoryRepo
}

func NewCategoryService(i *do.Injector) (CategoryService, error) {
	return &category{
		repo: do.MustInvoke[CategoryRepo](i),
	}, nil
}

func (p *category) Create(ctx context.Context, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, category)
}

func (p *category) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CategoryResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *category) GetAll(ctx context.Context, query string) ([]models.CategoryResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *category) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CategoriesWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *category) UpdateById(ctx context.Context, id primitive.ObjectID, category *dtos.CreateCategoryDto) (*models.CategoryResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, category)
}

func (p *category) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *category) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
