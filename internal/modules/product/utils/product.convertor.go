package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

func convertToProductAttributes(candidateAttributes []dtos.AttributeDto) []models.ProductAttribute {
	// Convert attributes from DTO to model
	attributes := make([]models.ProductAttribute, len(candidateAttributes))
	for i, attrDto := range candidateAttributes {
		attributes[i] = models.ProductAttribute{
			Code:  attrDto.Code,
			Value: attrDto.Value,
			Type:  attrDto.Type,
		}
	}

	return attributes
}

func ConvertToPriceItems(priceItems []dtos.PriceItemDto) []models.PriceItem {
	priceValues := make([]models.PriceItem, len(priceItems))
	for i, priceItem := range priceItems {
		priceValues[i] = models.PriceItem{
			CurrencyCode: priceItem.CurrencyCode,
			Value:        priceItem.Value,
			UnitId:       priceItem.UnitId,
		}
	}
	return priceValues
}

func convertToDiscount(discount *dtos.DiscountDto) *models.Discount {
	if discount == nil {
		return nil
	}
	return &models.Discount{
		Value:        discount.Value,
		CurrencyCode: discount.CurrencyCode,
		Type:         discount.Type,
	}
}

// ConvertToProductModel converts a CreateProductDto to a Product model
func ConvertToProductModel(createProductDto *dtos.CreateProductDto, additionalOptions types.AdditionalOptions) *models.Product {

	// Convert prices from DTO to model
	prices := make([]models.Price, len(createProductDto.Prices))
	for i, priceDto := range createProductDto.Prices {
		prices[i] = models.Price{
			BasePrice: models.BasePrice{
				Name:           priceDto.Name,
				BranchIds:      priceDto.BranchIds,
				Items:          ConvertToPriceItems(priceDto.Items),
				CustomerGroups: priceDto.CustomerGroups,
				Platforms:      priceDto.Platforms,
				Discount:       convertToDiscount(priceDto.Discount),
				IsDefault:      priceDto.IsDefault,
			},
		}
	}

	productType := product_enums.PRODUCT_TYPE_SIMPLE
	//Don't change the product type if it's inventory else it will be auto detected
	if createProductDto.ProductType != "" && createProductDto.ProductType == product_enums.PRODUCT_TYPE_INVENTORY {
		productType = createProductDto.ProductType
	} else {
		// Auto-detect product type based on attributes
		for _, attr := range createProductDto.Attributes {
			if attr.Type == product_enums.PRODUCT_ATTRIBUTE_TYPE_ARRAY {

				// Check if the attribute is configurable, if not, don't change the product type
				attribute := product_models.AttributeResponse{}
				for _, storedAttr := range *createProductDto.StoreAttributes {
					if storedAttr.Code == attr.Code {
						attribute = storedAttr
						break
					}
				}

				if !attribute.Configurable {
					continue
				}

				// If the attribute is configurable, check if the value is an array and has more than 1 item
				if attr.Code != string(product_enums.SYSTEM_BRANCH) {
					if arr, ok := attr.Value.([]interface{}); ok && len(arr) > 1 {
						productType = product_enums.PRODUCT_TYPE_CONFIGURABLE
						break
					}
				}
			}
		}
	}

	product := &models.Product{
		BaseProduct: models.BaseProduct{
			VendorId:          createProductDto.VendorId,
			CategoryIds:       createProductDto.CategoryIds,
			UnitIds:           createProductDto.UnitIds,
			Attributes:        convertToProductAttributes(createProductDto.Attributes),
			Name:              createProductDto.Name,
			Description:       createProductDto.Description,
			Images:            createProductDto.Images,
			Prices:            prices,
			Flag:              createProductDto.Flag,
			ProductType:       productType,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}
	return product
}
