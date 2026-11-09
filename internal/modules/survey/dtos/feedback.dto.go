package survey_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFeedbackDTO struct {
	Phone      string                 `json:"phone" validate:"required"`
	ClientName string                 `json:"clientName" validate:"required"`
	Email      *string                `json:"email" validate:"omitempty,email"`
	Birthday   string                 `json:"birthday"`
	WaiterId   primitive.ObjectID     `json:"waiterId" validate:"required"`
	BranchId   *primitive.ObjectID    `json:"branchId" validate:"omitempty"`
	VendorId   primitive.ObjectID     `json:"vendorId,required"`
	Points     float64                `json:"points" validate:"required"`
	Answers    []survey_models.Answer `json:"answers" validate:"required"`
	Notes      string                 `json:"notes"`
	WaiterRate float64                `json:"waiterRate" validate:"required"`
}

func (dto *CreateFeedbackDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Extract question IDs and check for duplicates
	questionIds := util.Map(dto.Answers, func(answer survey_models.Answer) primitive.ObjectID {
		return answer.QuestionId
	})

	if util.HasDuplicateValues(questionIds) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Duplicate answers for the same question are not allowed",
		}
	}

	// Initialize repos
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	vendorRepo := base_repos.NewBaseRepo[profile_models.VendorProfile, profile_models.ProfileResponse](args.DB, enums.PROFILE)
	waiterRepo := base_repos.NewBaseRepo[survey_models.WaiterModel, survey_models.WaiterResponse](args.DB, enums.WAITER)
	questionRepo := base_repos.NewBaseRepo[survey_models.QuestionModel, survey_models.QuestionResponse](args.DB, enums.QUESTION)

	_, err := vendorRepo.GetById(args.Context, dto.VendorId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Vendor not found",
		}
	}

	// Validate branch exists
	_, err = branchRepo.GetById(args.Context, *dto.BranchId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Branch not found",
		}
	}

	// Validate waiter exists
	waiter, err := waiterRepo.GetById(args.Context, dto.WaiterId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Waiter not found",
		}
	}

	if waiter.BranchId.Hex() != dto.BranchId.Hex() {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Waiter does not belong to the specified branch",
		}
	}

	if err != nil || waiter == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Waiter does not belong to the specified branch",
		}
	}

	// Validate questions exist
	for _, answer := range dto.Answers {
		_, err := questionRepo.GetById(args.Context, answer.QuestionId, []bson.M{})
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Question not found",
			}
		}
	}

	return nil, nil
}

type UpdateFeedbackDTO struct {
	Phone      string                 `json:"phone" validate:"required"`
	ClientName string                 `json:"clientName" validate:"required"`
	Email      string                 `json:"email" validate:"required,email"`
	Birthday   string                 `json:"birthday"`
	WaiterId   primitive.ObjectID     `json:"waiterId" validate:"required"`
	BranchId   *primitive.ObjectID    `json:"branchId" validate:"required"`
	VendorId   primitive.ObjectID     `json:"vendorId,required"`
	Points     float64                `json:"points" validate:"required"`
	Answers    []survey_models.Answer `json:"answers" validate:"required"`
	Notes      string                 `json:"notes"`
	WaiterRate float64                `json:"waiterRate" validate:"required"`
}

func (dto *UpdateFeedbackDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Extract question IDs and check for duplicates
	questionIds := util.Map(dto.Answers, func(answer survey_models.Answer) primitive.ObjectID {
		return answer.QuestionId
	})

	if util.HasDuplicateValues(questionIds) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Duplicate answers for the same question are not allowed",
		}
	}

	// Initialize repos
	waiterRepo := base_repos.NewBaseRepo[survey_models.WaiterModel, survey_models.WaiterResponse](args.DB, enums.WAITER)
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	questionRepo := base_repos.NewBaseRepo[survey_models.QuestionModel, survey_models.QuestionResponse](args.DB, enums.QUESTION)

	// Validate waiter exists
	_, err := waiterRepo.GetById(args.Context, dto.WaiterId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Waiter not found",
		}
	}

	// Validate branch exists
	_, err = branchRepo.GetById(args.Context, *dto.BranchId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Branch not found",
		}
	}

	// Check that waiter belongs to the specified branch
	waiter, err := waiterRepo.FindOne(args.Context, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"_id":      dto.WaiterId,
				"branchId": dto.BranchId,
			}),
		},
	})
	if err != nil || waiter == nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Waiter does not belong to the specified branch",
		}
	}

	// Validate questions exist
	for _, answer := range dto.Answers {
		_, err := questionRepo.GetById(args.Context, answer.QuestionId, []bson.M{})
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Question not found",
			}
		}
	}

	return nil, nil
}
