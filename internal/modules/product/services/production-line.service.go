package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductionLineRepo interface {
	Create(ctx context.Context, productionLine *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductionLineResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductionLinesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, productionLine *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type ProductionLineService interface {
	Create(ctx context.Context, productionLine *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductionLineResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductionLinesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, productionLine *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type productionLine struct {
	repo ProductionLineRepo
}

func NewProductionLineService(i *do.Injector) (ProductionLineService, error) {
	return &productionLine{
		repo: do.MustInvoke[ProductionLineRepo](i),
	}, nil
}

func (p *productionLine) Create(ctx context.Context, productionLine *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, productionLine)
}

func (p *productionLine) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductionLineResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *productionLine) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductionLinesWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *productionLine) UpdateById(ctx context.Context, id primitive.ObjectID, product *dtos.CreateProductionLineDto) (*models.ProductionLineResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, product)
}

func (p *productionLine) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *productionLine) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
