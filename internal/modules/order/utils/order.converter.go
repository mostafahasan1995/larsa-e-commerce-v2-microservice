package order_utils

import (
	"fmt"
	"net/http"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateOrderDTOToModel(dto *order_dto.CreateOrderDTO, cfg *types.AppCfg) (order_models.Order, *helpers.ApiError) {
	now := time.Now()
	userId := cfg.User.Id
	var userFullName string

	if cfg.User != nil && cfg.User.UserData.FirstName != "" {
		userFullName = fmt.Sprintf("%s %s", cfg.User.UserData.FirstName, cfg.User.UserData.LastName)
	} else {
		userFullName = "System"
	}
	userType := cfg.User.Type

	// Validate currency code
	if dto.CurrencyCode == "" {
		return order_models.Order{}, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Currency code is required",
		}
	}

	code := GenerateOrderCode()

	order := order_models.Order{
		BaseOrder: order_models.BaseOrder{
			MetaDataBaseModel: models.MetaDataBaseModel{
				Id:        primitive.NewObjectID(),
				CreatedAt: now,
				CreatedBy: &userId,
				UpdatedAt: now,
				UpdatedBy: &userId,
			},
			Code:                           code,
			Type:                           dto.Type,
			TotalDiscountValue:             dto.TotalDiscountValue,
			SubTotalPrice:                  dto.SubTotalPrice,
			TotalPrice:                     dto.TotalPrice,
			CurrencyCode:                   dto.CurrencyCode,
			InvoiceNumber:                  dto.InvoiceNumber,
			Status:                         dto.Status,
			QuotationId:                    dto.QuotationId,
			CouponCode:                     dto.CouponCode,
			Note:                           dto.Note,
			TotalPriceAfterProductDiscount: dto.TotalPriceAfterProductDiscount,
			CouponDiscountValue:            dto.CouponDiscountValue,
			LoyaltyDiscountValue:           dto.LoyaltyDiscountValue,
			CartId:                         dto.CartId,
			MigratedFrom:                   dto.MigratedFrom,
		},
	}

	if dto.Type == order_enums.PICKUP {
		if dto.PickupInfo == nil {
			return order_models.Order{}, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Pickup info is required",
			}
		}
		order.PickupInfo = MapPickupInfo(dto.PickupInfo)
		order.DeliveryInfo = nil
	}

	if dto.Type == order_enums.DELIVERY {
		if dto.DeliveryInfo == nil {
			return order_models.Order{}, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Delivery info is required",
			}
		}
		order.DeliveryInfo = MapDeliveryInfo(dto.DeliveryInfo, dto.Type)
		order.PickupInfo = nil
	}

	order.PaymentInfo = mapPaymentInfo(dto.PaymentInfo)

	order.OrderItems = mapOrderItems(dto.OrderItems)

	if dto.BuyerInfo.Id != primitive.NilObjectID {
		order.BuyerInfo = MapParticipantInfo(&dto.BuyerInfo)
	}

	if dto.VendorInfo != nil {
		order.VendorInfo = MapParticipantInfo(dto.VendorInfo)
	}
	order.Documents = mapDocuments(dto.Documents)
	order.Comments = mapComments(dto.Comments, now, userId, userFullName, userType)
	initialLog := createInitialLogHistory(now)
	order.LogHistory = []order_models.LogHistory{initialLog}
	return order, nil
}

func MapPickupInfo(dto *order_dto.CreatePickupInfoDTO) *order_models.PickupInfo {
	if dto == nil {
		return nil
	}
	return &order_models.PickupInfo{
		PickupDateTime: dto.PickupDateTime,
		PickupBranchId: dto.PickupBranchId,
		EstimationTime: dto.EstimationTime,
	}
}

func MapShipmentDestination(dto *order_dto.CreateShipmentDestinationDTO) order_models.ShipmentDestination {
	if dto == nil {
		return order_models.ShipmentDestination{}
	}
	model := order_models.ShipmentDestination{
		Address:   dto.Address,
		Area:      dto.Area,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
		Country:   dto.Country,
		City:      dto.City,
		Notes:     dto.Notes,
	}
	return model
}

