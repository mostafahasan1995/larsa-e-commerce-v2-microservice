package contact_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"time"

	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type AgentBasicInfo struct {
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	CompanyName       string `bson:"companyName" json:"companyName"`
	LicenseNumber     string `bson:"licenseNumber" json:"licenseNumber"`
	YearsOfExperience int    `bson:"yearsOfExperience" json:"yearsOfExperience"`
}

type AgentBusinessInfo struct {
	CommercialRegistrationNumber string   `bson:"commercialRegistrationNumber" json:"commercialRegistrationNumber"`
	TaxIdentificationNumber      string   `bson:"taxIdentificationNumber" json:"taxIdentificationNumber"`
	PhoneNumber                  string   `bson:"phoneNumber" json:"phoneNumber"`
	HeadOfficeAddress            string   `bson:"headOfficeAddress" json:"headOfficeAddress"`
	OfficialEmail                string   `bson:"officialEmail" json:"officialEmail"`
	AuthorizedCountryCodes       []string `bson:"authorizedCountryCodes" json:"authorizedCountryCodes"`
	WebsiteLink                  string   `bson:"websiteLink" json:"websiteLink"`
}

type AgentFinancialInfo struct {
	ShippingProfitPercentage     float64 `bson:"shippingProfitPercentage" json:"shippingProfitPercentage"`
	SubscriptionProfitPercentage float64 `bson:"subscriptionProfitPercentage" json:"subscriptionProfitPercentage"`
}

type AgentContractInfo struct {
	StartDate   time.Time            `bson:"startDate" json:"startDate"`
	EndDate     time.Time            `bson:"endDate" json:"endDate"`
	Attachments []types.FileFieldDto `bson:"attachments" json:"attachments"`
}

type AgentBankAccountInfo struct {
	BankName          string `bson:"bankName" json:"bankName"`
	AccountNumber     string `bson:"accountNumber" json:"accountNumber"`
	AccountHolderName string `bson:"accountHolderName" json:"accountHolderName"`
	IBAN              string `bson:"iban" json:"iban"`
	SWIFTCode         string `bson:"swiftCode" json:"swiftCode"`
}

type BaseAgent struct {
	base_models.MetaDataBaseModel `bson:",inline"`

	ResponsibleForLogistics bool `bson:"responsibleForLogistics" json:"responsibleForLogistics"`
	Active                  bool `bson:"active" json:"active"`

	// Basic Information
	BasicInformation AgentBasicInfo `bson:"basicInformation" json:"basicInformation"`

	// Business & Contact Information
	BusinessInfo AgentBusinessInfo `bson:"businessInfo" json:"businessInfo"`

	// Financial Information
	FinancialInfo AgentFinancialInfo `bson:"financialInfo" json:"financialInfo"`

	// Contract Information
	ContractInfo AgentContractInfo `bson:"contractInfo" json:"contractInfo"`

	// Bank Account Information
	BankAccountInfo AgentBankAccountInfo `bson:"bankAccountInfo" json:"bankAccountInfo"`

	// Attachments
	Attachments []types.FileFieldDto `bson:"attachments" json:"attachments"`
}

type Agent struct {
	BaseAgent `bson:",inline"`
}

type AgentResponse struct {
	BaseAgent `bson:",inline"`
}

type AgentsWithPagination = types.DataWithPagination[AgentResponse]
