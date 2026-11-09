package product_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddonRepo interface defines the repository operations for Addon
type AddonRepo interface {
	Create(ctx context.Context, addon *dtos.CreateAddonDto) (*models.AddonResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, addon *dtos.UpdateAddonDto) (*models.AddonResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

// AddonService interface defines the service operations for Addon
type AddonService interface {
	Create(ctx context.Context, addon *dtos.CreateAddonDto) (*models.AddonResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, addon *dtos.UpdateAddonDto) (*models.AddonResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

// addon struct implements the AddonService interface
type addon struct {
	repo AddonRepo
}

// NewAddonService creates a new instance of AddonService
func NewAddonService(i *do.Injector) (AddonService, error) {
	return &addon{
		repo: do.MustInvoke[AddonRepo](i),
	}, nil
}

// Create creates a new addon
func (a *addon) Create(ctx context.Context, addonDto *dtos.CreateAddonDto) (*models.AddonResponse, *helpers.ApiError) {
	return a.repo.Create(ctx, addonDto)
}

// GetById retrieves an addon by ID
func (a *addon) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError) {
	return a.repo.GetById(ctx, id, query)
}

// GetWithPagination retrieves paginated addons
func (a *addon) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError) {
	return a.repo.GetWithPagination(ctx, skip, limit, query)
}

// GetAll retrieves all addons
func (a *addon) GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError) {
	return a.repo.GetAll(ctx, query)
}

// UpdateById updates an addon by ID
func (a *addon) UpdateById(ctx context.Context, id primitive.ObjectID, addonDto *dtos.UpdateAddonDto) (*models.AddonResponse, *helpers.ApiError) {
	return a.repo.UpdateById(ctx, id, addonDto)
}

// DeleteById deletes an addon by ID
func (a *addon) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return a.repo.DeleteById(ctx, id)
}

// BulkDelete deletes multiple addons by IDs
func (a *addon) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return a.repo.BulkDelete(ctx, ids)
}
