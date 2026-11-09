package product_dtos

import (
	"encoding/json"
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

type AttributeDto struct {
	Code  string                             `json:"code" validate:"required"`
	Value interface{}                        `json:"value" validate:"required"`
	Type  product_enums.ProductAttributeType `json:"type" validate:"required,oneof=number string boolean array object"`
}

type CreateProductDto struct {
	Name        types.Localizable[string]  `json:"name" validate:"required,localizable"`
	Description types.Localizable[string]  `json:"description" validate:"required,localizable"`
	CategoryIds []primitive.ObjectID       `json:"categoryIds,omitempty" validate:"omitempty,min=1,dive,required"`
	UnitIds     []primitive.ObjectID       `json:"unitIds,omitempty" validate:"omitempty,min=1,dive,required"`
	Attributes  []AttributeDto             `json:"attributes,omitempty" validate:"omitempty,min=0,dive"`
	Flag        *product_enums.ProductFlag `json:"flag,omitempty" validate:"omitempty,oneof=hot_pick trending innovative new most_popular low_moq oem"`
	Images      []types.FileFieldDto       `json:"images,omitempty" validate:"omitempty,dive,required"`
	Prices      []CreatePriceDto           `json:"prices" validate:"required,min=1,dive,required"`
	VendorId    *primitive.ObjectID        `json:"vendorId,omitempty" validate:"omitempty"`
	ProductType product_enums.ProductType  `json:"productType,omitempty" validate:"omitempty"`

	BranchIds       *[]primitive.ObjectID               `json:"-"`
	AddonIds        *[]primitive.ObjectID               `json:"-"`
	StoreAttributes *[]product_models.AttributeResponse `json:"-"`
	Id              *primitive.ObjectID                 `json:"id,omitempty"`
}

func (dto *CreateProductDto) CustomValidator(args middleware.CustomValidatorArgs) ([]middleware.CustomValidatorContext, *helpers.ApiError) {
	isDefaults := []bool{}
	for _, price := range dto.Prices {
		if price.IsDefault != nil && *price.IsDefault {
			isDefaults = append(isDefaults, true)
		}

	}

	if len(isDefaults) > 1 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Only one default price is allowed",
		}
	}

	unitRepo := base_repos.NewBaseRepo[unit_models.Unit, unit_models.UnitResponse](args.DB, enums.UNIT)
	categoryRepo := base_repos.NewBaseRepo[config_models.Category, config_models.CategoryResponse](args.DB, enums.CATEGORY)
	branchRepo := base_repos.NewBaseRepo[product_models.Branch, product_models.BranchResponse](args.DB, enums.BRANCH)
	attributeRepo := base_repos.NewBaseRepo[product_models.Attribute, product_models.AttributeResponse](args.DB, enums.ATTRIBUTE)

	if err := checkProductAttributeSize(dto.Attributes); err != nil {
		return nil, err
	}

	err := ValidateProductAttributes(dto.Attributes)
	if err != nil {
		return nil, err
	}

	if len(dto.Attributes) > 0 {
		productAttribute := []product_models.ProductAttribute{}

		for _, attribute := range dto.Attributes {
			productAttribute = append(productAttribute, product_models.ProductAttribute{
				Code:  attribute.Code,
				Value: attribute.Value,
				Type:  attribute.Type,
			})
		}
		extractProductAttributeKeysResult, err := ExtractProductAttributeKeys(args.Context, args.Injector, productAttribute)
		if err != nil {
			return nil, err
		}

		dto.Attributes = util.Map(dto.Attributes, func(dtoAttribute AttributeDto) AttributeDto {
			switch dtoAttribute.Code {
			case string(product_enums.SYSTEM_BRANCH):
				dtoAttribute.Type = product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY
				dtoAttribute.Value = extractProductAttributeKeysResult.BranchIds
			case string(product_enums.SYSTEM_ADDON):
				dtoAttribute.Type = product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY
				dtoAttribute.Value = extractProductAttributeKeysResult.AddonIds
			}
			return dtoAttribute
		})

		if extractProductAttributeKeysResult.BranchIds != nil {
			dto.BranchIds = &extractProductAttributeKeysResult.BranchIds

			_, err = branchRepo.GetByIds(args.Context, *dto.BranchIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
			if err != nil {
				return nil, err
			}
		}

		if extractProductAttributeKeysResult.AddonIds != nil {
			dto.AddonIds = &extractProductAttributeKeysResult.AddonIds
		}

	}

	unitIds := util.ExtractUniqueValues(dto.UnitIds)
	categoryIds := util.ExtractUniqueValues(dto.CategoryIds)

	_, err = unitRepo.GetByIds(args.Context, unitIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
	if err != nil {
		return nil, err
	}

	_, err = categoryRepo.GetByIds(args.Context, categoryIds, []bson.M{}, base_utils.AdditionalOptions{ReturnError: true})
	if err != nil {
		return nil, err
	}

	attributeCodes := util.ExtractUniqueValues(util.Map(dto.Attributes, func(attribute AttributeDto) string {
		return attribute.Code
	}))
	storedAttributes, err := attributeRepo.GetAll(args.Context, []bson.M{
		{
			"$match": bson.M{
				"code": bson.M{
					"$in": attributeCodes,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	dto.StoreAttributes = &storedAttributes

	return nil, nil
}

// LinkProductionLineDto represents the data needed to link a product to a production line
type LinkProductionLineDto struct {
	ProductionLineId primitive.ObjectID `json:"productionLineId" validate:"required"`
}

/*
* check if any product attribute value has a size more than 50KB then return false
* if the size is less than 50KB then return true
 */
func checkProductAttributeSize(attributes []AttributeDto) *helpers.ApiError {
	const maxSize = 25 * 1024 // 50KB in bytes
	for _, attribute := range attributes {
		data, err := json.Marshal(attribute.Value)
		if err != nil {
			// If marshaling fails, treat as exceeding size (or you could choose to skip or handle differently)
			return helpers.DefaultError(err)
		}
		if len(data) > maxSize {
			return &helpers.ApiError{
				StatusCode: http.StatusRequestEntityTooLarge,
				Message:    "Product attribute value size is too large",
				Details: &[]types.ValidationError{
					{
						Field:   attribute.Code,
						Message: "Product attribute value size is too large",
					},
				},
			}
		}
	}

	return nil
}

type SaveInventoryProductsDto struct {
	ConfigurableProductId primitive.ObjectID `json:"configurableProductId" validate:"required"`
	InventoryProducts     []CreateProductDto `json:"inventoryProducts" validate:"required,min=1,dive,required"`
}
