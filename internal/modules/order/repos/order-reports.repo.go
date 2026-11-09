package order_repos

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
)

type OrderReportsRepo interface {
	GetOrdersByFilter(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError)
	GetOrderStatsByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) (*order_models.OrderStats, *helpers.ApiError)
	GetOrdersByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError)
	GetOrdersByItemFilter(
		ctx context.Context,
		filter *order_dto.OrderItemFilterDTO,
	) ([]order_models.OrderResponse, *helpers.ApiError)
	GetOrdersByFilterSeals(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError)
	GetOrdersWithStats(ctx context.Context, filter *order_dto.OrderReportDTO) (*order_models.PaginatedCustomersWithStats, *helpers.ApiError)
	GetAggregatedBranchSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedBranchSales, *helpers.ApiError)
	GetAggregatedBranchSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedBranchSales, *helpers.ApiError)
	GetProductSalesAggregationV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError)
	GetProductSalesAggregationV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError)
	GetProductSalesAggregatedReportV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError)
	GetAggregatedProductSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedProductSales, *helpers.ApiError)
	GetAggregatedProductSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedProductSales, *helpers.ApiError)
	// GetReportDependProductV3(ctx context.Context, filter *order_dto.OrderReportDTO) (map[string][]order_models.AggregatedProductSales, *helpers.ApiError)
}
type orderReportsRepository struct {
	*base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
}

func NewOrderReportsRepository(i *do.Injector) (OrderReportsRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &orderReportsRepository{
		BaseRepo: base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
	}, nil
}

func (repo *orderReportsRepository) GetAggregatedProductSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedProductSales, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	var dateFormat string
	switch filter.Period {
	case "This Week":
		dateFormat = "%w"
	case "This Month":
		dateFormat = "%d"
	case "This Year":
		dateFormat = "%b"
	default:
		dateFormat = "%d"
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{Key: "$match", Value: func() bson.M {
			if filter.CategoryId != "" {
				catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
				return bson.M{"orderItems.product.categoryIds": catId}
			}
			return bson.M{}
		}()}},
		bson.D{{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"groupKey": bson.M{"$dateToString": bson.M{"format": dateFormat, "date": "$createdAt"}},
					"itemId":   "$orderItems.itemId",
					"name":     "$orderItems.name.en",
				},
				"salesValue": bson.M{"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$or": []bson.M{
							{"$in": bson.A{"$status", revenueEligibleStatuses}},
							{"$and": []bson.M{
								{"$eq": bson.A{"$status", "pending"}},
								{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
							}},
						}},
						"$orderItems.totalPrice",
						0,
					},
				}},
				"quantity": bson.M{"$sum": "$orderItems.quantity"},
			},
		}},
		bson.D{{
			Key: "$project", Value: bson.M{
				"groupKey":    "$_id.groupKey",
				"itemId":      "$_id.itemId",
				"productName": "$_id.name",
				"salesValue":  "$salesValue",
				"quantity":    "$quantity",
			},
		}},
		bson.D{{Key: "$sort", Value: bson.M{"groupKey": 1}}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err2 := collection.Aggregate(ctxWithTimeout, pipeline)
	if err2 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to aggregate product sales"}
	}

	var result []order_models.AggregatedProductSales
	if err2 = cursor.All(ctxWithTimeout, &result); err2 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode product sales results"}
	}

	return result, nil
}

