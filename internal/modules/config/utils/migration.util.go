package config_utils

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ValidateMigrationRequest validates the migration request parameters
func ValidateMigrationRequest(fromDb, toDb, collectionName string) *helpers.ApiError {
	if fromDb == "" {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Source database name is required",
		}
	}

	if toDb == "" {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Target database name is required",
		}
	}

	if collectionName == "" {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Collection name is required",
		}
	}

	if fromDb == toDb {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "Source and target databases cannot be the same",
		}
	}

	return nil
}

// GetMigrationStats returns statistics about the migration operation
func GetMigrationStats(ctx context.Context, db *mongo.Database, collectionName string) (int64, error) {
	return db.Collection(collectionName).CountDocuments(ctx, bson.M{})
}

// LogMigrationOperation logs the migration operation details
func LogMigrationOperation(fromDb, toDb, collectionName string, documentCount int64, duration time.Duration) {
	fmt.Printf("Migration completed: %s.%s -> %s.%s (%d documents in %v)\n",
		fromDb, collectionName, toDb, collectionName, documentCount, duration)
}
