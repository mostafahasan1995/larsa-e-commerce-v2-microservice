package stats_repos

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	stats_dto "larsa-e-commerce-microservice-v2/internal/modules/stats/dto"
	"math"
	"strings"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatsOrderRepository struct {
	*base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
}

// addBranchFilter adds branch filtering to a match condition if branchId is provided
func addBranchFilter(match bson.M, branchId *string) {
	if branchId != nil && *branchId != "" {
		if branchObjId, err := primitive.ObjectIDFromHex(*branchId); err == nil {
			match["$or"] = bson.A{
				bson.M{"pickupInfo.pickupBranch": branchObjId},
				bson.M{"deliveryInfo.b2cInfo.branchId": branchObjId},
				bson.M{"orderItems.branchId": branchObjId},
			}
		}
	}
}

// addBranchFilterToD adds branch filtering to a bson.D match condition if branchId is provided
func addBranchFilterToD(matchConditions *bson.D, branchId *string) {
	if branchId != nil && *branchId != "" {
		if branchObjId, err := primitive.ObjectIDFromHex(*branchId); err == nil {
			*matchConditions = append(*matchConditions, bson.E{
				Key: "$or",
				Value: bson.A{
					bson.M{"pickupInfo.pickupBranch": branchObjId},
					bson.M{"deliveryInfo.b2cInfo.branchId": branchObjId},
					bson.M{"orderItems.branchId": branchObjId},
				},
			})
		}
	}
}

type StatsOrderRepo interface {
	// Orders
	GetActiveCustomersCount(ctx context.Context, startDate, endDate time.Time, active bool, branchId *string) (int, error)
	CountOrders(ctx context.Context, startDate, endDate time.Time, branchId *string) (int, error)
	GetOrdersComparisonByDateRange(ctx context.Context, startDate, endDate time.Time, branchId *string) (int64, int64, float64, error)
	GetGrowthRateComparison(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, float64, error)

	// Sales
	GetSalesComparisonByDateRange(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, int, error)
	GetSalesGrowthRateComparison(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, float64, error)
	GetTotalSalesSinceLaunch(ctx context.Context, branchId *string) (float64, error)

	// Companies
	GetTopCompanies(ctx context.Context, startDate, endDate time.Time, branchId *string, limit, skip int64) (types.DataWithPagination[stats_dto.TopCompanyStat], error)
	GetCompanyOrders(ctx context.Context, companyId primitive.ObjectID, startDate, endDate *time.Time, branchId *string, page, limit int) (*stats_dto.CompanyOrdersWithPagination, error)

	// Customers activity
	GetCustomersActivity(ctx context.Context, filter stats_dto.CustomersActivityFilterDto) ([]stats_dto.CustomerActivityRecord, int64, *primitive.ObjectID, error)

	// Coupons usage
	GetTopCoupons(ctx context.Context, filter stats_dto.TopCouponsFilterDto) (types.DataWithPagination[stats_dto.CouponUsageStat], *string, error)
	GetCouponOrders(ctx context.Context, couponCode string, startDate, endDate *time.Time, branchId *string, orderSource, status *string, page, limit int) (*stats_dto.CouponOrdersWithPagination, error)

	// Customer Profile
	GetCustomerProfile(ctx context.Context, customerId primitive.ObjectID) (*stats_dto.CustomerProfileDto, error)
	GetCustomerOrdersSummary(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time) (*stats_dto.OrdersSummaryDto, error)
	GetCustomerOrders(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time, orderStatus, branchId, channel, paymentMethod *string, skip, limit int64) (*types.DataWithPagination[stats_dto.CustomerOrderDto], error)
	GetCustomerCouponUsageSummary(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time) (*stats_dto.CouponUsageSummaryDto, error)
	GetCustomerCouponUsageDetails(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time, skip, limit int64) ([]stats_dto.CustomerCouponUsageDto, int64, error)

	CountPayment(ctx context.Context, status order_enums.PaymentStatus, startDate, endDate time.Time) (int, error)

	// Charts
	GetCouponRedemptionsPerBranch(ctx context.Context, startDate, endDate *time.Time, couponCode, paymentMethod *string) ([]stats_dto.CouponRedemptionPerBranchData, error)
	GetCouponRedemptionsByCustomerType(ctx context.Context, startDate, endDate *time.Time, couponCode *string) ([]stats_dto.CouponRedemptionsByCustomerTypeData, error)
	GetCustomersCouponRedeemsBar(ctx context.Context, startDate, endDate *time.Time, branchId *string, limit, skip int64) ([]stats_dto.CustomersCouponRedeemsBarData, error)
}

func NewStatsOrderRepo(i *do.Injector) (StatsOrderRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &StatsOrderRepository{
		BaseRepo: base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
	}, nil
}

