package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"

	"github.com/samber/do"
)

type DiscountRepo interface {
	HandleNewProductsDiscount(ctx context.Context, handleNewProductsDiscountDto dtos.HandleNewProductsDiscountDto) *helpers.ApiError
}

type DiscountService interface {
	HandleNewProductsDiscount(ctx context.Context, handleNewProductsDiscountDto dtos.HandleNewProductsDiscountDto) *helpers.ApiError
}

type discount struct {
	repo DiscountRepo
}

func NewDiscountService(i *do.Injector) (DiscountService, error) {
	return &discount{
		repo: do.MustInvoke[DiscountRepo](i),
	}, nil
}

func (d *discount) HandleNewProductsDiscount(ctx context.Context, handleNewProductsDiscountDto dtos.HandleNewProductsDiscountDto) *helpers.ApiError {
	return d.repo.HandleNewProductsDiscount(ctx, handleNewProductsDiscountDto)
}
