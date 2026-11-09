package quotation_models

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
)

type BaseQuotationBuyerInfo struct {
	BuyerProfileId primitive.ObjectID `bson:"buyerProfileId" json:"buyerProfileId"`

	ContactPerson     string `bson:"contactPerson" json:"contactPerson"`
	PhoneNumber       string `bson:"phoneNumber" json:"phoneNumber"`
	EmailAddress      string `bson:"emailAddress" json:"emailAddress"`
	CompanyName       string `bson:"companyName" json:"companyName"`
	CompanyAddress    string `bson:"companyAddress" json:"companyAddress"`
	BusinessRegNo     string `bson:"businessRegNo" json:"businessRegNo"`
	TaxRegistrationNo string `bson:"taxRegistrationNo" json:"taxRegistrationNo"`
}

type QuotationBuyerInfo struct {
	BaseQuotationBuyerInfo `bson:",inline"`
}

type QuotationBuyerInfoResponse struct {
	BaseQuotationBuyerInfo `bson:",inline"`

	BuyerProfile *profile_models.ProfileResponse `bson:"buyerProfile" json:"buyerProfile"`
}

type QuotationPrice struct {
	Amount       float64 `bson:"amount" json:"amount"`
	CurrencyCode string  `bson:"currencyCode" json:"currencyCode"`
}

//

type BaseQuotationProduct struct {
	ProductId       primitive.ObjectID          `bson:"productId" json:"productId"`
	Attributes      []QuotationProductAttribute `bson:"attributes" json:"attributes"`
	Price           QuotationPrice              `bson:"price" json:"price"`
	DiscountPercent float64                     `bson:"discountPercent" json:"discountPercent"`
	Quantity        uint64                      `bson:"quantity" json:"quantity"`
}

type QuotationProduct struct {
	BaseQuotationProduct `bson:",inline"`
}

type QuotationProductResponse struct {
	BaseQuotationProduct `bson:",inline"`
	Product              *product_models.ProductResponse `bson:"product" json:"product"`
}

type QuotationProductAttribute struct {
	Code  string                                        `bson:"code" json:"code"`
	Value interface{}                                   `bson:"value" json:"value"`
	Type  quotation_enums.QuotationProductAttributeType `bson:"type" json:"type"`
}

//

type QuotationDeliveryLocation struct {
	CountryIso3 string   `bson:"countryIso3" json:"countryIso3" validate:"required"`
	City        *string  `bson:"city" json:"city"`
	Area        *string  `bson:"area" json:"area"`
	Address     *string  `bson:"address" json:"address"`
	Latitude    *float64 `bson:"latitude" json:"latitude"`
	Longitude   *float64 `bson:"longitude" json:"longitude"`
}

type Delivery struct {
	//TODO: check if Delivery Date is equal to ProposedDeliveryDate.
	ProposedDeliveryDate time.Time `bson:"proposedDeliveryDate" json:"proposedDeliveryDate" validate:"required"`

	//TODO: make sure where to get FromLocation from, I think it is the same as vendor location, if it is the same, we need to check weather we need to save it here or not.
	DeliveryTerms quotation_enums.Incoterms `bson:"deliveryTerms" json:"deliveryTerms" validate:"required"`

	FromLocation QuotationDeliveryLocation `bson:"fromLocation" json:"fromLocation" validate:"required"`

	//TODO: we need to set the To delivery location when a vendor replay to the quotation
	ToLocation QuotationDeliveryLocation `bson:"toLocation" json:"toLocation" validate:"required"`

	ShippingMethod quotation_enums.ShippingMethod `json:"shippingMethod"`
}

type QuotationFee struct {
	DeliveryFee   QuotationDeliveryFee `bson:"deliveryFee" json:"deliveryFee"`
	SuperAdminFee QuotationPrice       `bson:"superAdminFee" json:"superAdminFee"`

	TotalFee QuotationPrice `bson:"totalFee" json:"totalFee"`
}

