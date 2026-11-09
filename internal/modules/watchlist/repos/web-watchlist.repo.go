package watchlist_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	watchlist_dtos "larsa-e-commerce-microservice-v2/internal/modules/watchlist/dtos"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"
	watchlist_filters "larsa-e-commerce-microservice-v2/internal/modules/watchlist/filters"
	watchlist_models "larsa-e-commerce-microservice-v2/internal/modules/watchlist/models"
	watchlist_services "larsa-e-commerce-microservice-v2/internal/modules/watchlist/services"
	watchlist_utils "larsa-e-commerce-microservice-v2/internal/modules/watchlist/utils"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebWatchListRepo struct {
	*base_repos.BaseRepo[watchlist_models.WatchList, watchlist_models.WatchListResponse]
}

func NewWebWatchListRepo(i *do.Injector) (watchlist_services.WebWatchListRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebWatchListRepo{
		BaseRepo: base_repos.NewBaseRepo[watchlist_models.WatchList, watchlist_models.WatchListResponse](db, enums.WATCH_LIST),
	}, nil
}

func (repo *WebWatchListRepo) Create(ctx context.Context, dto *watchlist_dtos.CreateWatchListDTO) (*watchlist_models.WatchListResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Convert DTO to model
	watchlistItem := watchlist_utils.ConvertToWatchListModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	// Set the current user as the owner of this watchlist item
	watchlistItem.UserId = cfg.User.Id

	// Check if this item already exists in the user's watchlist
	existingItem, err := repo.GetByItemIdAndType(ctx, watchlistItem.ItemId, watchlistItem.ItemType)
	if err == nil && existingItem != nil {
		return existingItem, nil // Return the existing item if found
	}

	return repo.BaseRepo.Create(ctx, *watchlistItem)
}

func (repo *WebWatchListRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*watchlist_models.WatchListsWithPagination, *helpers.ApiError) {
	// Base match stage
	matchStage := bson.M{}

	// Add current user filter - users can only see their own watchlist
	cfg, _ := util.GetReqAppCfg(ctx)
	matchStage["userId"] = cfg.User.Id

	// Create filters from query parameters
	watchlistFilters, _ := watchlist_filters.NewWatchListFilters(query)

	// Build pipeline with filters
	pipeline := watchlistFilters.BuildPipeline(matchStage, ctx, true)

	// Get paginated results
	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &watchlistFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *WebWatchListRepo) DeleteByFilter(ctx context.Context, filter bson.M) *helpers.ApiError {
	if filter == nil {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Filter cannot be nil",
		}
	}

	// Get the collection
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return err
	}

	// Delete the items matching the filter
	deletedCount, deleteErr := collection.DeleteMany(ctx, filter)
	if deleteErr != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete watchlist items: " + deleteErr.Error(),
		}
	}
	if deletedCount.DeletedCount == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "No items found to delete",
		}
	}

	return nil
}

// GetByItemIdAndType retrieves a watchlist item by item ID and type for the current user
func (repo *WebWatchListRepo) GetByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) (*watchlist_models.WatchListResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Use WatchListFilters to build the match filter
	filters := watchlist_filters.WatchListFilters{
		UserId:   &cfg.User.Id,
		ItemId:   &itemId,
		ItemType: &itemType,
	}
	matchFilter := filters.BuildMatchFilter(ctx, bson.M{})

	// Get the collection
	collection, err := repo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	// Find the item using the generated filter
	var watchlistItem watchlist_models.WatchList
	if err := collection.FindOne(ctx, matchFilter).Decode(&watchlistItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found, return nil without error
		}
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get watchlist item: " + err.Error(),
		}
	}

	// Convert to response
	response := watchlist_models.WatchListResponse{
		BaseWatchList: watchlistItem.BaseWatchList,
	}

	return &response, nil
}

// DeleteByItemIdAndType deletes a watchlist item by item ID and type for the current user
func (repo *WebWatchListRepo) DeleteByItemIdAndType(ctx context.Context, itemId primitive.ObjectID, itemType watchlist_enums.WatchItemType) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Use WatchListFilters to build the match filter
	filters := watchlist_filters.WatchListFilters{
		UserId:   &cfg.User.Id,
		ItemId:   &itemId,
		ItemType: &itemType,
	}
	matchFilter := filters.BuildMatchFilter(ctx, bson.M{})

	// Use the generalized DeleteByFilter method with the generated filter
	return repo.DeleteByFilter(ctx, matchFilter)
}
