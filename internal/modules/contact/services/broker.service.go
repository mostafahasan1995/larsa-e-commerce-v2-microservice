package contact_services

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models"
)

type BrokerRepo interface {
	Create(ctx context.Context, createBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BrokerResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BrokersWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, updateBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type BrokerService struct {
	repo BrokerRepo
}

func NewBrokerService(i *do.Injector) (BrokerService, error) {
	repo := do.MustInvoke[BrokerRepo](i)
	return BrokerService{repo: repo}, nil
}

func (s BrokerService) Create(ctx context.Context, createBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError) {
	return s.repo.Create(ctx, createBrokerDto)
}

func (s BrokerService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BrokerResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, query)
}

func (s BrokerService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BrokersWithPagination, *helpers.ApiError) {
	return s.repo.GetWithPagination(ctx, skip, limit, query)
}

func (s BrokerService) UpdateById(ctx context.Context, id primitive.ObjectID, updateBrokerDto *dtos.CreateBrokerDto) (*models.BrokerResponse, *helpers.ApiError) {
	return s.repo.UpdateById(ctx, id, updateBrokerDto)
}

func (s BrokerService) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id)
}

func (s BrokerService) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return s.repo.BulkDelete(ctx, ids)
}
