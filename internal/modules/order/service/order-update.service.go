package order_service

import (
	"context"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_dto "larsa-e-commerce-microservice-v2/internal/modules/order/dto"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	order_models "larsa-e-commerce-microservice-v2/internal/modules/order/models"
	order_repos "larsa-e-commerce-microservice-v2/internal/modules/order/repos"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderUpdateService interface {
	UpdateOrder(ctx context.Context, id primitive.ObjectID, dto *order_dto.UpdateOrderDTO) (*order_models.OrderResponse, *helpers.ApiError)
}

type orderUpdateService struct {
	repo order_repos.OrderRepo
}

func NewOrderUpdateService(i *do.Injector) (OrderUpdateService, error) {
	repo := do.MustInvoke[order_repos.OrderRepo](i)
	return &orderUpdateService{
		repo: repo,
	}, nil
}

func (s *orderUpdateService) UpdateOrder(ctx context.Context, id primitive.ObjectID, dto *order_dto.UpdateOrderDTO) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Get current order
	currentOrder, apiErr := s.repo.GetById(ctx, id, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	updates, apiErr := s.convertDTOToUpdates(ctx, dto, currentOrder)
	if apiErr != nil {
		return nil, apiErr
	}

	return s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)
}

func (s *orderUpdateService) convertDTOToUpdates(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse) (map[string]interface{}, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	updates := make(map[string]interface{})

	// Handle document update
	if dto.Document != nil {
		if err := s.applyDocumentUpdate(ctx, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	// Handle comment update
	if dto.Comment != nil {
		if err := s.applyCommentUpdate(ctx, cfg, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	// Handle participant info update
	if dto.ParticipantInfo != nil {
		if err := s.applyParticipantUpdate(ctx, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	// Handle payment update
	if dto.PaymentInfo != nil {
		if err := s.applyPaymentUpdate(ctx, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	// Handle pickup info update
	if dto.PickupInfo != nil {
		if err := s.applyPickupUpdate(ctx, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	// Handle delivery info update
	if dto.DeliveryInfo != nil {
		if err := s.applyDeliveryUpdate(ctx, dto, currentOrder, updates); err != nil {
			return nil, err
		}
	}

	return updates, nil
}

// applyDocumentUpdate handles document updates and log history
func (s *orderUpdateService) applyDocumentUpdate(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	document := order_models.Document{
		File:  dto.Document.File,
		Title: dto.Document.Title,
	}
	updates["documents"] = append(currentOrder.Documents, document)
	description := "Document added: " + dto.Document.Title
	_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.DOCUMENT_ADDED, description)
	if apiErr != nil {
		return apiErr
	}
	return nil
}

// applyCommentUpdate handles comment updates and log history
func (s *orderUpdateService) applyCommentUpdate(ctx context.Context, cfg *types.AppCfg, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	userName := ""
	if cfg.User.Type == enums.Vendor {
		userName = currentOrder.VendorInfo.FullName
	} else if cfg.User.Type == enums.Buyer {
		userName = currentOrder.BuyerInfo.FullName
	}

	comment := order_models.Comment{
		Comment:  dto.Comment.Comment,
		Date:     time.Now(),
		UserName: userName,
		UserType: cfg.User.Type,
		UserId:   cfg.User.Id,
	}
	updates["comments"] = append(currentOrder.Comments, comment)
	_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.COMMENT_ADDED, "Comment added")
	if apiErr != nil {
		return apiErr
	}
	return nil
}

// applyParticipantUpdate handles participant updates and related log history
func (s *orderUpdateService) applyParticipantUpdate(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	participantInfo := order_utils.MapParticipantInfo(dto.ParticipantInfo)
	if dto.ParticipantInfo.Type == order_enums.BUYER {
		updates["buyerInfo"] = participantInfo
		_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.STATUS_CHANGED, "Buyer information updated")
		if apiErr != nil {
			return apiErr
		}
	} else if dto.ParticipantInfo.Type == order_enums.VENDOR {
		updates["vendorInfo"] = participantInfo
		_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.STATUS_CHANGED, "Vendor information updated")
		if apiErr != nil {
			return apiErr
		}
	}
	return nil
}

// applyPaymentUpdate handles payment updates and status calculations
func (s *orderUpdateService) applyPaymentUpdate(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	// Process all payments
	for _, paymentDTO := range dto.PaymentInfo.Payments {
		payment := order_models.Payments{
			PaymentMethod: paymentDTO.PaymentMethod,
			PaymentStatus: paymentDTO.PaymentStatus,
			PaymentDate:   paymentDTO.PaymentDate,
			PaymentAmount: paymentDTO.Amount,
			Attachment:    paymentDTO.Attachment,
			PaidTo:        paymentDTO.PaidTo,
		}

		if paymentDTO.PaymentNumber != nil {
			payment.PaymentNumber = *paymentDTO.PaymentNumber
			for i, p := range currentOrder.PaymentInfo.Payments {
				if p.PaymentNumber == *paymentDTO.PaymentNumber {
					currentOrder.PaymentInfo.Payments[i] = payment
					updates["paymentInfo.payments"] = currentOrder.PaymentInfo.Payments
					description := "Updated payment: " + *paymentDTO.PaymentNumber
					_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.PAYMENT_UPDATED, description)
					if apiErr != nil {
						return apiErr
					}
					break
				}
			}
		} else {
			paymentNumber := "PMT-" + primitive.NewObjectID().Hex()[:8]
			payment.PaymentNumber = paymentNumber
			updates["paymentInfo.payments"] = append(currentOrder.PaymentInfo.Payments, payment)
			description := "Payment added: " + paymentNumber + " (" + string(payment.PaymentMethod) + ") for amount " + fmt.Sprintf("%.2f", payment.PaymentAmount)
			_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.PAYMENT_ADDED, description)
			if apiErr != nil {
				return apiErr
			}
		}
		//TODO: add invoice to payment
	}

	// Calculate total paid amount and update payment status automatically
	var updatedPayments []order_models.Payments
	if paymentsUpdate, exists := updates["paymentInfo.payments"]; exists {
		if payments, ok := paymentsUpdate.([]order_models.Payments); ok {
			updatedPayments = payments
		} else {
			updatedPayments = currentOrder.PaymentInfo.Payments
		}
	} else {
		updatedPayments = currentOrder.PaymentInfo.Payments
	}

	totalPaidAmount := 0.0
	for _, payment := range updatedPayments {
		if payment.PaymentStatus == order_enums.PAID {
			totalPaidAmount += payment.PaymentAmount
		}
	}

	var newPaymentStatus order_enums.PaymentStatus
	if totalPaidAmount >= currentOrder.TotalPrice {
		newPaymentStatus = order_enums.PAID
	} else if totalPaidAmount > 0 {
		newPaymentStatus = order_enums.PARTIALLY_PAID
	} else {
		newPaymentStatus = order_enums.UNPAID
	}

	if currentOrder.PaymentInfo.Status != newPaymentStatus {
		updates["paymentInfo.status"] = newPaymentStatus
		description := fmt.Sprintf("Payment status automatically updated to %s (Total paid: %.2f / Total price: %.2f)",
			string(newPaymentStatus), totalPaidAmount, currentOrder.TotalPrice)
		_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.PAYMENT_UPDATED, description)
		if apiErr != nil {
			return apiErr
		}

	}

	if dto.PaymentInfo.Status != "" && currentOrder.PaymentInfo.Status != dto.PaymentInfo.Status {
		updates["paymentInfo.status"] = dto.PaymentInfo.Status
		description := fmt.Sprintf("Payment status manually updated from %s to %s",
			string(currentOrder.PaymentInfo.Status), string(dto.PaymentInfo.Status))
		_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.PAYMENT_UPDATED, description)
		if apiErr != nil {
			return apiErr
		}

	}
	return nil
}

func (s *orderUpdateService) applyPickupUpdate(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	updates["pickupInfo"] = order_utils.MapPickupInfo(dto.PickupInfo)
	_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.STATUS_CHANGED, "Pickup information updated")
	if apiErr != nil {
		return apiErr
	}
	return nil
}

func (s *orderUpdateService) applyDeliveryUpdate(ctx context.Context, dto *order_dto.UpdateOrderDTO, currentOrder *order_models.OrderResponse, updates map[string]interface{}) *helpers.ApiError {
	updates["deliveryInfo"] = order_utils.MapDeliveryInfo(dto.DeliveryInfo, currentOrder.Type)
	_, apiErr := s.addLogHistory(ctx, currentOrder.Id, order_enums.SHIPMENT_UPDATED, "Delivery information updated")
	if apiErr != nil {
		return apiErr
	}
	return nil
}

func (s *orderUpdateService) addLogHistory(ctx context.Context, orderId primitive.ObjectID, action order_enums.LogActionType, description string) (*order_models.OrderResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	now := time.Now()
	logEntry := order_models.LogHistory{
		Action:      action,
		Description: description,
		Date:        now,
	}

	order, err := s.repo.GetById(ctx, orderId, nil)
	if err != nil {
		return nil, err
	}

	var updates map[string]interface{}

	if order.LogHistory == nil {
		updates = map[string]interface{}{
			"$set": bson.M{
				"logHistory": []order_models.LogHistory{logEntry},
			},
		}
	} else {
		updates = map[string]interface{}{
			"$push": bson.M{
				"logHistory": logEntry,
			},
		}
	}

	return s.repo.UpdateFieldsById(ctx, orderId, updates, cfg.User.Id)
}
