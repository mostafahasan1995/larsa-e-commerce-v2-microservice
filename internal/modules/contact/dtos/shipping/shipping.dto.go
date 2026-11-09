package contact_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"time"
)

type CreateBasicInformationDto struct {
	Name                         types.Localizable[string] `json:"name" validate:"required,localizable"`
	Type                         types.Localizable[string] `json:"type" validate:"required,localizable"`
	CommercialRegistrationNumber int64                     `json:"commercialRegistrationNumber" validate:"required"`
	TaxIdentificationNumber      int64                     `json:"taxIdentificationNumber" validate:"required"`
}

type CreateShippingContractDto struct {
	Email             string   `json:"email" validate:"required,email"`
	HeadOfficeAddress string   `json:"headOfficeAddress" validate:"required"`
	Phone             string   `json:"phone" validate:"required"`
	CountryCodes      []string `json:"countryCodes" validate:"required"`
	Website           string   `json:"website" validate:"required"`
}

type CreateShippingServiceDto struct {
	LandFreight                bool `json:"landFreight"`
	AirFreight                 bool `json:"airFreight"`
	SeaFreight                 bool `json:"seaFreight"`
	WarehousingAndDistribution bool `json:"warehousingAndDistribution"`
}

type CreateShippingFinancialDto struct {
	StartDate       time.Time            `json:"startDate" validate:"required"`
	EndDate         time.Time            `json:"endDate" validate:"required"`
	BankInformation string               `json:"bankInformation" validate:"required"`
	Attachments     []types.FileFieldDto `json:"attachments" validate:"required,min=1,dive"`
}

type CreateShippingDto struct {
	BasicInformation CreateBasicInformationDto  `json:"basicInformation" validate:"required"`
	Service          CreateShippingServiceDto   `json:"service" validate:"required"`
	Contract         CreateShippingContractDto  `json:"contract" validate:"required"`
	Financial        CreateShippingFinancialDto `json:"financial" validate:"required"`
}
