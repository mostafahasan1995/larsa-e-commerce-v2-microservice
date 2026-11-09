package recipe_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	recipe_dto "larsa-e-commerce-microservice-v2/internal/modules/recipe/dtos"
	recipe_models "larsa-e-commerce-microservice-v2/internal/modules/recipe/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateRecipeDTOToModel(dto *recipe_dto.CreateRecipeDTO, userId primitive.ObjectID) (recipe_models.Recipe, *helpers.ApiError) {
	now := time.Now()

	recipe := recipe_models.Recipe{
		BaseRecipe: recipe_models.BaseRecipe{
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id:        primitive.NewObjectID(),
				CreatedAt: now,
				CreatedBy: &userId,
				UpdatedAt: now,
				UpdatedBy: &userId,
			},
			Name:            dto.Name,
			Description:     dto.Description,
			HowToCook:       dto.HowToCook,
			ImagesHowToCook: dto.ImagesHowToCook,
			Calories:        dto.Calories,
			CookingTime:     dto.CookingTime,
			DifficultyLevel: dto.DifficultyLevel,
			RecipeType:      dto.RecipeType,
			// Ingredients:     dto.Ingredients,
			ProductItems: convertRecipeProductItemsDTO(dto.ProductItems),
			Images:       dto.Images,
		},
	}

	return recipe, nil
}

func ConvertUpdateRecipeDTOToModel(dto *recipe_dto.UpdateRecipeDTO, id primitive.ObjectID, userId primitive.ObjectID) (recipe_models.Recipe, *helpers.ApiError) {
	now := time.Now()

	// Convert ingredients

	recipe := recipe_models.Recipe{
		BaseRecipe: recipe_models.BaseRecipe{
			MetaDataBaseModel: base_models.MetaDataBaseModel{
				Id:        id,
				UpdatedAt: now,
				UpdatedBy: &userId,
			},
			Name:            dto.Name,
			Description:     dto.Description,
			HowToCook:       dto.HowToCook,
			ImagesHowToCook: dto.ImagesHowToCook,
			Calories:        dto.Calories,
			CookingTime:     dto.CookingTime,
			DifficultyLevel: dto.DifficultyLevel,
			RecipeType:      dto.RecipeType,
			Ingredients:     dto.Ingredients,
			ProductItems:    convertRecipeProductItemsDTO(dto.ProductItems),
			Images:          dto.Images,
		},
	}

	return recipe, nil
}

func convertRecipeProductItemsDTO(items []recipe_dto.RecipeProductItemDTO) []recipe_models.RecipeProductItem {
	if len(items) == 0 {
		return nil
	}
	result := make([]recipe_models.RecipeProductItem, 0, len(items))
	for _, item := range items {
		result = append(result, recipe_models.RecipeProductItem{
			ProductId: item.ProductId,
			UnitId:    item.UnitId,
			Quantity:  item.Quantity,
		})
	}
	return result
}
