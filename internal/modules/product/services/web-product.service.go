package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"

	"github.com/samber/do"
)

type WebProductRepo interface {
	GetById(ctx context.Context, args product_utils.GetByIdArgs) (*models.ProductResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, args product_utils.GetWithPaginationArgs) (*models.ProductsWithPagination, *helpers.ApiError)
}

type WebProductService interface {
	GetById(ctx context.Context, args product_utils.GetByIdArgs) (*models.ProductResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, args product_utils.GetWithPaginationArgs) (*models.ProductsWithPagination, *helpers.ApiError)
}

type webProduct struct {
	repo WebProductRepo
}

func NewWebProductService(i *do.Injector) (WebProductService, error) {
	return &webProduct{
		repo: do.MustInvoke[WebProductRepo](i),
	}, nil
}

func (p *webProduct) GetById(ctx context.Context, args product_utils.GetByIdArgs) (*models.ProductResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, args)
}

func (p *webProduct) GetWithPagination(ctx context.Context, args product_utils.GetWithPaginationArgs) (*models.ProductsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, args)
}
