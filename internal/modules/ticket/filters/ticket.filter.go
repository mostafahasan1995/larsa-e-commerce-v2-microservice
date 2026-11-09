package ticket_filters

import (
	"context"
	"encoding/json"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	ticket_enums "larsa-e-commerce-microservice-v2/internal/modules/ticket/enums"
	ticket_utils "larsa-e-commerce-microservice-v2/internal/modules/ticket/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketFilters struct {
	filters.BaseFilters

	Status      *ticket_enums.TicketStatus   `bson:"status,omitempty" json:"status,omitempty"`
	Priority    *ticket_enums.TicketPriority `bson:"priority,omitempty" json:"priority,omitempty"`
	TypeId      *primitive.ObjectID          `bson:"typeId,omitempty" json:"typeId,omitempty"`
	RequesterId *primitive.ObjectID          `bson:"requesterId,omitempty" json:"requesterId,omitempty"`
	BuyerId     *primitive.ObjectID          `bson:"buyerId,omitempty" json:"buyerId,omitempty"`
	VendorId    *primitive.ObjectID          `bson:"vendorId,omitempty" json:"vendorId,omitempty"`

	QuotationId *primitive.ObjectID `bson:"quotationId,omitempty" json:"quotationId,omitempty"`
	InvoiceId   *primitive.ObjectID `bson:"invoiceId,omitempty" json:"invoiceId,omitempty"`
	OrderId     *primitive.ObjectID `bson:"orderId,omitempty" json:"orderId,omitempty"`
}

func NewTicketFilters(query string) (*TicketFilters, error) {
	ticketFilters := TicketFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &ticketFilters); err != nil {
			return nil, err
		}
	}
	return &ticketFilters, nil
}

func (ticketFilters *TicketFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if ticketFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	baseFilters := &ticketFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.TICKET)

	if ticketFilters.Status != nil {
		matchStage["status"] = *ticketFilters.Status
	}

	if ticketFilters.Priority != nil {
		matchStage["priority"] = *ticketFilters.Priority
	}

	if ticketFilters.TypeId != nil {
		matchStage["typeId"] = *ticketFilters.TypeId
	}

	if ticketFilters.RequesterId != nil {
		matchStage["requesterId"] = *ticketFilters.RequesterId
	}

	if ticketFilters.BuyerId != nil {
		matchStage["buyerId"] = *ticketFilters.BuyerId
	}

	if ticketFilters.VendorId != nil {
		matchStage["vendorId"] = *ticketFilters.VendorId
	}

	if ticketFilters.QuotationId != nil {
		matchStage["quotationId"] = *ticketFilters.QuotationId
	}

	if ticketFilters.InvoiceId != nil {
		matchStage["invoiceId"] = *ticketFilters.InvoiceId
	}

	if ticketFilters.OrderId != nil {
		matchStage["orderId"] = *ticketFilters.OrderId
	}

	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = filters.BuildPopulatePipeline(ticketFilters.PopFields, enums.TICKET, pipeline)
		pipeline = filters.BuildSortingPipeline(ticketFilters.Sorts, pipeline)
		pipeline = append(pipeline, ticket_utils.TicketPopulateCommentsUsers...)
	}

	return pipeline
}
