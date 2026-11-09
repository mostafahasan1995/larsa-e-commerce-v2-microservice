package recipe_dto

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	recipe_enums "larsa-e-commerce-microservice-v2/internal/modules/recipe/enums"
	unit_models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateRecipeDTO struct {
	Name            types.Localizable[string] `json:"name" validate:"required"`
	Description     types.Localizable[string] `json:"description,omitempty"`
	HowToCook       types.Localizable[string] `json:"howToCook,omitempty"`
	ImagesHowToCook []types.FileFieldDto      `json:"imagesHowToCook,omitempty"`

	Calories        *int                          `json:"calories,omitempty" validate:"omitempty,min=0"`
	CookingTime     *int                          `json:"cookingTime,omitempty" validate:"omitempty,min=1"` // in minutes
	DifficultyLevel *recipe_enums.DifficultyLevel `json:"difficultyLevel,omitempty" validate:"omitempty,oneof=easy medium hard"`
	RecipeType      *recipe_enums.RecipeType      `json:"recipeType,omitempty" validate:"omitempty,oneof=breakfast lunch dinner"`

	Ingredients  types.LocalizableArray[string] `json:"ingredients" validate:"required,dive"`
	ProductItems []RecipeProductItemDTO         `json:"productItems,omitempty" validate:"omitempty,min=1,dive"`
	Images       []types.FileFieldDto           `json:"images,omitempty"`
}

type UpdateRecipeDTO struct {
	Name            types.Localizable[string]     `json:"name" validate:"required"`
	Description     types.Localizable[string]     `json:"description,omitempty"`
	HowToCook       types.Localizable[string]     `json:"howToCook,omitempty"`
	ImagesHowToCook []types.FileFieldDto          `json:"imagesHowToCook,omitempty"`
	Calories        *int                          `json:"calories,omitempty" validate:"omitempty,min=0"`
	CookingTime     *int                          `json:"cookingTime,omitempty" validate:"omitempty,min=1"` // in minutes
	DifficultyLevel *recipe_enums.DifficultyLevel `json:"difficultyLevel,omitempty" validate:"omitempty,oneof=easy medium hard"`
	RecipeType      *recipe_enums.RecipeType      `json:"recipeType,omitempty" validate:"omitempty,oneof=breakfast lunch dinner"`

	Ingredients  types.LocalizableArray[string] `json:"ingredients" validate:"required,dive"`
	ProductItems []RecipeProductItemDTO         `json:"productItems,omitempty" validate:"omitempty,min=1,dive"`
	Images       []types.FileFieldDto           `json:"images,omitempty" validate:"required,dive"`
}

type RecipeProductItemDTO struct {
	ProductId primitive.ObjectID  `json:"productId" validate:"required"`
	UnitId    *primitive.ObjectID `json:"unitId,omitempty" validate:"omitempty"`
	Quantity  float64             `json:"quantity" validate:"required,gt=0"`
}

func (dto *CreateRecipeDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Validate product items if provided
	if len(dto.ProductItems) > 0 {
		productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
		unitRepo := base_repos.NewBaseRepo[unit_models.Unit, unit_models.UnitResponse](args.DB, enums.UNIT)

		productIds := make([]primitive.ObjectID, 0, len(dto.ProductItems))
		unitIds := make([]primitive.ObjectID, 0)
		for _, item := range dto.ProductItems {
			productIds = append(productIds, item.ProductId)
			if item.UnitId != nil {
				unitIds = append(unitIds, *item.UnitId)
			}
			if item.Quantity <= 0 {
				return nil, &helpers.ApiError{StatusCode: 400, Message: "Quantity must be greater than 0"}
			}
		}

		// Validate products
		products, err := productRepo.GetByIds(args.Context, productIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
		if err != nil {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "Failed to validate products"}
		}
		if len(products) != len(productIds) {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "One or more product IDs are invalid"}
		}

		// Validate units if provided
		if len(unitIds) > 0 {
			_, err := unitRepo.GetByIds(args.Context, unitIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
			if err != nil {
				return nil, &helpers.ApiError{StatusCode: 400, Message: "Failed to validate units"}
			}
		}
	}

	// Validate difficulty level if provided
	if dto.DifficultyLevel != nil && !recipe_enums.IsValidDifficultyLevel(*dto.DifficultyLevel) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid difficulty level",
		}
	}

	// Validate recipe type if provided
	if dto.RecipeType != nil && !recipe_enums.IsValidRecipeType(*dto.RecipeType) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid recipe type",
		}
	}

	return nil, nil
}

func (dto *UpdateRecipeDTO) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Validate product items if provided
	if len(dto.ProductItems) > 0 {
		productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
		unitRepo := base_repos.NewBaseRepo[unit_models.Unit, unit_models.UnitResponse](args.DB, enums.UNIT)

		productIds := make([]primitive.ObjectID, 0, len(dto.ProductItems))
		unitIds := make([]primitive.ObjectID, 0)
		for _, item := range dto.ProductItems {
			productIds = append(productIds, item.ProductId)
			if item.UnitId != nil {
				unitIds = append(unitIds, *item.UnitId)
			}
			if item.Quantity <= 0 {
				return nil, &helpers.ApiError{StatusCode: 400, Message: "Quantity must be greater than 0"}
			}
		}

		products, err := productRepo.GetByIds(args.Context, productIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
		if err != nil {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "Failed to validate products"}
		}
		if len(products) != len(productIds) {
			return nil, &helpers.ApiError{StatusCode: 400, Message: "One or more product IDs are invalid"}
		}
		if len(unitIds) > 0 {
			_, err := unitRepo.GetByIds(args.Context, unitIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
			if err != nil {
				return nil, &helpers.ApiError{StatusCode: 400, Message: "Failed to validate units"}
			}
		}
	}

	// Validate difficulty level if provided
	if dto.DifficultyLevel != nil && !recipe_enums.IsValidDifficultyLevel(*dto.DifficultyLevel) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid difficulty level",
		}
	}

	// Validate recipe type if provided
	if dto.RecipeType != nil && !recipe_enums.IsValidRecipeType(*dto.RecipeType) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid recipe type",
		}
	}

	return nil, nil
}
