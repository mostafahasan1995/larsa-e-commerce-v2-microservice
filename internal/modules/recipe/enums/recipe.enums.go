package recipe_enums

// DifficultyLevel represents the difficulty level of a recipe
type DifficultyLevel string

const (
	EASY   DifficultyLevel = "easy"
	MEDIUM DifficultyLevel = "medium"
	HARD   DifficultyLevel = "hard"
)

// GetDifficultyLevels returns all available difficulty levels
func GetDifficultyLevels() []DifficultyLevel {
	return []DifficultyLevel{EASY, MEDIUM, HARD}
}

// IsValidDifficultyLevel checks if the provided difficulty level is valid
func IsValidDifficultyLevel(level DifficultyLevel) bool {
	switch level {
	case EASY, MEDIUM, HARD:
		return true
	default:
		return false
	}
}

// RecipeType represents the type/category of a recipe
type RecipeType string

const (
	BREAKFAST RecipeType = "breakfast"
	LUNCH     RecipeType = "lunch"
	DINNER    RecipeType = "dinner"
)

// GetRecipeTypes returns all available recipe types
func GetRecipeTypes() []RecipeType {
	return []RecipeType{BREAKFAST, LUNCH, DINNER}
}

// IsValidRecipeType checks if the provided recipe type is valid
func IsValidRecipeType(recipeType RecipeType) bool {
	switch recipeType {
	case BREAKFAST, LUNCH, DINNER:
		return true
	default:
		return false
	}
}
