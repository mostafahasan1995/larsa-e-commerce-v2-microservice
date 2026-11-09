package seed_utils

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	seed_location_data "larsa-e-commerce-microservice-v2/internal/modules/seed/data/location"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedCountries(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.SEED_COUNTRY.GetCollectionName())
	seed_location_data.InitCountries()

	// First, delete existing countries
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete existing countries: %v", err)
	}

	// Convert countries to []interface{}
	if len(seed_location_data.Countries) == 0 {
		fmt.Println("No countries to seed.")
		return nil // Nothing to insert
	}
	docs := make([]interface{}, len(seed_location_data.Countries))
	for i := range seed_location_data.Countries {
		docs[i] = seed_location_data.Countries[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed countries: %v", err)
	}

	fmt.Println("Successfully seeded countries.")
	return nil
}

func SeedCities(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.SEED_CITY.GetCollectionName())
	seed_location_data.InitCities()

	// First, delete existing cities
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete existing cities: %v", err)
	}

	// Convert cities to []interface{}
	if len(seed_location_data.Cities) == 0 {
		fmt.Println("No cities to seed.")
		return nil // Nothing to insert
	}
	docs := make([]interface{}, len(seed_location_data.Cities))
	for i := range seed_location_data.Cities {
		docs[i] = seed_location_data.Cities[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed cities: %v", err)
	}

	fmt.Println("Successfully seeded cities.")
	return nil
}

func SeedAreas(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.SEED_AREA.GetCollectionName())
	seed_location_data.InitAreas()

	// First, delete existing areas
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete existing areas: %v", err)
	}

	// Convert areas to []interface{}
	if len(seed_location_data.Areas) == 0 {
		fmt.Println("No areas to seed.")
		return nil // Nothing to insert
	}
	docs := make([]interface{}, len(seed_location_data.Areas))
	for i := range seed_location_data.Areas {
		docs[i] = seed_location_data.Areas[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed areas: %v", err)
	}

	fmt.Println("Successfully seeded areas.")
	return nil
}
