package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseReview struct {
	Rate      uint8              `bson:"rate" json:"rate"`
	Comment   string             `bson:"comment" json:"comment"`
	ProductId primitive.ObjectID `bson:"productId" json:"productId"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
}
type Review struct {
	base_models.MetaDataBaseModel `bson:",inline"`
	BaseReview                    `bson:",inline"`
}

type ReviewResponse struct {
	base_models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseReview                            `bson:",inline"`
}

type ReviewsWithPagination = types.DataWithPagination[ReviewResponse]
