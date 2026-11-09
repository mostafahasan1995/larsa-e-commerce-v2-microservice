package cart_handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_services "larsa-e-commerce-microservice-v2/internal/modules/carts/services"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
	product_handlers "larsa-e-commerce-microservice-v2/internal/modules/product/handlers"
)

type CartItemHandler struct {
	cartService        cart_services.CartService
	cartItemService    cart_services.CartItemService
	validationInstance *validator.Validate
}

func NewCartItemHandler(i *do.Injector, r *chi.Mux) {
	h := CartItemHandler{
		cartService:        do.MustInvoke[cart_services.CartService](i),
		cartItemService:    do.MustInvoke[cart_services.CartItemService](i),
		validationInstance: do.MustInvoke[*validator.Validate](i),
	}

	r.Route(enums.CART_ITEM.GetWebRoute(), func(r chi.Router) {
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(cart_dtos.UpsertCartItemDto{}, middleware.BODY, i),
		).Post("/", helpers.Make(h.Create))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(cart_dtos.BulkUpsertCartItemDto{}, middleware.BODY, i),
		).Post("/bulk", helpers.Make(h.CreateBulk))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(types.GetByIdDto{}, middleware.PARAM, i),
			middleware.ValidateDTO(cart_dtos.UpsertCartItemDto{}, middleware.BODY, i),
		).Put("/{id}", helpers.Make(h.UpdateById))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			})).Get("/mine", helpers.Make(h.GetMine))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			})).Delete("/mine", helpers.Make(h.Clear))
		r.With(
			middleware.AdminAuth(middleware.Restriction{
				BuyerCapability: enums.CapabilitiesRegistry.General,
			}),
			middleware.ValidateDTO(cart_dtos.MigrateCartDto{}, middleware.BODY, i),
		).Post("/migrate", helpers.Make(h.Migrate))

	})
}

func (h *CartItemHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[cart_dtos.UpsertCartItemDto](r)
	if err != nil {
		return err
	}

	createErr := h.cartItemService.Create(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, createErr)
		return nil
	}

	cart, err := h.getMine(r)
	if err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}
	response, _ := json.MarshalContext(ctx, cart)
	w.Write(response)

	return nil
}

func (h *CartItemHandler) CreateBulk(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[cart_dtos.BulkUpsertCartItemDto](r)
	if err != nil {
		return err
	}

	createErr := h.cartItemService.CreateBulk(ctx, dto)
	if createErr != nil {
		helpers.HandleApiError(w, err)
		return nil
	}

	cart, err := h.getMine(r)
	if err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}
	response, _ := json.MarshalContext(ctx, cart)
	w.Write(response)

	return nil
}

func (h *CartItemHandler) UpdateById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	id := chi.URLParam(r, "id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	dto, err := middleware.GetValidatedDTO[cart_dtos.UpsertCartItemDto](r)
	if err != nil {
		return err
	}

	updateErr := h.cartItemService.UpdateById(ctx, objectId, dto)
	if updateErr != nil {
		helpers.HandleApiError(w, updateErr)
		return nil
	}

	cart, err := h.getMine(r)
	if err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}
	response, _ := json.MarshalContext(ctx, cart)
	w.Write(response)

	return nil
}

func (h *CartItemHandler) GetMine(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	cart, err := h.getMine(r)
	if err != nil {
		helpers.HandleApiError(w, err)
		return nil
	}
	response, _ := json.MarshalContext(ctx, cart)
	w.Write(response)

	return nil
}

func (h *CartItemHandler) Clear(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	clearErr := h.cartItemService.Clear(ctx)
	if clearErr != nil {
		helpers.HandleApiError(w, clearErr)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *CartItemHandler) getMine(r *http.Request) (*cart_models.CartResponse, *helpers.ApiError) {

	query := r.URL.Query().Get("query")
	couponCode := r.URL.Query().Get("couponCode")
	currencyCode := r.URL.Query().Get("currencyCode")
	if currencyCode == "" {
		currencyCode = configuration.EnvironmentVariables.DefaultCurrencyCode
	}
	branchId := r.URL.Query().Get("branchId")
	zoneId := r.URL.Query().Get("zoneId")
	addressId := r.URL.Query().Get("addressId")
	calculateList := r.URL.Query().Get("calculatePriceList") == "true"
	orderType := r.URL.Query().Get("orderType")

	var orderTypeEnum order_enums.OrderType = order_enums.PICKUP
	if orderType != "" {
		orderTypeEnum = order_enums.OrderType(orderType)
	}
	ctx, _ := util.AddCtxAppCfg(r)

	cart, getErr := h.cartService.GetMine(ctx, &cart_shared.GetMineArgs{
		Query:         query,
		CurrencyCode:  currencyCode,
		BranchId:      branchId,
		AddressId:     addressId,
		CouponCode:    couponCode,
		CalculateList: calculateList,
		ZoneId:        zoneId,
		OrderType:     orderTypeEnum,
	})
	if getErr != nil {
		return nil, getErr
	}

	if cart == nil {
		return &cart_models.CartResponse{
			CartItems: []cart_shared.CartItemResponse{},
			BaseCart: cart_models.BaseCart{
				CartItemIds:     []primitive.ObjectID{},
				UserId:          primitive.NilObjectID,
				IsActive:        false,
				TotalQuantities: 0,
			},
			Calculation: &cart_shared.CalculateCartItemsPricesResponse{},
		}, nil
	}

	for _, cartItem := range cart.CartItems {
		unitId := ""
		if cartItem.UnitId != nil {
			unitId = cartItem.UnitId.Hex()
		}
		if cartItem.Product != nil {
			price := product_handlers.GetPrice(ctx, product_handlers.GetPriceArgs{
				Product:       *cartItem.Product,
				BranchId:      branchId,
				CurrencyCode:  currencyCode,
				UnitId:        unitId,
				CalculateList: calculateList,
			})
			cartItem.Product.Price = price.Price
		}
	}

	return cart, nil
}

// Migrate handles the cart and address migration endpoint
// It accepts an oldUserId in the request body and force migrates the cart and addresses from that user to the current user
// If the current user already has a cart, it will be deleted and replaced with the old user's cart
// Addresses are migrated without changing their IDs
func (h *CartItemHandler) Migrate(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := util.AddCtxAppCfg(r)

	dto, err := middleware.GetValidatedDTO[cart_dtos.MigrateCartDto](r)
	if err != nil {
		fmt.Println("error getting validatedDTO migrate cart", err)
		return err
	}

	fmt.Println("migate dto result", dto)

	migrateErr := h.cartService.MigrateCartAndAddresses(ctx, dto.OldUserId)
	fmt.Println("migate cart and address result", migrateErr)
	if migrateErr != nil {
		fmt.Println("migrateErr", migrateErr)
		helpers.HandleApiError(w, migrateErr)
		return nil
	}

	fmt.Println("success migrate cart and address", migrateErr)

	response := map[string]string{
		"message": "Cart and addresses migrated successfully",
	}
	responseBytes, marshalErr := json.MarshalContext(ctx, response)
	if marshalErr != nil {
		fmt.Println("error marshaling response", marshalErr)
		helpers.HandleApiError(w, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to marshal response",
		})
		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return nil
}