type QuotationDeliveryFee struct {
	// Set by to agent
	CustomsClearanceFee   *QuotationPrice `bson:"customsClearanceFee,omitempty" json:"customsClearanceFee,omitempty"`
	ImportDeclarationFees *QuotationPrice `bson:"importDeclarationFees,omitempty" json:"importDeclarationFees,omitempty"`
	DeliveryToBuyerFee    *QuotationPrice `bson:"deliveryToBuyerFee,omitempty" json:"deliveryToBuyerFee,omitempty"`

	// Set by from agent
	ExportDeclarationFees *QuotationPrice `bson:"exportDeclarationFees,omitempty" json:"exportDeclarationFees,omitempty"`

	// Set by shipping company
	FreightCharges *QuotationPrice `bson:"freightCharges,omitempty" json:"freightCharges,omitempty"`
	Insurance      *QuotationPrice `bson:"insurance,omitempty" json:"insurance,omitempty"`
}

type Payment struct {
	PreferredPaymentMethod quotation_enums.PaymentMethod   `bson:"preferredPaymentMethod" json:"preferredPaymentMethod"`
	PaymentSchedule        quotation_enums.PaymentSchedule `bson:"paymentSchedule" json:"paymentSchedule"`
	AdvancePaymentPercent  *float64                        `bson:"advancePaymentPercent,omitempty" json:"advancePaymentPercent,omitempty"`
}

type WarrantyAndValidity struct {
	Details                   string    `bson:"details" json:"details"`
	ReturnReplacementPolicies string    `bson:"returnReplacementPolicies" json:"returnReplacementPolicies"`
	QuotationValidityPeriod   time.Time `bson:"quotationValidityPeriod" json:"quotationValidityPeriod"`
}

type PackagingRequirements struct {
	PackingDetails         string    `bson:"packingDetails" json:"packingDetails"`
	LabelRequirement       string    `bson:"labelRequirement" json:"labelRequirement"`
	RequiredCertifications string    `bson:"requiredCertifications" json:"requiredCertifications"`
	SampleRequirement      bool      `bson:"sampleRequirement" json:"sampleRequirement"`
	SampleDeliveryDeadline time.Time `bson:"sampleDeliveryDeadline" json:"sampleDeliveryDeadline"`
}

type AdditionalNotes struct {
	RemarksComments string                        `bson:"remarksComments" json:"remarksComments"`
	Commitment      bool                          `bson:"commitment" json:"commitment"`
	SpecialNotes    []quotation_enums.SpecialNote `bson:"specialNotes" json:"specialNotes"`
}

type Signature struct {
	OfferResponsibleName string                   `bson:"offerResponsibleName" json:"offerResponsibleName"`
	Position             quotation_enums.Position `bson:"position" json:"position"`
	OfferSubmissionDate  time.Time                `bson:"offerSubmissionDate" json:"offerSubmissionDate"`
	SignatureDoc         types.FileFieldDto       `bson:"signatureDoc" json:"signatureDoc"`
}

type RejectedReason struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type Comment struct {
	Comment     string              `bson:"comment" json:"comment"`
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	CreatedById *primitive.ObjectID `bson:"createdById" json:"createdById"`
}

type QuotationHistory struct {
	Date           time.Time                              `bson:"date" json:"date"`
	Action         quotation_enums.QuotationHistoryAction `bson:"action" json:"action"`
	AdditionalInfo *string                                `bson:"additionalInfo,omitempty" json:"additionalInfo,omitempty"`
}

// GenerateCode generates a unique code for the quotation in the format: ()-year-(4 digit serial number)
func (q *Quotation) GenerateCode() string {
	year := time.Now().Year()
	uuid := uuid.New().String()[:4]
	if q.Type == quotation_enums.REQUEST {
		return "RFQ-" + strconv.Itoa(year) + "-" + uuid
	} else if q.Type == quotation_enums.RESPONSE {
		return "VQ-" + strconv.Itoa(year) + "-" + uuid
	} else {
		return "GQ-" + strconv.Itoa(year) + "-" + uuid
	}
}
