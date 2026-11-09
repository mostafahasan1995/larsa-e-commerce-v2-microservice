package survey_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebQuestionService interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*survey_models.QuestionResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError)
}

type webQuestion struct {
	repo QuestionRepo
}

func NewWebQuestionService(i *do.Injector) (WebQuestionService, error) {
	return &webQuestion{
		repo: do.MustInvoke[QuestionRepo](i),
	}, nil
}

func (q *webQuestion) GetById(ctx context.Context, id primitive.ObjectID) (*survey_models.QuestionResponse, *helpers.ApiError) {
	return q.repo.GetById(ctx, id, "")
}

func (q *webQuestion) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.QuestionsWithPagination, *helpers.ApiError) {
	return q.repo.GetWithPagination(ctx, skip, limit, query)
}

func (q *webQuestion) GetAll(ctx context.Context, query string) ([]survey_models.QuestionResponse, *helpers.ApiError) {
	return q.repo.GetAll(ctx, query)
}
