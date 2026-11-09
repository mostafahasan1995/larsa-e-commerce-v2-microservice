package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
)

type WebBannerRepo interface {
	GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError)
}

type WebBannerService interface {
	GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError)
}

type webBanner struct {
	repo WebBannerRepo
}

func NewWebBannerService(i *do.Injector) (WebBannerService, error) {
	return &webBanner{
		repo: do.MustInvoke[WebBannerRepo](i),
	}, nil
}

func (p *webBanner) GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}
