package product_services

import (
	"context"
	"encoding/json"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	product_filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	"reflect"
	"strings"

	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryProductService interface {
	GenerateInventoryProducts(ctx context.Context, configurableProduct *product_models.ProductResponse) ([]product_models.ProductResponse, *helpers.ApiError)
	IsConfigurableProduct(product *product_models.ProductResponse) bool
	CreateAndSaveInventoryProducts(ctx context.Context, configurableProduct *product_models.ProductResponse) ([]product_models.ProductResponse, *helpers.ApiError)
}

type inventoryProductService struct {
	productRepo          ProductRepo
	attributeRepo        AttributeRepo
	branchProductService BranchProductService
}

func NewInventoryProductService(i *do.Injector) (InventoryProductService, error) {
	return &inventoryProductService{
		productRepo:          do.MustInvoke[ProductRepo](i),
		attributeRepo:        do.MustInvoke[AttributeRepo](i),
		branchProductService: do.MustInvoke[BranchProductService](i),
	}, nil
}

func (s *inventoryProductService) IsConfigurableProduct(product *product_models.ProductResponse) bool {
	for _, attr := range product.Attributes {
		if attr.Type == product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY {
			if arr, ok := attr.Value.([]interface{}); ok && len(arr) > 1 {
				return true
			}
		}
	}
	return false
}

func (s *inventoryProductService) GenerateInventoryProducts(ctx context.Context, configurableProduct *product_models.ProductResponse) ([]product_models.ProductResponse, *helpers.ApiError) {
	configurableAttributes, err := s.getConfigurableAttributes(ctx, configurableProduct.Attributes)
	if err != nil {
		return nil, err
	}

	if len(configurableAttributes) == 0 {
		return []product_models.ProductResponse{}, nil
	}

	combinations := s.generateCombinations(configurableAttributes)

	inventoryProducts := make([]product_models.ProductResponse, 0, len(combinations))

	for _, combination := range combinations {
		inventoryProduct := s.createInventoryProduct(configurableProduct, combination)
		inventoryProducts = append(inventoryProducts, inventoryProduct)
	}

	return inventoryProducts, nil
}

func (s *inventoryProductService) getConfigurableAttributes(ctx context.Context, attributes []product_models.ProductAttribute) ([]product_models.ProductAttribute, *helpers.ApiError) {
	var configurableAttrs []product_models.ProductAttribute
	var storedAttrs []product_models.AttributeResponse

	attributeCodes := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		attributeCodes = append(attributeCodes, attr.Code)
	}

	attributeFilters := product_filters.AttributeFilters{
		Codes: &attributeCodes,
	}
	attributeFiltersString, err := json.Marshal(attributeFilters)
	if err != nil {
		return nil, helpers.DefaultError(err)
	}

	storedAttrs, getAllErr := s.attributeRepo.GetAll(ctx, string(attributeFiltersString))
	if getAllErr != nil {
		return nil, helpers.DefaultError(getAllErr)
	}

	for _, attr := range attributes {
		if attr.Type == product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY {
			storedAttr := product_models.AttributeResponse{}

			//Check if the attribute is not configurable skip it
			for _, _storedAttr := range storedAttrs {
				if _storedAttr.Code == attr.Code {
					storedAttr = _storedAttr
					break
				}
			}
			if !storedAttr.Configurable {
				continue
			}

			//
			if attr.Code != "system_branch" {
				if arr, ok := attr.Value.(primitive.A); ok {
					if len(arr) > 1 {
						configurableAttrs = append(configurableAttrs, attr)
					}
				} else {
					fmt.Println("attr.Value type", reflect.TypeOf(attr.Value), "attr.Value", attr.Value)
				}
			}
		}
	}

	return configurableAttrs, nil
}

func (s *inventoryProductService) generateCombinations(attributes []product_models.ProductAttribute) []map[string]interface{} {
	if len(attributes) == 0 {
		return []map[string]interface{}{}
	}

	attributeValues := make([][]interface{}, len(attributes))
	for i, attr := range attributes {
		if arr, ok := attr.Value.(primitive.A); ok {
			attributeValues[i] = arr
		}
	}

	return s.cartesianProduct(attributeValues, attributes)
}

// cartesianProduct generates all possible combinations of product attribute values.
// This function is used to create product variants by combining different attribute values.
//
// Algorithm explanation:
//  1. Base case: If no values, return empty result
//  2. Single attribute: Create maps with each value for the single attribute
//  3. Multiple attributes: Recursively generate combinations for remaining attributes,
//     then combine each value of the first attribute with each sub-combination
//
// Example: If we have:
// - Color: ["Red", "Blue"]
// - Size: ["S", "M"]
// - Material: ["Cotton"]
//
// Result will be:
// [
//
//	{"color": "Red", "size": "S", "material": "Cotton"},
//	{"color": "Red", "size": "M", "material": "Cotton"},
//	{"color": "Blue", "size": "S", "material": "Cotton"},
//	{"color": "Blue", "size": "M", "material": "Cotton"}
//
// ]
func (s *inventoryProductService) cartesianProduct(values [][]interface{}, attributes []product_models.ProductAttribute) []map[string]interface{} {
	// Base case: if no values provided, return empty result
	if len(values) == 0 {
		return []map[string]interface{}{}
	}

	// Base case: if only one attribute, create maps with each value
	if len(values) == 1 {
		result := make([]map[string]interface{}, len(values[0]))
		for i, value := range values[0] {
			result[i] = map[string]interface{}{
				attributes[0].Code: value,
			}
		}
		return result
	}

	// Recursive case: generate combinations for remaining attributes
	subCombinations := s.cartesianProduct(values[1:], attributes[1:])

	// Combine each value of the first attribute with each sub-combination
	result := make([]map[string]interface{}, 0, len(values[0])*len(subCombinations))

	for _, value := range values[0] {
		for _, subCombination := range subCombinations {
			// Create a new combination map
			combination := make(map[string]interface{})
			// Add the current attribute value
			combination[attributes[0].Code] = value
			// Merge with sub-combination values
			for k, v := range subCombination {
				combination[k] = v
			}
			result = append(result, combination)
		}
	}

	return result
}

