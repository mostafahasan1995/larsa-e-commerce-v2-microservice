package cart_dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MigrateCartDto represents the request body for cart migration
// It contains the ID of the old user whose cart should be migrated to the current user
type MigrateCartDto struct {
	OldUserId primitive.ObjectID `json:"oldUserId" validate:"required"`
}
