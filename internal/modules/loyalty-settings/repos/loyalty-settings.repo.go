package loyalty_settings_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	loyalty_settings_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/models"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type loyaltySettingsRepo struct {
	*base_repos.BaseRepo[loyalty_settings_models.LoyaltySettings, loyalty_settings_models.LoyaltySettingsResponse]
}

func NewLoyaltySettingsRepo(i *do.Injector) (loyalty_settings_services.LoyaltySettingsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &loyaltySettingsRepo{
		BaseRepo: base_repos.NewBaseRepo[loyalty_settings_models.LoyaltySettings, loyalty_settings_models.LoyaltySettingsResponse](
			db,
			enums.LOYALTY_SETTINGS,
		),
	}, nil
}

func (repo *loyaltySettingsRepo) DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError {
	return repo.BaseRepo.DeleteById(ctx, id, filter)
}

func (repo *loyaltySettingsRepo) GetActiveSettings(ctx context.Context) (*loyalty_settings_models.LoyaltySettingsResponse, *helpers.ApiError) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"isActive": true,
			},
		},
		{
			"$sort": bson.M{
				"createdAt": -1,
			},
		},
		{
			"$limit": 1,
		},
	}

	results, err := repo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 404,
			Message:    "No active loyalty settings found",
		}
	}

	return &results[0], nil
}
