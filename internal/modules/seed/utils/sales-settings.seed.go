package seed_utils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	seed_data "larsa-e-commerce-microservice-v2/internal/modules/seed/data"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SeedSalesSettings(ctx context.Context, db *mongo.Database, defaultUserId primitive.ObjectID) error {
	collection := db.Collection(enums.SALES_SETTINGS.GetCollectionName())

	seed_data.InitDefaultOrderStatuses()
	defaultOrderStatuses := seed_data.DefaultOrderStatuses

	now := time.Now()
	salesSettingsDoc := settings_sales_models.SalesSettings{
		BaseSalesSettings: settings_sales_models.BaseSalesSettings{
			MetaDataBaseModel: models.MetaDataBaseModel{
				UpdatedAt: now,
				UpdatedBy: &defaultUserId,
			},
			OrderStatuses:                defaultOrderStatuses,
			ShowVendorInfoInOrderDetails: false,
			OrderFlowMode:                settings_sales_enums.ORDER_FLOW_MODE_B2C,
			LogisticSettings: settings_sales_models.LogisticSettings{
				EnablePickup:   false,
				EnableDelivery: false,
				Pickup: &settings_sales_models.LogisticPickup{
					BranchName: false,
					PickUpDate: false,
					PickUpTime: false,
				},
				Delivery: &settings_sales_models.LogisticDelivery{
					ShipmentMethod:   false,
					TrackingNumber:   false,
					TransportCompany: false,
					ShipmentDate:     false,
					DeliveryFee:      false,
					CustomAddress:    false,
				},
			},
			RecentActivity: settings_sales_models.RecentActivity{
				EnableRecentActivity:     false,
				NumberOfProcessedOrders:  false,
				NumberOfIncompleteOrders: false,
				NumberOfCustomers:        false,
				UnpaidAmount:             false,
				OverdueAmount:            false,
			},
			EnableDiscountOnSales:           false,
			EnablePartialPayment:            false,
			EnableSalesNotes:                false,
			EnableAutoGenerateInvoiceNumber: false,
			EnableInvoiceEmailNotification:  false,
			InvoiceTerms: settings_sales_models.InvoiceTerms{
				DefaultInvoiceTerms:     30,
				DefaultInvoiceTermsUnit: settings_sales_enums.INVOICE_TERM_UNIT_DAYS,
			},
			SalesTaxesRate: 0.0,
		},
	}

	var existingSettings settings_sales_models.SalesSettings
	err := collection.FindOne(ctx, bson.M{}).Decode(&existingSettings)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check for existing sales settings: %v", err)
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.Id = primitive.NewObjectID()
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.CreatedAt = now
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.CreatedBy = &defaultUserId
	} else {
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.Id = existingSettings.Id
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.CreatedAt = existingSettings.CreatedAt
		salesSettingsDoc.BaseSalesSettings.MetaDataBaseModel.CreatedBy = existingSettings.CreatedBy
		salesSettingsDoc.OrderStatuses = defaultOrderStatuses
	}

	filter := bson.M{"_id": salesSettingsDoc.Id}
	update := bson.M{"$set": salesSettingsDoc}
	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to seed sales settings: %v", err)
	}

	fmt.Println("Successfully seeded sales settings with default order statuses.")
	return nil
}
