package config_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	config_enums "larsa-e-commerce-microservice-v2/internal/modules/config/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SEO represents the SEO-related fields for a category
type SEODto struct {
	MetaTitle       *types.Localizable[string]      `json:"metaTitle" validate:"omitempty"`
	MetaDescription *types.Localizable[string]      `json:"metaDescription" validate:"omitempty"`
	MetaKeywords    *types.LocalizableArray[string] `json:"metaKeywords" validate:"omitempty"`
	OgTitle         *types.Localizable[string]      `json:"ogTitle" validate:"omitempty"`
	OgDescription   *types.Localizable[string]      `json:"ogDescription" validate:"omitempty"`
}

type CreateCategoryDto struct {
	Name         types.Localizable[string]              `json:"name" validate:"required"`
	Description  *types.Localizable[string]             `json:"description,omitempty" validate:"omitempty"`
	Codes        *types.LocalizableArray[string]        `json:"codes" validate:"omitempty,dive,required"`
	ParentId     *primitive.ObjectID                    `json:"parentId,omitempty" validate:"omitempty"`
	Position     int                                    `json:"position" validate:"required"`
	Visibilities []config_enums.CategoryVisibility      `json:"visibilities" validate:"required,min=1,dive,oneof=categories_list search_result categories_carousel invisible"`
	Image        types.Localizable[types.FileFieldDto]  `json:"image,omitempty" validate:"required"`
	Icon         *types.Localizable[types.FileFieldDto] `json:"icon,omitempty" validate:"omitempty"`
	SEO          *SEODto                                `json:"seo,omitempty" validate:"omitempty"`
}
