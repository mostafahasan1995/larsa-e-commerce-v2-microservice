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

type QuestionRepository struct {
	*repo.BaseRepo[survey_models.QuestionModel, survey_models.QuestionResponse]
}

func NewQuestionRepository(i *do.Injector) (survey_services.QuestionRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &QuestionRepository{
		BaseRepo: repo.NewBaseRepo[survey_models.QuestionModel, survey_models.QuestionResponse](db, enums.QUESTION),
	}, nil
}

func (repo *QuestionRepository) Create(ctx context.Context, questionModel survey_models.QuestionModel) (*survey_models.QuestionResponse, *helpers.ApiError) {
	return repo.BaseRepo.Create(ctx, questionModel)
}

func (repo *QuestionRepository) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.QuestionResponse, *helpers.ApiError) {
	// Create filters from query parameters
	questionFilters, _ := survey_filters.NewQuestionFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := questionFilters.BuildPipeline(matchStage, ctx, true)

	question, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return question, nil
}

func (repo *QuestionRepository) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	questionFilters, _ := survey_filters.NewQuestionFilters(query)

	// Build pipeline with filters
	pipeline := questionFilters.BuildPipeline(matchStage, ctx, true)

	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &questionFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *QuestionRepository) GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	questionFilters, _ := survey_filters.NewQuestionFilters(query)

	// Build pipeline with filters
	pipeline := questionFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *QuestionRepository) UpdateById(ctx context.Context, id primitive.ObjectID, question survey_models.QuestionModel, filter bson.M) (*survey_models.QuestionResponse, *helpers.ApiError) {
	// If filter is empty, create a basic one
	if filter == nil {
		filter = bson.M{}
	}

	return repo.BaseRepo.UpdateById(ctx, id, question, filter)
}

func (repo *QuestionRepository) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	return repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
}
