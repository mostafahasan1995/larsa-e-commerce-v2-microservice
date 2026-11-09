package quotation_types

import (
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	"time"
)

// QuotationFilter represents filters for quotation queries
type QuotationFilter struct {
	VendorName      string    `json:"vendorName,omitempty"`
	VendorRating    float64   `json:"vendorRating,omitempty"`
	QuotationNumber string    `json:"quotationNumber,omitempty"`
	ProductType     string    `json:"productType,omitempty"`
	MinPrice        float64   `json:"minPrice,omitempty"`
	MaxPrice        float64   `json:"maxPrice,omitempty"`
	DeliveryDate    time.Time `json:"deliveryDate,omitempty"`
	SubmissionDate  time.Time `json:"submissionDate,omitempty"`
}

// QuotationRequestFilter represents filters for quotation request queries
type QuotationRequestFilter struct {
	CustomerName           string `json:"customerName,omitempty"`
	RequestForQuotationNum string `json:"requestForQuotationNum,omitempty"`
	ProductType            string `json:"productType,omitempty"`
	QuotationState         string `json:"quotationState,omitempty"`
}

// QuotationRejectionResponse represents the response when rejecting a quotation
type QuotationRejectionResponse struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// QuotationSummary represents a summary of quotations by status
type QuotationSummary struct {
	PendingBuyer int `json:"pendingBuyer"`
	Accepted     int `json:"accepted"`
	Rejected     int `json:"rejected"`
	Expired      int `json:"expired"`
	Archived     int `json:"archived"`
}

// QuotationActivity represents a recent activity for a quotation
type QuotationActivity struct {
	Date        time.Time                              `json:"date"`
	QuotationId string                                 `json:"quotationId"`
	Status      quotation_enums.QuotationStatus        `json:"status"`
	BuyerName   string                                 `json:"buyerName"`
	Activity    quotation_enums.QuotationHistoryAction `json:"activity"`
}

// ProductRequestStats represents statistics for product requests
type ProductRequestStats struct {
	ProductId    string      `json:"productId"`
	ProductName  interface{} `json:"productName"`
	RequestCount int         `json:"requestCount"`
}

// QuotationInsights represents insights about buyer responses
type QuotationInsights struct {
	PositiveFeedbackPercentage float64 `json:"positiveFeedbackPercentage"`
	NegativeFeedbackPercentage float64 `json:"negativeFeedbackPercentage"`
}
