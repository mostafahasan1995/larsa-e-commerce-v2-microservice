package cart_services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	_ "time/tzdata" // Embed timezone database for server deployment

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	cart_utils "larsa-e-commerce-microservice-v2/internal/modules/carts/utils"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_service "larsa-e-commerce-microservice-v2/internal/modules/order/service"
	payment_services "larsa-e-commerce-microservice-v2/internal/modules/payment/services"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	settings_sales_services "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckoutRepo interface {
	CalculateCartItemsPrices(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*cart_shared.CalculateCartItemsPricesResponse, *helpers.ApiError)
}

type CheckoutService interface {
	CalculateCartItemsPrices(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*cart_shared.CalculateCartItemsPricesResponse, *helpers.ApiError)
	ConvertCartToOrder(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*order_models.OrderResponse, *helpers.ApiError)
	CheckWorkingHours(ctx context.Context, branchId *primitive.ObjectID) *helpers.ApiError
}

type checkout struct {
	repo                   CheckoutRepo
	orderService           order_service.OrderService
	salesSettingsService   settings_sales_services.SalesSettingsService
	userCouponRepo         base_repos.BaseRepo[coupon_models.UserCoupon, coupon_models.UserCouponResponse]
	webCouponRepo          services.WebCouponRepo
	couponRuleCheckService services.CouponRuleCheckerService
	cartRepo               base_repos.BaseRepo[cart_models.Cart, cart_models.CartResponse]
	cartItemRepo           base_repos.BaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse]
	paymentService         payment_services.PaymentService
	branchRepo             base_repos.BaseRepo[product_models.Branch, product_models.BranchResponse]
}

func NewCheckoutService(i *do.Injector) (CheckoutService, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &checkout{
		repo:                   do.MustInvoke[CheckoutRepo](i),
		orderService:           do.MustInvoke[order_service.OrderService](i),
		salesSettingsService:   do.MustInvoke[settings_sales_services.SalesSettingsService](i),
		userCouponRepo:         *base_repos.NewBaseRepo[coupon_models.UserCoupon, coupon_models.UserCouponResponse](db, enums.USER_COUPON),
		webCouponRepo:          do.MustInvoke[services.WebCouponRepo](i),
		couponRuleCheckService: *do.MustInvoke[*services.CouponRuleCheckerService](i),
		cartRepo:               *base_repos.NewBaseRepo[cart_models.Cart, cart_models.CartResponse](db, enums.CART),
		cartItemRepo:           *base_repos.NewBaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse](db, enums.CART_ITEM),
		paymentService:         do.MustInvoke[payment_services.PaymentService](i),
		branchRepo:             *base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](db, enums.BRANCH),
	}, nil
}

func (checkout *checkout) CalculateCartItemsPrices(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*cart_shared.CalculateCartItemsPricesResponse, *helpers.ApiError) {
	return checkout.repo.CalculateCartItemsPrices(ctx, createCheckoutDto)
}

