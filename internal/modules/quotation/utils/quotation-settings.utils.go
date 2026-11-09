package quotation_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
)

// GetDefaultQuotationSettings returns the default quotation settings configuration
func GetDefaultQuotationSettings() quotation_dto.QuotationSettingsDTO {
	return quotation_dto.QuotationSettingsDTO{
		Attributes: []quotation_dto.QuotationProductAttributeSettingsDTO{
			{
				Code: "material",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_STRING_ARRAY,
				Name: types.Localizable[string]{
					enums.EN: "Material",
					enums.AR: "المادة",
					enums.TR: "Malzeme",
				},
				Choices: []string{"Plastic", "Metal", "Wood", "Paper", "Glass", "Other"},
			},
			{
				Code: "colors",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_STRING_ARRAY,
				Name: types.Localizable[string]{
					enums.EN: "Colors",
					enums.AR: "اللون",
					enums.TR: "Renk",
				},
				Choices: []string{"Red", "Blue", "Green", "Black", "White", "Yellow"},
			},
			{
				Code: "quality_standards",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_STRING_ARRAY,
				Name: types.Localizable[string]{
					enums.EN: "Quality Standards",
					enums.AR: "المعايير الجودة",
					enums.TR: "Kalite Standartları",
				},
				Choices: []string{"ISO3000", "GMP", "HACCP", "FDA", "HALAL", "Kosher"},
			},
			{
				Code: "specific_features",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_TEXT,
				Name: types.Localizable[string]{
					enums.EN: "Specific Features",
					enums.AR: "الميزات المحددة",
					enums.TR: "Özel Özellikler",
				},
				Choices: []string{},
			},
			{
				Code: "length",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER,
				Name: types.Localizable[string]{
					enums.EN: "Length",
					enums.AR: "الطول",
					enums.TR: "Uzunluk",
				},
				Choices: []string{},
			},
			{
				Code: "width",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER,
				Name: types.Localizable[string]{
					enums.EN: "Width",
					enums.AR: "العرض",
					enums.TR: "Genişlik",
				},
				Choices: []string{},
			},
			{
				Code: "height",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER,
				Name: types.Localizable[string]{
					enums.EN: "Height",
					enums.AR: "الإرتفاع",
					enums.TR: "Yükseklik",
				},
				Choices: []string{},
			},
			{
				Code: "weight",
				Type: quotation_enums.QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER,
				Name: types.Localizable[string]{
					enums.EN: "Weight",
					enums.AR: "الوزن",
					enums.TR: "Ağırlık",
				},
				Choices: []string{},
			},
		},
		ShippingMethods: []quotation_enums.ShippingMethod{
			quotation_enums.SEA,
			quotation_enums.LAND,
			quotation_enums.INTERNATIONAL_AIR,
			quotation_enums.LOCAL_AIR,
		},
		PaymentMethods: []quotation_enums.PaymentMethod{
			quotation_enums.BANK_TRANSFER,
			quotation_enums.CREDIT_CARD,
			quotation_enums.CASH,
			quotation_enums.LETTER_OF_CREDIT,
		},
		PaymentSchedules: []quotation_enums.PaymentSchedule{
			quotation_enums.FINAL_PAYMENT_UPON_SHIPMENT,
			quotation_enums.ADVANCE_PAYMENT,
			quotation_enums.INSTALLMENTS,
		},
	}
}
