package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebAddonRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError)
}

type WebAddonService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError)
}

type webAddon struct {
	repo WebAddonRepo
}

func NewWebAddonService(i *do.Injector) (WebAddonService, error) {
	return &webAddon{
		repo: do.MustInvoke[WebAddonRepo](i),
	}, nil
}

func (a *webAddon) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError) {
	return a.repo.GetById(ctx, id, query)
}

func (a *webAddon) GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError) {
	return a.repo.GetAll(ctx, query)
}

func (a *webAddon) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError) {
	return a.repo.GetWithPagination(ctx, skip, limit, query)
}
