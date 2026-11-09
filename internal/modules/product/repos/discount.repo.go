package product_repos

import (
	"context"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type DiscountRepo struct {
	productRepo *base_repos.BaseRepo[models.Product, models.ProductResponse]
}

func NewDiscountRepo(i *do.Injector) (services.DiscountRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &DiscountRepo{
		productRepo: base_repos.NewBaseRepo[models.Product, models.ProductResponse](db, enums.PRODUCT),
	}, nil
}

func (repo *DiscountRepo) HandleNewProductsDiscount(ctx context.Context, handleNewProductsDiscountDto dtos.HandleNewProductsDiscountDto) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	priceFilter := bson.M{}
	if handleNewProductsDiscountDto.IsDefault {
		priceFilter["isDefault"] = bson.M{"$eq": true}
	}
	if handleNewProductsDiscountDto.BranchId != nil {
		priceFilter["branchIds"] = bson.M{"$in": []primitive.ObjectID{*handleNewProductsDiscountDto.BranchId}}
	}
	filter := bson.M{
		"_id": bson.M{"$in": handleNewProductsDiscountDto.ProductIds},
		"prices": bson.M{
			"$elemMatch": priceFilter,
		},
	}

	arrayFilterPipeline := bson.M{}

	if handleNewProductsDiscountDto.IsDefault {
		arrayFilterPipeline["price.isDefault"] = bson.M{"$eq": true}
	}

	if handleNewProductsDiscountDto.BranchId != nil {
		arrayFilterPipeline["price.branchIds"] = *handleNewProductsDiscountDto.BranchId
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			arrayFilterPipeline,
		},
	})

	update := bson.M{
		"$set": bson.M{
			"prices.$[price].discount": handleNewProductsDiscountDto.Discount,
			"updatedAt":                time.Now(),
			"updatedBy":                cfg.User.Id,
		},
	}

	collection, err := repo.productRepo.GetCollection(ctx)
	if err != nil {
		return err
	}
	collection.UpdateMany(ctx, filter, update, arrayFilters)

	return nil
}
