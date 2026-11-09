package profile_services

import (
	"context"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateInviteUser creates a new user based on invitation and creates a profile
func (s *profileService) CreateInviteUser(ctx context.Context, dto *dtos.CreateInviteUserDto) (*models.ProfileResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		cfg = &types.AppCfg{Lang: "en"} // Default language if no config
	}

	// Check if user already exists with this email
	existingUser, userErr := s.usersgw.GetUserByEmail(ctx, *dto.Email)
	if userErr == nil && existingUser != nil {
		// User already exists, try to get their profile
		userObjId, objIdErr := primitive.ObjectIDFromHex(existingUser.Id)
		if objIdErr != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Invalid user ID format",
			}
		}
		existingProfile, profileErr := s.profileRepo.GetProfileByUserId(ctx, userObjId)
		if profileErr == nil && existingProfile != nil {
			return existingProfile, nil
		}
		//create profile for the user
		buyerDto := &dtos.CreateBuyerProfileDto{
			BaseProfileDto: dtos.BaseProfileDto{
				FirstName: &existingUser.FirstName,
				LastName:  &existingUser.LastName,
				Email:     dto.Email,
			},
		}
		profile, profileErr := s.profileRepo.CreateBuyerProfile(ctx, buyerDto)
		if profileErr != nil {
			return nil, profileErr
		}
		return profile, nil
	}

	// Generate a temporary password for the invited user
	tempPassword := util.GeneratePassword(12, 2, 2, 2)

	// Set default values if not provided
	firstName := "Invited"
	lastName := "User"

	if dto.FirstName != nil && *dto.FirstName != "" {
		firstName = *dto.FirstName
	}
	if dto.LastName != nil && *dto.LastName != "" {
		lastName = *dto.LastName
	}

	// Create user data
	newUser := &types.PostUserData{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        *dto.Email,
		Password:     tempPassword,
		Roles:        []primitive.ObjectID{},
		Capabilities: []primitive.ObjectID{},
	}

	// Create the user
	userId, userErr := s.usersgw.AddUser(ctx, newUser)
	if userErr != nil {
		return nil, &helpers.ApiError{
			Message:    "Error creating invited user: " + userErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if userId.IsZero() {
		return nil, &helpers.ApiError{
			Message:    "Invited user not created",
			StatusCode: http.StatusInternalServerError,
		}
	}

	cfg.User = &types.User{
		Id:       userId,
		Type:     enums.Buyer,
		UserData: common.User{Id: userId.Hex()},
	}
	ctx = context.WithValue(context.WithoutCancel(ctx), util.ReqAppCfg, cfg)

	cfg, err = util.GetReqAppCfg(ctx)

	userIdString := userId.Hex()
	// Create buyer profile (default)
	buyerDto := &dtos.CreateBuyerProfileDto{
		BaseProfileDto: dtos.BaseProfileDto{
			FirstName:   &firstName,
			LastName:    &lastName,
			Email:       dto.Email,
			CompanyName: dto.CompanyName,
			Phone:       dto.Phone,
			UserId:      &userIdString,
		},
	}
	profile, profileErr := s.profileRepo.CreateBuyerProfile(ctx, buyerDto)
	if profileErr != nil {
		return nil, profileErr
	}

	// Send welcome email for newly created user
	go s.sendWelcomeEmail(ctx, *dto.Email, firstName, tempPassword, cfg.Lang)
	return profile, nil

}
