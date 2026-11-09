package profile_repos

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	profile_utils "larsa-e-commerce-microservice-v2/internal/modules/profile/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateGuestProfile creates a new guest profile
func (repo *ProfileRepo) CreateGuestProfile(ctx context.Context, dto *dtos.CreateGuestProfileDto, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError) {
	profile := profile_utils.ConvertToGuestProfileModel(dto, userId)
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	_, err := collection.InsertOne(ctx, profile)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	result, getErr := repo.GetById(ctx, profile.Id, "")
	if getErr != nil {
		return nil, getErr
	}
	return result, nil
}
