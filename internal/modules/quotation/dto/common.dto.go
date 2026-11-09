package quotation_dto

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuotationPrice struct {
	Amount       float64 `json:"amount" binding:"required,min=0"`
	CurrencyCode string  `json:"currencyCode" binding:"required"`
}
type QuotationProductDTO struct {
	ProductId       primitive.ObjectID             `json:"productId" binding:"required"`
	Quantity        uint64                         `json:"quantity" binding:"required,min=1"`
	Price           QuotationPrice                 `json:"price" binding:"required"`
	DiscountPercent float64                        `json:"discountPercent"`
	Attributes      []QuotationProductAttributeDTO `json:"attributes" binding:"required"`
}

type QuotationDeliveryLocationDTO struct {
	CountryIso3 string   `json:"countryIso3" binding:"required"`
	City        *string  `json:"city" binding:"required"`
	Area        *string  `json:"area" binding:"required"`
	Address     *string  `json:"address" binding:"required"`
	Latitude    *float64 `json:"latitude" binding:"required"`
	Longitude   *float64 `json:"longitude" binding:"required"`
}

type DeliveryDTO struct {
	ProposedDeliveryDate time.Time                      `json:"proposedDeliveryDate" `
	DeliveryTerms        quotation_enums.Incoterms      `json:"deliveryTerms"`
	ToLocation           QuotationDeliveryLocationDTO   `json:"toLocation"`
	FromLocation         QuotationDeliveryLocationDTO   `json:"fromLocation"`
	ShippingMethod       quotation_enums.ShippingMethod `json:"shippingMethod"`
}

type PaymentDTO struct {
	PreferredPaymentMethod quotation_enums.PaymentMethod   `json:"preferredPaymentMethod" `
	PaymentSchedule        quotation_enums.PaymentSchedule `json:"paymentSchedule" `
	AdvancePaymentPercent  float64                         `json:"advancePaymentPercent" `
}

type WarrantyAndValidityDTO struct {
	Details                   string    `json:"details" `
	ReturnReplacementPolicies string    `json:"returnReplacementPolicies" `
	QuotationValidityPeriod   time.Time `json:"quotationValidityPeriod" `
}

type AdditionalNotesDTO struct {
	RemarksComments string                        `json:"remarksComments"`
	Commitment      bool                          `json:"commitment"`
	SpecialNotes    []quotation_enums.SpecialNote `json:"specialNotes"`
}

type SignatureDTO struct {
	OfferResponsibleName string                   `json:"offerResponsibleName" `
	Position             quotation_enums.Position `json:"position" `
	OfferSubmissionDate  time.Time                `json:"offerSubmissionDate" `
	SignatureDoc         *types.FileFieldDto      `json:"signatureDoc"`
}

type CommentDTO struct {
	Comment string `json:"comment" validate:"required"`
}

type QuotationBuyerInfoDTO struct {
	UserId            primitive.ObjectID `json:"userId"  bson:"userId"`
	CompanyName       string             `json:"companyName" bson:"companyName" validate:"required"`
	ContactPerson     string             `json:"contactPerson" bson:"contactPerson" validate:"required"`
	CompanyAddress    string             `json:"companyAddress" bson:"companyAddress" validate:"required"`
	PhoneNumber       string             `json:"phoneNumber" bson:"phoneNumber" validate:"required"`
	EmailAddress      string             `json:"emailAddress" bson:"emailAddress" validate:"required"`
	BusinessRegNo     string             `json:"businessRegNo" bson:"businessRegNo" validate:"required"`
	TaxRegistrationNo string             `json:"taxRegistrationNo" bson:"taxRegistrationNo" validate:"required"`
}

type QuotationProductAttributeDTO struct {
	Code  string      `json:"code" binding:"required"`
	Value interface{} `json:"value" binding:"required"`

	Type quotation_enums.QuotationProductAttributeType `json:"type" binding:"required"`
}

type PackagingRequirementsDTO struct {
	PackingDetails         string    `json:"packingDetails"`
	LabelRequirement       string    `json:"labelRequirement"`
	RequiredCertifications string    `json:"requiredCertifications"`
	SampleRequirement      bool      `json:"sampleRequirement"`
	SampleDeliveryDeadline time.Time `json:"sampleDeliveryDeadline"`
}
