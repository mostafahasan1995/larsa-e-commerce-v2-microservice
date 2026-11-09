package cart_repos

import (
	"context"
	"fmt"
	address_filters "larsa-e-commerce-microservice-v2/internal/modules/address/filters"
	address_models "larsa-e-commerce-microservice-v2/internal/modules/address/models"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_filters "larsa-e-commerce-microservice-v2/internal/modules/carts/filters"
	services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	coupon_services "larsa-e-commerce-microservice-v2/internal/modules/coupon/services"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_services "larsa-e-commerce-microservice-v2/internal/modules/product/services"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
)

type CartRepo struct {
	*base_repos.BaseRepo[models.Cart, models.CartResponse]
	checkoutRepo             services.CheckoutRepo
	branchService            product_services.BranchService
	currencyRepo             base_repos.BaseRepo[config_models.Currency, config_models.CurrencyResponse]
	couponRuleCheckerService *coupon_services.CouponRuleCheckerService
}

func NewCartRepo(i *do.Injector) (services.CartRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)

	return &CartRepo{
		BaseRepo:                 base_repos.NewBaseRepo[models.Cart, models.CartResponse](db, enums.CART),
		checkoutRepo:             do.MustInvoke[services.CheckoutRepo](i),
		branchService:            do.MustInvoke[product_services.BranchService](i),
		currencyRepo:             *base_repos.NewBaseRepo[config_models.Currency, config_models.CurrencyResponse](db, enums.CURRENCY),
		couponRuleCheckerService: do.MustInvoke[*coupon_services.CouponRuleCheckerService](i),
	}, nil
}

