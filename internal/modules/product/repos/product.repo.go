package product_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"net/http"
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
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type ProductRepo struct {
	*base_repos.BaseRepo[models.Product, models.ProductResponse]
	branchProductRepo services.BranchProductRepo
}

func NewProductRepo(i *do.Injector) (services.ProductRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &ProductRepo{
		BaseRepo:          base_repos.NewBaseRepo[models.Product, models.ProductResponse](db, enums.PRODUCT),
		branchProductRepo: do.MustInvoke[services.BranchProductRepo](i),
	}, nil
}

func (repo *ProductRepo) Create(ctx context.Context, createProductDto *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	product := product_utils.ConvertToProductModel(createProductDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})

	// Set VendorId based on user type
	if cfg.User.Type == enums.Vendor {
		product.VendorId = &cfg.User.Id
	}

	var createdProduct *models.ProductResponse

	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		linkProductToBranchesErr := repo.linkProductToBranches(ctx, product, createProductDto.BranchIds)
		if linkProductToBranchesErr != nil {
			return nil, linkProductToBranchesErr
		}

		storedProduct, createErr := repo.BaseRepo.Create(sessCtx, *product)
		if createErr != nil {
			return nil, createErr
		}

		createdProduct = storedProduct

		return storedProduct, nil

	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, helpers.DefaultError(sessionErr)
	}

	return createdProduct, nil
}

func (repo *ProductRepo) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.ProductResponse, *helpers.ApiError) {
	// Create filters from query parameters
	productFilters, _ := filters.NewProductFilters(query)
	productFilters.IncludeAllProductTypes = &[]bool{true}[0]
	// Build pipeline with filters
	matchStage := bson.M{}

	pipeline := productFilters.BuildPipeline(matchStage, ctx, filters.BuildPipelineArgs{
		AddExtraOperations: true,
	})

	product, findErr := repo.BaseRepo.GetById(ctx, id, pipeline)
	if findErr != nil {
		return nil, findErr
	}

	return product, nil
}

func (repo *ProductRepo) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.ProductsWithPagination, *helpers.ApiError) {
	// Base match stage with deletedAt filter
	matchStage := bson.M{}

	// Create filters from query parameters
	productFilters, _ := filters.NewProductFilters(query)

	// Build pipeline with filters
	pipeline := productFilters.BuildPipeline(matchStage, ctx, filters.BuildPipelineArgs{
		AddExtraOperations: true,
	})

	results, err := repo.BaseRepo.GetWithPagination(ctx, helpers.BuildPaginationPipelineProps{
		BaseFilters:    &productFilters.BaseFilters,
		FilterPipeline: pipeline,
		Skip:           skip,
		Limit:          limit,
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ProductRepo) UpdateById(ctx context.Context, id primitive.ObjectID, updateProductDto *dtos.CreateProductDto) (*models.ProductResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Build pipeline with filters
	matchStage := bson.M{}

	product := product_utils.ConvertToProductModel(updateProductDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})
	product.Id = id
	var updatedProduct *models.ProductResponse

	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		linkProductToBranchesErr := repo.linkProductToBranches(ctx, product, updateProductDto.BranchIds)
		if linkProductToBranchesErr != nil {
			return nil, linkProductToBranchesErr
		}

		storedProduct, updateErr := repo.BaseRepo.UpdateById(sessCtx, id, *product, matchStage)
		if updateErr != nil {
			return nil, updateErr
		}

		updatedProduct = storedProduct
		return storedProduct, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, helpers.DefaultError(sessionErr)
	}

	return updatedProduct, nil
}

