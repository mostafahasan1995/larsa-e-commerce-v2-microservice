package cart_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConvertToCartModelArgs struct {
	Dto    []cart_shared.CartItemResponse
	CartId primitive.ObjectID
}

func ConvertToCartModel(args ConvertToCartModelArgs, additionalOptions types.AdditionalOptions) cart_models.Cart {

	totalQuantities := uint64(0)
	cartItemIds := []primitive.ObjectID{}

	for _, item := range args.Dto {
		totalQuantities += item.Quantity
		cartItemIds = append(cartItemIds, *item.Id)
	}
	cart := cart_models.Cart{
		BaseCart: cart_models.BaseCart{
			UserId:          additionalOptions.UserId,
			CartItemIds:     cartItemIds,
			TotalQuantities: totalQuantities,
			IsActive:        true,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}
	cart.MetaDataBaseModel.Id = args.CartId
	return cart
}

// ConvertCheckoutToCreateOrderDto converts a checkout DTO and calculated prices to an order DTO
func ConvertCheckoutToCreateOrderDto(
	createCheckoutDto *dtos.CreateCheckoutDto,
	calculateCartItemsPrices *cart_shared.CalculateCartItemsPricesResponse,
	salesSettings *settings_sales_models.SalesSettingsResponse,
	userId primitive.ObjectID,
	userType enums.UserType,
) (*order_dto.CreateOrderDTO, *helpers.ApiError) {

	order := &order_dto.CreateOrderDTO{}

	order.Type = createCheckoutDto.OrderType
	order.TotalDiscountValue = &calculateCartItemsPrices.TotalDiscountValue
	order.SubTotalPrice = calculateCartItemsPrices.SubTotalPrice
	order.TotalPrice = calculateCartItemsPrices.TotalPrice
	order.CurrencyCode = *createCheckoutDto.CurrencyCode
	order.CouponCode = createCheckoutDto.CouponCode
	order.Note = createCheckoutDto.Notes
	order.CouponDiscountValue = calculateCartItemsPrices.CouponValue
	order.LoyaltyDiscountValue = calculateCartItemsPrices.LoyaltyValue
	order.DeliveryFee = calculateCartItemsPrices.DeliveryFee
	order.CartId = createCheckoutDto.CartId

	order.PickupInfo = &order_dto.CreatePickupInfoDTO{
		PickupBranchId: createCheckoutDto.BranchId,
	}

	buyerLocation := order_dto.CreateShipmentDestinationDTO{}
	if createCheckoutDto.OrderType == order_enums.DELIVERY {
		if createCheckoutDto.Address == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Address is required for delivery",
			}
		}

		longitude := 0.0
		latitude := 0.0
		if createCheckoutDto.Address.Location != nil {
			longitude = createCheckoutDto.Address.Location.Coordinates[0]
			latitude = createCheckoutDto.Address.Location.Coordinates[1]
		}

		buyerLocation = order_dto.CreateShipmentDestinationDTO{
			Country:   createCheckoutDto.Address.Country,
			City:      createCheckoutDto.Address.City,
			Area:      createCheckoutDto.Address.Street,
			Address:   createCheckoutDto.Address.AddressLine1,
			Notes:     createCheckoutDto.Address.Notes,
			Longitude: &longitude,
			Latitude:  &latitude,
		}

		if salesSettings.OrderFlowMode == settings_sales_enums.ORDER_FLOW_MODE_B2B {
			order.DeliveryInfo = &order_dto.CreateDeliveryInfoDTO{
				DeliveryType: order_enums.B2B,
				B2BInfo: &order_dto.CreateB2BInfoDTO{
					ShipmentDetails: &order_dto.CreateB2BShipmentDetailsDTO{
						ShipmentDestination: order_dto.CreateShipmentDestinationDTO{
							Country:   createCheckoutDto.Address.Country,
							City:      createCheckoutDto.Address.City,
							Area:      createCheckoutDto.Address.Street,
							Address:   createCheckoutDto.Address.AddressLine1,
							Notes:     createCheckoutDto.Address.Notes,
							Longitude: &longitude,
							Latitude:  &latitude,
						},
					},
				},
			}
		} else {
			order.DeliveryInfo = &order_dto.CreateDeliveryInfoDTO{
				DeliveryType: order_enums.B2C,
				B2CInfo: &order_dto.CreateB2CInfoDTO{
					CustomerAddress: order_dto.CreateShipmentDestinationDTO{
						Country:   createCheckoutDto.Address.Country,
						City:      createCheckoutDto.Address.City,
						Area:      createCheckoutDto.Address.Street,
						Address:   createCheckoutDto.Address.AddressLine1,
						Notes:     createCheckoutDto.Address.Notes,
						Longitude: &longitude,
						Latitude:  &latitude,
					},
					DeliveryFee:  calculateCartItemsPrices.DeliveryFee,
					BranchId:     createCheckoutDto.BranchId,
					DeliveryTime: &calculateCartItemsPrices.EstimatedTime,
				},
			}
		}
	}

	order.PaymentInfo = &order_dto.CreatePaymentInfoDTO{
		Status:        order_enums.PaymentStatus(order_enums.UNPAID),
		PaymentMethod: createCheckoutDto.PaymentMethod,
	}

	var _userType order_enums.OrderParticipantType
	if userType == enums.Vendor {
		_userType = order_enums.VENDOR
	} else if userType == enums.Buyer {
		_userType = order_enums.BUYER
	}

	order.BuyerInfo = order_dto.CreateParticipantInfoDTO{
		Id:       userId,
		Type:     _userType,
		Location: buyerLocation,
	}

	// Handle PersonalDetails if provided
	if createCheckoutDto.PersonalDetailsDto != (dtos.PersonalDetailsDto{}) {
		// Build full name safely, checking for nil pointers
		var fullName string
		if createCheckoutDto.PersonalDetailsDto.FirstName != nil {
			fullName += *createCheckoutDto.PersonalDetailsDto.FirstName
		}
		if createCheckoutDto.PersonalDetailsDto.LastName != nil {
			if fullName != "" {
				fullName += " "
			}
			fullName += *createCheckoutDto.PersonalDetailsDto.LastName
		}
		order.BuyerInfo.FullName = fullName
		order.BuyerInfo.Email = createCheckoutDto.PersonalDetailsDto.Email
		order.BuyerInfo.CountryCode = createCheckoutDto.PersonalDetailsDto.CountryCode
		order.BuyerInfo.Phone = createCheckoutDto.PersonalDetailsDto.Mobile
		// Map alternative mobile if provided
		order.BuyerInfo.AlternativeMobile = createCheckoutDto.PersonalDetailsDto.AlternativeMobile
		order.BuyerInfo.AlternativeMobileCode = createCheckoutDto.PersonalDetailsDto.AlternativeMobileCode
	}

	order.OrderItems = calculateCartItemsPrices.OrderItems

	return order, nil
}
