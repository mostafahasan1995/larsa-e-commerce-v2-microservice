package seed_data

import (
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	country_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	"time"

	"github.com/pariz/gountries"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Currencies = []country_models.Currency{}

func InitCurrencies() {

	query := gountries.New()
	currencies := []country_models.Currency{}

	for _, country := range query.FindAllCountries() {

		for _, currency := range country.Currencies {

			countries := []string{country.Alpha3}

			for _, existingCurrency := range currencies {
				if existingCurrency.Code == currency {
					currencies = util.Filter(currencies, func(el country_models.Currency, index int) bool {
						return el.Code != existingCurrency.Code
					})
					countries = append(countries, existingCurrency.CountryIso3s...)
				}
			}

			currency := country_models.Currency{
				BaseCurrency: country_models.BaseCurrency{
					MetaDataBaseModel: models.MetaDataBaseModel{
						Id:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Code:         currency,
					CountryIso3s: util.ExtractUniqueValues(countries),
				},
			}

			currencies = append(currencies, currency)
		}
	}

	Currencies = currencies
}
