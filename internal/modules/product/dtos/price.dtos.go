package product_dtos

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	base_repos "larsa-e-commerce-microservice-v2/internal/modules/base/repos"
	base_utils "larsa-e-commerce-microservice-v2/internal/modules/base/utils"
	config_models "larsa-e-commerce-microservice-v2/internal/modules/config/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	unit_models "larsa-e-commerce-microservice-v2/internal/modules/unit/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiscountDto struct {
	Value        float64                    `json:"value" validate:"omitempty,min=0"`
	CurrencyCode *string                    `json:"currencyCode" validate:"omitempty"`
	Type         product_enums.DiscountType `json:"type" validate:"required,oneof=percentage fixed"`
}

type PriceItemDto struct {
	Value        float64             `json:"value" validate:"required"`
	CurrencyCode string              `json:"currencyCode" validate:"required"`
	UnitId       *primitive.ObjectID `json:"unitId" validate:"omitempty"`
}

type CreatePriceDto struct {
	Name           *types.Localizable[string]     `json:"name" validate:"omitempty"`
	BranchIds      *[]primitive.ObjectID          `json:"branchIds" validate:"omitempty,min=0,dive"`
	Items          []PriceItemDto                 `json:"items" validate:"required,min=1,dive,required"`
	CustomerGroups *[]product_enums.CustomerGroup `json:"customerGroups,omitempty" validate:"omitempty,dive,oneof=loyalty engaged"`
	Platforms      *[]product_enums.PricePlatform `json:"platforms,omitempty" validate:"omitempty,dive,oneof=web mobile pos"`
	Discount       *DiscountDto                   `json:"discount,omitempty" validate:"omitempty"`
	IsDefault      *bool                          `json:"isDefault,omitempty" validate:"omitempty"`
}

func (dto *CreatePriceDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	if dto.IsDefault != nil && *dto.IsDefault {
		dto.BranchIds = nil
		dto.CustomerGroups = nil
		dto.Platforms = nil
	} else {
		isDefault := false
		dto.IsDefault = &isDefault

	}

	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	currencyRepo := base_repos.NewBaseRepo[config_models.Currency, config_models.CurrencyResponse](args.DB, enums.CURRENCY)
	unitRepo := base_repos.NewBaseRepo[unit_models.Unit, unit_models.UnitResponse](args.DB, enums.UNIT)

	if dto.BranchIds != nil {
		_, err := branchRepo.GetByIds(args.Context, *dto.BranchIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
		if err != nil {
			return nil, err
		}
	}

	unitIds := make([]primitive.ObjectID, 0)
	currencyCodes := make([]string, 0)
	for _, item := range dto.Items {
		if item.UnitId != nil {
			unitIds = append(unitIds, *item.UnitId)
			currencyCodes = append(currencyCodes, item.CurrencyCode)
		}
	}
	_, err := unitRepo.GetByIds(args.Context, unitIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
	if err != nil {
		return nil, err
	}

	currencies, err := currencyRepo.GetAll(args.Context, []bson.M{
		{
			"$match": bson.M{
				"code": bson.M{
					"$in": util.ExtractUniqueValues(currencyCodes),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(currencies) != len(util.ExtractUniqueValues(currencyCodes)) {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "Some currencies are not found",
		}
	}

	return nil, nil
}
