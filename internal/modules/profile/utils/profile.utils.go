package profile_utils

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	profile_enums "larsa-e-commerce-microservice-v2/internal/modules/profile/enums"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertRegistrationCertificationDtoToModel converts a RegistrationCertificationInfoDto to a RegistrationCertificationInfo model
func ConvertRegistrationCertificationDtoToModel(dto dtos.RegistrationCertificationInfoDto) models.RegistrationCertificationInfo {
	return models.RegistrationCertificationInfo{
		CompanyName:                dto.CompanyName,
		CommercialRegistrationNum:  dto.CommercialRegistrationNum,
		ExpirationDate:             dto.ExpirationDate,
		IssuingDate:                dto.IssuingDate,
		IssuingAuthority:           dto.IssuingAuthority,
		IssuingCountry:             dto.IssuingCountry,
		BusinessActivityTypes:      dto.BusinessActivityTypes,
		TaxIdentificationNumber:    dto.TaxIdentificationNumber,
		SpecialLicenseNumber:       dto.SpecialLicenseNumber,
		CertificateFile:            dto.CertificateFile,
		RelatedLicensesFile:        dto.RelatedLicensesFile,
		SupportingCertificatesFile: dto.SupportingCertificatesFile,
		DataAccuracyDeclaration:    dto.DataAccuracyDeclaration,
	}
}

// ConvertQualityCertificationDtoToModel converts a QualityCertificationInfoDto to a QualityCertificationInfo model
func ConvertQualityCertificationDtoToModel(dto dtos.QualityCertificationInfoDto) models.QualityCertificationInfo {
	return models.QualityCertificationInfo{
		CertificateName:         dto.CertificateName,
		CertificateType:         dto.CertificateType,
		CertificateNumber:       dto.CertificateNumber,
		IssuanceDate:            dto.IssuanceDate,
		ExpiryDate:              dto.ExpiryDate,
		GrantingAuthority:       dto.GrantingAuthority,
		CountryOfIssuance:       dto.CountryOfIssuance,
		ScopeOfCoverage:         dto.ScopeOfCoverage,
		ProductsLinesCovered:    dto.ProductsLinesCovered,
		AdditionalDetails:       dto.AdditionalDetails,
		CertificateFile:         dto.CertificateFile,
		SupportingFile:          dto.SupportingFile,
		DisplayOnProfile:        dto.DisplayOnProfile,
		RelatedCertifications:   dto.RelatedCertifications,
		IsRenewedAnnually:       dto.IsRenewedAnnually,
		RenewalDetails:          dto.RenewalDetails,
		DataAccuracyDeclaration: dto.DataAccuracyDeclaration,
	}
}

// ConvertLicenseDtoToModel converts a LicenseInfoDto to a LicenseInfo model
func ConvertLicenseDtoToModel(dto dtos.LicenseInfoDto) models.LicenseInfo {
	return models.LicenseInfo{
		LicenseName:             dto.LicenseName,
		LicenseNumber:           dto.LicenseNumber,
		LicenseClass:            dto.LicenseClass,
		ExpirationDate:          dto.ExpirationDate,
		IssuingDate:             dto.IssuingDate,
		IssuingAuthority:        dto.IssuingAuthority,
		IssuingCountry:          dto.IssuingCountry,
		Products:                dto.Products,
		Description:             dto.Description,
		LicenseFile:             dto.LicenseFile,
		SupportingFile:          dto.SupportingFile,
		DataAccuracyDeclaration: dto.DataAccuracyDeclaration,
	}
}

// ConvertBaseProfileDtoToModel converts a BaseProfileDto to a BaseProfile model
func ConvertBaseProfileDtoToModel(dto dtos.BaseProfileDto, userId primitive.ObjectID, userType enums.UserType, profileType profile_enums.ProfileType) models.BaseProfile {
	now := time.Now()
	id, err := primitive.ObjectIDFromHex(*dto.UserId)
	if err != nil {
		id = primitive.NewObjectID()
	}

	return models.BaseProfile{
		FirstName:              dto.FirstName,
		LastName:               dto.LastName,
		MobileCode:             dto.MobileCode,
		MobileNumber:           dto.MobileNumber,
		Gender:                 dto.Gender,
		DateOfBirth:            dto.DateOfBirth,
		CompanyName:            dto.CompanyName,
		CompanyWebsite:         dto.CompanyWebsite,
		Industry:               dto.Industry,
		Country:                dto.Country,
		City:                   dto.City,
		Area:                   dto.Area,
		Address:                dto.Address,
		Phone:                  dto.Phone,
		Email:                  dto.Email,
		ResponsiblePerson:      dto.ResponsiblePerson,
		ResponsibleRole:        dto.ResponsibleRole,
		InterestedInAgents:     dto.InterestedInAgents,
		InterestedInExports:    dto.InterestedInExports,
		InterestedInSeekAgents: dto.InterestedInSeekAgents,
		AdditionalInfo:         dto.AdditionalInfo,
		UserType:               userType,
		ProfileType:            &profileType,
		MetaDataBaseModel: base_models.MetaDataBaseModel{
			Id:        id,
			CreatedAt: now,
			UpdatedAt: now,
			CreatedBy: &userId,
			UpdatedBy: &userId,
		},
	}
}

