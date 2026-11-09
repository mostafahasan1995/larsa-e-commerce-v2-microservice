package config_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type TicketTypeFilter struct {
	base_filters.BaseFilters
	Codes    []string `json:"codes,omitempty"`
	IsActive *bool    `json:"isActive,omitempty"`
}

func NewTicketTypeFilter(query string) (*TicketTypeFilter, error) {
	ticketTypeFilter := TicketTypeFilter{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &ticketTypeFilter); err != nil {
			return nil, err
		}
	}

	return &ticketTypeFilter, nil
}

func (ticketTypeFilter *TicketTypeFilter) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if ticketTypeFilter == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &ticketTypeFilter.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.TICKET_TYPE)

	if len(ticketTypeFilter.Codes) > 0 {
		matchStage["code"] = bson.M{
			"$in": ticketTypeFilter.Codes,
		}
	}

	if ticketTypeFilter.IsActive != nil {
		matchStage["isActive"] = *ticketTypeFilter.IsActive
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(ticketTypeFilter.PopFields, enums.TICKET_TYPE, pipeline)
		pipeline = base_filters.BuildSortingPipeline(ticketTypeFilter.Sorts, pipeline)
	}

	return pipeline
}
