package survey_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	survey_enums "larsa-e-commerce-microservice-v2/internal/modules/survey/enums"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateQuestionDTO struct {
	Type        survey_enums.QuestionType   `json:"type" validate:"required,oneof=rate freeText multiChoice"`
	Order       float64                     `json:"order" validate:"required"`
	Question    types.Localizable[string]   `json:"question" validate:"required"`
	Options     types.Localizable[[]string] `json:"options"`
	BranchIds   []primitive.ObjectID        `json:"branchIds"`
	AllBranches bool                        `json:"allBranches"`
	VendorId    *primitive.ObjectID         `json:"vendorId,omitempty"`
}

func (dto *CreateQuestionDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Initialize branch repo
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)

	// Initialize vendor repo
	vendorRepo := base_repos.NewBaseRepo[profile_models.VendorProfile, profile_models.ProfileResponse](args.DB, enums.PROFILE)

	// Validate branches based on allBranches flag
	if !dto.AllBranches {
		if len(dto.BranchIds) == 0 {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "branchIds is required when allBranches is false",
			}
		}

		// Validate all branches exist
		for _, branchId := range dto.BranchIds {
			_, err := branchRepo.GetById(args.Context, branchId, []bson.M{})
			if err != nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusNotFound,
					Message:    "Branch not found",
				}
			}
		}
	}

	// Validate vendor exists
	if dto.VendorId != nil {
		_, err := vendorRepo.GetById(args.Context, *dto.VendorId, []bson.M{})
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Vendor not found",
			}
		}
	}
	// For multi-choice questions, validate that options are provided
	if dto.Type == survey_enums.MultiChoiceType {
		hasOptions := false

		// Check if any language has options
		for _, options := range dto.Options {
			if len(options) > 0 {
				hasOptions = true
				break
			}
		}

		if !hasOptions {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Multi-choice questions must have options",
			}
		}
	} else {
		if dto.Options != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Questions must not have options",
			}
		}
	}

	return nil, nil
}

type UpdateQuestionDTO struct {
	Type        survey_enums.QuestionType   `json:"type" validate:"required"`
	Order       float64                     `json:"order" validate:"required"`
	Question    types.Localizable[string]   `json:"question" validate:"required"`
	Options     types.Localizable[[]string] `json:"options"`
	BranchIds   []primitive.ObjectID        `json:"branchIds"`
	AllBranches bool                        `json:"allBranches"`
	VendorId    *primitive.ObjectID         `json:"vendorId,omitempty"`
}

func (dto *UpdateQuestionDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Initialize branch repo
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)

	// Validate branches based on allBranches flag
	if !dto.AllBranches {
		if len(dto.BranchIds) == 0 {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "branchIds is required when allBranches is false",
			}
		}

		// Validate all branches exist
		for _, branchId := range dto.BranchIds {
			_, err := branchRepo.GetById(args.Context, branchId, []bson.M{})
			if err != nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusNotFound,
					Message:    "Branch not found",
				}
			}
		}
	}

	// For multi-choice questions, validate that options are provided
	if dto.Type == survey_enums.MultiChoiceType {
		hasOptions := false

		// Check if any language has options
		for _, options := range dto.Options {
			if len(options) > 0 {
				hasOptions = true
				break
			}
		}

		if !hasOptions {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Multi-choice questions must have options",
			}
		}
	}

	return nil, nil
}
