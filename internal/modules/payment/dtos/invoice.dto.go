package payment_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
	payment_models "larsa-e-commerce-microservice-v2/internal/modules/payment/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInvoiceDto struct {
	InvoiceNumber   string                             `json:"invoiceNumber" validate:"required"`
	Date            time.Time                          `json:"date" validate:"required"`
	User            string                             `json:"user" validate:"required"`
	Company         string                             `json:"company" validate:"required"`
	Region          string                             `json:"region" validate:"required"`
	Details         string                             `json:"details" validate:"required"`
	Status          payment_enums.InvoiceStatus        `json:"status" validate:"required,oneof=draft pending approved rejected completed"`
	PaymentStatus   payment_enums.InvoicePaymentStatus `json:"paymentStatus" validate:"required,oneof=pending paid partially overdue"`
	Amount          float64                            `json:"amount" validate:"required,gt=0"`
	Discount        float64                            `json:"discount" validate:"gte=0"`
	PaidAmount      float64                            `json:"paidAmount" validate:"gte=0"`
	DueAmount       float64                            `json:"dueAmount" validate:"gte=0"`
	PaymentMethod   string                             `json:"paymentMethod" validate:"required"`
	TransactionId   primitive.ObjectID                 `json:"transactionId" validate:"required"`
	ReferenceNumber string                             `json:"referenceNumber" validate:"required"`
	AccountManager  string                             `json:"accountManager" validate:"required"`
	SalesRep        string                             `json:"salesRep" validate:"required"`
	SalesSource     string                             `json:"salesSource" validate:"required"`
	Currency        string                             `json:"currency" validate:"required"`
}

func (dto *CreateInvoiceDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	transactionRepo := base_repos.NewBaseRepo[payment_models.Transaction, payment_models.TransactionResponse](args.DB, enums.TRANSACTION)
	_, apiErr := transactionRepo.GetById(args.Context, dto.TransactionId, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}

	return nil, nil
}
