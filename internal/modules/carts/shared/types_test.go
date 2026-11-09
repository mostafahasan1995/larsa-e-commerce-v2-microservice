package cart_shared

import (
	"testing"

	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
)

func TestGetMineArgs_HasAddressIdAndDeliveryParams(t *testing.T) {
	// GetMineArgs uses addressId from query (user's chosen address) for delivery fee and cart
	args := GetMineArgs{
		AddressId:     "507f1f77bcf86cd799439011",
		BranchId:      "507f1f77bcf86cd799439012",
		ZoneId:        "507f1f77bcf86cd799439013",
		OrderType:     order_enums.DELIVERY,
		CurrencyCode:  "USD",
		CalculateList: true,
	}
	if args.AddressId != "507f1f77bcf86cd799439011" {
		t.Errorf("expected AddressId from query, got %s", args.AddressId)
	}
	if args.BranchId != "507f1f77bcf86cd799439012" {
		t.Errorf("expected BranchId, got %s", args.BranchId)
	}
	if args.OrderType != order_enums.DELIVERY {
		t.Error("expected OrderType DELIVERY")
	}
}
