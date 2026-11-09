package loyalty_settings_dtos

// Not
type CreateLoyaltySettingsDTO struct {
	IsActive                   bool    `json:"isActive" validate:"required"`
	OrdersRequiredForReward    int64   `json:"ordersRequiredForReward" validate:"required,min=1"`
	RewardPeriodMonths         int64   `json:"rewardPeriodMonths" validate:"required,min=1"`
	RewardDiscountPercentage   float64 `json:"rewardDiscountPercentage" validate:"required,min=0,max=100"`
	RewardOrdersLimit          int64   `json:"rewardOrdersLimit" validate:"required,min=1"`
	RewardPeriodDurationMonths int64   `json:"rewardPeriodDurationMonths" validate:"required,min=1"`
}