func (repo *orderReportsRepository) GetAggregatedProductSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedProductSales, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	// Determine the date format based on period type
	var dateFormat string
	switch strings.ToLower(filter.PeriodType) {
	case "daily":
		dateFormat = "%Y-%m-%d"
	case "weekly":
		dateFormat = "%Y-%U" // Year-Week number
	case "monthly":
		dateFormat = "%Y-%m" // Year-Month
	case "yearly":
		dateFormat = "%Y" // Just year
	default:
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid PeriodType"}
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{Key: "$match", Value: func() bson.M {
			if filter.CategoryId != "" {
				catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
				return bson.M{"orderItems.product.categoryIds": catId}
			}
			return bson.M{}
		}()}},
		// Convert price values to decimal to ensure consistent typing
		bson.D{{Key: "$addFields", Value: bson.M{
			"orderItems.totalPriceDecimal": bson.M{"$toDecimal": "$orderItems.totalPrice"},
			"orderItems.priceDecimal":      bson.M{"$toDecimal": "$orderItems.price.value"},
		}}},
		bson.D{{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"groupKey": bson.M{"$dateToString": bson.M{"format": dateFormat, "date": "$createdAt"}},
					"itemId":   "$orderItems.itemId",
					"name":     "$orderItems.name.en",
				},
				"salesValue": bson.M{"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$or": []bson.M{
							{"$in": bson.A{"$status", revenueEligibleStatuses}},
							{"$and": []bson.M{
								{"$eq": bson.A{"$status", "pending"}},
								{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
							}},
						}},
						"$orderItems.totalPriceDecimal",
						0,
					},
				}},
				"quantity": bson.M{"$sum": "$orderItems.quantity"},
				"price":    bson.M{"$avg": "$orderItems.priceDecimal"},
			},
		}},
		bson.D{{
			Key: "$project", Value: bson.M{
				"groupKey":    "$_id.groupKey",
				"itemId":      "$_id.itemId",
				"productName": "$_id.name",
				"salesValue":  1,
				"quantity":    1,
				"price":       1,
			},
		}},
		bson.D{{Key: "$sort", Value: bson.M{"groupKey": 1}}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to aggregate product sales: " + err.Error()}
	}

	var result []order_models.AggregatedProductSales
	if err = cursor.All(ctxWithTimeout, &result); err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode product sales results: " + err.Error()}
	}

	return result, nil
}

func (repo *orderReportsRepository) buildOrderMatchFilter(filter *order_dto.OrderReportDTO) (bson.M, error) {
	var startDate, endDate time.Time
	var err error

	if filter.StartDate != "" {
		startDate, err = time.Parse(time.RFC3339, filter.StartDate)
		if err != nil {
			return nil, err
		}
	}
	if filter.EndDate != "" {
		endDate, err = time.Parse(time.RFC3339, filter.EndDate)
		if err != nil {
			return nil, err
		}
	}

	match := bson.M{
		"$or": []bson.M{
			{"deletedAt": bson.M{"$exists": false}},
			{"deletedAt": nil},
		},
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": startDate, "$lte": endDate}
	} else if !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": startDate}
	} else if !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": endDate}
	}

	if filter.PaymentMethod != "" {
		match["paymentInfo.paymentMethod"] = filter.PaymentMethod
	}
	if len(filter.Status) > 0 {
		match["status"] = bson.M{"$in": filter.Status}
	}
	if filter.DeliveryType != "" {
		match["deliveryInfo.deliveryType"] = filter.DeliveryType
	}

	if !filter.IncludeOnlineUnPaidOrders {
		if _, exists := match["$and"]; !exists {
			match["$and"] = []bson.M{}
		}
		match["$and"] = append(match["$and"].([]bson.M), bson.M{
			"$nor": []bson.M{
				{"paymentInfo.paymentMethod": order_enums.ONLINE_PAYMENT, "paymentInfo.status": bson.M{"$ne": order_enums.PAID}},
			},
		})
	}

	// Handle BranchId filter
	if filter.BranchId != "" {
		id, err := primitive.ObjectIDFromHex(filter.BranchId)
		if err == nil {
			branchCondition := bson.M{
				"$or": []bson.M{
					{"pickupInfo.pickupBranchId": filter.BranchId},
					{"deliveryInfo.b2cInfo.branchId": filter.BranchId},
					{"deliveryInfo.b2cInfo.branchId": id},
					{"pickupInfo.pickupBranchId": id},
					{"orderItems.branchId": id},
				},
			}
			if _, exists := match["$and"]; !exists {
				match["$and"] = []bson.M{}
			}
			match["$and"] = append(match["$and"].([]bson.M), branchCondition)
		}
	}

	// Handle ProductId filter
	if filter.ProductId != "" {
		productObjId, err := primitive.ObjectIDFromHex(filter.ProductId)
		if err == nil {
			if _, exists := match["$and"]; !exists {
				match["$and"] = []bson.M{}
			}
			match["$and"] = append(match["$and"].([]bson.M), bson.M{
				"orderItems": bson.M{
					"$elemMatch": bson.M{"itemId": productObjId},
				},
			})
		}
	}

	// Handle CategoryId filter
	if filter.CategoryId != "" {
		categoryObjId, err := primitive.ObjectIDFromHex(filter.CategoryId)
		if err == nil {
			if _, exists := match["$and"]; !exists {
				match["$and"] = []bson.M{}
			}
			match["$and"] = append(match["$and"].([]bson.M), bson.M{
				"orderItems": bson.M{
					"$elemMatch": bson.M{"product.categoryIds": categoryObjId},
				},
			})
		}
	}

	return match, nil
}

