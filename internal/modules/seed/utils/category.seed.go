package seed_utils

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	seeds_data "larsa-e-commerce-microservice-v2/internal/modules/seed/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedCategories(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.CATEGORY.GetCollectionName())
	seeds_data.InitCategories()

	var categoryIds []primitive.ObjectID
	for _, category := range seeds_data.Categories {
		categoryIds = append(categoryIds, category.Id)
	}

	// First, delete existing categories
	_, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": categoryIds}})
	if err != nil {
		return fmt.Errorf("failed to delete existing categories: %v", err)
	}

	// Convert categories to []interface{}
	docs := make([]interface{}, len(seeds_data.Categories))
	for i := range seeds_data.Categories {
		docs[i] = seeds_data.Categories[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %v", err)
	}

	return nil
}
