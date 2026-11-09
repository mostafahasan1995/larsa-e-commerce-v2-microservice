package payment_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
)

type BasePayout struct {
	Amount            float64 `bson:"amount" json:"amount"`
	BankName          string  `bson:"bankName" json:"bankName"`
	BankAccountNumber string  `bson:"bankAccountNumber" json:"bankAccountNumber"`
	AccountHolderName string  `bson:"accountHolderName" json:"accountHolderName"`
	BankBranchAddress string  `bson:"bankBranchAddress" json:"bankBranchAddress"`
	SwiftCode         string  `bson:"swiftCode" json:"swiftCode"`
	TransferReason    string  `bson:"transferReason" json:"transferReason"`
	Notes             string  `bson:"notes" json:"notes"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Payout struct {
	BasePayout `bson:",inline"`
}

type PayoutResponse struct {
	BasePayout `bson:",inline"`
}

type PayoutsWithPagination = types.DataWithPagination[PayoutResponse]
