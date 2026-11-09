package seed_utils

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	seeds_data "larsa-e-commerce-microservice-v2/internal/modules/seed/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedProducts(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.PRODUCT.GetCollectionName())

	// First, delete existing products
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to delete existing products: %v", err)
	}

	// Convert seed data to products models
	products := make([]product_models.Product, 0, len(seeds_data.Products))
	for _, seed := range seeds_data.Products {
		product := product_utils.ConvertToProductModel(&seed, types.AdditionalOptions{
			UserId: defaultUserId,
			Mode:   types.CREATE,
		})

		products = append(products, *product)
	}

	// Convert products to []interface{}
	docs := make([]interface{}, len(products))
	for i := range products {
		docs[i] = products[i]
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed products: %v", err)
	}

	return nil
}
