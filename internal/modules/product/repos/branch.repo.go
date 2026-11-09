package product_repos

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
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
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
)

type BranchRepo struct {
	*base_repos.BaseRepo[models.Branch, models.BranchResponse]
	branchProductRepo services.BranchProductRepo
	productRepo       *base_repos.BaseRepo[models.Product, models.ProductResponse]
}

func NewBranchRepo(i *do.Injector) (services.BranchRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &BranchRepo{
		BaseRepo:          base_repos.NewBaseRepo[models.Branch, models.BranchResponse](db, enums.BRANCH),
		branchProductRepo: do.MustInvoke[services.BranchProductRepo](i),
		productRepo:       base_repos.NewBaseRepo[models.Product, models.ProductResponse](db, enums.PRODUCT),
	}, nil
}

func (repo *BranchRepo) Create(ctx context.Context, createBranchDto dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	branch := utils.ConvertToBranchModel(createBranchDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	// Set VendorId based on user type
	if cfg.User.Type == enums.Vendor {
		branch.VendorId = &cfg.User.Id
	}

	var branchResponse *models.BranchResponse

	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		linkBranchToProductsErr := repo.linkBranchToProducts(ctx, branch.Id, createBranchDto.Products)
		if linkBranchToProductsErr != nil {
			return nil, linkBranchToProductsErr
		}

		createdBranch, createErr := repo.BaseRepo.Create(sessCtx, branch)
		if createErr != nil {
			return nil, createErr
		}

		branchResponse = createdBranch
		return createdBranch, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, helpers.DefaultError(sessionErr)
	}

	return branchResponse, nil
}

func (repo *BranchRepo) GetById(ctx context.Context, id primitive.ObjectID, branchFilters filters.BranchFilters) (*models.BranchResponse, *helpers.ApiError) {
	matchStage := bson.M{}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, true)
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *BranchRepo) GetAll(ctx context.Context, query string, branchFilters filters.BranchFilters) ([]models.BranchResponse, *helpers.ApiError) {
	matchStage := bson.M{}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, true)
	return repo.BaseRepo.GetAll(ctx, pipeline)
}

func (repo *BranchRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string, branchFilters filters.BranchFilters) (*models.BranchesWithPagination, *helpers.ApiError) {
	matchStage := bson.M{}
	pipeline := branchFilters.BuildPipeline(matchStage, ctx, false)
	return repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &branchFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
}

func (repo *BranchRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateBranchDto dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	branch := utils.ConvertToBranchModel(updateBranchDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	var branchResponse *models.BranchResponse

	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		linkBranchToProductsErr := repo.linkBranchToProducts(ctx, id, updateBranchDto.Products)
		if linkBranchToProductsErr != nil {
			return nil, linkBranchToProductsErr
		}

		matchStage := bson.M{}

		updatedBranch, updateErr := repo.BaseRepo.UpdateById(ctx, id, branch, matchStage)
		if updateErr != nil {
			return nil, updateErr
		}

		branchResponse = updatedBranch
		return updatedBranch, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, helpers.DefaultError(sessionErr)
	}

	return branchResponse, nil
}

func (repo *BranchRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
	if err != nil {
		return err
	}

	return nil
}

func (repo *BranchRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
	if bulkDeleteErr != nil {
		return bulkDeleteErr
	}
	return nil
}