func (r *StatsOrderRepository) GetActiveCustomersCount(ctx context.Context, startDate, endDate time.Time, active bool, branchId *string) (int, error) {
	var pipeline mongo.Pipeline

	if active {
		// Build match conditions
		matchConditions := bson.D{
			{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			{Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}

		// Add branch filtering if provided
		addBranchFilterToD(&matchConditions, branchId)

		pipeline = mongo.Pipeline{
			bson.D{{Key: "$match", Value: matchConditions}},
			bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$buyerInfo.id"}}}},
			bson.D{{Key: "$count", Value: "customersCount"}},
		}
	} else {
		// Inactive within range:
		// 1. Customers who had orders before startDate but none in [startDate, endDate]
		// 2. Customers who have accounts but never placed any orders

		// Get database name from context
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return 0, fmt.Errorf("failed to get app config: %v", err)
		}

		// Get all customers from customers collection
		customersColl := r.BaseRepo.GetDB().Database(cfg.Db).Collection("ecommerceProfile")
		customersCursor, err := customersColl.Aggregate(ctx, mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			}}},
			bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}}}},
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get customers: %v", err)
		}
		defer customersCursor.Close(ctx)

		var allCustomers []struct {
			Id primitive.ObjectID `bson:"_id"`
		}
		if err := customersCursor.All(ctx, &allCustomers); err != nil {
			return 0, fmt.Errorf("failed to decode customers: %v", err)
		}

		// Build match conditions for active customers
		activeMatchConditions := bson.D{
			{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: startDate}, {Key: "$lte", Value: endDate}}},
		}

		// Add branch filtering if provided
		addBranchFilterToD(&activeMatchConditions, branchId)

		// Get customers who ordered during the date range
		ordersColl := r.BaseRepo.GetDB().Database(cfg.Db).Collection("orders")
		activeCustomersCursor, err := ordersColl.Aggregate(ctx, mongo.Pipeline{
			bson.D{{Key: "$match", Value: activeMatchConditions}},
			bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$buyerInfo.id"}}}},
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get active customers: %v", err)
		}
		defer activeCustomersCursor.Close(ctx)

		var activeCustomers []struct {
			Id primitive.ObjectID `bson:"_id"`
		}
		if err := activeCustomersCursor.All(ctx, &activeCustomers); err != nil {
			return 0, fmt.Errorf("failed to decode active customers: %v", err)
		}

		// Create maps for efficient lookup
		activeCustomerMap := make(map[primitive.ObjectID]bool)
		for _, customer := range activeCustomers {
			activeCustomerMap[customer.Id] = true
		}

		// Count inactive customers (customers with accounts but no orders in the range)
		inactiveCount := 0
		for _, customer := range allCustomers {
			if !activeCustomerMap[customer.Id] {
				inactiveCount++
			}
		}

		return inactiveCount, nil
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	cursor, aggErr := coll.Aggregate(ctx, pipeline)
	if aggErr != nil {
		return 0, fmt.Errorf("aggregation failed: %v", aggErr)
	}

	var result []struct {
		CustomerCount int `bson:"customersCount"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, fmt.Errorf("failed to read aggregation results: %v", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].CustomerCount, nil
}

func (r *StatsOrderRepository) CountOrders(ctx context.Context, startDate, endDate time.Time, branchId *string) (int, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	match := bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	// Add branch filtering if provided
	addBranchFilter(match, branchId)

	count, errCount := coll.CountDocuments(ctx, match, options.Count().SetHint("_id_"))
	if errCount != nil {
		return 0, errCount
	}
	return int(count), nil
}

func (r *StatsOrderRepository) GetOrdersComparisonByDateRange(ctx context.Context, startDate, endDate time.Time, branchId *string) (int64, int64, float64, error) {
	if !startDate.Before(endDate) {
		return 0, 0, 0, fmt.Errorf("startDate must be before endDate")
	}

	// Calculate previous period
	delta := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-delta)

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, 0, 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Count for current period
	currentCount, err := coll.CountDocuments(ctx, bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
		"createdAt": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	})
	if err != nil {
		return 0, 0, 0, err
	}

	// Count for previous period
	previousCount, err := coll.CountDocuments(ctx, bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
		"createdAt": bson.M{
			"$gte": prevStart,
			"$lt":  prevEnd,
		},
	})
	if err != nil {
		return 0, 0, 0, err
	}

	// percentage change: ((current - previous) / previous) * 100
	var percentage float64
	if previousCount == 0 {
		percentage = 0
	} else {
		percentage = (float64(currentCount-previousCount) / float64(previousCount)) * 100
	}

	return currentCount, previousCount, percentage, nil
}

func (r *StatsOrderRepository) GetGrowthRateComparison(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, float64, error) {
	if !startDate.Before(endDate) {
		return 0, 0, 0, fmt.Errorf("startDate must be before endDate")
	}
	delta := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-delta)
	_, _, currPercentage, err := r.GetOrdersComparisonByDateRange(ctx, startDate, endDate, branchId)
	fmt.Println("curr: ", currPercentage)
	if err != nil {
		return 0, 0, 0, err
	}
	_, _, prevPercentage, err := r.GetOrdersComparisonByDateRange(ctx, prevStart, prevEnd, branchId)
	fmt.Println("prev: ", prevPercentage)
	if err != nil {
		return 0, 0, 0, err
	}

	var percentageChange float64
	if prevPercentage == 0 {
		if currPercentage == 0 {
			percentageChange = 0
		} else {
			// If there was no previous percentage but now there is, treat as 100% increase
			percentageChange = 100
		}
	} else {
		percentageChange = (currPercentage - prevPercentage) / math.Abs(prevPercentage) * 100
	}

	return currPercentage, prevPercentage, percentageChange, nil
}

func (r *StatsOrderRepository) GetSalesComparisonByDateRange(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, int, error) {
	if !startDate.Before(endDate) {
		return 0, 0, 0, fmt.Errorf("startDate must be before endDate")
	}
	// Calculate previous period
	delta := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-delta)

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, 0, 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	sumSales := func(start, end time.Time) (float64, error) {
		pipeline := mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
				{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
				{Key: "createdAt", Value: bson.D{
					{Key: "$gte", Value: start},
					{Key: "$lt", Value: end},
				}},
			}}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalSales", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
			}}},
		}
		cursor, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)
		var result []struct {
			TotalSales float64 `bson:"totalSales"`
		}
		if err := cursor.All(ctx, &result); err != nil {
			return 0, err
		}
		if len(result) == 0 {
			return 0, nil
		}
		return result[0].TotalSales, nil
	}

	currentSales, err := sumSales(startDate, endDate)
	if err != nil {
		return 0, 0, 0, err
	}
	previousSales, err := sumSales(prevStart, prevEnd)
	if err != nil {
		return 0, 0, 0, err
	}

	var percentage float64
	if previousSales == 0 {
		if currentSales == 0 {
			percentage = 0
		} else {
			percentage = 0
		}
	} else {
		percentage = (currentSales - previousSales) / previousSales * 100
	}
	rounded := int(roundFloat(percentage))
	return currentSales, previousSales, rounded, nil
}

func (r *StatsOrderRepository) GetSalesGrowthRateComparison(ctx context.Context, startDate, endDate time.Time, branchId *string) (float64, float64, float64, error) {
	if !startDate.Before(endDate) {
		return 0, 0, 0, fmt.Errorf("startDate must be before endDate")
	}
	delta := endDate.Sub(startDate)
	prevEnd := startDate
	prevStart := prevEnd.Add(-delta)
	prevPrevEnd := prevStart
	prevPrevStart := prevPrevEnd.Add(-delta)

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, 0, 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	sumSales := func(start, end time.Time) (float64, error) {
		pipeline := mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
				{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
				{Key: "createdAt", Value: bson.D{
					{Key: "$gte", Value: start},
					{Key: "$lt", Value: end},
				}},
			}}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "totalSales", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
			}}},
		}
		cursor, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)
		var result []struct {
			TotalSales float64 `bson:"totalSales"`
		}
		if err := cursor.All(ctx, &result); err != nil {
			return 0, err
		}
		if len(result) == 0 {
			return 0, nil
		}
		return result[0].TotalSales, nil
	}

	currentSales, err := sumSales(startDate, endDate)
	if err != nil {
		return 0, 0, 0, err
	}
	previousSales, err := sumSales(prevStart, prevEnd)
	if err != nil {
		return 0, 0, 0, err
	}
	prevPrevSales, err := sumSales(prevPrevStart, prevPrevEnd)
	if err != nil {
		return 0, 0, 0, err
	}

	var currentGrowth, prevGrowth float64
	if previousSales == 0 {
		currentGrowth = 0
	} else {
		currentGrowth = (currentSales - previousSales) / previousSales * 100
	}
	if prevPrevSales == 0 {
		prevGrowth = 0
	} else {
		prevGrowth = (previousSales - prevPrevSales) / prevPrevSales * 100
	}

	var percentBetween float64
	if prevGrowth == 0 {
		if currentGrowth == 0 {
			percentBetween = 0
		} else {
			percentBetween = 100
		}
	} else {
		percentBetween = (currentGrowth - prevGrowth) / absFloat(prevGrowth) * 100
	}

	return currentGrowth, prevGrowth, percentBetween, nil
}

func (r *StatsOrderRepository) GetTotalSalesSinceLaunch(ctx context.Context, branchId *string) (float64, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Build match conditions
	matchConditions := bson.D{
		{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
		{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
	}

	// Add branch filtering if provided
	addBranchFilterToD(&matchConditions, branchId)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchConditions}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalSales", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		TotalSales float64 `bson:"totalSales"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0].TotalSales, nil
}

// GetTopCompanies returns top buyers (companies) ranked by orders count and revenue in the current month
func (r *StatsOrderRepository) GetTopCompanies(ctx context.Context, startDate, endDate time.Time, branchId *string, limit, skip int64) (types.DataWithPagination[stats_dto.TopCompanyStat], error) {
	fmt.Println("Skip: ", skip)
	fmt.Println("Limit: ", limit)
	var result types.DataWithPagination[stats_dto.TopCompanyStat]

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return result, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// --- Step 1: Get total count of companies ---
	countPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			{Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
		}}},
		bson.D{{Key: "$count", Value: "totalCount"}},
	}

	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return result, err
	}
	var countRes []struct {
		TotalCount int64 `bson:"totalCount"`
	}
	if err := countCursor.All(ctx, &countRes); err != nil {
		return result, err
	}

	var totalCount int64
	if len(countRes) > 0 {
		totalCount = countRes[0].TotalCount
	}

	// --- Step 2: Fetch paginated data ---
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			{Key: "createdAt", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "ordersThisMonth", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "revenueThisMonth", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
			{Key: "companyName", Value: bson.D{{Key: "$first", Value: "$buyerInfo.fullName"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "ordersThisMonth", Value: -1},
			{Key: "revenueThisMonth", Value: -1},
		}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	var monthResults []struct {
		Id               primitive.ObjectID `bson:"_id"`
		OrdersThisMonth  int                `bson:"ordersThisMonth"`
		RevenueThisMonth float64            `bson:"revenueThisMonth"`
		CompanyName      *string            `bson:"companyName"`
	}
	if err := cursor.All(ctx, &monthResults); err != nil {
		return result, err
	}

	if len(monthResults) == 0 {
		result.Data = []stats_dto.TopCompanyStat{}
		result.Pagination = types.Pagination{
			TotalPages: 0,
			PerPage:    int64(limit),
			TotalCount: totalCount,
		}
		return result, nil
	}

	// --- Step 3: Fetch all-time totals for paginated companies ---
	ids := make([]primitive.ObjectID, 0, len(monthResults))
	for _, m := range monthResults {
		ids = append(ids, m.Id)
	}

	allTimePipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
			{Key: "buyerInfo.id", Value: bson.D{{Key: "$in", Value: ids}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "totalOrders", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "totalRevenue", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
		}}},
	}

	allCursor, aerr := coll.Aggregate(ctx, allTimePipeline)
	if aerr != nil {
		return result, aerr
	}
	defer allCursor.Close(ctx)

	type aggAll struct {
		Id           primitive.ObjectID `bson:"_id"`
		TotalOrders  int                `bson:"totalOrders"`
		TotalRevenue float64            `bson:"totalRevenue"`
	}
	var allAgg []aggAll
	if err := allCursor.All(ctx, &allAgg); err != nil {
		return result, err
	}

	allMap := make(map[primitive.ObjectID]aggAll, len(allAgg))
	for _, a := range allAgg {
		allMap[a.Id] = a
	}

	// --- Step 4: Merge data ---
	finalData := make([]stats_dto.TopCompanyStat, 0, len(monthResults))
	for _, m := range monthResults {
		agg := allMap[m.Id]
		finalData = append(finalData, stats_dto.TopCompanyStat{
			CompanyId:        m.Id.Hex(),
			CompanyName:      m.CompanyName,
			OrdersThisMonth:  m.OrdersThisMonth,
			TotalOrders:      agg.TotalOrders,
			RevenueThisMonth: m.RevenueThisMonth,
			TotalRevenue:     agg.TotalRevenue,
		})
	}

	// --- Step 5: Fill result ---
	result = types.DataWithPagination[stats_dto.TopCompanyStat]{
		Data: finalData,
		Pagination: types.Pagination{
			TotalPages: math.Ceil(float64(totalCount) / float64(limit)),
			PerPage:    int64(limit),
			TotalCount: totalCount,
		},
	}

	return result, nil
}

