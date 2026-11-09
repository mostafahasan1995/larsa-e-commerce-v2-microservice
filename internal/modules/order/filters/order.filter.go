package order_filters

import (
	"context"
	"encoding/json"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderFilters struct {
	filters.BaseFilters

	QuotationId   *primitive.ObjectID        `bson:"quotationId,omitempty" json:"quotationId,omitempty"`
	Status        *string                    `bson:"status,omitempty" json:"status,omitempty"`
	PaymentStatus *order_enums.PaymentStatus `bson:"paymentStatus,omitempty" json:"paymentStatus,omitempty"`
	Code          *string                    `bson:"code,omitempty" json:"code,omitempty"`
	ProductId     *primitive.ObjectID        `bson:"productId,omitempty" json:"productId,omitempty"`

	VendorId                  *primitive.ObjectID `bson:"vendorId,omitempty" json:"vendorId,omitempty"`
	BuyerId                   *primitive.ObjectID `bson:"buyerId,omitempty" json:"buyerId,omitempty"`
	BranchId                  *primitive.ObjectID `bson:"branchId,omitempty" json:"branchId,omitempty"`
	IncludeOnlineUnPaidOrders bool                `bson:"includeOnlineUnPaidOrders,omitempty" json:"includeOnlineUnPaidOrders,omitempty"`
}

func NewOrderFilters(query string) (*OrderFilters, error) {
	orderFilters := OrderFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &orderFilters); err != nil {
			return nil, err
		}
	}
	return &orderFilters, nil
}

func (orderFilters *OrderFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) ([]bson.M, []bson.M) {
	if orderFilters == nil {
		return []bson.M{{"$match": matchStage}}, []bson.M{}
	}

	pipeline := make([]bson.M, 0)

	baseFilters := &orderFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.ORDER)

	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if orderFilters.QuotationId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"quotationId": *orderFilters.QuotationId}})
	}

	if orderFilters.Status != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"status": *orderFilters.Status}})
	}

	if orderFilters.PaymentStatus != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"paymentInfo.status": *orderFilters.PaymentStatus}})
	}

	if orderFilters.Code != nil && *orderFilters.Code != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"code": bson.M{"$regex": *orderFilters.Code, "$options": "i"}}})
	}

	if orderFilters.ProductId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"products.productId": *orderFilters.ProductId}})
	}

	if orderFilters.VendorId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"vendorInfo.id": *orderFilters.VendorId}})
	}

	if orderFilters.BuyerId != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"buyerInfo.id": *orderFilters.BuyerId}})
	}

	if orderFilters.IncludeOnlineUnPaidOrders {
		// When true, include all orders (no additional filtering needed)
		// The base filters and other conditions will handle the rest
	} else {
		// Exclude online unpaid orders when IncludeOnlineUnPaidOrders is false
		pipeline = append(pipeline,
			bson.M{"$match": bson.M{"$nor": []bson.M{
				{"paymentInfo.paymentMethod": order_enums.ONLINE_PAYMENT, "paymentInfo.status": bson.M{"$ne": order_enums.PAID}},
			}}},
		)
	}

	if orderFilters.BranchId != nil {
		branchFilter := bson.M{
			"$or": []bson.M{
				{"pickupInfo.pickupBranchId": orderFilters.BranchId.Hex()},
				{"pickupInfo.pickupBranchId": *orderFilters.BranchId},
				{"deliveryInfo.b2cInfo.branchId": orderFilters.BranchId.Hex()},
				{"deliveryInfo.b2cInfo.branchId": *orderFilters.BranchId},
			},
		}
		pipeline = append(pipeline, bson.M{"$match": branchFilter})
	}

	postPipeline := make([]bson.M, 0)
	if addExtraOperations {
		postPipeline = base_filters.BuildPopulatePipeline(orderFilters.BaseFilters.PopFields, enums.ORDER, postPipeline)

		pipeline = base_filters.BuildSortingPipeline(orderFilters.BaseFilters.Sorts, pipeline)

		if orderFilters.BaseFilters.PopFields != nil {
			for _, field := range *orderFilters.BaseFilters.PopFields {
				if field == "orderItems" {
					// Use optimized product population to reduce document size
					postPipeline = append(postPipeline, order_utils.OptimizedPopulateOrderItemsProducts...)
					break
				}
			}
		}

		// Add projection to exclude heavy fields and reduce document size
		postPipeline = append(postPipeline, bson.M{
			"$project": bson.M{
				// Exclude only heavy fields that aren't needed - everything else included by default
				"orderItems.product.reviews":        0,
				"orderItems.product.productionLine": 0,
				"orderItems.product.categories":     0,
				"orderItems.product.branches":       0,
				"orderItems.product.metadata":       0,
				"orderItems.product.outOfStock":     0,
				"orderItems.product.technicalSpecs": 0,
				"orderItems.product.attributes":     0,
			},
		})

	}

	return pipeline, postPipeline
}
