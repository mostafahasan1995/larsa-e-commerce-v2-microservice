package order_repos

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderWebRepo interface {
	GetById(ctx context.Context, id primitive.ObjectID, filter *order_filters.OrderFilters) (*order_models.OrderResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter order_filters.OrderFilters) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError)
}

type orderWebRepository struct {
	*base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
}

func NewOrderWebRepository(i *do.Injector) (OrderWebRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &orderWebRepository{
		BaseRepo: base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
	}, nil
}

func (repo *orderWebRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter order_filters.OrderFilters) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError) {
	pipeline, postPipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		PostPipeline:   postPipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *orderWebRepository) GetById(ctx context.Context, id primitive.ObjectID, filter *order_filters.OrderFilters) (*order_models.OrderResponse, *helpers.ApiError) {
	pipeline, postPipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetById(ctx, id, append(pipeline, postPipeline...))
}

func (repo *orderWebRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	update := bson.M{}

	for key, value := range updateFields {
		if key[0] == '$' {
			update[key] = value
			delete(updateFields, key)
		}
	}

	if len(updateFields) > 0 {
		updateFields["updatedAt"] = time.Now()
		updateFields["updatedBy"] = userId
		update["$set"] = updateFields
	}

	result, updateErr := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if updateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update order: " + updateErr.Error(),
		}
	}

	if result.MatchedCount == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 404,
			Message:    "Order not found",
		}
	}

	return repo.GetById(ctx, id, nil)
}
