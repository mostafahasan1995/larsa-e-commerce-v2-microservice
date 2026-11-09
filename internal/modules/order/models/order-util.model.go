package order_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	cart_enums "larsa-e-commerce-microservice-v2/internal/modules/carts/enums"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_shared "larsa-e-commerce-microservice-v2/internal/modules/order/shared"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PickupInfo struct {
	PickupDateTime *time.Time `bson:"pickupDateTime,omitempty" json:"pickupDateTime,omitempty"`
	//EstimationTime in minutes
	EstimationTime *int64                 `bson:"estimationTime,omitempty" json:"estimationTime,omitempty"`
	PickupBranchId *primitive.ObjectID    `bson:"pickupBranchId,omitempty" json:"pickupBranchId,omitempty"`
	PickupBranch   *product_models.Branch `bson:"pickupBranch,omitempty" json:"pickupBranch,omitempty"`
}

type DeliveryInfo struct {
	DeliveryType order_enums.DeliveryType `bson:"deliveryType" json:"deliveryType"`

	B2BInfo *B2BInfo `bson:"b2bInfo,omitempty" json:"b2bInfo,omitempty"`
	B2CInfo *B2CInfo `bson:"b2cInfo,omitempty" json:"b2cInfo,omitempty"`
}

type B2BInfo struct {
	TrackingNumber     string              `bson:"trackingNumber" json:"trackingNumber"`
	TransportCompanyId *primitive.ObjectID `bson:"transportCompany,omitempty" json:"transportCompany,omitempty"`
	DeliveryFee        float64             `bson:"deliveryFee" json:"deliveryFee"`

	ShipmentDetails *ShipmentDetails `bson:"shipmentDetails" json:"shipmentDetails"`

	TransitionStatus *settings_sales_enums.TransitionStatusType `bson:"transitionStatus,omitempty" json:"transitionStatus,omitempty"`
}

type B2CInfo struct {
	CustomerAddress ShipmentDestination    `bson:"customerAddress" json:"customerAddress"`
	DeliveryFee     float64                `bson:"deliveryFee" json:"deliveryFee"`
	DeliveryTime    *uint64                `bson:"deliveryTime,omitempty" json:"deliveryTime,omitempty"`
	BranchId        *primitive.ObjectID    `bson:"branchId,omitempty" json:"branchId,omitempty"`
	Branch          *product_models.Branch `bson:"branch,omitempty" json:"branch,omitempty"`
}

//.....//

type PaymentInfo struct {
	Status        order_enums.PaymentStatus `bson:"status" json:"status"`
	PaymentMethod order_enums.PaymentMethod `bson:"paymentMethod" json:"paymentMethod"`
	Payments      []Payments                `bson:"payments" json:"payments"`
}

type Payments struct {
	PaymentMethod order_enums.PaymentMethod `bson:"paymentMethod" json:"paymentMethod"`
	PaymentStatus order_enums.PaymentStatus `bson:"paymentStatus" json:"paymentStatus"`
	PaymentDate   time.Time                 `bson:"paymentDate" json:"paymentDate"`
	PaymentAmount float64                   `bson:"paymentAmount" json:"paymentAmount"`
	PaymentNumber string                    `bson:"paymentNumber" json:"paymentNumber"`
	PaidTo        *order_enums.PayTo        `bson:"paidTo,omitempty" json:"paidTo,omitempty"`
	Attachment    *types.FileFieldDto       `bson:"attachment,omitempty" json:"attachment,omitempty"`
}

//......//

type OrderAddon struct {
	Label  types.Localizable[string] `bson:"label" json:"label"`
	Values []OrderAddonValue         `bson:"values" json:"values"`
}

type OrderAddonValue struct {
	Label    types.Localizable[string] `bson:"label" json:"label"`
	Quantity uint64                    `bson:"quantity" json:"quantity"`
}

type OrderItem struct {
	ItemId                primitive.ObjectID                     `bson:"itemId" json:"itemId"`
	Name                  types.Localizable[string]              `bson:"name" json:"name"`
	Description           *types.Localizable[string]             `bson:"description,omitempty" json:"description,omitempty"`
	Images                *[]types.FileFieldDto                  `bson:"images,omitempty" json:"images,omitempty"`
	Quantity              uint64                                 `bson:"quantity" json:"quantity"`
	SubTotalPrice         float64                                `bson:"subTotalPrice" json:"subTotalPrice"`
	TotalPrice            float64                                `bson:"totalPrice" json:"totalPrice"`
	DiscountValue         float64                                `bson:"discountValue" json:"discountValue"`
	Discount              *product_models.Discount               `bson:"discount,omitempty" json:"discount,omitempty"`
	Price                 product_models.PriceItem               `bson:"price" json:"price"`
	Attributes            []map[string]any                       `bson:"attributes" json:"attributes"`
	BranchId              *primitive.ObjectID                    `bson:"branchId,omitempty" json:"branchId,omitempty"`
	Type                  cart_enums.CartItemType                `bson:"type" json:"type"`
	Addons                []OrderAddon                           `json:"addons"`
	AddonTotalPrice       float64                                `json:"addonTotalPrice"`
	Product               *product_models.ProductResponse        `bson:"product,omitempty" json:"product,omitempty"`
	AdditionalInformation *order_shared.AdditionalInformationDto `bson:"additionalInformation,omitempty" json:"additionalInformation,omitempty"`
}

//......//

type OrderParticipantInfo struct {
	Id                    primitive.ObjectID               `bson:"id" json:"id"`
	FullName              string                           `bson:"fullName" json:"fullName"`
	Email                 *string                          `bson:"email,omitempty" json:"email,omitempty"`
	CountryCode           *string                          `bson:"countryCode,omitempty" json:"countryCode,omitempty"`
	Phone                 *string                          `bson:"phone,omitempty" json:"phone,omitempty"`
	AlternativeMobile     *string                          `bson:"alternativeMobile,omitempty" json:"alternativeMobile,omitempty"`
	AlternativeMobileCode *string                          `bson:"alternativeMobileCode,omitempty" json:"alternativeMobileCode,omitempty"`
	Location              *ShipmentDestination             `bson:"location,omitempty" json:"location,omitempty"`
	Type                  order_enums.OrderParticipantType `bson:"type" json:"type"`
	Vendor                *profile_models.BaseProfile      `bson:"vendor,omitempty" json:"vendor,omitempty"`
	Buyer                 *profile_models.BaseProfile      `bson:"buyer,omitempty" json:"buyer,omitempty"`
}

//......//

type LogHistory struct {
	Action      order_enums.LogActionType `bson:"action" json:"action"`
	Description string                    `bson:"description" json:"description"`
	Date        time.Time                 `bson:"date" json:"date"`
}

//......//

type Document struct {
	File  types.FileFieldDto `bson:"file" json:"file"`
	Title string             `bson:"title" json:"title"`
}

//......//

type Comment struct {
	Comment  string             `bson:"comment" json:"comment"`
	Date     time.Time          `bson:"date" json:"date"`
	UserName string             `bson:"userName" json:"userName"`
	UserType enums.UserType     `bson:"userType" json:"userType"`
	UserId   primitive.ObjectID `bson:"userId" json:"userId"`
}
