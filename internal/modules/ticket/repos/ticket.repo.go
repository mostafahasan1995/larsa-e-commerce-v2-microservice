package ticket_repos

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	ticket_filters "larsa-e-commerce-microservice-v2/internal/modules/ticket/filters"
	ticket_models "larsa-e-commerce-microservice-v2/internal/modules/ticket/models"
	ticket_service "larsa-e-commerce-microservice-v2/internal/modules/ticket/service"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketRepo struct {
	*base_repos.BaseRepo[ticket_models.Ticket, ticket_models.TicketResponse]
}

func NewTicketRepository(i *do.Injector) (ticket_service.TicketRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ticketRepo{
		BaseRepo: base_repos.NewBaseRepo[ticket_models.Ticket, ticket_models.TicketResponse](db, enums.TICKET),
	}, nil
}

func (repo *ticketRepo) Create(ctx context.Context, ticket *ticket_models.Ticket) (*ticket_models.TicketResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, *ticket)
}

func (repo *ticketRepo) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*ticket_models.TicketResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *ticketRepo) GetAll(ctx context.Context, filter []bson.M) ([]ticket_models.TicketResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetAll(ctx, filter)
}

func (repo *ticketRepo) GetWithPagination(ctx context.Context, skip, limit int64, filter *ticket_filters.TicketFilters) (*ticket_models.TicketsWithPagination, *helpers.ApiError) {
	pipeline := filter.BuildPipeline(bson.M{}, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *ticketRepo) GetMyTickets(ctx context.Context, skip, limit int64, filter *ticket_filters.TicketFilters) (*ticket_models.TicketsWithPagination, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	pipeline := repo.GenerateUserTicketPipeline(cfg)
	pipeline = append(pipeline, filter.BuildPipeline(bson.M{}, ctx, true)...)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *ticketRepo) Update(ctx context.Context, id primitive.ObjectID, ticket *ticket_models.Ticket) (*ticket_models.TicketResponse, *helpers.ApiError) {
	return repo.BaseRepo.UpdateById(ctx, id, *ticket, bson.M{})
}

func (repo *ticketRepo) Count(ctx context.Context, filter []bson.M) (int64, *helpers.ApiError) {
	collection, apiErr := repo.BaseRepo.GetCollection(ctx)
	if apiErr != nil {
		return 0, apiErr
	}

	// Convert pipeline to a single match stage for count
	matchStage := bson.M{}
	for _, stage := range filter {
		if m, ok := stage["$match"]; ok {
			if matchStage["$or"] != nil {
				// If we already have an $or condition, merge with existing conditions
				existingOr := matchStage["$or"].([]bson.M)
				if newOr, ok := m.(bson.M)["$or"]; ok {
					matchStage["$or"] = append(existingOr, newOr.([]bson.M)...)
				} else {
					matchStage["$or"] = append(existingOr, m.(bson.M))
				}
			} else {
				// First match stage
				matchStage = m.(bson.M)
			}
		}
	}

	count, err := collection.CountDocuments(ctx, matchStage)
	if err != nil {
		return 0, helpers.DefaultError(err)
	}
	return count, nil
}

// Util function to generate pipeline for user ticket
func (repo *ticketRepo) GenerateUserTicketPipeline(cfg *types.AppCfg) []bson.M {
	if cfg.User.Type == enums.SuperAdmin {
		return []bson.M{}
	}

	conditions := []bson.M{
		{"requesterId": cfg.User.Id}, // User created the ticket
	}

	switch cfg.User.Type {
	case enums.Vendor:
		conditions = append(conditions, bson.M{
			"vendorId": cfg.User.Id,
		})
	case enums.Member, enums.Normal:
		conditions = append(conditions, bson.M{
			"buyerId": cfg.User.Id,
		})
	}

	return []bson.M{
		{
			"$match": bson.M{
				"$or": conditions,
			},
		},
	}
}
