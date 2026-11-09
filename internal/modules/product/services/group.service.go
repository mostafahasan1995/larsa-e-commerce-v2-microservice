package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRepo interface {
	Create(ctx context.Context, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.GroupResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.GroupResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.GroupsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type GroupService interface {
	Create(ctx context.Context, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.GroupResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.GroupResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.GroupsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type group struct {
	repo GroupRepo
}

func NewGroupService(i *do.Injector) (GroupService, error) {
	return &group{
		repo: do.MustInvoke[GroupRepo](i),
	}, nil
}

func (p *group) Create(ctx context.Context, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, group)
}

func (p *group) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.GroupResponse, *helpers.ApiError) {
	return p.repo.GetById(ctx, id, query)
}

func (p *group) GetAll(ctx context.Context, query string) ([]models.GroupResponse, *helpers.ApiError) {
	return p.repo.GetAll(ctx, query)
}

func (p *group) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.GroupsWithPagination, *helpers.ApiError) {
	return p.repo.GetWithPagination(ctx, skip, limit, query)
}

func (p *group) UpdateById(ctx context.Context, id primitive.ObjectID, group *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, group)
}

func (p *group) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *group) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
