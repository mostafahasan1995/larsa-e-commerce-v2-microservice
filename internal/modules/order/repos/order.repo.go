package order_repos

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepo interface {
	Create(ctx context.Context, order order_models.Order) (*order_models.OrderResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, filter *order_filters.OrderFilters) (*order_models.OrderResponse, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]order_models.OrderResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter order_filters.OrderFilters) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, order order_models.Order, filter bson.M) (*order_models.OrderResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}, userId primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError)
	GetRecentActivity(ctx context.Context, userId primitive.ObjectID, userType enums.UserType, branchId *primitive.ObjectID, filter order_filters.OrderFilters) (map[string]float64, *helpers.ApiError)
	GetBuyerActivity(ctx context.Context, buyerId primitive.ObjectID, branchId *primitive.ObjectID, filter order_filters.OrderFilters) (*order_dto.BuyerActivityResponse, *helpers.ApiError)
	DeleteAllOrders(ctx context.Context) *helpers.ApiError
}

type orderRepository struct {
	*base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
}

func NewOrderRepository(i *do.Injector) (OrderRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &orderRepository{
		BaseRepo: base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
	}, nil
}

// TODO: just for testing, should be removed
func (repo *orderRepository) DeleteAllOrders(ctx context.Context) *helpers.ApiError {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get collection: " + err.Error(),
		}
	}

	_, deleteErr := collection.DeleteMany(ctx, bson.M{})
	if deleteErr != nil {
		return &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to delete all orders: " + deleteErr.Error(),
		}
	}
	return nil
}

func (repo *orderRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter order_filters.OrderFilters) (*types.DataWithPagination[order_models.OrderResponse], *helpers.ApiError) {
	pipeline, postPipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetWithOptimizedPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		PostPipeline:   postPipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *orderRepository) GetById(ctx context.Context, id primitive.ObjectID, filter *order_filters.OrderFilters) (*order_models.OrderResponse, *helpers.ApiError) {
	pipeline, postPipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetById(ctx, id, append(pipeline, postPipeline...))
}

func (repo *orderRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*order_models.OrderResponse, *helpers.ApiError) {
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
			StatusCode: http.StatusNotFound,
			Message:    "Order not found",
		}
	}
	return repo.GetById(ctx, id, nil)
}

func (repo *orderRepository) GetRecentActivity(ctx context.Context, userId primitive.ObjectID, userType enums.UserType, branchId *primitive.ObjectID, filter order_filters.OrderFilters) (map[string]float64, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	// Build additional filters from OrderFilters
	var additionalFilters bson.M
	pipeline, _ := filter.BuildPipeline(bson.M{}, ctx, false)
	if len(pipeline) > 0 {
		// Extract the match stage from the first pipeline stage
		if matchStage, exists := pipeline[0]["$match"]; exists {
			if filters, ok := matchStage.(bson.M); ok {
				additionalFilters = filters
			}
		}
	}

	// Build base match stage based on user type
	matchStage := buildUserMatchStage(userId, userType, branchId)

	// Merge additional filters if provided
	for key, value := range additionalFilters {
		matchStage[key] = value
	}

	// Define all pipelines
	pipelines := map[string][]bson.M{
		"processedOrdersTotal":  buildProcessedOrdersPipeline(matchStage),
		"incompleteOrdersTotal": buildIncompleteOrdersPipeline(matchStage),
		"unpaidAmount":          buildUnpaidAmountPipeline(matchStage),
		"overdueAmount":         buildOverdueAmountPipeline(matchStage),
	}

	// Execute all pipelines in parallel and collect results
	stats := make(map[string]float64)
	type pipelineResult struct {
		key    string
		result float64
		err    *helpers.ApiError
	}

	// Create a channel to collect results with timeout context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resultChan := make(chan pipelineResult, len(pipelines))

	var wg sync.WaitGroup
	wg.Add(len(pipelines))

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Execute all pipelines concurrently with panic recovery
	for key, pipeline := range pipelines {
		go func(k string, p []bson.M) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					resultChan <- pipelineResult{
						key:    k,
						result: 0,
						err: &helpers.ApiError{
							StatusCode: 500,
							Message:    fmt.Sprintf("panic in pipeline %s: %v", k, r),
						},
					}
				}
			}()
			result, err := executePipeline(ctxWithTimeout, collection, p)
			resultChan <- pipelineResult{
				key:    k,
				result: result,
				err:    err,
			}
		}(key, pipeline)
	}

	// Collect all results with timeout
	for i := 0; i < len(pipelines); i++ {
		select {
		case res := <-resultChan:
			if res.err != nil {
				return nil, res.err
			}
			stats[res.key] = res.result
		case <-ctxWithTimeout.Done():
			return nil, &helpers.ApiError{
				StatusCode: 408,
				Message:    "Timeout waiting for activity results",
			}
		}
	}

	return stats, nil
}

