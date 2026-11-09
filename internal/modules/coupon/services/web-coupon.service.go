package coupon_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"

	"github.com/samber/do"
)

type WebCouponRepo interface {
	GetMine(ctx context.Context) (*[]models.CouponResponse, *helpers.ApiError)
	ConvertToUserCoupon(ctx context.Context, coupon *models.CouponResponse) *helpers.ApiError
}

type WebCouponService interface {
	GetMine(ctx context.Context) (*[]models.CouponResponse, *helpers.ApiError)
	CheckCouponRedemption(ctx context.Context, couponCode string, subTotal *float64, orderType *order_enums.OrderType) (*models.CouponResponse, *helpers.ApiError)
}

type webCoupon struct {
	repo                   WebCouponRepo
	couponRuleCheckService CouponRuleCheckerService
}

func NewWebCouponService(i *do.Injector) (WebCouponService, error) {
	return &webCoupon{
		repo:                   do.MustInvoke[WebCouponRepo](i),
		couponRuleCheckService: *do.MustInvoke[*CouponRuleCheckerService](i),
	}, nil
}

func (webCoupon *webCoupon) GetMine(ctx context.Context) (*[]models.CouponResponse, *helpers.ApiError) {
	return webCoupon.repo.GetMine(ctx)
}

func (webCoupon *webCoupon) CheckCouponRedemption(ctx context.Context, couponCode string, subTotal *float64, orderType *order_enums.OrderType) (*models.CouponResponse, *helpers.ApiError) {
	return webCoupon.couponRuleCheckService.CheckCouponRedemption(ctx, couponCode, nil, subTotal, orderType)
}

func (webCoupon *webCoupon) ConvertToUserCoupon(ctx context.Context, coupon *models.CouponResponse) *helpers.ApiError {
	return webCoupon.repo.ConvertToUserCoupon(ctx, coupon)
}
