package profile_repos

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

// DeleteAllGuestProfiles deletes all profiles where firstName equals "Guest"
func (repo *ProfileRepo) DeleteAllGuestProfiles(ctx context.Context) (int64, *helpers.ApiError) {
	collection, _ := repo.BaseRepo.GetCollection(ctx)

	// Filter to find all profiles with firstName = "Guest"
	filter := bson.M{
		"firstName": "Guest",
	}

	// Delete all matching documents
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, helpers.DefaultError(err)
	}

	return result.DeletedCount, nil
}
