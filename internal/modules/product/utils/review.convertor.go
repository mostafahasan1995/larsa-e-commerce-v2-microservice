package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToReviewModel(productId primitive.ObjectID, createReviewDto dtos.CreateReviewDto, additionalOptions types.AdditionalOptions) models.Review {

	review := models.Review{
		BaseReview: models.BaseReview{
			ProductId: productId,
			UserId:    additionalOptions.UserId,
			Rate:      createReviewDto.Rate,
			Comment:   createReviewDto.Comment,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}

	return review
}
