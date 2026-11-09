package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dto "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
)

func ConvertToCityModel(dto *config_dto.CreateCityDto, additionalOptions types.AdditionalOptions) config_location_models.City {
	city := config_location_models.City{
		BaseCity: config_location_models.BaseCity{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Name:              dto.Name,
			IsActive:          dto.IsActive,
			CountryId:         dto.CountryId,
		},
	}

	return city
}
