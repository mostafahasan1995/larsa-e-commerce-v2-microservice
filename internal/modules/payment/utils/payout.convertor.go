package payment_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/payment/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
)

func ConvertToPayoutModel(dto *dtos.CreatePayoutDto, additionalOptions types.AdditionalOptions) *models.Payout {

	return &models.Payout{
		BasePayout: models.BasePayout{
			BankName:          dto.BankName,
			BankAccountNumber: dto.BankAccountNumber,
			AccountHolderName: dto.AccountHolderName,
			BankBranchAddress: dto.BankBranchAddress,
			SwiftCode:         dto.SwiftCode,
			TransferReason:    dto.TransferReason,
			Notes:             dto.Notes,
			Amount:            dto.Amount,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
}
