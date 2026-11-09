package cart_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	cart_enums "larsa-e-commerce-microservice-v2/internal/modules/carts/enums"
	order_shared "larsa-e-commerce-microservice-v2/internal/modules/order/shared"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_repos "larsa-e-commerce-microservice-v2/internal/modules/product/repos"
	product_shared "larsa-e-commerce-microservice-v2/internal/modules/product/shared"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	unit_repos "larsa-e-commerce-microservice-v2/internal/modules/unit/repos"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductAddonItemDto struct {
	Id     primitive.ObjectID     `json:"id" validate:"required"`
	Values []ProductAddonValueDto `json:"values" validate:"required,min=1"`
}

type ProductAddonValueDto struct {
	Id       primitive.ObjectID `json:"id" validate:"required"`
	Quantity uint64             `json:"quantity" validate:"min=1"`
}

type UpsertCartItemDto struct {
	ItemId                primitive.ObjectID                     `json:"itemId" validate:"required"`
	Quantity              uint64                                 `json:"quantity" validate:"min=0"`
	Type                  cart_enums.CartItemType                `json:"type" validate:"required,oneof=product"`
	BranchId              *primitive.ObjectID                    `json:"branchId,omitempty"`
	UnitId                *primitive.ObjectID                    `json:"unitId,omitempty"`
	Addons                *[]ProductAddonItemDto                 `json:"addons,omitempty" validate:"omitempty,nonemptysliceptr,dive"`
	AdditionalInformation *order_shared.AdditionalInformationDto `json:"additionalInformation,omitempty"`

	VendorId *primitive.ObjectID `json:"-"`
}

