package config_utils

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_enums "larsa-e-commerce-microservice-v2/internal/modules/config/enums"
	models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

func convertSEODtoToSEO(seoDto *dtos.SEODto) *models.SEO {
	if seoDto == nil {
		return nil
	}

	return &models.SEO{
		MetaTitle:       seoDto.MetaTitle,
		MetaDescription: seoDto.MetaDescription,
		MetaKeywords:    seoDto.MetaKeywords,
		OgTitle:         seoDto.OgTitle,
		OgDescription:   seoDto.OgDescription,
	}
}

// ConvertToCategoryModel converts a CreateCategoryDto to a Category model
func ConvertToCategoryModel(createCategoryDto *dtos.CreateCategoryDto, additionalOptions types.AdditionalOptions) *models.Category {

	now := time.Now()

	status := models.Status{
		Name:        config_enums.APPROVED,
		TriggeredBy: &additionalOptions.UserId,
		Timestamp:   now,
	}

	category := &models.Category{
		BaseCategory: models.BaseCategory{
			Name:         createCategoryDto.Name,
			Description:  createCategoryDto.Description,
			Codes:        createCategoryDto.Codes,
			ParentId:     createCategoryDto.ParentId,
			Position:     createCategoryDto.Position,
			Visibilities: createCategoryDto.Visibilities,
			Image:        &createCategoryDto.Image,
			Icon:         createCategoryDto.Icon,
			SEO:          convertSEODtoToSEO(createCategoryDto.SEO),
			Status:       &status,

			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
	return category
}
