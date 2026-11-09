package seed_utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	settings_general_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/general/enums"
	settings_general_models "larsa-e-commerce-microservice-v2/internal/modules/settings/general/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SeedGeneralSettings(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.GENERAL_SETTINGS.GetCollectionName())

	now := time.Now()
	generalSettingsDoc := settings_general_models.GeneralSettings{
		BaseGeneralSettings: settings_general_models.BaseGeneralSettings{
			MetaDataBaseModel: models.MetaDataBaseModel{
				UpdatedAt: now,
				UpdatedBy: &defaultUserId,
			},
			SupportedLanguages: []settings_general_enums.SupportedLanguage{
				settings_general_enums.SUPPORTED_LANGUAGE_ENGLISH,
				settings_general_enums.SUPPORTED_LANGUAGE_ARABIC,
			},
			FeedbackCoupon: settings_general_models.FeedbackCouponSettings{
				IsEnabled:          true,
				DiscountPercentage: 13.0,
				ExpirationDays:     2,
			},
		},
	}

	var existingSettings settings_general_models.GeneralSettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&existingSettings)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check for existing general settings: %v", err)
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.Id = primitive.NewObjectID()
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.CreatedAt = now
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.CreatedBy = &defaultUserId
	} else {
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.Id = existingSettings.Id
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.CreatedAt = existingSettings.CreatedAt
		generalSettingsDoc.BaseGeneralSettings.MetaDataBaseModel.CreatedBy = existingSettings.CreatedBy
	}

	filter := bson.M{"_id": generalSettingsDoc.Id}
	update := bson.M{"$set": generalSettingsDoc}
	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to seed general settings: %v", err)
	}

	fmt.Println("Successfully seeded general settings.")
	return nil
}
