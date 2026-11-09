package payment_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseInvoice struct {
	InvoiceNumber   string                             `bson:"invoiceNumber" json:"invoiceNumber"`
	Date            time.Time                          `bson:"date" json:"date"`
	User            string                             `bson:"user" json:"user"`
	Company         string                             `bson:"company" json:"company"`
	Region          string                             `bson:"region" json:"region"`
	Details         string                             `bson:"details" json:"details"`
	Status          payment_enums.InvoiceStatus        `bson:"status" json:"status"`
	PaymentStatus   payment_enums.InvoicePaymentStatus `bson:"paymentStatus" json:"paymentStatus"`
	Amount          float64                            `bson:"amount" json:"amount"`
	Discount        float64                            `bson:"discount" json:"discount"`
	PaidAmount      float64                            `bson:"paidAmount" json:"paidAmount"`
	DueAmount       float64                            `bson:"dueAmount" json:"dueAmount"`
	PaymentMethod   string                             `bson:"paymentMethod" json:"paymentMethod"`
	TransactionId   primitive.ObjectID                 `bson:"transactionId" json:"transactionId"`
	ReferenceNumber string                             `bson:"referenceNumber" json:"referenceNumber"`
	AccountManager  string                             `bson:"accountManager" json:"accountManager"`
	SalesSource     string                             `bson:"salesSource" json:"salesSource"`
	SalesRep        string                             `bson:"salesRep" json:"salesRep"`
	Currency        string                             `bson:"currency" json:"currency"`

	base_models.MetaDataBaseModel `bson:",inline"`
}

type Invoice struct {
	BaseInvoice `bson:",inline"`
}

type InvoiceResponse struct {
	BaseInvoice `bson:",inline"`
}

type InvoicesWithPagination = types.DataWithPagination[InvoiceResponse]
