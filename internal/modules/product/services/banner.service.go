package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BannerRepo interface {
	Create(ctx context.Context, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BannerResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type BannerService interface {
	Create(ctx context.Context, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BannerResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type banner struct {
	repo BannerRepo
}

func NewBannerService(i *do.Injector) (BannerService, error) {
	return &banner{
		repo: do.MustInvoke[BannerRepo](i),
	}, nil
}

func (p *banner) Create(ctx context.Context, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, banner)
}

func (p *banner) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BannerResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *banner) GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *banner) UpdateById(ctx context.Context, id primitive.ObjectID, banner *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, banner)
}

func (p *banner) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *banner) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
