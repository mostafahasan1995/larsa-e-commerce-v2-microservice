package survey_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	repo "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackRepository struct {
	*repo.BaseRepo[survey_models.FeedbackModel, survey_models.FeedbackResponse]
}

func NewFeedbackRepository(i *do.Injector) (survey_services.FeedbackRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &FeedbackRepository{
		BaseRepo: repo.NewBaseRepo[survey_models.FeedbackModel, survey_models.FeedbackResponse](db, enums.FEEDBACK),
	}, nil
}

func (repo *FeedbackRepository) Create(ctx context.Context, feedbackModel survey_models.FeedbackModel) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, feedbackModel)
}

func (repo *FeedbackRepository) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	// Create filters from query parameters
	feedbackFilters, _ := survey_filters.NewFeedbackFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := feedbackFilters.BuildPipeline(matchStage, ctx, true)

	feedback, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return feedback, nil
}

func (repo *FeedbackRepository) GetWithPagination(ctx context.Context, skip, limit int64, filters survey_filters.FeedbackFilters) (*survey_models.FeedbacksWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Build pipeline with filters
	pipeline := filters.BuildPipeline(matchStage, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &filters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *FeedbackRepository) GetAll(ctx context.Context, filters survey_filters.FeedbackFilters) ([]survey_models.FeedbackResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Build pipeline with filters
	pipeline := filters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *FeedbackRepository) UpdateById(ctx context.Context, id primitive.ObjectID, feedback survey_models.FeedbackModel, filter bson.M) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{}
	}

	return repo.BaseRepo.UpdateById(ctx, id, feedback, filter)
}

func (repo *FeedbackRepository) DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{"deletedAt": bson.M{"$exists": false}}
	}

	return repo.BaseRepo.DeleteById(ctx, id, filter)
}

func (repo *FeedbackRepository) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
}
