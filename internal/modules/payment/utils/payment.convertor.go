package payment_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
)

func ConvertToPaymentModel(dto *dtos.CreatePaymentDto, additionalOptions types.AdditionalOptions) *models.Payment {
	return &models.Payment{
		BasePayment: models.BasePayment{
			PaymentMethod:     dto.PaymentMethod,
			BankName:          dto.BankName,
			IBAN:              dto.IBAN,
			HolderName:        dto.HolderName,
			Status:            dto.Status,
			InvoiceId:         dto.InvoiceId,
			ReferenceNumber:   dto.ReferenceNumber,
			Amount:            dto.Amount,
			Account:           dto.Account,
			Note:              dto.Note,
			Attachments:       dto.Attachments,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
