package ticket_dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	ticket_enums "larsa-e-commerce-microservice-v2/internal/modules/ticket/enums"
)

type CreateTicketDto struct {
	Title       string                      `json:"title" validate:"required"`
	TypeId      primitive.ObjectID          `json:"typeId" validate:"required"`
	Priority    ticket_enums.TicketPriority `json:"priority" validate:"required"`
	Description string                      `json:"description" validate:"required"`

	Attachment *types.FileFieldDto `json:"attachment,omitempty"`

	QuotationId *primitive.ObjectID       `json:"quotationId,omitempty"`
	InvoiceId   *primitive.ObjectID       `json:"invoiceId,omitempty"`
	OrderId     *primitive.ObjectID       `json:"orderId,omitempty"`
	Status      ticket_enums.TicketStatus `json:"status"`

	BuyerId  *primitive.ObjectID `json:"buyerId,omitempty"`
	VendorId *primitive.ObjectID `json:"vendorId,omitempty"`
}

type AddCommentDto struct {
	Text       string              `json:"text" validate:"required"`
	Attachment *types.FileFieldDto `json:"attachment,omitempty"`
}

type TicketSummaryResponse struct {
	TotalTickets      int64   `json:"totalTickets"`
	OpenTickets       int64   `json:"openTickets"`
	ClosedTickets     int64   `json:"closedTickets"`
	AvgResolutionTime float64 `json:"avgResolutionTime"` // in minutes
}

func (dto *CreateTicketDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	ticketTypeRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.TICKET_TYPE)
	_, err := ticketTypeRepo.GetById(args.Context, dto.TypeId, nil)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid ticket type",
		}
	}

	userRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.PROFILE)
	if dto.BuyerId != nil {
		_, err = userRepo.GetById(args.Context, *dto.BuyerId, nil)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Invalid buyer",
			}
		}
	}

	if dto.VendorId != nil {
		_, err = userRepo.GetById(args.Context, *dto.VendorId, nil)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Invalid vendor",
			}
		}
	}

	if dto.QuotationId != nil {
		quotationRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.QUOTATION)
		_, err = quotationRepo.GetById(args.Context, *dto.QuotationId, nil)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Invalid quotation",
			}
		}
	}

	if dto.InvoiceId != nil {
		invoiceRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.INVOICE)
		_, err = invoiceRepo.GetById(args.Context, *dto.InvoiceId, nil)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Invalid invoice",
			}
		}
	}

	if dto.OrderId != nil {
		orderRepo := base_repos.NewBaseRepo[any, any](args.DB, enums.ORDER)
		_, err = orderRepo.GetById(args.Context, *dto.OrderId, nil)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    "Invalid order",
			}
		}
	}

	return nil, nil
}
