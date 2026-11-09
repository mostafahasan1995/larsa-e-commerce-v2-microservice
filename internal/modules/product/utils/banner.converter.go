package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

func ConvertToBannerModel(dto *dtos.CreateBannerDto, additionalOptions types.AdditionalOptions) models.Banner {

	banner := models.Banner{
		BaseBanner: models.BaseBanner{
			Name:         &dto.Name,
			Description:  dto.Description,
			PrimaryImage: &dto.PrimaryImage,
			Images:       dto.Images,
			Url:          dto.Url,
			IsActive:     dto.IsActive,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}

	if dto.Item != nil {
		banner.Item = &models.BannerItem{
			Id:   dto.Item.Id,
			Type: dto.Item.Type,
		}
	}

	return banner
}
