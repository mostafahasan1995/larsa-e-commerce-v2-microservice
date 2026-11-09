package profile_services

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
)

// DeleteAllGuestProfiles deletes all guest profiles (where firstName equals "Guest")
func (s *profileService) DeleteAllGuestProfiles(ctx context.Context) (int64, *helpers.ApiError) {
	return s.profileRepo.DeleteAllGuestProfiles(ctx)
}
