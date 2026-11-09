package quotation_service

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_services "larsa-e-commerce-microservice-v2/internal/modules/profile/services"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationOrderConvertorService interface {
	ConvertQuotationToOrder(ctx context.Context, quotation *quotation_models.QuotationResponse) (*order_models.Order, *helpers.ApiError)
}

type quotationOrderConvertor struct {
	profileService profile_services.ProfileService
}

func NewQuotationOrderConvertorService(i *do.Injector) (QuotationOrderConvertorService, error) {
	return &quotationOrderConvertor{
		profileService: do.MustInvoke[profile_services.ProfileService](i),
	}, nil
}

func (s *quotationOrderConvertor) ConvertQuotationToOrder(ctx context.Context, quotation *quotation_models.QuotationResponse) (*order_models.Order, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	metaData := models.MetaDataBaseModel{
		CreatedBy: &cfg.User.Id,
		UpdatedBy: &cfg.User.Id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Id:        primitive.NewObjectID(),
	}

	calculatePricesResult := s.calculatePrices(quotation)
	subTotalPrice := calculatePricesResult.SubTotalPrice
	currencyCode := calculatePricesResult.CurrencyCode
	totalPrice := calculatePricesResult.TotalPrice
	totalDiscountValue := calculatePricesResult.DiscountValue
	deliveryFee := calculatePricesResult.DeliveryFee

	quotationModel := order_models.Order{
		BaseOrder: order_models.BaseOrder{
			MetaDataBaseModel: metaData,
			QuotationId:       &quotation.Id,

			Type:         order_enums.DELIVERY,
			DeliveryInfo: s.mapDeliveryInfo(quotation, deliveryFee),
			PaymentInfo: &order_models.PaymentInfo{
				Status:   order_enums.UNPAID,
				Payments: []order_models.Payments{},
			},

			Code: order_utils.GenerateOrderCode(),

			Status: string(settings_sales_enums.SYSTEM_STATUS_TYPE_PENDING),

			SubTotalPrice: subTotalPrice,
			CurrencyCode:  currencyCode,

			TotalPrice:         totalPrice,
			TotalDiscountValue: &totalDiscountValue,

			BuyerInfo: s.mapBuyerInfo(quotation),

			VendorInfo: s.mapVendorInfo(ctx, quotation),

			OrderItems: s.mapOrderItems(quotation),

			Documents:  s.mapDocuments(quotation),
			Comments:   s.mapComments(quotation),
			LogHistory: s.mapHistory(quotation),
		},
	}

	return &quotationModel, nil
}

type CalculatePricesResult struct {
	SubTotalPrice float64
	TotalPrice    float64
	DiscountValue float64
	DeliveryFee   float64
	CurrencyCode  string
}

func (s *quotationOrderConvertor) calculatePrices(quotation *quotation_models.QuotationResponse) CalculatePricesResult {
	subTotalPrice := 0.0
	totalPrice := 0.0
	totalDiscountValue := 0.0

	for _, product := range quotation.Products {
		//TODO: We need to handle different currencies of products
		subTotalPrice += product.Price.Amount * float64(product.Quantity)
		totalDiscountValue += product.DiscountPercent * float64(product.Quantity)
	}

	if quotation.DiscountAmount != nil {
		totalDiscountValue += quotation.DiscountAmount.Amount
	}

	totalPrice = subTotalPrice - totalDiscountValue

	deliveryFee := 0.0
	if quotation.Fee.DeliveryFee.DeliveryToBuyerFee != nil {
		deliveryFee += quotation.Fee.DeliveryFee.DeliveryToBuyerFee.Amount
	}
	if quotation.Fee.DeliveryFee.CustomsClearanceFee != nil {
		deliveryFee += quotation.Fee.DeliveryFee.CustomsClearanceFee.Amount
	}
	if quotation.Fee.DeliveryFee.ImportDeclarationFees != nil {
		deliveryFee += quotation.Fee.DeliveryFee.ImportDeclarationFees.Amount
	}
	if quotation.Fee.DeliveryFee.ExportDeclarationFees != nil {
		deliveryFee += quotation.Fee.DeliveryFee.ExportDeclarationFees.Amount
	}
	if quotation.Fee.DeliveryFee.FreightCharges != nil {
		deliveryFee += quotation.Fee.DeliveryFee.FreightCharges.Amount
	}
	if quotation.Fee.DeliveryFee.Insurance != nil {
		deliveryFee += quotation.Fee.DeliveryFee.Insurance.Amount
	}

	return CalculatePricesResult{
		SubTotalPrice: subTotalPrice,
		TotalPrice:    totalPrice,
		DiscountValue: totalDiscountValue,
		DeliveryFee:   deliveryFee,
		//TODO: We need to handle different currencies of products
		CurrencyCode: "USD",
	}
}

