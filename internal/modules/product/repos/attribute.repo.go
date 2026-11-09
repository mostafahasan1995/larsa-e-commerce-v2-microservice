package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

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

type AttributeRepo struct {
	*base_repos.BaseRepo[models.Attribute, models.AttributeResponse]
}

func NewAttributeRepo(i *do.Injector) (services.AttributeRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)

	repo := &AttributeRepo{
		BaseRepo: base_repos.NewBaseRepo[models.Attribute, models.AttributeResponse](db, enums.ATTRIBUTE),
	}

	return repo, nil
}

func (repo *AttributeRepo) Create(ctx context.Context, createAttributeDto *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	attribute := utils.ConvertToAttributeModel(createAttributeDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	return repo.BaseRepo.Create(ctx, *attribute)
}

func (repo *AttributeRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.AttributeResponse, *helpers.ApiError) {
	// Create filters from query parameters
	attributeFilters, _ := filters.NewAttributeFilters(query)

	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := attributeFilters.BuildPipeline(matchStage, ctx, true)

	attribute, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return attribute, nil
}

func (repo *AttributeRepo) GetAll(ctx context.Context, query string) ([]models.AttributeResponse, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	attributeFilters, _ := filters.NewAttributeFilters(query)

	// Build pipeline with filters
	pipeline := attributeFilters.BuildPipeline(matchStage, ctx, true)

	results, err := repo.BaseRepo.GetAll(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *AttributeRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.AttributesWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	attributeFilters, _ := filters.NewAttributeFilters(query)

	// Build pipeline with filters
	pipeline := attributeFilters.BuildPipeline(matchStage, ctx, false)

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &attributeFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *AttributeRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateAttributeDto *dtos.CreateAttributeDto) (*models.AttributeResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	candidateAttribute := utils.ConvertToAttributeModel(updateAttributeDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	storedAttribute, findErr := repo.GetById(ctx, id, "")
	if findErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Attribute not found",
		}
	}

	if storedAttribute.Tag == product_enums.SYSTEM {
		candidateAttribute.Code = storedAttribute.Code
		candidateAttribute.Tag = product_enums.SYSTEM
		candidateAttribute.GroupId = storedAttribute.GroupId

		switch *storedAttribute.RefField {
		case product_enums.BRANCH:
			candidateAttribute.RefField = storedAttribute.RefField
		case product_enums.ADDON:
			candidateAttribute.RefField = storedAttribute.RefField
		}

		uniqueBelongsTo := make(map[product_enums.BelongTo]bool)
		for _, belongTo := range storedAttribute.BelongsTo {
			uniqueBelongsTo[belongTo] = true
		}
		for _, belongTo := range candidateAttribute.BelongsTo {
			uniqueBelongsTo[belongTo] = true
		}
		candidateAttribute.BelongsTo = []product_enums.BelongTo{}
		for belongTo := range uniqueBelongsTo {
			candidateAttribute.BelongsTo = append(candidateAttribute.BelongsTo, belongTo)
		}

	}

	return repo.BaseRepo.UpdateById(ctx, id, *candidateAttribute, matchStage)
}

func (repo *AttributeRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	// err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}, "tag": product_enums.CUSTOM})
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

func (repo *AttributeRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}, "tag": product_enums.CUSTOM})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}
