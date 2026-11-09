package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dto "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
)

func ConvertToCountryModel(dto *config_dto.CreateCountryDto, additionalOptions types.AdditionalOptions) config_location_models.Country {
	country := config_location_models.Country{
		BaseCountry: config_location_models.BaseCountry{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Name:              dto.Name,
			IsActive:          dto.IsActive,
			Iso3:              dto.Iso3,
		},
	}

	return country
}
