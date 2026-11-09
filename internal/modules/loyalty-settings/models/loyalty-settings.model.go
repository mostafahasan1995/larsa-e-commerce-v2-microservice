package loyalty_settings_models

import (
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BaseLoyaltySettings struct {
	models.MetaDataBaseModel `bson:",inline"`

	// Program Status
	IsActive bool `bson:"isActive" json:"isActive"`

	// Order-based Reward Configuration
	OrdersRequiredForReward    int64   `bson:"ordersRequiredForReward" json:"ordersRequiredForReward"`       // e.g., 30
	RewardPeriodMonths         int64   `bson:"rewardPeriodMonths" json:"rewardPeriodMonths"`                 // e.g., 3
	RewardDiscountPercentage   float64 `bson:"rewardDiscountPercentage" json:"rewardDiscountPercentage"`     // e.g., 10%
	RewardOrdersLimit          int64   `bson:"rewardOrdersLimit" json:"rewardOrdersLimit"`                   // e.g., 30
	RewardPeriodDurationMonths int64   `bson:"rewardPeriodDurationMonths" json:"rewardPeriodDurationMonths"` // e.g., 3
}

type LoyaltySettings struct {
	BaseLoyaltySettings `bson:",inline"`
}

type LoyaltySettingsResponse struct {
	BaseLoyaltySettings `bson:",inline"`
}