func (s *quotationOrderConvertor) mapOrderItems(quotation *quotation_models.QuotationResponse) []order_models.OrderItem {
	orderItems := make([]order_models.OrderItem, len(quotation.Products))
	for i, product := range quotation.Products {
		attributes := make([]map[string]any, len(product.Attributes))
		for j, attribute := range product.Attributes {
			attributes[j] = map[string]any{
				attribute.Code: attribute.Value,
			}
		}

		// Initialize default values
		var name types.Localizable[string]
		var description *types.Localizable[string]
		var images *[]types.FileFieldDto
		var subTotalPrice float64
		var discount *product_models.Discount
		var totalPrice float64

		// Safe access to product data with nil checks
		if product.Product != nil {
			name = product.Product.Name
			description = &product.Product.Description
			images = &product.Product.Images

			// Safe access to price data with nil checks
			if product.Product.Price != nil {
				subTotalPrice = product.Product.Price.SubTotalPrice
				discount = product.Product.Price.Discount
				totalPrice = product.Product.Price.TotalPrice
			} else {
				// Fallback to quotation product price if Product.Price is nil
				subTotalPrice = product.Price.Amount
				totalPrice = product.Price.Amount
			}
		} else {
			// Fallback values when Product is nil - create a simple localizable name
			name = types.Localizable[string]{enums.EN: "Unknown Product"}
			subTotalPrice = product.Price.Amount
			totalPrice = product.Price.Amount
		}

		orderItems[i] = order_models.OrderItem{
			ItemId:        product.ProductId,
			Name:          name,
			Description:   description,
			Images:        images,
			Quantity:      product.Quantity,
			SubTotalPrice: subTotalPrice,
			Discount:      discount,
			TotalPrice:    totalPrice,
			Attributes:    attributes,
			Price: product_models.PriceItem{
				Value:        product.Price.Amount,
				CurrencyCode: product.Price.CurrencyCode,
			},
		}
	}

	return orderItems
}

func (s *quotationOrderConvertor) mapDeliveryInfo(quotation *quotation_models.QuotationResponse, deliveryFee float64) *order_models.DeliveryInfo {
	shippingCompanyId := &primitive.ObjectID{}
	if quotation.ShippingCompanyId != nil {
		shippingCompanyId = quotation.ShippingCompanyId
	} else {
		shippingCompanyId = quotation.AgentFromId
	}

	deliveryInfo := order_models.DeliveryInfo{
		DeliveryType: order_enums.B2B,
		B2BInfo: &order_models.B2BInfo{
			DeliveryFee:        deliveryFee,
			TransportCompanyId: shippingCompanyId,
			TrackingNumber:     order_utils.GenerateOrderCode(),
			//We will set the shipment details during the order processing
			ShipmentDetails: &order_models.ShipmentDetails{
				ShipmentMethod: string(quotation.Delivery.ShippingMethod),
				ShipmentOrigin: order_models.ShipmentDestination{
					Country: &quotation.Delivery.FromLocation.CountryIso3,
					City:    quotation.Delivery.FromLocation.City,
					Area:    quotation.Delivery.FromLocation.Area,
					Address: quotation.Delivery.FromLocation.Address,
				},
				ShipmentDestination: order_models.ShipmentDestination{
					Country: &quotation.Delivery.ToLocation.CountryIso3,
					City:    quotation.Delivery.ToLocation.City,
					Area:    quotation.Delivery.ToLocation.Area,
					Address: quotation.Delivery.ToLocation.Address,
				},
			},
		},
	}

	return &deliveryInfo
}

func (s *quotationOrderConvertor) mapHistory(quotation *quotation_models.QuotationResponse) []order_models.LogHistory {
	history := order_models.LogHistory{
		Action:      order_enums.CREATED,
		Description: "Order converted from quotation",
		Date:        time.Now(),
	}
	return []order_models.LogHistory{history}
}

func (s *quotationOrderConvertor) mapDocuments(quotation *quotation_models.QuotationResponse) []order_models.Document {
	documents := make([]order_models.Document, len(quotation.Attachments))
	for i, attachment := range quotation.Attachments {
		documents[i] = order_models.Document{
			File: attachment,
		}
	}

	return documents
}

func (s *quotationOrderConvertor) mapComments(quotation *quotation_models.QuotationResponse) []order_models.Comment {
	comments := make([]order_models.Comment, len(quotation.Comments))
	for i, comment := range quotation.Comments {
		comments[i] = order_models.Comment{
			Comment: comment.Comment,
		}
	}
	return comments
}

func (s *quotationOrderConvertor) mapBuyerInfo(quotation *quotation_models.QuotationResponse) order_models.OrderParticipantInfo {
	buyerLocation := order_models.ShipmentDestination{}
	if quotation.QuotationBuyerInfo.BuyerProfile != nil {
		buyerLocation = order_models.ShipmentDestination{
			Country: quotation.QuotationBuyerInfo.BuyerProfile.Country,
			City:    quotation.QuotationBuyerInfo.BuyerProfile.City,
			Area:    quotation.QuotationBuyerInfo.BuyerProfile.Area,
			Address: quotation.QuotationBuyerInfo.BuyerProfile.Address,
		}
	}

	return order_models.OrderParticipantInfo{
		Id:       *quotation.CreatedBy,
		FullName: quotation.QuotationBuyerInfo.ContactPerson,
		Email:    &quotation.QuotationBuyerInfo.EmailAddress,
		Phone:    &quotation.QuotationBuyerInfo.PhoneNumber,
		Type:     order_enums.BUYER,
		Location: &buyerLocation,
	}
}

func (s *quotationOrderConvertor) mapVendorInfo(ctx context.Context, quotation *quotation_models.QuotationResponse) order_models.OrderParticipantInfo {
	if quotation.VendorProfile == nil {
		return order_models.OrderParticipantInfo{}
	}

	vendorLocation := order_models.ShipmentDestination{
		Country: quotation.VendorProfile.Country,
		City:    quotation.VendorProfile.City,
		Area:    quotation.VendorProfile.Area,
		Address: quotation.VendorProfile.Address,
	}

	profile, _ := s.profileService.GetById(ctx, quotation.VendorProfile.Id, "")
	return order_models.OrderParticipantInfo{
		Id:       profile.Id,
		FullName: *profile.FirstName + " " + *profile.LastName,
		Email:    profile.Email,
		Phone:    profile.Phone,
		Type:     order_enums.VENDOR,
		Location: &vendorLocation,
	}
}