func (repo *orderReportsRepository) GetOrdersWithStats(ctx context.Context, filter *order_dto.OrderReportDTO) (*order_models.PaginatedCustomersWithStats, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 5000 {
		filter.Limit = 5000
	}

	skip := (filter.Page - 1) * filter.Limit

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := []bson.M{
		{"$match": match},
		{"$facet": bson.M{
			"customers": []bson.M{
				{"$group": bson.M{
					"_id":           "$buyerInfo.id",
					"customerName":  bson.M{"$first": "$buyerInfo.fullName"},
					"customerEmail": bson.M{"$first": "$buyerInfo.email"},
					"customerPhone": bson.M{"$first": "$buyerInfo.phone"},
					"locationCity":  bson.M{"$first": "$buyerInfo.location.city"},
					"totalOrder":    bson.M{"$sum": 1},
				}},
				{"$sort": bson.M{"totalOrder": -1}},
				{"$skip": skip},
				{"$limit": filter.Limit},
			},
			"stats": []bson.M{
				{"$group": bson.M{
					"_id":            nil,
					"totalPurchases": bson.M{"$sum": 1},
					"totalQuantity":  bson.M{"$sum": bson.M{"$sum": "$orderItems.quantity"}},
					"totalAmount": bson.M{"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$or": []bson.M{
								{"$in": bson.A{"$status", revenueEligibleStatuses}},
								{"$and": []bson.M{
									{"$eq": bson.A{"$status", "pending"}},
									{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
								}},
							}},
							"$totalPrice",
							0,
						},
					}},
					"uniqueCustomers": bson.M{"$addToSet": "$buyerInfo.id"},
				}},
				{"$project": bson.M{
					"totalPurchases":      1,
					"totalQuantity":       1,
					"totalAmount":         1,
					"customersPurchasing": bson.M{"$size": "$uniqueCustomers"},
				}},
			},
		}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Aggregate failed: " + err1.Error()}
	}
	defer cursor.Close(ctxWithTimeout)

	var results []order_models.PaginatedCustomersWithStats
	if err := cursor.All(ctxWithTimeout, &results); err != nil || len(results) == 0 {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode aggregation result: " + err1.Error()}
	}

	if len(results[0].Stats) == 0 {
		results[0].Stats = []order_models.OrderStat{{
			CustomersPurchasing: 0,
			TotalPurchases:      0,
			TotalQuantity:       0,
			TotalAmount:         0,
		}}
	}

	return &results[0], nil
}
func (repo *orderReportsRepository) GetOrderStatsByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) (*order_models.OrderStats, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{
			Key: "$facet", Value: bson.M{
				"overallStats": mongo.Pipeline{
					bson.D{{
						Key: "$group", Value: bson.M{
							"_id": nil,
							"totalSales": bson.M{"$sum": bson.M{
								"$cond": bson.A{
									bson.M{"$or": []bson.M{
										{"$in": bson.A{"$status", revenueEligibleStatuses}},
										{"$and": []bson.M{
											{"$eq": bson.A{"$status", "pending"}},
											{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
										}},
									}},
									"$totalPrice",
									0,
								},
							}},
							"orderCount": bson.M{"$sum": 1},
						},
					}},
				},
				"topSellingItem": mongo.Pipeline{
					bson.D{{Key: "$unwind", Value: "$orderItems"}},
					bson.D{{
						Key: "$group", Value: bson.M{
							"_id":      "$orderItems.itemId",
							"name":     bson.M{"$first": "$orderItems.name.en"},
							"quantity": bson.M{"$sum": "$orderItems.quantity"},
						},
					}},
					bson.D{{Key: "$sort", Value: bson.M{"quantity": -1}}},
					bson.D{{Key: "$limit", Value: 1}},
				},
			},
		}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Aggregation failed: " + err.Error()}
	}
	defer cursor.Close(ctxWithTimeout)

	var result struct {
		OverallStats []struct {
			TotalSales float64 `bson:"totalSales"`
			OrderCount int     `bson:"orderCount"`
		} `bson:"overallStats"`
		TopSellingItem []struct {
			ID       primitive.ObjectID `bson:"_id"`
			Name     string             `bson:"name"`
			Quantity int                `bson:"quantity"`
		} `bson:"topSellingItem"`
	}

	if cursor.Next(ctxWithTimeout) {
		if err := cursor.Decode(&result); err != nil {
			return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode stats: " + err.Error()}
		}
	}

	totalSales := 0.0
	orderCount := 0
	avg := 0.0

	if len(result.OverallStats) > 0 {
		totalSales = result.OverallStats[0].TotalSales
		orderCount = result.OverallStats[0].OrderCount
		if orderCount > 0 {
			avg = totalSales / float64(orderCount)
		}
	}

	stats := &order_models.OrderStats{
		TotalSales:     totalSales,
		NumberOfOrders: orderCount,
		AvgOrderValue:  avg,
	}

	if len(result.TopSellingItem) > 0 {
		stats.TopSellingItemId = result.TopSellingItem[0].ID.Hex()
		stats.TopSellingName = result.TopSellingItem[0].Name
		stats.TopQuantity = result.TopSellingItem[0].Quantity
	}

	return stats, nil
}

