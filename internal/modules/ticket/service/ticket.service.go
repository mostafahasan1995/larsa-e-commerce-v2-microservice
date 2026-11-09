package ticket_service

import (
	"context"
	"encoding/json"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	ticket_dto "larsa-e-commerce-microservice-v2/internal/modules/ticket/dto"
	ticket_enums "larsa-e-commerce-microservice-v2/internal/modules/ticket/enums"
	ticket_filters "larsa-e-commerce-microservice-v2/internal/modules/ticket/filters"
	ticket_models "larsa-e-commerce-microservice-v2/internal/modules/ticket/models"
	ticket_utils "larsa-e-commerce-microservice-v2/internal/modules/ticket/utils"

	"git.larsa.io/mahdawi/microservices-commons.git/nats"
)

type TicketRepo interface {
	Create(ctx context.Context, ticket *ticket_models.Ticket) (*ticket_models.TicketResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*ticket_models.TicketResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *ticket_filters.TicketFilters) (*ticket_models.TicketsWithPagination, *helpers.ApiError)
	GetMyTickets(ctx context.Context, skip, limit int64, filter *ticket_filters.TicketFilters) (*ticket_models.TicketsWithPagination, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, ticket *ticket_models.Ticket) (*ticket_models.TicketResponse, *helpers.ApiError)
	Count(ctx context.Context, filter []bson.M) (int64, *helpers.ApiError)
	GetAll(ctx context.Context, filter []bson.M) ([]ticket_models.TicketResponse, *helpers.ApiError)
	GenerateUserTicketPipeline(cfg *types.AppCfg) []bson.M
}

type TicketService interface {
	Create(ctx context.Context, dto ticket_dto.CreateTicketDto) (*ticket_models.TicketResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID) (*ticket_models.TicketResponse, *helpers.ApiError)
	GetMyTickets(ctx context.Context, skip, limit int64, query string) (*ticket_models.TicketsWithPagination, *helpers.ApiError)
	AddComment(ctx context.Context, id primitive.ObjectID, dto ticket_dto.AddCommentDto) (*ticket_models.TicketResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto ticket_dto.CreateTicketDto) (*ticket_models.TicketResponse, *helpers.ApiError)
	GetTicketSummary(ctx context.Context) (*ticket_dto.TicketSummaryResponse, *helpers.ApiError)
}

type ticketService struct {
	ticketRepository  TicketRepo
	ticketTypeService config_services.TicketTypeService
}

func NewTicketService(i *do.Injector) (TicketService, error) {
	return &ticketService{
		ticketRepository: do.MustInvoke[TicketRepo](i),
	}, nil
}

func (s *ticketService) Create(ctx context.Context, dto ticket_dto.CreateTicketDto) (*ticket_models.TicketResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	ticket, err := ticket_utils.ConvertCreateTicketDTOToModel(&dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	ticket.RequesterId = cfg.User.Id
	ticket.Status = ticket_enums.PENDING
	if err != nil {
		return nil, err
	}

	return s.ticketRepository.Create(ctx, &ticket)
}

func (s *ticketService) GetById(ctx context.Context, id primitive.ObjectID) (*ticket_models.TicketResponse, *helpers.ApiError) {
	ticketFilters, err := ticket_filters.NewTicketFilters("")
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}
	return s.ticketRepository.GetById(ctx, id, ticketFilters.BuildPipeline(bson.M{}, ctx, true))
}

func (s *ticketService) GetMyTickets(ctx context.Context, skip, limit int64, query string) (*ticket_models.TicketsWithPagination, *helpers.ApiError) {
	ticketFilters, err := ticket_filters.NewTicketFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "invalid filter format: " + err.Error(),
		}
	}
	return s.ticketRepository.GetMyTickets(ctx, skip, limit, ticketFilters)
}

func (s *ticketService) AddComment(ctx context.Context, id primitive.ObjectID, dto ticket_dto.AddCommentDto) (*ticket_models.TicketResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	comment, err := ticket_utils.ConvertAddCommentDTOToModel(&dto, cfg)
	if err != nil {
		return nil, err
	}

	ticket, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	ticket.Comments = append(ticket.Comments, comment)
	ticketModel := ticket_models.Ticket{
		BaseTicket: ticket.BaseTicket,
	}

	return s.ticketRepository.Update(ctx, id, &ticketModel)
}

func (s *ticketService) Update(ctx context.Context, id primitive.ObjectID, dto ticket_dto.CreateTicketDto) (*ticket_models.TicketResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	existingTicket, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	ticketModel, err := ticket_utils.ConvertCreateTicketDTOToModel(&dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	ticketModel.RequesterId = existingTicket.RequesterId
	ticketModel.Comments = existingTicket.Comments

	if ticketModel.Status != existingTicket.Status {
		s.SendNotification(ctx, &ticketModel, "Ticket status updated", "The status of ticket "+ticketModel.Title+" has been updated to "+string(ticketModel.Status))
	}

	if err != nil {
		return nil, err
	}

	return s.ticketRepository.Update(ctx, id, &ticketModel)
}

func (s *ticketService) GetTicketSummary(ctx context.Context) (*ticket_dto.TicketSummaryResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	userTicketPipeline := s.ticketRepository.GenerateUserTicketPipeline(cfg)

	totalTickets, err := s.ticketRepository.Count(ctx, userTicketPipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}
	openTicketPipeline := append(userTicketPipeline, bson.M{
		"$match": bson.M{
			"status": bson.M{
				"$in": []string{
					string(ticket_enums.PENDING),
					string(ticket_enums.IN_PROGRESS),
				},
			},
		},
	})
	openTickets, err := s.ticketRepository.Count(ctx, openTicketPipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	closedTicketPipeline := append(userTicketPipeline, bson.M{
		"$match": bson.M{
			"status": ticket_enums.CLOSED,
		},
	})
	closedTickets, err := s.ticketRepository.Count(ctx, closedTicketPipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	// Get all closed tickets to calculate resolution time
	closedTicketsList, err := s.ticketRepository.GetAll(ctx, closedTicketPipeline)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	// Calculate average resolution time
	var totalResolutionTime float64
	for _, ticket := range closedTicketsList {
		resolutionTime := ticket.UpdatedAt.Sub(ticket.CreatedAt).Minutes()
		totalResolutionTime += resolutionTime
	}
	avgResolutionTime := 0.0
	if len(closedTicketsList) > 0 {
		avgResolutionTime = totalResolutionTime / float64(len(closedTicketsList))
	}

	return &ticket_dto.TicketSummaryResponse{
		TotalTickets:      totalTickets,
		OpenTickets:       openTickets,
		ClosedTickets:     closedTickets,
		AvgResolutionTime: avgResolutionTime,
	}, nil
}

func (s *ticketService) SendNotification(ctx context.Context, ticketModel *ticket_models.Ticket, title, message string) {
	ticketType, _ := s.ticketTypeService.GetById(ctx, ticketModel.TypeId, "")
	if ticketType.InternalNotification {
		for _, team := range ticketType.TeamNotify {
			// Create email notification message
			emailMsg := map[string]interface{}{
				"to":      team.Email,
				"subject": title,
				"body":    message,
				"headers": map[string]string{
					"xClient":  "ecommerce",
					"xService": "notification",
				},
			}

			// Convert to JSON
			msgBytes, err := json.Marshal(emailMsg)
			if err != nil {
				continue
			}

			// Publish to NATS
			nats.Publish("notification:email", msgBytes)
			nats.PublishJet("notification:email", msgBytes)
		}
	}
}
