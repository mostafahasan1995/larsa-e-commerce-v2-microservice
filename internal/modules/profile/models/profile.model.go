package profile_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	profile_enums "larsa-e-commerce-microservice-v2/internal/modules/profile/enums"
)

// CertificationInfo represents certification details for vendors
type RegistrationCertificationInfo struct {
	CompanyName                *string             `bson:"companyName" json:"companyName"`
	CommercialRegistrationNum  *string             `bson:"commercialRegNum" json:"commercialRegNum"`
	ExpirationDate             *time.Time          `bson:"expirationDate" json:"expirationDate"`
	IssuingDate                *time.Time          `bson:"issuingDate" json:"issuingDate"`
	IssuingAuthority           *string             `bson:"issuingAuthority" json:"issuingAuthority"`
	IssuingCountry             *string             `bson:"issuingCountry" json:"issuingCountry"`
	BusinessActivityTypes      []string            `bson:"businessActivityTypes" json:"businessActivityTypes"`
	TaxIdentificationNumber    *string             `bson:"taxIdentificationNumber,omitempty" json:"taxIdentificationNumber,omitempty"`
	SpecialLicenseNumber       *string             `bson:"specialLicenseNumber,omitempty" json:"specialLicenseNumber,omitempty"`
	CertificateFile            *types.FileFieldDto `bson:"certificateFile" json:"certificateFile"`
	RelatedLicensesFile        *types.FileFieldDto `bson:"relatedLicensesFile,omitempty" json:"relatedLicensesFile,omitempty"`
	SupportingCertificatesFile *types.FileFieldDto `bson:"supportingCertificatesFile,omitempty" json:"supportingCertificatesFile,omitempty"`
	DataAccuracyDeclaration    *bool               `bson:"dataAccuracyDeclaration" json:"dataAccuracyDeclaration"`
}

// LicenseInfo represents license details for vendors
type LicenseInfo struct {
	LicenseName             *string             `bson:"licenseName" json:"licenseName"`
	LicenseNumber           *string             `bson:"licenseNumber" json:"licenseNumber"`
	LicenseClass            *string             `bson:"licenseClass" json:"licenseClass"`
	ExpirationDate          *time.Time          `bson:"expirationDate" json:"expirationDate"`
	IssuingDate             *time.Time          `bson:"issuingDate" json:"issuingDate"`
	IssuingAuthority        *string             `bson:"issuingAuthority" json:"issuingAuthority"`
	IssuingCountry          *string             `bson:"issuingCountry" json:"issuingCountry"`
	Products                []string            `bson:"products" json:"products"`
	Description             *string             `bson:"description" json:"description"`
	LicenseFile             *types.FileFieldDto `bson:"licenseFile" json:"licenseFile"`
	SupportingFile          *types.FileFieldDto `bson:"supportingFile" json:"supportingFile"`
	DataAccuracyDeclaration *bool               `bson:"dataAccuracyDeclaration" json:"dataAccuracyDeclaration"`
}

// QualityCertificationInfo represents quality certification details for vendors
type QualityCertificationInfo struct {
	CertificateName         *string             `bson:"certificateName" json:"certificateName"`
	CertificateType         *string             `bson:"certificateType" json:"certificateType"`
	CertificateNumber       *string             `bson:"certificateNumber" json:"certificateNumber"`
	IssuanceDate            *time.Time          `bson:"issuanceDate" json:"issuanceDate"`
	ExpiryDate              *time.Time          `bson:"expiryDate,omitempty" json:"expiryDate,omitempty"`
	GrantingAuthority       *string             `bson:"grantingAuthority" json:"grantingAuthority"`
	CountryOfIssuance       *string             `bson:"countryOfIssuance" json:"countryOfIssuance"`
	ScopeOfCoverage         []string            `bson:"scopeOfCoverage" json:"scopeOfCoverage"`
	ProductsLinesCovered    []string            `bson:"productsLinesCovered" json:"productsLinesCovered"`
	AdditionalDetails       []string            `bson:"additionalDetails,omitempty" json:"additionalDetails,omitempty"`
	CertificateFile         *types.FileFieldDto `bson:"certificateFile" json:"certificateFile"`
	SupportingFile          *types.FileFieldDto `bson:"supportingFile" json:"supportingFile"`
	DisplayOnProfile        *bool               `bson:"displayOnProfile" json:"displayOnProfile"`
	RelatedCertifications   []string            `bson:"relatedCertifications,omitempty" json:"relatedCertifications,omitempty"`
	IsRenewedAnnually       *bool               `bson:"isRenewedAnnually,omitempty" json:"isRenewedAnnually,omitempty"`
	RenewalDetails          []string            `bson:"renewalDetails,omitempty" json:"renewalDetails,omitempty"`
	DataAccuracyDeclaration *bool               `bson:"dataAccuracyDeclaration" json:"dataAccuracyDeclaration"`
}

