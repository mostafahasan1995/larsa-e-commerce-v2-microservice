package payment_dtos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreatePaymentDto struct {
	PaymentMethod   payment_enums.PaymentMethod `json:"paymentMethod" validate:"required,oneof=bank_transfer"`
	BankName        string                      `json:"bankName" validate:"required"`
	IBAN            string                      `json:"iban" validate:"required"`
	HolderName      string                      `json:"holderName" validate:"required"`
	Status          payment_enums.PaymentStatus `json:"status" validate:"required,oneof=pending approved rejected completed"`
	InvoiceId       primitive.ObjectID          `json:"invoiceId" validate:"required"`
	ReferenceNumber string                      `json:"referenceNumber" validate:"required"`
	Amount          float64                     `json:"amount" validate:"required,gt=0"`
	Account         string                      `json:"account" validate:"required"`
	Note            string                      `json:"note,omitempty"`
	Attachments     []types.FileFieldDto        `json:"attachments" validate:"required,min=1,dive"`
}

func (dto *CreatePaymentDto) CustomValidator(ctx context.Context, db *mongo.Client) *helpers.ApiError {
	invoiceRepo := base_repos.NewBaseRepo[payment_models.Invoice, payment_models.InvoiceResponse](db, enums.INVOICE)
	_, apiErr := invoiceRepo.GetById(ctx, dto.InvoiceId, []bson.M{})
	if apiErr != nil {
		return apiErr
	}

	return nil
}
