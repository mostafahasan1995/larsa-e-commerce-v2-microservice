package feature_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	feature_dto "larsa-e-commerce-microservice-v2/internal/modules/feature/dto"
	feature_filters "larsa-e-commerce-microservice-v2/internal/modules/feature/filters"
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	feature_utils "larsa-e-commerce-microservice-v2/internal/modules/feature/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeatureRepo interface {
	Create(ctx context.Context, feature feature_models.Feature) (*feature_models.FeatureResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, feature feature_models.Feature, filter bson.M) (*feature_models.FeatureResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*feature_models.FeatureResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter feature_filters.FeatureFilters) (*feature_models.FeaturesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]feature_models.FeatureResponse, *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}, userId primitive.ObjectID) (*feature_models.FeatureResponse, *helpers.ApiError)
}

type FeatureService interface {
	Create(ctx context.Context, dto *feature_dto.CreateFeatureDTO) (*feature_models.FeatureResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *feature_dto.UpdateFeatureDTO) (*feature_models.FeatureResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID) (*feature_models.FeatureResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*feature_models.FeaturesWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]feature_models.FeatureResponse, *helpers.ApiError)
	UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*feature_models.FeatureResponse, *helpers.ApiError)
}

type feature struct {
	repo FeatureRepo
}

func NewFeatureService(i *do.Injector) (FeatureService, error) {
	return &feature{
		repo: do.MustInvoke[FeatureRepo](i),
	}, nil
}

func (s *feature) Create(ctx context.Context, dto *feature_dto.CreateFeatureDTO) (*feature_models.FeatureResponse, *helpers.ApiError) {
	config, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}
	feature, apiErr := feature_utils.ConvertCreateDTOToModel(dto, config.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}
	result, apiErr := s.repo.Create(ctx, feature)
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *feature) Update(ctx context.Context, id primitive.ObjectID, dto *feature_dto.UpdateFeatureDTO) (*feature_models.FeatureResponse, *helpers.ApiError) {
	config, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}
	feature, apiErr := feature_utils.ConvertUpdateDTOToModel(dto, id, config.User.Id)
	if apiErr != nil {
		return nil, apiErr
	}

	result, apiErr := s.repo.UpdateById(ctx, id, feature, bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *feature) GetById(ctx context.Context, id primitive.ObjectID) (*feature_models.FeatureResponse, *helpers.ApiError) {
	return s.repo.GetById(ctx, id, []bson.M{})
}

func (s *feature) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id, bson.M{})
}

func (s *feature) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*feature_models.FeaturesWithPagination, *helpers.ApiError) {
	filter, err := feature_filters.NewFeatureFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *feature) GetAll(ctx context.Context, query string) ([]feature_models.FeatureResponse, *helpers.ApiError) {
	filter, err := feature_filters.NewFeatureFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	pipeline := []bson.M{}
	pipeline = filter.BuildPipeline(pipeline, ctx, false)

	return s.repo.GetAll(ctx, pipeline)
}

func (s *feature) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*feature_models.FeatureResponse, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to get request app config",
		}
	}

	return s.repo.UpdateFieldsById(ctx, id, updates, cfg.User.Id)
}
