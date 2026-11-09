package loyalty_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoyaltyPointsHistory struct {
	OrderId     primitive.ObjectID `bson:"orderId" json:"orderId"`
	Description string             `bson:"description" json:"description"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}

type BaseLoyaltyPoints struct {
	models.MetaDataBaseModel `bson:",inline"`
	UserId                   primitive.ObjectID `bson:"userId" json:"userId"`

	// Order Counting
	OrdersInCurrentPeriod int64 `bson:"ordersInCurrentPeriod" json:"ordersInCurrentPeriod"`

	// Period Tracking
	CurrentPeriodStart *time.Time `bson:"currentPeriodStart,omitempty" json:"currentPeriodStart,omitempty"`
	CurrentPeriodEnd   *time.Time `bson:"currentPeriodEnd,omitempty" json:"currentPeriodEnd,omitempty"`

	// Reward Status
	HasEarnedReward  bool       `bson:"hasEarnedReward" json:"hasEarnedReward"`
	RewardStartDate  *time.Time `bson:"rewardStartDate,omitempty" json:"rewardStartDate,omitempty"`
	RewardEndDate    *time.Time `bson:"rewardEndDate,omitempty" json:"rewardEndDate,omitempty"`
	RewardOrdersUsed int64      `bson:"rewardOrdersUsed" json:"rewardOrdersUsed"`

	// History
	History       []LoyaltyPointsHistory `bson:"history" json:"history"`
	LastOrderDate *time.Time             `bson:"lastOrderDate,omitempty" json:"lastOrderDate,omitempty"`
}

type LoyaltyPoints struct {
	BaseLoyaltyPoints `bson:",inline"`
}

type LoyaltyPointsResponse struct {
	models.OptionalMetaDataBaseModel `bson:",inline"`
	BaseLoyaltyPoints                `bson:",inline"`
}

type LoyaltyPointsWithPagination = types.DataWithPagination[LoyaltyPointsResponse]
