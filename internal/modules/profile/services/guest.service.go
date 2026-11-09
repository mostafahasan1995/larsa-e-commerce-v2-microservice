package profile_services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getNextGuestId gets and increments the guest ID counter from settings
func (s *profileService) getNextGuestId(ctx context.Context) (uint64, *helpers.ApiError) {
	setting, err := s.settingRepo.FindOneAndIncreaseSerialNumber(ctx, enums.PROFILE.GetCollectionName())
	if err != nil {
		return 0, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get next guest ID: " + err.Message,
		}
	}
	return setting.CurrentSerialNumber, nil
}

// CreateGuestProfile creates a new guest profile
func (s *profileService) CreateGuestProfile(ctx *context.Context, dto *dtos.CreateGuestProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	guestPassword := util.GeneratePassword(12, 2, 2, 2)

	// Get next guest ID
	guestId, err := s.getNextGuestId(*ctx)
	if err != nil {
		return nil, err
	}
	uniqueId := strconv.FormatUint(guestId, 10)

	guestEmail := "guest_" + uniqueId + "@guest.com"
	guestFirstName := "Guest"
	guestLastName := uniqueId

	//By default when we create a new user in the user service we assign a default role to the user
	// Create user data
	newUser := &types.PostUserData{
		FirstName:    guestFirstName,
		LastName:     guestLastName,
		Email:        guestEmail,
		Password:     guestPassword,
		Roles:        []primitive.ObjectID{},
		Capabilities: []primitive.ObjectID{},
	}
	dto.Email = &guestEmail
	dto.FirstName = &guestFirstName
	dto.LastName = &guestLastName

	// Create the user
	userId, userErr := s.usersgw.AddUser(*ctx, newUser)
	fmt.Println("User ID: ", userId)
	fmt.Println("User error: ", userErr)
	if userErr != nil {
		return nil, &helpers.ApiError{
			Message:    "Error creating user: " + userErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if userId.IsZero() {
		return nil, &helpers.ApiError{
			Message:    "User not created",
			StatusCode: http.StatusInternalServerError,
		}
	}
	_, userErr = s.usersgw.Login(ctx, newUser)

	if userErr != nil {
		return nil, &helpers.ApiError{
			Message:    "Error Login user: " + userErr.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	// Now create the guest profile using the repository
	return s.profileRepo.CreateGuestProfile(*ctx, dto, userId)
}

// GetProfileByFingerprintId retrieves a profile by fingerprint ID
func (s *profileService) GetProfileByFingerprintId(ctx context.Context, fingerprintId string) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.GetProfileByFingerprintId(ctx, fingerprintId)
}
