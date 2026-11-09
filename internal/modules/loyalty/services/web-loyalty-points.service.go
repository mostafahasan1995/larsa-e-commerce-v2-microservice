package loyalty_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	loyalty_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebLoyaltyPointsRepo interface {
	GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter bson.M) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
}

type WebLoyaltyPointsService interface {
	GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter bson.M) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
}

type webLoyaltyPointsService struct {
	repo WebLoyaltyPointsRepo
}

func NewWebLoyaltyPointsService(i *do.Injector) (WebLoyaltyPointsService, error) {
	return &webLoyaltyPointsService{
		repo: do.MustInvoke[WebLoyaltyPointsRepo](i),
	}, nil
}

func (s *webLoyaltyPointsService) GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	return s.repo.GetByUserId(ctx, userId)
}

func (s *webLoyaltyPointsService) UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter bson.M) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	return s.repo.UpdateById(ctx, id, loyaltyPoints, filter)
}
