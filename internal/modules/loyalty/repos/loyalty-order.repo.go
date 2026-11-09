package loyalty_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	loyalty_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty/models"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type loyaltyOrderRepo struct {
	*base_repos.BaseRepo[loyalty_models.LoyaltyPoints, loyalty_models.LoyaltyPointsResponse]
}

func NewLoyaltyOrderRepo(i *do.Injector) (loyalty_services.LoyaltyOrderRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &loyaltyOrderRepo{
		BaseRepo: base_repos.NewBaseRepo[loyalty_models.LoyaltyPoints, loyalty_models.LoyaltyPointsResponse](
			db,
			enums.LOYALTY_POINTS,
		),
	}, nil
}

func (repo *loyaltyOrderRepo) Create(ctx context.Context, loyaltyPoints loyalty_models.LoyaltyPoints) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, loyaltyPoints)
}

func (repo *loyaltyOrderRepo) GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"userId": userId,
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
			Message:    "Loyalty points not found for user",
		}
	}

	return &results[0], nil
}

func (repo *loyaltyOrderRepo) UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter map[string]interface{}) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, loyaltyPoints, bson.M(filter))
}
