package order_dto

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	cart_enums "larsa-e-commerce-microservice-v2/internal/modules/carts/enums"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_shared "larsa-e-commerce-microservice-v2/internal/modules/order/shared"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOrderDTO struct {
	Type order_enums.OrderType `json:"type" validate:"required"`

	PickupInfo *CreatePickupInfoDTO `json:"pickupInfo,omitempty" validate:"omitempty"`

	DeliveryInfo *CreateDeliveryInfoDTO `json:"deliveryInfo,omitempty" validate:"omitempty"`

	TotalDiscountValue             *float64              `json:"totalDiscountValue,omitempty"`
	SubTotalPrice                  float64               `json:"subTotalPrice"`
	TotalPrice                     float64               `json:"totalPrice"`
	CouponDiscountValue            float64               `json:"couponDiscountValue"`
	LoyaltyDiscountValue           float64               `json:"loyaltyDiscountValue"`
	DeliveryFee                    float64               `json:"deliveryFee"`
	TotalPriceAfterProductDiscount float64               `json:"totalPriceAfterProductDiscount"`
	CurrencyCode                   string                `json:"currencyCode" validate:"required"`
	PaymentInfo                    *CreatePaymentInfoDTO `json:"paymentInfo" validate:"required"`
	CartId                         *primitive.ObjectID   `json:"cartId,omitempty"`

	InvoiceNumber *string `json:"invoiceNumber,omitempty"`

	Status string `json:"status" validate:"required"`

	Note *string `json:"note,omitempty"`

	BuyerInfo  CreateParticipantInfoDTO  `json:"buyerInfo" validate:"required"`
	VendorInfo *CreateParticipantInfoDTO `json:"vendorInfo" validate:"required"`
	OrderItems []CreateOrderItemDto      `json:"orderItems" validate:"required,dive"`

	Documents []CreateDocumentDTO `json:"documents,omitempty" validate:"omitempty,dive"`
	Comments  []CreateCommentDTO  `json:"comments,omitempty" validate:"omitempty,dive"`

	QuotationId *primitive.ObjectID `json:"quotationId,omitempty"`
	CouponCode  *string             `json:"couponCode,omitempty"`

	MigratedFrom *string `json:"migratedFrom,omitempty"`
}

type CreatePickupInfoDTO struct {
	PickupDateTime *time.Time          `json:"pickupDateTime,omitempty"`
	PickupBranchId *primitive.ObjectID `json:"pickupBranchId,omitempty"`
	EstimationTime *int64              `json:"estimationTime,omitempty"`
}

type CreateDeliveryInfoDTO struct {
	DeliveryType order_enums.DeliveryType `json:"deliveryType" validate:"required,oneof=b2b b2c"`
	B2BInfo      *CreateB2BInfoDTO        `json:"b2bInfo,omitempty" validate:"omitempty"`
	B2CInfo      *CreateB2CInfoDTO        `json:"b2cInfo,omitempty" validate:"omitempty"`
}

type CreateB2BInfoDTO struct {
	TrackingNumber     string                       `json:"trackingNumber" validate:"required"`
	TransportCompanyId *primitive.ObjectID          `json:"transportCompany,omitempty"`
	DeliveryFee        float64                      `json:"deliveryFee" validate:"required"`
	ShipmentDetails    *CreateB2BShipmentDetailsDTO `json:"shipmentDetails" validate:"required"`
}

type CreateB2CInfoDTO struct {
	CustomerAddress CreateShipmentDestinationDTO `json:"customerAddress" validate:"required"`
	DeliveryFee     float64                      `json:"deliveryFee" validate:"required"`
	DeliveryTime    *uint64                      `json:"deliveryTime,omitempty"`
	BranchId        *primitive.ObjectID          `json:"branchId,omitempty"`
}

