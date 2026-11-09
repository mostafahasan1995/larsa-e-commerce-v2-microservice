package recipe_models

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	base_models "larsa-e-commerce-microservice-v2/internal/modules/base/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	recipe_enums "larsa-e-commerce-microservice-v2/internal/modules/recipe/enums"
	unit_models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseRecipe struct {
	base_models.MetaDataBaseModel `bson:",inline"`

	Name            types.Localizable[string]     `bson:"name" json:"name"`
	Description     types.Localizable[string]     `bson:"description,omitempty" json:"description,omitempty"`
	HowToCook       types.Localizable[string]     `bson:"howToCook,omitempty" json:"howToCook,omitempty"`
	ImagesHowToCook []types.FileFieldDto          `bson:"imagesHowToCook,omitempty" json:"imagesHowToCook,omitempty"`
	Calories        *int                          `bson:"calories,omitempty" json:"calories,omitempty"`
	CookingTime     *int                          `bson:"cookingTime,omitempty" json:"cookingTime,omitempty"` // in minutes
	DifficultyLevel *recipe_enums.DifficultyLevel `bson:"difficultyLevel,omitempty" json:"difficultyLevel,omitempty"`
	RecipeType      *recipe_enums.RecipeType      `bson:"recipeType,omitempty" json:"recipeType,omitempty"`

	Ingredients  types.LocalizableArray[string] `bson:"ingredients" json:"ingredients"`
	ProductItems []RecipeProductItem            `bson:"productItems" json:"productItems"`
	Images       []types.FileFieldDto           `bson:"images,omitempty" json:"images,omitempty"`
}

type Recipe struct {
	BaseRecipe `bson:",inline"`
}

type RecipeResponse struct {
	BaseRecipe `bson:",inline"`
}

type RecipesWithPagination = types.DataWithPagination[RecipeResponse]

// RecipeProductItem represents a product entry inside a recipe with its unit and quantity
type RecipeProductItem struct {
	ProductId primitive.ObjectID              `bson:"productId" json:"productId"`
	UnitId    *primitive.ObjectID             `bson:"unitId,omitempty" json:"unitId,omitempty"`
	Quantity  float64                         `bson:"quantity" json:"quantity"`
	Product   *product_models.ProductResponse `bson:"product,omitempty" json:"product,omitempty"`
	Unit      *unit_models.UnitResponse       `bson:"unit,omitempty" json:"unit,omitempty"`
}
