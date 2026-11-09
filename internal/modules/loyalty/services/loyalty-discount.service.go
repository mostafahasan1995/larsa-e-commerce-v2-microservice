package loyalty_services

import (
	"context"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"
	loyalty_dtos "larsa-e-commerce-microservice-v2/internal/modules/loyalty/dtos"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoyaltyDiscountService interface {
	CheckDiscountEligibility(ctx context.Context, userId primitive.ObjectID) (*loyalty_dtos.DiscountEligibility, error)
}

type loyaltyDiscountService struct {
	loyaltySettingsService loyalty_settings_services.LoyaltySettingsService
	loyaltyPointsService   WebLoyaltyPointsService
}

func NewLoyaltyDiscountService(i *do.Injector) (LoyaltyDiscountService, error) {
	return &loyaltyDiscountService{
		loyaltySettingsService: do.MustInvoke[loyalty_settings_services.LoyaltySettingsService](i),
		loyaltyPointsService:   do.MustInvoke[WebLoyaltyPointsService](i),
	}, nil
}

func (s *loyaltyDiscountService) CheckDiscountEligibility(ctx context.Context, userId primitive.ObjectID) (*loyalty_dtos.DiscountEligibility, error) {
	settings, err := s.loyaltySettingsService.GetActiveSettings(ctx)
	if err != nil {
		return nil, err
	}

	userLoyalty, err := s.loyaltyPointsService.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Check if user has earned reward and is in reward period
	if !userLoyalty.HasEarnedReward || userLoyalty.RewardEndDate == nil || time.Now().After(*userLoyalty.RewardEndDate) {
		return &loyalty_dtos.DiscountEligibility{IsEligible: false}, nil
	}

	// Check if user hasn't exceeded reward order limit
	if userLoyalty.RewardOrdersUsed >= settings.RewardOrdersLimit {
		return &loyalty_dtos.DiscountEligibility{IsEligible: false}, nil
	}

	return &loyalty_dtos.DiscountEligibility{
		IsEligible:         true,
		DiscountPercentage: settings.RewardDiscountPercentage,
		OrdersRemaining:    settings.RewardOrdersLimit - userLoyalty.RewardOrdersUsed,
		PeriodEndDate:      userLoyalty.RewardEndDate,
	}, nil
}
