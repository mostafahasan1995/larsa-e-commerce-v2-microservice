package loyalty_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_settings_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/models"
	loyalty_settings_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/services"
	loyalty_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	profile_enums "larsa-e-commerce-microservice-v2/internal/modules/profile/enums"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoyaltyOrderService interface {
	ProcessOrderForLoyaltyPoints(ctx context.Context, order *order_models.Order) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
	countOrderForLoyalty(ctx context.Context, userId primitive.ObjectID, orderId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
}

type LoyaltyOrderRepo interface {
	Create(ctx context.Context, loyaltyPoints loyalty_models.LoyaltyPoints) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
	GetByUserId(ctx context.Context, userId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, loyaltyPoints loyalty_models.LoyaltyPoints, filter map[string]interface{}) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError)
}

type loyaltyOrderService struct {
	repo                   LoyaltyOrderRepo
	profileService         profile_services.ProfileService
	loyaltySettingsService loyalty_settings_services.LoyaltySettingsService
}

func NewLoyaltyOrderService(i *do.Injector) (LoyaltyOrderService, error) {
	return &loyaltyOrderService{
		repo:                   do.MustInvoke[LoyaltyOrderRepo](i),
		profileService:         do.MustInvoke[profile_services.ProfileService](i),
		loyaltySettingsService: do.MustInvoke[loyalty_settings_services.LoyaltySettingsService](i),
	}, nil
}

func (s *loyaltyOrderService) ProcessOrderForLoyaltyPoints(ctx context.Context, order *order_models.Order) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	userId := order.BuyerInfo.Id
	if userId.IsZero() {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid user ID in order",
		}
	}

	profile, err := s.profileService.GetById(ctx, order.BuyerInfo.Id, "")
	if err == nil && profile != nil && profile.ProfileType != nil {
		if *profile.ProfileType == profile_enums.Guest {
			fmt.Println("config.User.UserData.FirstName is Guest")
			return nil, nil
		}
	} else {
		fmt.Println("profile not found", err)
	}

	// Check if order is paid (any payment type)
	if order.PaymentInfo.Status == order_enums.PAID ||
		order.PaymentInfo.Status == order_enums.PARTIALLY_PAID {

		// Count this order for loyalty tracking
		return s.countOrderForLoyalty(ctx, userId, order.Id)
	}

	return nil, &helpers.ApiError{
		StatusCode: 400,
		Message:    "Order payment not completed",
	}
}

func (s *loyaltyOrderService) countOrderForLoyalty(ctx context.Context, userId primitive.ObjectID, orderId primitive.ObjectID) (*loyalty_models.LoyaltyPointsResponse, *helpers.ApiError) {
	settings, err := s.loyaltySettingsService.GetActiveSettings(ctx)
	if err != nil {
		return nil, err
	}

	if !settings.IsActive {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Loyalty program is not active",
		}
	}

	// Get or create user loyalty record
	userLoyalty, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		if err.StatusCode == 404 {
			// Create new loyalty record
			userLoyalty = s.createNewLoyaltyRecord(ctx, userId, settings)
		} else {
			return nil, err
		}
	}

	// Check if we need to reset periods
	userLoyalty = s.checkAndResetPeriods(ctx, userLoyalty, settings)

	// Count this order
	userLoyalty.OrdersInCurrentPeriod++
	now := time.Now()
	userLoyalty.LastOrderDate = &now

	if !userLoyalty.HasEarnedReward && userLoyalty.OrdersInCurrentPeriod >= settings.OrdersRequiredForReward {
		userLoyalty = s.grantReward(ctx, userLoyalty, settings)
	} else if userLoyalty.HasEarnedReward && userLoyalty.RewardOrdersUsed >= settings.RewardOrdersLimit {
		userLoyalty = s.resetRewardCycle(ctx, userLoyalty, orderId, "limit reached")
	} else if userLoyalty.HasEarnedReward {
		userLoyalty.RewardOrdersUsed++
	}

	// Update history
	historyEntry := loyalty_models.LoyaltyPointsHistory{
		OrderId:     orderId,
		Description: "Order counted for loyalty",
		Timestamp:   now,
	}
	userLoyalty.History = append(userLoyalty.History, historyEntry)

	// Convert response to model for update
	cfg, _ := util.GetReqAppCfg(ctx)
	loyaltyPoints := loyalty_models.LoyaltyPoints{
		BaseLoyaltyPoints: userLoyalty.BaseLoyaltyPoints,
	}
	loyaltyPoints.Id = *userLoyalty.Id

	// Set metadata fields (fallback to buyer userId when request user is not set)
	userIdForMeta := userId
	if cfg != nil && cfg.User != nil && cfg.User.Id != primitive.NilObjectID {
		userIdForMeta = cfg.User.Id
	}
	metaData := util.CreateMetaDataBaseModel(types.AdditionalOptions{
		Mode:   types.UPDATE,
		UserId: userIdForMeta,
	})
	loyaltyPoints.MetaDataBaseModel = metaData

	// Save updated record
	return s.repo.UpdateById(ctx, *userLoyalty.Id, loyaltyPoints, bson.M{})
}

