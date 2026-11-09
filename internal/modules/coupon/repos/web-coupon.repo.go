package coupon_repos

import (
	"context"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	coupon_utils "larsa-e-commerce-microservice-v2/internal/modules/coupon/utils"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
)

type WebCouponRepo struct {
	*base_repos.BaseRepo[models.Coupon, models.CouponResponse]
	userCouponRepo base_repos.BaseRepo[models.UserCoupon, models.UserCouponResponse]
}

func NewWebCouponRepo(i *do.Injector) (services.WebCouponRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebCouponRepo{
		BaseRepo:       base_repos.NewBaseRepo[models.Coupon, models.CouponResponse](db, enums.COUPON),
		userCouponRepo: *base_repos.NewBaseRepo[models.UserCoupon, models.UserCouponResponse](db, enums.USER_COUPON),
	}, nil
}

func (repo *WebCouponRepo) GetMine(ctx context.Context) (*[]models.CouponResponse, *helpers.ApiError) {

	pipeline := coupon_utils.GetMineCouponPipeline(ctx, coupon_utils.GetMineCouponPipelineArgs{
		Code:                  nil,
		CalculateTotalUsages:  true,
		CheckCouponRedemption: true,
		AppliedAutomatically:  util.BoolPtr(true),
	})
	coupons, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return &coupons, nil

}

func (repo *WebCouponRepo) ConvertToUserCoupon(ctx context.Context, coupon *models.CouponResponse) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	usages := []coupon_models.UserCouponUsage{{
		Timestamp: time.Now(),
	}}

	// Only try to access UserCoupon if it exists
	if coupon.UserCoupon != nil {
		if coupon.UserCoupon.Usages != nil {
			usages = append(usages, *coupon.UserCoupon.Usages...)
		}
	}

	newUserCoupon := coupon_models.UserCoupon{
		BaseUserCoupon: coupon_models.BaseUserCoupon{
			UserId:      cfg.User.Id,
			CouponCode:  coupon.Code,
			Usages:      &usages,
			TotalUsages: 1, // Default to 1 for new users
		},
	}

	// If UserCoupon exists, update the total usages
	if coupon.UserCoupon != nil {
		newUserCoupon.TotalUsages = uint64(coupon.UserCoupon.TotalUsages) + 1
		newUserCoupon.MetaDataBaseModel = util.CreateMetaDataBaseModel(types.AdditionalOptions{
			UserId: cfg.User.Id,
			Mode:   types.UPDATE,
		})
	} else {
		newUserCoupon.MetaDataBaseModel = util.CreateMetaDataBaseModel(types.AdditionalOptions{
			UserId: cfg.User.Id,
			Mode:   types.CREATE,
		})
	}

	err := repo.userCouponRepo.Upsert(ctx, base_repos.UpsertArgs[coupon_models.UserCoupon, coupon_models.UserCouponResponse]{
		Filter: []bson.M{
			{"$match": bson.M{
				"userId":     cfg.User.Id,
				"couponCode": coupon.Code,
			}},
		},
		Create: []coupon_models.UserCoupon{newUserCoupon},
		Update: []coupon_models.UserCoupon{newUserCoupon},
	})
	if err != nil {
		return err
	}

	return nil
}
