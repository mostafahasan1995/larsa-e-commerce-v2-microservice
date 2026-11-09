package config_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_enums "larsa-e-commerce-microservice-v2/internal/modules/config/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SEO represents the SEO-related fields for a category
type SEO struct {
	MetaTitle       *types.Localizable[string]      `bson:"metaTitle,omitempty" json:"metaTitle,omitempty"`
	MetaDescription *types.Localizable[string]      `bson:"metaDescription,omitempty" json:"metaDescription,omitempty"`
	MetaKeywords    *types.LocalizableArray[string] `bson:"metaKeywords,omitempty" json:"metaKeywords,omitempty"`
	OgTitle         *types.Localizable[string]      `bson:"ogTitle,omitempty" json:"ogTitle,omitempty"`
	OgDescription   *types.Localizable[string]      `bson:"ogDescription,omitempty" json:"ogDescription,omitempty"`
}

// Status represents the approval status of a category
type Status struct {
	Name        config_enums.StatusName `bson:"name" json:"name"`
	TriggeredBy *primitive.ObjectID     `bson:"triggeredBy,omitempty" json:"triggeredBy,omitempty"`
	Timestamp   time.Time               `bson:"timestamp" json:"timestamp"`
}

type BaseCategory struct {
	Name         types.Localizable[string]              `bson:"name" json:"name"`
	Description  *types.Localizable[string]             `bson:"description,omitempty" json:"description,omitempty"`
	Codes        *types.LocalizableArray[string]        `bson:"codes,omitempty" json:"codes,omitempty"`
	ParentId     *primitive.ObjectID                    `bson:"parentId,omitempty" json:"parentId,omitempty"`
	Position     int                                    `bson:"position" json:"position"`
	Visibilities []config_enums.CategoryVisibility      `bson:"visibilities" json:"visibilities"`
	Image        *types.Localizable[types.FileFieldDto] `bson:"image,omitempty" json:"image,omitempty"`
	Icon         *types.Localizable[types.FileFieldDto] `bson:"icon,omitempty" json:"icon,omitempty"`
	SEO          *SEO                                   `bson:"seo,omitempty" json:"seo,omitempty"`
	Status       *Status                                `bson:"status,omitempty" json:"status,omitempty"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Category struct {
	BaseCategory `bson:",inline"`
}

type CategoryResponse struct {
	BaseCategory  `bson:",inline"`
	Parent        *CategoryResponse `bson:"parent,omitempty" json:"parent,omitempty"`
	ChildrenCount *int              `bson:"childrenCount" json:"childrenCount"`
	ProductsCount *int              `bson:"productsCount" json:"productsCount"`
}

type CategoriesWithPagination = types.DataWithPagination[CategoryResponse]
