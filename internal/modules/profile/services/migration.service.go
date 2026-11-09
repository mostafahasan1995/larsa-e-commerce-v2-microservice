package profile_services

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MigrateUsersToProfilesDirect migrates user data using direct database access instead of gateway
func (s *profileService) MigrateUsersToProfilesDirect(ctx context.Context, limit int64, offset int64) (*types.MigrationResponse, *helpers.ApiError) {
	// Get database client for direct access to user collection
	cfg, _ := util.GetReqAppCfg(ctx)

	// Access user collection directly
	userCollection := s.db.Database(cfg.Db).Collection("users")

	findOptions := options.Find().SetLimit(limit).SetSkip(offset)
	// Find all users directly from database
	cursor, err := userCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get users from database: " + err.Error(),
		}
	}
	defer cursor.Close(ctx)

	// Decode users directly using our custom schema
	var users []profile_models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode users: " + err.Error(),
		}
	}

	if len(users) > 0 {
		fmt.Println("First user:", users[0].Id)
		fmt.Println("First user:", users[0].Email)
	}

	response := &types.MigrationResponse{
		TotalUsersProcessed: int64(len(users)),
		ProfilesCreated:     0,
		ProfilesSkipped:     0,
		Errors:              []string{},
		Message:             "Migration completed using direct database access",
	}

	for _, user := range users {
		if user.FirstName == "Guest" {
			response.ProfilesSkipped++
			continue
		}

		userId := user.Id
		existingProfile, _ := s.profileRepo.GetProfileByUserId(ctx, userId)

		if existingProfile != nil {
			response.ProfilesSkipped++
			continue
		}

		var profileErr *helpers.ApiError
		createAsVendor := user.UserType == "vendor"
		profileErr = s.createProfileFromCustomUser(ctx, &user, userId, createAsVendor)

		if profileErr != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Failed to create profile for user %s: %s", user.Email, profileErr.Message))
			response.ProfilesSkipped++
		} else {
			response.ProfilesCreated++
		}
	}

	return response, nil
}

func (s *profileService) createProfileFromCustomUser(ctx context.Context, user *profile_models.User, userId primitive.ObjectID, isVendor bool) *helpers.ApiError {
	userIdString := userId.Hex()

	if isVendor {
		// Create vendor profile
		vendorDto := &dtos.CreateVendorProfileDto{
			BaseProfileDto: dtos.BaseProfileDto{
				FirstName:   &user.FirstName,
				LastName:    &user.LastName,
				Email:       &user.Email,
				Gender:      &user.Gender,
				DateOfBirth: user.Birthday,
				UserId:      &userIdString,
			},
		}
		fmt.Println("Creating vendor profile")
		_, err := s.profileRepo.CreateVendorProfile(ctx, vendorDto)
		fmt.Println("vendor err", err)
		return err
	} else {
		// Create buyer profile
		buyerDto := &dtos.CreateBuyerProfileDto{
			BaseProfileDto: dtos.BaseProfileDto{
				FirstName:   &user.FirstName,
				LastName:    &user.LastName,
				Email:       &user.Email,
				Gender:      &user.Gender,
				DateOfBirth: user.Birthday,
				UserId:      &userIdString,
			},
		}
		fmt.Println("Creating buyer profile")
		_, err := s.profileRepo.CreateBuyerProfile(ctx, buyerDto)
		fmt.Println("buyer err", err)
		return err
	}
}
