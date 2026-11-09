package feature_service

import (
	"context"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	feature_dto "larsa-e-commerce-microservice-v2/internal/modules/feature/dto"
	feature_filters "larsa-e-commerce-microservice-v2/internal/modules/feature/filters"
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	feature_utils "larsa-e-commerce-microservice-v2/internal/modules/feature/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeatureGroupRepo interface {
	Create(ctx context.Context, group feature_models.FeatureGroupModel) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, group feature_models.FeatureGroupModel, filter bson.M) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filter feature_filters.FeatureGroupFilters) (*feature_models.FeatureGroupsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, pipeline []bson.M) ([]feature_models.FeatureGroupResponse, *helpers.ApiError)
}

type FeatureGroupService interface {
	Create(ctx context.Context, dto *feature_dto.CreateFeatureGroupDTO) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *feature_dto.UpdateFeatureGroupDTO) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID) (*feature_models.FeatureGroupResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*feature_models.FeatureGroupsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]feature_models.FeatureGroupResponse, *helpers.ApiError)
}

type featureGroup struct {
	repo           FeatureGroupRepo
	featureService FeatureService
}

func NewFeatureGroupService(i *do.Injector) (FeatureGroupService, error) {
	return &featureGroup{
		repo:           do.MustInvoke[FeatureGroupRepo](i),
		featureService: do.MustInvoke[FeatureService](i),
	}, nil
}

func (s *featureGroup) Create(ctx context.Context, dto *feature_dto.CreateFeatureGroupDTO) (*feature_models.FeatureGroupResponse, *helpers.ApiError) {
	group, apiErr := feature_utils.ConvertCreateGroupDTOToModel(dto, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	result, apiErr := s.repo.Create(ctx, group)
	if apiErr != nil {
		return nil, apiErr
	}

	for _, feature := range dto.Features {
		if feature.Active {
			_, _apiErr := s.featureService.UpdateFieldsById(ctx, feature.Id, map[string]interface{}{
				"groupId": result.Id,
			})
			if _apiErr != nil {
				return nil, _apiErr
			}
		}
	}
	return result, nil
}

func (s *featureGroup) Update(ctx context.Context, id primitive.ObjectID, dto *feature_dto.UpdateFeatureGroupDTO) (*feature_models.FeatureGroupResponse, *helpers.ApiError) {
	group, apiErr := feature_utils.ConvertUpdateGroupDTOToModel(dto, id, ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	for _, feature := range dto.Features {
		if feature.Active {
			_, _apiErr := s.featureService.UpdateFieldsById(ctx, feature.Id, map[string]interface{}{
				"groupId": id,
			})
			if _apiErr != nil {
				return nil, _apiErr
			}
		} else {
			_feature, _ := s.featureService.GetById(ctx, feature.Id)
			if _feature.GroupId == &id {
				_, _apiErr := s.featureService.UpdateFieldsById(ctx, feature.Id, map[string]interface{}{
					"groupId": nil,
				})
				if _apiErr != nil {
					return nil, _apiErr
				}
			}
		}
	}

	result, apiErr := s.repo.UpdateById(ctx, id, group, bson.M{"deletedAt": bson.M{"$exists": false}})
	if apiErr != nil {
		return nil, apiErr
	}
	return result, nil
}

func (s *featureGroup) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return s.repo.DeleteById(ctx, id, bson.M{})
}

func (s *featureGroup) GetById(ctx context.Context, id primitive.ObjectID) (*feature_models.FeatureGroupResponse, *helpers.ApiError) {
	featureGroup, apiErr := s.repo.GetById(ctx, id, []bson.M{})
	if apiErr != nil {
		return nil, apiErr
	}
	return featureGroup, apiErr
}

func (s *featureGroup) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*types.DataWithPagination[feature_models.FeatureGroupResponse], *helpers.ApiError) {
	filter, err := feature_filters.NewFeatureGroupFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}
	return s.repo.GetWithPagination(ctx, skip, limit, *filter)
}

func (s *featureGroup) GetAll(ctx context.Context, query string) ([]feature_models.FeatureGroupResponse, *helpers.ApiError) {
	filter, err := feature_filters.NewFeatureGroupFilters(query)
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
