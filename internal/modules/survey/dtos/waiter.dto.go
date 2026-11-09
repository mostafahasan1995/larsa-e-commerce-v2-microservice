package survey_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	profile_models "larsa-e-commerce-microservice-v2/internal/modules/profile/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// notes
type CreateWaiterDTO struct {
	Name     types.Localizable[string] `json:"name" validate:"required"`
	Image    *types.FileFieldDto       `json:"image"`
	BranchId primitive.ObjectID        `json:"branchId" validate:"required"`
	VendorId *primitive.ObjectID       `json:"vendorId,omitempty"`
}

func (dto *CreateWaiterDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Initialize branch repo
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	vendorRepo := base_repos.NewBaseRepo[profile_models.VendorProfile, profile_models.ProfileResponse](args.DB, enums.PROFILE)

	// Validate branch exists
	_, err := branchRepo.GetById(args.Context, dto.BranchId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Branch not found",
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

	return nil, nil
}

type UpdateWaiterDTO struct {
	Name     types.Localizable[string] `json:"name" validate:"required"`
	Image    *types.FileFieldDto       `json:"image"`
	BranchId primitive.ObjectID        `json:"branchId" validate:"required"`
	VendorId *primitive.ObjectID       `json:"vendorId,omitempty"`
}

func (dto *UpdateWaiterDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Initialize branch repo
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)

	// Validate branch exists
	_, err := branchRepo.GetById(args.Context, dto.BranchId, []bson.M{})
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Branch not found",
		}
	}

	// Validate name has at least one language defined
	if len(dto.Name) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Name must have at least one language defined",
		}
	}

	return nil, nil
}