func (repo *BranchRepo) UpdateZones(ctx context.Context, id primitive.ObjectID, zones []models.BranchZone) (*models.BranchResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	updateDoc := bson.M{
		"$set": bson.M{
			"zones":     zones,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}
	_, updateErr := collection.UpdateOne(ctx, filter, updateDoc)
	if updateErr != nil {
		return nil, helpers.DefaultError(updateErr)
	}

	// Return the updated branch
	matchStage := bson.M{}
	pipeline := []bson.M{{"$match": matchStage}}
	return repo.BaseRepo.GetById(ctx, id, pipeline)
}

func (repo *BranchRepo) upsertProductBranchIdAttribute(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	productUpdateManyFilter := base_utils.AddSoftDeletedFilter(bson.M{
		"_id": upsertProductBranchDto.ProductId,
	})

	product, productErr := repo.productRepo.FindOne(ctx, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"_id":             upsertProductBranchDto.ProductId,
				"attributes.code": product_enums.SYSTEM_BRANCH,
				"attributes.type": product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY,
			}),
		},
	})
	if productErr != nil {
		return productErr
	}

	productUpdateMany := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	// First try to update existing attribute
	arrayFilters := []interface{}{
		bson.M{
			"elem.code": product_enums.SYSTEM_BRANCH,
			"elem.type": product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY,
		},
	}
	collection, _ := repo.productRepo.GetCollection(ctx)

	if product != nil {

		updateResult, updateErr := collection.UpdateOne(ctx, productUpdateManyFilter, bson.M{
			"$set": productUpdateMany["$set"],
			"$addToSet": bson.M{
				"attributes.$[elem].value": upsertProductBranchDto.BranchId,
			},
		}, options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters}))
		if updateErr != nil {
			return helpers.DefaultError(updateErr)
		}

		// Log the update result
		if updateResult.MatchedCount == 0 {
			return helpers.DefaultError(fmt.Errorf("no products found to update"))
		}

		if updateResult.ModifiedCount == 0 {
			return helpers.DefaultError(fmt.Errorf("no products were updated"))
		}

	} else {
		// If no documents were modified, insert new attribute
		_, insertErr := collection.UpdateOne(ctx, productUpdateManyFilter, bson.M{
			"$set": productUpdateMany["$set"],
			"$addToSet": bson.M{
				"attributes": bson.M{
					"code":  product_enums.SYSTEM_BRANCH,
					"type":  product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY,
					"value": []primitive.ObjectID{upsertProductBranchDto.BranchId},
				},
			},
		})
		if insertErr != nil {
			return helpers.DefaultError(insertErr)

		}
	}

	updateBranchProductPriceErr := repo.updateBranchProductPrice(ctx, upsertProductBranchDto)
	if updateBranchProductPriceErr != nil {
		return updateBranchProductPriceErr
	}

	return nil
}

