package cart_repos

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	cart_utils "larsa-e-commerce-microservice-v2/internal/modules/carts/utils"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	utils "larsa-e-commerce-microservice-v2/internal/modules/carts/utils"
)

type CartItemRepo struct {
	*base_repos.BaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse]
	cartRepo *base_repos.BaseRepo[models.Cart, models.CartResponse]
}

func NewCartItemRepo(i *do.Injector) (services.CartItemRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &CartItemRepo{
		BaseRepo: base_repos.NewBaseRepo[cart_shared.CartItem, cart_shared.CartItemResponse](db, enums.CART_ITEM),
		cartRepo: base_repos.NewBaseRepo[models.Cart, models.CartResponse](db, enums.CART),
	}, nil
}

type HandleCartItemResult struct {
	StoredCart cart_models.CartResponse
	CartItems  []cart_shared.CartItemResponse
}

func (repo *CartItemRepo) HandleCartItem(sessCtx context.Context, upsertCartItemDto *dtos.UpsertCartItemDto) (*HandleCartItemResult, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(sessCtx)

	storedCart, err := repo.cartRepo.FindOne(sessCtx, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(bson.M{
				"userId":   cfg.User.Id,
				"isActive": true,
			}),
		},
	})
	if err != nil {
		return nil, err
	}

	cartItemsFilter := bson.M{
		"userId":   cfg.User.Id,
		"isActive": true,
	}

	if storedCart != nil {
		cartItemsFilter["cartId"] = *storedCart.Id
	}

	cartItems, err := repo.BaseRepo.GetAll(sessCtx, []bson.M{
		{
			"$match": base_utils.AddSoftDeletedFilter(cartItemsFilter),
		},
	})
	if err != nil {
		return nil, err
	}

	result := &HandleCartItemResult{
		CartItems: cartItems,
	}

	if storedCart != nil {
		result.StoredCart = *storedCart
	}

	return result, nil
}

func (repo *CartItemRepo) Create(ctx context.Context, upsertCartItemDto *dtos.UpsertCartItemDto) *helpers.ApiError {

	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		if upsertCartItemDto.Quantity == 0 {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Quantity must be greater than 0",
			}
		}

		handleCartItemResult, handleCartItemErr := repo.HandleCartItem(ctx, upsertCartItemDto)
		if handleCartItemErr != nil {
			return nil, handleCartItemErr
		}

		var storedCart *models.CartResponse
		if handleCartItemResult.StoredCart.Id != nil {
			storedCart = &handleCartItemResult.StoredCart
		}

		repo.handleCreation(ctx, HandleCreationArgs{
			UpsertCartItemDto: upsertCartItemDto,
			CartItems:         handleCartItemResult.CartItems,
			StoredCart:        storedCart,
		})

		return nil, nil

	}, nil)

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

func (repo *CartItemRepo) UpdateById(ctx context.Context, cartItemId primitive.ObjectID, upsertCartItemDto *dtos.UpsertCartItemDto) *helpers.ApiError {

	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)

	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		storedCartItem, getErr := repo.BaseRepo.GetById(sessCtx, cartItemId, []bson.M{
			{
				"$match": base_utils.AddSoftDeletedFilter(bson.M{
					"isActive": true,
				}),
			},
		})
		if getErr != nil {
			return nil, getErr
		}

		handleCartItemResult, handleCartItemErr := repo.HandleCartItem(ctx, upsertCartItemDto)
		if handleCartItemErr != nil {
			return nil, handleCartItemErr
		}

		repo.handleUpdate(ctx, HandleUpdateArgs{
			UpsertCartItemDto: upsertCartItemDto,
			CartItems:         handleCartItemResult.CartItems,
			CartItem:          *storedCartItem,
		})

		return nil, nil

	}, nil)

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return helpers.DefaultError(sessionErr)
	}

	return nil
}

type HandleCreationArgs struct {
	UpsertCartItemDto *dtos.UpsertCartItemDto
	CartItems         []cart_shared.CartItemResponse
	StoredCart        *models.CartResponse
}

