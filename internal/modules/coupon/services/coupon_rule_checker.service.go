package coupon_services

import (
	"context"
	"fmt"
	"net/http"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	config_repos "larsa-e-commerce-microservice-v2/internal/modules/config/repos"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	coupon_enums "larsa-e-commerce-microservice-v2/internal/modules/coupon/enums"
	coupon_filters "larsa-e-commerce-microservice-v2/internal/modules/coupon/filters"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	coupon_utils "larsa-e-commerce-microservice-v2/internal/modules/coupon/utils"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_filters "larsa-e-commerce-microservice-v2/internal/modules/order/filters"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"github.com/robfig/cron/v3"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CouponRuleChecker interface {
	CheckCouponRedemption(ctx context.Context, CouponCode string, appliedAutomatically *bool, subTotal *float64, orderType *order_enums.OrderType) (*coupon_models.CouponResponse, *helpers.ApiError)
	GenerateUserCouponBasedOnCouponRules(databaseName string)
	Start()
}

type CouponRuleCheckerService struct {
	userCouponRepo base_repos.BaseRepo[coupon_models.UserCoupon, coupon_models.UserCouponResponse]
	couponRepo     base_repos.BaseRepo[coupon_models.Coupon, coupon_models.CouponResponse]
	branchRepo     base_repos.BaseRepo[product_models.Branch, product_models.BranchResponse]
	settingRepo    config_services.SettingRepo
	orderRepo      base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
}

func NewCouponRuleChecker(i *do.Injector) (*CouponRuleCheckerService, error) {
	db := do.MustInvoke[*mongo.Client](i)

	settingRepo, err := config_repos.NewSettingRepo(i)
	if err != nil {
		return nil, err
	}
	return &CouponRuleCheckerService{
		settingRepo:    settingRepo,
		branchRepo:     *base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](db, enums.BRANCH),
		couponRepo:     *base_repos.NewBaseRepo[coupon_models.Coupon, coupon_models.CouponResponse](db, enums.COUPON),
		userCouponRepo: *base_repos.NewBaseRepo[coupon_models.UserCoupon, coupon_models.UserCouponResponse](db, enums.USER_COUPON),
		orderRepo:      *base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
	}, nil
}

func (checker *CouponRuleCheckerService) Start() {
	// Run the check every hour
	c := cron.New()
	c.AddFunc("@every 30m", func() {
		fmt.Println("Checking coupon rules")
		dbs, err := checker.settingRepo.GetRelevantDbs()
		if err != nil {
			helpers.DefaultError(err)
		}
		for _, db := range dbs {
			// if db == "604674966fe156529e3e03e9" {
			checker.GenerateUserCouponBasedOnCouponRules(nil, db)
			// }
		}
	})

	c.Start()
}

func (checker *CouponRuleCheckerService) GenerateUserCouponBasedOnCouponRules(ctx *context.Context, databaseName string) {
	if ctx == nil {
		context := checker.settingRepo.AddContext(databaseName)
		ctx = &context
	}
	couponFilters, err := coupon_filters.NewCouponFilters("")
	if err != nil {
		helpers.DefaultError(err)
		return
	}
	matchStage := bson.M{
		"$and": []bson.M{
			{
				"type":          coupon_enums.AUTO_GENERATED,
				"rule.isActive": true,
			},
			{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$lt": []interface{}{"$rule.startDate", "$$NOW"}},
						{"$gt": []interface{}{"$rule.endDate", "$$NOW"}},
					},
				},
			},
		},
	}
	pipeline := couponFilters.BuildPipeline(matchStage, *ctx, false)

	coupons, getAllErr := checker.couponRepo.GetAll(*ctx, pipeline)
	if getAllErr != nil {
		helpers.DefaultError(getAllErr)
	}

	if len(coupons) > 0 {
		for _, coupon := range coupons {

			var conditions = coupon.Rule.Conditions.Values

			pipeline := coupon_utils.BuildOrderCouponRuleConditionPipeline(*ctx, coupon_utils.BuildOrderCouponRuleConditionPipelineArgs{
				Conditions: conditions,
				Coupon:     coupon,
			})

			orders, _ := checker.orderRepo.FindAll(*ctx, pipeline.NumberOfOrdersConditionPipeline)
			userIds := []primitive.ObjectID{}
			if len(orders) > 0 {
				for _, order := range orders {
					var orderMap bson.M
					orderBytes, _ := bson.Marshal(order)
					bson.Unmarshal(orderBytes, &orderMap)
					if rawUserId, ok := orderMap["userId"]; ok {
						switch v := rawUserId.(type) {
						case primitive.ObjectID:
							userIds = append(userIds, v)
						case string:
							if oid, err := primitive.ObjectIDFromHex(v); err == nil {
								userIds = append(userIds, oid)
							}
						}
					}
				}
			}

			pipeline.MatchStage = append(pipeline.MatchStage, bson.M{
				"buyerInfo.id": bson.M{"$in": userIds},
			})
			orderPipeline := []bson.M{
				{
					"$match": bson.M{
						pipeline.Operator: pipeline.MatchStage,
					},
				},
			}

			err := checker.createUserCouponsFromOrders(*ctx, coupon_utils.CreateUserCouponsFromOrdersArgs{
				Coupon:        coupon,
				OrderPipeline: orderPipeline,
			})
			if err != nil {
				helpers.DefaultError(err)
			}

		}
	}
}

