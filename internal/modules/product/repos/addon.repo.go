package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddonRepo implements the AddonRepo interface
type AddonRepo struct {
	*base_repos.BaseRepo[models.Addon, models.AddonResponse]
}

// NewAddonRepo creates a new instance of AddonRepo
func NewAddonRepo(i *do.Injector) (services.AddonRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &AddonRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Addon, models.AddonResponse](db, enums.ADDON),
	}, nil
}

// Create creates a new addon
func (repo *AddonRepo) Create(ctx context.Context, createAddonDto *dtos.CreateAddonDto) (*models.AddonResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	addon := utils.ConvertToAddonModel(createAddonDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	return repo.BaseRepo.Create(ctx, *addon)
}

// GetById retrieves an addon by ID
func (repo *AddonRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AddonResponse, *helpers.ApiError) {
	// Build pipeline for aggregation
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
	}

	addon, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return addon, nil
}

// GetWithPagination retrieves addons with pagination
func (repo *AddonRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AddonsWithPagination, *helpers.ApiError) {
	// Create addon filters from query string
	addonFilters, _ := filters.NewAddonFilters(query)

	// Build match stage
	matchStage := bson.M{}

	// Build pipeline using filters
	pipeline := addonFilters.BuildPipeline(matchStage, ctx, true)

	results, apiErr := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &addonFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if apiErr != nil {
		return nil, apiErr
	}

	return results, nil
}

// UpdateById updates an addon by ID
func (repo *AddonRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateAddonDto *dtos.UpdateAddonDto) (*models.AddonResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Create an addon model with updated values
	addonValues := make([]models.AddonValue, len(updateAddonDto.AddonValues))
	for i, valueDto := range updateAddonDto.AddonValues {
		addonValues[i] = utils.ConvertUpdateAddonValueDtoToModel(valueDto)
	}

	// Create the addon with updated values
	updatedAddon := models.Addon{
		BaseAddon: models.BaseAddon{
			Label:        updateAddonDto.Label,
			MinSelection: updateAddonDto.MinSelection,
			MaxSelection: updateAddonDto.MaxSelection,
			AddonValues:  addonValues,
		},
	}

	// Set update metadata
	updatedAddon.UpdatedAt = time.Now()
	updatedAddon.UpdatedBy = &cfg.User.Id

	return repo.BaseRepo.UpdateById(ctx, id, updatedAddon, bson.M{})
}

// DeleteById deletes an addon by ID
func (repo *AddonRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
}

// BulkDelete deletes multiple addons by IDs
func (repo *AddonRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
}

func (repo *AddonRepo) GetAll(ctx context.Context, query string) ([]models.AddonResponse, *helpers.ApiError) {
	// Create addon filters from query string
	addonFilters, _ := filters.NewAddonFilters(query)

	// Build match stage
	matchStage := bson.M{}

	// Build pipeline using filters
	pipeline := addonFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}
