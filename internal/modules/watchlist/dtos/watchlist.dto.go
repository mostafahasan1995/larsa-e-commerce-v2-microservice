package watchlist_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateWatchListDTO struct {
	ItemId   primitive.ObjectID
	ItemType watchlist_enums.WatchItemType
}

func (dto *CreateWatchListDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	switch dto.ItemType {
	case watchlist_enums.PRODUCT:
		productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
		product, apiErr := productRepo.GetById(args.Context, dto.ItemId, []bson.M{})
		if apiErr != nil || product == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Product not found",
			}
		}

	default:
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Unsupported item type",
		}
	}

	return nil, nil
}
