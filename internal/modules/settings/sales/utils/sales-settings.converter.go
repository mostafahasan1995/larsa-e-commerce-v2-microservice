package settings_sales_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	settings_sales_dtos "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/dtos"
	settings_sales_models "larsa-e-commerce-microservice-v2/internal/modules/settings/sales/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertUpdateDTOToSalesSettingsModel(
	ctx context.Context,
	dto settings_sales_dtos.UpdateSalesSettingsDTO,
	existingModel *settings_sales_models.SalesSettings,
) (*settings_sales_models.SalesSettings, *helpers.ApiError) {
	now := time.Now()
	cfg, _ := util.GetReqAppCfg(ctx)
	userId := cfg.User.Id

	modelToUpdate := settings_sales_models.SalesSettings{}
	if existingModel != nil {
		modelToUpdate = *existingModel
		modelToUpdate.UpdatedBy = &userId
		modelToUpdate.UpdatedAt = now
	} else {
		modelToUpdate.Id = primitive.NewObjectID()
		modelToUpdate.CreatedBy = &userId
		modelToUpdate.CreatedAt = now
		modelToUpdate.UpdatedBy = &userId
		modelToUpdate.UpdatedAt = now
	}

	modelToUpdate.OrderStatuses = make([]settings_sales_models.OrderStatus, len(dto.OrderStatuses))
	for i, osDto := range dto.OrderStatuses {
		modelToUpdate.OrderStatuses[i] = settings_sales_models.OrderStatus{
			Title:                        osDto.Title,
			Color:                        osDto.Color,
			Order:                        osDto.Order,
			IsSystemStatus:               osDto.IsSystemStatus,
			SystemStatusType:             osDto.SystemStatusType,
			AppliesTo:                    osDto.AppliesTo,
			MoveForward:                  osDto.MoveForward,
			MoveForwardOnlyIfPaymentDone: osDto.MoveForwardOnlyIfPaymentDone,
			MoveBackward:                 osDto.MoveBackward,
			CancelOrder:                  osDto.CancelOrder,
			RefundOrder:                  osDto.RefundOrder,
			AssignToDriver:               osDto.AssignToDriver,
			ReopenOrder:                  osDto.ReopenOrder,
			ConfirmDelivery:              osDto.ConfirmDelivery,
			ConfirmDeliveryTemplate:      osDto.ConfirmDeliveryTemplate,
			IsActive:                     osDto.IsActive,
		}
	}

	if dto.ShowVendorInfoInOrderDetails != nil {
		modelToUpdate.ShowVendorInfoInOrderDetails = *dto.ShowVendorInfoInOrderDetails
	}

	if dto.LogisticSettings != nil {
		modelToUpdate.LogisticSettings.EnablePickup = *dto.LogisticSettings.EnablePickup
		modelToUpdate.LogisticSettings.EnableDelivery = *dto.LogisticSettings.EnableDelivery
		if dto.LogisticSettings.Pickup != nil {
			modelToUpdate.LogisticSettings.Pickup = &settings_sales_models.LogisticPickup{
				BranchName: *dto.LogisticSettings.Pickup.BranchName,
				PickUpDate: *dto.LogisticSettings.Pickup.PickUpDate,
				PickUpTime: *dto.LogisticSettings.Pickup.PickUpTime,
			}
		}
		if dto.LogisticSettings.Delivery != nil {
			modelToUpdate.LogisticSettings.Delivery = &settings_sales_models.LogisticDelivery{
				ShipmentMethod:   *dto.LogisticSettings.Delivery.ShipmentMethod,
				TrackingNumber:   *dto.LogisticSettings.Delivery.TrackingNumber,
				TransportCompany: *dto.LogisticSettings.Delivery.TransportCompany,
				ShipmentDate:     *dto.LogisticSettings.Delivery.ShipmentDate,
				DeliveryFee:      *dto.LogisticSettings.Delivery.DeliveryFee,
				CustomAddress:    *dto.LogisticSettings.Delivery.CustomAddress,
			}
		}
	}

	if dto.RecentActivity != nil {
		modelToUpdate.RecentActivity.EnableRecentActivity = *dto.RecentActivity.EnableRecentActivity
		modelToUpdate.RecentActivity.NumberOfProcessedOrders = *dto.RecentActivity.NumberOfProcessedOrders
		modelToUpdate.RecentActivity.NumberOfIncompleteOrders = *dto.RecentActivity.NumberOfIncompleteOrders
		modelToUpdate.RecentActivity.NumberOfCustomers = *dto.RecentActivity.NumberOfCustomers
		modelToUpdate.RecentActivity.UnpaidAmount = *dto.RecentActivity.UnpaidAmount
		modelToUpdate.RecentActivity.OverdueAmount = *dto.RecentActivity.OverdueAmount
	}

	if dto.InvoiceTerms != nil {
		if dto.InvoiceTerms.DefaultInvoiceTerms != nil {
			modelToUpdate.InvoiceTerms.DefaultInvoiceTerms = *dto.InvoiceTerms.DefaultInvoiceTerms
		}
		if dto.InvoiceTerms.DefaultInvoiceTermsUnit != nil {
			modelToUpdate.InvoiceTerms.DefaultInvoiceTermsUnit = *dto.InvoiceTerms.DefaultInvoiceTermsUnit
		}
	}

	if dto.EnableDiscountOnSales != nil {
		modelToUpdate.EnableDiscountOnSales = *dto.EnableDiscountOnSales
	}
	if dto.EnablePartialPayment != nil {
		modelToUpdate.EnablePartialPayment = *dto.EnablePartialPayment
	}
	if dto.SalesTaxesRate != nil {
		modelToUpdate.SalesTaxesRate = *dto.SalesTaxesRate
	}
	if dto.EnableSalesNotes != nil {
		modelToUpdate.EnableSalesNotes = *dto.EnableSalesNotes
	}
	if dto.EnableAutoGenerateInvoiceNumber != nil {
		modelToUpdate.EnableAutoGenerateInvoiceNumber = *dto.EnableAutoGenerateInvoiceNumber
	}
	if dto.EnableInvoiceEmailNotification != nil {
		modelToUpdate.EnableInvoiceEmailNotification = *dto.EnableInvoiceEmailNotification
	}

	// Flow mode
	if dto.OrderFlowMode != nil {
		modelToUpdate.OrderFlowMode = *dto.OrderFlowMode
	}

	return &modelToUpdate, nil
}
