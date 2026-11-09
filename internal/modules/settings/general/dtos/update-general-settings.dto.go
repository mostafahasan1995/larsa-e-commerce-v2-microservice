package settings_general_dtos

import (
	settings_general_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/general/enums"
	settings_general_models "larsa-e-commerce-microservice-v2/internal/modules/settings/general/models"
)

type UpdateGeneralSettingsDTO struct {
	SupportedLanguages []settings_general_enums.SupportedLanguage     `json:"supportedLanguages" validate:"required"`
	FeedbackCoupon     settings_general_models.FeedbackCouponSettings `json:"feedbackCoupon" validate:"required"`
}
