package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"strconv"

	"github.com/iancoleman/strcase"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
)

type GroupRepo struct {
	*base_repos.BaseRepo[models.Group, models.GroupResponse]
	settingService config_services.SettingService
}

func NewGroupRepo(i *do.Injector) (services.GroupRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &GroupRepo{
		BaseRepo:       base_repos.NewBaseRepo[models.Group, models.GroupResponse](db, enums.GROUP),
		settingService: do.MustInvoke[config_services.SettingService](i),
	}, nil
}

func (repo *GroupRepo) Create(ctx context.Context, createGroupDto *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	setting, err := repo.settingService.FindOneAndIncreaseSerialNumber(ctx, enums.GROUP.GetCollectionName())
	if err != nil {
		return nil, err
	}

	// Get the English name for the code generation, with fallback
	englishName := createGroupDto.Name[enums.EN]
	if englishName == "" {
		// If English name is not available, try to get any available name
		for _, name := range createGroupDto.Name {
			if name != "" {
				englishName = name
				break
			}
		}
		// If still no name, use a default
		if englishName == "" {
			englishName = "group"
		}
	}

	code := strcase.ToSnake(englishName) + "_" + strconv.FormatUint(setting.CurrentSerialNumber, 10)
	createGroupDto.Code = &code

	group := utils.ConvertToGroupModel(createGroupDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	return repo.BaseRepo.Create(ctx, *group)
}

func (repo *GroupRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.GroupResponse, *helpers.ApiError) {
	// Create filters from query parameters
	groupFilters, _ := filters.NewGroupFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := groupFilters.BuildPipeline(matchStage, ctx, true)

	group, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return group, nil
}

func (repo *GroupRepo) GetAll(ctx context.Context, query string) ([]models.GroupResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	groupFilters, _ := filters.NewGroupFilters(query)

	// Build pipeline with filters
	pipeline := groupFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *GroupRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.GroupsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	groupFilters, _ := filters.NewGroupFilters(query)

	// Build pipeline with filters
	pipeline := groupFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &groupFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *GroupRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateGroupDto *dtos.CreateGroupDto) (*models.GroupResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	storedGroup, err := repo.GetById(ctx, id, "")
	if err != nil {
		return nil, err
	}

	updateGroupDto.Code = &storedGroup.Code

	updatedGroup := utils.ConvertToGroupModel(updateGroupDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, *updatedGroup, matchStage)
}

func (repo *GroupRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *GroupRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
