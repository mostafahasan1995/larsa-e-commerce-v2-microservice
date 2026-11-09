package settings_general_repos

import (
	"context"
	"errors"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	settings_general_models "larsa-e-commerce-microservice-v2/internal/modules/settings/general/models"
	settings_general_services "larsa-e-commerce-microservice-v2/internal/modules/settings/general/services" // For the interface

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type generalSettingsRepository struct {
	*base_repos.BaseRepo[settings_general_models.GeneralSettings, settings_general_models.GeneralSettingsResponse]
}

func NewGeneralSettingsRepository(i *do.Injector) (settings_general_services.GeneralSettingsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	baseRepo := base_repos.NewBaseRepo[settings_general_models.GeneralSettings, settings_general_models.GeneralSettingsResponse](db, enums.GENERAL_SETTINGS)
	return &generalSettingsRepository{BaseRepo: baseRepo}, nil
}

func (r *generalSettingsRepository) Get(ctx context.Context) (*settings_general_models.GeneralSettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var settings settings_general_models.GeneralSettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// It's acceptable for settings to not exist initially.
			return nil, nil
		}
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get general settings: " + err.Error(),
		}
	}
	return &settings, nil
}

func (r *generalSettingsRepository) Update(ctx context.Context, model *settings_general_models.GeneralSettings) (*settings_general_models.GeneralSettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var existingSettings settings_general_models.GeneralSettings
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
			Message:    "Error fetching general settings for update: " + err.Error(),
		}
	}

	updatePayload := bson.M{"$set": model}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedSettingsDoc settings_general_models.GeneralSettings
	if err := collection.FindOneAndUpdate(ctx, filter, updatePayload, opts).Decode(&updatedSettingsDoc); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update or create general settings: " + err.Error(),
		}
	}

	return &updatedSettingsDoc, nil
}
