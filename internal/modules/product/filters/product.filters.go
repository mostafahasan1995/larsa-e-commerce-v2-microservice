package product_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DateRange represents a date range with greater than or equal and less than or equal
type DateRange struct {
	Gte *time.Time `json:"gte"`
	Lte *time.Time `json:"lte"`
}

// PriceRange represents a price range with minimum and maximum values
type PriceRange struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

type ProductAttributeFilter struct {
	Code     string         `json:"code"`
	Operator enums.Operator `json:"operator"`
	Value    interface{}    `json:"value"`
	Key      *string        `json:"key"`
}

type ProductFilters struct {
	base_filters.BaseFilters
	Name                   *types.Localizable[string] `json:"name,omitempty"`
	VendorId               *primitive.ObjectID        `json:"vendorId,omitempty"`
	Flag                   *string                    `json:"flag,omitempty"`
	BranchId               *primitive.ObjectID        `json:"branchId,omitempty"`
	CategoryIds            *[]primitive.ObjectID      `json:"categoryIds,omitempty"`
	Attributes             *[]ProductAttributeFilter  `json:"attributes,omitempty"`
	Discount               *bool                      `json:"discount,omitempty"`
	Prices                 *[]PriceRange              `json:"prices,omitempty"`
	OutOfStock             *product_enums.OutOfStock  `json:"outOfStock,omitempty"`
	UnitIds                *[]primitive.ObjectID      `json:"unitIds,omitempty"`
	ProductType            *product_enums.ProductType `json:"productType,omitempty"`
	IncludeAllProductTypes *bool                      `json:"includeAllProductTypes,omitempty"`
	IncludeAttributes      *bool                      `json:"includeAttributes,omitempty"`
}

type BuildPipelineArgs struct {
	AddExtraOperations bool
	BranchId           *primitive.ObjectID
	CurrencyCode       *string
	UnitId             *primitive.ObjectID
}

func NewProductFilters(query string) (*ProductFilters, error) {
	ProductFilters := ProductFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &ProductFilters); err != nil {
			return nil, err
		}
	}

	return &ProductFilters, nil
}

// BuildPipeline builds a MongoDB aggregation pipeline based on the provided filters
func (productFilters *ProductFilters) BuildPipeline(matchStage bson.M, ctx context.Context, args BuildPipelineArgs) []bson.M {
	if productFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &productFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.PRODUCT)

	if productFilters.IncludeAllProductTypes == nil || !*productFilters.IncludeAllProductTypes {
		if productFilters.ProductType != nil {
			matchStage["productType"] = productFilters.ProductType
		} else {
			matchStage["productType"] = bson.M{"$ne": product_enums.PRODUCT_TYPE_INVENTORY}
		}
	}

	// Filter by vendor ID if provided in filters
	if productFilters.VendorId != nil {
		matchStage["vendorId"] = productFilters.VendorId
	}

	if productFilters.CategoryIds != nil {
		// Get all child category IDs including the original ones
		allCategoryIds, err := config_utils.GetChildCategoryIds(ctx, *productFilters.CategoryIds)
		if err != nil {
			// If there's an error getting child categories, fall back to original category IDs
			matchStage["categoryIds"] = bson.M{"$in": productFilters.CategoryIds}
		} else {
			// Use all category IDs (original + children)
			matchStage["categoryIds"] = bson.M{"$in": allCategoryIds}
		}
	}

	if productFilters.Flag != nil {
		matchStage["flag"] = productFilters.Flag
	}

	attributeConditions := make([]bson.M, 0)
	if productFilters.Attributes != nil {
		for _, attribute := range *productFilters.Attributes {
			if attribute.Code == string(product_enums.SYSTEM_BRANCH) || attribute.Code == string(product_enums.SYSTEM_ADDON) {
				continue
			}

			var operator = "$" + string(attribute.Operator)
			keyName := "value"
			if attribute.Key != nil {
				keyName = keyName + "." + *attribute.Key
			}
			attributesFilter := bson.M{"attributes": bson.M{
				"$elemMatch": bson.M{
					keyName: bson.M{operator: attribute.Value},
					"code":  attribute.Code,
				},
			}}
			attributeConditions = append(attributeConditions, attributesFilter)
		}
	}

	if len(attributeConditions) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": attributeConditions}})
	}

	if productFilters.BranchId != nil {
		pipeline = append(pipeline, product_utils.ProductRelatedBranchPipeline(*productFilters.BranchId, productFilters.OutOfStock)...)
	}

	// Filter by vendor ID based on user type
	//TODO move to service
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil {
		if cfg != nil && cfg.User != nil {
			if cfg.User.Type == enums.Vendor {
				matchStage["vendorId"] = cfg.User.Id
			}
		}
	}

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	if args.CurrencyCode != nil {
		pipeline = append(pipeline, product_shared.ExtractPriceWithCurrencyPipeline(args.BranchId, *args.CurrencyCode, args.UnitId)...)

		// Normalize discounts
		pipeline = append(pipeline, product_shared.ConvertFixedDiscountToPercentagePipeline(*args.CurrencyCode)...)

		// Apply boolean discount filter
		if productFilters.Discount != nil {
			if *productFilters.Discount {
				// Lookup active promotions and match products having either direct discount or active promotion
				pipeline = append(pipeline, product_utils.ActivePromotionLookupPipeline(args.BranchId)...)
				pipeline = append(pipeline, bson.M{"$match": bson.M{
					"$or": []bson.M{
						{"convertedDiscount.value": bson.M{"$gt": 0}},
						{"hasActivePromotion": true},
					},
				}})
			}
		}

		if productFilters.Prices != nil {
			priceFilterPipeline := make([]bson.M, 0)
			for _, price := range *productFilters.Prices {
				rangeFilter := bson.M{}

				if price.Min != nil {
					rangeFilter["$gte"] = price.Min
				}
				if price.Max != nil {
					rangeFilter["$lte"] = price.Max
				}
				priceFilterPipeline = append(priceFilterPipeline, bson.M{
					"$match": bson.M{
						"priceItem.value": rangeFilter,
					},
				})
			}
			if len(priceFilterPipeline) > 0 {
				pipeline = append(pipeline, priceFilterPipeline...)
			}
		}

		if productFilters.UnitIds != nil {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"priceItem.unitId": bson.M{"$in": productFilters.UnitIds}}})
		}
	}

	if args.AddExtraOperations {
		if cfg != nil && cfg.User != nil {
			pipeline = append(pipeline, bson.M{"$addFields": bson.M{"userId": cfg.User.Id}})
		}

		pipeline = base_filters.BuildPopulatePipeline(productFilters.PopFields, enums.PRODUCT, pipeline, ctx)

		pipeline = base_filters.BuildSortingPipeline(productFilters.Sorts, pipeline)
	}

	return pipeline
}
