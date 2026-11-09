package settings_inventory_repos

import (
	"context"
	"errors"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	settings_inventory_models "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/models"
	settings_inventory_services "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory/services" // For the interface

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type inventorySettingsRepository struct {
	*base_repos.BaseRepo[settings_inventory_models.InventorySettings, settings_inventory_models.InventorySettingsResponse]
}

func NewInventorySettingsRepository(i *do.Injector) (settings_inventory_services.InventorySettingsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	baseRepo := base_repos.NewBaseRepo[settings_inventory_models.InventorySettings, settings_inventory_models.InventorySettingsResponse](db, enums.INVENTORY_SETTINGS)
	return &inventorySettingsRepository{BaseRepo: baseRepo}, nil
}

func (r *inventorySettingsRepository) Get(ctx context.Context) (*settings_inventory_models.InventorySettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var settings settings_inventory_models.InventorySettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// It's acceptable for settings to not exist initially.
			return nil, nil
		}
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get inventory settings: " + err.Error(),
		}
	}
	return &settings, nil
}

func (r *inventorySettingsRepository) Update(ctx context.Context, model *settings_inventory_models.InventorySettings) (*settings_inventory_models.InventorySettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var existingSettings settings_inventory_models.InventorySettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&existingSettings)

	filter := bson.M{}

	if err == nil {
		model.Id = existingSettings.Id
		filter = bson.M{"_id": existingSettings.Id}
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		filter = bson.M{}
	} else {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Error fetching inventory settings for update: " + err.Error(),
		}
	}

	updatePayload := bson.M{"$set": model}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedSettingsDoc settings_inventory_models.InventorySettings
	if err := collection.FindOneAndUpdate(ctx, filter, updatePayload, opts).Decode(&updatedSettingsDoc); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update or create inventory settings: " + err.Error(),
		}
	}

	return &updatedSettingsDoc, nil
}
