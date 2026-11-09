package payment_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
)

type CreateTransactionDto struct {
	Type     payment_enums.TransactionType     `json:"type" validate:"required,oneof=complaints request inquiry"`
	Subject  types.Localizable[string]         `json:"subject" validate:"required,localizable"`
	Status   payment_enums.TransactionStatus   `json:"status" validate:"required,oneof=closed pending draft in_progress open"`
	Priority payment_enums.TransactionPriority `json:"priority" validate:"required,oneof=high medium low"`
}
