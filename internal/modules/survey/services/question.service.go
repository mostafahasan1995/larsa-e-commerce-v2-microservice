package survey_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_utils "larsa-e-commerce-microservice-v2/internal/modules/survey/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionRepo interface {
	Create(ctx context.Context, question survey_models.QuestionModel) (*survey_models.QuestionResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, question survey_models.QuestionModel, filter bson.M) (*survey_models.QuestionResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.QuestionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError)
}

type QuestionService interface {
	Create(ctx context.Context, dto *survey_dtos.CreateQuestionDTO) (*survey_models.QuestionResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateQuestionDTO) (*survey_models.QuestionResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.QuestionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError)
}

type question struct {
	repo QuestionRepo
}

func NewQuestionService(i *do.Injector) (QuestionService, error) {
	return &question{
		repo: do.MustInvoke[QuestionRepo](i),
	}, nil
}

func (q *question) Create(ctx context.Context, dto *survey_dtos.CreateQuestionDTO) (*survey_models.QuestionResponse, *helpers.ApiError) {
	questionModel, err := survey_utils.ConvertCreateQuestionDTOToModel(dto, ctx)
	if err != nil {
		return nil, err
	}

	return q.repo.Create(ctx, questionModel)
}

func (q *question) Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateQuestionDTO) (*survey_models.QuestionResponse, *helpers.ApiError) {
	questionModel, err := survey_utils.ConvertUpdateQuestionDTOToModel(dto, id, ctx)
	if err != nil {
		return nil, err
	}

	return q.repo.UpdateById(ctx, id, questionModel, bson.M{})
}

func (q *question) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return q.repo.DeleteById(ctx, id)
}

func (q *question) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.QuestionResponse, *helpers.ApiError) {
	return q.repo.GetById(ctx, id, query)
}

func (q *question) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError) {
	return q.repo.GetWithPagination(ctx, skip, limit, query)
}

func (q *question) GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError) {
	return q.repo.GetAll(ctx, query)
}
