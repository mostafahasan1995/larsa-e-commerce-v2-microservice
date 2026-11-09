package config_repos

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MigrationRepo struct {
	db *mongo.Client
}

func NewMigrationRepo(i *do.Injector) (services.MigrationRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &MigrationRepo{
		db: db,
	}, nil
}

func (repo *MigrationRepo) MigrateData(ctx context.Context, dto config_dtos.DataMigrationDto) *helpers.ApiError {
	startTime := time.Now()

	// Validate request parameters
	if validationErr := config_utils.ValidateMigrationRequest(dto.FromDatabase, dto.ToDatabase, dto.CollectionName); validationErr != nil {
		return validationErr
	}

	// Validate that both databases exist
	fromDb := repo.db.Database(dto.FromDatabase)
	toDb := repo.db.Database(dto.ToDatabase)

	// Check if source collection exists and has data
	fromCollection := fromDb.Collection(dto.CollectionName)
	count, err := fromCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("Error counting documents in source collection: %v", err),
		}
	}

	if count == 0 {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    fmt.Sprintf("Source collection '%s' in database '%s' is empty", dto.CollectionName, dto.FromDatabase),
		}
	}

	// Get all documents from source collection
	cursor, err := fromCollection.Find(ctx, bson.M{}, &options.FindOptions{})
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("Error reading from source collection: %v", err),
		}
	}
	defer cursor.Close(ctx)

	// Prepare documents for insertion
	var documents []interface{}
	if err = cursor.All(ctx, &documents); err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("Error decoding documents: %v", err),
		}
	}

	if len(documents) == 0 {
		return &helpers.ApiError{
			StatusCode: 400,
			Message:    "No documents found to migrate",
		}
	}

	// Clear target collection if it exists
	toCollection := toDb.Collection(dto.CollectionName)
	_, err = toCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("Error clearing target collection: %v", err),
		}
	}

	// Insert documents into target collection
	if len(documents) > 0 {
		_, err = toCollection.InsertMany(ctx, documents)
		if err != nil {
			return &helpers.ApiError{
				StatusCode: 500,
				Message:    fmt.Sprintf("Error inserting documents into target collection: %v", err),
			}
		}
	}

	// Log migration operation
	duration := time.Since(startTime)
	config_utils.LogMigrationOperation(dto.FromDatabase, dto.ToDatabase, dto.CollectionName, count, duration)

	return nil
}
