package appointment_filters

import (
	"context"
	"encoding/json"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	appointment_enums "larsa-e-commerce-microservice-v2/internal/modules/appointment/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppointmentFilters defines filters for querying appointments
type AppointmentFilters struct {
	base_filters.BaseFilters
	Id        *primitive.ObjectID                  `json:"id" form:"id"`
	VendorId  *primitive.ObjectID                  `json:"vendorId" form:"vendorId"`
	BuyerId   *primitive.ObjectID                  `json:"buyerId" form:"buyerId"`
	Status    *appointment_enums.AppointmentStatus `json:"status" form:"status"`
	StartDate *string                              `json:"startDate" form:"startDate"`
	EndDate   *string                              `json:"endDate" form:"endDate"`
	SortBy    string                               `json:"sortBy" form:"sortBy"`
	SortOrder string                               `json:"sortOrder" form:"sortOrder"`
}

func NewAppointmentFilters(query string) (*AppointmentFilters, error) {
	if query == "" {
		return &AppointmentFilters{}, nil
	}

	var filters AppointmentFilters
	if err := json.Unmarshal([]byte(query), &filters); err != nil {
		return nil, err
	}

	return &filters, nil
}

func (appointmentFilters *AppointmentFilters) BuildPipeline(pipeline []bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	matchStage := bson.M{}

	// Apply base filters
	baseFilters := &appointmentFilters.BaseFilters

	if ctx != nil {
		baseFilters.ApplyBaseFilters(ctx, matchStage, enums.APPOINTMENT)
	} else {
		matchStage["deletedAt"] = bson.M{
			"$exists": false,
		}
	}

	if appointmentFilters.Id != nil {
		matchStage["_id"] = *appointmentFilters.Id
	}

	if appointmentFilters.VendorId != nil {
		matchStage["vendorId"] = *appointmentFilters.VendorId
	}

	if appointmentFilters.BuyerId != nil {
		matchStage["buyerId"] = *appointmentFilters.BuyerId
	}

	if appointmentFilters.Status != nil {
		matchStage["status"] = *appointmentFilters.Status
	}

	if appointmentFilters.StartDate != nil && appointmentFilters.EndDate != nil {
		matchStage["date"] = bson.M{
			"$gte": *appointmentFilters.StartDate,
			"$lte": *appointmentFilters.EndDate,
		}
	} else if appointmentFilters.StartDate != nil {
		matchStage["date"] = bson.M{
			"$gte": *appointmentFilters.StartDate,
		}
	} else if appointmentFilters.EndDate != nil {
		matchStage["date"] = bson.M{
			"$lte": *appointmentFilters.EndDate,
		}
	}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.M{
			"$match": matchStage,
		})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(appointmentFilters.PopFields, enums.APPOINTMENT, pipeline)

		pipeline = base_filters.BuildSortingPipeline(appointmentFilters.Sorts, pipeline)

	}
	return pipeline
}