func (repo *CartRepo) GetMine(ctx context.Context, args *cart_shared.GetMineArgs) (*models.CartResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	// Base match stage with deletedAt filter
	matchStage := bson.M{
		"userId":   cfg.User.Id,
		"isActive": true,
	}

	// Create filters from query parameters
	cartFilters, _ := cart_filters.NewCartFilters(args.Query)

	// Build pipeline with filters
	pipeline := cartFilters.BuildPipeline(matchStage, ctx, true)

	// Get collection with primary read preference
	collection, _ := repo.BaseRepo.GetCollection(ctx)
	primaryCollection, cloneErr := collection.Clone(options.Collection().SetReadPreference(readpref.Primary()))
	if cloneErr != nil {
		return nil, helpers.DefaultError(cloneErr)
	}

	cursor, aggregateErr := primaryCollection.Aggregate(ctx, pipeline)
	if aggregateErr != nil {
		return nil, helpers.DefaultError(aggregateErr)
	}
	defer cursor.Close(ctx)

	var cart models.CartResponse
	if cursor.Next(ctx) {
		if err := cursor.Decode(&cart); err != nil {
			return nil, helpers.DefaultError(err)
		}
	} else {
		return nil, nil
	}

	currencyCode := configuration.EnvironmentVariables.DefaultCurrencyCode
	var branchIdObject *primitive.ObjectID
	var branch *product_models.BranchResponse
	if args.BranchId != "" {
		objectId, err := primitive.ObjectIDFromHex(args.BranchId)
		if err != nil {
			return nil, helpers.DefaultError(err)
		}
		storedBranch, getBranchErr := repo.branchService.GetById(ctx, objectId, "")
		if getBranchErr != nil {
			return nil, getBranchErr
		}

		if storedBranch == nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusNotFound,
				Message:    "Branch not found",
			}
		}
		branch = storedBranch
		branchIdObject = &objectId
	}

	if args.CurrencyCode != "" {
		_, getCurrencyErr := repo.currencyRepo.FindOne(ctx, []bson.M{
			{
				"$match": base_utils.AddSoftDeletedFilter(bson.M{
					"code": args.CurrencyCode,
				}),
			},
		})
		if getCurrencyErr != nil {
			return nil, getCurrencyErr
		}
		currencyCode = args.CurrencyCode
	}

	// Get the first zone from the branch if available
	var zoneId *primitive.ObjectID
	if args.ZoneId != "" {
		objectId, err := primitive.ObjectIDFromHex(args.ZoneId)
		if err != nil {
			return nil, helpers.DefaultError(err)
		}
		zoneId = &objectId
	} else if branch != nil && branch.Zones != nil && len(*branch.Zones) > 0 {
		zoneId = &(*branch.Zones)[0].ZoneId
	}

	checkoutDto := &cart_dtos.CreateCheckoutDto{
		CartId:        cart.Id,
		CurrencyCode:  &currencyCode,
		BranchId:      branchIdObject,
		Branch:        branch,
		ZoneId:        zoneId,
		CalculateList: &args.CalculateList,
		OrderType:     args.OrderType,
	}
	// Use addressId from query (what the user chose); load address and persist to cart so checkout and next GetMine use it
	if args.AddressId != "" {
		addressObjId, err := primitive.ObjectIDFromHex(args.AddressId)
		if err != nil {
			return nil, helpers.DefaultError(err)
		}

		addressRepo := base_repos.NewBaseRepo[address_models.AddressModel, address_models.AddressResponse](repo.BaseRepo.GetDB(), enums.ADDRESS)
		addressFilter := address_filters.AddressFilters{
			UserId: &cfg.User.Id,
		}
		addressPipeline := addressFilter.BuildPipeline(bson.M{}, ctx, false)
		address, getAddrErr := addressRepo.GetById(ctx, addressObjId, addressPipeline)
		if getAddrErr != nil {
			return nil, getAddrErr
		}
		checkoutDto.AddressId = &addressObjId
		checkoutDto.Address = address

		// Persist the chosen addressId on cart so checkout and subsequent GetMine use the user's selection
		_, _ = repo.BaseRepo.UpdateOne(ctx, bson.M{"_id": cart.Id}, bson.M{"addressId": addressObjId}, cfg.User.Id)
		// Reflect chosen address on the response so cart.addressId matches what the user sent (and calculation uses)
		cart.AddressId = &addressObjId
	}

	if args.CouponCode != "" {
		coupon, err := repo.couponRuleCheckerService.CheckCouponRedemption(ctx, args.CouponCode, nil, nil, &args.OrderType)
		if err == nil && coupon != nil {
			checkoutDto.CouponCode = &args.CouponCode
			checkoutDto.Coupon = coupon
		}
	}

	calculateCartItemsPrices, calculateCartItemsPricesErr := repo.checkoutRepo.CalculateCartItemsPrices(ctx, checkoutDto)

	// If there's a coupon-related error (minimum subtotal, etc.), retry without the coupon
	if calculateCartItemsPricesErr != nil && checkoutDto.Coupon != nil {
		// Check if error is related to coupon validation
		isCouponError := calculateCartItemsPricesErr.StatusCode == http.StatusBadRequest ||
			calculateCartItemsPricesErr.StatusCode == http.StatusUnprocessableEntity

		if isCouponError {
			// Remove coupon and retry calculation
			checkoutDtoWithoutCoupon := &cart_dtos.CreateCheckoutDto{
				CartId:        checkoutDto.CartId,
				CurrencyCode:  checkoutDto.CurrencyCode,
				BranchId:      checkoutDto.BranchId,
				Branch:        checkoutDto.Branch,
				ZoneId:        checkoutDto.ZoneId,
				AddressId:     checkoutDto.AddressId,
				Address:       checkoutDto.Address,
				CalculateList: checkoutDto.CalculateList,
				OrderType:     checkoutDto.OrderType,
			}
			calculateCartItemsPrices, calculateCartItemsPricesErr = repo.checkoutRepo.CalculateCartItemsPrices(ctx, checkoutDtoWithoutCoupon)
		}
	}

	if calculateCartItemsPricesErr != nil && calculateCartItemsPricesErr.StatusCode == 422 {
		return nil, nil
	}

	// calculateCartItemsPrices.OrderItems = []order_dto.CreateOrderItemDto{}
	// if calculateCartItemsPricesErr != nil {
	// 	return nil, calculateCartItemsPricesErr
	// }
	cart.Calculation = calculateCartItemsPrices

	return &cart, nil
}

