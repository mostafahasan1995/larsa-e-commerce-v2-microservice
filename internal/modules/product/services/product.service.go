package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRepo interface {
	Create(ctx context.Context, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	GetAll(ctx context.Context, filters []bson.M) ([]models.ProductResponse, *helpers.ApiError)
	LinkToProductionLine(ctx context.Context, id primitive.ObjectID, linkProductionLineDto dtos.LinkProductionLineDto) *helpers.ApiError
	DeleteDuplicatesByName(ctx context.Context) *helpers.ApiError
	SaveInventoryProducts(ctx context.Context, saveInventoryProductsDto *dtos.SaveInventoryProductsDto) (*models.ProductResponse, *helpers.ApiError)
	UpdateAttributeCode(ctx context.Context, oldCode string, newCode string) *helpers.ApiError
}

type ProductService interface {
	Create(ctx context.Context, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	GetAll(ctx context.Context, filters []bson.M) ([]models.ProductResponse, *helpers.ApiError)
	LinkToProductionLine(ctx context.Context, id primitive.ObjectID, linkProductionLineDto dtos.LinkProductionLineDto) *helpers.ApiError
	DeleteDuplicatesByName(ctx context.Context) *helpers.ApiError
	GenerateInventoryProducts(ctx context.Context, product *models.ProductResponse) ([]models.ProductResponse, *helpers.ApiError)
	SaveInventoryProducts(ctx context.Context, saveInventoryProductsDto *dtos.SaveInventoryProductsDto) (*models.ProductResponse, *helpers.ApiError)
	UpdateAttributeCode(ctx context.Context, oldCode string, newCode string) *helpers.ApiError
}

type product struct {
	repo                    ProductRepo
	branchProductService    BranchProductService
	inventoryProductService InventoryProductService
}

func NewProductService(i *do.Injector) (ProductService, error) {
	return &product{
		repo:                    do.MustInvoke[ProductRepo](i),
		branchProductService:    do.MustInvoke[BranchProductService](i),
		inventoryProductService: do.MustInvoke[InventoryProductService](i),
	}, nil
}

func (p *product) Create(ctx context.Context, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, product)
}

func (p *product) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductResponse, *helpers.ApiError) {
	product, err := p.repo.GetById(ctx, id, query)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *product) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductsWithPagination, *helpers.ApiError) {
	result, err := p.repo.GetWithPagination(ctx, skip, limit, query)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *product) UpdateById(ctx context.Context, id primitive.ObjectID, product *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, product)
}

func (p *product) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *product) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}

func (p *product) GetAll(ctx context.Context, filters []bson.M) ([]models.ProductResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, filters)
}

func (p *product) LinkToProductionLine(ctx context.Context, id primitive.ObjectID, linkProductionLineDto dtos.LinkProductionLineDto) *helpers.ApiError {
	return p.repo.LinkToProductionLine(ctx, id, linkProductionLineDto)
}

func (p *product) DeleteDuplicatesByName(ctx context.Context) *helpers.ApiError {
	return p.repo.DeleteDuplicatesByName(ctx)
}

func (p *product) GenerateInventoryProducts(ctx context.Context, product *models.ProductResponse) ([]models.ProductResponse, *helpers.ApiError) {
	return p.inventoryProductService.GenerateInventoryProducts(ctx, product)
}

func (p *product) SaveInventoryProducts(ctx context.Context, saveInventoryProductsDto *dtos.SaveInventoryProductsDto) (*models.ProductResponse, *helpers.ApiError) {
	return p.repo.SaveInventoryProducts(ctx, saveInventoryProductsDto)
}

func (p *product) UpdateAttributeCode(ctx context.Context, oldCode string, newCode string) *helpers.ApiError {
	return p.repo.UpdateAttributeCode(ctx, oldCode, newCode)
}
