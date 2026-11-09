package quotation_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
)

type ShippingEmail struct {
	Subject string `bson:"subject" json:"subject"`
	Body    string `bson:"body" json:"body"`
}

type Cost struct {
	Amount       float64 `bson:"amount" json:"amount"`
	CurrencyCode string  `bson:"currencyCode" json:"currencyCode"`
}

type PlatformCost struct {
	Amount       float64                          `bson:"amount" json:"amount"`
	CurrencyCode string                           `bson:"currencyCode" json:"currencyCode"`
	Type         quotation_enums.PlatformCostType `bson:"type" json:"type"`
}

type ShippingHistory struct {
	Status    quotation_enums.ShippingStatus `bson:"status" json:"status"`
	Timestamp time.Time                      `bson:"timestamp" json:"timestamp"`
	UserId    primitive.ObjectID             `bson:"userId" json:"userId"`
}

type ShippingAgent struct {
	Method           *quotation_enums.ShippingMethod `bson:"method,omitempty" json:"method,omitempty"`
	CountryIso3      string                          `bson:"countryIso3,omitempty" json:"countryIso3,omitempty"`
	AgentId          primitive.ObjectID              `bson:"agentId" json:"agentId"`
	CustomsClearance *Cost                           `bson:"customsClearance,omitempty" json:"customsClearance,omitempty"`
	Logistic         *Cost                           `bson:"logistic,omitempty" json:"logistic,omitempty"`
	Insurance        *Cost                           `bson:"insurance,omitempty" json:"insurance,omitempty"`
}

type BaseShipping struct {
	models.MetaDataBaseModel `bson:",inline"`

	Type              quotation_enums.ShippingType   `bson:"type" json:"type"`
	From              ShippingAgent                  `bson:"from" json:"from"`
	To                ShippingAgent                  `bson:"to" json:"to"`
	Services          []string                       `bson:"services" json:"services"`
	Email             *ShippingEmail                 `bson:"email,omitempty" json:"email,omitempty"`
	QuotationId       primitive.ObjectID             `bson:"quotationId" json:"quotationId"`
	ShipmentCompanyId *primitive.ObjectID            `bson:"shipmentCompanyId,omitempty" json:"shipmentCompanyId,omitempty"`
	Platform          *PlatformCost                  `bson:"platform,omitempty" json:"platform,omitempty"`
	Status            quotation_enums.ShippingStatus `bson:"status" json:"status"`
	Histories         []ShippingHistory              `bson:"histories" json:"histories"`
}

type Shipping struct {
	BaseShipping `bson:",inline"`
}

type ShippingResponse struct {
	BaseShipping `bson:",inline"`
}
