package cart_repos

import (
	"context"
	"fmt"
	cart_enums "larsa-e-commerce-microservice-v2/internal/modules/carts/enums"
	cart_filters "larsa-e-commerce-microservice-v2/internal/modules/carts/filters"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	loyalty_services "larsa-e-commerce-microservice-v2/internal/modules/loyalty/services"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_shared "larsa-e-commerce-microservice-v2/internal/modules/order/shared"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/filters"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
)

type CheckoutRepo struct {
	branchProductRepo      *base_repos.BaseRepo[product_models.BranchProduct, product_models.BranchProductResponse]
	cartItemRepo           *base_repos.BaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse]
	cartRepo               *base_repos.BaseRepo[cart_models.Cart, cart_models.CartResponse]
	orderRepo              *base_repos.BaseRepo[order_models.Order, order_models.OrderResponse]
	productRepo            *base_repos.BaseRepo[product_models.Product, product_models.ProductResponse]
	loyaltyDiscountService loyalty_services.LoyaltyDiscountService
}

func NewCheckoutRepo(i *do.Injector) (services.CheckoutRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CheckoutRepo{
		branchProductRepo:      base_repos.NewBaseRepo[product_models.BranchProduct, product_models.BranchProductResponse](db, enums.BRANCH_PRODUCT),
		cartItemRepo:           base_repos.NewBaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse](db, enums.CART_ITEM),
		cartRepo:               base_repos.NewBaseRepo[cart_models.Cart, cart_models.CartResponse](db, enums.CART),
		orderRepo:              base_repos.NewBaseRepo[order_models.Order, order_models.OrderResponse](db, enums.ORDER),
		productRepo:            base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](db, enums.PRODUCT),
		loyaltyDiscountService: do.MustInvoke[loyalty_services.LoyaltyDiscountService](i),
	}, nil
}