func (repo *orderReportsRepository) GetOrdersByFilterV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	pipeline := []bson.M{{"$match": match}, {"$sort": bson.M{"createdAt": -1}}}

	limit := 50
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	skip := 0
	if filter.Page > 0 {
		skip = (filter.Page - 1) * limit
	}
	if limit > 5000 {
		limit = 5000
	}
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	orders, apiErr := repo.BaseRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}
	return orders, nil
}

func (repo *orderReportsRepository) GetAggregatedBranchSalesV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedBranchSales, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	// Determine the date format based on period type
	var dateFormat string
	switch strings.ToLower(filter.PeriodType) {
	case "daily":
		dateFormat = "%Y-%m-%d"
	case "weekly":
		dateFormat = "%Y-%U" // Year-Week number
	case "monthly":
		dateFormat = "%Y-%m" // Year-Month
	case "yearly":
		dateFormat = "%Y" // Just year
	default:
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid PeriodType"}
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"branchId": bson.M{
				"$switch": bson.M{
					"branches": bson.A{
						bson.M{"case": bson.M{"$eq": bson.A{"$type", "pickup"}}, "then": "$pickupInfo.pickupBranchId"},
						bson.M{"case": bson.M{"$eq": bson.A{"$type", "delivery"}}, "then": "$deliveryInfo.b2cInfo.branchId"},
					},
					"default": nil,
				},
			},
		}}},
		bson.D{{Key: "$match", Value: bson.M{"$and": []bson.M{
			{"branchId": bson.M{"$exists": true}},
			{"branchId": bson.M{"$ne": nil}},
			{"branchId": bson.M{"$ne": ""}},
		}}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"groupKey": bson.M{"$dateToString": bson.M{"format": dateFormat, "date": "$createdAt"}},
				"branchId": "$branchId",
			},
			"salesValue": bson.M{"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$or": []bson.M{
						{"$in": bson.A{"$status", revenueEligibleStatuses}},
						{"$and": []bson.M{
							{"$eq": bson.A{"$status", "pending"}},
							{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
						}},
					}},
					"$totalPrice",
					0,
				},
			}},
			"orderSet": bson.M{"$addToSet": "$_id"},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"groupKey":       "$_id.groupKey",
			"branchId":       "$_id.branchId",
			"salesValue":     1,
			"numberOfOrders": bson.M{"$size": "$orderSet"},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"groupKey": 1}}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to aggregate orders"}
	}

	var result []order_models.AggregatedBranchSales
	if err = cursor.All(ctxWithTimeout, &result); err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode results"}
	}

	return result, nil
}