// GetCompanyOrders returns paginated compact orders for a buyer (company)
func (r *StatsOrderRepository) GetCompanyOrders(ctx context.Context, companyId primitive.ObjectID, startDate, endDate *time.Time, branchId *string, page, limit int) (*stats_dto.CompanyOrdersWithPagination, error) {
	// Build match
	match := bson.M{
		"status":       bson.M{"$ne": "refunded"},
		"deletedAt":    bson.M{"$exists": false},
		"buyerInfo.id": companyId,
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	// Add branch filtering if provided
	addBranchFilter(match, branchId)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	skip := int64((page - 1) * limit)

	// Count and fetch data
	countPipeline := []bson.M{
		{"$match": match},
		{"$count": "totalCount"},
	}
	dataPipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": int64(limit)},
		{"$project": bson.M{
			"orderDate":     "$createdAt",
			"status":        "$status",
			"amount":        "$totalPrice",
			"couponApplied": bson.M{"$cond": []interface{}{bson.M{"$ifNull": []interface{}{"$couponCode", nil}}, true, false}},
			"customerName":  "$buyerInfo.fullName",
			"orderCode":     "$code",
		}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Execute count
	countCursor, cerr := coll.Aggregate(ctx, countPipeline)
	if cerr != nil {
		return nil, cerr
	}
	var cresult []struct {
		TotalCount int64 `bson:"totalCount"`
	}
	if err := countCursor.All(ctx, &cresult); err != nil {
		return nil, err
	}
	var totalCount int64
	if len(cresult) > 0 {
		totalCount = cresult[0].TotalCount
	}

	// Execute data
	dataCursor, derr := coll.Aggregate(ctx, dataPipeline)
	if derr != nil {
		return nil, derr
	}
	var data []stats_dto.CompanyOrderDetail
	if err := dataCursor.All(ctx, &data); err != nil {
		return nil, err
	}

	// Build response
	resp := &stats_dto.CompanyOrdersWithPagination{
		Data: data,
		Pagination: types.Pagination{
			TotalPages: func() float64 {
				if limit == 0 {
					return 0
				}
				return float64((totalCount + int64(limit) - 1) / int64(limit))
			}(),
			PerPage:    int64(limit),
			TotalCount: totalCount,
		},
	}
	return resp, nil
}

// GetCustomersActivity aggregates customers with orders >= minOrders in current month
func (r *StatsOrderRepository) GetCustomersActivity(ctx context.Context, filter stats_dto.CustomersActivityFilterDto) ([]stats_dto.CustomerActivityRecord, int64, *primitive.ObjectID, error) {
	// Parse dates - support both ISO 8601 and simple date formats
	startDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.StartDate)
	if err != nil {
		// Fallback to simple date format
		startDate, err = time.Parse("2006-01-02", filter.StartDate)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("invalid startDate: %v", err)
		}
	}
	endDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.EndDate)
	if err != nil {
		// Fallback to simple date format
		endDate, err = time.Parse("2006-01-02", filter.EndDate)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("invalid endDate: %v", err)
		}
	}

	minOrders := filter.MinOrders
	if minOrders <= 0 {
		minOrders = 5
	}

	var branchId *string
	if filter.BranchId != "" {
		branchId = &filter.BranchId
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, 0, nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Build match conditions
	matchConditions := bson.D{
		{Key: "status", Value: bson.D{{Key: "$ne", Value: "refunded"}}},
		{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
		{Key: "createdAt", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	// Add branch filtering if provided
	addBranchFilterToD(&matchConditions, branchId)

	baseStages := mongo.Pipeline{
		// Match stage - uses index
		bson.D{{Key: "$match", Value: matchConditions}},
		// Calculate per-order totals without unwind
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{Key: "orderQty", Value: bson.M{"$sum": "$orderItems.quantity"}},
				{Key: "hasCoupon", Value: bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$couponCode", nil}},
						true,
						false,
					},
				}},
			},
		}},
		// Add branch identification logic
		bson.D{{
			Key: "$addFields",
			Value: bson.D{
				{Key: "currentBranchId", Value: bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{"case": bson.M{"$eq": bson.A{"$type", "pickup"}}, "then": "$pickupInfo.pickupBranch"},
							bson.M{"case": bson.M{"$eq": bson.A{"$type", "delivery"}}, "then": "$deliveryInfo.b2cInfo.branchId"},
						},
						"default": bson.M{
							"$cond": bson.A{
								bson.M{"$and": bson.A{
									bson.M{"$isArray": "$orderItems"},
									bson.M{"$gt": bson.A{bson.M{"$size": "$orderItems"}, 0}},
								}},
								bson.M{"$arrayElemAt": bson.A{"$orderItems.branchId", 0}},
								nil,
							},
						},
					},
				}},
			},
		}},
		// Lookup branch information
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "ecommerceV2Branches",
				"localField":   "currentBranchId",
				"foreignField": "_id",
				"as":           "branchInfo",
			},
		}},
		// Add branch name field
		bson.D{{
			Key: "$addFields",
			Value: bson.D{
				{Key: "branchName", Value: bson.M{
					"$let": bson.M{
						"vars": bson.M{
							"branch": bson.M{"$arrayElemAt": bson.A{"$branchInfo", 0}},
						},
						"in": bson.M{
							"$cond": bson.A{
								bson.M{"$gt": bson.A{bson.M{"$size": "$branchInfo"}, 0}},
								"$$branch.name.en",
								"Unknown Branch",
							},
						},
					},
				}},
			},
		}},
		// Group directly by buyer
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$buyerInfo.id"},
				{Key: "customerName", Value: bson.M{"$first": "$buyerInfo.fullName"}},
				{Key: "companyName", Value: bson.M{"$first": "$buyerInfo.buyer.companyName"}},
				{Key: "email", Value: bson.M{"$first": "$buyerInfo.email"}},
				{Key: "phoneNumber", Value: bson.M{"$first": "$buyerInfo.phone"}},
				{Key: "totalOrdersCount", Value: bson.M{"$sum": 1}},
				{Key: "totalQuantityOrdered", Value: bson.M{"$sum": "$orderQty"}},
				{Key: "totalMonthSalesAmount", Value: bson.M{"$sum": "$totalPrice"}},
				{Key: "lastOrderDate", Value: bson.M{"$max": "$createdAt"}},
				{Key: "usedCouponsCount", Value: bson.M{
					"$sum": bson.M{"$cond": bson.A{
						bson.M{"$eq": bson.A{"$hasCoupon", true}},
						1,
						0,
					}},
				}},
				{Key: "couponDetails", Value: bson.M{
					"$push": bson.M{
						"$cond": bson.A{
							bson.M{"$ifNull": bson.A{"$couponCode", nil}},
							bson.M{
								"code":          "$couponCode",
								"usageDate":     "$createdAt",
								"discountValue": bson.M{"$ifNull": bson.A{"$totalDiscountValue", 0}},
							},
							"$$REMOVE",
						},
					},
				}},
				{Key: "branchNames", Value: bson.M{
					"$addToSet": bson.M{
						"$cond": bson.A{
							bson.M{"$and": bson.A{
								bson.M{"$ne": bson.A{"$branchName", nil}},
								bson.M{"$ne": bson.A{"$branchName", "Unknown Branch"}},
							}},
							"$branchName",
							"$$REMOVE",
						},
					},
				}},
			},
		}},
		// Add lastBranchOrderedFrom field
		bson.D{{
			Key: "$addFields",
			Value: bson.D{
				{Key: "lastBranchOrderedFrom", Value: bson.M{
					"$cond": bson.A{
						bson.M{"$gt": bson.A{bson.M{"$size": "$branchNames"}, 0}},
						bson.M{"$arrayElemAt": bson.A{"$branchNames", -1}},
						"Unknown Branch",
					},
				}},
			},
		}},
		// Apply order count threshold
		bson.D{{
			Key: "$match",
			Value: bson.D{
				{Key: "totalOrdersCount", Value: bson.D{{Key: "$gte", Value: minOrders}}},
			},
		}},
		// Calculate derived fields
		bson.D{{
			Key: "$addFields",
			Value: bson.D{
				{Key: "averageOrderValue", Value: bson.M{
					"$cond": bson.A{
						bson.M{"$gt": bson.A{"$totalOrdersCount", 0}},
						bson.M{"$divide": bson.A{"$totalMonthSalesAmount", "$totalOrdersCount"}},
						0,
					},
				}},
				{Key: "customerType", Value: bson.M{
					"$cond": bson.A{
						bson.M{"$gte": bson.A{"$totalOrdersCount", 2}},
						"returning",
						"new",
					},
				}},
			},
		}},
	}

	// Add additional filters after the base stages
	additionalFilters := bson.D{}

	// Customer name filter (case-insensitive partial match)
	if filter.CustomerName != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "customerName",
			Value: bson.M{"$regex": filter.CustomerName, "$options": "i"},
		})
	}

	// Email filter (case-insensitive partial match)
	if filter.Email != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "email",
			Value: bson.M{"$regex": filter.Email, "$options": "i"},
		})
	}

	// Phone number filter (case-insensitive partial match)
	if filter.PhoneNumber != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "phoneNumber",
			Value: bson.M{"$regex": filter.PhoneNumber, "$options": "i"},
		})
	}

	// Max orders filter
	if filter.MaxOrders != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalOrdersCount",
			Value: bson.M{"$lte": *filter.MaxOrders},
		})
	}

	// Total quantity filters
	if filter.MinTotalQuantity != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalQuantityOrdered",
			Value: bson.M{"$gte": *filter.MinTotalQuantity},
		})
	}
	if filter.MaxTotalQuantity != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalQuantityOrdered",
			Value: bson.M{"$lte": *filter.MaxTotalQuantity},
		})
	}

	// Total sales amount filters
	if filter.MinTotalSalesAmount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalMonthSalesAmount",
			Value: bson.M{"$gte": *filter.MinTotalSalesAmount},
		})
	}
	if filter.MaxTotalSalesAmount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalMonthSalesAmount",
			Value: bson.M{"$lte": *filter.MaxTotalSalesAmount},
		})
	}

	// Average order value filters
	if filter.MinAverageOrderValue != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "averageOrderValue",
			Value: bson.M{"$gte": *filter.MinAverageOrderValue},
		})
	}
	if filter.MaxAverageOrderValue != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "averageOrderValue",
			Value: bson.M{"$lte": *filter.MaxAverageOrderValue},
		})
	}

	// Used coupons count filters
	if filter.MinUsedCouponsCount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "usedCouponsCount",
			Value: bson.M{"$gte": *filter.MinUsedCouponsCount},
		})
	}
	if filter.MaxUsedCouponsCount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "usedCouponsCount",
			Value: bson.M{"$lte": *filter.MaxUsedCouponsCount},
		})
	}

	// Customer type filter
	if filter.CustomerType != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "customerType",
			Value: filter.CustomerType,
		})
	}

	// Last branch ordered from filter (case-insensitive partial match)
	if filter.LastBranchOrderedFrom != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "lastBranchOrderedFrom",
			Value: bson.M{"$regex": filter.LastBranchOrderedFrom, "$options": "i"},
		})
	}

	// Last order date filters - support both ISO 8601 and simple date formats
	if filter.MinLastOrderDate != "" {
		minLastOrderDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.MinLastOrderDate)
		if err != nil {
			// Fallback to simple date format
			minLastOrderDate, err = time.Parse("2006-01-02", filter.MinLastOrderDate)
		}
		if err == nil {
			additionalFilters = append(additionalFilters, bson.E{
				Key:   "lastOrderDate",
				Value: bson.M{"$gte": minLastOrderDate},
			})
		}
	}
	if filter.MaxLastOrderDate != "" {
		maxLastOrderDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.MaxLastOrderDate)
		if err != nil {
			// Fallback to simple date format
			maxLastOrderDate, err = time.Parse("2006-01-02", filter.MaxLastOrderDate)
		}
		if err == nil {
			additionalFilters = append(additionalFilters, bson.E{
				Key:   "lastOrderDate",
				Value: bson.M{"$lte": maxLastOrderDate},
			})
		}
	}

	// Count pipeline
	countPipeline := append(mongo.Pipeline{}, baseStages...)

	// Apply additional filters to count pipeline as well
	if len(additionalFilters) > 0 {
		countPipeline = append(countPipeline, bson.D{{Key: "$match", Value: additionalFilters}})
	}

	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "totalCount"}})

	// Top customer pipeline - find customer with highest total sales amount (without minOrders filter)
	// We need to exclude the last stage of baseStages which is the minOrders filter
	topCustomerPipeline := make(mongo.Pipeline, len(baseStages)-1)
	copy(topCustomerPipeline, baseStages[:len(baseStages)-1])
	topCustomerPipeline = append(topCustomerPipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "totalMonthSalesAmount", Value: -1}}}},
		bson.D{{Key: "$limit", Value: 1}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}}}},
	)

	// Data pipeline with sort and pagination
	dataPipeline := append(mongo.Pipeline{}, baseStages...)

	// Apply additional filters if any
	if len(additionalFilters) > 0 {
		dataPipeline = append(dataPipeline, bson.D{{Key: "$match", Value: additionalFilters}})
	}

	// Sorting whitelist mapping
	sortField := "totalMonthSalesAmount"
	switch filter.SortBy {
	case "totalOrdersCount", "totalMonthSalesAmount", "totalQuantityOrdered", "customerName", "lastOrderDate", "usedCouponsCount", "averageOrderValue":
		sortField = filter.SortBy
	}
	dir := -1
	if strings.ToLower(filter.SortDir) == "asc" {
		dir = 1
	}
	dataPipeline = append(dataPipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: dir}}}})
	dataPipeline = append(dataPipeline,
		bson.D{{Key: "$skip", Value: filter.Skip}},
		bson.D{{Key: "$limit", Value: filter.Limit}},
	)

	// Execute count
	countCursor, cerr := coll.Aggregate(ctx, countPipeline)
	if cerr != nil {
		return nil, 0, nil, cerr
	}

	// Execute top customer query
	var topCustomerId *primitive.ObjectID
	topCustomerCursor, terr := coll.Aggregate(ctx, topCustomerPipeline)
	if terr != nil {
		return nil, 0, nil, terr
	}
	var topCustomerResult []bson.M
	if err := topCustomerCursor.All(ctx, &topCustomerResult); err != nil {
		return nil, 0, nil, err
	}
	if len(topCustomerResult) > 0 {
		if customerId, ok := topCustomerResult[0]["_id"].(primitive.ObjectID); ok {
			topCustomerId = &customerId
		}
	}
	var countRes []struct {
		TotalCount int64 `bson:"totalCount"`
	}
	if err := countCursor.All(ctx, &countRes); err != nil {
		return nil, 0, nil, err
	}
	var totalCount int64
	if len(countRes) > 0 {
		totalCount = countRes[0].TotalCount
	}

	// Compute threshold for high-frequency customers (top 10% by totalOrdersCount)
	top10Count := int(math.Ceil(float64(totalCount) * 0.1))
	if top10Count < 1 {
		top10Count = 0
	}
	highFreqThreshold := 0
	if top10Count > 0 {
		thresholdPipeline := append(mongo.Pipeline{}, baseStages...)
		thresholdPipeline = append(thresholdPipeline,
			bson.D{{Key: "$sort", Value: bson.D{{Key: "totalOrdersCount", Value: -1}, {Key: "totalMonthSalesAmount", Value: -1}}}},
			bson.D{{Key: "$limit", Value: top10Count}},
			bson.D{{Key: "$group", Value: bson.D{
				{Key: "_id", Value: nil},
				{Key: "minTopCount", Value: bson.M{"$min": "$totalOrdersCount"}},
			}}},
		)

		thrCursor, terr := coll.Aggregate(ctx, thresholdPipeline)
		if terr == nil {
			var thrRes []struct {
				MinTopCount int `bson:"minTopCount"`
			}
			if err := thrCursor.All(ctx, &thrRes); err == nil && len(thrRes) > 0 {
				highFreqThreshold = thrRes[0].MinTopCount
			}
		}
	}

	// Execute data
	cursor, err := coll.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, nil, err
	}
	defer cursor.Close(ctx)

	type agg struct {
		CustomerId            primitive.ObjectID `bson:"_id"`
		CustomerName          string             `bson:"customerName"`
		CompanyName           *string            `bson:"companyName"`
		Email                 *string            `bson:"email"`
		PhoneNumber           *string            `bson:"phoneNumber"`
		TotalOrdersCount      int                `bson:"totalOrdersCount"`
		TotalQuantityOrdered  int                `bson:"totalQuantityOrdered"`
		TotalMonthSalesAmount float64            `bson:"totalMonthSalesAmount"`
		LastOrderDate         time.Time          `bson:"lastOrderDate"`
		UsedCouponsCount      int                `bson:"usedCouponsCount"`
		BranchNames           []string           `bson:"branchNames"`
		LastBranchOrderedFrom string             `bson:"lastBranchOrderedFrom"`
		AverageOrderValue     float64            `bson:"averageOrderValue"`
		CustomerType          string             `bson:"customerType"`
		CouponDetails         []struct {
			Code          string    `bson:"code"`
			UsageDate     time.Time `bson:"usageDate"`
			DiscountValue float64   `bson:"discountValue"`
		} `bson:"couponDetails"`
	}

	var rows []agg
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, 0, nil, err
	}

	result := make([]stats_dto.CustomerActivityRecord, 0, len(rows))
	for _, rrow := range rows {
		details := make([]stats_dto.CouponUsageDetail, 0, len(rrow.CouponDetails))
		for _, d := range rrow.CouponDetails {
			details = append(details, stats_dto.CouponUsageDetail{
				Code:          d.Code,
				UsageDate:     d.UsageDate,
				DiscountValue: d.DiscountValue,
			})
		}

		customerType := rrow.CustomerType
		if highFreqThreshold > 0 && rrow.TotalOrdersCount >= highFreqThreshold {
			customerType = "high-frequency"
		}

		result = append(result, stats_dto.CustomerActivityRecord{
			CustomerId:            rrow.CustomerId,
			CustomerName:          rrow.CustomerName,
			CompanyName:           rrow.CompanyName,
			Email:                 rrow.Email,
			PhoneNumber:           rrow.PhoneNumber,
			TotalOrdersCount:      rrow.TotalOrdersCount,
			TotalQuantityOrdered:  rrow.TotalQuantityOrdered,
			TotalMonthSalesAmount: rrow.TotalMonthSalesAmount,
			LastOrderDate:         rrow.LastOrderDate,
			UsedCouponsCount:      rrow.UsedCouponsCount,
			LastBranchOrderedFrom: rrow.LastBranchOrderedFrom,
			AverageOrderValue:     rrow.AverageOrderValue,
			CustomerType:          customerType,
			CouponDetails:         details,
		})
	}
	return result, totalCount, topCustomerId, nil
}

