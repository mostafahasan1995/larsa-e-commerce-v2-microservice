package quotation_report_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_types "larsa-e-commerce-microservice-v2/internal/modules/quotation/types"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuotationSummaryRepo interface {
	GetQuickSummary(ctx context.Context, vendorId primitive.ObjectID) (*quotation_types.QuotationSummary, *helpers.ApiError)
	GetRecentActivity(ctx context.Context, vendorId primitive.ObjectID, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.QuotationActivity, *helpers.ApiError)
	GetMostRequestedProducts(ctx context.Context, vendorId primitive.ObjectID, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.ProductRequestStats, *helpers.ApiError)
	GetExpiringQuotations(ctx context.Context, vendorId primitive.ObjectID) ([]*quotation_models.QuotationResponse, *helpers.ApiError)
}

type QuotationSummaryRepository struct {
	*repo.BaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse]
}

func NewQuotationSummaryRepository(i *do.Injector) (QuotationSummaryRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &QuotationSummaryRepository{
		BaseRepo: repo.NewBaseRepo[quotation_models.Quotation, quotation_models.QuotationResponse](db, enums.QUOTATION),
	}, nil
}

func (repo *QuotationSummaryRepository) GetQuickSummary(ctx context.Context, vendorId primitive.ObjectID) (*quotation_types.QuotationSummary, *helpers.ApiError) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"vendorProfileId": vendorId,
			},
		},
		{
			"$group": bson.M{
				"_id": "$status",
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to aggregate quotations",
		}
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode aggregation results",
		}
	}
	summary := &quotation_types.QuotationSummary{}
	for _, doc := range results {
		status := quotation_enums.QuotationStatus(doc["_id"].(string))
		count := doc["count"].(int32)

		switch status {
		case quotation_enums.SENT_TO_BUYER:
			summary.PendingBuyer = int(count)
		case quotation_enums.APPROVED:
			summary.Accepted = int(count)
		case quotation_enums.REJECTED:
			summary.Rejected = int(count)
		case quotation_enums.CONVERTED_TO_ORDER:
			summary.Archived = int(count)
		}
	}

	return summary, nil
}

func (repo *QuotationSummaryRepository) GetRecentActivity(ctx context.Context, vendorId primitive.ObjectID, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.QuotationActivity, *helpers.ApiError) {
	pipeline := []bson.M{}

	if fromDate != nil && toDate != nil {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"updatedAt": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
				},
			},
		})
	}
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         enums.PROFILE.GetCollectionName(),
			"localField":   "quotationBuyerInfo.buyerProfileId",
			"foreignField": "_id",
			"as":           "quotationBuyerInfo.buyerProfile",
		},
	},
	)

	// Add $unwind to convert the array from lookup to a single object
	pipeline = append(pipeline, bson.M{
		"$unwind": bson.M{
			"path":                       "$quotationBuyerInfo.buyerProfile",
			"preserveNullAndEmptyArrays": true,
		},
	})

	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	// Also add vendor ID to the match
	pipeline = append([]bson.M{
		{
			"$match": bson.M{
				"vendorProfileId": vendorId,
			},
		},
	}, pipeline...)

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to aggregate quotations" + aggregateErr.Error(),
		}
	}
	defer cursor.Close(ctx)

	var quotations []quotation_models.QuotationResponse
	if err := cursor.All(ctx, &quotations); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode aggregation results" + err.Error(),
		}
	}

	activities := make([]*quotation_types.QuotationActivity, len(quotations))
	for i, quotation := range quotations {
		activities[i] = &quotation_types.QuotationActivity{
			Date:        quotation.UpdatedAt,
			Status:      quotation.Status,
			Activity:    quotation.History[len(quotation.History)-1].Action,
			QuotationId: quotation.Code,
			BuyerName:   quotation.QuotationBuyerInfo.BaseQuotationBuyerInfo.ContactPerson,
		}
	}

	return activities, nil
}

func (repo *QuotationSummaryRepository) GetMostRequestedProducts(ctx context.Context, vendorId primitive.ObjectID, fromDate *time.Time, toDate *time.Time) ([]*quotation_types.ProductRequestStats, *helpers.ApiError) {
	pipeline := []bson.M{}

	if fromDate != nil && toDate != nil {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"vendorProfileId": vendorId,
				"updatedAt": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
				},
			},
		})
	}
	pipeline = append(pipeline, bson.M{
		"$unwind": "$products",
	})
	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": "$products.productId",
			"count": bson.M{
				"$sum": 1,
			},
		},
	})
	pipeline = append(pipeline,
		bson.M{
			"$lookup": bson.M{
				"from":         enums.PRODUCT.GetCollectionName(),
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "product",
			},
		})

	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}
	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to aggregate quotations",
		}
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode aggregation results",
		}
	}

	products := make([]*quotation_types.ProductRequestStats, len(results))
	for i, result := range results {
		productArray, ok := result["product"].(primitive.A)
		if !ok || len(productArray) == 0 {
			continue
		}
		product, ok := productArray[0].(bson.M)
		if !ok {
			continue
		}
		products[i] = &quotation_types.ProductRequestStats{
			ProductId:    result["_id"].(primitive.ObjectID).Hex(),
			ProductName:  product["name"],
			RequestCount: int(result["count"].(int32)),
		}
	}

	return products, nil
}

func (repo *QuotationSummaryRepository) GetExpiringQuotations(ctx context.Context, vendorId primitive.ObjectID) ([]*quotation_models.QuotationResponse, *helpers.ApiError) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"vendorProfileId": vendorId,
				"deadline": bson.M{
					"$lte": time.Now(),
				},
			},
		},
		{
			"$sort": bson.M{
				"deadline": 1,
			},
		},
		{
			"$limit": 5,
		},
	}

	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, aggregateErr := collection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to aggregate quotations",
		}
	}
	defer cursor.Close(ctx)

	var quotations []quotation_models.QuotationResponse
	if err := cursor.All(ctx, &quotations); err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to decode aggregation results",
		}
	}

	quotationPtrs := make([]*quotation_models.QuotationResponse, len(quotations))
	for i := range quotations {
		quotationPtrs[i] = &quotations[i]
	}

	return quotationPtrs, nil
}