// CheckWorkingHours validates if the current time is within the branch's working hours
func (checkout *checkout) CheckWorkingHours(ctx context.Context, branchId *primitive.ObjectID) *helpers.ApiError {
	if branchId == nil {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Branch ID is required",
		}
	}

	// Get branch information
	branch, err := checkout.branchRepo.GetById(ctx, *branchId, []bson.M{})
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Failed to get branch: %v", err),
		}
	}

	if branch == nil {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Branch not found",
		}
	}

	// Check if branch is active
	if !branch.IsActive {
		return &helpers.ApiError{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Branch is not active",
		}
	}

	// Get current time in branch's timezone
	now := time.Now().UTC()
	branchTime := now

	// Convert to branch's timezone using location if available
	if branch.TimeZone.Location != "" {
		location, err := time.LoadLocation(branch.TimeZone.Location)
		if err != nil {
			log.Printf("[CHECKOUT] Failed to load timezone location %s: %v", branch.TimeZone.Location, err)
			// Fallback to offset-based conversion
			if branch.TimeZone.UtcOffset != 0 {
				offsetDuration := time.Duration(branch.TimeZone.UtcOffset) * time.Second
				branchTime = now.Add(offsetDuration)
			}
		} else {
			branchTime = now.In(location)
		}
	} else if branch.TimeZone.UtcOffset != 0 {
		// Fallback to offset-based conversion if no location specified
		offsetDuration := time.Duration(branch.TimeZone.UtcOffset) * time.Second
		branchTime = now.Add(offsetDuration)
	}

	log.Printf("[CHECKOUT] Checking working hours for branch %s at %s (UTC: %s)", branchId.Hex(), branchTime.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))

	// Check if branch has opening days configured
	if branch.OpeningDays == nil || len(branch.OpeningDays.Days) == 0 {
		log.Printf("[CHECKOUT] No opening days configured for branch %s - allowing checkout", branchId.Hex())
		return nil // Allow if no restrictions configured
	}

	// Get current day of week
	currentDay := branchTime.Weekday()

	// Check if current day has opening hours
	daySchedule, exists := branch.OpeningDays.Days[currentDay]
	if !exists || len(daySchedule.TimeRanges) == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusServiceUnavailable,
			Message:    fmt.Sprintf("Branch is closed on %s", currentDay.String()),
		}
	}

	// Convert current time to seconds since midnight
	currentTimeMilliSeconds := uint64(branchTime.Hour()*3600+branchTime.Minute()*60+branchTime.Second()) * 1000

	fmt.Println("currentTimeMilliSeconds", currentTimeMilliSeconds)

	// Check if current time falls within any opening range
	for _, timeRange := range daySchedule.TimeRanges {
		if currentTimeMilliSeconds >= timeRange.From && currentTimeMilliSeconds <= timeRange.To {
			log.Printf("[CHECKOUT] Branch %s is open - allowing checkout", branchId.Hex())
			return nil
		}
	}

	// If we get here, we're outside working hours
	nextOpening := checkout.findNextOpening(branch.OpeningDays, currentDay, currentTimeMilliSeconds)

	log.Printf("[CHECKOUT] Branch %s is closed - checkout blocked. Next opening: %s", branchId.Hex(), nextOpening)

	return &helpers.ApiError{
		StatusCode: http.StatusServiceUnavailable,
		Message:    fmt.Sprintf("Branch is currently closed. Next opening: %s", nextOpening),
	}
}

// findNextOpening finds the next opening time for the branch
func (checkout *checkout) findNextOpening(openingDays *product_models.OpeningDays, currentDay time.Weekday, currentTimeMilliSeconds uint64) string {
	// Check remaining ranges for today
	if dayRanges, exists := openingDays.Days[currentDay]; exists {
		for _, timeRange := range dayRanges.TimeRanges {
			if timeRange.From > currentTimeMilliSeconds {
				hours := timeRange.From / 1000 / 3600
				minutes := (timeRange.From / 1000 % 3600) / 60
				return fmt.Sprintf("Today at %02d:%02d", hours, minutes)
			}
		}
	}

	// Check next days
	for i := 1; i <= 7; i++ {
		nextDay := (currentDay + time.Weekday(i)) % 7
		if dayRanges, exists := openingDays.Days[nextDay]; exists && len(dayRanges.TimeRanges) > 0 {
			timeRange := dayRanges.TimeRanges[0]
			hours := timeRange.From / 1000 / 3600
			minutes := (timeRange.From / 1000 % 3600) / 60
			return fmt.Sprintf("%s at %02d:%02d", nextDay.String(), hours, minutes)
		}
	}

	return "No upcoming openings found"
}

