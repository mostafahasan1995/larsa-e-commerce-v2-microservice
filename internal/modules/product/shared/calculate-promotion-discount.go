package product_shared

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	promotion_services "larsa-e-commerce-microservice-v2/internal/modules/promotion/services"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculatePromotionDiscountArgs struct {
	ProductId    primitive.ObjectID
	BranchId     *primitive.ObjectID
	Branch       *product_models.BranchResponse
	TotalPrice   float64
	CurrencyCode string
}

type CalculatePromotionDiscountResult struct {
	ActivePromotion          *product_models.ProductPromotionInfo
	TotalPriceAfterPromotion float64
	PromotionDiscountValue   float64
}

// getBranchTimezone converts branch timezone information to time.Location
func getBranchTimezone(branch *product_models.BranchResponse) *time.Location {
	if branch == nil || branch.TimeZone.Location == "" {
		return time.UTC
	}

	// Try to load the timezone location
	loc, err := time.LoadLocation(branch.TimeZone.Location)
	if err != nil {
		// If location string is invalid, fall back to UTC
		return time.UTC
	}

	return loc
}

func CalculatePromotionDiscount(ctx context.Context, args CalculatePromotionDiscountArgs) (*CalculatePromotionDiscountResult, *helpers.ApiError) {
	// Get promotion service
	promotionService := do.MustInvoke[promotion_services.PromotionService](util.GetInjector(ctx))

	// Get timezone from branch if available, otherwise use UTC
	timezone := getBranchTimezone(args.Branch)

	// Get the single best active promotion for the product and branch
	activePromotion, err := promotionService.GetActivePromotionsForProductAndBranch(ctx, args.ProductId, args.BranchId, timezone)
	if err != nil {
		return nil, err
	}

	// If no active promotions, return original price
	if activePromotion.Id == primitive.NilObjectID {
		return &CalculatePromotionDiscountResult{
			ActivePromotion:          nil,
			TotalPriceAfterPromotion: args.TotalPrice,
			PromotionDiscountValue:   0,
		}, nil
	}

	// Get the single best promotion (highest discount percentage)
	bestPromotion := activePromotion

	// Calculate promotion discount
	promotionDiscountValue := args.TotalPrice * bestPromotion.DiscountPercentage / 100
	totalPriceAfterPromotion := args.TotalPrice - promotionDiscountValue

	// Create promotion info
	promotionInfo := &product_models.ProductPromotionInfo{
		PromotionId:            bestPromotion.Id,
		Name:                   bestPromotion.Name,
		Description:            bestPromotion.Description,
		DiscountPercentage:     bestPromotion.DiscountPercentage,
		StartDate:              bestPromotion.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		EndDate:                bestPromotion.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		PromotionDiscountValue: promotionDiscountValue,
	}

	return &CalculatePromotionDiscountResult{
		ActivePromotion:          promotionInfo,
		TotalPriceAfterPromotion: totalPriceAfterPromotion,
		PromotionDiscountValue:   promotionDiscountValue,
	}, nil
}