// ConvertToVendorProfileModel converts a CreateVendorProfileDto to a VendorProfile model
func ConvertToVendorProfileModel(dto *dtos.CreateVendorProfileDto, userId primitive.ObjectID) *models.VendorProfile {
	baseProfile := ConvertBaseProfileDtoToModel(dto.BaseProfileDto, userId, enums.Vendor, profile_enums.Vendor)
	RegistrationCertifications := make([]models.RegistrationCertificationInfo, len(dto.RegistrationCertifications))
	for i, cert := range dto.RegistrationCertifications {
		RegistrationCertifications[i] = ConvertRegistrationCertificationDtoToModel(cert)
	}
	QualityCertifications := make([]models.QualityCertificationInfo, len(dto.QualityCertifications))
	for i, cert := range dto.QualityCertifications {
		QualityCertifications[i] = ConvertQualityCertificationDtoToModel(cert)
	}
	Licenses := make([]models.LicenseInfo, len(dto.Licenses))
	for i, license := range dto.Licenses {
		Licenses[i] = ConvertLicenseDtoToModel(license)
	}

	return &models.VendorProfile{
		BaseProfile:                baseProfile,
		RegistrationCertifications: RegistrationCertifications,
		QualityCertifications:      QualityCertifications,
		Licenses:                   Licenses,
	}
}

// ConvertToBuyerProfileModel converts a CreateBuyerProfileDto to a BuyerProfile model
func ConvertToBuyerProfileModel(dto *dtos.CreateBuyerProfileDto, userId primitive.ObjectID) *models.BuyerProfile {
	baseProfile := ConvertBaseProfileDtoToModel(dto.BaseProfileDto, userId, enums.Member, profile_enums.Buyer)

	return &models.BuyerProfile{
		BaseProfile: baseProfile,
	}
}

// UpdateBaseProfileFromDto updates a Profile model from a BaseProfileDto
func UpdateBaseProfileFromDto(profile *models.BaseProfile, dto dtos.BaseProfileDto, updatedBy primitive.ObjectID) {
	now := time.Now()
	profile.FirstName = dto.FirstName
	profile.LastName = dto.LastName
	profile.MobileCode = dto.MobileCode
	profile.MobileNumber = dto.MobileNumber
	profile.Gender = dto.Gender
	profile.DateOfBirth = dto.DateOfBirth
	profile.CompanyName = dto.CompanyName
	profile.CompanyWebsite = dto.CompanyWebsite
	profile.Industry = dto.Industry
	profile.Country = dto.Country
	profile.City = dto.City
	profile.Area = dto.Area
	profile.Address = dto.Address
	profile.Phone = dto.Phone
	profile.Email = dto.Email
	profile.ResponsiblePerson = dto.ResponsiblePerson
	profile.ResponsibleRole = dto.ResponsibleRole
	profile.InterestedInAgents = dto.InterestedInAgents
	profile.InterestedInExports = dto.InterestedInExports
	profile.InterestedInSeekAgents = dto.InterestedInSeekAgents
	profile.AdditionalInfo = dto.AdditionalInfo
	profile.MetaDataBaseModel.UpdatedAt = now
	profile.MetaDataBaseModel.UpdatedBy = &updatedBy
}

// ConvertToGuestProfileModel converts a CreateGuestProfileDto to a GuestProfile model
func ConvertToGuestProfileModel(dto *dtos.CreateGuestProfileDto, userId primitive.ObjectID) *models.GuestProfile {
	userIdString := userId.Hex()
	baseProfile := ConvertBaseProfileDtoToModel(dtos.BaseProfileDto{
		UserId:    &userIdString,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}, userId, enums.Guest, profile_enums.Guest)

	profileType := profile_enums.Guest
	baseProfile.ProfileType = &profileType
	baseProfile.UserType = enums.Guest
	baseProfile.MetaDataBaseModel.Id = userId

	return &models.GuestProfile{
		BaseProfile:   baseProfile,
		FingerprintId: dto.FingerprintId,
	}
}
