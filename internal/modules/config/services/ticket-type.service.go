package config_services

import (
	"context"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	config_dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils"
)

type TicketTypeService interface {
	Create(ctx context.Context, dto *config_dtos.CreateTicketTypeDto) (*config_models.TicketTypeResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_models.TicketTypeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]config_models.TicketTypeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[config_models.TicketTypeResponse], *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, dto *config_dtos.CreateTicketTypeDto) (*config_models.TicketTypeResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
}

type TicketTypeRepo interface {
	Create(ctx context.Context, model config_models.TicketType) (*config_models.TicketTypeResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_models.TicketTypeResponse, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]config_models.TicketTypeResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.TicketTypeFilter) (*types.DataWithPagination[config_models.TicketTypeResponse], *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, model *config_models.TicketType) (*config_models.TicketTypeResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
}

type ticketTypeService struct {
	ticketTypeRepo TicketTypeRepo
}

func NewTicketTypeService(i *do.Injector) (TicketTypeService, error) {
	return &ticketTypeService{
		ticketTypeRepo: do.MustInvoke[TicketTypeRepo](i),
	}, nil
}

func (s *ticketTypeService) Create(ctx context.Context, dto *config_dtos.CreateTicketTypeDto) (*config_models.TicketTypeResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	ticketType := config_utils.ConvertToTicketTypeModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return s.ticketTypeRepo.Create(ctx, *ticketType)
}

func (s *ticketTypeService) GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_models.TicketTypeResponse, *helpers.ApiError) {
	filters, err := config_filters.NewTicketTypeFilter(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	pipeline := filters.BuildPipeline(bson.M{}, ctx, true)
	return s.ticketTypeRepo.GetById(ctx, id, pipeline)
}

func (s *ticketTypeService) GetAll(ctx context.Context, query string) ([]config_models.TicketTypeResponse, *helpers.ApiError) {
	filters, err := config_filters.NewTicketTypeFilter(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	pipeline := filters.BuildPipeline(bson.M{}, ctx, false)
	result, apiErr := s.ticketTypeRepo.GetAll(ctx, pipeline)
	if apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

func (s *ticketTypeService) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[config_models.TicketTypeResponse], *helpers.ApiError) {
	filters, err := config_filters.NewTicketTypeFilter(query)
	if err != nil {
		apiErr := helpers.NewApiError(400, err)
		return nil, &apiErr
	}

	return s.ticketTypeRepo.GetWithPagination(ctx, skip, limit, filters)
}

func (s *ticketTypeService) UpdateById(ctx context.Context, id primitive.ObjectID, dto *config_dtos.CreateTicketTypeDto) (*config_models.TicketTypeResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	ticketType := config_utils.ConvertToTicketTypeModel(dto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return s.ticketTypeRepo.UpdateById(ctx, id, ticketType)
}

func (s *ticketTypeService) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.ticketTypeRepo.DeleteById(ctx, id)
}
