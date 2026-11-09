package loyalty_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	loyalty_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty/models"
	web_loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type webLoyaltyPointsRepo struct {
	*base_repos.BaseRepo[loyalty_models.LoyaltyPoints, loyalty_models.LoyaltyPointsResponse]
}

func NewWebLoyaltyPointsRepo(i *do.Injector) (web_loyalty_services.WebLoyaltyPointsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &webLoyaltyPointsRepo{
		BaseRepo: base_repos.NewBaseRepo[loyalty_models.LoyaltyPoints, loyalty_models.LoyaltyPointsResponse](
			db,
			enums.LOYALTY_POINTS,
		),
	}, nil
}

func (repo *webLoyaltyPointsRepo) GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userId}},
	}

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &loyalty_models.LoyaltyPointsResponse{
			BaseLoyaltyPoints: loyalty_models.BaseLoyaltyPoints{
				UserId:                userId,
				OrdersInCurrentPeriod: 0,
				CurrentPeriodStart:    nil,
				CurrentPeriodEnd:      nil,
				HasEarnedReward:       false,
				RewardStartDate:       nil,
				RewardEndDate:         nil,
				RewardOrdersUsed:      0,
				History:               []loyalty_models.LoyaltyPointsHistory{},
				LastOrderDate:         nil,
			},
		}, nil
	}

	return &results[0], nil
}

func (repo *webLoyaltyPointsRepo) UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter bson.M) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, loyaltyPoints, filter)
}
