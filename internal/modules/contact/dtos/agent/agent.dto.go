package contact_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"time"
)

type AgentBasicInfoDto struct {
	FirstName         string `json:"firstName" validate:"required"`
	LastName          string `json:"lastName" validate:"required"`
	CompanyName       string `json:"companyName" validate:"required"`
	LicenseNumber     string `json:"licenseNumber" validate:"required"`
	YearsOfExperience int    `json:"yearsOfExperience" validate:"required,min=0"`
}

type AgentBusinessInfoDto struct {
	CommercialRegistrationNumber string   `json:"commercialRegistrationNumber" validate:"required"`
	TaxIdentificationNumber      string   `json:"taxIdentificationNumber" validate:"required"`
	PhoneNumber                  string   `json:"phoneNumber" validate:"required"`
	HeadOfficeAddress            string   `json:"headOfficeAddress" validate:"required"`
	OfficialEmail                string   `json:"officialEmail" validate:"required,email"`
	AuthorizedCountryCodes       []string `json:"authorizedCountryCodes" validate:"required,min=1,dive,required"`
	WebsiteLink                  string   `json:"websiteLink" validate:"required,url"`
}

type AgentFinancialInfoDto struct {
	ShippingProfitPercentage     float64 `json:"shippingProfitPercentage" validate:"required,min=0,max=100"`
	SubscriptionProfitPercentage float64 `json:"subscriptionProfitPercentage" validate:"required,min=0,max=100"`
}

type AgentContractInfoDto struct {
	StartDate   time.Time            `json:"startDate" validate:"required"`
	EndDate     time.Time            `json:"endDate" validate:"required,gtfield=StartDate"`
	Attachments []types.FileFieldDto `json:"attachments" validate:"required,min=1,dive,required"`
}

type AgentBankAccountInfoDto struct {
	BankName          string `json:"bankName" validate:"required"`
	AccountNumber     string `json:"accountNumber" validate:"required"`
	AccountHolderName string `json:"accountHolderName" validate:"required"`
	Iban              string `json:"iban" validate:"required"`
	SwiftCode         string `json:"swiftCode" validate:"required"`
}

type CreateAgentDto struct {
	ResponsibleForLogistics bool                    `json:"responsibleForLogistics"`
	Active                  bool                    `json:"active"`
	BasicInformation        AgentBasicInfoDto       `json:"basicInformation" validate:"required"`
	BusinessInfo            AgentBusinessInfoDto    `json:"businessInfo" validate:"required"`
	FinancialInfo           AgentFinancialInfoDto   `json:"financialInfo" validate:"required"`
	ContractInfo            AgentContractInfoDto    `json:"contractInfo" validate:"required"`
	BankAccountInfo         AgentBankAccountInfoDto `json:"bankAccountInfo" validate:"required"`
	Attachments             []types.FileFieldDto    `json:"attachments" validate:"required,dive,required"`
}