func buildUserMatchStage(userId primitive.ObjectID, userType enums.UserType, branchId *primitive.ObjectID) bson.M {
	matchStage := bson.M{}
	if userType == enums.Vendor {
		matchStage["vendorInfo.id"] = userId
	} else if userType == enums.Buyer {
		matchStage["buyerInfo.id"] = userId
	}

	if branchId != nil {
		if *branchId != primitive.NilObjectID {
			matchStage["$or"] = []bson.M{
				{"pickupInfo.pickupBranchId": branchId},
				{"deliveryInfo.b2cInfo.branchId": branchId},
			}
		}
	}

	return matchStage
}

// s
func buildProcessedOrdersPipeline(matchStage bson.M) []bson.M {
	revenueEligibleStatuses := []string{"completed", "accepted"}

	return []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"deletedAt": bson.M{"$exists": false}},
					{"deletedAt": nil},
				},
				"$and": []bson.M{matchStage},
			},
		},
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"status": bson.M{"$in": revenueEligibleStatuses}},
					{"status": "pending", "paymentInfo.status": order_enums.PAID},
				},
			},
		},
		{
			"$match": bson.M{
				"$nor": []bson.M{
					{"paymentInfo.paymentMethod": order_enums.ONLINE_PAYMENT, "paymentInfo.status": bson.M{"$ne": order_enums.PAID}},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$totalPrice"},
			},
		},
	}
}

func buildIncompleteOrdersPipeline(matchStage bson.M) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"deletedAt": bson.M{"$exists": false}},
					{"deletedAt": nil},
				},
				"$and": []bson.M{matchStage},
			},
		},
		{
			"$addFields": bson.M{
				"statusLower": bson.M{"$toLower": "$status"},
			},
		},
		{
			"$match": bson.M{
				"statusLower": bson.M{"$ne": settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED},
			},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$totalPrice"},
			},
		},
	}
}

func buildUnpaidAmountPipeline(matchStage bson.M) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"deletedAt": bson.M{"$exists": false}},
					{"deletedAt": nil},
				},
				"$and": []bson.M{matchStage},
			},
		},
		{
			"$addFields": bson.M{
				"statusLower": bson.M{"$toLower": "$status"},
			},
		},
		{
			"$match": bson.M{
				"statusLower": bson.M{"$ne": settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED},
			},
		},
		{
			"$match": bson.M{
				"paymentInfo.payments.paymentAmount": bson.M{"$exists": true},
			},
		},
		{
			"$match": bson.M{
				"paymentInfo.payments.paymentStatus": bson.M{"$ne": order_enums.PAID},
			},
		},
		{
			"$unwind": "$paymentInfo.payments",
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$paymentInfo.payments.paymentAmount"},
			},
		},
	}
}

func buildOverdueAmountPipeline(matchStage bson.M) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					{"deletedAt": bson.M{"$exists": false}},
					{"deletedAt": nil},
				},
				"$and": []bson.M{matchStage},
			},
		},
		{
			"$unwind": "$paymentInfo.payments",
		},
		{
			"$match": bson.M{
				"paymentInfo.payments.paymentStatus": bson.M{"$ne": order_enums.PAID},
				"paymentInfo.payments.dueDate": bson.M{
					"$lt": time.Now(),
				},
			},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$paymentInfo.payments.paymentAmount"},
			},
		},
	}
}

func executePipeline(ctx context.Context, collection *mongo.Collection, pipeline []bson.M) (float64, *helpers.ApiError) {
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to execute pipeline: " + err.Error(),
		}
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to process pipeline results: " + err.Error(),
		}
	}

	if len(result) > 0 {
		if total, ok := result[0]["total"].(float64); ok {
			return total, nil
		}
	}

	return 0, nil
}

