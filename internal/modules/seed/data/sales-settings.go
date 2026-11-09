package seed_data

import (
	settings_sales_enums "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/enums"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"
)

var DefaultOrderStatuses = []settings_sales_models.OrderStatus{}

func InitDefaultOrderStatuses() {
	pending := settings_sales_enums.SYSTEM_STATUS_TYPE_PENDING
	canceled := settings_sales_enums.SYSTEM_STATUS_TYPE_CANCELED
	failed := settings_sales_enums.SYSTEM_STATUS_TYPE_FAILED
	refund := settings_sales_enums.SYSTEM_STATUS_TYPE_REFUND
	completed := settings_sales_enums.SYSTEM_STATUS_TYPE_COMPLETED
	DefaultOrderStatuses = []settings_sales_models.OrderStatus{
		{
			Title:                        "pending",
			Color:                        "#FFA500", // Orange
			SystemStatusType:             &pending,
			Order:                        1,
			IsSystemStatus:               true,
			MoveForward:                  true,
			MoveForwardOnlyIfPaymentDone: false,
			MoveBackward:                 false,
			CancelOrder:                  false,
			RefundOrder:                  false,
			AssignToDriver:               false,
			IsActive:                     true,
		},
		{
			Title:                        "canceled",
			Color:                        "#FF0000", // Red
			SystemStatusType:             &canceled,
			Order:                        2,
			IsSystemStatus:               true,
			MoveForward:                  false,
			MoveForwardOnlyIfPaymentDone: false,
			MoveBackward:                 false,
			CancelOrder:                  false,
			RefundOrder:                  false,
			AssignToDriver:               false,
			IsActive:                     true,
		},
		{
			Title:                        "failed",
			Color:                        "#8B0000", // Dark Red
			SystemStatusType:             &failed,
			Order:                        3,
			IsSystemStatus:               true,
			MoveForward:                  false,
			MoveForwardOnlyIfPaymentDone: false,
			MoveBackward:                 false,
			CancelOrder:                  false,
			RefundOrder:                  false,
			AssignToDriver:               false,
			IsActive:                     true,
		},
		{
			Title:                        "refund",
			Color:                        "#8B0000", // Dark Red
			SystemStatusType:             &refund,
			Order:                        4,
			IsSystemStatus:               true,
			MoveForward:                  false,
			MoveForwardOnlyIfPaymentDone: false,
			MoveBackward:                 false,
			CancelOrder:                  false,
			RefundOrder:                  false,
			AssignToDriver:               false,
			IsActive:                     true,
		},
		{
			Title:                        "completed",
			Color:                        "#008000", // Green
			SystemStatusType:             &completed,
			Order:                        5,
			IsSystemStatus:               true,
			MoveForward:                  false,
			MoveForwardOnlyIfPaymentDone: false,
			MoveBackward:                 false,
			CancelOrder:                  false,
			RefundOrder:                  false,
			AssignToDriver:               false,
			IsActive:                     true,
		},
	}
}
