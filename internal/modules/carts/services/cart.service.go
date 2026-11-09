package cart_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	address_services "larsa-e-commerce-microservice-v2/internal/modules/address/services"
	cart_models "larsa-e-commerce-microservice-v2/internal/modules/carts/models"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepo interface {
	GetMine(ctx context.Context, args *cart_shared.GetMineArgs) (*cart_models.CartResponse, *helpers.ApiError)
	MigrateCart(ctx context.Context, oldUserId, newUserId primitive.ObjectID) *helpers.ApiError
}

type CartService interface {
	GetMine(ctx context.Context, args *cart_shared.GetMineArgs) (*cart_models.CartResponse, *helpers.ApiError)
	MigrateCartAndAddresses(ctx context.Context, oldUserId primitive.ObjectID) *helpers.ApiError
}

type cart struct {
	repo           CartRepo
	addressService address_services.AddressService
}

func NewCartService(i *do.Injector) (CartService, error) {
	return &cart{
		repo:           do.MustInvoke[CartRepo](i),
		addressService: do.MustInvoke[address_services.AddressService](i),
	}, nil
}

func (p *cart) GetMine(ctx context.Context, args *cart_shared.GetMineArgs) (*cart_models.CartResponse, *helpers.ApiError) {
	return p.repo.GetMine(ctx, args)
}

func (p *cart) MigrateCartAndAddresses(ctx context.Context, oldUserId primitive.ObjectID) *helpers.ApiError {
	cfg, cfgErr := util.GetReqAppCfg(ctx)
	if cfgErr != nil {
		return &helpers.ApiError{StatusCode: 500, Message: "Failed to get user context"}
	}

	newUserId := cfg.User.Id
	migrateAddressErr := p.addressService.MigrateAddresses(ctx, oldUserId, newUserId)
	if migrateAddressErr != nil {
		fmt.Println("migrateAddressErr", migrateAddressErr)
		return migrateAddressErr
	}

	migrateCartErr := p.repo.MigrateCart(ctx, oldUserId, newUserId)
	if migrateCartErr != nil {
		fmt.Println("migrateCartErr", migrateCartErr)
		return migrateCartErr
	}
	return nil
}
