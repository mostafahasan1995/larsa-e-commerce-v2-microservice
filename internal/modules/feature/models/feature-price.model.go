package feature_models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
)

type FeaturePriceCurrency struct {
	CurrencyId primitive.ObjectID `bson:"currencyId" json:"currencyId"`
	Price      float64            `bson:"price" json:"price"`
}

type FeaturePrice struct {
	Active     bool                   `bson:"active" json:"active"`
	Currencies []FeaturePriceCurrency `bson:"currencies" json:"currencies"`
}

//

type FeaturePriceCurrencyResponse struct {
	CurrencyId primitive.ObjectID      `bson:"currencyId" json:"currencyId"`
	Price      float64                 `bson:"price" json:"price"`
	Currency   *config_models.Currency `bson:"currency" json:"currency"`
}

type FeaturePriceResponse struct {
	Active     bool                   `bson:"active" json:"active"`
	Currencies []FeaturePriceCurrency `bson:"currencies" json:"currencies"`
}
