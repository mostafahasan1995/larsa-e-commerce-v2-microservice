package config_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dto "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

func ConvertToCurrencyModel(dto *config_dto.CreateCurrencyDto, additionalOptions types.AdditionalOptions) config_models.Currency {

	currency := config_models.Currency{
		BaseCurrency: config_models.BaseCurrency{

			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Name:              &dto.Name,
			Symbol:            &dto.Symbol,
		},
	}

	return currency
}