func (repo *BranchRepo) updateBranchProductPrice(ctx context.Context, upsertProductBranchDto dtos.UpsertProductBranchDto) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// --- Update product price for the branch if OverridePrice is provided ---
	if upsertProductBranchDto.OverridePrice != nil {
		productCollection, _ := repo.productRepo.GetCollection(ctx)
		priceItem := models.PriceItem{
			Value:        upsertProductBranchDto.OverridePrice.Value,
			CurrencyCode: upsertProductBranchDto.OverridePrice.CurrencyCode,
			UnitId:       upsertProductBranchDto.OverridePrice.UnitId,
		}

		// Build the items filter for the default price
		currencyCodeFilter := bson.M{
			"currencyCode": priceItem.CurrencyCode,
			"value":        priceItem.Value,
		}
		if priceItem.UnitId != nil {
			currencyCodeFilter["unitId"] = priceItem.UnitId
		}

		/*
		* We need to make sure that the new price does not exist in the default price nor in the the same branch,
		* to avoid pushing the same price twice.
		* Timestamp: 2025-06-29 17:42:00
		 */
		fmt.Println("currencyCodeFilter", currencyCodeFilter)
		productPriceCount, _ := productCollection.CountDocuments(ctx, base_utils.AddSoftDeletedFilter(bson.M{
			"_id": upsertProductBranchDto.ProductId,
			"prices": bson.M{
				"$elemMatch": bson.M{
					"$or": []bson.M{
						{
							"branchIds": bson.M{"$in": []primitive.ObjectID{upsertProductBranchDto.BranchId}},
							"isDefault": false,
						},
						{
							"items": bson.M{
								"$elemMatch": currencyCodeFilter,
							},
							"isDefault": true,
						},
					},
				},
			},
		}))

		newPriceFilter := base_utils.AddSoftDeletedFilter(bson.M{
			"_id": upsertProductBranchDto.ProductId,
			"prices": bson.M{
				"$elemMatch": bson.M{
					"branchIds": bson.M{"$nin": []primitive.ObjectID{upsertProductBranchDto.BranchId}},
				},
			},
		})

		newPrice := models.Price{
			BasePrice: models.BasePrice{
				BranchIds: &[]primitive.ObjectID{upsertProductBranchDto.BranchId},
				Items:     []models.PriceItem{priceItem},
				IsDefault: util.BoolPtr(false),
			},
		}

		/*
		* We only update the price if the price item is different from the default price item and
		* If there is no price for this branch.
		* Timestamp: 2025-06-29 17:42:00
		 */

		matchedCount := int64(0)
		if productPriceCount == 0 {
			newPriceResult, err := productCollection.UpdateOne(ctx, newPriceFilter, bson.M{
				"$push": bson.M{"prices": newPrice},
				"$set":  bson.M{"updatedAt": time.Now(), "updatedBy": cfg.User.Id},
			})
			if err != nil {
				return helpers.DefaultError(err)
			}
			matchedCount = newPriceResult.MatchedCount
		}

		if matchedCount == 1 {
			return nil
		}

		existsPriceItemFilter := bson.M{
			"currencyCode": priceItem.CurrencyCode,
			"value":        bson.M{"$ne": priceItem.Value},
		}
		if priceItem.UnitId != nil {
			existsPriceItemFilter["unitId"] = bson.M{"$eq": priceItem.UnitId}
		}

		existsFilter := base_utils.AddSoftDeletedFilter(bson.M{
			"_id": upsertProductBranchDto.ProductId,
			"prices": bson.M{
				"$elemMatch": bson.M{
					"branchIds": bson.M{"$in": []primitive.ObjectID{upsertProductBranchDto.BranchId}},
					"items": bson.M{
						"$elemMatch": existsPriceItemFilter,
					},
				},
			},
		})

		priceItemArrayFilter := bson.M{
			"item.currencyCode": bson.M{"$eq": priceItem.CurrencyCode},
			"item.value":        bson.M{"$ne": priceItem.Value},
		}
		if priceItem.UnitId != nil {
			priceItemArrayFilter["item.unitId"] = bson.M{"$eq": priceItem.UnitId}
		}
		arrayFiltersInterface := []interface{}{
			bson.M{
				"price.branchIds": bson.M{"$in": []primitive.ObjectID{upsertProductBranchDto.BranchId}},
			},
			priceItemArrayFilter,
		}

		arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: arrayFiltersInterface,
		})
		update := bson.M{
			"$set": bson.M{
				"updatedAt":                           time.Now(),
				"updatedBy":                           cfg.User.Id,
				"prices.$[price].items.$[item].value": priceItem.Value,
			},
		}

		/*
		* If there is a price and price item belong to the same branch, we update the price item value
		* NOTE: It's an update to an existing price item not pushing new one
		* Timestamp: 2025-06-29 17:42:00
		 */
		updatedPriceItemResult, err := productCollection.UpdateOne(ctx, existsFilter, update, arrayFilters)
		if err != nil {
			return helpers.DefaultError(err)
		}

		if updatedPriceItemResult.MatchedCount == 1 {
			return nil
		}

		/*
		* NOTE: If you noticed that we are using items.currencyCode instead of currencyCode,
		* For some reason, adding $elemMatch to the items like before won't work.
		* Timestamp: 2025-07-01 17:42:00
		 */
		doesPriceItemExistFilter := bson.M{
			"branchIds":          upsertProductBranchDto.BranchId,
			"items.currencyCode": bson.M{"$ne": priceItem.CurrencyCode},
		}
		if priceItem.UnitId != nil {
			doesPriceItemExistFilter["items.unitId"] = bson.M{"$eq": priceItem.UnitId}
		}

		doesPriceExistFilter := base_utils.AddSoftDeletedFilter(bson.M{
			"_id": upsertProductBranchDto.ProductId,
			"prices": bson.M{
				"$elemMatch": doesPriceItemExistFilter,
			},
		})

		pushArrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"price.branchIds": bson.M{"$in": []primitive.ObjectID{upsertProductBranchDto.BranchId}},
				},
			},
		})

		// Push the new item
		pushUpdate := bson.M{
			"$push": bson.M{
				"prices.$[price].items": priceItem,
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
				"updatedBy": cfg.User.Id,
			},
		}

		/*
		* To push a new price item to the price item array, It must belong to the same branch
		* and the price item must be different from the price item in the array
		 */
		_, pushErr := productCollection.UpdateOne(ctx, doesPriceExistFilter, pushUpdate, pushArrayFilters)
		if pushErr != nil {
			return helpers.DefaultError(pushErr)
		}

	}
	return nil
}

