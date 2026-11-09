package unit_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebUnitRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError)
}

type WebUnitService interface {
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError)
}

type webUnit struct {
	repo WebUnitRepo
}

func NewWebUnitService(i *do.Injector) (WebUnitService, error) {
	return &webUnit{
		repo: do.MustInvoke[WebUnitRepo](i),
	}, nil
}

func (p *webUnit) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *webUnit) GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *webUnit) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}
