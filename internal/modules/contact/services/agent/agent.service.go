package contact_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/contact/dtos/agent"
	agent_filters "larsa-e-commerce-microservice-v2/internal/modules/contact/filters/agent"
	models "larsa-e-commerce-microservice-v2/internal/modules/contact/models/agent"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentRepo interface {
	Create(ctx context.Context, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, filter *agent_filters.AgentFilters) (*models.AgentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *agent_filters.AgentFilters) (*models.AgentsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type AgentService interface {
	Create(ctx context.Context, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AgentResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AgentsWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
}

type agent struct {
	repo AgentRepo
}

func NewAgentService(i *do.Injector) (AgentService, error) {
	return &agent{
		repo: do.MustInvoke[AgentRepo](i),
	}, nil
}

func (p *agent) Create(ctx context.Context, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError) {
	return p.repo.Create(ctx, agent)
}

func (p *agent) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AgentResponse, *helpers.ApiError) {
	filter, err := agent_filters.NewAgentFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}

	return p.repo.GetById(ctx, id, filter)
}

func (p *agent) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AgentsWithPagination, *helpers.ApiError) {
	filter, err := agent_filters.NewAgentFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	return p.repo.GetWithPagination(ctx, skip, limit, filter)
}

func (p *agent) UpdateById(ctx context.Context, id primitive.ObjectID, agent *dtos.CreateAgentDto) (*models.AgentResponse, *helpers.ApiError) {
	return p.repo.UpdateById(ctx, id, agent)
}

func (p *agent) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *agent) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}