func (repo *ProductRepo) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// For vendor users, verify they own the product
	if cfg.User.Type == enums.Vendor {
		// Get the product to check vendorId
		product, getErr := repo.GetById(ctx, id, "")
		if getErr != nil {
			return getErr
		}

		// Vendors can only delete their own products
		if product.VendorId != &cfg.User.Id {
			return &helpers.ApiError{
				StatusCode: http.StatusForbidden,
				Message:    "not authorized to delete this product",
			}
		}
	}

	// Start a session for the transaction
	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := repo.BaseRepo.DeleteById(ctx, id, bson.M{"deletedAt": bson.M{"$exists": false}})
		if err != nil {
			return nil, err
		}

		unlinkProductErr := repo.unlinkProductFromProductionLine(ctx, sessCtx, id, nil)
		if unlinkProductErr != nil {
			return nil, unlinkProductErr
		}
		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *ProductRepo) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// For vendor users, verify they own all the products they're trying to delete
	if cfg.User.Type == enums.Vendor {
		for _, id := range ids {
			// Get the product to check vendorId
			product, getErr := repo.GetById(ctx, id, "")
			if getErr != nil {
				return getErr
			}

			// Vendors can only delete their own products
			if product.VendorId != &cfg.User.Id {
				return &helpers.ApiError{
					StatusCode: http.StatusForbidden,
					Message:    "not authorized to delete product with id: " + id.Hex(),
				}
			}
		}
	}

	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		bulkDeleteErr := repo.BaseRepo.BulkDelete(ctx, ids, bson.M{"deletedAt": bson.M{"$exists": false}})
		if bulkDeleteErr != nil {
			return nil, bulkDeleteErr
		}

		for _, productId := range ids {

			unlinkProductErr := repo.unlinkProductFromProductionLine(ctx, sessCtx, productId, nil)
			if unlinkProductErr != nil {
				return nil, unlinkProductErr
			}
		}
		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

// LinkToProductionLine links a product to a production line
func (repo *ProductRepo) LinkToProductionLine(ctx context.Context, id primitive.ObjectID, linkProductionLineDto dtos.LinkProductionLineDto) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// For vendor users, verify they own the product
	if cfg.User.Type == enums.Vendor {
		// Get the product to check vendorId
		product, getErr := repo.GetById(ctx, id, "")
		if getErr != nil {
			return getErr
		}

		// Vendors can only link their own products
		if product.VendorId != &cfg.User.Id {
			return &helpers.ApiError{
				StatusCode: http.StatusForbidden,
				Message:    "not authorized to update this product",
			}
		}
	}

	// Start a session for the transaction
	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	// Start the transaction
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Get collections
		productCollection, _ := repo.BaseRepo.GetCollection(ctx)
		// Check if product exists and is not deleted
		productFilter := bson.M{
			"_id":       id,
			"deletedAt": bson.M{"$exists": false},
		}

		// Update product with production line ID
		productUpdate := bson.M{
			"$set": bson.M{
				"productionLineId": linkProductionLineDto.ProductionLineId,
				"updatedAt":        time.Now(),
				"updatedBy":        cfg.User.Id,
			},
		}
		updatedOneProudct, updateErr := productCollection.UpdateOne(sessCtx, productFilter, productUpdate)
		if updateErr != nil {
			return nil, helpers.DefaultError(updateErr)
		}

		if updatedOneProudct.MatchedCount == 0 {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Product not found",
			}
		}

		productionLineLinkErr := repo.LinkProductInProductionLine(ctx, sessCtx, id, linkProductionLineDto)
		if productionLineLinkErr != nil {
			return nil, productionLineLinkErr
		}
		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *ProductRepo) LinkProductInProductionLine(ctx context.Context, sessCtx context.Context, productId primitive.ObjectID, linkProductionLineDto dtos.LinkProductionLineDto) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	productionLineCollection := repo.BaseRepo.GetDB().Database(cfg.Db).Collection(enums.PRODUCTION_LINE.GetCollectionName())

	// Check if production line exists and is not deleted
	productionLineUpdateOneFilter := bson.M{
		// "_id":       productionLineId,
		"deletedAt": bson.M{"$exists": false},
	}

	// Update production line with product ID
	productionLineUpdate := bson.M{
		"$addToSet": bson.M{
			"productIds": productId,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}

	updatedOneProductionLine, updateErr := productionLineCollection.UpdateOne(sessCtx, productionLineUpdateOneFilter, productionLineUpdate)
	if updateErr != nil {
		return helpers.DefaultError(updateErr)
	}

	if updatedOneProductionLine.MatchedCount == 0 {
		return &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Production Line not found",
		}
	}

	// Remove Old One
	unlinkProductErr := repo.unlinkProductFromProductionLine(ctx, sessCtx, productId, &linkProductionLineDto.ProductionLineId)
	if unlinkProductErr != nil {
		return unlinkProductErr
	}

	return nil
}