// GetTopCoupons aggregates coupon usage stats within date range
func (r *StatsOrderRepository) GetTopCoupons(ctx context.Context, filter stats_dto.TopCouponsFilterDto) (types.DataWithPagination[stats_dto.CouponUsageStat], *string, error) {
	var result types.DataWithPagination[stats_dto.CouponUsageStat]

	// Parse dates - support both ISO 8601 and simple date formats
	startDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.StartDate)
	if err != nil {
		// Fallback to simple date format
		startDate, err = time.Parse("2006-01-02", filter.StartDate)
		if err != nil {
			return result, nil, fmt.Errorf("invalid startDate: %v", err)
		}
	}
	endDate, err := time.Parse("2006-01-02T15:04:05.000Z", filter.EndDate)
	if err != nil {
		// Fallback to simple date format
		endDate, err = time.Parse("2006-01-02", filter.EndDate)
		if err != nil {
			return result, nil, fmt.Errorf("invalid endDate: %v", err)
		}
	}

	var branchId *string
	if filter.BranchId != "" {
		branchId = &filter.BranchId
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return result, nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Count distinct coupons used in range
	matchStage := bson.M{
		"status":     bson.M{"$ne": "refunded"},
		"deletedAt":  bson.M{"$exists": false},
		"createdAt":  bson.M{"$gte": startDate, "$lte": endDate},
		"couponCode": bson.M{"$ne": nil},
	}

	// Add branch filter if provided
	if branchId != nil && *branchId != "" {
		addBranchFilter(matchStage, branchId)
	}

	countPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$couponCode"}}}},
		bson.D{{Key: "$count", Value: "totalCount"}},
	}

	countCursor, cerr := coll.Aggregate(ctx, countPipeline)
	if cerr != nil {
		return result, nil, cerr
	}
	var countRes []struct {
		TotalCount int64 `bson:"totalCount"`
	}
	if err := countCursor.All(ctx, &countRes); err != nil {
		return result, nil, err
	}
	var totalCount int64
	if len(countRes) > 0 {
		totalCount = countRes[0].TotalCount
	}

	// Top coupon pipeline - find coupon with highest redeem count
	topCouponPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$couponCode"},
			{Key: "redeemCount", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "redeemCount", Value: -1}}}},
		bson.D{{Key: "$limit", Value: 1}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	// Execute top coupon query
	var topCouponCode *string
	topCouponCursor, terr := coll.Aggregate(ctx, topCouponPipeline)
	if terr != nil {
		return result, nil, terr
	}
	var topCouponResult []bson.M
	if err := topCouponCursor.All(ctx, &topCouponResult); err != nil {
		return result, nil, err
	}
	if len(topCouponResult) > 0 {
		if couponCode, ok := topCouponResult[0]["_id"].(string); ok {
			topCouponCode = &couponCode
		}
	}

	// Data pipeline
	dataPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchStage}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$couponCode"},
			{Key: "redeemCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "uniqueCustomers", Value: bson.D{{Key: "$addToSet", Value: "$buyerInfo.id"}}},
			{Key: "totalDiscountValue", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$couponDiscountValue", 0}}}}}},
		}}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "uniqueCustomers", Value: bson.D{{Key: "$size", Value: "$uniqueCustomers"}}},
			{Key: "couponValuePerUse", Value: bson.D{{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$redeemCount", 0}}},
				0,
				bson.D{{Key: "$divide", Value: bson.A{"$totalDiscountValue", "$redeemCount"}}},
			}}}},
		}}},
	}

	// Add additional filters after the base stages
	additionalFilters := bson.D{}

	// Coupon code filter (case-insensitive partial match)
	if filter.CouponCode != "" {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "_id",
			Value: bson.M{"$regex": filter.CouponCode, "$options": "i"},
		})
	}

	// Redeem count filters
	if filter.MinRedeemCount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "redeemCount",
			Value: bson.M{"$gte": *filter.MinRedeemCount},
		})
	}
	if filter.MaxRedeemCount != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "redeemCount",
			Value: bson.M{"$lte": *filter.MaxRedeemCount},
		})
	}

	// Unique customers filters
	if filter.MinUniqueCustomers != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "uniqueCustomers",
			Value: bson.M{"$gte": *filter.MinUniqueCustomers},
		})
	}
	if filter.MaxUniqueCustomers != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "uniqueCustomers",
			Value: bson.M{"$lte": *filter.MaxUniqueCustomers},
		})
	}

	// Total discount value filters
	if filter.MinTotalDiscountValue != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalDiscountValue",
			Value: bson.M{"$gte": *filter.MinTotalDiscountValue},
		})
	}
	if filter.MaxTotalDiscountValue != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "totalDiscountValue",
			Value: bson.M{"$lte": *filter.MaxTotalDiscountValue},
		})
	}

	// Coupon value per use filters
	if filter.MinCouponValuePerUse != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "couponValuePerUse",
			Value: bson.M{"$gte": *filter.MinCouponValuePerUse},
		})
	}
	if filter.MaxCouponValuePerUse != nil {
		additionalFilters = append(additionalFilters, bson.E{
			Key:   "couponValuePerUse",
			Value: bson.M{"$lte": *filter.MaxCouponValuePerUse},
		})
	}

	// Apply additional filters if any
	if len(additionalFilters) > 0 {
		dataPipeline = append(dataPipeline, bson.D{{Key: "$match", Value: additionalFilters}})
	}

	// Sorting whitelist mapping
	sortField := "redeemCount"
	switch filter.SortBy {
	case "redeemCount", "uniqueCustomers", "totalDiscountValue", "couponValuePerUse":
		sortField = filter.SortBy
	}
	dir := -1
	if strings.ToLower(filter.SortDir) == "asc" {
		dir = 1
	}
	dataPipeline = append(dataPipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: sortField, Value: dir}}}})
	dataPipeline = append(dataPipeline,
		bson.D{{Key: "$skip", Value: filter.Skip}},
		bson.D{{Key: "$limit", Value: filter.Limit}},
	)

	cursor, derr := coll.Aggregate(ctx, dataPipeline)
	if derr != nil {
		return result, nil, derr
	}
	defer cursor.Close(ctx)

	type agg struct {
		CouponCode         *string `bson:"_id"`
		RedeemCount        int     `bson:"redeemCount"`
		UniqueCustomers    int     `bson:"uniqueCustomers"`
		TotalDiscountValue float64 `bson:"totalDiscountValue"`
		CouponValuePerUse  float64 `bson:"couponValuePerUse"`
	}
	var rows []agg
	if err := cursor.All(ctx, &rows); err != nil {
		return result, nil, err
	}

	// Get unique coupon codes for batch lookup
	couponCodes := make([]string, 0, len(rows))
	for _, rrow := range rows {
		if rrow.CouponCode != nil {
			couponCodes = append(couponCodes, *rrow.CouponCode)
		}
	}

	// Batch lookup coupon details
	couponDetailsMap := make(map[string]struct {
		ValidityStartDate *time.Time
		ValidityEndDate   *time.Time
		CouponType        string
		DiscountType      string
		IsActive          *bool
	})

	if len(couponCodes) > 0 {
		// Get configuration for database access
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			// If we can't get config, continue without coupon details
		} else {
			// Single efficient query to get all coupon details
			couponColl := r.BaseRepo.GetDB().Database(cfg.Db).Collection("ecommerceV2Coupons")
			couponCursor, err := couponColl.Find(ctx, bson.M{
				"code":      bson.M{"$in": couponCodes},
				"deletedAt": bson.M{"$exists": false},
			}, options.Find().SetProjection(bson.M{
				"code":            1,
				"usage.startDate": 1,
				"usage.endDate":   1,
				"type":            1,
				"discount.type":   1,
				"isActive":        1,
			}))

			if err == nil {
				defer couponCursor.Close(ctx)

				type couponDetail struct {
					Code  string `bson:"code"`
					Usage *struct {
						StartDate *time.Time `bson:"startDate"`
						EndDate   *time.Time `bson:"endDate"`
					} `bson:"usage"`
					Type     *string `bson:"type"`
					Discount *struct {
						Type *string `bson:"type"`
					} `bson:"discount"`
					IsActive *bool `bson:"isActive"`
				}

				var coupons []couponDetail
				if err := couponCursor.All(ctx, &coupons); err == nil {
					for _, coupon := range coupons {
						detail := struct {
							ValidityStartDate *time.Time
							ValidityEndDate   *time.Time
							CouponType        string
							DiscountType      string
							IsActive          *bool
						}{
							CouponType:   "unknown",
							DiscountType: "unknown",
							IsActive:     coupon.IsActive,
						}

						if coupon.Usage != nil {
							detail.ValidityStartDate = coupon.Usage.StartDate
							detail.ValidityEndDate = coupon.Usage.EndDate
						}

						if coupon.Type != nil {
							detail.CouponType = *coupon.Type
						}

						if coupon.Discount != nil && coupon.Discount.Type != nil {
							detail.DiscountType = *coupon.Discount.Type
						}

						couponDetailsMap[coupon.Code] = detail
					}
				}
			}
		}
	}

	// Process results with coupon details
	data := make([]stats_dto.CouponUsageStat, 0, len(rows))
	currentTime := time.Now()

	for _, rrow := range rows {
		code := ""
		if rrow.CouponCode != nil {
			code = *rrow.CouponCode
		}

		// Get coupon details from map
		couponDetail, exists := couponDetailsMap[code]

		// Calculate status and set defaults
		status := "unknown"
		couponType := "unknown"
		var validityStartDate, validityEndDate *time.Time

		if exists {
			couponType = couponDetail.CouponType
			validityStartDate = couponDetail.ValidityStartDate
			validityEndDate = couponDetail.ValidityEndDate

			// Calculate status based on dates and active flag
			if couponDetail.IsActive != nil && !*couponDetail.IsActive {
				status = "inactive"
			} else if validityStartDate != nil && validityEndDate != nil {
				if currentTime.Before(*validityStartDate) {
					status = "scheduled"
				} else if currentTime.After(*validityEndDate) {
					status = "expired"
				} else {
					status = "active"
				}
			} else if couponDetail.IsActive != nil && *couponDetail.IsActive {
				status = "active"
			}
		}

		data = append(data, stats_dto.CouponUsageStat{
			CouponCode:         code,
			RedeemCount:        rrow.RedeemCount,
			UniqueCustomers:    rrow.UniqueCustomers,
			TotalDiscountValue: rrow.TotalDiscountValue,
			CouponValuePerUse:  rrow.CouponValuePerUse,
			ValidityStartDate:  validityStartDate,
			ValidityEndDate:    validityEndDate,
			CouponType:         couponType,
			Status:             status,
		})
	}

	result = types.DataWithPagination[stats_dto.CouponUsageStat]{
		Data: data,
		Pagination: types.Pagination{
			TotalPages: func() float64 {
				if filter.Limit == 0 {
					return 0
				}
				return math.Ceil(float64(totalCount) / float64(filter.Limit))
			}(),
			PerPage:    filter.Limit,
			TotalCount: totalCount,
		},
	}
	return result, topCouponCode, nil
}