func mapPacking(dtos []order_dto.CreatePackingDTO) []order_models.Packing {
	if dtos == nil {
		return []order_models.Packing{}
	}
	packings := make([]order_models.Packing, 0, len(dtos))
	for _, pDto := range dtos {
		packings = append(packings, order_models.Packing{
			Name:        pDto.Name,
			Content:     pDto.Content,
			Dimensions:  pDto.Dimensions,
			Unit:        pDto.Unit,
			TotalWeight: pDto.TotalWeight,
		})
	}
	return packings
}

func mapCertificateOfOrigin(dto *order_dto.CreateCertificateOriginDTO) *order_models.CertificateOfOrigin {
	if dto == nil {
		return nil
	}
	return &order_models.CertificateOfOrigin{
		File:            dto.File,
		CountryOfOrigin: dto.CountryOfOrigin,
	}
}

func mapProductSpecifications(dto *order_dto.CreateProductSpecsDTO) *order_models.ProductSpecifications {
	if dto == nil {
		return nil
	}
	return &order_models.ProductSpecifications{
		Details: dto.Details,
		File:    dto.File,
	}
}

func mapExportLicense(dto *order_dto.CreateExportLicenseDTO) *order_models.ExportLicense {
	if dto == nil {
		return nil
	}
	return &order_models.ExportLicense{
		File:            dto.File,
		LicenseRequired: dto.LicenseRequired,
	}
}

func mapShipmentComplianceAndGoodsDetails(dto *order_dto.CreateShipmentComplianceAndGoodsDetailsDTO) *order_models.ShipmentComplianceAndGoodsDetails {
	if dto == nil {
		return nil
	}
	details := &order_models.ShipmentComplianceAndGoodsDetails{
		QualityComplianceCert: dto.QualityComplianceCert,
	}
	details.Packing = mapPacking(dto.Packing)
	details.CertificateOfOrigin = mapCertificateOfOrigin(dto.CertificateOfOrigin)
	details.ProductSpecifications = mapProductSpecifications(dto.ProductSpecifications)
	details.ExportLicense = mapExportLicense(dto.ExportLicense)
	return details
}

func mapShipmentDetails(dto *order_dto.CreateB2BShipmentDetailsDTO) *order_models.ShipmentDetails {
	if dto == nil {
		return nil
	}
	details := &order_models.ShipmentDetails{
		ShipmentMethod:         dto.ShipmentMethod,
		ShipmentDate:           dto.ShipmentDate,
		ShipmentStatus:         dto.ShipmentStatus,
		EstimationDeliveryDate: dto.EstimationDeliveryDate,
		DorToDorShipmentCost:   dto.DorToDorShipmentCost,
	}
	details.ShipmentOrigin = MapShipmentDestination(&dto.ShipmentOrigin)
	details.ShipmentDestination = MapShipmentDestination(&dto.ShipmentDestination)
	details.ShipmentComplianceAndGoodsDetails = mapShipmentComplianceAndGoodsDetails(dto.ComplianceAndGoodsDetails)
	return details
}

func mapB2BInfo(dto *order_dto.CreateB2BInfoDTO) *order_models.B2BInfo {
	if dto == nil {
		return nil
	}
	b2bInfo := &order_models.B2BInfo{
		TrackingNumber:     dto.TrackingNumber,
		TransportCompanyId: dto.TransportCompanyId,
		DeliveryFee:        dto.DeliveryFee,
	}
	b2bInfo.ShipmentDetails = mapShipmentDetails(dto.ShipmentDetails)
	return b2bInfo
}

func mapB2CInfo(dto *order_dto.CreateB2CInfoDTO) *order_models.B2CInfo {
	if dto == nil {
		return nil
	}
	b2cInfo := &order_models.B2CInfo{
		DeliveryFee:     dto.DeliveryFee,
		BranchId:        dto.BranchId,
		DeliveryTime:    dto.DeliveryTime,
		CustomerAddress: MapShipmentDestination(&dto.CustomerAddress),
	}
	return b2cInfo
}

func MapDeliveryInfo(dto *order_dto.CreateDeliveryInfoDTO, orderType order_enums.OrderType) *order_models.DeliveryInfo {
	if dto == nil {
		return nil
	}
	deliveryInfo := &order_models.DeliveryInfo{
		DeliveryType: dto.DeliveryType,
	}
	if dto.DeliveryType == order_enums.B2B {
		deliveryInfo.B2BInfo = mapB2BInfo(dto.B2BInfo)
		deliveryInfo.B2CInfo = nil
	} else if dto.DeliveryType == order_enums.B2C {
		deliveryInfo.B2CInfo = mapB2CInfo(dto.B2CInfo)
		deliveryInfo.B2BInfo = nil
	}
	return deliveryInfo
}

