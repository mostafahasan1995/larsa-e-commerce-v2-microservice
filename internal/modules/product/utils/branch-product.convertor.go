package product_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

func ConvertToBranchProductModel(createBranchProductDto dtos.UpsertProductBranchDto, additionalOptions types.AdditionalOptions) models.BranchProduct {

	branchProduct := models.BranchProduct{
		BranchId:          createBranchProductDto.BranchId,
		ProductId:         createBranchProductDto.ProductId,
		MetaDataBaseModel: util.CreateMetaDataBaseModel(additionalOptions),
		BaseBranchProduct: models.BaseBranchProduct{
			OverrideDescription: createBranchProductDto.OverrideDescription,
			IsActive:            createBranchProductDto.IsActive,
			OutOfStock:          createBranchProductDto.OutOfStock,
		},
	}

	if createBranchProductDto.OverridePrice != nil {
		branchProduct.BaseBranchProduct.OverridePrice = &models.PriceItem{
			Value:        createBranchProductDto.OverridePrice.Value,
			CurrencyCode: createBranchProductDto.OverridePrice.CurrencyCode,
			UnitId:       createBranchProductDto.OverridePrice.UnitId,
		}
	}

	return branchProduct
}
