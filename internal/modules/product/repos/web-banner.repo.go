package product_repos

import (
	"context"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type WebBannerRepo struct {
	*base_repos.BaseRepo[models.Banner, models.BannerResponse]
}

func NewWebBannerRepo(i *do.Injector) (services.WebBannerRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &WebBannerRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Banner, models.BannerResponse](db, enums.BANNER),
	}, nil
}

func (repo *WebBannerRepo) GetAll(ctx context.Context, query string) ([]models.BannerResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	bannerFilters, _ := filters.NewBannerFilters(query)
	bannerFilters.IsActive = &[]bool{true}[0]

	// Build pipeline with filters
	pipeline := bannerFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}
