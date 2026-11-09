package order_utils

import (
	"fmt"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
)

// GetOrderExportColumns returns the column definitions for order export
func GetOrderExportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("orderCode", "Order Code", excel_services.IntPtr(50)),
		excel_services.NewRequiredStringColumn("status", "Status", excel_services.IntPtr(50)),
		excel_services.NewFloatColumn("totalPrice", "Total Price", true, excel_services.Float64Ptr(0), nil),
		excel_services.NewFloatColumn("totalPaid", "Total Paid", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewRequiredStringColumn("buyerName", "Buyer Name", excel_services.IntPtr(100)),
		excel_services.NewEmailColumn("buyerEmail", "Buyer Email", false),
		excel_services.NewOptionalStringColumn("buyerPhone", "Buyer Phone", "", excel_services.IntPtr(20)),
		excel_services.NewEnumColumn("paymentStatus", "Payment Status", true,
			[]string{"PAID", "UNPAID", "PARTIALLY_PAID"}, "UNPAID"),
		excel_services.NewOptionalStringColumn("deliveryAddress", "Delivery Address", "", excel_services.IntPtr(200)),
		excel_services.NewOptionalStringColumn("deliveryCity", "Delivery City", "", excel_services.IntPtr(50)),
		excel_services.NewOptionalStringColumn("estimationTime", "Estimation Time", "", excel_services.IntPtr(50)),
		excel_services.NewFloatColumn("deliveryFee", "Delivery Fee", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewOptionalStringColumn("deliveryTime", "Delivery Time", "", excel_services.IntPtr(50)),
		excel_services.NewOptionalStringColumn("branchName", "Branch Name", "", excel_services.IntPtr(100)),
		excel_services.NewFloatColumn("totalDiscountValue", "Total Discount Value", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewRequiredStringColumn("currencyCode", "Currency Code", excel_services.IntPtr(10)),
		excel_services.NewOptionalStringColumn("couponCode", "Coupon Code", "", excel_services.IntPtr(50)),
		excel_services.NewFloatColumn("couponDiscountValue", "Coupon Discount Value", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewOptionalStringColumn("note", "Note", "", excel_services.IntPtr(500)),
		excel_services.NewOptionalStringColumn("orderItems", "Order Items", "", excel_services.IntPtr(1000)),
		excel_services.NewOptionalStringColumn("addons", "Addons", "", excel_services.IntPtr(1000)),
		excel_services.NewFloatColumn("addonTotalPrice", "Addon Total Price", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewDateTimeColumn("createdAt", "Created Date", true, "2006-01-02 15:04:05"),
		excel_services.NewDateTimeColumn("updatedAt", "Updated Date", true, "2006-01-02 15:04:05"),
	}
}

// ConvertOrderToExportData converts an order response to export data map
func ConvertOrderToExportData(order order_models.OrderResponse) map[string]interface{} {

	data := map[string]interface{}{
		"orderCode":     order.Code,
		"status":        order.Status,
		"totalPrice":    order.TotalPrice,
		"buyerName":     order.BuyerInfo.FullName,
		"paymentStatus": string(order.PaymentInfo.Status),
		"createdAt":     order.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt":     order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Add buyer contact info
	if order.BuyerInfo.Email != nil {
		data["buyerEmail"] = *order.BuyerInfo.Email
	} else {
		data["buyerEmail"] = ""
	}
	if order.BuyerInfo.Phone != nil {
		data["buyerPhone"] = *order.BuyerInfo.Phone
	} else {
		data["buyerPhone"] = ""
	}

	// Add delivery info
	if order.DeliveryInfo != nil && order.DeliveryInfo.B2CInfo != nil {
		if order.DeliveryInfo.B2CInfo.CustomerAddress.Address != nil {
			data["deliveryAddress"] = *order.DeliveryInfo.B2CInfo.CustomerAddress.Address
		} else {
			data["deliveryAddress"] = ""
		}
		if order.DeliveryInfo.B2CInfo.CustomerAddress.City != nil {
			data["deliveryCity"] = *order.DeliveryInfo.B2CInfo.CustomerAddress.City
		} else {
			data["deliveryCity"] = ""
		}
		if order.DeliveryInfo.B2CInfo.CustomerAddress.Notes != nil {
			data["deliveryNotes"] = *order.DeliveryInfo.B2CInfo.CustomerAddress.Notes
		} else {
			data["deliveryNotes"] = ""
		}
		data["deliveryFee"] = order.DeliveryInfo.B2CInfo.DeliveryFee
		if order.DeliveryInfo.B2CInfo.DeliveryTime != nil {
			data["deliveryTime"] = fmt.Sprintf("%d minutes", *order.DeliveryInfo.B2CInfo.DeliveryTime)
		} else {
			data["deliveryTime"] = ""
		}
		data["estimationTime"] = ""
	} else if order.DeliveryInfo != nil && order.DeliveryInfo.B2BInfo != nil && order.DeliveryInfo.B2BInfo.ShipmentDetails != nil {
		if order.DeliveryInfo.B2BInfo.ShipmentDetails.ShipmentDestination.Address != nil {
			data["deliveryAddress"] = *order.DeliveryInfo.B2BInfo.ShipmentDetails.ShipmentDestination.Address
		} else {
			data["deliveryAddress"] = ""
		}
		if order.DeliveryInfo.B2BInfo.ShipmentDetails.ShipmentDestination.City != nil {
			data["deliveryCity"] = *order.DeliveryInfo.B2BInfo.ShipmentDetails.ShipmentDestination.City
		} else {
			data["deliveryCity"] = ""
		}

		// Add B2B specific fields
		data["deliveryFee"] = order.DeliveryInfo.B2BInfo.DeliveryFee
		data["estimationTime"] = ""
		data["deliveryTime"] = ""
	} else {
		data["deliveryAddress"] = ""
		data["deliveryCity"] = ""
		data["deliveryFee"] = 0.0
		data["estimationTime"] = ""
		data["deliveryTime"] = ""
	}

	// Add estimation time from pickup info
	if order.PickupInfo != nil && order.PickupInfo.EstimationTime != nil {
		data["estimationTime"] = fmt.Sprintf("%d minutes", *order.PickupInfo.EstimationTime)
	} else if data["estimationTime"] == "" {
		data["estimationTime"] = ""
	}

	// Add branch name from pickup info
	if order.PickupInfo != nil && order.PickupInfo.PickupBranch != nil {
		data["branchName"] = order.PickupInfo.PickupBranch.Name
	} else {
		data["branchName"] = ""
	}

	// Add discount and coupon info
	if order.TotalDiscountValue != nil {
		data["totalDiscountValue"] = *order.TotalDiscountValue
	} else {
		data["totalDiscountValue"] = 0.0
	}

	// Add currency code
	data["currencyCode"] = order.CurrencyCode

	// Add coupon info
	if order.CouponCode != nil {
		data["couponCode"] = *order.CouponCode
	} else {
		data["couponCode"] = ""
	}
	data["couponDiscountValue"] = order.CouponDiscountValue

	// Add note
	if order.Note != nil {
		data["note"] = *order.Note
	} else {
		data["note"] = ""
	}

	// Calculate total paid amount
	totalPaid := 0.0
	if order.PaymentInfo != nil {
		for _, payment := range order.PaymentInfo.Payments {
			if payment.PaymentStatus == order_enums.PAID {
				totalPaid += payment.PaymentAmount
			}
		}
	}
	data["totalPaid"] = totalPaid

	// Process order items
	orderItemsText := ""
	addonTotalPrice := 0.0
	addonsText := ""

	if order.OrderItems != nil {
		for i, item := range order.OrderItems {
			if i > 0 {
				orderItemsText += "; "
			}
			// Use Name field which is Localizable[string] - get default language value
			productName := ""
			if len(item.Name) > 0 {
				for _, value := range item.Name {
					productName = value
					break // Take first available language
				}
			}
			orderItemsText += fmt.Sprintf("%s (Qty: %d)", productName, item.Quantity)

			// Process addons for this item
			if item.Addons != nil {
				for j, addon := range item.Addons {
					if j == 0 && i == 0 {
						// First addon of first item
					} else if j == 0 {
						addonsText += "; "
					} else {
						addonsText += ", "
					}

					// Get addon label (Localizable[string])
					addonLabel := ""
					if len(addon.Label) > 0 {
						for _, value := range addon.Label {
							addonLabel = value
							break // Take first available language
						}
					}

					// Process addon values
					for k, addonValue := range addon.Values {
						if k > 0 {
							addonsText += ", "
						}

						// Get addon value label (Localizable[string])
						addonValueLabel := ""
						if len(addonValue.Label) > 0 {
							for _, value := range addonValue.Label {
								addonValueLabel = value
								break // Take first available language
							}
						}

						addonsText += fmt.Sprintf("%s: %s (Qty: %d)", addonLabel, addonValueLabel, addonValue.Quantity)
					}
				}
			}

			// Add item's addon total price
			addonTotalPrice += item.AddonTotalPrice
		}
	}

	data["orderItems"] = orderItemsText
	data["addons"] = addonsText
	data["addonTotalPrice"] = addonTotalPrice

	return data
}
