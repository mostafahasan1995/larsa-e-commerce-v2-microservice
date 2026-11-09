package profile_services

import (
	"context"
	"fmt"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	profile_enums "larsa-e-commerce-microservice-v2/internal/modules/profile/enums"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
)

// UpdateMyProfile updates the authenticated user's profile
func (s *profileService) UpdateMyProfile(ctx context.Context, dto *dtos.CreateProfileDto) (*models.ProfileResponse, *helpers.ApiError) {

	cfg, _ := util.GetReqAppCfg(ctx)
	userId := cfg.User.Id
	userIdString := userId.Hex()
	// Get current profile
	storedProfile, getErr := s.profileRepo.GetProfileByUserId(ctx, userId)
	dto.BaseProfileDto.UserId = &userIdString

	// If profile doesn't exist, create a new one based on user type
	if getErr != nil || storedProfile == nil {

		// Determine user type from context
		if cfg.User.Type == enums.Vendor {
			// Create vendor profile
			vendorDto := &dtos.CreateVendorProfileDto{
				BaseProfileDto:             dto.BaseProfileDto,
				RegistrationCertifications: dto.RegistrationCertifications,
				QualityCertifications:      dto.QualityCertifications,
				Licenses:                   dto.Licenses,
			}
			if vendorDto.RegistrationCertifications == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "certification is required",
				}
			}
			if vendorDto.QualityCertifications == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "quality certification is required",
				}
			}
			if vendorDto.Licenses == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "license is required",
				}
			}

			return s.profileRepo.CreateVendorProfile(ctx, vendorDto)
		}
		// Create buyer profile
		buyerDto := &dtos.CreateBuyerProfileDto{
			BaseProfileDto: dto.BaseProfileDto,
		}
		return s.profileRepo.CreateBuyerProfile(ctx, buyerDto)

	}

	profileId := storedProfile.Id
	if profileId.IsZero() {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "profile ID is missing",
		}
	}

	// Determine if it's a vendor or buyer profile based on profile type
	var result *models.ProfileResponse

	if cfg.User.Type == enums.Vendor {
		// Convert UpdateProfileDto to CreateVendorProfileDto
		vendorDto := &dtos.CreateVendorProfileDto{
			BaseProfileDto:             dto.BaseProfileDto,
			RegistrationCertifications: dto.RegistrationCertifications,
			QualityCertifications:      dto.QualityCertifications,
			Licenses:                   dto.Licenses,
		}
		if vendorDto.RegistrationCertifications == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "certification is required",
			}
		}
		if vendorDto.QualityCertifications == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "quality certification is required",
			}
		}
		if vendorDto.Licenses == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "license is required",
			}
		}

		result, getErr = s.profileRepo.UpdateVendorProfile(ctx, profileId, vendorDto)
	} else {
		// Convert UpdateProfileDto to CreateBuyerProfileDto
		buyerDto := &dtos.CreateBuyerProfileDto{
			BaseProfileDto: dto.BaseProfileDto,
		}
		result, getErr = s.profileRepo.UpdateBuyerProfile(ctx, profileId, buyerDto)
		result.RegistrationCertifications = []models.RegistrationCertificationInfo{}
		result.QualityCertifications = []models.QualityCertificationInfo{}
		result.Licenses = []models.LicenseInfo{}
	}

	return result, getErr
}

// GetMe retrieves the authenticated user's profile
func (s *profileService) GetMe(ctx context.Context) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	userType := cfg.User.Type

	profile, err := s.profileRepo.GetMe(ctx)
	fmt.Println("get me user type", userType)

	if err == nil && profile != nil {
		return profile, nil
	}

	if err != nil && err.StatusCode != http.StatusNotFound {
		return nil, err
	}

	// Profile not found, create one based on user context

	userId := cfg.User.Id
	userIdString := userId.Hex()

	// Get user details from gateway to populate profile
	user, _ := s.usersgw.GetUserById(ctx, userId.Hex())
	// Create profile based on user type
	if userType == enums.Vendor {
		// Create vendor profile with minimal required fields
		vendorDto := &dtos.CreateVendorProfileDto{
			BaseProfileDto: dtos.BaseProfileDto{
				FirstName:   &user.FirstName,
				LastName:    &user.LastName,
				Email:       &user.Email,
				UserId:      &userIdString,
				Gender:      &user.Gender,
				DateOfBirth: &user.Birthday,
			},
		}
		return s.profileRepo.CreateVendorProfile(ctx, vendorDto)
	} else if userType == enums.Agent {
		agentType := profile_enums.Agent
		userType := enums.Agent
		return &models.ProfileResponse{
			FirstName:   &user.FirstName,
			LastName:    &user.LastName,
			Email:       &user.Email,
			Gender:      &user.Gender,
			DateOfBirth: &user.Birthday,
			ProfileType: &agentType,
			UserType:    &userType,
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id: userId,
			},
		}, nil
	} else if userType == enums.SuperAdmin {
		superAdminType := profile_enums.SuperAdmin
		userType := enums.Agent
		return &models.ProfileResponse{
			FirstName:   &user.FirstName,
			LastName:    &user.LastName,
			Email:       &user.Email,
			Gender:      &user.Gender,
			DateOfBirth: &user.Birthday,
			ProfileType: &superAdminType,
			UserType:    &userType,
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id: userId,
			},
		}, nil
	} else if userType == enums.ShippingCompany {
		shippingCompanyType := profile_enums.ShippingCompany
		userType := enums.ShippingCompany
		return &models.ProfileResponse{
			FirstName:   &user.FirstName,
			LastName:    &user.LastName,
			Email:       &user.Email,
			Gender:      &user.Gender,
			DateOfBirth: &user.Birthday,
			ProfileType: &shippingCompanyType,
			UserType:    &userType,
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id: userId,
			},
		}, nil
	}

	fmt.Println("Create buyer profile")
	// Create buyer profile
	buyerDto := &dtos.CreateBuyerProfileDto{
		BaseProfileDto: dtos.BaseProfileDto{
			FirstName:   &user.FirstName,
			LastName:    &user.LastName,
			Email:       &user.Email,
			UserId:      &userIdString,
			Gender:      &user.Gender,
			DateOfBirth: &user.Birthday,
		},
	}

	if user.Id != "" {
		return s.profileRepo.CreateBuyerProfileFromUser(ctx, buyerDto, userId)
	}

	return s.profileRepo.CreateBuyerProfile(ctx, buyerDto)

}
