package cart_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	address_filters "larsa-e-commerce-microservice-v2/internal/modules/address/filters"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	coupon_services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCheckoutDto struct {
	CartId       *primitive.ObjectID `json:"-"`
	BranchId     *primitive.ObjectID `json:"branchId" validate:"required"`
	CurrencyCode *string             `json:"currencyCode" validate:"omitempty"`

	AddressId *primitive.ObjectID `json:"addressId" validate:"omitempty"`
	/*
	* NOTE: ZoneId is required, but We are using optional
	* because this dto is being used in multiple places and in some places we don't want
	* to pass them
	 */
	ZoneId     *primitive.ObjectID   `json:"zoneId" validate:"omitempty"`
	CouponCode *string               `json:"couponCode" validate:"omitempty"`
	OrderType  order_enums.OrderType `json:"orderType" validate:"required,oneof=delivery pickup"`
	Notes      *string               `json:"notes" validate:"omitempty"`

	PersonalDetailsDto PersonalDetailsDto        `json:"personalDetailsDto"`
	PaymentMethod      order_enums.PaymentMethod `json:"paymentMethod" validate:"required,oneof=cash_terminal cash_on_delivery online_payment"`

	Coupon        *coupon_models.CouponResponse   `json:"-"`
	Branch        *product_models.BranchResponse  `json:"-"`
	Address       *address_models.AddressResponse `json:"-"`
	CalculateList *bool                           `json:"-" validate:"omitempty"`
}

type PersonalDetailsDto struct {
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Email       *string `json:"email" validate:"omitempty,email"`
	CountryCode *string `json:"countryCode"`
	Mobile      *string `json:"mobile"`
	//Alternative Mobile is optional
	AlternativeMobile     *string `json:"alternativeMobile" validate:"omitempty"`
	AlternativeMobileCode *string `json:"alternativeMobileCode" validate:"omitempty"`
}

