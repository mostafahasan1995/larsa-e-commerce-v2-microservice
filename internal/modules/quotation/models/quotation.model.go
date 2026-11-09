package quotation_models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
)

type BaseQuotation struct {
	models.BaseModel `bson:",inline"`

	Title  string                          `bson:"title" json:"title"`
	Code   string                          `bson:"code" json:"code"`
	Type   quotation_enums.QuotationType   `bson:"type" json:"type"`
	Status quotation_enums.QuotationStatus `bson:"status" json:"status"`

	Deadline *time.Time `bson:"deadline,omitempty" json:"deadline,omitempty"`

	Delivery Delivery     `bson:"delivery" json:"delivery"`
	Fee      QuotationFee `bson:"fee" json:"fee"`
	Payment  Payment      `bson:"payment" json:"payment"`

	DiscountAmount *QuotationPrice `bson:"discountAmount" json:"discountAmount"`

	WarrantyAndValidity WarrantyAndValidity `bson:"warrantyAndValidity" json:"warrantyAndValidity"`

	PackagingRequirements PackagingRequirements `bson:"packagingRequirements" json:"packagingRequirements"`

	AdditionalNotes AdditionalNotes `bson:"additionalNotes" json:"additionalNotes"`

	Attachments []types.FileFieldDto `bson:"attachments" json:"attachments"`

	Signature *Signature `bson:"signature,omitempty" json:"signature,omitempty"`

	RejectedReason *RejectedReason `bson:"rejectedReason,omitempty" json:"rejectedReason,omitempty"`

	Comments []Comment          `bson:"comments,omitempty" json:"comments,omitempty"`
	History  []QuotationHistory `bson:"history,omitempty" json:"history,omitempty"`

	RequestQuotationId *primitive.ObjectID `bson:"requestQuotationId,omitempty" json:"requestQuotationId,omitempty"`
	GeneralQuotationId *primitive.ObjectID `bson:"generalQuotationId,omitempty" json:"generalQuotationId,omitempty"`

	FromCountryIso3 *string `bson:"fromCountryIso3,omitempty" json:"fromCountryIso3,omitempty"`
	ToCountryIso3   *string `bson:"toCountryIso3,omitempty" json:"toCountryIso3,omitempty"`

	VendorProfileId *primitive.ObjectID `bson:"vendorProfileId,omitempty" json:"vendorProfileId,omitempty"`

	AgentFromId *primitive.ObjectID `bson:"agentFromId,omitempty" json:"agentFromId,omitempty"`
	//ShippingCompanyId could be null if the AgentFrom is responsible for the shipping
	ShippingCompanyId *primitive.ObjectID `bson:"shippingCompanyId,omitempty" json:"shippingCompanyId,omitempty"`
	AgentToId         *primitive.ObjectID `bson:"agentToId,omitempty" json:"agentToId,omitempty"`
}

type Quotation struct {
	BaseQuotation      `bson:",inline"`
	QuotationBuyerInfo *QuotationBuyerInfo `bson:"quotationBuyerInfo" json:"quotationBuyerInfo"`
	Products           []QuotationProduct  `bson:"products" json:"products"`
}

type QuotationResponse struct {
	BaseQuotation      `bson:",inline"`
	QuotationBuyerInfo QuotationBuyerInfoResponse      `bson:"quotationBuyerInfo" json:"quotationBuyerInfo"`
	VendorProfile      *profile_models.ProfileResponse `bson:"vendorProfile" json:"vendorProfile,omitempty"`

	Products []QuotationProductResponse `bson:"products" json:"products"`

	SubmittedQuotationCount int `bson:"submittedQuotationCount" json:"submittedQuotationCount"`
}

type QuotationsWithPagination = types.DataWithPagination[QuotationResponse]
