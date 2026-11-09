package watchlist_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseWatchList struct {
	UserId   primitive.ObjectID            `bson:"userId" json:"userId"`
	ItemId   primitive.ObjectID            `bson:"itemId" json:"itemId"`
	ItemType watchlist_enums.WatchItemType `bson:"itemType" json:"itemType"`

	// Embedded metadata model
	models.MetaDataBaseModel `bson:",inline"`
}

type WatchList struct {
	BaseWatchList `bson:",inline"`
}

type WatchListResponse struct {
	BaseWatchList `bson:",inline"`
	Product       *product_models.ProductResponse `bson:"product,omitempty" json:"product,omitempty"`
}

type WatchListsWithPagination = types.DataWithPagination[WatchListResponse]