// GetCouponOrders returns paginated orders where the coupon was applied
func (r *StatsOrderRepository) GetCouponOrders(ctx context.Context, couponCode string, startDate, endDate *time.Time, branchId *string, orderSource, status *string, page, limit int) (*stats_dto.CouponOrdersWithPagination, error) {
	match := bson.M{
		"status":     bson.M{"$ne": "refunded"},
		"deletedAt":  bson.M{"$exists": false},
		"couponCode": couponCode,
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}
	if status != nil && *status != "" {
		match["status"] = *status
	}
	// Order source field is not defined on Order; skip filtering by source

	// Add branch filtering if provided
	addBranchFilter(match, branchId)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	skip := int64((page - 1) * limit)

	countPipeline := []bson.M{{"$match": match}, {"$count": "totalCount"}}

	dataPipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"createdAt": -1}},
		{"$skip": skip},
		{"$limit": int64(limit)},
		{"$project": bson.M{
			"orderId":        bson.M{"$toString": "$code"},
			"orderDate":      "$createdAt",
			"orderSource":    nil,
			"customerName":   "$buyerInfo.fullName",
			"originalAmount": "$subTotalPrice",
			"discountAmount": bson.M{"$ifNull": bson.A{"$couponDiscountValue", 0}},
			"finalAmount":    "$totalPrice",
			"orderStatus":    "$status",
		}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	countCursor, cerr := coll.Aggregate(ctx, countPipeline)
	if cerr != nil {
		return nil, cerr
	}
	var cresult []struct {
		TotalCount int64 `bson:"totalCount"`
	}
	if err := countCursor.All(ctx, &cresult); err != nil {
		return nil, err
	}
	var totalCount int64
	if len(cresult) > 0 {
		totalCount = cresult[0].TotalCount
	}

	dataCursor, derr := coll.Aggregate(ctx, dataPipeline)
	if derr != nil {
		return nil, derr
	}
	var data []stats_dto.CouponOrderDetail
	if err := dataCursor.All(ctx, &data); err != nil {
		return nil, err
	}

	resp := &stats_dto.CouponOrdersWithPagination{
		Data: data,
		Pagination: types.Pagination{
			TotalPages: func() float64 {
				if limit == 0 {
					return 0
				}
				return float64((totalCount + int64(limit) - 1) / int64(limit))
			}(),
			PerPage:    int64(limit),
			TotalCount: totalCount,
		},
	}
	return resp, nil
}