// convert cart to order
func (checkout *checkout) ConvertCartToOrder(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	salesSettings, err := checkout.salesSettingsService.GetSalesSettings(ctx)
	if err != nil {
		return nil, err
	}

	calculateCartItemsPrices, err := checkout.CalculateCartItemsPrices(ctx, createCheckoutDto)
	if err != nil {
		return nil, err
	}

	// cart :=

	// Convert checkout DTO to order DTO using the extracted utility function
	order, err := cart_utils.ConvertCheckoutToCreateOrderDto(
		createCheckoutDto,
		calculateCartItemsPrices,
		salesSettings,
		cfg.User.Id,
		cfg.User.Type,
	)
	if err != nil {
		return nil, err
	}

	session, _ := checkout.cartRepo.GetDB().StartSession()
	defer session.EndSession(ctx)
	var createdOrder *order_models.OrderResponse
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if createCheckoutDto.Coupon != nil {
			err := checkout.webCouponRepo.ConvertToUserCoupon(sessCtx, createCheckoutDto.Coupon)

			if err != nil {
				return nil, err
			}
		}

		cartItemFilter := bson.M{
			"cartId":   createCheckoutDto.CartId,
			"branchId": bson.M{"$ne": createCheckoutDto.BranchId},
			"isActive": true,
		}

		fmt.Println("ConvertCartToOrder cartItemFilter", cartItemFilter)
		toBeMovedCartItems, err := checkout.cartItemRepo.GetAll(sessCtx, []bson.M{
			{
				"$match": cartItemFilter,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(toBeMovedCartItems) > 0 {
			newCart := cart_utils.ConvertToCartModel(cart_utils.ConvertToCartModelArgs{
				Dto:    toBeMovedCartItems,
				CartId: primitive.NewObjectID(),
			}, types.AdditionalOptions{
				UserId: cfg.User.Id,
				Mode:   types.CREATE,
			})
			createdCart, createNewCartErr := checkout.cartRepo.Create(sessCtx, newCart)
			if createNewCartErr != nil {
				return nil, createNewCartErr
			}
			checkout.cartItemRepo.UpdateMany(sessCtx, cartItemFilter, map[string]interface{}{
				"cartId": createdCart.Id,
			}, cfg.User.Id)
		}

		if createCheckoutDto.PaymentMethod == order_enums.ONLINE_PAYMENT {
			//clone the cart
			if err := checkout.CloneCart(sessCtx, createCheckoutDto.CartId, createCheckoutDto.BranchId); err != nil {
				return nil, err
			}
		}

		checkout.cartRepo.UpdateFieldsById(sessCtx, *createCheckoutDto.CartId, map[string]interface{}{
			"isActive": false,
		}, cfg.User.Id)
		checkout.cartItemRepo.UpdateMany(sessCtx, bson.M{"cartId": createCheckoutDto.CartId, "branchId": createCheckoutDto.BranchId}, map[string]interface{}{
			"isActive": false,
		}, cfg.User.Id)

		/*
		 * For some reason passing &sessCtx won't work even though it's a pointer.
		 * It works if the other side is a not pointer which is not possible in this case.
		 * So we need to create a new context and pass session context to it.
		 */
		// var pointerSessionCtx context.Context
		// pointerSessionCtx = sessCtx
		// checkout.couponRuleCheckService.GenerateUserCouponBasedOnCouponRules(&pointerSessionCtx, cfg.Db)

		err = checkout.CreateInvoiceForOrder(sessCtx, createdOrder)
		if err != nil {
			return nil, err
		}

		createdOrder, err = checkout.orderService.Create(sessCtx, order)
		if err != nil {
			return nil, err
		}

		return createdOrder, nil

	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to checkout",
		}
	}

	return createdOrder, nil
}

func (checkout *checkout) CreateInvoiceForOrder(ctx context.Context, order *order_models.OrderResponse) *helpers.ApiError {
	//TODO: create invoice for the order
	return nil
}

func (checkout *checkout) CloneCart(ctx context.Context, cartId, branchId *primitive.ObjectID) error {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": *cartId,
			},
		},
	}

	var carts []cart_models.Cart
	if err := checkout.cartRepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &carts)
	}); err != nil {
		return err
	}

	if len(carts) == 0 {
		return &helpers.ApiError{
			StatusCode: 404,
			Message:    "Cart not found",
		}
	}

	itemsPipeline := []bson.M{
		{
			"$match": bson.M{
				"cartId":   *cartId,
				"branchId": *branchId,
			},
		},
	}

	var items []cart_shared.CartItem
	if err := checkout.cartItemRepo.Aggregate(ctx, itemsPipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &items)
	}); err != nil {
		return err
	}

	if len(items) == 0 {
		return &helpers.ApiError{
			StatusCode: 404,
			Message:    "Cart items not found",
		}
	}

	newCart := carts[0]
	newCart.IsClone = true
	newCart.RefId = *cartId
	newCart.Id = primitive.NewObjectID()

	writeOps := []mongo.WriteModel{}
	newCartItemIds := []primitive.ObjectID{}

	for _, item := range items {
		id := primitive.NewObjectID()
		newItem := item
		newItem.Id = id
		newItem.CartId = newCart.Id
		insertOp := mongo.NewInsertOneModel().SetDocument(newItem)
		writeOps = append(writeOps, insertOp)
		newCartItemIds = append(newCartItemIds, id)
	}

	newCart.CartItemIds = newCartItemIds

	if _, err := checkout.cartRepo.BulkWrite(ctx, []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(newCart),
	}); err != nil {
		return err
	}

	if _, err := checkout.cartItemRepo.BulkWrite(ctx, writeOps); err != nil {
		return err
	}

	return nil
}
