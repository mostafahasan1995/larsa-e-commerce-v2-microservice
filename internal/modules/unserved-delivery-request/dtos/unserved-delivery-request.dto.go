package unserved_delivery_request_dtos

// UnservedDeliveryRequestAggregateItem represents one row in the aggregated response
type UnservedDeliveryRequestAggregateItem struct {
	Country *string `json:"country,omitempty"`
	City    *string `json:"city,omitempty"`
	Area    *string `json:"area,omitempty"`
	Date    string  `json:"date,omitempty"` // Format: 2006-01-02
	Count   int64   `json:"count"`
}

// UnservedDeliveryRequestAggregateResponse is the response for the aggregate endpoint
type UnservedDeliveryRequestAggregateResponse struct {
	GroupBy string                              `json:"groupBy"` // "country", "city", "area", "date"
	Data    []UnservedDeliveryRequestAggregateItem `json:"data"`
}