func (checker *CouponRuleCheckerService) createUserCouponsFromOrders(ctx context.Context, args coupon_utils.CreateUserCouponsFromOrdersArgs) *helpers.ApiError {
	orderFilters, _ := order_filters.NewOrderFilters("")
	orderPipeline, _ := orderFilters.BuildPipeline(bson.M{}, ctx, false)
	orderPipeline = append(orderPipeline, args.OrderPipeline...)
	orders, err := checker.orderRepo.GetAll(ctx, orderPipeline)
	if err != nil {
		return helpers.DefaultError(err)
	}

	if len(orders) == 0 {
		return nil
	}

	userCoupons := []coupon_models.UserCoupon{}
	userIds := []primitive.ObjectID{}
	couponCodes := []string{}
	for _, order := range orders {
		userIds = append(userIds, order.BuyerInfo.Id)
		couponCodes = append(couponCodes, args.Coupon.Code)
	}

	matchStage := bson.M{
		"userId":     bson.M{"$in": userIds},
		"couponCode": bson.M{"$in": couponCodes},
	}

	userCouponFilters, _ := coupon_filters.NewUserCouponFilters("")
	userCouponPipeline := userCouponFilters.BuildPipeline(matchStage, ctx, false)
	storedUserCoupons, err := checker.userCouponRepo.GetAll(ctx, userCouponPipeline)
	if err != nil {
		return helpers.DefaultError(err)
	}

	storedUserCouponsMap := make(map[primitive.ObjectID]bool)
	for _, userCoupon := range storedUserCoupons {
		storedUserCouponsMap[userCoupon.UserId] = true
	}

	for _, order := range orders {
		if !storedUserCouponsMap[order.BuyerInfo.Id] {
			userCoupons = append(userCoupons, coupon_utils.ConvertToUserCouponModel(coupon_utils.ConvertToUserCouponModelArgs{
				UserId:     order.BuyerInfo.Id,
				CouponCode: args.Coupon.Code,
			}, types.AdditionalOptions{}))
		}
	}

	/*
	* Only create coupons if we have new ones to create
	 */
	if len(userCoupons) > 0 {
		_, createErr := checker.userCouponRepo.CreateMany(ctx, userCoupons)
		if createErr != nil {
			return helpers.DefaultError(createErr)
		}
	}

	return nil
}

func (checker *CouponRuleCheckerService) CheckCouponRedemption(ctx context.Context, CouponCode string, appliedAutomatically *bool, subTotal *float64, orderType *order_enums.OrderType) (*coupon_models.CouponResponse, *helpers.ApiError) {

	getMineCouponPipelineArgs := coupon_utils.GetMineCouponPipelineArgs{
		Code:                  &CouponCode,
		CheckCouponRedemption: false,
		CalculateTotalUsages:  true,
	}
	if appliedAutomatically != nil {
		getMineCouponPipelineArgs.AppliedAutomatically = appliedAutomatically
	}
	getMineCouponPipelineResult := coupon_utils.GetMineCouponPipeline(ctx, getMineCouponPipelineArgs)
	coupon, err := checker.couponRepo.FindOne(ctx, getMineCouponPipelineResult)
	if err != nil {
		return nil, err
	}

	if coupon == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Coupon not found",
		}
	}

	checkCouponRedemptionErr := coupon_utils.CheckCouponRedemption(coupon_utils.CheckCouponRedemptionArgs{
		Coupon:    coupon,
		SubTotal:  subTotal,
		OrderType: orderType,
	})

	if checkCouponRedemptionErr != nil {
		return nil, checkCouponRedemptionErr
	}

	return coupon, nil
}
