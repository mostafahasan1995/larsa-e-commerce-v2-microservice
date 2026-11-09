package survey_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	coupon_dtos "larsa-e-commerce-microservice-v2/internal/modules/coupon/dtos"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	coupon_services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	settings_general_services "larsa-e-commerce-microservice-v2/internal/modules/settings/general/services"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_utils "larsa-e-commerce-microservice-v2/internal/modules/survey/utils"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebFeedbackService interface {
	Create(ctx context.Context, dto *survey_dtos.CreateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError)
}

type webFeedback struct {
	repo                   FeedbackRepo
	usersgw                *gateway.UsersGw
	emailgw                *gateway.EmailGw
	profileService         profile_services.ProfileService
	couponService          coupon_services.CouponService
	emailTemplateService   SurveyEmailTemplateService
	couponRepo             base_repos.BaseRepo[coupon_models.Coupon, coupon_models.CouponResponse]
	generalSettingsService settings_general_services.GeneralSettingsService
}

func NewWebFeedbackService(i *do.Injector) (WebFeedbackService, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &webFeedback{
		repo:                   do.MustInvoke[FeedbackRepo](i),
		usersgw:                do.MustInvoke[*gateway.UsersGw](i),
		emailgw:                do.MustInvoke[*gateway.EmailGw](i),
		profileService:         do.MustInvoke[profile_services.ProfileService](i),
		couponService:          do.MustInvoke[coupon_services.CouponService](i),
		emailTemplateService:   do.MustInvoke[SurveyEmailTemplateService](i),
		couponRepo:             *base_repos.NewBaseRepo[coupon_models.Coupon, coupon_models.CouponResponse](db, enums.COUPON),
		generalSettingsService: do.MustInvoke[settings_general_services.GeneralSettingsService](i),
	}, nil
}

func (f *webFeedback) Create(ctx context.Context, dto *survey_dtos.CreateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	feedbackModel, err := survey_utils.ConvertCreateFeedbackDTOToModel(dto, ctx)
	if err != nil {
		return nil, err
	}

	feedback, createErr := f.repo.Create(ctx, feedbackModel)
	if createErr != nil {
		return nil, createErr
	}

	// Generate coupon if settings allow it
	couponCode, generateErr := f.generateFeedbackCoupon(ctx, dto)
	if generateErr != nil {
		// Log error but don't fail the feedback creation
		fmt.Printf("Error generating feedback coupon: %v\n", generateErr)
	} else if couponCode != nil {
		feedback.CouponCode = couponCode
	}

	return feedback, nil
}

func (f *webFeedback) generateFeedbackCoupon(ctx context.Context, dto *survey_dtos.CreateFeedbackDTO) (*string, error) {
	// Get general settings
	settings, err := f.generalSettingsService.GetGeneralSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get general settings: %v", err)
	}

	// Check if feedback coupon is enabled
	if !settings.FeedbackCoupon.IsEnabled {
		return nil, nil
	}

	// Generate unique coupon code
	couponCode := survey_utils.GenerateCouponCode()

	orderType := order_enums.DELIVERY
	// Create coupon DTO
	couponDto := &coupon_dtos.CreateCouponDto{
		Title: types.Localizable[string]{
			enums.EN: "Feedback Reward",
			enums.AR: "مكافأة التقييم",
		},
		Code:      couponCode,
		Type:      coupon_enums.FEEDBACK_GENERATED,
		AppliedOn: coupon_enums.GENERAL,
		Description: &types.Localizable[string]{
			enums.EN: "Thank you for your feedback!",
			enums.AR: "شكراً لك على تقييمك!",
		},
		OrderType: types.Nullable[order_enums.OrderType]{Value: &orderType},
		Discount: coupon_dtos.CouponDiscountDto{
			DiscountDto: product_dtos.DiscountDto{
				Type:         product_enums.PERCENTAGE,
				Value:        settings.FeedbackCoupon.DiscountPercentage,
				CurrencyCode: &[]string{"USD"}[0], // Default currency
			},
		},
		Expiration: &coupon_dtos.UsageDto{
			Type:  coupon_enums.DAYS,
			Value: uint16(settings.FeedbackCoupon.ExpirationDays),
		},
		TotalUsages: &[]uint64{1}[0], // One-time use
		Usage: &coupon_dtos.CouponUsageDto{
			MaxUsagePerCustomer: &[]uint16{1}[0], // One-time use per customer
		},
		AppliedAutomatically: false,
	}

	// Get the existing AppCfg from the context
	cfg, appCfgErr := util.GetReqAppCfg(ctx)
	if appCfgErr != nil {
		return nil, fmt.Errorf("failed to get app config: %v", appCfgErr)
	}

	// Create a context with system user for coupon creation
	systemUser := &types.User{
		Id: primitive.NilObjectID, // System user ID
		UserData: common.User{
			Id:    primitive.NilObjectID.Hex(),
			Email: "system@larsa.io",
		},
		Type: enums.SuperAdmin,
	}

	// Create new context with system user and update AppCfg
	ctx = util.SetReqUser(ctx, systemUser)

	// Create new AppCfg with system user
	newCfg := &types.AppCfg{
		Db:      cfg.Db,
		Lang:    cfg.Lang,
		Hp:      cfg.Hp,
		User:    systemUser,
		ReqCaps: cfg.ReqCaps,
	}

	// Set the new AppCfg in context
	ctx = context.WithValue(ctx, util.ReqAppCfg, newCfg)

	// Create the coupon with proper context
	_, createErr := f.couponService.Create(ctx, couponDto)
	if createErr != nil {
		return nil, fmt.Errorf("failed to create feedback coupon: %v", createErr)
	}

	return &couponCode, nil
}