func (s *loyaltyOrderService) createNewLoyaltyRecord(ctx context.Context, userId primitive.ObjectID, settings *loyalty_settings_models.LoyaltySettingsResponse) *loyalty_models.LoyaltyPointsResponse {
	cfg, _ := util.GetReqAppCfg(ctx)
	now := time.Now()
	periodEnd := now.AddDate(0, int(settings.RewardPeriodMonths), 0)

	loyaltyPoints := loyalty_models.LoyaltyPoints{
		BaseLoyaltyPoints: loyalty_models.BaseLoyaltyPoints{
			UserId:                userId,
			OrdersInCurrentPeriod: 0,
			CurrentPeriodStart:    &now,
			CurrentPeriodEnd:      &periodEnd,
			HasEarnedReward:       false,
			RewardOrdersUsed:      0,
			History:               []loyalty_models.LoyaltyPointsHistory{},
		},
	}

	// Set metadata fields for create (fallback to buyer userId when request user is not set)
	userIdForMeta := userId
	if cfg != nil && cfg.User != nil && cfg.User.Id != primitive.NilObjectID {
		userIdForMeta = cfg.User.Id
	}
	metaData := util.CreateMetaDataBaseModel(types.AdditionalOptions{
		Mode:   types.CREATE,
		UserId: userIdForMeta,
	})
	loyaltyPoints.MetaDataBaseModel = metaData

	result, err := s.repo.Create(ctx, loyaltyPoints)
	if err != nil {
		return nil
	}

	return result
}

func (s *loyaltyOrderService) checkAndResetPeriods(ctx context.Context, userLoyalty *loyalty_models.LoyaltyPointsResponse, settings *loyalty_settings_models.LoyaltySettingsResponse) *loyalty_models.LoyaltyPointsResponse {
	now := time.Now()

	// Check if current period has expired
	if userLoyalty.CurrentPeriodEnd != nil && now.After(*userLoyalty.CurrentPeriodEnd) {
		// Reset current period
		userLoyalty.OrdersInCurrentPeriod = 0
		userLoyalty.CurrentPeriodStart = &now
		periodEnd := now.AddDate(0, int(settings.RewardPeriodMonths), 0)
		userLoyalty.CurrentPeriodEnd = &periodEnd
	}

	// Check if reward period has expired
	if userLoyalty.RewardEndDate != nil && now.After(*userLoyalty.RewardEndDate) {
		// Reset reward status
		userLoyalty.HasEarnedReward = false
		userLoyalty.RewardStartDate = nil
		userLoyalty.RewardEndDate = nil
		userLoyalty.RewardOrdersUsed = 0
	}

	return userLoyalty
}

func (s *loyaltyOrderService) grantReward(ctx context.Context, userLoyalty *loyalty_models.LoyaltyPointsResponse, settings *loyalty_settings_models.LoyaltySettingsResponse) *loyalty_models.LoyaltyPointsResponse {
	now := time.Now()

	userLoyalty.HasEarnedReward = true
	userLoyalty.RewardStartDate = &now
	rewardEnd := now.AddDate(0, int(settings.RewardPeriodDurationMonths), 0)
	userLoyalty.RewardEndDate = &rewardEnd
	userLoyalty.RewardOrdersUsed = 0

	return userLoyalty
}

func (s *loyaltyOrderService) resetRewardCycle(ctx context.Context, userLoyalty *loyalty_models.LoyaltyPointsResponse, orderId primitive.ObjectID, reason string) *loyalty_models.LoyaltyPointsResponse {
	now := time.Now()

	userLoyalty.HasEarnedReward = false
	userLoyalty.RewardStartDate = nil
	userLoyalty.RewardEndDate = nil
	userLoyalty.RewardOrdersUsed = 0
	userLoyalty.OrdersInCurrentPeriod = 0
	userLoyalty.CurrentPeriodStart = &now

	settings, err := s.loyaltySettingsService.GetActiveSettings(ctx)
	if err == nil && settings.IsActive {
		periodEnd := now.AddDate(0, int(settings.RewardPeriodMonths), 0)
		userLoyalty.CurrentPeriodEnd = &periodEnd
	}

	resetHistoryEntry := loyalty_models.LoyaltyPointsHistory{
		OrderId:     orderId,
		Description: fmt.Sprintf("Reward cycle reset - %s. Order counting restarted from zero.", reason),
		Timestamp:   now,
	}
	userLoyalty.History = append(userLoyalty.History, resetHistoryEntry)

	return userLoyalty
}
