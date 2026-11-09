package watchlist_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	watchlist_enums "larsa-e-commerce-microservice-v2/internal/modules/watchlist/enums"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WatchListFilters struct {
	// Base filters (common filters like date ranges, etc.)
	filters.BaseFilters

	// WatchList specific filters
	UserId   *primitive.ObjectID            `json:"userId,omitempty"`
	ItemId   *primitive.ObjectID            `json:"itemId,omitempty"`
	ItemType *watchlist_enums.WatchItemType `json:"itemType,omitempty"`
}

func NewWatchListFilters(query string) (WatchListFilters, error) {
	f := WatchListFilters{}
	if query == "" {
		return f, nil
	}

	if err := json.Unmarshal([]byte(query), &f); err != nil {
		return f, err
	}

	return f, nil
}

// BuildMatchFilter creates a bson.M document representing the $match stage
// based on the filter fields.
func (watchListFilters *WatchListFilters) BuildMatchFilter(ctx context.Context, matchStage bson.M) bson.M {

	// Apply base filters that affect the match stage
	watchListFilters.ApplyBaseFilters(ctx, matchStage, enums.WATCH_LIST)

	if watchListFilters.UserId != nil {
		matchStage["userId"] = *watchListFilters.UserId
	}

	if watchListFilters.ItemId != nil {
		matchStage["itemId"] = *watchListFilters.ItemId
	}

	if watchListFilters.ItemType != nil {
		matchStage["itemType"] = *watchListFilters.ItemType
	}

	return matchStage
}

func (watchListFilters *WatchListFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if watchListFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Build the core match filter using the dedicated method
	builtMatchFilter := watchListFilters.BuildMatchFilter(ctx, matchStage)

	if len(builtMatchFilter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": builtMatchFilter})
	}

	if addExtraOperations {
		pipeline = filters.BuildPopulatePipeline(watchListFilters.PopFields, enums.WATCH_LIST, pipeline)
		pipeline = filters.BuildSortingPipeline(watchListFilters.Sorts, pipeline)
	}

	return pipeline
}
