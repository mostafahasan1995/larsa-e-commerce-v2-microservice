package config_utils

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetChildCategoryIds retrieves all child category IDs for the given category IDs
// This includes the original category IDs plus all their descendants
func GetChildCategoryIds(ctx context.Context, categoryIds []primitive.ObjectID) ([]primitive.ObjectID, *helpers.ApiError) {
	if len(categoryIds) == 0 {
		return []primitive.ObjectID{}, nil
	}

	// Get the database client from the injector
	i := util.GetInjector(ctx)
	db := do.MustInvoke[*mongo.Client](i)

	// Get configuration to access database name
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get app configuration: " + err.Error(),
		}
	}

	collection := db.Database(cfg.Db).Collection(enums.CATEGORY.GetCollectionName())

	// Build the pipeline to get all child category IDs
	pipeline := GetCategoryChildrenIdsPipeline(categoryIds)

	// Execute the aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get child category IDs: " + err.Error(),
		}
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode child category IDs: " + err.Error(),
		}
	}

	// If no results, return the original category IDs
	if len(results) == 0 {
		return categoryIds, nil
	}

	// Extract category IDs from the result
	var allCategoryIds []primitive.ObjectID
	if result := results[0]; result != nil {
		if categoryIdsInterface, exists := result["categoryIds"]; exists {
			if categoryIdsArray, ok := categoryIdsInterface.(primitive.A); ok {
				for _, id := range categoryIdsArray {
					if objectId, ok := id.(primitive.ObjectID); ok {
						allCategoryIds = append(allCategoryIds, objectId)
					}
				}
			}
		}
	}

	// If no child IDs found, return the original category IDs
	if len(allCategoryIds) == 0 {
		return categoryIds, nil
	}

	return allCategoryIds, nil
}
