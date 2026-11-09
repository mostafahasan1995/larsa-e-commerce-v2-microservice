package watchlist_utils

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	watchlist_dtos "larsa-e-commerce-microservice-v2/internal/modules/watchlist/dtos"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"
	watchlist_models "larsa-e-commerce-microservice-v2/internal/modules/watchlist/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var urlPathToItemType = map[string]watchlist_enums.WatchItemType{
	"products": watchlist_enums.PRODUCT,
}

// ConvertToWatchListModel converts a CreateWatchListDTO to a WatchList model
func ConvertToWatchListModel(dto *watchlist_dtos.CreateWatchListDTO, options types.AdditionalOptions) *watchlist_models.WatchList {
	now := time.Now()

	watchList := &watchlist_models.WatchList{
		BaseWatchList: watchlist_models.BaseWatchList{
			ItemId:   dto.ItemId,
			ItemType: dto.ItemType,
			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        primitive.NewObjectID(),
				CreatedAt: now,
				CreatedBy: &options.UserId,
				UpdatedAt: now,
				UpdatedBy: &options.UserId,
			},
		},
	}

	return watchList
}

// ParseItemTypeFromURL extracts and validates the item type from the URL parameters
// It uses a mapping from URL paths (singular or plural) to enum values
func ParseItemTypeFromURL(r *http.Request) (watchlist_enums.WatchItemType, error) {
	itemTypeStr := chi.URLParam(r, "itemType")
	itemTypeLower := strings.ToLower(itemTypeStr)

	if itemType, exists := urlPathToItemType[itemTypeLower]; exists {
		return itemType, nil
	}

	// Create a list of supported types for better error messages
	var supportedTypes []string
	for path := range urlPathToItemType {
		supportedTypes = append(supportedTypes, path)
	}

	return "", fmt.Errorf("unsupported item type: '%s'. Supported types are: %s",
		itemTypeStr,
		strings.Join(supportedTypes, ", "))
}

// ValidateWatchListItem checks if an item exists and is valid to be added to the watchlist
func ValidateWatchListItem(i *do.Injector) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract item type and ID
			itemType, err := ParseItemTypeFromURL(r)
			if err != nil {
				helpers.HandleApiError(w, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				})
				return
			}

			id := chi.URLParam(r, "id")
			itemId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				helpers.HandleApiError(w, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				})
				return
			}

			// Create DTO
			dto := &watchlist_dtos.CreateWatchListDTO{
				ItemId:   itemId,
				ItemType: itemType,
			}

			// Get DB client
			db := do.MustInvoke[*mongo.Client](i)
			ctxAppCfg, _ := util.AddCtxAppCfg(r)

			// Run the CustomValidator from the DTO
			_, apiErr := dto.CustomValidator(middleware.CustomValidatorArgs{
				Context:  ctxAppCfg,
				DB:       db,
				Request:  r,
				Injector: i,
			})

			if apiErr != nil {
				helpers.HandleApiError(w, apiErr)
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}
