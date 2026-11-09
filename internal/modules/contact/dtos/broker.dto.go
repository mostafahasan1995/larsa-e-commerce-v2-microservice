package contact_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"time"
)

type BankAccountDto struct {
	BankName          types.Localizable[string] `json:"bankName" validate:"required,localizable"`
	AccountNumber     string                    `json:"accountNumber" validate:"required"`
	AccountHolderName string                    `json:"accountHolderName" validate:"required"`
	IBAN              string                    `json:"iban" validate:"required"`
	SWIFTCode         string                    `json:"swiftCode" validate:"required"`
}

type BrokerBasicInformationDto struct {
	FullName                     types.Localizable[string] `json:"fullName" validate:"required,localizable"`
	CompanyName                  types.Localizable[string] `json:"companyName" validate:"required,localizable"`
	CustomsLicenseNumber         string                    `json:"customsLicenseNumber" validate:"required"`
	CommercialRegistrationNumber string                    `json:"commercialRegistrationNumber" validate:"required"`
	TaxIdentificationNumber      string                    `json:"taxIdentificationNumber" validate:"required"`
	AuthorizedCountry            string                    `json:"authorizedCountry" validate:"required"`
	YearsOfExperience            string                    `json:"yearsOfExperience" validate:"required"`
}

type BrokerContractDto struct {
	HeadOfficeAddress string `json:"headOfficeAddress" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	PhoneNumber       string `json:"phoneNumber" validate:"required"`
}

type BrokerFinancialDto struct {
	StartDate   time.Time            `json:"startDate" validate:"required"`
	EndDate     time.Time            `json:"endDate" validate:"required"`
	Attachments []types.FileFieldDto `json:"attachments" validate:"required,min=1,dive,required"`
}

type CreateBrokerDto struct {
	BasicInformation BrokerBasicInformationDto `json:"basicInformation" validate:"required"`
	BankAccount      BankAccountDto            `json:"bankAccount" validate:"required"`
	Financial        BrokerFinancialDto        `json:"financial" validate:"required"`
	Contract         BrokerContractDto         `json:"contract" validate:"required"`
}