// TODO: function need to be updated, this is not the correct way to get the buyer activity, we need to get more details about how to calculate the buyer activity
func (repo *orderRepository) GetBuyerActivity(ctx context.Context, buyerId primitive.ObjectID, branchId *primitive.ObjectID, filter order_filters.OrderFilters) (*order_dto.BuyerActivityResponse, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	//TODO: based on sales settings if it was B2B we should add more filters to get only for specific vendor
	//If it was B2C we should get all orders

	// Match stage to filter orders by buyer ID
	matchStage := bson.M{}

	if branchId != nil {
		if *branchId != primitive.NilObjectID {
			matchStage["$or"] = []bson.M{
				{"pickupInfo.pickupBranchId": branchId},
				{"deliveryInfo.b2cInfo.branchId": branchId},
			}
		}
	}

	filterPipeline, _ := filter.BuildPipeline(matchStage, ctx, false)

	// Get total orders count
	totalCountPipeline := append(filterPipeline, bson.M{
		"$count": "total",
	})

	// Get delivered orders count
	deliveredCountPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"status": settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED,
		},
	}, bson.M{
		"$count": "total",
	})

	// Get cancelled orders count
	cancelledCountPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"status": settings_sales_enums.SYSTEM_STATUS_TYPE_CANCELED,
		},
	}, bson.M{
		"$count": "total",
	})

	// Get pending orders count
	pendingCountPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"status": settings_sales_enums.SYSTEM_STATUS_TYPE_PENDING,
		},
	}, bson.M{
		"$count": "total",
	})

	// Get rejected orders count
	rejectedCountPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"status": settings_sales_enums.SYSTEM_STATUS_TYPE_FAILED,
		},
	}, bson.M{
		"$count": "total",
	})

	// Get accepted orders count
	acceptedCountPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"status": settings_sales_enums.SYSTEM_STATUS_TYPE_ACCEPTED,
		},
	}, bson.M{
		"$count": "total",
	})

	// Get total spent amount (match order report logic: ignore canceled, failed, rejected, refunded)
	revenueEligibleStatuses := []string{"completed", "accepted"}

	totalSpentPipeline := append(filterPipeline, bson.M{
		"$match": bson.M{
			"$or": []bson.M{
				{"status": bson.M{"$in": revenueEligibleStatuses}},
				{"status": "pending", "paymentInfo.status": order_enums.PAID},
			},
		},
	}, bson.M{
		"$match": bson.M{
			"$nor": []bson.M{
				{"paymentInfo.paymentMethod": order_enums.ONLINE_PAYMENT, "paymentInfo.status": bson.M{"$ne": order_enums.PAID}},
			},
		},
	}, bson.M{
		"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$totalPrice"},
		},
	})

	// Execute all pipelines in parallel
	type activityResult struct {
		name   string
		result interface{}
		err    error
	}

	// Create channel for results with timeout context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resultChan := make(chan activityResult, 7)

	var wg sync.WaitGroup
	wg.Add(7)

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Execute all pipelines concurrently with panic recovery
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"total", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, totalCountPipeline)
		if err != nil {
			resultChan <- activityResult{"total", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"total", nil, err}
			return
		}
		resultChan <- activityResult{"total", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"delivered", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, deliveredCountPipeline)
		if err != nil {
			resultChan <- activityResult{"delivered", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"delivered", nil, err}
			return
		}
		resultChan <- activityResult{"delivered", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"cancelled", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, cancelledCountPipeline)
		if err != nil {
			resultChan <- activityResult{"cancelled", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"cancelled", nil, err}
			return
		}
		resultChan <- activityResult{"cancelled", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"pending", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, pendingCountPipeline)
		if err != nil {
			resultChan <- activityResult{"pending", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"pending", nil, err}
			return
		}
		resultChan <- activityResult{"pending", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"rejected", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, rejectedCountPipeline)
		if err != nil {
			resultChan <- activityResult{"rejected", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"rejected", nil, err}
			return
		}
		resultChan <- activityResult{"rejected", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"accepted", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, acceptedCountPipeline)
		if err != nil {
			resultChan <- activityResult{"accepted", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"accepted", nil, err}
			return
		}
		resultChan <- activityResult{"accepted", result, nil}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				resultChan <- activityResult{"spent", nil, fmt.Errorf("panic: %v", r)}
			}
		}()
		cursor, err := collection.Aggregate(ctxWithTimeout, totalSpentPipeline)
		if err != nil {
			resultChan <- activityResult{"spent", nil, err}
			return
		}
		defer cursor.Close(ctxWithTimeout)

		var result []bson.M
		if err := cursor.All(ctxWithTimeout, &result); err != nil {
			resultChan <- activityResult{"spent", nil, err}
			return
		}
		resultChan <- activityResult{"spent", result, nil}
	}()

	// Initialize response
	activity := &order_dto.BuyerActivityResponse{
		TotalOrders:      0,
		DeliveredOrders:  0,
		CancelledOrders:  0,
		PendingOrders:    0,
		RejectedOrders:   0,
		AcceptedOrders:   0,
		TotalSpentAmount: 0,
	}

	// Collect all results with timeout
	for i := 0; i < 7; i++ {
		select {
		case res := <-resultChan:
			if res.err != nil {
				return nil, &helpers.ApiError{
					StatusCode: 500,
					Message:    "Failed to get " + res.name + " data: " + res.err.Error(),
				}
			}

			result, ok := res.result.([]bson.M)
			if !ok {
				continue
			}
			if len(result) > 0 {
				switch res.name {
				case "total":
					if total, ok := result[0]["total"].(int32); ok {
						activity.TotalOrders = int(total)
					}
				case "delivered":
					if total, ok := result[0]["total"].(int32); ok {
						activity.DeliveredOrders = int(total)
					}
				case "cancelled":
					if total, ok := result[0]["total"].(int32); ok {
						activity.CancelledOrders = int(total)
					}
				case "pending":
					if total, ok := result[0]["total"].(int32); ok {
						activity.PendingOrders = int(total)
					}
				case "rejected":
					if total, ok := result[0]["total"].(int32); ok {
						activity.RejectedOrders = int(total)
					}
				case "accepted":
					if total, ok := result[0]["total"].(int32); ok {
						activity.AcceptedOrders = int(total)
					}
				case "spent":
					if total, ok := result[0]["total"].(float64); ok {
						activity.TotalSpentAmount = total
					}
				}
			}
		case <-ctxWithTimeout.Done():
			return nil, &helpers.ApiError{
				StatusCode: 408,
				Message:    "Timeout waiting for activity results",
			}
		}
	}

	return activity, nil
}
