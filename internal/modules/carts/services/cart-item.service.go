package cart_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItemRepo interface {
	Create(ctx context.Context, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError
	UpdateById(ctx context.Context, cartItemId primitive.ObjectID, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError
	Clear(ctx context.Context) *helpers.ApiError
}

type CartItemService interface {
	Create(ctx context.Context, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError
	CreateBulk(ctx context.Context, bulkUpsertCartItemDto *cart_dtos.BulkUpsertCartItemDto) *helpers.ApiError
	UpdateById(ctx context.Context, cartItemId primitive.ObjectID, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError
	Clear(ctx context.Context) *helpers.ApiError
}

type cartItem struct {
	repo CartItemRepo
}

func NewCartItemService(i *do.Injector) (CartItemService, error) {
	return &cartItem{
		repo: do.MustInvoke[CartItemRepo](i),
	}, nil
}

func (p *cartItem) Create(ctx context.Context, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError {
	return p.repo.Create(ctx, upsertCartItemDto)
}

func (p *cartItem) CreateBulk(ctx context.Context, bulkUpsertCartItemDto *cart_dtos.BulkUpsertCartItemDto) *helpers.ApiError {
	for _, item := range bulkUpsertCartItemDto.Items {
		err := p.Create(ctx, &item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *cartItem) UpdateById(ctx context.Context, cartItemId primitive.ObjectID, upsertCartItemDto *cart_dtos.UpsertCartItemDto) *helpers.ApiError {
	return p.repo.UpdateById(ctx, cartItemId, upsertCartItemDto)
}

func (p *cartItem) Clear(ctx context.Context) *helpers.ApiError {
	return p.repo.Clear(ctx)
}
