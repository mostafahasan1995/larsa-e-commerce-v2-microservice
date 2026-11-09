package ticket_utils

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	ticket_dto "larsa-e-commerce-microservice-v2/internal/modules/ticket/dto"
	ticket_models "larsa-e-commerce-microservice-v2/internal/modules/ticket/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateTicketDTOToModel(dto *ticket_dto.CreateTicketDto, additionalOptions types.AdditionalOptions) (ticket_models.Ticket, *helpers.ApiError) {
	ticket := ticket_models.Ticket{
		BaseTicket: ticket_models.BaseTicket{
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
			Title:             dto.Title,
			TypeId:            dto.TypeId,
			Priority:          dto.Priority,
			Description:       dto.Description,
			Attachment:        dto.Attachment,
			Status:            dto.Status,

			QuotationId: dto.QuotationId,
			InvoiceId:   dto.InvoiceId,
			OrderId:     dto.OrderId,

			RequesterId: additionalOptions.UserId,
			BuyerId:     dto.BuyerId,
			VendorId:    dto.VendorId,

			Comments: []ticket_models.Comment{},
		},
	}

	return ticket, nil
}

func ConvertAddCommentDTOToModel(dto *ticket_dto.AddCommentDto, cfg *types.AppCfg) (ticket_models.Comment, *helpers.ApiError) {
	now := time.Now()
	userId := cfg.User.Id

	comment := ticket_models.Comment{
		Text:       dto.Text,
		Attachment: dto.Attachment,
		CreatedAt:  primitive.NewDateTimeFromTime(now),
		UserId:     userId,
	}

	return comment, nil
}
