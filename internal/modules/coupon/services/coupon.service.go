package coupon_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponRepo interface {
	Create(ctx context.Context, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CouponResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CouponResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CouponsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type CouponService interface {
	Create(ctx context.Context, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CouponResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.CouponResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CouponsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type coupon struct {
	repo CouponRepo
}

func NewCouponService(i *do.Injector) (CouponService, error) {
	return &coupon{
		repo: do.MustInvoke[CouponRepo](i),
	}, nil
}

func (p *coupon) Create(ctx context.Context, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, coupon)
}

func (p *coupon) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.CouponResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *coupon) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.CouponsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *coupon) GetAll(ctx context.Context, query string) ([]models.CouponResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *coupon) UpdateById(ctx context.Context, id primitive.ObjectID, coupon *dtos.CreateCouponDto) (*models.CouponResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, coupon)
}

func (p *coupon) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *coupon) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