func (repo *orderReportsRepository) GetAggregatedBranchSalesV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.AggregatedBranchSales, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	var dateFormat string
	switch filter.Period {
	case "This Week":
		dateFormat = "%w"
	case "This Month":
		dateFormat = "%d"
	case "This Year":
		dateFormat = "%b"
	default:
		dateFormat = "%d"
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{
			Key: "$group", Value: bson.M{
				"_id": bson.M{
					"groupKey": bson.M{"$dateToString": bson.M{"format": dateFormat, "date": "$createdAt"}},
					"branchId": "$orderItems.branchId",
				},
				"salesValue": bson.M{"$sum": bson.M{
					"$cond": bson.A{
						bson.M{"$or": []bson.M{
							{"$in": bson.A{"$status", revenueEligibleStatuses}},
							{"$and": []bson.M{
								{"$eq": bson.A{"$status", "pending"}},
								{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
							}},
						}},
						"$orderItems.totalPrice",
						0,
					},
				}},
				"orderSet": bson.M{"$addToSet": "$_id"},
			},
		}},
		bson.D{{
			Key: "$project", Value: bson.M{
				"groupKey":       "$_id.groupKey",
				"branchId":       "$_id.branchId",
				"salesValue":     1,
				"numberOfOrders": bson.M{"$size": "$orderSet"},
			},
		}},
		bson.D{{Key: "$sort", Value: bson.M{"groupKey": 1}}},
	}

	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err2 := collection.Aggregate(ctxWithTimeout, pipeline)
	if err2 != nil {
		fmt.Println("err2222222222222222222222222")
		fmt.Println(err2)
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to aggregate orders"}
	}
	var result []order_models.AggregatedBranchSales
	if err2 = cursor.All(ctxWithTimeout, &result); err2 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode results"}
	}

	return result, nil
}

func (repo *orderReportsRepository) GetProductSalesAggregationV2(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError) {
	matchStage, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{Key: "$match", Value: func() bson.M {
			if filter.CategoryId != "" {
				catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
				return bson.M{"orderItems.product.categoryIds": catId}
			}
			return bson.M{}
		}()}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":      "$orderItems.name.en",
			"quantity": bson.M{"$sum": "$orderItems.quantity"},
			"salesValue": bson.M{"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$or": []bson.M{
						{"$in": bson.A{"$status", revenueEligibleStatuses}},
						{"$and": []bson.M{
							{"$eq": bson.A{"$status", "pending"}},
							{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
						}},
					}},
					"$orderItems.totalPrice",
					0,
				},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"name":       "$_id",
			"quantity":   1,
			"salesValue": 1,
			"_id":        0,
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"salesValue": -1}}},
	}

	var results []order_models.ProductSales
	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Aggregation error: " + err.Error()}
	}
	if err = cursor.All(ctxWithTimeout, &results); err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Cursor error: " + err.Error()}
	}
	return results, nil
}