func mapPaymentInfo(dto *order_dto.CreatePaymentInfoDTO) *order_models.PaymentInfo {
	if dto == nil {
		return &order_models.PaymentInfo{Payments: []order_models.Payments{}}
	}
	paymentInfo := &order_models.PaymentInfo{
		Status:        dto.Status,
		PaymentMethod: dto.PaymentMethod,
	}
	if dto.Payments != nil {
		payments := make([]order_models.Payments, 0, len(dto.Payments))
		for _, pDto := range dto.Payments {
			payments = append(payments, order_models.Payments{
				PaymentMethod: pDto.PaymentMethod,
				PaymentStatus: pDto.PaymentStatus,
				PaymentDate:   pDto.PaymentDate,
				PaymentAmount: pDto.Amount,
				PaidTo:        pDto.PaidTo,
				Attachment:    pDto.Attachment,
				PaymentNumber: "PMT-" + primitive.NewObjectID().Hex()[:8],
			})
		}
		paymentInfo.Payments = payments
	} else {
		paymentInfo.Payments = []order_models.Payments{}
	}
	return paymentInfo
}

func mapOrderItems(dtos []order_dto.CreateOrderItemDto) []order_models.OrderItem {
	var orderItems []order_models.OrderItem
	if dtos != nil {
		orderItems = make([]order_models.OrderItem, 0, len(dtos))
		for _, pDto := range dtos {
			orderItem := order_models.OrderItem{
				ItemId:                pDto.ItemId,
				BranchId:              pDto.BranchId,
				Type:                  pDto.Type,
				Quantity:              pDto.Quantity,
				Attributes:            pDto.Attributes,
				Name:                  pDto.Name,
				Price:                 pDto.Price,
				Discount:              pDto.Discount,
				DiscountValue:         pDto.DiscountValue,
				SubTotalPrice:         pDto.SubTotalPrice,
				TotalPrice:            pDto.TotalPrice,
				Addons:                pDto.Addons,
				Product:               pDto.Product,
				AddonTotalPrice:       pDto.AddonTotalPrice,
				AdditionalInformation: pDto.AdditionalInformation,
			}
			orderItems = append(orderItems, orderItem)
		}
	}
	return orderItems
}

func MapParticipantInfo(dto *order_dto.CreateParticipantInfoDTO) order_models.OrderParticipantInfo {
	if dto == nil {
		return order_models.OrderParticipantInfo{}
	}
	location := MapShipmentDestination(&dto.Location)
	return order_models.OrderParticipantInfo{
		Id:                    dto.Id,
		FullName:              dto.FullName,
		Email:                 dto.Email,
		Phone:                 dto.Phone,
		AlternativeMobile:     dto.AlternativeMobile,
		AlternativeMobileCode: dto.AlternativeMobileCode,
		Type:                  dto.Type,
		CountryCode:           dto.CountryCode,
		Location:              &location,
	}
}

func mapDocuments(dtos []order_dto.CreateDocumentDTO) []order_models.Document {
	if len(dtos) == 0 {
		return []order_models.Document{}
	}
	documents := make([]order_models.Document, 0, len(dtos))
	for _, dDto := range dtos {
		documents = append(documents, order_models.Document{
			File:  dDto.File,
			Title: dDto.Title,
		})
	}
	return documents
}

func mapComments(dtos []order_dto.CreateCommentDTO, now time.Time, userId primitive.ObjectID, userFullName string, userType enums.UserType) []order_models.Comment {
	if len(dtos) == 0 {
		return []order_models.Comment{}
	}
	comments := make([]order_models.Comment, 0, len(dtos))
	for _, cDto := range dtos {
		comments = append(comments, order_models.Comment{
			Comment:  cDto.Comment,
			Date:     now,
			UserId:   userId,
			UserName: userFullName,
			UserType: userType,
		})
	}
	return comments
}

func createInitialLogHistory(now time.Time) order_models.LogHistory {
	return order_models.LogHistory{
		Action:      order_enums.CREATED,
		Description: "Order created",
		Date:        now,
	}
}
