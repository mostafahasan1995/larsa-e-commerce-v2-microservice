package unit_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/unit/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitRepo interface {
	Create(ctx context.Context, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type UnitService interface {
	Create(ctx context.Context, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type unit struct {
	repo UnitRepo
}

func NewUnitService(i *do.Injector) (UnitService, error) {
	return &unit{
		repo: do.MustInvoke[UnitRepo](i),
	}, nil
}

func (p *unit) Create(ctx context.Context, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, unit)
}

func (p *unit) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.UnitResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *unit) GetAll(ctx context.Context, query string) ([]models.UnitResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *unit) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.UnitsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *unit) UpdateById(ctx context.Context, id primitive.ObjectID, unit *dtos.CreateUnitDto) (*models.UnitResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, unit)
}

func (p *unit) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *unit) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
