package watchlist_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	watchlist_dtos "larsa-e-commerce-microservice-v2/internal/modules/watchlist/dtos"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"
	watchlist_models "larsa-e-commerce-microservice-v2/internal/modules/watchlist/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WebWatchListRepo defines the interface for the watchlist repository
// This is implemented by the watchlist_repos.WebWatchListRepo struct
type WebWatchListRepo interface {
	// Core CRUD operations
	Create(ctx context.Context, dto *watchlist_dtos.CreateWatchListDTO) (*watchlist_models.WatchListResponse, *helpers.ApiError)
	// GetById removed - we query by itemId and itemType instead
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*watchlist_models.WatchListsWithPagination, *helpers.ApiError)

	// Filter-based operations
	DeleteByFilter(ctx context.Context, filter bson.M) *helpers.ApiError

	// Item and type specific operations
	GetByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) (*watchlist_models.WatchListResponse, *helpers.ApiError)
	DeleteByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) *helpers.ApiError
}

// WebWatchListService defines the public interface for the watchlist service
type WebWatchListService interface {
	// Public API operations
	Create(ctx context.Context, dto *watchlist_dtos.CreateWatchListDTO) (*watchlist_models.WatchListResponse, *helpers.ApiError)
	// GetById removed - we query by itemId and itemType instead
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*watchlist_models.WatchListsWithPagination, *helpers.ApiError)
	GetByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) (*watchlist_models.WatchListResponse, *helpers.ApiError)
	DeleteByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) *helpers.ApiError
}

// webWatchList implements the WatchListService interface
type webWatchList struct {
	repo WebWatchListRepo
}

func NewWebWatchListService(i *do.Injector) (WebWatchListService, error) {
	return &webWatchList{
		repo: do.MustInvoke[WebWatchListRepo](i),
	}, nil
}

func (w *webWatchList) Create(ctx context.Context, dto *watchlist_dtos.CreateWatchListDTO) (*watchlist_models.WatchListResponse, *helpers.ApiError) {
	return w.repo.Create(ctx, dto)
}

// GetById removed - we query by itemId and itemType instead

func (w *webWatchList) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*watchlist_models.WatchListsWithPagination, *helpers.ApiError) {
	return w.repo.GetWithPagination(ctx, skip, limit, query)
}

func (w *webWatchList) GetByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) (*watchlist_models.WatchListResponse, *helpers.ApiError) {
	return w.repo.GetByItemIdAndType(ctx, itemId, itemType)
}

func (w *webWatchList) DeleteByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) *helpers.ApiError {
	// Simply delegate to the repository method that already has the filter logic
	return w.repo.DeleteByItemIdAndType(ctx, itemId, itemType)
}