func (repo *CartItemRepo) handleCreation(sessCtx context.Context, args HandleCreationArgs) *helpers.ApiError {

	cfg, _ := util.GetReqAppCfg(sessCtx)

	var cartId = primitive.NewObjectID()
	if args.StoredCart != nil {
		cartId = *args.StoredCart.Id
	}

	cartItem := utils.ConvertToCartItemModel(*args.UpsertCartItemDto, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.CREATE,
	})
	cartItem.CartId = cartId
	createdCartItem, err := repo.BaseRepo.Create(sessCtx, cartItem)
	if err != nil {
		return err
	}

	cartItems := append(args.CartItems, *createdCartItem)

	mode := types.CREATE
	if args.StoredCart != nil {
		mode = types.UPDATE
	}
	cart := utils.ConvertToCartModel(cart_utils.ConvertToCartModelArgs{
		Dto:    cartItems,
		CartId: cartId,
	}, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   mode,
	})

	if args.StoredCart == nil {
		repo.cartRepo.Create(sessCtx, cart)
	} else {
		repo.cartRepo.UpdateById(sessCtx, *args.StoredCart.Id, cart, bson.M{})
	}

	return nil
}

type HandleUpdateArgs struct {
	UpsertCartItemDto *dtos.UpsertCartItemDto
	CartItems         []cart_shared.CartItemResponse
	CartItem          cart_shared.CartItemResponse
}

func (repo *CartItemRepo) handleUpdate(sessCtx context.Context, args HandleUpdateArgs) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(sessCtx)
	collection, err := repo.BaseRepo.GetCollection(sessCtx)
	if err != nil {
		return err
	}

	// If the quantity is 0 that means the user wants to delete the cart item
	if args.UpsertCartItemDto.Quantity == 0 {
		_, _ = collection.DeleteOne(sessCtx, bson.M{
			"_id":    args.CartItem.Id,
			"userId": cfg.User.Id,
		})
	} else {
		cartItem := utils.ConvertToCartItemModel(*args.UpsertCartItemDto, types.AdditionalOptions{
			UserId: cfg.User.Id,
			Mode:   types.UPDATE,
		})
		cartItem.CartId = args.CartItem.CartId
		_, err := repo.BaseRepo.UpdateById(sessCtx, *args.CartItem.Id, cartItem, bson.M{})
		if err != nil {
			return err
		}
	}

	cart := utils.ConvertToCartModel(cart_utils.ConvertToCartModelArgs{
		Dto:    args.CartItems,
		CartId: args.CartItem.CartId,
	}, types.AdditionalOptions{
		UserId: cfg.User.Id,
		Mode:   types.UPDATE,
	})

	cartFilter := base_utils.AddSoftDeletedFilter(bson.M{
		"userId":   cfg.User.Id,
		"isActive": true,
	})
	updateField :=
		bson.M{
			"cartItemIds":     cart.CartItemIds,
			"totalQuantities": cart.TotalQuantities,
			"updatedAt":       cart.UpdatedAt,
			"updatedBy":       cart.UpdatedBy,
			"userId":          cart.UserId,
		}
	repo.cartRepo.UpdateOne(sessCtx, cartFilter, updateField, cfg.User.Id)

	return nil
}

func (repo *CartItemRepo) Clear(ctx context.Context) *helpers.ApiError {
	cfg, _ := util.GetReqAppCfg(ctx)

	collection, err := repo.BaseRepo.GetCollection(ctx)
	if err != nil {
		return err
	}

	session, _ := repo.BaseRepo.GetDB().StartSession()
	defer session.EndSession(ctx)
	_, sessionErr := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {

		cart, err := repo.cartRepo.FindOne(sessCtx, []bson.M{
			{
				"$match": base_utils.AddSoftDeletedFilter(bson.M{
					"userId":   cfg.User.Id,
					"isActive": true,
				}),
			},
		})
		if err != nil {
			return nil, err
		}
		if cart == nil {
			return nil, nil
		}

		collection.DeleteMany(sessCtx, bson.M{
			"cartId": *cart.Id,
			"userId": cfg.User.Id,
		})

		repo.cartRepo.UpdateOne(sessCtx, bson.M{
			"_id":      *cart.Id,
			"isActive": true,
		}, bson.M{
			"cartItemIds": []bson.M{},
		}, cfg.User.Id)

		return nil, nil
	})

	if sessionErr != nil {
		if apiErr, ok := sessionErr.(*helpers.ApiError); ok {
			return apiErr
		}
		return &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to checkout",
		}
	}

	return nil
}
