package contact_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	customer_enums "larsa-e-commerce-microservice-v2/internal/modules/contact/enums"

	"go.mongodb.org/mongo-driver/bson"
)

type CustomerFilters struct {
	base_filters.BaseFilters
	Status            *customer_enums.CustomerStatus `json:"status,omitempty"`
	ProfileLink       *string                        `json:"profileLink,omitempty"`
	Email             *bool                          `json:"email,omitempty"`
	SMS               *bool                          `json:"sms,omitempty"`
	Whatsapp          *bool                          `json:"whatsapp,omitempty"`
	Dashboard         *bool                          `json:"dashboard,omitempty"`
	Plan              *string                        `json:"plan,omitempty"`
	UnlimitedDuration *bool                          `json:"unlimitedDuration,omitempty"`
}

func NewCustomerFilters(query string) (*CustomerFilters, error) {
	customerFilters := CustomerFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &customerFilters); err != nil {
			return nil, err
		}
	}
	return &customerFilters, nil
}

func (customerFilters *CustomerFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if customerFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &customerFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.CUSTOMER)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if customerFilters.Status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"basicInformation.status": customerFilters.Status}})
	}

	if customerFilters.ProfileLink != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"storeDetails.profileLink": customerFilters.ProfileLink}})
	}

	if customerFilters.Email != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"notificationMethods.email": customerFilters.Email}})
	}

	if customerFilters.SMS != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"notificationMethods.sms": customerFilters.SMS}})
	}

	if customerFilters.Whatsapp != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"notificationMethods.whatsapp": customerFilters.Whatsapp}})
	}

	if customerFilters.Dashboard != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"notificationMethods.dashboard": customerFilters.Dashboard}})
	}

	if customerFilters.Plan != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"subscription.plan": customerFilters.Plan}})
	}

	if customerFilters.UnlimitedDuration != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"subscription.unlimitedDuration": customerFilters.UnlimitedDuration}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(baseFilters.PopFields, enums.CUSTOMER, pipeline)

		pipeline = base_filters.BuildSortingPipeline(baseFilters.Sorts, pipeline)
	}

	return pipeline
}