type CreateB2BShipmentDetailsDTO struct {
	ShipmentMethod            string                                      `json:"shipmentMethod" validate:"required"`
	ShipmentDate              *time.Time                                  `json:"shipmentDate,omitempty"`
	ShipmentStatus            order_enums.ShipmentStatus                  `json:"shipmentStatus" validate:"required"`
	EstimationDeliveryDate    *time.Time                                  `json:"estimationDeliveryDate,omitempty"`
	DorToDorShipmentCost      float64                                     `json:"dorToDorShipmentCost" validate:"required"`
	ShipmentOrigin            CreateShipmentDestinationDTO                `json:"shipmentOrigin" validate:"required"`
	ShipmentDestination       CreateShipmentDestinationDTO                `json:"shipmentDestination" validate:"required"`
	ComplianceAndGoodsDetails *CreateShipmentComplianceAndGoodsDetailsDTO `json:"complianceAndGoodsDetails,omitempty" validate:"omitempty"`
}

type CreateShipmentComplianceAndGoodsDetailsDTO struct {
	Packing               []CreatePackingDTO          `json:"packing,omitempty" validate:"omitempty,dive"`
	CertificateOfOrigin   *CreateCertificateOriginDTO `json:"certificateOfOrigin,omitempty" validate:"omitempty"`
	ProductSpecifications *CreateProductSpecsDTO      `json:"productSpecifications,omitempty" validate:"omitempty"`
	QualityComplianceCert *types.FileFieldDto         `json:"qualityComplianceCert,omitempty"`
	ExportLicense         *CreateExportLicenseDTO     `json:"exportLicense,omitempty" validate:"omitempty"`
}

type CreatePaymentInfoDTO struct {
	Status        order_enums.PaymentStatus `json:"status" validate:"required"`
	PaymentMethod order_enums.PaymentMethod `json:"paymentMethod" validate:"required"`
	Payments      []CreatePaymentDTO        `json:"payments" validate:"required,dive"`
}

type CreateParticipantInfoDTO struct {
	Id                    primitive.ObjectID               `json:"id" validate:"required"`
	FullName              string                           `json:"fullName" validate:"required"`
	Type                  order_enums.OrderParticipantType `json:"type" validate:"required"`
	Email                 *string                          `json:"email,omitempty"`
	CountryCode           *string                          `json:"countryCode,omitempty"`
	Phone                 *string                          `json:"phone,omitempty"`
	AlternativeMobile     *string                          `json:"alternativeMobile,omitempty"`
	AlternativeMobileCode *string                          `json:"alternativeMobileCode,omitempty"`
	Location              CreateShipmentDestinationDTO     `json:"location" validate:"required"`
}

type CreateOrderItemDto struct {
	ItemId                primitive.ObjectID                     `json:"itemId" validate:"required"`
	Name                  types.Localizable[string]              `json:"name" validate:"required,localizable"`
	Product               *product_models.ProductResponse        `json:"product" validate:"required"`
	Quantity              uint64                                 `json:"quantity" validate:"required"`
	SubTotalPrice         float64                                `json:"subTotalPrice" validate:"required"`
	TotalPrice            float64                                `json:"totalPrice" validate:"required"`
	DiscountValue         float64                                `json:"discountValue"`
	Discount              *product_models.Discount               `json:"discount,omitempty"`
	Price                 product_models.PriceItem               `json:"price"`
	Attributes            []map[string]any                       `json:"attributes"`
	BranchId              *primitive.ObjectID                    `json:"branchId" validate:"required"`
	Type                  cart_enums.CartItemType                `json:"type" validate:"required"`
	Addons                []order_models.OrderAddon              `json:"addons"`
	AddonTotalPrice       float64                                `json:"addonTotalPrice"`
	AdditionalInformation *order_shared.AdditionalInformationDto `json:"additionalInformation,omitempty"`
}

