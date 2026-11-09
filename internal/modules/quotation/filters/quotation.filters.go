package quotation_filters

import (
	"context"
	"encoding/json"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationFilters struct {
	base_filters.BaseFilters
	Type     *quotation_enums.QuotationType `json:"type,omitempty"`
	Deadline *time.Time                     `json:"deadline,omitempty"`

	Code            *string                           `json:"code,omitempty"`
	Status          *quotation_enums.QuotationStatus  `json:"status,omitempty"`
	ExcludeStatuses []quotation_enums.QuotationStatus `json:"excludeStatuses,omitempty"`

	//
	VendorProfileId    *primitive.ObjectID `json:"vendorProfileId,omitempty"`
	BuyerProfileId     *primitive.ObjectID `json:"buyerProfileId,omitempty"`
	RequestQuotationId *primitive.ObjectID `json:"requestQuotationId,omitempty"`
	GeneralQuotationId *primitive.ObjectID `json:"generalQuotationId,omitempty"`

	// Party filters
	AgentFromId       *primitive.ObjectID `json:"agentFromId,omitempty"`
	AgentToId         *primitive.ObjectID `json:"agentToId,omitempty"`
	ShippingCompanyId *primitive.ObjectID `json:"shippingCompanyId,omitempty"`

	// Any party user id (matches any of agentFromId, agentToId, or shippingCompanyId)
	AnyPartyUserId *primitive.ObjectID `json:"anyPartyUserId,omitempty"`
}

func NewQuotationFilters(query string) (*QuotationFilters, error) {
	quotationFilters := QuotationFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &quotationFilters); err != nil {
			return nil, err
		}
	}
	return &quotationFilters, nil
}

func (quotationFilters *QuotationFilters) BuildPipeline(ctx context.Context, matchStage bson.M, addExtraOperations bool) []bson.M {
	if quotationFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &quotationFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.QUOTATION)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// Apply type filter
	if quotationFilters.Type != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": *quotationFilters.Type}})
	}

	// Apply deadline filter
	if quotationFilters.Deadline != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"deadline": *quotationFilters.Deadline}})
	}

	// Apply code filter
	if quotationFilters.Code != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"code": *quotationFilters.Code}})
	}

	// Apply status filter
	if quotationFilters.Status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": *quotationFilters.Status}})
	}

	// Apply exclude statuses filter
	if len(quotationFilters.ExcludeStatuses) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": bson.M{"$nin": quotationFilters.ExcludeStatuses}}})
	}

	// Apply vendor profile id filter
	if quotationFilters.VendorProfileId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"vendorProfileId": *quotationFilters.VendorProfileId}})
	}

	// Apply buyer profile id filter
	if quotationFilters.BuyerProfileId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"quotationBuyerInfo.buyerProfileId": *quotationFilters.BuyerProfileId}})
	}

	// Apply general quotation id filter
	if quotationFilters.GeneralQuotationId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"generalQuotationId": *quotationFilters.GeneralQuotationId}})
	}

	// Apply request quotation id filter
	if quotationFilters.RequestQuotationId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"requestQuotationId": *quotationFilters.RequestQuotationId}})
	}

	// Apply agent/shipping specific filters
	if quotationFilters.AgentFromId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"agentFromId": *quotationFilters.AgentFromId}})
	}
	if quotationFilters.AgentToId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"agentToId": *quotationFilters.AgentToId}})
	}
	if quotationFilters.ShippingCompanyId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"shippingCompanyId": *quotationFilters.ShippingCompanyId}})
	}

	// Apply "any party" OR filter
	if quotationFilters.AnyPartyUserId != nil {
		id := quotationFilters.AnyPartyUserId
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$or": []bson.M{
			{"agentFromId": *id},
			{"agentToId": *id},
			{"shippingCompanyId": *id},
		}}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(quotationFilters.PopFields, enums.QUOTATION, pipeline)
		pipeline = base_filters.BuildSortingPipeline(quotationFilters.Sorts, pipeline)
		pipeline = append(pipeline, quotation_utils.QuotationCalculateSubmittedQuotationCount...)
		//TODO: check if we need to populate products
		pipeline = append(pipeline, quotation_utils.QuotationPopulateProducts...)
	}

	return pipeline
}
