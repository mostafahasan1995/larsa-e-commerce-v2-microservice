package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttributeRepo interface {
	Create(ctx context.Context, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AttributeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AttributeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AttributesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type AttributeService interface {
	Create(ctx context.Context, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AttributeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AttributeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AttributesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type attribute struct {
	repo           AttributeRepo
	productService ProductService
}

func NewAttributeService(i *do.Injector) (AttributeService, error) {
	return &attribute{
		repo:           do.MustInvoke[AttributeRepo](i),
		productService: do.MustInvoke[ProductService](i),
	}, nil
}

func (p *attribute) Create(ctx context.Context, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, attribute)
}

func (p *attribute) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AttributeResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *attribute) GetAll(ctx context.Context, query string) ([]models.AttributeResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *attribute) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AttributesWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *attribute) UpdateById(ctx context.Context, id primitive.ObjectID, attribute *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError) {
	stored, getErr := p.repo.GetById(ctx, id, "")
	if getErr != nil {
		return nil, getErr
	}

	updated, updateErr := p.repo.UpdateById(ctx, id, attribute)
	if updateErr != nil {
		return nil, updateErr
	}

	if stored != nil && updated != nil && stored.Code != updated.Code {
		if err := p.productService.UpdateAttributeCode(ctx, stored.Code, updated.Code); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (p *attribute) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *attribute) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
