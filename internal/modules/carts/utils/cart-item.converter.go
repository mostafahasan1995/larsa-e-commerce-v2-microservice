package cart_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	cart_dtos "larsa-e-commerce-microservice-v2/internal/modules/carts/dtos"
	cart_shared "larsa-e-commerce-microservice-v2/internal/modules/carts/shared"
)

func ConvertToCartItemModel(dto cart_dtos.UpsertCartItemDto, additionalOptions types.AdditionalOptions) cart_shared.CartItem {

	cartItem := cart_shared.CartItem{
		BaseCartItem: cart_shared.BaseCartItem{
			ItemId:   dto.ItemId,
			Quantity: dto.Quantity,
			Type:     dto.Type,
			UserId:   additionalOptions.UserId,
			VendorId: dto.VendorId,
			BranchId: dto.BranchId,
			UnitId:   dto.UnitId,
			IsActive: true,
		},
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
	}

	if dto.Addons != nil {
		addons := convertToCartAddonModel(*dto.Addons)
		cartItem.Addons = &addons
	}

	if dto.AdditionalInformation != nil {
		cartItem.AdditionalInformation = &cart_shared.AdditionalInformation{
			Image: dto.AdditionalInformation.Image,
			Logo:  dto.AdditionalInformation.Logo,
			Note:  dto.AdditionalInformation.Note,
		}
	}

	return cartItem
}

func convertToCartAddonModel(cartAddon []cart_dtos.ProductAddonItemDto) []cart_shared.CartAddon {
	addonModels := []cart_shared.CartAddon{}
	for _, addon := range cartAddon {
		addonModels = append(addonModels, cart_shared.CartAddon{
			Id:     addon.Id,
			Values: convertToCartAddonValueModel(addon.Values),
		})
	}

	return addonModels
}

func convertToCartAddonValueModel(cartAddonValue []cart_dtos.ProductAddonValueDto) []cart_shared.CartAddonValue {
	cartAddonValueModels := []cart_shared.CartAddonValue{}
	for _, addonValue := range cartAddonValue {
		cartAddonValueModels = append(cartAddonValueModels, cart_shared.CartAddonValue{
			Id:       addonValue.Id,
			Quantity: addonValue.Quantity,
		})
	}

	return cartAddonValueModels
}
