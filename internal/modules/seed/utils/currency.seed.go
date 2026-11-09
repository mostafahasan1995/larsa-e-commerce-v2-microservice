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

func SeedCurrencies(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.CURRENCY.GetCollectionName())
	seeds_data.InitCurrencies()

	// First, delete currencies
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete existing currencies: %v", err)
	}

	// Convert countries to []interface{}
	docs := make([]interface{}, len(seeds_data.Currencies))
	for i := range seeds_data.Currencies {
		docs[i] = seeds_data.Currencies[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed currencies: %v", err)
	}

	return nil
}