func (dto *UpsertCartItemDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	productRepo := base_repos.NewBaseRepo[product_models.Product, product_models.ProductResponse](args.DB, enums.PRODUCT)
	addonRepo := base_repos.NewBaseRepo[product_models.Addon, product_models.AddonResponse](args.DB, enums.ADDON)
	unitRepo, _ := unit_repos.NewUnitRepo(args.Injector)
	branchProductRepo, _ := product_repos.NewBranchProductRepo(args.Injector)
	var extractProductAttributeKeysResult product_dtos.ExtractProductAttributeKeysResult

	if dto.Type != cart_enums.CartItemTypeProduct {
		if dto.Addons != nil {
			return nil, &helpers.ApiError{
				StatusCode: http.StatusBadRequest,
				Message:    "Addons are only allowed for product items",
			}
		}
	}

	if dto.UnitId != nil {
		_, err := unitRepo.GetById(args.Context, *dto.UnitId, "")
		if err != nil {
			return nil, err
		}
	}

	if dto.Type == cart_enums.CartItemTypeProduct {

		if dto.BranchId != nil {
			branchProduct, err := branchProductRepo.GetByBranchIdAndProductId(args.Context, product_utils.GetByBranchIdAndProductIdArgs{
				BranchId:  *dto.BranchId,
				ProductId: dto.ItemId,
			}, []bson.M{
				{
					"$match": bson.M{
						"isActive": true,
					},
				},
			})
			if err != nil {
				return nil, err
			}

			if branchProduct.OutOfStock {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Product is out of stock",
				}
			}
		}

		productFilter := product_filters.ProductFilters{IncludeAllProductTypes: &[]bool{true}[0]}
		productPipeline := productFilter.BuildPipeline(bson.M{}, args.Context, product_filters.BuildPipelineArgs{
			AddExtraOperations: true,
		})
		product, err := productRepo.GetById(args.Context, dto.ItemId, productPipeline)
		if err != nil {
			return nil, err
		}

		if product.VendorId != nil {
			dto.VendorId = product.VendorId
		}

		if dto.UnitId != nil {
			doesUnitIdExist := util.Find(product.Prices, func(price product_models.Price) bool {
				item := util.Find(price.Items, func(item product_models.PriceItem) bool {
					if item.UnitId == nil {
						return false
					}
					return item.UnitId.Hex() == dto.UnitId.Hex()
				})
				return item != nil
			})

			if doesUnitIdExist == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    "Unit is not allowed to be added to cart because it's not in the product's price list",
					Details: &[]types.ValidationError{
						{
							Message: "Unit with id " + dto.UnitId.Hex() + " is not allowed to be added to cart because it's not in the product's price list",
						},
					},
				}
			}
		}

		if dto.BranchId != nil || dto.Addons != nil {
			res, err := product_dtos.ExtractProductAttributeKeys(args.Context, args.Injector, product.Attributes)
			if err != nil {
				return nil, err
			}
			extractProductAttributeKeysResult = *res
		}

		if dto.BranchId != nil {
			doesBranchExist := util.Find(extractProductAttributeKeysResult.BranchIds, func(branchId primitive.ObjectID) bool {
				return branchId.Hex() == dto.BranchId.Hex()
			})

			if doesBranchExist == nil {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Branch is not allowed to be added to cart",
				}
			}
		}

		// Ensure product has a valid price for the selected branch or default before adding to cart
		currencyCode := configuration.EnvironmentVariables.DefaultCurrencyCode
		_, priceErr := product_shared.ExtractPrice(product_shared.ExtractPriceArgs{
			ProductId:    product.Id,
			Prices:       product.Prices,
			BranchId:     dto.BranchId,
			CurrencyCode: currencyCode,
			UnitId:       dto.UnitId,
		})
		if priceErr != nil {
			return nil, priceErr
		}

		/*
		* Addons validation.
		* 1. Check if addon ids are unique
		* 2. Check if addon values are unique
		* 3. Check if addon that has been sent exist in the database (not soft deleted)
		* 4. Check if addon values that has been sent exist in the database (There's no soft deleted on addon values because it's sub document).
		 */

		if dto.Addons != nil {
			addonIds := []primitive.ObjectID{}
			for _, addon := range *dto.Addons {
				addonIds = append(addonIds, addon.Id)
			}
			uniqueAddonIds := util.ExtractUniqueValues(addonIds)
			if len(uniqueAddonIds) != len(addonIds) {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Some addon ids are duplicated",
				}
			}

			addonValueIds := []primitive.ObjectID{}
			for _, dtoAddon := range *dto.Addons {
				for _, value := range dtoAddon.Values {
					addonValueIds = append(addonValueIds, value.Id)
				}
			}

			uniqueAddonValueIds := util.ExtractUniqueValues(addonValueIds)
			if len(uniqueAddonValueIds) != len(addonValueIds) {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Some addon values are duplicated",
				}
			}

			allowedAddons := util.Filter(uniqueAddonIds, func(addonId primitive.ObjectID, index int) bool {
				doesAddonExist := util.Find(extractProductAttributeKeysResult.AddonIds, func(allowedAddonId primitive.ObjectID) bool {

					return addonId.Hex() == allowedAddonId.Hex()
				})

				return doesAddonExist != nil
			})

			if len(allowedAddons) != len(uniqueAddonIds) {
				return nil, &helpers.ApiError{
					StatusCode: http.StatusBadRequest,
					Message:    "Product have addons that are not allowed to be added to cart",
					Details: &[]types.ValidationError{
						{
							Message: "Product with id " + dto.ItemId.Hex() + " have addons that are not allowed to be added to cart",
						},
					},
				}
			}

			addonFilter := product_filters.AddonFilters{}
			addonPipeline := addonFilter.BuildPipeline(bson.M{}, args.Context, false)
			storedAddons, err := addonRepo.GetByIds(args.Context, uniqueAddonIds, addonPipeline, base_utils.AdditionalOptions{
				ReturnError: true,
			})
			if err != nil {
				return nil, err
			}

			for _, dtoAddon := range *dto.Addons {

				storedAddon := util.Find(storedAddons, func(addon product_models.AddonResponse) bool {
					return addon.Id.Hex() == dtoAddon.Id.Hex()
				})

				if len(dtoAddon.Values) < int(storedAddon.MinSelection) || len(dtoAddon.Values) > int(storedAddon.MaxSelection) {
					return nil, &helpers.ApiError{
						StatusCode: http.StatusBadRequest,
						Message:    "addon and values length must be between " + strconv.Itoa(int(storedAddon.MinSelection)) + " and " + strconv.Itoa(int(storedAddon.MaxSelection)),
						Details: &[]types.ValidationError{
							{
								Message: "addon with id " + storedAddon.Id.Hex() + " and values length must be between " + strconv.Itoa(int(storedAddon.MinSelection)) + " and " + strconv.Itoa(int(storedAddon.MaxSelection)),
							},
						},
					}
				}

				addonValues := util.Filter(dtoAddon.Values, func(dtoValue ProductAddonValueDto, index int) bool {
					storedAddonValue := util.Find(storedAddon.AddonValues, func(addonValue product_models.AddonValue) bool {
						return addonValue.Id.Hex() == dtoValue.Id.Hex()
					})
					return storedAddonValue != nil
				})

				if len(addonValues) != len(dtoAddon.Values) {
					return nil, &helpers.ApiError{
						StatusCode: http.StatusBadRequest,
						Message:    "Some addon values are does not exist",
					}
				}

				for _, addonValue := range addonValues {
					storedAddonValue := util.Find(storedAddon.AddonValues, func(dtoValue product_models.AddonValue) bool {
						return addonValue.Id.Hex() == dtoValue.Id.Hex()
					})

					if storedAddonValue.MaxSelection != nil && addonValue.Quantity > *storedAddonValue.MaxSelection {
						return nil, &helpers.ApiError{
							StatusCode: http.StatusBadRequest,
							Message:    "addon value must be less than " + strconv.Itoa(int(*storedAddonValue.MaxSelection)),
							Details: &[]types.ValidationError{
								{
									Message: "addon value with id " + storedAddonValue.Id.Hex() + " must be less than " + strconv.Itoa(int(*storedAddonValue.MaxSelection)),
								},
							},
						}
					}

					if storedAddonValue.MinSelection != nil && addonValue.Quantity < *storedAddonValue.MinSelection {
						return nil, &helpers.ApiError{
							StatusCode: http.StatusBadRequest,
							Message:    "addon value with id must be greater than " + strconv.Itoa(int(*storedAddonValue.MinSelection)),
							Details: &[]types.ValidationError{
								{
									Message: "addon value with id " + storedAddonValue.Id.Hex() + " must be greater than " + strconv.Itoa(int(*storedAddonValue.MinSelection)),
								},
							},
						}
					}

				}
			}

		}

	} else if dto.Type == cart_enums.CartItemTypeSubscription {
		// TODO: Implement subscription validation
	}

	return nil, nil
}

type BulkUpsertCartItemDto struct {
	Items []UpsertCartItemDto `json:"items" validate:"required,min=1,dive"`
}

func (dto *BulkUpsertCartItemDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	// Validate each item in the bulk request
	for _, item := range dto.Items {
		_, err := item.CustomValidator(args)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