func roundFloat(f float64) float64 {
	if f >= 0 {
		return float64(int(f + 0.5))
	}
	return float64(int(f - 0.5))
}

func absFloat(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// Customer Profile Methods

func (r *StatsOrderRepository) GetCustomerProfile(ctx context.Context, customerId primitive.ObjectID) (*stats_dto.CustomerProfileDto, error) {
	// Get customer profile from contact module or profile module
	// For now, we'll get it from the order data
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "buyerInfo.id", Value: customerId},
			{Key: "deletedAt", Value: bson.D{{Key: "$exists", Value: false}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "fullName", Value: bson.D{{Key: "$first", Value: "$buyerInfo.fullName"}}},
			{Key: "companyName", Value: bson.D{{Key: "$first", Value: "$buyerInfo.buyer.companyName"}}},
			{Key: "email", Value: bson.D{{Key: "$first", Value: "$buyerInfo.buyer.email"}}},
			{Key: "phoneNumber", Value: bson.D{{Key: "$first", Value: "$buyerInfo.buyer.phone"}}},
		}}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []struct {
		FullName    string  `bson:"fullName"`
		CompanyName *string `bson:"companyName"`
		Email       *string `bson:"email"`
		PhoneNumber *string `bson:"phoneNumber"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	return &stats_dto.CustomerProfileDto{
		Id:          customerId,
		FullName:    result[0].FullName,
		CompanyName: result[0].CompanyName,
		Email:       result[0].Email,
		PhoneNumber: result[0].PhoneNumber,
	}, nil
}

func (r *StatsOrderRepository) GetCustomerOrdersSummary(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time) (*stats_dto.OrdersSummaryDto, error) {
	match := bson.M{
		"buyerInfo.id": customerId,
		"status":       bson.M{"$ne": "refunded"},
		"deletedAt":    bson.M{"$exists": false},
	}

	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "totalQuantity", Value: bson.D{{Key: "$sum", Value: "$orderItems.quantity"}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalOrders", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "totalOrderValue", Value: bson.D{{Key: "$sum", Value: "$totalPrice"}}},
			{Key: "totalQuantityPurchased", Value: bson.D{{Key: "$sum", Value: "$totalQuantity"}}},
			{Key: "ordersPerChannel", Value: bson.D{{Key: "$push", Value: "$type"}}},
			{Key: "ordersValuePerChannel", Value: bson.D{{Key: "$push", Value: bson.M{"channel": "$type", "value": "$totalPrice"}}}},
		}}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []struct {
		TotalOrders            int      `bson:"totalOrders"`
		TotalOrderValue        float64  `bson:"totalOrderValue"`
		TotalQuantityPurchased int      `bson:"totalQuantityPurchased"`
		OrdersPerChannel       []string `bson:"ordersPerChannel"`
		OrdersValuePerChannel  []bson.M `bson:"ordersValuePerChannel"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &stats_dto.OrdersSummaryDto{
			TotalOrders:            0,
			TotalOrderValue:        0,
			OrdersPerChannel:       make(map[string]int),
			OrdersValuePerChannel:  make(map[string]float64),
			TotalQuantityPurchased: 0,
			UncompletedOrdersCount: 0,
		}, nil
	}

	// Count uncompleted orders
	uncompletedMatch := bson.M{
		"buyerInfo.id": customerId,
		"status":       bson.M{"$nin": []string{"completed", "refunded", "cancelled"}},
		"deletedAt":    bson.M{"$exists": false},
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		uncompletedMatch["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	}

	uncompletedPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: uncompletedMatch}},
		bson.D{{Key: "$count", Value: "count"}},
	}

	uncompletedCursor, err := coll.Aggregate(ctx, uncompletedPipeline)
	if err != nil {
		return nil, err
	}

	var uncompletedResult []struct {
		Count int64 `bson:"count"`
	}
	if err := uncompletedCursor.All(ctx, &uncompletedResult); err != nil {
		return nil, err
	}

	uncompletedCount := int64(0)
	if len(uncompletedResult) > 0 {
		uncompletedCount = uncompletedResult[0].Count
	}

	// Process channel data
	ordersPerChannel := make(map[string]int)
	ordersValuePerChannel := make(map[string]float64)

	for _, channel := range result[0].OrdersPerChannel {
		ordersPerChannel[channel]++
	}

	for _, channelValue := range result[0].OrdersValuePerChannel {
		if channel, ok := channelValue["channel"].(string); ok {
			if value, ok := channelValue["value"].(float64); ok {
				ordersValuePerChannel[channel] += value
			}
		}
	}

	return &stats_dto.OrdersSummaryDto{
		TotalOrders:            result[0].TotalOrders,
		TotalOrderValue:        result[0].TotalOrderValue,
		OrdersPerChannel:       ordersPerChannel,
		OrdersValuePerChannel:  ordersValuePerChannel,
		TotalQuantityPurchased: result[0].TotalQuantityPurchased,
		UncompletedOrdersCount: int(uncompletedCount),
	}, nil
}

