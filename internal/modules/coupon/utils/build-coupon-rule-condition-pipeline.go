package coupon_utils

import (
	"context"
	coupon_dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuildOrderCouponRuleConditionPipelineArgs struct {
	Conditions []coupon_models.CouponRuleCondition
	Coupon     coupon_models.CouponResponse
}

type BuildOrderCouponRuleConditionPipelineResult struct {
	MatchStage                      []bson.M
	NumberOfOrdersConditionPipeline []bson.M
	Operator                        string
}

func BuildOrderCouponRuleConditionPipeline(ctx context.Context, args BuildOrderCouponRuleConditionPipelineArgs) BuildOrderCouponRuleConditionPipelineResult {
	matchStage := []bson.M{}
	numberOfOrdersConditionPipeline := []bson.M{}
	for _, condition := range args.Conditions {

		switch condition.Type {
		case coupon_enums.NUMBER_OF_ORDERS:
			{
				var operator = "$" + string(condition.Operator)
				numberOfOrdersConditionPipeline = append(numberOfOrdersConditionPipeline, NumberOfOrdersConditionPipeline(operator, condition.Value)...)
			}

		case coupon_enums.BRANCH:
			{
				operator := "$" + string(condition.Operator)
				matchStage = append(matchStage, bson.M{"branchId": bson.M{operator: condition.Value}})
			}
		case coupon_enums.ORDER_PAYMENT_STATUS:
			{
				operator := "$" + string(condition.Operator)
				matchStage = append(matchStage, bson.M{"paymentInfo.status": bson.M{operator: condition.Value}})

			}
		case coupon_enums.DELIVERY_METHOD:
			{

				operator := "$" + string(condition.Operator)
				matchStage = append(matchStage, bson.M{"type": bson.M{operator: condition.Value}})
			}
		case coupon_enums.ORDER_TIMEFRAME:
			{
				operator := "$" + string(condition.Operator)
				usage := coupon_dtos.UsageDto{}

				if usageMap, ok := condition.Value.(primitive.D); ok {
					for _, v := range usageMap {
						switch v.Key {
						case "value":
							if val, ok := v.Value.(float64); ok {
								usage.Value = uint16(val)
							} else if val, ok := v.Value.(int64); ok {
								usage.Value = uint16(val)
							} else if val, ok := v.Value.(int); ok {
								usage.Value = uint16(val)
							}
						case "type":
							if val, ok := v.Value.(string); ok {
								usage.Type = coupon_enums.CouponExpirationType(val)
							}
						}
					}
				}

				extractUsageDtoTimeResult := coupon_dtos.ExtractUsageDtoTime(coupon_dtos.ExtractUsageDtoTimeArgs{
					Usage:    usage,
					Operator: condition.Operator,
				})
				matchStage = append(matchStage, bson.M{"createdAt": bson.M{operator: extractUsageDtoTimeResult}})
			}
		case coupon_enums.SUB_TOTAL_AMOUNT:
			{
				operator := "$" + string(condition.Operator)
				priceItem := product_dtos.PriceItemDto{}

				if priceMap, ok := condition.Value.(primitive.D); ok {
					for _, v := range priceMap {
						switch v.Key {
						case "value":
							if val, ok := v.Value.(float64); ok {
								priceItem.Value = val
							} else if val, ok := v.Value.(int64); ok {
								priceItem.Value = float64(val)
							} else if val, ok := v.Value.(int); ok {
								priceItem.Value = float64(val)
							}
						case "currencyCode":
							if val, ok := v.Value.(string); ok {
								priceItem.CurrencyCode = val
							}
						}
					}
				}

				matchStage = append(matchStage, bson.M{"subTotalPrice": bson.M{operator: priceItem.Value}})
				matchStage = append(matchStage, bson.M{"currencyCode": priceItem.CurrencyCode})
			}
		}
	}

	var operator = extractCouponRuleAllowedIf(args.Coupon.Rule.Conditions.AllowedIf, args.Coupon.Rule.Conditions.When)

	return BuildOrderCouponRuleConditionPipelineResult{
		MatchStage:                      matchStage,
		NumberOfOrdersConditionPipeline: numberOfOrdersConditionPipeline,
		Operator:                        operator,
	}
}

/*
* extractCouponRuleAllowedIf determines the MongoDB operator to use based on the coupon rule's allowedIf condition and when flag.
* It converts the logical conditions (ALL, ANY, NONE) into appropriate MongoDB operators ($and, $or, $nor).
*
* Parameters:
*   - allowedIf: The coupon rule condition type (ALL, ANY, NONE) that determines how conditions should be combined
*   - when: A boolean flag that inverts the logic when false
*
* Returns:
*   - string: The MongoDB operator to use in the aggregation pipeline ($and, $or, or $nor)
*
* The function implements the following logic:
* When when=true:
*   - ALL -> $and (all conditions must match)
*   - ANY -> $or (any condition can match)
*   - NONE -> $nor (no conditions should match)
*
* When when=false:
*   - ALL -> $nor (none of the conditions should match)
*   - ANY -> $and (all conditions must not match)
*   - NONE -> $and (all conditions must not match)
 */
func extractCouponRuleAllowedIf(allowedIf coupon_enums.CouponRuleConditionAllowedIf, when bool) string {

	if when {
		switch allowedIf {
		case coupon_enums.ALL:
			return "$and"
		case coupon_enums.ANY:
			return "$or"
		case coupon_enums.NONE:
			return "$nor"
		}
	} else {
		switch allowedIf {
		case coupon_enums.ALL:
			return "$nor"
		case coupon_enums.ANY:
			return "$and"
		case coupon_enums.NONE:
			return "$and"
		}
	}

	return "$and"
}
