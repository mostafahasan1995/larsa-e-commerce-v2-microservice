package seed_data

import (
	"encoding/json"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	config_enums "larsa-e-commerce-microservice-v2/internal/modules/config/enums"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OldCategory represents the structure of the JSON data
type OldCategory struct {
	ID           string              `json:"_id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ParentId     string              `json:"parentId,omitempty"`
	PictureFile  *types.FileFieldDto `json:"pictureFile,omitempty"`
	IconFile     *types.FileFieldDto `json:"iconFile,omitempty"`
	Translations struct {
		AR struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"ar"`
		EN struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"en"`
	} `json:"translations"`
}

// LoadCategoriesFromJSON loads categories from either the embedded JSON file or filesystem
func LoadCategoriesFromJSON(jsonPath string) ([]config_models.Category, error) {
	// If embedded file fails, try to read from filesystem
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read categories file (both embedded and filesystem): %v", err)
	}

	var oldCategories []OldCategory
	if err := json.Unmarshal(data, &oldCategories); err != nil {
		return nil, err
	}

	var categories []config_models.Category
	for _, old := range oldCategories {
		category := config_models.Category{
			BaseCategory: config_models.BaseCategory{
				Name: types.Localizable[string]{
					enums.EN: old.Translations.EN.Name,
					enums.AR: old.Translations.AR.Name,
				},
				Description: &types.Localizable[string]{
					enums.EN: old.Translations.EN.Description,
					enums.AR: old.Translations.AR.Description,
				},
				Status: &config_models.Status{
					Name:        config_enums.APPROVED,
					TriggeredBy: nil,
					Timestamp:   time.Now(),
				},
				Visibilities: []config_enums.CategoryVisibility{
					config_enums.CATEGORIES_LIST,
					config_enums.SEARCH_RESULT,
					config_enums.CATEGORIES_CAROUSEL,
					config_enums.INVISIBLE,
				},
				MetaDataBaseModel: models.MetaDataBaseModel{
					Id:        func() primitive.ObjectID { id, _ := primitive.ObjectIDFromHex(old.ID); return id }(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		if old.ParentId != "" {
			parentID, _ := primitive.ObjectIDFromHex(old.ParentId)
			category.ParentId = &parentID
		}

		if old.PictureFile != nil {
			category.Image = &types.Localizable[types.FileFieldDto]{
				enums.EN: *old.PictureFile,
				enums.AR: *old.PictureFile,
			}
		}

		if old.IconFile != nil {
			category.Icon = &types.Localizable[types.FileFieldDto]{
				enums.EN: *old.IconFile,
				enums.AR: *old.IconFile,
			}
		}
		categories = append(categories, category)
	}

	return categories, nil
}

var Categories = []config_models.Category{}

// Initialize categories from JSON
func InitCategories() {
	// Get the path relative to the current file
	_, currentFile, _, _ := runtime.Caller(0)
	jsonPath := filepath.Join(filepath.Dir(currentFile), "json/old-categories.json")

	loadedCategories, err := LoadCategoriesFromJSON(jsonPath)
	if err != nil {
		// Handle error appropriately for your application
		panic(fmt.Sprintf("Failed to load categories: %v", err))
	}
	Categories = loadedCategories
}
