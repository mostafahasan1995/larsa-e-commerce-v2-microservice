package loyalty_settings_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	loyalty_settings_dtos "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/dtos"
	loyalty_settings_models "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings/models"
)

func ConvertCreateLoyaltySettingsDTOToModel(dto *loyalty_settings_dtos.CreateLoyaltySettingsDTO, additionalOptions types.AdditionalOptions) loyalty_settings_models.LoyaltySettings {
	return loyalty_settings_models.LoyaltySettings{
		BaseLoyaltySettings: loyalty_settings_models.BaseLoyaltySettings{
			MetaDataBaseModel:          util.CreateMetaDataBaseModel(additionalOptions),
			IsActive:                   dto.IsActive,
			OrdersRequiredForReward:    dto.OrdersRequiredForReward,
			RewardPeriodMonths:         dto.RewardPeriodMonths,
			RewardDiscountPercentage:   dto.RewardDiscountPercentage,
			RewardOrdersLimit:          dto.RewardOrdersLimit,
			RewardPeriodDurationMonths: dto.RewardPeriodDurationMonths,
		},
	}
}
