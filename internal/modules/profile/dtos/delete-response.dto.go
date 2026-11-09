package profile_dtos

// DeleteGuestProfilesResponseDto represents the response for deleting guest profiles
type DeleteGuestProfilesResponseDto struct {
	DeletedCount int64  `json:"deletedCount"`
	Message      string `json:"message"`
}
