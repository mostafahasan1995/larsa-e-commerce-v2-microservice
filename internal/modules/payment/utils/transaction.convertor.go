package payment_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
)

func ConvertToTransactionModel(dto *dtos.CreateTransactionDto, additionalOptions types.AdditionalOptions) *models.Transaction {
	return &models.Transaction{
		BaseTransaction: models.BaseTransaction{
			Type:              dto.Type,
			Status:            dto.Status,
			Priority:          dto.Priority,
			Subject:           dto.Subject,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
