package watchlist_handlers

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	watchlist_dtos "larsa-e-commerce-microservice-v2/internal/modules/watchlist/dtos"
	watchlist_filters "larsa-e-commerce-microservice-v2/internal/modules/watchlist/filters"
	watchlist_services "larsa-e-commerce-microservice-v2/internal/modules/watchlist/services"
	watchlist_utils "larsa-e-commerce-microservice-v2/internal/modules/watchlist/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebWatchListHandler struct {
	webWatchListService watchlist_services.WebWatchListService
	validationInstance  *validator.Validate
}

func NewWebWatchListHandler(i *do.Injector, r *chi.Mux) {
	h := WebWatchListHandler{
		webWatchListService: do.MustInvoke[watchlist_services.WebWatchListService](i),
		validationInstance:  do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.WATCH_LIST.GetWebRoute(), func(r chi.Router) {
		// Dynamic route for any item type
		r.Route("/{itemType}", func(r chi.Router) {
			// Add item to watchlist
			r.With(
				middleware.AdminAuth(middleware.Restriction{BuyerCapability: enums.CapabilitiesRegistry.General}),
				middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
				watchlist_utils.ValidateWatchListItem(i),
			).Post("/{id}", helpers.Make(h.AddItemToWatchList))

			// Remove item from watchlist
			r.With(
				middleware.AdminAuth(middleware.Restriction{BuyerCapability: enums.CapabilitiesRegistry.General}),
				middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			).Delete("/{id}", h.RemoveItemFromWatchList)

			// Get items of a specific type in watchlist
			r.With(
				middleware.AdminAuth(middleware.Restriction{BuyerCapability: enums.CapabilitiesRegistry.General}),
				middleware.ValidateDTO(watchlist_filters.WatchListFilters{}, middleware.QUERY, i),
			).Get("/", h.GetItemTypeWatchList)
		})

		// Root endpoint for all watchlist items
		r.With(
			middleware.AdminAuth(middleware.Restriction{BuyerCapability: enums.CapabilitiesRegistry.General}),
			middleware.ValidateDTO(watchlist_filters.WatchListFilters{}, middleware.QUERY, i),
		).Get("/", h.GetAllWatchList)
	})
}

// parseItemTypeFromURL function has been moved to watchlist_utils.ParseItemTypeFromURL

// AddItemToWatchList adds an item to the user's watchlist
func (h *WebWatchListHandler) AddItemToWatchList(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Extract and validate the item type using the utility function
	itemType, err := watchlist_utils.ParseItemTypeFromURL(r)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return nil
	}

	// Get item ID from URL param
	id := chi.URLParam(r, "id")
	itemId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid item ID format",
		})
		return nil
	}

	dto := &watchlist_dtos.CreateWatchListDTO{
		ItemId:   itemId,
		ItemType: itemType,
	}

	//validation must be here

	// Add to watchlist
	result, createErr := h.webWatchListService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	response, _ := json.MarshalContext(ctx, *result)
	w.Write(response)
	return nil
}

// RemoveItemFromWatchList removes an item from the user's watchlist
func (h *WebWatchListHandler) RemoveItemFromWatchList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Extract and validate the item type using the utility function
	itemType, err := watchlist_utils.ParseItemTypeFromURL(r)
	if err != nil {
		helpers.HandleApiError(w, helpers.DefaultError(err))
		return
	}

	// Get item ID from URL param
	id := chi.URLParam(r, "id")
	itemId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid item ID format",
		})
		return
	}

	// Delete from watchlist
	if deleteErr := h.webWatchListService.DeleteByItemIdAndType(ctx, itemId, itemType); deleteErr != nil {
		helpers.HandleApiError(w, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetItemTypeWatchList retrieves watchlist items of a specific type
func (h *WebWatchListHandler) GetItemTypeWatchList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	// Extract and validate the item type using the utility function
	itemType, err := watchlist_utils.ParseItemTypeFromURL(r)
	if err != nil {
		helpers.HandleApiError(w, helpers.DefaultError(err))
		return
	}

	query := r.URL.Query().Get("query")

	var filters watchlist_filters.WatchListFilters
	if query != "" {
		if err := json.Unmarshal([]byte(query), &filters); err != nil {
			helpers.HandleApiError(w, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid query format",
			})
			return
		}
	}

	// Apply the specified item type filter
	filters.ItemType = &itemType

	filterBytes, _ := json.Marshal(filters)
	modifiedQuery := string(filterBytes)

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webWatchListService.GetWithPagination(ctx, skip, limit, modifiedQuery)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}

// GetAllWatchList retrieves all watchlist items for the current user
func (h *WebWatchListHandler) GetAllWatchList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)
	query := r.URL.Query().Get("query")

	skip, limit, errGetPaginate := util.Paginate(r)
	if errGetPaginate != nil {
		helpers.DefaultError(errGetPaginate)
		return
	}

	result, getErr := h.webWatchListService.GetWithPagination(ctx, skip, limit, query)
	if getErr != nil {
		helpers.HandleApiError(w, getErr)
		return
	}

	response, _ := json.MarshalContext(ctx, result)
	w.Write(response)
}
