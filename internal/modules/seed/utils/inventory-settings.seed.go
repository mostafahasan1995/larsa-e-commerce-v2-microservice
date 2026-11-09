package seed_utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	settings_inventory_models "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SeedInventorySettings(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.INVENTORY_SETTINGS.GetCollectionName())

	now := time.Now()
	inventorySettingsDoc := settings_inventory_models.InventorySettings{
		BaseInventorySettings: settings_inventory_models.BaseInventorySettings{
			MetaDataBaseModel: models.MetaDataBaseModel{
				UpdatedAt: now,
				UpdatedBy: &defaultUserId,
			},
			UnitsActive: false,
		},
	}

	var existingSettings settings_inventory_models.InventorySettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&existingSettings)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check for existing inventory settings: %v", err)
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.Id = primitive.NewObjectID()
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.CreatedAt = now
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.CreatedBy = &defaultUserId
	} else {
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.Id = existingSettings.Id
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.CreatedAt = existingSettings.CreatedAt
		inventorySettingsDoc.BaseInventorySettings.MetaDataBaseModel.CreatedBy = existingSettings.CreatedBy
	}

	filter := bson.M{"_id": inventorySettingsDoc.Id}
	update := bson.M{"$set": inventorySettingsDoc}
	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to seed inventory settings: %v", err)
	}

	fmt.Println("Successfully seeded inventory settings.")
	return nil
}