func (checkout *CheckoutRepo) CalculateCartItemsPrices(ctx context.Context, createCheckoutDto *dtos.CreateCheckoutDto) (*cart_shared.CalculateCartItemsPricesResponse, *helpers.ApiError) {
	fmt.Println("createCheckoutDto", createCheckoutDto)

	checkoutFilter := cart_filters.CartItemFilters{
		CartId:   createCheckoutDto.CartId,
		BranchId: createCheckoutDto.BranchId,
		BaseFilters: filters.BaseFilters{
			PopFields: &[]string{"itemId"},
		},
	}

	fmt.Println("checkoutFilter", checkoutFilter)

	pipeline := checkoutFilter.BuildPipeline(bson.M{}, ctx, true)

	cartItems, err := checkout.cartItemRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Cart is empty",
		}
	}

	var productIds []primitive.ObjectID
	var productCartItems []cart_shared.CartItemResponse

	/*
	 * NOTE: Right now we're only need to get the product ids for the products in the cart
	 * We'll need to add the subscription ids for the subscriptions in the cart
	 * when the time comes.
	 * Timestamp: 2025-05-12 15:42:00
	 */
	for _, cartItem := range cartItems {
		if cartItem.Type == cart_enums.CartItemTypeProduct {
			productIds = append(productIds, cartItem.ItemId)
			productCartItems = append(productCartItems, cartItem)
		}
	}

	/*
	 * NOTE: Right now we're only need to get the product ids for the products in the cart
	 * We'll need to add the subscription ids for the subscriptions in the cart
	 * when the time comes.
	 * We'll get the product even if it's not available in the branch or soft deleted.
	 * because we've checked the availability in the add to cart step.
	 * And we've the product availability (soft deleted) in the cart item adding step.
	 * Timestamp: 2025-05-27 15:42:00
	 */
	storedProducts, err := checkout.productRepo.GetAll(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{"$in": productIds},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	productDiscountValue := 0.0
	couponDiscountValue := 0.0
	loyaltyDiscountValue := 0.0

	totalDiscountValue := 0.0
	totalPrice := 0.0
	subTotalPrice := 0.0

	deliveryFee := 0.0
	estimatedTime := uint64(0)
	totalQuantities := uint64(0)
	orderItems := []order_dto.CreateOrderItemDto{}

	for _, storedProduct := range storedProducts {
		productCartItems = util.Map(productCartItems, func(cartItemResponse cart_shared.CartItemResponse) cart_shared.CartItemResponse {

			if cartItemResponse.ItemId == storedProduct.Id {
				cartItemResponse.Product = &storedProduct
			}

			return cartItemResponse
		})
	}

	for _, productCartItem := range productCartItems {

		if productCartItem.Product == nil {
			continue
		}

		calculateProductPriceResult, err := product_shared.CalculateProductPrice(ctx, &product_shared.CalculateCartItemsPricesArgs{
			Product:      *productCartItem.Product,
			CurrencyCode: *createCheckoutDto.CurrencyCode,
			Quantity:     productCartItem.Quantity,
			BranchId:     createCheckoutDto.BranchId,
			Branch:       createCheckoutDto.Branch,
			Addons:       productCartItem.Addons,
			UnitId:       productCartItem.UnitId,
		})

		if err != nil {
			return nil, err
		}

		addons := []order_models.OrderAddon{}
		if calculateProductPriceResult.Addons != nil {
			for _, addon := range *calculateProductPriceResult.Addons {
				values := []order_models.OrderAddonValue{}
				for _, value := range addon.AddonValues {
					addonValue := util.Find(calculateProductPriceResult.AddonValues, func(addonValue product_models.ProductAddonValue) bool {
						return addonValue.Id == value.Id
					})
					if addonValue != nil {
						values = append(values, order_models.OrderAddonValue{
							Label:    value.Label,
							Quantity: addonValue.Quantity,
						})
					}
				}
				addons = append(addons, order_models.OrderAddon{
					Label:  addon.Label,
					Values: values,
				})
			}
		}

		var product = product_models.ProductResponse{
			BaseProduct: product_models.BaseProduct{
				Name:             productCartItem.Product.Name,
				Description:      productCartItem.Product.Description,
				Images:           productCartItem.Product.Images,
				ProductionLineId: productCartItem.Product.ProductionLineId,
				Prices:           productCartItem.Product.Prices,
				Flag:             productCartItem.Product.Flag,
				VendorId:         productCartItem.Product.VendorId,
				CategoryIds:      productCartItem.Product.CategoryIds,
			},
		}

		orderItem := order_dto.CreateOrderItemDto{
			Quantity:        productCartItem.Quantity,
			ItemId:          productCartItem.ItemId,
			Name:            productCartItem.Product.Name,
			Product:         &product,
			BranchId:        productCartItem.BranchId,
			Type:            productCartItem.Type,
			SubTotalPrice:   calculateProductPriceResult.SubTotalPrice,
			TotalPrice:      calculateProductPriceResult.TotalPrice,
			DiscountValue:   calculateProductPriceResult.DiscountValue,
			Discount:        calculateProductPriceResult.Discount,
			Price:           calculateProductPriceResult.PriceItem,
			Addons:          addons,
			AddonTotalPrice: calculateProductPriceResult.AddonTotalPrice,
		}
		if productCartItem.AdditionalInformation != nil {
			orderItem.AdditionalInformation = &order_shared.AdditionalInformationDto{
				Image: productCartItem.AdditionalInformation.Image,
				Logo:  productCartItem.AdditionalInformation.Logo,
				Note:  productCartItem.AdditionalInformation.Note,
			}
		}
		orderItems = append(orderItems, orderItem)

		subTotalPrice = (subTotalPrice + calculateProductPriceResult.SubTotalPrice)
		productDiscountValue = (productDiscountValue + calculateProductPriceResult.DiscountValue)
		totalQuantities = totalQuantities + productCartItem.Quantity
	}

	// Delivery fee and estimated time: check deliveryPricingMode first, then zone vs distance.
	
	if createCheckoutDto.Branch != nil && createCheckoutDto.OrderType == order_enums.DELIVERY {
		mode := createCheckoutDto.Branch.DeliveryPricingMode
		if mode == "" {
			mode = string(product_enums.DeliveryPricingModeZone)
		}
		if mode == string(product_enums.DeliveryPricingModeDistance) {
			// Distance mode: delivery fee from distance; estimatedTime from zone (ZoneId from DTO or first zone)
			fee, _, feeErr := product_shared.CalculateDeliveryFee(product_shared.DeliveryFeeArgs{
				Branch:       createCheckoutDto.Branch,
				ZoneId:       createCheckoutDto.ZoneId,
				Address:      createCheckoutDto.Address,
				OrderType:    createCheckoutDto.OrderType,
				CurrencyCode: createCheckoutDto.CurrencyCode,
			})
			if feeErr != nil {
				return nil, feeErr
			}
			deliveryFee = fee
			// EstimatedTime and branchZone for min-order check: use ZoneId from DTO, or first zone (same as add-to-cart)
			if createCheckoutDto.Branch.Zones != nil && len(*createCheckoutDto.Branch.Zones) > 0 {
				var branchZone *product_models.BranchZone
				if createCheckoutDto.ZoneId != nil {
					branchZone = util.Find(*createCheckoutDto.Branch.Zones, func(zone product_models.BranchZone) bool {
						return zone.ZoneId == *createCheckoutDto.ZoneId
					})
				}
				if branchZone == nil {
					branchZone = &(*createCheckoutDto.Branch.Zones)[0]
				}
				estimatedTime = branchZone.EstimatedTime
			}
			} else if createCheckoutDto.Branch.Zones != nil && createCheckoutDto.ZoneId != nil {
				// Zone mode: delivery fee and estimated time from zone
				branchZone := util.Find(*createCheckoutDto.Branch.Zones, func(zone product_models.BranchZone) bool {
					return zone.ZoneId == *createCheckoutDto.ZoneId
				})
				if branchZone != nil {
					if branchZone.DeliveryFee.CurrencyCode == *createCheckoutDto.CurrencyCode {
						deliveryFee = branchZone.DeliveryFee.Value
					}
					estimatedTime = branchZone.EstimatedTime
				}
			}
		}
	totalDiscountValue = totalDiscountValue + productDiscountValue
	totalPrice = subTotalPrice + deliveryFee - totalDiscountValue

	if createCheckoutDto.Coupon != nil {
		// Validate minimum subtotal if coupon has this requirement
		if createCheckoutDto.Coupon.Usage != nil && createCheckoutDto.Coupon.Usage.MinimumSubtotal != nil {
			if subTotalPrice < *createCheckoutDto.Coupon.Usage.MinimumSubtotal {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Order subtotal does not meet the minimum requirement for this coupon",
					Details: &[]types.ValidationError{
						{
							Field:   "coupon_code",
							Message: fmt.Sprintf("Minimum subtotal required: %.2f, but current subtotal is: %.2f", *createCheckoutDto.Coupon.Usage.MinimumSubtotal, subTotalPrice),
							Value:   createCheckoutDto.Coupon.Code,
						},
					},
				}
			}
		}

		couponResult := cart_shared.CalculateCoupon(ctx, cart_shared.CalculateCouponArgs{
			CurrencyCode: *createCheckoutDto.CurrencyCode,
			Coupon:       *createCheckoutDto.Coupon,
			TotalPrice:   totalPrice,
		})

		couponDiscountValue = couponResult.CouponDiscountValue
		totalDiscountValue = totalDiscountValue + couponDiscountValue
		totalPrice = totalPrice - couponDiscountValue
	}

	// Get user ID from context for loyalty discount calculation
	cfg, _ := util.GetReqAppCfg(ctx)
	userId := cfg.User.Id

	// Apply loyalty discount after coupon discount
	loyaltyResult := cart_shared.CalculateLoyaltyDiscount(ctx, cart_shared.CalculateLoyaltyDiscountArgs{
		TotalPrice:   totalPrice,
		CurrencyCode: *createCheckoutDto.CurrencyCode,
		UserId:       userId,
	}, checkout.loyaltyDiscountService)

	loyaltyDiscountValue = loyaltyResult.LoyaltyDiscountValue
	totalDiscountValue = totalDiscountValue + loyaltyDiscountValue
	totalPrice = totalPrice - loyaltyDiscountValue

	return &cart_shared.CalculateCartItemsPricesResponse{
		SubTotalPrice:             subTotalPrice,
		TotalPrice:                totalPrice,
		DeliveryFee:               deliveryFee,
		EstimatedTime:             estimatedTime,
		DiscountValue:             productDiscountValue,
		CouponValue:               couponDiscountValue,
		LoyaltyValue:              loyaltyDiscountValue,
		LoyaltyDiscountPercentage: loyaltyResult.DiscountPercentage,
		LoyaltyOrdersRemaining:    loyaltyResult.OrdersRemaining,
		TotalDiscountValue:        totalDiscountValue,
		CartId:                    *createCheckoutDto.CartId,
		BranchId:                  createCheckoutDto.BranchId,
		ZoneId:                    createCheckoutDto.ZoneId,
		CurrencyCode:              createCheckoutDto.CurrencyCode,
		CouponCode:                createCheckoutDto.CouponCode,
		AddressId:                 createCheckoutDto.AddressId,
		OrderType:                 createCheckoutDto.OrderType,
		OrderItems:                orderItems,
		Note:                      createCheckoutDto.Notes,
		// CartItems:                 ,
	}, nil
}
