package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
)

func ConvertToAreaModel(dto config_dtos.CreateAreaDto, additionalOptions types.AdditionalOptions) config_location_models.Area {
	area := config_location_models.Area{
		BaseArea: config_location_models.BaseArea{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Name:              dto.Name,
			IsActive:          dto.IsActive,
			CityId:            dto.CityId,
		},
	}

	return area
}