func (r *StatsOrderRepository) GetCustomerOrders(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time, orderStatus, branchId, channel, paymentMethod *string, skip, limit int64) (*types.DataWithPagination[stats_dto.CustomerOrderDto], error) {
	match := bson.M{
		"buyerInfo.id": customerId,
		"status":       bson.M{"$ne": "refunded"},
		"deletedAt":    bson.M{"$exists": false},
	}

	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	if orderStatus != nil && *orderStatus != "" {
		match["status"] = *orderStatus
	}

	if branchId != nil && *branchId != "" {
		if branchObjId, err := primitive.ObjectIDFromHex(*branchId); err == nil {
			match["branchId"] = branchObjId
		}
	}

	if channel != nil && *channel != "" {
		match["type"] = *channel
	}

	if paymentMethod != nil && *paymentMethod != "" {
		match["paymentInfo.paymentMethod"] = *paymentMethod
	}

	// Build pipeline with filters
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"orderId":        bson.M{"$toString": "$code"},
			"orderDate":      "$createdAt",
			"orderStatus":    "$status",
			"salesChannel":   "$type",
			"branchName":     "$branch.name",
			"invoiceTotal":   "$totalPrice",
			"paymentMethod":  "$paymentInfo.paymentMethod",
			"paymentStatus":  "$paymentInfo.status",
			"usedCouponCode": "$couponCode",
		}}},
	}

	// Get collection
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	// Execute aggregation
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("aggregation failed: %v", err),
		}
	}
	defer cursor.Close(ctx)

	var orders []stats_dto.CustomerOrderDto
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed to decode results: %v", err),
		}
	}

	// Count total documents
	countPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$count", Value: "total"}},
	}

	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("count aggregation failed: %v", err),
		}
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed to decode count results: %v", err),
		}
	}

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Calculate pagination info
	totalPages := float64(0)
	if total > 0 && limit > 0 {
		totalPages = float64((total + limit - 1) / limit) // Ceiling division
	}

	return &types.DataWithPagination[stats_dto.CustomerOrderDto]{
		Data: orders,
		Pagination: types.Pagination{
			TotalCount: total,
			PerPage:    limit,
			TotalPages: totalPages,
		},
	}, nil
}

func (r *StatsOrderRepository) GetCustomerCouponUsageSummary(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time) (*stats_dto.CouponUsageSummaryDto, error) {
	match := bson.M{
		"buyerInfo.id": customerId,
		"status":       bson.M{"$ne": "refunded"},
		"deletedAt":    bson.M{"$exists": false},
		"couponCode":   bson.M{"$exists": true, "$ne": ""},
	}

	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalCoupons", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "totalDiscount", Value: bson.D{{Key: "$sum", Value: "$couponDiscountValue"}}},
			{Key: "lastCouponUsed", Value: bson.D{{Key: "$max", Value: "$couponCode"}}},
		}}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []struct {
		TotalCoupons   int     `bson:"totalCoupons"`
		TotalDiscount  float64 `bson:"totalDiscount"`
		LastCouponUsed *string `bson:"lastCouponUsed"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &stats_dto.CouponUsageSummaryDto{
			TotalCoupons:   0,
			TotalDiscount:  0,
			LastCouponUsed: nil,
		}, nil
	}

	return &stats_dto.CouponUsageSummaryDto{
		TotalCoupons:   result[0].TotalCoupons,
		TotalDiscount:  result[0].TotalDiscount,
		LastCouponUsed: result[0].LastCouponUsed,
	}, nil
}

func (r *StatsOrderRepository) GetCustomerCouponUsageDetails(ctx context.Context, customerId primitive.ObjectID, startDate, endDate *time.Time, skip, limit int64) ([]stats_dto.CustomerCouponUsageDto, int64, error) {
	match := bson.M{
		"buyerInfo.id": customerId,
		"status":       bson.M{"$ne": "refunded"},
		"deletedAt":    bson.M{"$exists": false},
		"couponCode":   bson.M{"$exists": true, "$ne": ""},
	}

	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	// Count total
	countPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$count", Value: "total"}},
	}

	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, 0, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Get data
	dataPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.M{
			"couponCode":    "$couponCode",
			"discountValue": "$couponDiscountValue",
			"invoiceTotal":  "$totalPrice",
			"dateUsed":      "$createdAt",
			"orderId":       bson.M{"$toString": "$code"},
			"channel":       "$type",
			"branch":        "$branch.name",
		}}},
	}

	dataCursor, err := coll.Aggregate(ctx, dataPipeline)
	if err != nil {
		return nil, 0, err
	}

	var coupons []stats_dto.CustomerCouponUsageDto
	if err := dataCursor.All(ctx, &coupons); err != nil {
		return nil, 0, err
	}

	return coupons, total, nil
}

func (r *StatsOrderRepository) CountPayment(ctx context.Context, status order_enums.PaymentStatus, startDate, endDate time.Time) (int, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, fmt.Errorf("something wrong happened while getting collection: %v", apiErr)
	}
	if !startDate.Before(endDate) {
		return 0, fmt.Errorf("startDate must be before endDate")
	}

	filter := bson.M{
		"status":             bson.M{"$ne": "refunded"},
		"deletedAt":          bson.M{"$exists": false},
		"paymentInfo.status": string(status),
		"createdAt": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// Charts Methods
func (r *StatsOrderRepository) GetCouponRedemptionsPerBranch(ctx context.Context, startDate, endDate *time.Time, couponCode, paymentMethod *string) ([]stats_dto.CouponRedemptionPerBranchData, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Build match conditions
	match := bson.M{
		"status":     bson.M{"$ne": "refunded"},
		"deletedAt":  bson.M{"$exists": false},
		"couponCode": bson.M{"$exists": true, "$ne": ""},
	}

	// Add date range filter
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}

	// Add coupon code filter
	if couponCode != nil && *couponCode != "" {
		match["couponCode"] = *couponCode
	}

	// Add payment method filter
	if paymentMethod != nil && *paymentMethod != "" {
		match["paymentInfo.paymentMethod"] = *paymentMethod
	}

	// Aggregation pipeline to get coupon redemptions per branch
	pipeline := mongo.Pipeline{
		// Match orders with coupons
		bson.D{{Key: "$match", Value: match}},

		// Add fields to identify branch from different sources based on order type
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "branchId", Value: bson.M{
				"$switch": bson.M{
					"branches": bson.A{
						bson.M{"case": bson.M{"$eq": bson.A{"$type", "pickup"}}, "then": "$pickupInfo.pickupBranch"},
						bson.M{"case": bson.M{"$eq": bson.A{"$type", "delivery"}}, "then": "$deliveryInfo.b2cInfo.branchId"},
					},
					"default": bson.M{
						"$cond": bson.A{
							bson.M{"$and": bson.A{
								bson.M{"$isArray": "$orderItems"},
								bson.M{"$gt": bson.A{bson.M{"$size": "$orderItems"}, 0}},
							}},
							bson.M{"$arrayElemAt": bson.A{"$orderItems.branchId", 0}},
							nil,
						},
					},
				},
			}},
		}}},

		// Filter out orders without branch information
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "branchId", Value: bson.D{{Key: "$ne", Value: nil}}},
		}}},

		// Lookup branch information
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "ecommerceV2Branches",
			"localField":   "branchId",
			"foreignField": "_id",
			"as":           "branchInfo",
		}}},

		// Group by branch
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$branchId"},
			{Key: "branchName", Value: bson.M{
				"$first": bson.M{
					"$cond": bson.A{
						bson.M{"$gt": bson.A{bson.M{"$size": "$branchInfo"}, 0}},
						bson.M{
							"$let": bson.M{
								"vars": bson.M{
									"branch": bson.M{"$arrayElemAt": bson.A{"$branchInfo", 0}},
								},
								"in": "$$branch.name.en",
							},
						},
						"Unknown Branch",
					},
				},
			}},
			{Key: "totalCouponRedemptions", Value: bson.M{"$sum": 1}},
			{Key: "totalDiscountValue", Value: bson.M{
				"$sum": bson.M{"$ifNull": bson.A{"$couponDiscountValue", 0}},
			}},
		}}},

		// Project final structure
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "branchId", Value: bson.M{"$toString": "$_id"}},
			{Key: "branchName", Value: "$branchName"},
			{Key: "totalCouponRedemptions", Value: "$totalCouponRedemptions"},
			{Key: "totalDiscountValue", Value: "$totalDiscountValue"},
		}}},

		// Sort by redemptions count descending
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "totalCouponRedemptions", Value: -1},
			{Key: "totalDiscountValue", Value: -1},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %v", err)
	}
	defer cursor.Close(ctx)

	var results []stats_dto.CouponRedemptionPerBranchData
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %v", err)
	}

	return results, nil
}

func (r *StatsOrderRepository) GetCouponRedemptionsByCustomerType(ctx context.Context, startDate, endDate *time.Time, couponCode *string) ([]stats_dto.CouponRedemptionsByCustomerTypeData, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// 1) Compute high-frequency threshold (top 10% by totalOrdersCount in range)
	matchAll := bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		matchAll["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		matchAll["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		matchAll["createdAt"] = bson.M{"$lte": *endDate}
	}

	thresholdPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchAll}},
		bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$buyerInfo.id"}, {Key: "totalOrdersCount", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "totalOrdersCount", Value: -1}}}},
	}
	thrCursor, err := coll.Aggregate(ctx, thresholdPipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed (threshold): %v", err)
	}
	var buyers []struct {
		TotalOrdersCount int `bson:"totalOrdersCount"`
	}
	if err := thrCursor.All(ctx, &buyers); err != nil {
		return nil, fmt.Errorf("decode failed (threshold): %v", err)
	}
	highFreqThreshold := 0
	totalBuyers := len(buyers)
	if totalBuyers > 0 {
		limitTop := int(math.Ceil(float64(totalBuyers) * 0.1))
		if limitTop > totalBuyers {
			limitTop = totalBuyers
		}
		minTop := 0
		if limitTop > 0 {
			minTop = buyers[limitTop-1].TotalOrdersCount
		}
		highFreqThreshold = minTop
	}

	// 2) Build match for ALL orders in the date range (not just coupon orders)
	matchAllOrders := bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		matchAllOrders["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		matchAllOrders["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		matchAllOrders["createdAt"] = bson.M{"$lte": *endDate}
	}

	// 3) Aggregate all customers, get their total orders and coupon orders, then classify
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchAllOrders}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "totalOrdersCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ordersWithCoupon", Value: func() interface{} {
				if couponCode != nil && *couponCode != "" {
					// Filter by specific coupon code
					return bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$couponCode", *couponCode}},
								1,
								0,
							},
						},
					}
				} else {
					// Count all coupon orders
					return bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{"$and": bson.A{
									bson.M{"$ne": bson.A{"$couponCode", nil}},
									bson.M{"$ne": bson.A{"$couponCode", ""}},
								}},
								1,
								0,
							},
						},
					}
				}
			}()},
		}}},
		// Only include customers who have used coupons
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "ordersWithCoupon", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "customerType", Value: bson.M{
				"$cond": bson.A{
					bson.M{"$gte": bson.A{"$totalOrdersCount", 2}},
					"returning",
					"new",
				},
			}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$customerType"},
			{Key: "customersSet", Value: bson.M{"$addToSet": "$_id"}},
			{Key: "ordersCount", Value: bson.M{"$sum": "$ordersWithCoupon"}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "customerType", Value: "$_id"},
			{Key: "customersCount", Value: bson.M{"$size": "$customersSet"}},
			{Key: "ordersCount", Value: 1},
			{Key: "_id", Value: 0},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %v", err)
	}
	defer cursor.Close(ctx)

	var rows []struct {
		CustomerType   string `bson:"customerType"`
		CustomersCount int    `bson:"customersCount"`
		OrdersCount    int    `bson:"ordersCount"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, fmt.Errorf("failed to decode results: %v", err)
	}

	// Now we need to get individual customer data to apply high-frequency classification
	// This matches the post-processing logic from GetCustomersActivity
	customersPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: matchAllOrders}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "totalOrdersCount", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "ordersWithCoupon", Value: func() interface{} {
				if couponCode != nil && *couponCode != "" {
					return bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{"$eq": bson.A{"$couponCode", *couponCode}},
								1,
								0,
							},
						},
					}
				} else {
					return bson.M{
						"$sum": bson.M{
							"$cond": bson.A{
								bson.M{"$and": bson.A{
									bson.M{"$ne": bson.A{"$couponCode", nil}},
									bson.M{"$ne": bson.A{"$couponCode", ""}},
								}},
								1,
								0,
							},
						},
					}
				}
			}()},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "ordersWithCoupon", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
	}

	customersCursor, err := coll.Aggregate(ctx, customersPipeline)
	if err != nil {
		return nil, fmt.Errorf("customers aggregation failed: %v", err)
	}
	defer customersCursor.Close(ctx)

	var customers []struct {
		CustomerId       string `bson:"_id"`
		TotalOrdersCount int    `bson:"totalOrdersCount"`
		OrdersWithCoupon int    `bson:"ordersWithCoupon"`
	}
	if err := customersCursor.All(ctx, &customers); err != nil {
		return nil, fmt.Errorf("failed to decode customer results: %v", err)
	}

	// Count customers by type (matching GetCustomersActivity logic)
	customerTypeCounts := map[string]int{
		"new":            0,
		"returning":      0,
		"high-frequency": 0,
	}
	orderTypeCounts := map[string]int{
		"new":            0,
		"returning":      0,
		"high-frequency": 0,
	}

	for _, customer := range customers {
		customerType := "new"
		if customer.TotalOrdersCount >= 2 {
			customerType = "returning"
		}
		// Apply high-frequency classification (matching GetCustomersActivity)
		if highFreqThreshold > 0 && customer.TotalOrdersCount >= highFreqThreshold {
			customerType = "high-frequency"
		}

		customerTypeCounts[customerType]++
		orderTypeCounts[customerType] += customer.OrdersWithCoupon
	}

	// Build response with the recalculated counts
	resp := make([]stats_dto.CouponRedemptionsByCustomerTypeData, 0)
	for customerType, customersCount := range customerTypeCounts {
		if customersCount > 0 { // Only include types that have customers
			resp = append(resp, stats_dto.CouponRedemptionsByCustomerTypeData{
				CustomerType:   customerType,
				CustomersCount: customersCount,
				OrdersCount:    orderTypeCounts[customerType],
			})
		}
	}

	return resp, nil
}