func (repo *orderReportsRepository) GetProductSalesAggregationV3(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.ProductSales, *helpers.ApiError) {
	matchStage, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{Key: "$match", Value: func() bson.M {
			if filter.CategoryId != "" {
				catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
				return bson.M{"orderItems.product.categoryIds": catId}
			}
			return bson.M{}
		}()}},
		// Convert price values to decimal to ensure consistent typing
		bson.D{{Key: "$addFields", Value: bson.M{
			"orderItems.totalPriceDecimal": bson.M{"$toDecimal": "$orderItems.totalPrice"},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":      "$orderItems.name.en",
			"quantity": bson.M{"$sum": "$orderItems.quantity"},
			"salesValue": bson.M{"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$or": []bson.M{
						{"$in": bson.A{"$status", revenueEligibleStatuses}},
						{"$and": []bson.M{
							{"$eq": bson.A{"$status", "pending"}},
							{"$eq": bson.A{"$paymentInfo.status", order_enums.PAID}},
						}},
					}},
					"$orderItems.totalPriceDecimal",
					0,
				},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"name":       "$_id",
			"quantity":   1,
			"salesValue": 1,
			"_id":        0,
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"salesValue": -1}}},
	}

	if filter.Page > 0 && filter.Limit > 0 {
		if filter.Limit > 5000 {
			filter.Limit = 5000
		}
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: (filter.Page - 1) * filter.Limit}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: filter.Limit}})
	}

	collection, err2 := repo.BaseRepo.GetCollection(ctx)
	if err2 != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err2.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Aggregation error: " + err.Error()}
	}

	// Handle decimal conversion
	var rawResults []bson.M
	if err = cursor.All(ctxWithTimeout, &rawResults); err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Failed to decode raw results: " + err.Error()}
	}

	// Convert to ProductSales struct, handling decimal conversion
	var results []order_models.ProductSales
	for _, raw := range rawResults {
		var name string
		if n, ok := raw["name"].(string); ok {
			name = n
		}
		item := order_models.ProductSales{
			Name: name,
		}

		// Handle quantity (should be int)
		if quantity, ok := raw["quantity"].(int64); ok {
			item.Quantity = uint64(quantity)
		}

		// Handle salesValue (decimal conversion)
		if salesValue, ok := raw["salesValue"].(primitive.Decimal128); ok {
			// Convert Decimal128 to float64 using String() and strconv.ParseFloat
			strVal := salesValue.String()
			floatVal, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				// If conversion fails, set to 0
				item.SalesValue = 0
			} else {
				item.SalesValue = floatVal
			}
		} else if salesValue, ok := raw["salesValue"].(float64); ok {
			// If it's already a float64
			item.SalesValue = salesValue
		}

		results = append(results, item)
	}

	return results, nil
}

func (repo *orderReportsRepository) GetOrdersByFilter(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError) {
	matchFilter, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	pipeline := []bson.M{{"$match": matchFilter}}
	orders, apiErr := repo.BaseRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}
	return orders, nil
}

