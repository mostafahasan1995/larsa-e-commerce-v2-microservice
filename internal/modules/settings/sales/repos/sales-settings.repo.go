package settings_sales_repos

import (
	"context"
	"errors"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services" // For the interface

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type salesSettingsRepository struct {
	*base_repos.BaseRepo[settings_sales_models.SalesSettings, settings_sales_models.SalesSettingsResponse]
}

func NewSalesSettingsRepository(i *do.Injector) (settings_sales_services.SalesSettingsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	baseRepo := base_repos.NewBaseRepo[settings_sales_models.SalesSettings, settings_sales_models.SalesSettingsResponse](db, enums.SALES_SETTINGS)
	return &salesSettingsRepository{BaseRepo: baseRepo}, nil
}

func (r *salesSettingsRepository) Get(ctx context.Context) (*settings_sales_models.SalesSettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var settings settings_sales_models.SalesSettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&settings)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// It's acceptable for settings to not exist initially.
			return nil, nil
		}
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get sales settings: " + err.Error(),
		}
	}
	return &settings, nil
}

func (r *salesSettingsRepository) Update(ctx context.Context, model *settings_sales_models.SalesSettings) (*settings_sales_models.SalesSettings, *helpers.ApiError) {
	collection, collErr := r.GetCollection(ctx)
	if collErr != nil {
		return nil, collErr
	}

	var existingSettings settings_sales_models.SalesSettings
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
			Message:    "Error fetching sales settings for update: " + err.Error(),
		}
	}

	updatePayload := bson.M{"$set": model}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedSettingsDoc settings_sales_models.SalesSettings
	if err := collection.FindOneAndUpdate(ctx, filter, updatePayload, opts).Decode(&updatedSettingsDoc); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update or create sales settings: " + err.Error(),
		}
	}

	return &updatedSettingsDoc, nil
}
