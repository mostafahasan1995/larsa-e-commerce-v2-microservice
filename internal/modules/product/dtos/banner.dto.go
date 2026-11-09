package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBannerItemDto struct {
	Id   primitive.ObjectID           `json:"id" validate:"required"`
	Type product_enums.BannerItemType `json:"type" validate:"required,oneof=product"`
}

type CreateBannerDto struct {
	Name         types.Localizable[string]                   `json:"name" validate:"required,localizable"`
	Description  *types.Localizable[string]                  `json:"description" validate:"required,localizable"`
	PrimaryImage types.Localizable[types.FileFieldDto]       `json:"primaryImage" validate:"required"`
	Images       *types.LocalizableArray[types.FileFieldDto] `json:"images" validate:"omitempty"`
	Url          *string                                     `json:"url" validate:"omitempty,url"`
	IsActive     bool                                        `json:"isActive" validate:"omitempty"`
	Item         *CreateBannerItemDto                        `json:"item" validate:"omitempty"`
}

func (dto *CreateBannerDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
	if dto.Item != nil {
		if dto.Item.Type == product_enums.BANNER_ITEM_TYPE_PRODUCT {
			product, err := productRepo.GetById(args.Context, dto.Item.Id, []bson.M{})
			if err != nil {
				return nil, err
			}

			if product == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Product not found",
				}
			}

		}
	}
	return nil, nil
}
