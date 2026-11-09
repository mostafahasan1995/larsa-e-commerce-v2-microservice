package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

// ConvertToProductionLineModel converts a CreateProductionLineDto to a ProductionLine model
func ConvertToProductionLineModel(createProductionLineDto *dtos.CreateProductionLineDto, additionalOptions types.AdditionalOptions) *models.ProductionLine {

	productionLine := &models.ProductionLine{
		BaseProductionLine: models.BaseProductionLine{

			VendorId:          createProductionLineDto.VendorId,
			Attributes:        convertToProductAttributes(createProductionLineDto.Attributes),
			Name:              createProductionLineDto.Name,
			Description:       createProductionLineDto.Description,
			CategoryIds:       createProductionLineDto.CategoryIds,
			MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		},
	}

	return productionLine
}
