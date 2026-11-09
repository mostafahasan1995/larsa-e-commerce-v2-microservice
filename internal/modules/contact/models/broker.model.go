package contact_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BrokerBasicInformation struct {
	FullName                     types.Localizable[string] `json:"fullName" bson:"fullName"`
	CompanyName                  types.Localizable[string] `json:"companyName" bson:"companyName"`
	CustomsLicenseNumber         string                    `json:"customsLicenseNumber" bson:"customsLicenseNumber"`
	CommercialRegistrationNumber string                    `json:"commercialRegistrationNumber" bson:"commercialRegistrationNumber"`
	TaxIdentificationNumber      string                    `json:"taxIdentificationNumber" bson:"taxIdentificationNumber"`
	AuthorizedCountry            string                    `json:"authorizedCountry" bson:"authorizedCountry"`
	YearsOfExperience            string                    `json:"yearsOfExperience" bson:"yearsOfExperience"`
}

type BrokerContract struct {
	HeadOfficeAddress string `json:"headOfficeAddress" bson:"headOfficeAddress"`
	Email             string `json:"email" bson:"email"`
	PhoneNumber       string `json:"phoneNumber" bson:"phoneNumber"`
}

type BrokerFinancial struct {
	StartDate   time.Time            `json:"startDate" bson:"startDate"`
	EndDate     time.Time            `json:"endDate" bson:"endDate"`
	Attachments []types.FileFieldDto `json:"attachments" bson:"attachments"`
}

type BankAccount struct {
	BankName          types.Localizable[string] `json:"bankName" bson:"bankName"`
	AccountNumber     string                    `json:"accountNumber" bson:"accountNumber"`
	AccountHolderName string                    `json:"accountHolderName" bson:"accountHolderName"`
	IBAN              string                    `json:"iban" bson:"iban"`
	SWIFTCode         string                    `json:"swiftCode" bson:"swiftCode"`
}

type BaseBroker struct {
	BasicInformation BrokerBasicInformation `json:"basicInformation" bson:"basicInformation"`
	BankAccount      BankAccount            `json:"bankAccount" bson:"bankAccount"`
	Financial        BrokerFinancial        `json:"financial" bson:"financial"`
	Contract         BrokerContract         `json:"contract" bson:"contract"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Broker struct {
	BaseBroker `bson:",inline"`
}

type BrokerResponse struct {
	BaseBroker `bson:",inline"`
}

type BrokersWithPagination = types.DataWithPagination[BrokerResponse]
