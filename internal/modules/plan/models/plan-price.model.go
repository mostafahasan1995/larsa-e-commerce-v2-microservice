package plan_models

import (
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasePlanPrice struct {
	CurrencyId primitive.ObjectID `bson:"currencyId" json:"currencyId"`
	Price      float64            `bson:"price" json:"price"`
}

type PlanPrice struct {
	BasePlanPrice `bson:",inline"`
}

type PlanPriceResponse struct {
	BasePlanPrice `bson:",inline"`
	Currency      *config_models.Currency `bson:"currencyObject" json:"currencyObject"`
}
