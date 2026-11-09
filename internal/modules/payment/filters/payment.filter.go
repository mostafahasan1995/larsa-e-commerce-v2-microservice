package payment_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	payment_enums "larsa-e-commerce-microservice-v2/internal/modules/payment/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type PaymentFilters struct {
	base_filters.BaseFilters
	PaymentMethod   *payment_enums.PaymentMethod `json:"paymentMethod,omitempty"`
	BankName        *string                      `json:"bankName,omitempty"`
	IBAN            *string                      `json:"iban,omitempty"`
	HolderName      *string                      `json:"holderName,omitempty"`
	Status          *payment_enums.PaymentStatus `json:"status,omitempty"`
	InvoiceID       *string                      `json:"invoiceId,omitempty"`
	ReferenceNumber *string                      `json:"referenceNumber,omitempty"`
	Account         *string                      `json:"account,omitempty"`
}

func NewPaymentFilters(query string) (*PaymentFilters, error) {
	paymentFilters := PaymentFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &paymentFilters); err != nil {
			return nil, err
		}
	}

	return &paymentFilters, nil
}

func (paymentFilters *PaymentFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if paymentFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &paymentFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PAYMENT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if paymentFilters.PaymentMethod != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"paymentMethod": paymentFilters.PaymentMethod}})
	}

	if paymentFilters.BankName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankName": paymentFilters.BankName}})
	}

	if paymentFilters.IBAN != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"iban": paymentFilters.IBAN}})
	}

	if paymentFilters.HolderName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"holderName": paymentFilters.HolderName}})
	}

	if paymentFilters.Status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": paymentFilters.Status}})
	}

	if paymentFilters.InvoiceID != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"invoiceId": paymentFilters.InvoiceID}})
	}

	if paymentFilters.ReferenceNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"referenceNumber": paymentFilters.ReferenceNumber}})
	}

	if paymentFilters.Account != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"account": paymentFilters.Account}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(paymentFilters.PopFields, enums.PAYMENT, pipeline)
		pipeline = base_filters.BuildSortingPipeline(paymentFilters.Sorts, pipeline)
	}

	return pipeline
}
