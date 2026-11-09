package contact_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BasicInformation struct {
	Name                         types.Localizable[string] `json:"name" bson:"name"`
	Type                         types.Localizable[string] `json:"type" bson:"type"`
	CommercialRegistrationNumber int64                     `json:"commercialRegistrationNumber" bson:"commercialRegistrationNumber"`
	TaxIdentificationNumber      int64                     `json:"taxIdentificationNumber" bson:"taxIdentificationNumber"`
}
type ShippingContract struct {
	Email             string   `json:"email" bson:"email"`
	HeadOfficeAddress string   `json:"headOfficeAddress" bson:"headOfficeAddress"`
	Phone             string   `json:"phone" bson:"phone"`
	CountryCodes      []string `json:"countryCodes" bson:"countryCodes"`
	Website           string   `json:"website" bson:"website"`
}

type ShippingService struct {
	LandFreight                bool `json:"landFreight" bson:"landFreight"`
	AirFreight                 bool `json:"airFreight" bson:"airFreight"`
	SeaFreight                 bool `json:"seaFreight" bson:"seaFreight"`
	WarehousingAndDistribution bool `json:"warehousingAndDistribution" bson:"warehousingAndDistribution"`
}

type ShippingFinancial struct {
	StartDate       time.Time            `json:"startDate" bson:"startDate"`
	EndDate         time.Time            `json:"endDate" bson:"endDate"`
	BankInformation string               `json:"bankInformation" bson:"bankInformation"`
	Attachments     []types.FileFieldDto `json:"attachments" bson:"attachments"`
}

type BaseShipping struct {
	BasicInformation BasicInformation  `json:"basicInformation" bson:"basicInformation"`
	Service          ShippingService   `json:"service" bson:"service"`
	Contract         ShippingContract  `json:"contract" bson:"contract"`
	Financial        ShippingFinancial `json:"financial" bson:"financial"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Shipping struct {
	BaseShipping `bson:",inline"`
}

type ShippingResponse struct {
	BaseShipping `bson:",inline"`
}

type ShippingsWithPagination = types.DataWithPagination[ShippingResponse]
