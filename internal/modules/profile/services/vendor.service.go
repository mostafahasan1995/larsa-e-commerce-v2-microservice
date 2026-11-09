package profile_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/profile/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateVendorProfile creates a new vendor profile
func (s *profileService) CreateVendorProfile(ctx context.Context, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.CreateVendorProfile(ctx, dto)
}

// UpdateVendorProfile updates an existing vendor profile
func (s *profileService) UpdateVendorProfile(ctx context.Context, id primitive.ObjectID, dto *dtos.CreateVendorProfileDto) (*models.ProfileResponse, *helpers.ApiError) {
	return s.profileRepo.UpdateVendorProfile(ctx, id, dto)
}