func (repo *BranchRepo) linkBranchToProducts(ctx context.Context, branchId primitive.ObjectID, linkBranchToProductDtoList *[]dtos.LinkBranchToProductDto) *helpers.ApiError {
	if linkBranchToProductDtoList == nil {
		return nil
	}

	cfg, _ := util.GetReqAppCfg(ctx)
	branchProductCollection := base_repos.NewBaseRepo[models.BranchProduct, models.BranchProductResponse](repo.BaseRepo.GetDB(), enums.BRANCH_PRODUCT)

	productBranches := make([]models.BranchProduct, len(*linkBranchToProductDtoList))
	productIds := make([]primitive.ObjectID, len(*linkBranchToProductDtoList))
	for i, linkBranchToProductDto := range *linkBranchToProductDtoList {
		productIds[i] = linkBranchToProductDto.ProductId
		dto := dtos.UpsertProductBranchDto{
			BaseLinkBranchProductDto: linkBranchToProductDto.BaseLinkBranchProductDto,
			BranchId:                 branchId,
			ProductId:                linkBranchToProductDto.ProductId,
		}

		productBranches[i] = utils.ConvertToBranchProductModel(dto, types.AdditionalOptions{
			UserId: cfg.User.Id,
			Mode:   types.UPDATE,
		})

		updateProductAttributesErr := repo.upsertProductBranchIdAttribute(ctx, dto)
		if updateProductAttributesErr != nil {
			return updateProductAttributesErr
		}
	}

	/*
	* Delete products that are not in the list (old ones)
	* Upsert products that are in the list (new ones)
	 */
	branchProductCollection.BulkDeleteFilter(ctx, bson.M{"branchId": branchId, "productId": bson.M{"$nin": productIds}})
	err := repo.branchProductRepo.BulkUpsert(ctx, productBranches)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BranchRepo) LinkAllProductsToTheFirstBranch(ctx context.Context) *helpers.ApiError {
	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		branch, err := repo.BaseRepo.FindOne(ctx, []bson.M{{"$match": bson.M{"deletedAt": bson.M{"$exists": false}}}})
		if err != nil {
			return nil, err
		}

		if branch == nil {
			return nil, helpers.DefaultError(fmt.Errorf("branch not found"))
		}

		products, err := repo.productRepo.GetAll(ctx, []bson.M{{"$match": bson.M{"deletedAt": bson.M{"$exists": false}}}})
		if err != nil {
			return nil, err
		}

		linkBranchToProductDtoList := make([]dtos.LinkBranchToProductDto, len(products))

		for i, product := range products {
			linkBranchToProductDtoList[i] = dtos.LinkBranchToProductDto{
				ProductId: product.Id,
				BaseLinkBranchProductDto: dtos.BaseLinkBranchProductDto{
					OverridePrice:       nil,
					OverrideDescription: nil,
					IsActive:            true,
					OutOfStock:          false,
				},
			}
		}

		linkBranchToProductsErr := repo.linkBranchToProducts(ctx, branch.Id, &linkBranchToProductDtoList)
		if linkBranchToProductsErr != nil {
			return nil, linkBranchToProductsErr
		}

		return nil, nil
	})

	if sessionErr != nil {
		return helpers.DefaultError(sessionErr)
	}

	return nil
}