func (r *StatsOrderRepository) GetCustomersCouponRedeemsBar(ctx context.Context, startDate, endDate *time.Time, branchId *string, limit, skip int64) ([]stats_dto.CustomersCouponRedeemsBarData, error) {
	coll, apiErr := r.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return nil, fmt.Errorf("failed to get collection: %v", apiErr.Message)
	}

	// Base match for all orders (not just coupon orders)
	match := bson.M{
		"status":    bson.M{"$ne": "refunded"},
		"deletedAt": bson.M{"$exists": false},
	}
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate, "$lte": *endDate}
	} else if startDate != nil && !startDate.IsZero() {
		match["createdAt"] = bson.M{"$gte": *startDate}
	} else if endDate != nil && !endDate.IsZero() {
		match["createdAt"] = bson.M{"$lte": *endDate}
	}
	// Optional branch filter using existing helper
	addBranchFilter(match, branchId)

	// Create lookup pipeline for checking coupon redemptions
	couponLookupPipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"$expr":      bson.M{"$eq": bson.A{"$buyerInfo.id", "$$customerId"}},
			"status":     bson.M{"$ne": "refunded"},
			"deletedAt":  bson.M{"$exists": false},
			"couponCode": bson.M{"$exists": true, "$ne": ""},
		}}},
	}

	// Apply same date and branch filters to coupon lookup
	if startDate != nil && endDate != nil && !startDate.IsZero() && !endDate.IsZero() {
		couponLookupPipeline = append(couponLookupPipeline, bson.D{{Key: "$match", Value: bson.M{"createdAt": bson.M{"$gte": *startDate, "$lte": *endDate}}}})
	} else if startDate != nil && !startDate.IsZero() {
		couponLookupPipeline = append(couponLookupPipeline, bson.D{{Key: "$match", Value: bson.M{"createdAt": bson.M{"$gte": *startDate}}}})
	} else if endDate != nil && !endDate.IsZero() {
		couponLookupPipeline = append(couponLookupPipeline, bson.D{{Key: "$match", Value: bson.M{"createdAt": bson.M{"$lte": *endDate}}}})
	}

	if branchId != nil && *branchId != "" {
		// Add branch filter to coupon lookup
		branchFilter := bson.M{
			"$or": bson.A{
				bson.M{"type": "pickup", "pickupInfo.pickupBranch": *branchId},
				bson.M{"type": "delivery", "deliveryInfo.b2cInfo.branchId": *branchId},
			},
		}
		couponLookupPipeline = append(couponLookupPipeline, bson.D{{Key: "$match", Value: branchFilter}})
	}

	couponLookupPipeline = append(couponLookupPipeline, bson.D{{Key: "$limit", Value: 1}}) // We just need to know if any exist

	pipeline := mongo.Pipeline{
		// Match all orders (not just coupon orders)
		bson.D{{Key: "$match", Value: match}},

		// Group by customer to get total sales
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$buyerInfo.id"},
			{Key: "customerName", Value: bson.M{"$first": "$buyerInfo.fullName"}},
			{Key: "totalSales", Value: bson.M{"$sum": "$totalPrice"}},
		}}},

		// Lookup to check if customer has any coupon redemptions
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":     "orders",
			"let":      bson.M{"customerId": "$_id"},
			"pipeline": couponLookupPipeline,
			"as":       "couponOrders",
		}}},

		// Add didRedeem field based on whether coupon orders exist
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "didRedeem", Value: bson.M{
				"$gt": bson.A{bson.M{"$size": "$couponOrders"}, 0},
			}},
		}}},

		// Project final structure
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "customerId", Value: bson.M{"$toString": "$_id"}},
			{Key: "customerName", Value: 1},
			{Key: "totalSales", Value: 1},
			{Key: "didRedeem", Value: 1},
			{Key: "_id", Value: 0},
		}}},

		// Sort by total sales (top customers first)
		bson.D{{Key: "$sort", Value: bson.D{{Key: "totalSales", Value: -1}}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %v", err)
	}
	defer cursor.Close(ctx)

	var rows []stats_dto.CustomersCouponRedeemsBarData
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, fmt.Errorf("failed to decode results: %v", err)
	}
	return rows, nil
}
