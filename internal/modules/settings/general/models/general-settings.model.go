package settings_general_models

import (
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	settings_general_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/general/enums"
)

type FeedbackCouponSettings struct {
	IsEnabled          bool    `bson:"isEnabled" json:"isEnabled"`
	DiscountPercentage float64 `bson:"discountPercentage" json:"discountPercentage"`
	ExpirationDays     int     `bson:"expirationDays" json:"expirationDays"`
}

type BaseGeneralSettings struct {
	models.MetaDataBaseModel `bson:",inline"`

	SupportedLanguages []settings_general_enums.SupportedLanguage `bson:"supportedLanguages" json:"supportedLanguages"`
	FeedbackCoupon     FeedbackCouponSettings                     `bson:"feedbackCoupon" json:"feedbackCoupon"`
}

type GeneralSettings struct {
	BaseGeneralSettings `bson:",inline"`
}

type GeneralSettingsResponse struct {
	BaseGeneralSettings `bson:",inline"`
}