// unlinkProductFromProductionLine unlinks a product from a production line
func (repo *ProductRepo) unlinkProductFromProductionLine(ctx context.Context, sessCtx context.Context, productId primitive.ObjectID, productionLineId *primitive.ObjectID) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)
	productionLineCollection := repo.BaseRepo.GetDB().Database(cfg.Db).Collection(enums.PRODUCTION_LINE.GetCollectionName())

	// Remove Old One
	productionLineUpdateOneFilter := bson.M{
		"productIds": productId,
		"deletedAt":  bson.M{"$exists": false},
	}
	if productionLineId != nil {
		productionLineUpdateOneFilter["_id"] = bson.M{"$ne": productionLineId}
	}
	productionLineUpdateOne := bson.M{
		"$pull": bson.M{
			"productIds": productId,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		},
	}
	_, updateManyErr := productionLineCollection.UpdateOne(sessCtx, productionLineUpdateOneFilter, productionLineUpdateOne)
	if updateManyErr != nil {
		return helpers.DefaultError(updateManyErr)
	}

	return nil
}

func (repo *ProductRepo) linkProductToBranches(ctx context.Context, product *models.Product, branchIds *[]primitive.ObjectID) *helpers.ApiError {
	if branchIds == nil {
		return nil
	}

	cfg, _ := util.GetReqAppCfg(ctx)

	branchProductCollection := base_repos.NewBaseRepo[models.BranchProduct, models.BranchProductResponse](repo.BaseRepo.GetDB(), enums.BRANCH_PRODUCT)
	productBranches := make([]models.BranchProduct, len(*branchIds))
	for i, branchId := range *branchIds {

		branchProduct, err := branchProductCollection.FindOne(ctx, []bson.M{
			{
				"$match": base_utils.AddSoftDeletedFilter(bson.M{
					"branchId":  branchId,
					"productId": product.Id,
				}),
			},
		})
		if err != nil {
			return err
		}

		currencyCode := configuration.EnvironmentVariables.DefaultCurrencyCode
		if branchProduct != nil && branchProduct.OverridePrice != nil {
			currencyCode = branchProduct.OverridePrice.CurrencyCode
		}

		var priceItem *models.PriceItem

		extractedPriceResult, _ := product_shared.ExtractPrice(product_shared.ExtractPriceArgs{
			ProductId:    product.Id,
			Prices:       product.Prices,
			BranchId:     &branchId,
			CurrencyCode: currencyCode,
			UnitId:       nil,
		})

		priceItem = extractedPriceResult.PriceItem

		isActive := true
		dto := dtos.UpsertProductBranchDto{
			ProductId: product.Id,
			BranchId:  branchId,
			BaseLinkBranchProductDto: dtos.BaseLinkBranchProductDto{
				IsActive:   isActive,
				OutOfStock: false,
			},
		}

		if priceItem != nil {
			dto.BaseLinkBranchProductDto.OverridePrice = &dtos.PriceItemDto{
				Value:        priceItem.Value,
				CurrencyCode: priceItem.CurrencyCode,
				UnitId:       priceItem.UnitId,
			}
		}

		if branchProduct != nil {
			dto.BaseLinkBranchProductDto.OverrideDescription = branchProduct.OverrideDescription
			dto.BaseLinkBranchProductDto.IsActive = branchProduct.IsActive
			dto.BaseLinkBranchProductDto.OutOfStock = branchProduct.OutOfStock

		}

		productBranches[i] = product_utils.ConvertToBranchProductModel(dto, types.AdditionalOptions{
			UserId: cfg.User.Id,
			Mode:   types.UPDATE,
		})
	}

	branchProductCollection.BulkDeleteFilter(ctx, bson.M{"productId": product.Id, "branchId": bson.M{"$nin": branchIds}})
	err := repo.branchProductRepo.BulkUpsert(ctx, productBranches)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductRepo) GetAll(ctx context.Context, filters []bson.M) ([]models.ProductResponse, *helpers.ApiError) {
	return repo.BaseRepo.GetAll(ctx, filters)
}

