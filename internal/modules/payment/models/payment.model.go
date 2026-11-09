package payment_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BasePayment struct {
	PaymentMethod   payment_enums.PaymentMethod `bson:"paymentMethod" json:"paymentMethod"`
	BankName        string                      `bson:"bankName" json:"bankName"`
	IBAN            string                      `bson:"iban" json:"iban"`
	HolderName      string                      `bson:"holderName" json:"holderName"`
	Status          payment_enums.PaymentStatus `bson:"status" json:"status"`
	InvoiceId       primitive.ObjectID          `bson:"invoiceId" json:"invoiceId"`
	ReferenceNumber string                      `bson:"referenceNumber" json:"referenceNumber"`
	Amount          float64                     `bson:"amount" json:"amount"`
	Account         string                      `bson:"account" json:"account"`
	Note            string                      `bson:"note,omitempty" json:"note,omitempty"`
	Attachments     []types.FileFieldDto        `bson:"attachments" json:"attachments"`

	base_models.MetaDataBaseModel `bson:",inline"`
}

type Payment struct {
	BasePayment `bson:",inline"`
}

type PaymentResponse struct {
	BasePayment `bson:",inline"`
}

type PaymentsWithPagination = types.DataWithPagination[PaymentResponse]
