package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"

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

type BannerRepo struct {
	*base_repos.BaseRepo[models.Banner, models.BannerResponse]
}

func NewBannerRepo(i *do.Injector) (services.BannerRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &BannerRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Banner, models.BannerResponse](db, enums.BANNER),
	}, nil
}

func (repo *BannerRepo) Create(ctx context.Context, createBannerDto *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	banner := utils.ConvertToBannerModel(createBannerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, banner)
}

func (repo *BannerRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BannerResponse, *helpers.ApiError) {

	// Create filters from query parameters
	bannerFilters, _ := filters.NewBannerFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := bannerFilters.BuildPipeline(matchStage, ctx, true)

	banner, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return banner, nil
}

func (repo *BannerRepo) GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}
	//
	// Create filters from query parameters
	bannerFilters, _ := filters.NewBannerFilters(query)

	// Build pipeline with filters
	pipeline := bannerFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *BannerRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateBannerDto *dtos.CreateBannerDto) (*models.BannerResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	matchStage := bson.M{}

	banner := utils.ConvertToBannerModel(updateBannerDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	return repo.BaseRepo.UpdateById(ctx, id, banner, matchStage)
}

func (repo *BannerRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {

	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *BannerRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