// Profile represents base fields for both vendor and buyer profiles
type BaseProfile struct {
	FirstName              *string                    `bson:"firstName" json:"firstName"`
	LastName               *string                    `bson:"lastName" json:"lastName"`
	MobileCode             *string                    `bson:"mobileCode" json:"mobileCode"`
	MobileNumber           *string                    `bson:"mobileNumber" json:"mobileNumber"`
	Gender                 *string                    `bson:"gender" json:"gender"`
	DateOfBirth            *time.Time                 `bson:"dateOfBirth" json:"dateOfBirth"`
	CompanyName            *string                    `bson:"companyName" json:"companyName"`
	CompanyWebsite         *string                    `bson:"companyWebsite" json:"companyWebsite"`
	Industry               []string                   `bson:"industry" json:"industry"`
	Country                *string                    `bson:"country" json:"country"`
	City                   *string                    `bson:"city" json:"city"`
	Area                   *string                    `bson:"area" json:"area"`
	Address                *string                    `bson:"address" json:"address"`
	Phone                  *string                    `bson:"phone" json:"phone"`
	Email                  *string                    `bson:"email" json:"email"`
	ResponsiblePerson      *string                    `bson:"responsiblePerson,omitempty" json:"responsiblePerson,omitempty"`
	ResponsibleRole        *string                    `bson:"responsibleRole,omitempty" json:"responsibleRole,omitempty"`
	InterestedInAgents     *string                    `bson:"interestedInAgents" json:"interestedInAgents"`
	InterestedInExports    *string                    `bson:"interestedInExports" json:"interestedInExports"`
	InterestedInSeekAgents *string                    `bson:"interestedInSeekAgents" json:"interestedInSeekAgents"`
	AdditionalInfo         *string                    `bson:"additionalInfo,omitempty" json:"additionalInfo,omitempty"`
	ProfileType            *profile_enums.ProfileType `bson:"profileType" json:"profileType"`
	UserType               enums.UserType             `bson:"userType" json:"userType"`
	// Embedded metadata model
	base_models.MetaDataBaseModel `bson:",inline"`
}

// VendorProfile extends the base Profile with vendor-specific fields
type VendorProfile struct {
	BaseProfile                `bson:",inline" json:",inline"`
	RegistrationCertifications []RegistrationCertificationInfo `bson:"registrationCertifications,omitempty" json:"registrationCertifications,omitempty"`
	QualityCertifications      []QualityCertificationInfo      `bson:"qualityCertifications,omitempty" json:"qualityCertifications,omitempty"`
	Licenses                   []LicenseInfo                   `bson:"licenses,omitempty" json:"licenses,omitempty"`
}

// BuyerProfile uses the base Profile without additional fields
type BuyerProfile struct {
	BaseProfile `bson:",inline" json:",inline"`
}

type CommonProfile struct {
	BaseProfile                `bson:",inline" json:",inline"`
	RegistrationCertifications []RegistrationCertificationInfo `bson:"registrationCertifications,omitempty" json:"registrationCertifications,omitempty"`
	QualityCertifications      []QualityCertificationInfo      `bson:"qualityCertifications,omitempty" json:"qualityCertifications,omitempty"`
	Licenses                   []LicenseInfo                   `bson:"licenses,omitempty" json:"licenses,omitempty"`
}

// ProfileResponse is the response structure for profile data
type ProfileResponse struct {
	FirstName                     *string                         `bson:"firstName" json:"firstName"`
	LastName                      *string                         `bson:"lastName" json:"lastName"`
	MobileCode                    *string                         `bson:"mobileCode" json:"mobileCode"`
	MobileNumber                  *string                         `bson:"mobileNumber" json:"mobileNumber"`
	Gender                        *string                         `bson:"gender" json:"gender"`
	DateOfBirth                   *time.Time                      `bson:"dateOfBirth" json:"dateOfBirth"`
	CompanyName                   *string                         `bson:"companyName" json:"companyName"`
	CompanyWebsite                *string                         `bson:"companyWebsite" json:"companyWebsite"`
	Industry                      []string                        `bson:"industry" json:"industry"`
	Country                       *string                         `bson:"country" json:"country"`
	City                          *string                         `bson:"city" json:"city"`
	Area                          *string                         `bson:"area" json:"area"`
	Address                       *string                         `bson:"address" json:"address"`
	Phone                         *string                         `bson:"phone" json:"phone"`
	Email                         *string                         `bson:"email" json:"email"`
	ResponsiblePerson             *string                         `bson:"responsiblePerson,omitempty" json:"responsiblePerson,omitempty"`
	ResponsibleRole               *string                         `bson:"responsibleRole,omitempty" json:"responsibleRole,omitempty"`
	InterestedInAgents            *string                         `bson:"interestedInAgents" json:"interestedInAgents"`
	InterestedInExports           *string                         `bson:"interestedInExports" json:"interestedInExports"`
	InterestedInSeekAgents        *string                         `bson:"interestedInSeekAgents" json:"interestedInSeekAgents"`
	AdditionalInfo                *string                         `bson:"additionalInfo,omitempty" json:"additionalInfo,omitempty"`
	ProfileType                   *profile_enums.ProfileType      `bson:"profileType" json:"profileType"`
	UserType                      *enums.UserType                 `bson:"userType" json:"userType"`
	FingerprintId                 *string                         `bson:"fingerprintId,omitempty" json:"fingerprintId,omitempty"`
	RegistrationCertifications    []RegistrationCertificationInfo `bson:"registrationCertifications,omitempty" json:"registrationCertifications,omitempty"`
	QualityCertifications         []QualityCertificationInfo      `bson:"qualityCertifications,omitempty" json:"qualityCertifications,omitempty"`
	Licenses                      []LicenseInfo                   `bson:"licenses,omitempty" json:"licenses,omitempty"`
	base_models.MetaDataBaseModel `bson:",inline"`
}

// ProfilesWithPagination represents paginated profile results
type ProfilesWithPagination = types.DataWithPagination[ProfileResponse]

// GuestProfile represents a guest user profile
type GuestProfile struct {
	BaseProfile   `bson:",inline" json:",inline"`
	FingerprintId *string `bson:"fingerprintId" json:"fingerprintId"`
}
