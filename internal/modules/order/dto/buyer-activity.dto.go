package order_dto

// BuyerActivityResponse contains statistics about a buyer's recent activity
type BuyerActivityResponse struct {
	TotalOrders      int     `json:"totalOrders"`
	DeliveredOrders  int     `json:"deliveredOrders"`
	CancelledOrders  int     `json:"cancelledOrders"`
	PendingOrders    int     `json:"pendingOrders"`
	RejectedOrders   int     `json:"rejectedOrders"`
	AcceptedOrders   int     `json:"acceptedOrders"`
	TotalSpentAmount float64 `json:"totalSpentAmount"`
}
