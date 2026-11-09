package feature_repos

import (
	"context"
	enums "larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	feature_filters "larsa-e-commerce-microservice-v2/internal/modules/feature/filters"
	feature_models "larsa-e-commerce-microservice-v2/internal/modules/feature/models"
	feature_service "larsa-e-commerce-microservice-v2/internal/modules/feature/service"
	feature_utils "larsa-e-commerce-microservice-v2/internal/modules/feature/utils"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/samber/do"
)

type FeatureRepository struct {
	*repo.BaseRepo[feature_models.Feature, feature_models.FeatureResponse]
}

func NewFeatureRepository(i *do.Injector) (feature_service.FeatureRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &FeatureRepository{
		BaseRepo: repo.NewBaseRepo[feature_models.Feature, feature_models.FeatureResponse](db, enums.FEATURE),
	}, nil
}

func (repo *FeatureRepository) GetById(ctx context.Context, id primitive.ObjectID, pipeline []bson.M) (*feature_models.FeatureResponse, *helpers.ApiError) {
	pipeline = append(pipeline, feature_utils.PopulateFeatureGroup...)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *FeatureRepository) GetWithPagination(ctx context.Context, skip, limit int64, filter feature_filters.FeatureFilters) (*feature_models.FeaturesWithPagination, *helpers.ApiError) {
	pipeline := filter.BuildPipeline([]bson.M{}, ctx, true)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filter.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *FeatureRepository) UpdateFieldsById(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}, userId primitive.ObjectID) (*feature_models.FeatureResponse, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	updateFields["updatedAt"] = time.Now()
	updateFields["updatedBy"] = userId

	result, updateErr := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	if updateErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to update document: " + updateErr.Error(),
		}
	}

	if result.MatchedCount == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Document not found",
		}
	}

	return repo.BaseRepo.GetById(ctx, id, []bson.M{})
}