func (repo *orderReportsRepository) GetProductSalesAggregatedReportV2(ctx context.Context, filter *order_dto.OrderReportDTO) (order_models.ProductSalesReportResult, *helpers.ApiError) {
	matchFilter, apiErr := repo.buildOrderMatchFilter(filter)
	if apiErr != nil {
		return order_models.ProductSalesReportResult{}, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + apiErr.Error()}
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 5000 {
		limit = 5000
	}

	revenueEligibleStatuses := []string{"completed", "accepted"}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchFilter}},
		bson.D{{Key: "$unwind", Value: "$orderItems"}},
		bson.D{{Key: "$project", Value: bson.M{
			"productName":     "$orderItems.name.en",
			"productId":       "$orderItems.itemId",
			"branchId":        "$orderItems.branchId",
			"categoryIds":     "$orderItems.product.categoryIds",
			"quantity":        "$orderItems.quantity",
			"orderItemTotal":  "$orderItems.totalPrice",
			"orderTotal":      "$totalPrice",
			"orderId":         "$_id",
			"customerName":    "$buyerInfo.fullName",
			"customerEmail":   "$buyerInfo.email",
			"customerAddress": "$buyerInfo.location.address",
			"paymentMethod":   "$paymentInfo.paymentMethod",
			"status":          "$status",
			"paymentStatus":   "$paymentInfo.status",
		}}},
		bson.D{{Key: "$facet", Value: bson.M{
			"report": bson.A{
				bson.M{"$match": func() bson.M {
					if filter.CategoryId != "" {
						catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
						return bson.M{"categoryIds": catId}
					}
					return bson.M{}
				}()},
				bson.M{"$group": bson.M{
					"_id":             "$productId",
					"productName":     bson.M{"$first": "$productName"},
					"branchId":        bson.M{"$first": "$branchId"},
					"quantity":        bson.M{"$sum": "$quantity"},
					"orderTotal":      bson.M{"$sum": "$orderItemTotal"},
					"customerName":    bson.M{"$first": "$customerName"},
					"customerEmail":   bson.M{"$first": "$customerEmail"},
					"customerAddress": bson.M{"$first": "$customerAddress"},
					"paymentMethod":   bson.M{"$first": "$paymentMethod"},
					"orderId":         bson.M{"$first": "$orderId"},
					"orderPrice":      bson.M{"$first": "$orderTotal"},
				}},
				bson.M{"$sort": bson.M{"quantity": -1}},
			},
			"stat": bson.A{
				bson.M{"$match": func() bson.M {
					if filter.CategoryId != "" {
						catId, _ := primitive.ObjectIDFromHex(filter.CategoryId)
						return bson.M{"categoryIds": catId}
					}
					return bson.M{}
				}()},
				bson.M{"$group": bson.M{
					"_id":  "$productId",
					"name": bson.M{"$first": "$productName"},
					"qty":  bson.M{"$sum": "$quantity"},
					"orderItems": bson.M{"$push": bson.M{
						"orderId":       "$orderId",
						"orderTotal":    "$orderTotal",
						"status":        "$status",
						"paymentStatus": "$paymentStatus",
					}},
				}},
				bson.M{"$sort": bson.M{"qty": -1}},
				bson.M{"$group": bson.M{
					"_id": nil,
					"products": bson.M{"$push": bson.M{
						"name": "$name",
						"qty":  "$qty",
					}},
					"allOrderItems": bson.M{"$push": "$orderItems"},
				}},
				bson.M{"$project": bson.M{
					"totalProductsSold":  bson.M{"$size": "$products"},
					"totalQuantity":      bson.M{"$sum": "$products.qty"},
					"bestSellingProduct": bson.M{"$arrayElemAt": []interface{}{"$products.name", 0}},
					"uniqueOrders": bson.M{
						"$reduce": bson.M{
							"input":        "$allOrderItems",
							"initialValue": []interface{}{},
							"in":           bson.M{"$concatArrays": []interface{}{"$$value", "$$this"}},
						},
					},
				}},
				bson.M{"$unwind": "$uniqueOrders"},
				bson.M{"$group": bson.M{
					"_id":                "$_id",
					"totalProductsSold":  bson.M{"$first": "$totalProductsSold"},
					"totalQuantity":      bson.M{"$first": "$totalQuantity"},
					"bestSellingProduct": bson.M{"$first": "$bestSellingProduct"},
					"orders":             bson.M{"$addToSet": "$uniqueOrders"},
				}},
				bson.M{"$project": bson.M{
					"totalProductsSold":  1,
					"totalQuantity":      1,
					"bestSellingProduct": 1,
					"totalSales": bson.M{"$sum": bson.M{
						"$cond": bson.A{
							bson.M{"$or": []bson.M{
								{"$in": bson.A{"$orders.status", revenueEligibleStatuses}},
								{"$and": []bson.M{
									{"$eq": bson.A{"$orders.status", "pending"}},
									{"$eq": bson.A{"$orders.paymentStatus", order_enums.PAID}},
								}},
							}},
							"$orders.orderTotal",
							0,
						},
					}},
				}},
			},
		}}},
	}
	collection, err1 := repo.BaseRepo.GetCollection(ctx)
	if err1 != nil {
		return order_models.ProductSalesReportResult{}, &helpers.ApiError{StatusCode: 500, Message: "Get Collection failed: " + err1.Error()}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctxWithTimeout, pipeline)
	if err != nil {
		return order_models.ProductSalesReportResult{}, &helpers.ApiError{StatusCode: 500, Message: "Aggregate error: " + err.Error()}
	}

	var results []bson.M
	if err := cursor.All(ctxWithTimeout, &results); err != nil || len(results) == 0 {
		return order_models.ProductSalesReportResult{}, &helpers.ApiError{StatusCode: 500, Message: "Failed to parse result: " + err1.Error()}
	}

	raw := results[0]

	var report []order_models.ProductSalesReport
	reportDocs, _ := raw["report"].(primitive.A)
	for _, r := range reportDocs {
		doc, ok := r.(bson.M)
		if !ok {
			continue
		}

		var productName, customerName, paymentMethod string
		if pn, ok := doc["productName"].(string); ok {
			productName = pn
		}
		if cn, ok := doc["customerName"].(string); ok {
			customerName = cn
		}
		if pm, ok := doc["paymentMethod"].(string); ok {
			paymentMethod = pm
		}

		report = append(report, order_models.ProductSalesReport{
			ProductName:     productName,
			BranchId:        toObjectIDPtr(doc["branchId"]),
			OrderTotal:      toFloat64(doc["orderTotal"]),
			Quantity:        toUInt64(doc["quantity"]),
			CustomerName:    customerName,
			CustomerEmail:   toStringPtr(doc["customerEmail"]),
			CustomerAddress: toStringPtr(doc["customerAddress"]),
			PaymentMethod:   order_enums.PaymentMethod(paymentMethod),
		})
	}

	var stat order_models.ProductSalesStat
	statDocs, _ := raw["stat"].(primitive.A)
	if len(statDocs) > 0 {
		doc, ok := statDocs[0].(bson.M)
		if ok {
			var bestSellingProduct string
			if bsp, ok := doc["bestSellingProduct"].(string); ok {
				bestSellingProduct = bsp
			}
			stat = order_models.ProductSalesStat{
				TotalProductsSold:  toInt(doc["totalProductsSold"]),
				TotalQuantity:      toUInt64(doc["totalQuantity"]),
				TotalSales:         toFloat64(doc["totalSales"]),
				BestSellingProduct: bestSellingProduct,
			}
		}
	}

	return order_models.ProductSalesReportResult{Report: report, Stat: stat}, nil
}