// Test
type CreateShipmentDestinationDTO struct {
	Country   *string  `json:"country,omitempty"`
	City      *string  `json:"city,omitempty"`
	Area      *string  `json:"area,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Notes     *string  `json:"notes,omitempty"`
}

type CreatePackingDTO struct {
	Name        string           `json:"name" validate:"required"`
	Content     string           `json:"content"`
	Dimensions  string           `json:"dimensions"`
	Unit        order_enums.Unit `json:"unit" validate:"required"`
	TotalWeight float64          `json:"totalWeight"`
}

type CreateCertificateOriginDTO struct {
	File            types.FileFieldDto `json:"file"`
	CountryOfOrigin string             `json:"countryOfOrigin"`
}

type CreateProductSpecsDTO struct {
	Details string             `json:"details"`
	File    types.FileFieldDto `json:"file"`
}

type CreateExportLicenseDTO struct {
	File            types.FileFieldDto `json:"file"`
	LicenseRequired bool               `json:"licenseRequired"`
}

type CreatePaymentDTO struct {
	Amount        float64                   `json:"amount" validate:"required"`
	PaymentStatus order_enums.PaymentStatus `json:"paymentStatus" validate:"required"`
	PaymentMethod order_enums.PaymentMethod `json:"paymentMethod" validate:"required"`
	PaymentDate   time.Time                 `json:"paymentDate"`
	PaymentNumber *string                   `json:"paymentNumber,omitempty"`
	Attachment    *types.FileFieldDto       `json:"attachment,omitempty"`
	PaidTo        *order_enums.PayTo        `json:"paidTo,omitempty"`
}

type CreateLogHistoryDTO struct {
	Action      order_enums.LogActionType `json:"action" validate:"required"`
	Description string                    `json:"description"`
	Date        time.Time                 `json:"date"`
}

type CreateDocumentDTO struct {
	File  types.FileFieldDto `json:"file" validate:"required"`
	Title string             `json:"title" validate:"required"`
}
type CreateCommentDTO struct {
	Comment string `json:"comment" validate:"required"`
}

type UpdateOrderStatusDTO struct {
	Status string `json:"status" validate:"required"`
}

type UpdateOrderDTO struct {
	Document        *CreateDocumentDTO        `json:"document,omitempty" validate:"omitempty"`
	Comment         *CreateCommentDTO         `json:"comment,omitempty" validate:"omitempty"`
	ParticipantInfo *CreateParticipantInfoDTO `json:"participantInfo,omitempty" validate:"omitempty"`
	PaymentInfo     *CreatePaymentInfoDTO     `json:"paymentInfo,omitempty" validate:"omitempty"`
	PickupInfo      *CreatePickupInfoDTO      `json:"pickupInfo,omitempty" validate:"omitempty"`
	DeliveryInfo    *CreateDeliveryInfoDTO    `json:"deliveryInfo,omitempty" validate:"omitempty"`
}

type PayOrderDTO struct {
	OrderId        primitive.ObjectID         `json:"orderId" validate:"required"`
	PaymentGateway order_enums.PaymentGateway `json:"paymentGateway" validate:"required"`
}

// type OrderReportDTO struct {
// 	StartDate     string `json:"startDate"`
// 	EndDate       string `json:"endDate"`
// 	PaymentMethod string `json:"paymentMethod"`
// 	DeliveryType  string `json:"deliveryType"`

// }

type OrderReportDTO struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	StartHour     string `json:"startHour"`
	EndHour       string `json:"endHour"`
	PaymentMethod string `json:"paymentMethod"`
	BranchId      string `json:"branchId"`
	ProductId     string `json:"productId"`
	DeliveryType  string `json:"deliveryType"` // e.g. "standard", "express"
	Page          int    `json:"page"`         // new
	Limit         int    `json:"limit"`        // new
	Period        string `json:"period"`

	Status                    []string `json:"status"`
	PeriodType                string   `json:"periodType"` // "weekly", "monthly", "yearly"
	Year                      int      `json:"year"`
	Month                     int      `json:"month"`                     // required for monthly & weekly
	Week                      int      `json:"week"`                      // required for weekly
	CategoryId                string   `json:"categoryId"`                // new
	IncludeOnlineUnPaidOrders bool     `json:"includeOnlineUnPaidOrders"` // include unpaid online orders (default false - excludes them)
}
type OrderItemFilterDTO struct {
	StartDate  string `json:"startDate"` // RFC3339
	EndDate    string `json:"endDate"`
	CategoryID string `json:"categoryId"`
	BranchID   string `json:"branchId"`
}

type UpdateAdminSeenDTO struct {
	AdminSeen bool `json:"adminSeen" validate:"required"`
}
