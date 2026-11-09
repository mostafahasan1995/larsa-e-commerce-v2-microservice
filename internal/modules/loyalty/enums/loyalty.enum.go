package loyalty_enums

type LoyaltyActionType string

const (
	POINTS_EARNED  LoyaltyActionType = "points_earned"
	POINTS_SPENT   LoyaltyActionType = "points_spent"
	REWARD_GRANTED LoyaltyActionType = "reward_granted"
	REWARD_EXPIRED LoyaltyActionType = "reward_expired"
)

type LoyaltyStatus string

const (
	ACTIVE   LoyaltyStatus = "active"
	INACTIVE LoyaltyStatus = "inactive"
	PENDING  LoyaltyStatus = "pending"
	EXPIRED  LoyaltyStatus = "expired"
)
