package payment_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
)

type BaseTransaction struct {
	Type     payment_enums.TransactionType     `bson:"type" json:"type"`
	Subject  types.Localizable[string]         `bson:"subject" json:"subject"`
	Status   payment_enums.TransactionStatus   `bson:"status" json:"status"`
	Priority payment_enums.TransactionPriority `bson:"priority" json:"priority"`

	base_models.MetaDataBaseModel `bson:",inline"`
}
type Transaction struct {
	BaseTransaction `bson:",inline"`
}

type TransactionResponse struct {
	BaseTransaction `bson:",inline"`
}

type TransactionsWithPagination = types.DataWithPagination[TransactionResponse]