func (repo *orderReportsRepository) GetOrdersByItemFilter(ctx context.Context, filter *order_dto.OrderItemFilterDTO) ([]order_models.OrderResponse, *helpers.ApiError) {
	var startDate, endDate time.Time
	var err error

	if filter.StartDate != "" {
		startDate, err = time.Parse(time.RFC3339, filter.StartDate)
		if err != nil {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid start date format"}
		}
	}
	if filter.EndDate != "" {
		endDate, err = time.Parse(time.RFC3339, filter.EndDate)
		if err != nil {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "Invalid end date format"}
		}
	}

	match := bson.M{
		"deletedAt": bson.M{"$exists": false},
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": startDate, "$lte": endDate}
	} else if !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": startDate}
	} else if !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": endDate}
	}

	if filter.BranchID != "" || filter.CategoryID != "" {
		elemMatch := bson.M{}
		if filter.BranchID != "" {
			id, _ := primitive.ObjectIDFromHex(filter.BranchID)
			elemMatch["branchId"] = id
		}
		if filter.CategoryID != "" {
			id, _ := primitive.ObjectIDFromHex(filter.CategoryID)
			elemMatch["product.categoryIds"] = bson.M{"$in": bson.A{id}}
		}
		match["orderItems"] = bson.M{"$elemMatch": elemMatch}
	}

	fmt.Println("matchhhhhhhhhhhhhhhhhhhhh")
	fmt.Println(match)
	pipeline := []bson.M{{"$match": match}}

	orders, apiErr := repo.BaseRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}
	return orders, nil
}

func (repo *orderReportsRepository) GetOrdersByFilterSeals(ctx context.Context, filter *order_dto.OrderReportDTO) ([]order_models.OrderResponse, *helpers.ApiError) {
	match, err := repo.buildOrderMatchFilter(filter)
	if err != nil {
		return nil, &helpers.ApiError{StatusCode: 500, Message: "Invalid filter: " + err.Error()}
	}

	pipeline := []bson.M{{"$match": match}}
	orders, apiErr := repo.BaseRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}
	return orders, nil
}

// helper converters used in reports
func toFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func toUInt64(i interface{}) uint64 {
	switch v := i.(type) {
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case float64:
		return uint64(v)
	default:
		return 0
	}
}

func toInt(i interface{}) int {
	switch v := i.(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func toObjectIDPtr(v interface{}) *primitive.ObjectID {
	if oid, ok := v.(primitive.ObjectID); ok {
		return &oid
	}
	return nil
}

func toStringPtr(v interface{}) *string {
	if s, ok := v.(string); ok {
		return &s
	}
	return nil
}
