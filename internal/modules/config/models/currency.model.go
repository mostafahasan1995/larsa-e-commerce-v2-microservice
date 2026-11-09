package config_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
)

type BaseCurrency struct {
	models.MetaDataBaseModel `bson:",inline"`
	Name                     *types.Localizable[string] `bson:"name,omitempty" json:"name,omitempty"`
	Code                     string                     `bson:"code" json:"code"`
	CountryIso3s             []string                   `bson:"countryIso3s" json:"countryIso3s"`
	Symbol                   *string                    `bson:"symbol,omitempty" json:"symbol,omitempty"`
}

type Currency struct {
	BaseCurrency `bson:",inline"`
}

type CurrencyResponse struct {
	BaseCurrency `bson:",inline"`
	Country      *config_location_models.CountryResponse `bson:"country,omitempty" json:"country,omitempty"`
}

type CurrenciesWithPagination = types.DataWithPagination[CurrencyResponse]
