package profile_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	profile_filters "larsa-e-commerce-microservice-v2/internal/modules/profile/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetById retrieves a profile by its ID
func (s *profileService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.GetById(ctx, id, query)
}

// GetWithPagination retrieves profiles with pagination
func (s *profileService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProfilesWithPagination, *helpers.ApiError) {
	return s.profileRepo.GetWithPagination(ctx, skip, limit, query)
}

// GetAll retrieves all profiles
func (s *profileService) GetAll(ctx context.Context, query string) ([]models.ProfileResponse, *helpers.ApiError) {
	profileFilters, err := profile_filters.NewProfileFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}
	return s.profileRepo.GetAll(ctx, profileFilters.BuildPipeline(bson.M{}, ctx, false))
}

// GetProfileByUserId retrieves a profile by user ID
func (s *profileService) GetProfileByUserId(ctx context.Context, userId primitive.ObjectID) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.GetProfileByUserId(ctx, userId)
}

// DeleteById deletes a profile by its ID
func (s *profileService) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.profileRepo.DeleteById(ctx, id)
}

// BulkDelete deletes multiple profiles by their IDs
func (s *profileService) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return s.profileRepo.BulkDelete(ctx, ids)
}