func (dto *CreateCheckoutDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {

	if dto.CurrencyCode == nil {
		dto.CurrencyCode = &configuration.EnvironmentVariables.DefaultCurrencyCode
	}

	addressRepo := base_repos.NewBaseRepo[address_models.AddressModel, address_models.AddressResponse](args.DB, enums.ADDRESS)
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	cartRepo := base_repos.NewBaseRepo[cart_models.Cart, cart_models.CartResponse](args.DB, enums.CART)
	currencyRepo := base_repos.NewBaseRepo[config_models.Currency, config_models.CurrencyResponse](args.DB, enums.CURRENCY)
	zoneRepo := base_repos.NewBaseRepo[config_models.Zone, config_models.ZoneResponse](args.DB, enums.ZONE)
	// zoneService := do.MustInvoke[config_services.ZoneService](args.Injector)
	couponRuleCheckerService, _ := coupon_services.NewCouponRuleChecker(args.Injector)
	if args.User == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnauthorized,
			Message:    "User context is required",
		}
	}
	cart, err := cartRepo.FindOne(args.Context, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"userId":   args.User.Id,
				"isActive": true,
			}),
		},
	})
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Cart not found",
		}
	}
	dto.CartId = cart.Id
	if dto.OrderType == order_enums.DELIVERY && dto.AddressId == nil && cart.AddressId != nil {
		dto.AddressId = cart.AddressId
	}

	if dto.OrderType == order_enums.DELIVERY && dto.AddressId == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "addressId is required for delivery (include it in the request body or set it when getting the cart with ?addressId=...)",
		}
	}

	if dto.BranchId != nil {
		branchFilter := product_filters.BranchFilters{}
		branchPipeline := branchFilter.BuildPipeline(bson.M{}, args.Context, false)
		branch, err := branchRepo.GetById(args.Context, *dto.BranchId, branchPipeline)
		if err != nil {
			return nil, err
		}

		dto.Branch = branch
	}

	if dto.AddressId != nil {

		if dto.OrderType == order_enums.DELIVERY {
			addressFilter := address_filters.AddressFilters{
				UserId: &args.User.Id,
			}
			addressPipeline := addressFilter.BuildPipeline(bson.M{}, args.Context, false)
			address, err := addressRepo.GetById(args.Context, *dto.AddressId, addressPipeline)
			if err != nil {
				return nil, err
			}

			dto.Address = address
		} else {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "AddressId can't be provided for pickup order type",
			}
		}

	}

	// if dto.ZoneId != nil && dto.Address != nil && dto.Address.Location != nil {
	// 	fmt.Printf("dto.Address: %v\n", dto.Address.Location.Coordinates)
	// 	zones, err := zoneService.FindZonesContainingPoint(args.Context, dto.Address.Location.Coordinates[1], dto.Address.Location.Coordinates[0])
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	fmt.Printf("zones: %v\n", zones)
	// 	if len(zones) == 0 {
	// 		return nil, &helpers.ApiError{
	// 			StatusCode: http.StatusNotFound,
	// 			Message:    "No zone found for the provided address",
	// 		}
	// 	}

	// 	zone := util.Find(zones, func(zone config_models.ZoneResponse) bool {
	// 		return zone.Id == *dto.ZoneId
	// 	})

	// 	if zone == nil {
	// 		return nil, &helpers.ApiError{
	// 			StatusCode: http.StatusNotFound,
	// 			Message:    "Zone not found for the provided address",
	// 		}
	// 	}
	// }
	if dto.OrderType == order_enums.DELIVERY && dto.Branch != nil {
		mode := dto.Branch.DeliveryPricingMode
		if mode == "" {
			mode = string(product_enums.DeliveryPricingModeZone)
		}

		switch mode {
		case string(product_enums.DeliveryPricingModeZone):
			// For zone-based branches, ZoneId must be provided and valid.
			if dto.ZoneId == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "zoneId is required for delivery when branch uses zone-based pricing",
				}
			}

			zoneFilter := config_filters.ZoneFilters{}
			zonePipeline := zoneFilter.BuildPipeline(bson.M{}, args.Context, false)
			if _, err = zoneRepo.GetById(args.Context, *dto.ZoneId, zonePipeline); err != nil {
				return nil, err
			}

		case string(product_enums.DeliveryPricingModeDistance):
			// Distance mode: ZoneId not required in DTO; use branch's first zone for estimatedTime/minOrderAmount (same as add-to-cart).
			if dto.ZoneId == nil && dto.Branch.Zones != nil && len(*dto.Branch.Zones) > 0 {
				firstZoneId := (*dto.Branch.Zones)[0].ZoneId
				dto.ZoneId = &firstZoneId
			}
			if dto.AddressId == nil || dto.Address == nil || dto.Address.Location == nil || len(dto.Address.Location.Coordinates) < 2 {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "addressId with valid location is required for delivery when branch uses distance-based pricing",
				}
			}
			// If client sent zoneId, validate it exists; if we set from first zone, skip (already from branch)
			if dto.ZoneId != nil {
				zoneFilter := config_filters.ZoneFilters{}
				zonePipeline := zoneFilter.BuildPipeline(bson.M{}, args.Context, false)
				if _, err = zoneRepo.GetById(args.Context, *dto.ZoneId, zonePipeline); err != nil {
					return nil, err
				}
			}
		}
	}

	if dto.CurrencyCode != nil {
		currency, err := currencyRepo.FindOne(args.Context, []bson.M{
			{
				"$match": base_utils.AddSoftDeletedFilter(bson.M{
					"code": *dto.CurrencyCode,
				}),
			},
		})
		if err != nil {
			return nil, err
		}

		if currency == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Currency not found",
			}
		}
	}

	if dto.CouponCode != nil {
		// Note: SubTotal is passed as nil here since it's calculated later in CalculateCartItemsPrices
		// The minimum subtotal validation can be added there if needed
		coupon, err := couponRuleCheckerService.CheckCouponRedemption(args.Context, *dto.CouponCode, nil, nil, &dto.OrderType)

		if err != nil {
			return nil, err
		}

		dto.Coupon = coupon
	}

	if dto.PersonalDetailsDto.Mobile != nil {
		if !util.IsValidPhoneNumber(*dto.PersonalDetailsDto.Mobile) {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid phone number",
			}
		}
	}

	return nil, nil
}
