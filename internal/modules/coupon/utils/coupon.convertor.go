package coupon_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	coupon_dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
)

func ConvertToCouponModel(createCouponDto coupon_dtos.CreateCouponDto, additionalOptions types.AdditionalOptions) coupon_models.Coupon {
	var upTo *[]product_models.PriceItem
	if createCouponDto.Discount.UpTo != nil {
		priceItems := product_utils.ConvertToPriceItems(*createCouponDto.Discount.UpTo)
		upTo = &priceItems
	}

	// Convert Nullable OrderType to *OrderType
	// Always set orderType - if not provided, set to nil (will be stored as null in DB)
	var orderType *order_enums.OrderType
	if createCouponDto.OrderType.IsProvided() {
		// If explicitly null or has value, use the value (which may be nil)
		orderType = createCouponDto.OrderType.GetValue()
	} else {
		// If not provided, set to nil (will be stored as null in DB for create operations)
		orderType = nil
	}

	coupon := coupon_models.Coupon{
		BaseCoupon: coupon_models.BaseCoupon{
			Title:       createCouponDto.Title,
			IsActive:    createCouponDto.IsActive,
			Code:        createCouponDto.Code,
			Description: createCouponDto.Description,
			Type:        createCouponDto.Type,
			AppliedOn:   createCouponDto.AppliedOn,
			OrderType:   orderType,
			Discount: product_models.Discount{
				Type:         createCouponDto.Discount.Type,
				Value:        createCouponDto.Discount.Value,
				CurrencyCode: createCouponDto.Discount.CurrencyCode,
				UpTo:         upTo,
			},
			TotalUsages: createCouponDto.TotalUsages,

			AppliedAutomatically: createCouponDto.AppliedAutomatically,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}

	if createCouponDto.Usage != nil {
		coupon.Usage = &coupon_models.CouponUsage{

			MaxUsagePerCustomer: createCouponDto.Usage.MaxUsagePerCustomer,
			MinimumSubtotal:     createCouponDto.Usage.MinimumSubtotal,
			StartDate:           createCouponDto.Usage.StartDate,
			EndDate:             createCouponDto.Usage.EndDate,
			RelatedTo:           createCouponDto.Usage.RelatedTo,
		}

		if createCouponDto.Usage.Interval != nil {
			coupon.Usage.Interval = &coupon_models.Usage{
				Type:  createCouponDto.Usage.Interval.Type,
				Value: createCouponDto.Usage.Interval.Value,
			}
		}
	}

	if createCouponDto.Expiration != nil {
		coupon.Expiration = &coupon_models.Usage{
			Type:  createCouponDto.Expiration.Type,
			Value: createCouponDto.Expiration.Value,
		}

	}
	if createCouponDto.Rule != nil {
		coupon.Rule = &coupon_models.CouponRule{
			IsActive:  createCouponDto.Rule.IsActive,
			Title:     createCouponDto.Rule.Title,
			StartDate: createCouponDto.Rule.StartDate,
			EndDate:   createCouponDto.Rule.EndDate,
			Priority:  createCouponDto.Rule.Priority,
			Conditions: coupon_models.CouponRuleConditions{
				AllowedIf: createCouponDto.Rule.Conditions.AllowedIf,
				When:      createCouponDto.Rule.Conditions.When,
				Values:    convertToCouponRuleConditions(createCouponDto.Rule.Conditions.Values, additionalOptions),
			},
		}

	}
	return coupon

}

func convertToCouponRuleConditions(couponRuleConditionsDto []coupon_dtos.CouponRuleConditionDto, additionalOptions types.AdditionalOptions) []coupon_models.CouponRuleCondition {
	couponRuleConditions := []coupon_models.CouponRuleCondition{}
	for _, couponRuleConditionDto := range couponRuleConditionsDto {
		couponRuleConditions = append(couponRuleConditions, coupon_models.CouponRuleCondition{
			Type:     couponRuleConditionDto.Type,
			Operator: couponRuleConditionDto.Operator,
			Value:    couponRuleConditionDto.Value,
		})
	}
	return couponRuleConditions
}
