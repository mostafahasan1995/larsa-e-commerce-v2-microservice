package survey_services

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	order_utils "larsa-e-commerce-microservice-v2/internal/modules/order/utils"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	survey_filters "larsa-e-commerce-microservice-v2/internal/modules/survey/filters"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	survey_utils "larsa-e-commerce-microservice-v2/internal/modules/survey/utils"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackRepo interface {
	Create(ctx context.Context, feedback survey_models.FeedbackModel) (*survey_models.FeedbackResponse, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, feedback survey_models.FeedbackModel, filter bson.M) (*survey_models.FeedbackResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID, filter bson.M) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.FeedbackResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, filters survey_filters.FeedbackFilters) (*survey_models.FeedbacksWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, filters survey_filters.FeedbackFilters) ([]survey_models.FeedbackResponse, *helpers.ApiError)
}

type FeedbackService interface {
	Create(ctx context.Context, dto *survey_dtos.CreateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError)
	Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError)
	Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.FeedbackResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.FeedbacksWithPagination, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]survey_models.FeedbackResponse, *helpers.ApiError)
	ExportFeedbacks(ctx context.Context) ([]byte, *helpers.ApiError)
}

type feedback struct {
	repo         FeedbackRepo
	excelService excel_services.ExcelService
}

func NewFeedbackService(i *do.Injector) (FeedbackService, error) {
	return &feedback{
		repo:         do.MustInvoke[FeedbackRepo](i),
		excelService: do.MustInvoke[excel_services.ExcelService](i),
	}, nil
}

func (f *feedback) Create(ctx context.Context, dto *survey_dtos.CreateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	feedbackModel, err := survey_utils.ConvertCreateFeedbackDTOToModel(dto, ctx)
	if err != nil {
		return nil, err
	}

	return f.repo.Create(ctx, feedbackModel)
}

func (f *feedback) Update(ctx context.Context, id primitive.ObjectID, dto *survey_dtos.UpdateFeedbackDTO) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	feedbackModel, err := survey_utils.ConvertUpdateFeedbackDTOToModel(dto, id, ctx)
	if err != nil {
		return nil, err
	}

	return f.repo.UpdateById(ctx, id, feedbackModel, bson.M{})
}

func (f *feedback) Delete(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return f.repo.DeleteById(ctx, id, bson.M{})
}

func (f *feedback) GetById(ctx context.Context, id primitive.ObjectID, query string) (*survey_models.FeedbackResponse, *helpers.ApiError) {
	return f.repo.GetById(ctx, id, query)
}

func (f *feedback) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*survey_models.FeedbacksWithPagination, *helpers.ApiError) {
	filters, err := survey_filters.NewFeedbackFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	// Check user branch access
	branchAccess, err := order_utils.CheckUserBranchAccess(ctx, util.GetInjector(ctx))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to check user branch access",
		}
	}

	// If user has limited access, set the branch filter
	if branchAccess.LimitAccess && branchAccess.BranchId != nil {
		filters.BranchId = branchAccess.BranchId
	}

	return f.repo.GetWithPagination(ctx, skip, limit, *filters)
}

func (f *feedback) GetAll(ctx context.Context, query string) ([]survey_models.FeedbackResponse, *helpers.ApiError) {
	filters, err := survey_filters.NewFeedbackFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid query",
		}
	}

	// Check user branch access
	branchAccess, err := order_utils.CheckUserBranchAccess(ctx, util.GetInjector(ctx))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to check user branch access",
		}
	}

	// If user has limited access, set the branch filter
	if branchAccess.LimitAccess && branchAccess.BranchId != nil {
		filters.BranchId = branchAccess.BranchId
	}

	return f.repo.GetAll(ctx, *filters)
}

func (f *feedback) ExportFeedbacks(ctx context.Context) ([]byte, *helpers.ApiError) {
	feedbacks, apiErr := f.GetAll(ctx, `{"popFields":["waiterId","branchId","answers"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	if len(feedbacks) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 404,
			Message:    "No feedbacks found",
		}
	}

	// Convert feedbacks to Excel data
	excelData, apiErr := survey_utils.ConvertFeedbacksToExcelData(ctx, feedbacks)
	if apiErr != nil {
		return nil, apiErr
	}

	// Create Excel service
	excelService := f.excelService

	// Define columns for export
	columns := survey_utils.GetFeedbackExportColumns()

	// Generate Excel file
	excelBytes, apiErr := excelService.Export(ctx, excelData, columns, "Feedbacks Export")
	if apiErr != nil {
		return nil, apiErr
	}

	return excelBytes, nil
}
