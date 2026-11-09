package product_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BannerItem struct {
	Id   primitive.ObjectID           `bson:"id" json:"id"`
	Type product_enums.BannerItemType `bson:"type" json:"type"`
}

type BaseBanner struct {
	Name         *types.Localizable[string]                  `bson:"name,omitempty" json:"name,omitempty"`
	Description  *types.Localizable[string]                  `bson:"description,omitempty" json:"description,omitempty"`
	PrimaryImage *types.Localizable[types.FileFieldDto]      `bson:"primaryImage,omitempty" json:"primaryImage,omitempty"`
	Images       *types.LocalizableArray[types.FileFieldDto] `bson:"images,omitempty" json:"images,omitempty"`
	Url          *string                                     `bson:"url,omitempty" json:"url,omitempty"`
	IsActive     bool                                        `bson:"isActive" json:"isActive"`
	Item         *BannerItem                                 `bson:"item,omitempty" json:"item,omitempty"`
}

type Banner struct {
	models.MetaDataBaseModel `bson:",inline"`
	BaseBanner               `bson:",inline"`
}

type BannerResponse struct {
	models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseBanner                       `bson:",inline"`
	Product                          *ProductResponse `bson:"product,omitempty" json:"product,omitempty"`
}
