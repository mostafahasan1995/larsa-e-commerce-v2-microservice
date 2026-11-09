package unserved_delivery_request_repos

import (
	"context"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	unserved_delivery_request_dtos "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/dtos"
	unserved_delivery_request_filters "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/filters"
	unserved_delivery_request_models "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request/models"
)

type UnservedDeliveryRequestRepo struct {
	*base_repos.BaseRepo[unserved_delivery_request_models.UnservedDeliveryRequest, unserved_delivery_request_models.UnservedDeliveryRequestResponse]
}

func NewUnservedDeliveryRequestRepo(i *do.Injector) (*UnservedDeliveryRequestRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &UnservedDeliveryRequestRepo{
		BaseRepo: base_repos.NewBaseRepo[unserved_delivery_request_models.UnservedDeliveryRequest, unserved_delivery_request_models.UnservedDeliveryRequestResponse](db, enums.UNSERVED_DELIVERY_REQUEST),
	}, nil
}

func (repo *UnservedDeliveryRequestRepo) Create(ctx context.Context, model unserved_delivery_request_models.UnservedDeliveryRequest) (*unserved_delivery_request_models.UnservedDeliveryRequestResponse, *helpers.ApiError) {
	if model.Id == primitive.NilObjectID {
		model.Id = primitive.NewObjectID()
	}
	now := time.Now()
	model.CreatedAt = now
	model.UpdatedAt = now

	_, err := repo.BaseRepo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	var response unserved_delivery_request_models.UnservedDeliveryRequestResponse
	modelBytes, _ := bson.Marshal(model)
	_ = bson.Unmarshal(modelBytes, &response)
	return &response, nil
}

func (repo *UnservedDeliveryRequestRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) (*types.DataWithPagination[unserved_delivery_request_models.UnservedDeliveryRequestResponse], *helpers.ApiError) {
	matchStage := bson.M{}
	pipeline := filter.BuildPipeline(matchStage, ctx, false)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

// GetAll returns all records matching the filter (for export). Uses same pipeline as list with sorting.
func (repo *UnservedDeliveryRequestRepo) GetAll(ctx context.Context, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) ([]unserved_delivery_request_models.UnservedDeliveryRequestResponse, *helpers.ApiError) {
	matchStage := bson.M{}
	pipeline := filter.BuildPipeline(matchStage, ctx, true)
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *UnservedDeliveryRequestRepo) AggregateByLocationAndDate(ctx context.Context, groupBy string, filter *unserved_delivery_request_filters.UnservedDeliveryRequestFilters) ([]unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateItem, *helpers.ApiError) {
	collection, apiErr := repo.GetCollection(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	cfg, _ := util.GetReqAppCfg(ctx)
	if cfg == nil || cfg.Db == "" {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "missing app config in context"}
	}

	// Build match stage from filters
	matchStage := bson.M{}
	base_utils.AddSoftDeletedFilter(matchStage)

	if filter.Country != nil && *filter.Country != "" {
		matchStage["country"] = bson.M{"$regex": *filter.Country, "$options": "i"}
	}
	if filter.City != nil && *filter.City != "" {
		matchStage["city"] = bson.M{"$regex": *filter.City, "$options": "i"}
	}
	if filter.Area != nil && *filter.Area != "" {
		matchStage["area"] = bson.M{"$regex": *filter.Area, "$options": "i"}
	}
	if filter.CreatedAt != nil {
		dateCond := bson.M{}
		if filter.CreatedAt.Gte != nil {
			dateCond["$gte"] = *filter.CreatedAt.Gte
		}
		if filter.CreatedAt.Lte != nil {
			dateCond["$lte"] = *filter.CreatedAt.Lte
		}
		if len(dateCond) > 0 {
			matchStage["createdAt"] = dateCond
		}
	}

	// Build group stage based on groupBy
	var groupId bson.M
	switch groupBy {
	case "country":
		groupId = bson.M{"country": bson.M{"$ifNull": bson.A{"$country", "unknown"}}}
	case "city":
		groupId = bson.M{
			"country": bson.M{"$ifNull": bson.A{"$country", "unknown"}},
			"city":    bson.M{"$ifNull": bson.A{"$city", "unknown"}},
		}
	case "area":
		groupId = bson.M{
			"country": bson.M{"$ifNull": bson.A{"$country", "unknown"}},
			"city":    bson.M{"$ifNull": bson.A{"$city", "unknown"}},
			"area":    bson.M{"$ifNull": bson.A{"$area", "unknown"}},
		}
	case "date":
		groupId = bson.M{
			"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
		}
	default:
		groupId = bson.M{"date": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}}}
	}

	pipeline := []bson.M{
		{"$match": matchStage},
		{"$group": bson.M{
			"_id":   groupId,
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	defer cursor.Close(ctx)

	var rawResults []struct {
		Id    bson.M `bson:"_id"`
		Count int64  `bson:"count"`
	}
	if err := cursor.All(ctx, &rawResults); err != nil {
		return nil, helpers.DefaultError(err)
	}

	results := make([]unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateItem, 0, len(rawResults))
	for _, r := range rawResults {
		item := unserved_delivery_request_dtos.UnservedDeliveryRequestAggregateItem{Count: r.Count}
		if c, ok := r.Id["country"].(string); ok {
			item.Country = &c
		}
		if c, ok := r.Id["city"].(string); ok {
			item.City = &c
		}
		if a, ok := r.Id["area"].(string); ok {
			item.Area = &a
		}
		if d, ok := r.Id["date"].(string); ok {
			item.Date = d
		}
		results = append(results, item)
	}

	return results, nil
}
