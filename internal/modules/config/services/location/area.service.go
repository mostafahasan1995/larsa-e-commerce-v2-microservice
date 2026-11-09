package config_location_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/config/dtos/location"
	config_filters "larsa-e-commerce-microservice-v2/internal/modules/config/filters/location"
	config_location_models "larsa-e-commerce-microservice-v2/internal/modules/config/models/location"
	config_utils "larsa-e-commerce-microservice-v2/internal/modules/config/utils/location"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AreaRepo interface {
	CreateArea(ctx context.Context, area config_location_models.Area) (*config_location_models.AreaResponse, *helpers.ApiError)
	UpdateArea(ctx context.Context, updateArea config_location_models.Area) (*config_location_models.AreaResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*config_location_models.AreaResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter *config_filters.AreaFilters) (*config_location_models.AreasWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]config_location_models.AreaResponse, *helpers.ApiError)
	GetAllSeedsAreas(ctx context.Context, cityName string) ([]config_location_models.Area, *helpers.ApiError)
}

type AreaService interface {
	Create(ctx context.Context, area *dtos.CreateAreaDto) (*config_location_models.AreaResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.AreaResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]config_location_models.AreaResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.AreasWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, area *dtos.CreateAreaDto) (*config_location_models.AreaResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	//Seeds
	GetAllSeedsAreas(ctx context.Context, areaName string) ([]config_location_models.Area, *helpers.ApiError)
}

type area struct {
	repo AreaRepo
}

func NewAreaService(i *do.Injector) (AreaService, error) {
	return &area{
		repo: do.MustInvoke[AreaRepo](i),
	}, nil
}

func (p *area) Create(ctx context.Context, area *dtos.CreateAreaDto) (*config_location_models.AreaResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	areaModel := config_utils.ConvertToAreaModel(*area, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	return p.repo.CreateArea(ctx, areaModel)
}

func (p *area) GetById(ctx context.Context, id primitive.ObjectID, query string) (*config_location_models.AreaResponse, *helpers.ApiError) {
	areaFilters, err := config_filters.NewAreaFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetById(ctx, id, areaFilters.BuildPipeline(bson.M{}, ctx, true))
}

func (p *area) GetAll(ctx context.Context, query string) ([]config_location_models.AreaResponse, *helpers.ApiError) {
	areaFilters, err := config_filters.NewAreaFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetAll(ctx, areaFilters.BuildPipeline(bson.M{}, ctx, false))
}

func (p *area) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*config_location_models.AreasWithPagination, *helpers.ApiError) {
	areaFilters, err := config_filters.NewAreaFilters(query)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	return p.repo.GetWithPagination(ctx, skip, limit, areaFilters)
}

func (p *area) UpdateById(ctx context.Context, id primitive.ObjectID, area *dtos.CreateAreaDto) (*config_location_models.AreaResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	areaModel := config_utils.ConvertToAreaModel(*area, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})
	areaModel.Id = id

	return p.repo.UpdateArea(ctx, areaModel)
}

func (p *area) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id, bson.M{})
}

func (p *area) GetAllSeedsAreas(ctx context.Context, cityName string) ([]config_location_models.Area, *helpers.ApiError) {
	areas, err := p.repo.GetAllSeedsAreas(ctx, cityName)
	if err != nil {
		return nil, err
	}

	return areas, nil
}
