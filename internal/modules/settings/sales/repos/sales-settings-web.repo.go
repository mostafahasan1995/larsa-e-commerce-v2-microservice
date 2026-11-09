package settings_sales_repos

import (
	"context"
	"errors"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SalesSettingsWebRepository struct {
	*base_repos.BaseRepo[settings_sales_models.SalesSettings, settings_sales_models.SalesSettingsResponse]
}

func NewSalesSettingsWebRepository(i *do.Injector) (settings_sales_services.SalesSettingsWebRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &SalesSettingsWebRepository{
		BaseRepo: base_repos.NewBaseRepo[settings_sales_models.SalesSettings, settings_sales_models.SalesSettingsResponse](db, enums.SALES_SETTINGS),
	}, nil
}

func (r *SalesSettingsWebRepository) Get(ctx context.Context) (*settings_sales_models.SalesSettings, *helpers.ApiError) {
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
