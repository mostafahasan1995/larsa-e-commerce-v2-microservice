package coupon_dtos

import (
	"net/http"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
)

type UsageDto struct {
	Type  coupon_enums.CouponExpirationType `json:"type" validate:"required,oneof=days weeks months years"`
	Value uint16                            `json:"value" validate:"required,min=0"`
}

type CouponUsageDto struct {
	Interval            *UsageDto                     `json:"interval" validate:"omitempty"`
	MaxUsagePerCustomer *uint16                       `json:"maxUsagePerCustomer" validate:"omitempty,min=0"`
	MinimumSubtotal     *float64                      `json:"minimumSubtotal" validate:"omitempty,min=0"`
	StartDate           *time.Time                    `json:"startDate" validate:"omitempty"`
	EndDate             *time.Time                    `json:"endDate" validate:"omitempty"`
	RelatedTo           *coupon_enums.CouponRelatedTo `json:"relatedTo,omitempty" validate:"omitempty,oneof=feedback"`
}

type CouponRuleConditionDto struct {
	Type     coupon_enums.CouponRuleConditionType `json:"type" validate:"required,oneof=order_payment_status delivery_method branch sub_total_amount order_timeframe number_of_orders"`
	Operator enums.Operator                       `json:"operator" validate:"required,oneof=eq ne gt gte lt lte in not_in"`
	Value    interface{}                          `json:"value" validate:"required"`
}

type CouponRuleConditionsDto struct {
	IsActive  bool                                      `json:"isActive"`
	AllowedIf coupon_enums.CouponRuleConditionAllowedIf `json:"allowedIf" validate:"required,oneof=all any none"`
	When      bool                                      `json:"when" `
	Values    []CouponRuleConditionDto                  `json:"values" validate:"required,min=1,dive"`
}

type CouponRuleDto struct {
	IsActive   bool                      `json:"isActive"`
	Title      types.Localizable[string] `json:"title" validate:"required"`
	StartDate  time.Time                 `json:"startDate" validate:"required,ltefield=EndDate"`
	EndDate    time.Time                 `json:"endDate" validate:"required,gtefield=StartDate"`
	Priority   uint16                    `json:"priority" validate:"required,gt=0"`
	Conditions CouponRuleConditionsDto   `json:"conditions" validate:"required"`
}

type CouponDiscountDto struct {
	product_dtos.DiscountDto
	UpTo *[]product_dtos.PriceItemDto `json:"upTo" validate:"omitempty,min=1,dive,required"`
}

type CreateCouponDto struct {
	Title                types.Localizable[string]             `json:"title" validate:"required"`
	IsActive             *bool                                 `json:"isActive"`
	Code                 string                                `json:"code" validate:"required"`
	Description          *types.Localizable[string]            `json:"description"`
	Type                 coupon_enums.CouponType               `json:"type" validate:"required,oneof=auto_generated manual feedback_generated"`
	OrderType            types.Nullable[order_enums.OrderType] `json:"orderType" validate:"omitempty"`
	AppliedOn            coupon_enums.CouponAppliedOn          `json:"appliedOn" validate:"required,oneof=general product"`
	Expiration           *UsageDto                             `json:"expiration" validate:"omitempty"`
	Discount             CouponDiscountDto                     `json:"discount" validate:"required"`
	TotalUsages          *uint64                               `json:"totalUsages" validate:"omitempty"`
	Usage                *CouponUsageDto                       `json:"usage" validate:"omitempty"`
	AppliedAutomatically bool                                  `json:"appliedAutomatically"`
	Rule                 *CouponRuleDto                        `json:"rule" validate:"omitempty"`
}

func (dto *CreateCouponDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {

	// Validate OrderType if provided
	if dto.OrderType.IsProvided() && dto.OrderType.GetValue() != nil {
		orderTypeValue := *dto.OrderType.GetValue()
		if orderTypeValue != order_enums.PICKUP && orderTypeValue != order_enums.DELIVERY {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "OrderType must be either 'pickup' or 'delivery'",
			}
		}
	}

	if dto.Type == coupon_enums.AUTO_GENERATED {

		if dto.Rule == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons must have rule",
			}
		}

		if dto.OrderType.IsProvided() {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons can not have orderType",
			}
		}

		// if dto.Expiration == nil {
		// 	return nil, &helpers.ApiError{
		// 		StatusCode: http.StatusBadRequest,
		// 		Message:    "Auto generated coupons must have expiration",
		// 	}
		// }

		if dto.IsActive != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons can not have isActive",
			}
		}

		if dto.Usage != nil && dto.Usage.StartDate != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons can not have startDate",
			}
		}

		if dto.Usage != nil && dto.Usage.EndDate != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons can not have endDate",
			}
		}

		if dto.Usage != nil && dto.Usage.MinimumSubtotal != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Auto generated coupons can not have minimumSubtotal",
			}
		}

	} else if dto.Type == coupon_enums.MANUAL {

		if dto.Rule != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Manual generated coupons can not have rule",
			}
		}

		if dto.Expiration != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Manual generated coupons can not have expiration",
			}
		}

		// if dto.Usage.StartDate == nil {
		// 	return nil, &helpers.ApiError{
		// 		StatusCode: http.StatusBadRequest,
		// 		Message:    "Manual generated coupons must have startDate",
		// 	}
		// }

		// if dto.Usage.EndDate == nil {
		// 	return nil, &helpers.ApiError{
		// 		StatusCode: http.StatusBadRequest,
		// 		Message:    "Manual generated coupons must have endDate",
		// 	}
		// }

		// if dto.Usage.MinimumSubtotal == nil {
		// 	return nil, &helpers.ApiError{
		// 		StatusCode: http.StatusBadRequest,
		// 		Message:    "Manual generated coupons must have minimumSubtotal",
		// 	}
		// }

	} else if dto.Type == coupon_enums.FEEDBACK_GENERATED {

		if dto.Rule != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have rule",
			}
		}

		if dto.OrderType.IsProvided() {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have orderType",
			}
		}

		if dto.IsActive != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have isActive",
			}
		}

		if dto.Usage != nil && dto.Usage.StartDate != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have startDate",
			}
		}

		if dto.Usage != nil && dto.Usage.EndDate != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have endDate",
			}
		}

		if dto.Usage != nil && dto.Usage.MinimumSubtotal != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Feedback generated coupons can not have minimumSubtotal",
			}
		}

	}

	// Validate date range when both StartDate and EndDate are provided
	if dto.Usage != nil && dto.Usage.StartDate != nil && dto.Usage.EndDate != nil {
		if dto.Usage.StartDate.After(*dto.Usage.EndDate) || dto.Usage.StartDate.Equal(*dto.Usage.EndDate) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Start date must be before end date",
			}
		}
	}

	if dto.Discount.Type == product_enums.PERCENTAGE && dto.Discount.Value > 100 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Discount value must be less than 100",
		}
	}

	if dto.Discount.UpTo != nil && dto.Discount.Type == product_enums.FIXED {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Fixed discount can not have upTo",
		}
	}

	if dto.Rule != nil {
		newConditions := []CouponRuleConditionDto{}
		for _, condition := range dto.Rule.Conditions.Values {
			couponRuleConditionDto, err := validateCouponRuleConditionTypeAndOperator(ValidateCouponRuleConditionTypeAndOperatorArgs{
				CandidateCouponRuleConditionDto: condition,
				CustomValidatorArgs:             args,
			})
			if err != nil {
				return nil, err
			}

			newConditions = append(newConditions, *couponRuleConditionDto)
		}

		dto.Rule.Conditions.Values = newConditions
	}
	return nil, nil
}
