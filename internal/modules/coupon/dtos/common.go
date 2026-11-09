package coupon_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var allowedConditionCouponTypeToOperatorMap = map[coupon_enums.CouponRuleConditionType][]enums.Operator{
	coupon_enums.NUMBER_OF_ORDERS: {
		enums.EQ,
		enums.NE,
		enums.GT,
		enums.GTE,
		enums.LT,
		enums.LTE,
	},
	coupon_enums.ORDER_PAYMENT_STATUS: {
		enums.EQ,
		enums.NE,
	},
	coupon_enums.DELIVERY_METHOD: {
		enums.EQ,
		enums.NE,
	},
	coupon_enums.BRANCH: {
		enums.EQ,
		enums.NE,
		enums.IN,
		enums.NOT_IN,
	},
	coupon_enums.ORDER_TIMEFRAME: {
		enums.LT,
		enums.LTE,
	},
	coupon_enums.SUB_TOTAL_AMOUNT: {
		enums.EQ,
		enums.NE,
		enums.GT,
		enums.GTE,
		enums.LT,
		enums.LTE,
	},
}

type ValidateCouponRuleConditionTypeAndOperatorArgs struct {
	CandidateCouponRuleConditionDto CouponRuleConditionDto
	CustomValidatorArgs             middleware.CustomValidatorArgs
}

func validateCouponRuleConditionTypeAndOperator(args ValidateCouponRuleConditionTypeAndOperatorArgs) (*CouponRuleConditionDto, *helpers.ApiError) {
	err := ""
	couponRuleConditionDto := args.CandidateCouponRuleConditionDto
	allowedOperators, ok := allowedConditionCouponTypeToOperatorMap[couponRuleConditionDto.Type]
	if !ok {
		err = "You've entered an invalid condition type"
		return nil, &helpers.ApiError{
			Message:    err,
			StatusCode: http.StatusBadRequest,
		}
	}

	isAllowed := util.Contains(allowedOperators, couponRuleConditionDto.Operator)
	if !isAllowed {
		err = "You've entered an invalid operator"
		return nil, &helpers.ApiError{
			Message:    err,
			StatusCode: http.StatusBadRequest,
		}
	}

	if couponRuleConditionDto.Value == nil {
		err = "You've entered an invalid value"
		return nil, &helpers.ApiError{
			Message:    err,
			StatusCode: http.StatusBadRequest,
		}
	}

	switch couponRuleConditionDto.Type {
	case coupon_enums.NUMBER_OF_ORDERS:
		{
			isNumber := util.IsNumber(couponRuleConditionDto.Value)
			if !isNumber {
				err = "Value must be a number"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			return &couponRuleConditionDto, nil
		}

	case coupon_enums.ORDER_PAYMENT_STATUS:
		{

			_, ok := couponRuleConditionDto.Value.(string)
			if !ok {
				err = "Value must be a string"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			validPaymentStatus := []string{
				string(order_enums.PAID),
				string(order_enums.PARTIALLY_PAID),
				string(order_enums.AGENT_PAID),
				string(order_enums.UNPAID),
			}

			if !util.Contains(validPaymentStatus, couponRuleConditionDto.Value.(string)) {
				err = "Value must be a valid payment status"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			return &couponRuleConditionDto, nil
		}

	case coupon_enums.DELIVERY_METHOD:
		{
			value, ok := couponRuleConditionDto.Value.(string)
			if !ok {
				err = "Value must be a string"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			validDeliveryMethods := []string{
				string(order_enums.PICKUP),
				string(order_enums.DELIVERY),
			}

			if !util.Contains(validDeliveryMethods, value) {
				err = "Value must be a valid delivery method"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			return &couponRuleConditionDto, nil
		}

	case coupon_enums.BRANCH:
		{
			isString := util.IsString(couponRuleConditionDto.Value) && couponRuleConditionDto.Operator == enums.EQ
			isArray := util.IsArray(couponRuleConditionDto.Value) && couponRuleConditionDto.Operator == enums.IN
			if !isString && !isArray {
				err = "Value must be a string or array"
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

			branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.CustomValidatorArgs.DB, enums.BRANCH)

			if isString {
				objectId, objectIdErr := primitive.ObjectIDFromHex(couponRuleConditionDto.Value.(string))
				if objectIdErr != nil {
					err = "Value must be a valid object id"
					return nil, &helpers.ApiError{
						Message:    err,
						StatusCode: http.StatusBadRequest,
					}
				}

				couponRuleConditionDto.Value = objectId
				_, getErr := branchRepo.GetById(args.CustomValidatorArgs.Context, objectId, []bson.M{})
				if getErr != nil {
					err = getErr.Message
					return nil, &helpers.ApiError{
						Message:    err,
						StatusCode: http.StatusBadRequest,
					}
				}

			} else if isArray {
				newValueAsSlice := []primitive.ObjectID{}
				valueSlice := couponRuleConditionDto.Value.([]interface{})
				for _, value := range valueSlice {
					strValue, ok := value.(string)
					if !ok {
						err = "Array values must be strings"
						return nil, &helpers.ApiError{
							Message:    err,
							StatusCode: http.StatusBadRequest,
						}
					}
					objectId, objectIdErr := primitive.ObjectIDFromHex(strValue)
					if objectIdErr != nil {
						err = "Value must be a valid object id"
						return nil, &helpers.ApiError{
							Message:    err,
							StatusCode: http.StatusBadRequest,
						}
					}

					newValueAsSlice = append(newValueAsSlice, objectId)
				}

				couponRuleConditionDto.Value = newValueAsSlice

				_, getErr := branchRepo.GetByIds(args.CustomValidatorArgs.Context, newValueAsSlice, []bson.M{}, base_utils.AdditionalOptions{
					ReturnError: true,
				})
				if getErr != nil {
					err = getErr.Message
					return nil, &helpers.ApiError{
						Message:    err,
						StatusCode: http.StatusBadRequest,
					}
				}
			}
		}

	case coupon_enums.ORDER_TIMEFRAME:
		{

			structErr := middleware.ValidateInstance(UsageDto{}, couponRuleConditionDto.Value, args.CustomValidatorArgs.Injector)
			if structErr != nil {
				err = structErr.Error()
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}
		}
	case coupon_enums.SUB_TOTAL_AMOUNT:
		{

			structErr := middleware.ValidateInstance(product_dtos.PriceItemDto{}, couponRuleConditionDto.Value, args.CustomValidatorArgs.Injector)
			if structErr != nil {
				err = structErr.Error()
				return nil, &helpers.ApiError{
					Message:    err,
					StatusCode: http.StatusBadRequest,
				}
			}

		}
	}

	if err != "" {
		return nil, &helpers.ApiError{
			Message:    err,
			StatusCode: http.StatusBadRequest,
		}
	}

	return &couponRuleConditionDto, nil
}

type ExtractUsageDtoTimeArgs struct {
	Usage    UsageDto
	Operator enums.Operator
}

func ExtractUsageDtoTime(args ExtractUsageDtoTimeArgs) time.Time {
	now := time.Now()
	var subtractedTime time.Time

	switch args.Usage.Type {
	case coupon_enums.DAYS:
		subtractedTime = now.AddDate(0, 0, -int(args.Usage.Value))
	case coupon_enums.WEEKS:
		subtractedTime = now.AddDate(0, 0, -int(args.Usage.Value*7))
	case coupon_enums.MONTHS:
		subtractedTime = now.AddDate(0, -int(args.Usage.Value), 0)
	case coupon_enums.YEARS:
		subtractedTime = now.AddDate(-int(args.Usage.Value), 0, 0)
	}

	return subtractedTime
}