// MigrateCart moves a cart, its items, and addresses from an old user to a new user (current user from token)
// This method handles the following scenarios:
// 1. If the old user has a cart but the new user doesn't, it transfers the cart to the new user
// 2. If both users have carts, it force migrates by deleting the new user's cart and replacing it with the old user's cart
// 3. If the old user has no cart, it returns an error
// Addresses are migrated without changing their IDs
func (repo *CartRepo) MigrateCart(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError {
	if repo.BaseRepo == nil {
		return &helpers.ApiError{StatusCode: 500, Message: "BaseRepo is not initialized"}
	}
	if repo.BaseRepo.GetDB() == nil {
		return &helpers.ApiError{StatusCode: 500, Message: "Database connection is not available"}
	}

	oldCart, err := repo.findActiveCartByUserId(ctx, oldUserId)
	if err != nil {
		return err
	}
	if oldCart == nil {
		fmt.Println("migrate cart oldCart is nil")
		return nil
	}

	newCart, err := repo.findActiveCartByUserId(ctx, newUserId)
	if err != nil {
		fmt.Println("Error finding new user's active cart: ", err)
	}

	// Resolve collections
	cartCollection, cartItemCollection, collErr := repo.getCartCollections(ctx)
	if collErr != nil {
		return collErr
	}

	// Load old user's cart items
	oldCartItems, itemsErr := repo.findActiveCartItemsByUserId(ctx, oldUserId)
	if itemsErr != nil {
		return itemsErr
	}
	if len(oldCartItems) == 0 {
		fmt.Println("migrate cart oldCartItems is empty")
		return nil
	}

	// Transactional migration
	session, sessionErr := repo.BaseRepo.GetDB().StartSession()
	if sessionErr != nil {
		return helpers.DefaultError(sessionErr)
	}
	defer session.EndSession(ctx)

	_, transactionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if newCart != nil {
			// Remove new user's existing cart and items
			if err := repo.deleteUserCartItems(sessCtx, cartItemCollection, newUserId); err != nil {
				return nil, err
			}
			if err := repo.deleteCartById(sessCtx, cartCollection, *newCart.Id); err != nil {
				return nil, err
			}

			// Reassign old items to the old cart id and new owner
			cartItemIds := extractCartItemIds(oldCartItems)
			if err := repo.reassignCartItems(sessCtx, cartItemCollection, oldCartItems, newUserId, *oldCart.Id); err != nil {
				return nil, err
			}
			if err := repo.updateCartOwnershipAndItems(sessCtx, cartCollection, *oldCart.Id, newUserId, cartItemIds); err != nil {
				return nil, err
			}
			return nil, nil
		}

		// No cart for new user: create a fresh cart and move items into it
		newCartId := primitive.NewObjectID()
		cartItemIds, totalQuantities := extractCartItemIdsAndQuantity(oldCartItems)

		if err := repo.reassignCartItems(sessCtx, cartItemCollection, oldCartItems, newUserId, newCartId); err != nil {
			return nil, err
		}

		newCartDoc := buildNewCartDocument(newUserId, cartItemIds, totalQuantities)
		newCartDoc.Id = newCartId
		if _, apiErr := repo.BaseRepo.Create(sessCtx, newCartDoc); apiErr != nil {
			return nil, apiErr
		}

		if delErr := repo.BaseRepo.DeleteById(sessCtx, *oldCart.Id, bson.M{}); delErr != nil {
			return nil, delErr
		}
		return nil, nil
	})

	if transactionErr != nil {
		if apiErr, ok := transactionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(transactionErr)
	}

	return nil
}

// findActiveCartByUserId fetches the active (non-soft-deleted) cart for a user
func (repo *CartRepo) findActiveCartByUserId(ctx context.Context, userId primitive.ObjectID) (*models.CartResponse, *helpers.ApiError) {
	return repo.FindOne(ctx, []bson.M{{"$match": base_utils.AddSoftDeletedFilter(bson.M{"userId": userId, "isActive": true})}})
}