func (repo *ProductRepo) DeleteDuplicatesByName(ctx context.Context) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Start a session for the transaction
	session, err := repo.BaseRepo.GetDB().StartSession()
	if err != nil {
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection, _ := repo.BaseRepo.GetCollection(ctx)

		// Build match stage to filter by vendor if user is vendor
		matchStage := bson.M{
			"deletedAt": bson.M{"$exists": false},
		}

		// Aggregation pipeline to find duplicates by name
		pipeline := []bson.M{
			{"$match": matchStage},
			{
				"$group": bson.M{
					"_id": bson.M{
						"nameEn": "$name.en",
						"nameAr": "$name.ar",
					},
					"products": bson.M{"$push": bson.M{
						"id":        "$_id",
						"createdAt": "$createdAt",
					}},
					"count": bson.M{"$sum": 1},
				},
			},
			{
				"$match": bson.M{
					"count": bson.M{"$gt": 1},
				},
			},
		}

		cursor, err := collection.Aggregate(sessCtx, pipeline)
		if err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to find duplicate products",
			}
		}
		defer cursor.Close(sessCtx)

		var duplicateGroups []struct {
			Id struct {
				NameEn string `bson:"nameEn"`
				NameAr string `bson:"nameAr"`
			} `bson:"_id"`
			Products []struct {
				Id        primitive.ObjectID `bson:"id"`
				CreatedAt time.Time          `bson:"createdAt"`
			} `bson:"products"`
			Count int `bson:"count"`
		}

		if err := cursor.All(sessCtx, &duplicateGroups); err != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusInternalServerError,
				Message:    "failed to decode duplicate products",
			}
		}

		var idsToDelete []primitive.ObjectID

		// For each group of duplicates, keep the oldest and mark others for deletion
		for _, group := range duplicateGroups {
			if len(group.Products) > 1 {
				// Sort by createdAt to find the oldest
				for i := 1; i < len(group.Products); i++ {
					for j := 0; j < len(group.Products)-i; j++ {
						if group.Products[j].CreatedAt.After(group.Products[j+1].CreatedAt) {
							group.Products[j], group.Products[j+1] = group.Products[j+1], group.Products[j]
						}
					}
				}

				// Add all but the first (oldest) to deletion list
				for i := 1; i < len(group.Products); i++ {
					idsToDelete = append(idsToDelete, group.Products[i].Id)
				}
			}
		}

		if len(idsToDelete) > 0 {
			// Soft delete duplicates
			deleteFilter := bson.M{
				"_id": bson.M{"$in": idsToDelete},
			}
			if cfg.User.Type == enums.Vendor {
				deleteFilter["vendorId"] = cfg.User.Id
			}

			update := bson.M{
				"$set": bson.M{
					"deletedAt": time.Now(),
					"updatedAt": time.Now(),
				},
			}

			_, updateErr := collection.UpdateMany(sessCtx, deleteFilter, update)
			if updateErr != nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusInternalServerError,
					Message:    "failed to delete duplicate products",
				}
			}
		}

		return map[string]interface{}{
			"deletedCount": len(idsToDelete),
		}, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *ProductRepo) SaveInventoryProducts(ctx context.Context, saveInventoryProductsDto *dtos.SaveInventoryProductsDto) (*models.ProductResponse, *helpers.ApiError) {
	// Get the configurable product first
	configurableProduct, err := repo.GetById(ctx, saveInventoryProductsDto.ConfigurableProductId, "")
	if err != nil {
		return nil, err
	}

	// Check if product is configurable
	if configurableProduct.ProductType != product_enums.PRODUCT_TYPE_CONFIGURABLE {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Product is not configurable",
		}
	}

	// Extract branch IDs from the configurable product to ensure inventory products
	// are linked to the same branches in the branch-product model
	var branchIds []primitive.ObjectID
	if configurableProduct.Branches != nil {
		for _, branch := range *configurableProduct.Branches {
			branchIds = append(branchIds, branch.Id)
		}
	}

	// Start a session for the transaction
	session, sessionErr := repo.BaseRepo.GetDB().StartSession()
	if sessionErr != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to start transaction",
		}
	}
	defer session.EndSession(ctx)

	var updatedProduct *models.ProductResponse
	_, transactionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Create inventory products
		inventoryProductIds := make([]primitive.ObjectID, 0, len(saveInventoryProductsDto.InventoryProducts))

		for _, inventoryProductDto := range saveInventoryProductsDto.InventoryProducts {
			// Set product type to inventory
			inventoryProductDto.ProductType = product_enums.PRODUCT_TYPE_INVENTORY

			// Ensure inventory products are linked to the same branches as the configurable product
			if len(branchIds) > 0 {
				inventoryProductDto.BranchIds = &branchIds
			}

			var createdInventoryProduct *models.ProductResponse
			var updateErr *helpers.ApiError
			var createErr *helpers.ApiError
			// Create the inventory product
			if inventoryProductDto.Id != nil && *inventoryProductDto.Id != saveInventoryProductsDto.ConfigurableProductId {
				createdInventoryProduct, updateErr = repo.UpdateById(sessCtx, *inventoryProductDto.Id, &inventoryProductDto)
				if updateErr != nil {
					return nil, updateErr
				}
			} else {
				createdInventoryProduct, createErr = repo.Create(sessCtx, &inventoryProductDto)
				if createErr != nil {
					return nil, createErr
				}
			}

			inventoryProductIds = append(inventoryProductIds, createdInventoryProduct.Id)
		}

		// Update the configurable product with inventory product IDs
		updateFilter := bson.M{"_id": saveInventoryProductsDto.ConfigurableProductId}
		updateData := bson.M{
			"$set": bson.M{
				"inventoryProductIds": inventoryProductIds,
				"updatedAt":           time.Now(),
			},
		}

		collection, collectionErr := repo.BaseRepo.GetCollection(sessCtx)
		if collectionErr != nil {
			return nil, collectionErr
		}

		collection.FindOneAndUpdate(sessCtx, updateFilter, updateData)

		// Get the updated configurable product
		updatedProduct, getErr := repo.GetById(sessCtx, saveInventoryProductsDto.ConfigurableProductId, "")
		if getErr != nil {
			return nil, getErr
		}

		return updatedProduct, nil
	})

	if transactionErr != nil {
		if apiErr, ok := transactionErr.(*helpers.ApiError); ok {
			return nil, apiErr
		}
		return nil, helpers.DefaultError(transactionErr)
	}

	return updatedProduct, nil
}

// UpdateAttributeCode updates product attributes array code values from oldCode to newCode.
// If vendorId is provided, the update is scoped to that vendor's products.
func (repo *ProductRepo) UpdateAttributeCode(ctx context.Context, oldCode string, newCode string) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	collection, collErr := repo.BaseRepo.GetCollection(ctx)
	if collErr != nil {
		return collErr
	}

	filter := bson.M{
		"deletedAt":       bson.M{"$exists": false},
		"attributes.code": oldCode,
	}

	// Use arrayFilters to target the matching attribute element
	update := bson.M{
		"$set": bson.M{
			"attributes.$[elem].code": newCode,
			"updatedAt":               time.Now(),
			"updatedBy":               cfg.User.Id,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.code": oldCode},
		},
	}

	_, err := collection.UpdateMany(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return helpers.DefaultError(err)
	}

	return nil
}
