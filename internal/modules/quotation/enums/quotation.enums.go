package quotation_enums

type QuotationStatus string
type QuotationHistoryAction string
type QuotationProductAttributeType string
type Incoterms string
type PaymentSchedule string
type SpecialNote string
type Position string
type QuotationType string

// TODO: Handle canceled status
const (
	PENDING_ADMIN      QuotationStatus = "pending_admin"
	PENDING_VENDOR     QuotationStatus = "pending_vendor"
	READY_TO_SEND      QuotationStatus = "ready_to_send"
	SENT_TO_BUYER      QuotationStatus = "sent_to_buyer"
	REJECTED           QuotationStatus = "rejected"
	APPROVED           QuotationStatus = "approved"
	CONVERT_TO_ORDER   QuotationStatus = "convert_to_order"
	CONVERTED_TO_ORDER QuotationStatus = "converted_to_order"
	DRAFT              QuotationStatus = "draft"
	CLOSED             QuotationStatus = "closed"
	OPEN               QuotationStatus = "open"
)

const (
	QUOTATION_PRODUCT_ATTRIBUTE_TYPE_TEXT         QuotationProductAttributeType = "text"
	QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER       QuotationProductAttributeType = "number"
	QUOTATION_PRODUCT_ATTRIBUTE_TYPE_BOOLEAN      QuotationProductAttributeType = "boolean"
	QUOTATION_PRODUCT_ATTRIBUTE_TYPE_STRING_ARRAY QuotationProductAttributeType = "string_array"
	QUOTATION_PRODUCT_ATTRIBUTE_TYPE_NUMBER_ARRAY QuotationProductAttributeType = "number_array"
)

const (
	DDP Incoterms = "DDP" // Delivered Duty Paid
	EXW Incoterms = "EXW" // Ex Works
	FOB Incoterms = "FOB" // Free on Board
	CIF Incoterms = "CIF" // Cost, Insurance and Freight
	CPT Incoterms = "CPT" // Carriage Paid To
	CIP Incoterms = "CIP" // Carriage and Insurance Paid To
)

type PaymentMethod string

const (
	BANK_TRANSFER    PaymentMethod = "bank_transfer"
	CREDIT_CARD      PaymentMethod = "credit_card"
	CASH             PaymentMethod = "cash"
	LETTER_OF_CREDIT PaymentMethod = "letter_of_credit"
)

const (
	FINAL_PAYMENT_UPON_SHIPMENT PaymentSchedule = "final_payment_upon_shipment"
	ADVANCE_PAYMENT             PaymentSchedule = "advance_payment"
	INSTALLMENTS                PaymentSchedule = "installments"
)

const (
	VALID_FOR_15_DAYS        SpecialNote = "valid_for_15_days"
	VALID_FOR_30_DAYS        SpecialNote = "valid_for_30_days"
	VALID_FOR_60_DAYS        SpecialNote = "valid_for_60_days"
	SUBJECT_TO_STOCK         SpecialNote = "subject_to_stock"
	PRICES_SUBJECT_TO_CHANGE SpecialNote = "prices_subject_to_change"
)

const (
	SALES_MANAGER    Position = "sales_manager"
	SALES_EXECUTIVE  Position = "sales_executive"
	BUSINESS_MANAGER Position = "business_manager"
	ACCOUNT_MANAGER  Position = "account_manager"
	REGIONAL_MANAGER Position = "regional_manager"
	DIRECTOR         Position = "director"
)

const (
	GENERAL  QuotationType = "general"
	RESPONSE QuotationType = "response"
	REQUEST  QuotationType = "request"
)

const (
	QUOTATION_CREATED   QuotationHistoryAction = "quotation_created"
	QUOTATION_UPDATED   QuotationHistoryAction = "quotation_updated"
	SUBMITTED_TO_BEEBUY QuotationHistoryAction = "submitted_to_beebuy"
	REVIEWED_BY_BEEBUY  QuotationHistoryAction = "reviewed_by_beebuy"
	QUOTATION_SENT      QuotationHistoryAction = "quotation_sent"
	QUOTATION_REJECTED  QuotationHistoryAction = "quotation_rejected"
	QUOTATION_ACCEPTED  QuotationHistoryAction = "quotation_accepted"
	ORDER_CONVERSION    QuotationHistoryAction = "order_conversion"
	COMMENT_ADDED       QuotationHistoryAction = "comment_added"
)