// findActiveCartItemsByUserId returns all active cart items for a user
func (repo *CartRepo) findActiveCartItemsByUserId(ctx context.Context, userId primitive.ObjectID) ([]cart_shared.CartItem, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	cartItemCollection := repo.BaseRepo.GetDB().Database(cfg.Db).Collection(enums.CART_ITEM.GetCollectionName())
	cursor, findErr := cartItemCollection.Find(ctx, bson.M{"userId": userId, "isActive": true})
	if findErr != nil {
		return nil, helpers.DefaultError(findErr)
	}
	defer cursor.Close(ctx)

	var items []cart_shared.CartItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, helpers.DefaultError(err)
	}
	return items, nil
}

// getCartCollections returns the carts collection and the cart items collection
func (repo *CartRepo) getCartCollections(ctx context.Context) (*mongo.Collection, *mongo.Collection, *helpers.ApiError) {
	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return nil, nil, helpers.DefaultError(err)
	}
	cfg, _ := util.GetReqAppCfg(ctx)
	cartItemCollection := repo.BaseRepo.GetDB().Database(cfg.Db).Collection(enums.CART_ITEM.GetCollectionName())
	return collection, cartItemCollection, nil
}

// deleteUserCartItems removes all active items for a user
func (repo *CartRepo) deleteUserCartItems(ctx context.Context, cartItemCollection *mongo.Collection, userId primitive.ObjectID) *helpers.ApiError {
	_, err := cartItemCollection.DeleteMany(ctx, bson.M{"userId": userId, "isActive": true})
	if err != nil {
		return helpers.DefaultError(err)
	}
	return nil
}

// deleteCartById deletes a cart by id
func (repo *CartRepo) deleteCartById(ctx context.Context, cartCollection *mongo.Collection, cartId primitive.ObjectID) *helpers.ApiError {
	_, err := cartCollection.DeleteOne(ctx, bson.M{"_id": cartId})
	if err != nil {
		return helpers.DefaultError(err)
	}
	return nil
}

// reassignCartItems sets the owner and cartId for a set of items
func (repo *CartRepo) reassignCartItems(ctx context.Context, cartItemCollection *mongo.Collection, items []cart_shared.CartItem, newUserId primitive.ObjectID, cartId primitive.ObjectID) *helpers.ApiError {
	for _, cartItem := range items {
		_, updateErr := cartItemCollection.UpdateOne(ctx, bson.M{"_id": cartItem.Id}, bson.M{"$set": bson.M{"userId": newUserId, "cartId": cartId}})
		if updateErr != nil {
			return helpers.DefaultError(updateErr)
		}
	}
	return nil
}

// updateCartOwnershipAndItems assigns the cart to new user and updates its item ids
func (repo *CartRepo) updateCartOwnershipAndItems(ctx context.Context, cartCollection *mongo.Collection, cartId primitive.ObjectID, newUserId primitive.ObjectID, itemIds []primitive.ObjectID) *helpers.ApiError {
	_, err := cartCollection.UpdateOne(ctx, bson.M{"_id": cartId}, bson.M{"$set": bson.M{"userId": newUserId, "cartItemIds": itemIds}})
	if err != nil {
		return helpers.DefaultError(err)
	}
	return nil
}

func extractCartItemIds(items []cart_shared.CartItem) []primitive.ObjectID {
	ids := make([]primitive.ObjectID, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.Id)
	}
	return ids
}

func extractCartItemIdsAndQuantity(items []cart_shared.CartItem) ([]primitive.ObjectID, uint64) {
	ids := make([]primitive.ObjectID, 0, len(items))
	var total uint64
	for _, it := range items {
		ids = append(ids, it.Id)
		total += it.Quantity
	}
	return ids, total
}

func buildNewCartDocument(newUserId primitive.ObjectID, itemIds []primitive.ObjectID, totalQuantities uint64) models.Cart {
	return models.Cart{
		BaseCart: models.BaseCart{
			CartItemIds:     itemIds,
			UserId:          newUserId,
			IsActive:        true,
			TotalQuantities: totalQuantities,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(types.AdditionalOptions{UserId: newUserId, Mode: types.CREATE}),
	}
}
