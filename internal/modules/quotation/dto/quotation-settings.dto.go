package quotation_dto

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
)

type QuotationProductAttributeSettingsDTO struct {
	Code    string                                        `bson:"code" json:"code"`
	Type    quotation_enums.QuotationProductAttributeType `bson:"type" json:"type"`
	Name    types.Localizable[string]                     `bson:"name" json:"name"`
	Choices []string                                      `bson:"choices" json:"choices"`
}

type QuotationSettingsDTO struct {
	Attributes       []QuotationProductAttributeSettingsDTO `json:"attributes"`
	ShippingMethods  []quotation_enums.ShippingMethod       `json:"shippingMethods"`
	PaymentMethods   []quotation_enums.PaymentMethod        `json:"paymentMethods"`
	PaymentSchedules []quotation_enums.PaymentSchedule      `json:"paymentSchedules"`
}
