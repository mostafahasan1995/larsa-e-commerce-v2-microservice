package payment_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
)

func ConvertToInvoiceModel(dto *dtos.CreateInvoiceDto, additionalOptions types.AdditionalOptions) *models.Invoice {
	return &models.Invoice{
		BaseInvoice: models.BaseInvoice{
			InvoiceNumber:     dto.InvoiceNumber,
			Date:              dto.Date,
			User:              dto.User,
			Company:           dto.Company,
			Region:            dto.Region,
			Details:           dto.Details,
			Status:            dto.Status,
			PaymentStatus:     dto.PaymentStatus,
			Amount:            dto.Amount,
			Discount:          dto.Discount,
			PaidAmount:        dto.PaidAmount,
			DueAmount:         dto.DueAmount,
			PaymentMethod:     dto.PaymentMethod,
			TransactionId:     dto.TransactionId,
			ReferenceNumber:   dto.ReferenceNumber,
			AccountManager:    dto.AccountManager,
			SalesRep:          dto.SalesRep,
			SalesSource:       dto.SalesSource,
			Currency:          dto.Currency,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