func (s *inventoryProductService) createInventoryProduct(baseProduct *product_models.ProductResponse, attributeCombination map[string]interface{}) product_models.ProductResponse {
	inventoryProduct := *baseProduct

	inventoryProduct.ProductType = product_enums.PRODUCT_TYPE_INVENTORY

	inventoryProduct.Name = s.generateInventoryProductName(baseProduct.Name, attributeCombination)

	inventoryProduct.Attributes = s.setInventoryProductAttributes(baseProduct.Attributes, attributeCombination)

	inventoryProduct.InventoryProducts = nil

	return inventoryProduct
}

func (s *inventoryProductService) generateInventoryProductName(baseName types.Localizable[string], attributeCombination map[string]interface{}) types.Localizable[string] {
	var suffixes []string
	for _, value := range attributeCombination {
		suffix := fmt.Sprintf("%v", value)
		suffixes = append(suffixes, suffix)
	}

	attributeSuffix := strings.Join(suffixes, ", ")

	result := make(types.Localizable[string])
	for lang, name := range baseName {
		result[lang] = fmt.Sprintf("%s (%s)", name, attributeSuffix)
	}

	return result
}

func (s *inventoryProductService) setInventoryProductAttributes(baseAttributes []product_models.ProductAttribute, combination map[string]interface{}) []product_models.ProductAttribute {
	result := make([]product_models.ProductAttribute, 0, len(baseAttributes))

	for _, attr := range baseAttributes {
		if value, exists := combination[attr.Code]; exists {
			result = append(result, product_models.ProductAttribute{
				Code:  attr.Code,
				Value: value,
				Type:  product_enums.PRODUCT_ATTRIBUTE_TYPE_STRING,
			})
		} else {
			result = append(result, attr)
		}
	}

	return result
}

func (s *inventoryProductService) CreateAndSaveInventoryProducts(ctx context.Context, configurableProduct *product_models.ProductResponse) ([]product_models.ProductResponse, *helpers.ApiError) {
	// Generate inventory products
	inventoryProducts, err := s.GenerateInventoryProducts(ctx, configurableProduct)
	if err != nil {
		return nil, err
	}

	if len(inventoryProducts) == 0 {
		return []product_models.ProductResponse{}, nil
	}

	// Convert inventory products to DTOs for creation
	createdInventoryProducts := make([]product_models.ProductResponse, 0, len(inventoryProducts))

	for _, inventoryProduct := range inventoryProducts {
		// Convert to CreateProductDto
		createDto := s.convertToCreateProductDto(&inventoryProduct, configurableProduct)

		// Create the inventory product
		createdProduct, createErr := s.productRepo.Create(ctx, createDto)
		if createErr != nil {
			return nil, createErr
		}

		createdInventoryProducts = append(createdInventoryProducts, *createdProduct)
	}

	return createdInventoryProducts, nil
}

func (s *inventoryProductService) convertToCreateProductDto(inventoryProduct *product_models.ProductResponse, configurableProduct *product_models.ProductResponse) *dtos.CreateProductDto {
	// Convert attributes
	attributes := make([]dtos.AttributeDto, len(inventoryProduct.Attributes))
	for i, attr := range inventoryProduct.Attributes {
		attributes[i] = dtos.AttributeDto{
			Code:  attr.Code,
			Value: attr.Value,
			Type:  attr.Type,
		}
	}

	// Convert prices (inherit from configurable product)
	prices := make([]dtos.CreatePriceDto, len(configurableProduct.Prices))
	for i, price := range configurableProduct.Prices {
		priceItems := make([]dtos.PriceItemDto, len(price.Items))
		for j, item := range price.Items {
			priceItems[j] = dtos.PriceItemDto{
				Value:        item.Value,
				CurrencyCode: item.CurrencyCode,
				UnitId:       item.UnitId,
			}
		}

		prices[i] = dtos.CreatePriceDto{
			Name:           price.Name,
			BranchIds:      price.BranchIds,
			Items:          priceItems,
			CustomerGroups: price.CustomerGroups,
			Platforms:      price.Platforms,
			IsDefault:      price.IsDefault,
		}
	}

	// Extract branch IDs from configurable product's branches
	var branchIds []primitive.ObjectID
	if configurableProduct.Branches != nil {
		for _, branch := range *configurableProduct.Branches {
			branchIds = append(branchIds, branch.Id)
		}
	}

	return &dtos.CreateProductDto{
		Name:        inventoryProduct.Name,
		Description: inventoryProduct.Description,
		CategoryIds: configurableProduct.CategoryIds,
		UnitIds:     configurableProduct.UnitIds,
		Attributes:  attributes,
		Flag:        configurableProduct.Flag,
		Images:      inventoryProduct.Images,
		Prices:      prices,
		VendorId:    configurableProduct.VendorId,
		BranchIds:   &branchIds,
	}
}
