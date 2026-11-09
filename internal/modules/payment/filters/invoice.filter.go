package payment_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type InvoiceFilters struct {
	base_filters.BaseFilters
	InvoiceNumber   *string                             `json:"invoiceNumber,omitempty"`
	Date            *time.Time                          `json:"date,omitempty"`
	User            *string                             `json:"user,omitempty"`
	Company         *string                             `json:"company,omitempty"`
	Region          *string                             `json:"region,omitempty"`
	Status          *payment_enums.InvoiceStatus        `json:"status,omitempty"`
	PaymentStatus   *payment_enums.InvoicePaymentStatus `json:"paymentStatus,omitempty"`
	PaymentMethod   *string                             `json:"paymentMethod,omitempty"`
	TransactionId   *string                             `json:"transactionId,omitempty"`
	ReferenceNumber *string                             `json:"referenceNumber,omitempty"`
	AccountManager  *string                             `json:"accountManager,omitempty"`
	SalesRep        *string                             `json:"salesRep,omitempty"`
	Currency        *string                             `json:"currency,omitempty"`
}

func NewInvoiceFilters(query string) (*InvoiceFilters, error) {
	invoiceFilters := InvoiceFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &invoiceFilters); err != nil {
			return nil, err
		}
	}

	return &invoiceFilters, nil
}

func (invoiceFilters *InvoiceFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if invoiceFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &invoiceFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.INVOICE)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if invoiceFilters.InvoiceNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"invoiceNumber": invoiceFilters.InvoiceNumber}})
	}

	if invoiceFilters.Date != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"date": invoiceFilters.Date}})
	}

	if invoiceFilters.User != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"user": invoiceFilters.User}})
	}

	if invoiceFilters.Company != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"company": invoiceFilters.Company}})
	}

	if invoiceFilters.Region != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"region": invoiceFilters.Region}})
	}

	if invoiceFilters.Status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": invoiceFilters.Status}})
	}

	if invoiceFilters.PaymentStatus != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"paymentStatus": invoiceFilters.PaymentStatus}})
	}

	if invoiceFilters.PaymentMethod != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"paymentMethod": invoiceFilters.PaymentMethod}})
	}

	if invoiceFilters.TransactionId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"transactionId": invoiceFilters.TransactionId}})
	}

	if invoiceFilters.ReferenceNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"referenceNumber": invoiceFilters.ReferenceNumber}})
	}

	if invoiceFilters.AccountManager != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"accountManager": invoiceFilters.AccountManager}})
	}

	if invoiceFilters.SalesRep != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"salesRep": invoiceFilters.SalesRep}})
	}

	if invoiceFilters.Currency != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"currency": invoiceFilters.Currency}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(invoiceFilters.PopFields, enums.INVOICE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(invoiceFilters.Sorts, pipeline)
	}

	return pipeline
}
