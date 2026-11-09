package profile_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBuyerProfile creates a new buyer profile
func (s *profileService) CreateBuyerProfile(ctx context.Context, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.CreateBuyerProfile(ctx, dto)
}

// UpdateBuyerProfile updates an existing buyer profile
func (s *profileService) UpdateBuyerProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateBuyerProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.UpdateBuyerProfile(ctx, id, dto)
}
