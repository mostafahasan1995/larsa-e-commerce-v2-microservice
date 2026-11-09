package profile_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"time"
)

// CreateVendorProfileDto represents the data required to create a vendor profile
type CreateVendorProfileDto struct {
	BaseProfileDto
	RegistrationCertifications []RegistrationCertificationInfoDto `json:"registrationCertifications,omitempty" validate:"required"`
	QualityCertifications      []QualityCertificationInfoDto      `json:"qualityCertifications,omitempty" validate:"required"`
	Licenses                   []LicenseInfoDto                   `json:"licenses,omitempty" validate:"required"`
}

type BaseProfileDto struct {
	UserId                 *string    `json:"userId,omitempty"`
	FirstName              *string    `json:"firstName" validate:"omitempty"`
	LastName               *string    `json:"lastName" validate:"omitempty"`
	MobileCode             *string    `json:"mobileCode" validate:"omitempty"`
	MobileNumber           *string    `json:"mobileNumber" validate:"omitempty"`
	Gender                 *string    `json:"gender" validate:"omitempty,oneof=male female other"`
	DateOfBirth            *time.Time `json:"dateOfBirth" validate:"omitempty"`
	CompanyName            *string    `json:"companyName" validate:"omitempty"`
	CompanyWebsite         *string    `json:"companyWebsite" validate:"omitempty,url"`
	Industry               []string   `json:"industry" validate:"omitempty,min=1"`
	Country                *string    `json:"country" validate:"omitempty"`
	City                   *string    `json:"city" validate:"omitempty"`
	Area                   *string    `json:"area" validate:"omitempty"`
	Address                *string    `json:"address" validate:"omitempty"`
	Phone                  *string    `json:"phone" validate:"omitempty"`
	Email                  *string    `json:"email" validate:"required,email"`
	ResponsiblePerson      *string    `json:"responsiblePerson,omitempty"`
	ResponsibleRole        *string    `json:"responsibleRole,omitempty"`
	InterestedInAgents     *string    `json:"interestedInAgents,omitempty"`
	InterestedInExports    *string    `json:"interestedInExports,omitempty"`
	InterestedInSeekAgents *string    `json:"interestedInSeekAgents,omitempty"`
	AdditionalInfo         *string    `json:"additionalInfo,omitempty"`
}

// CreateGuestProfileDto represents the data required to create a guest profile
type CreateGuestProfileDto struct {
	FingerprintId *string `json:"fingerprintId" validate:"required"`
	FirstName     *string
	LastName      *string
	Email         *string
}

// CreateInviteUserDto represents the data required to invite a user to the system
type CreateInviteUserDto struct {
	Email       *string `json:"email" validate:"required,email"`
	FirstName   *string `json:"firstName" validate:"omitempty"`
	LastName    *string `json:"lastName" validate:"omitempty"`
	CompanyName *string `json:"companyName" validate:"omitempty"`
	Phone       *string `json:"phone" validate:"omitempty"`
	UserType    *string `json:"userType" validate:"omitempty,oneof=vendor buyer"` // Default to buyer if not specified
}

type RegistrationCertificationInfoDto struct {
	CompanyName                *string             `json:"companyName" validate:"required"`
	CommercialRegistrationNum  *string             `json:"commercialRegNum" validate:"required"`
	ExpirationDate             *time.Time          `json:"expirationDate" validate:"required"`
	IssuingDate                *time.Time          `json:"issuingDate" validate:"required"`
	IssuingAuthority           *string             `json:"issuingAuthority" validate:"required"`
	IssuingCountry             *string             `json:"issuingCountry" validate:"required"`
	BusinessActivityTypes      []string            `json:"businessActivityTypes" validate:"required,min=1"`
	TaxIdentificationNumber    *string             `json:"taxIdentificationNumber,omitempty"`
	SpecialLicenseNumber       *string             `json:"specialLicenseNumber,omitempty"`
	CertificateFile            *types.FileFieldDto `json:"certificateFile" validate:"required"`
	RelatedLicensesFile        *types.FileFieldDto `json:"relatedLicensesFile,omitempty"`
	SupportingCertificatesFile *types.FileFieldDto `json:"supportingCertificatesFile,omitempty"`
	DataAccuracyDeclaration    *bool               `json:"dataAccuracyDeclaration" validate:"required"`
}

type LicenseInfoDto struct {
	LicenseName             *string             `json:"licenseName" validate:"required"`
	LicenseNumber           *string             `json:"licenseNumber" validate:"required"`
	LicenseClass            *string             `json:"licenseClass" validate:"required"`
	ExpirationDate          *time.Time          `json:"expirationDate" validate:"required"`
	IssuingDate             *time.Time          `json:"issuingDate" validate:"required"`
	IssuingAuthority        *string             `json:"issuingAuthority" validate:"required"`
	IssuingCountry          *string             `json:"issuingCountry" validate:"required"`
	Products                []string            `json:"products" validate:"required,min=1"`
	Description             *string             `json:"description" validate:"required"`
	LicenseFile             *types.FileFieldDto `json:"licenseFile" validate:"required"`
	SupportingFile          *types.FileFieldDto `json:"supportingFile,omitempty"`
	DataAccuracyDeclaration *bool               `json:"dataAccuracyDeclaration" validate:"required"`
}

type QualityCertificationInfoDto struct {
	CertificateName         *string             `json:"certificateName" validate:"required"`
	CertificateType         *string             `json:"certificateType" validate:"required"`
	CertificateNumber       *string             `json:"certificateNumber" validate:"required"`
	IssuanceDate            *time.Time          `json:"issuanceDate" validate:"required"`
	ExpiryDate              *time.Time          `json:"expiryDate,omitempty"`
	GrantingAuthority       *string             `json:"grantingAuthority" validate:"required"`
	CountryOfIssuance       *string             `json:"countryOfIssuance" validate:"required"`
	ScopeOfCoverage         []string            `json:"scopeOfCoverage" validate:"required,min=1"`
	ProductsLinesCovered    []string            `json:"productsLinesCovered" validate:"required,min=1"`
	AdditionalDetails       []string            `json:"additionalDetails,omitempty"`
	CertificateFile         *types.FileFieldDto `json:"certificateFile" validate:"required"`
	SupportingFile          *types.FileFieldDto `json:"supportingFile,omitempty"`
	DisplayOnProfile        *bool               `json:"displayOnProfile" validate:"required"`
	RelatedCertifications   []string            `json:"relatedCertifications,omitempty"`
	IsRenewedAnnually       *bool               `json:"isRenewedAnnually,omitempty"`
	RenewalDetails          []string            `json:"renewalDetails,omitempty"`
	DataAccuracyDeclaration *bool               `json:"dataAccuracyDeclaration" validate:"required"`
}

// CreateBuyerProfileDto represents the data required to create a buyer profile
type CreateBuyerProfileDto struct {
	BaseProfileDto `bson:",inline" json:",inline"`
	Password       *string `json:"password" `
}

// CreateProfileDto represents the data required to create a profile (either vendor or buyer)
type CreateProfileDto struct {
	BaseProfileDto             `bson:",inline" json:",inline"`
	RegistrationCertifications []RegistrationCertificationInfoDto `json:"registrationCertifications,omitempty" validate:"omitempty"`
	QualityCertifications      []QualityCertificationInfoDto      `json:"qualityCertifications,omitempty" validate:"omitempty"`
	Licenses                   []LicenseInfoDto                   `json:"licenses,omitempty" validate:"omitempty"`
}
