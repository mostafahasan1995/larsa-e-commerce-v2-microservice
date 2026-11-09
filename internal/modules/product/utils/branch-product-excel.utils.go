package product_utils

import (
	"context"
	"fmt"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetBranchProductImportColumns returns the column definitions for branch product import
func GetBranchProductImportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("productName", "Product Name (English)", excel_services.IntPtr(150)),
		excel_services.NewFloatColumn("overridePrice", "Override Price", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewOptionalStringColumn("overrideDescription", "Override Description", "", excel_services.IntPtr(500)),
		excel_services.NewBoolColumn("outOfStock", "Out of Stock", false, excel_services.BoolPtr(false)),
		excel_services.NewBoolColumn("isActive", "Active", false, excel_services.BoolPtr(true)),
	}
}

// GetBranchProductExportColumns returns the column definitions for branch product export
func GetBranchProductExportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("productName", "Product Name (English)", excel_services.IntPtr(150)),
		excel_services.NewFloatColumn("overridePrice", "Override Price", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewOptionalStringColumn("overrideDescription", "Override Description", "", excel_services.IntPtr(500)),
		excel_services.NewBoolColumn("outOfStock", "Out of Stock", false, excel_services.BoolPtr(false)),
		excel_services.NewBoolColumn("isActive", "Active", false, excel_services.BoolPtr(true)),
	}
}

// ConvertExcelDataToBranchProducts converts Excel import data to LinkBranchToProductDto objects
func ConvertExcelDataToBranchProducts(ctx context.Context, excelData []map[string]interface{}, products []product_models.ProductResponse, branchCurrencyCode string) ([]dtos.LinkBranchToProductDto, *helpers.ApiError) {
	branchProducts := make([]dtos.LinkBranchToProductDto, 0)

	// Create a map for product name to product ID lookup
	productNameToID := make(map[string]primitive.ObjectID)
	for _, product := range products {
		if product.Name != nil && product.Name["en"] != "" {
			productNameToID[product.Name["en"]] = product.Id
		}
	}

	for i, row := range excelData {
		// Get product name - Excel service already validated this as required string
		productName := row["productName"].(string)

		// Business logic validation: check if product exists
		productID, exists := productNameToID[productName]
		if !exists {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    fmt.Sprintf("Row %d: Product '%s' not found", i+1, productName),
			}
		}

		// Create base dto
		dto := dtos.LinkBranchToProductDto{
			ProductId: productID,
		}

		// Get override price - Excel service already validated this as optional float
		if overridePriceVal, exists := row["overridePrice"]; exists && overridePriceVal != nil {
			if overridePriceFloat := overridePriceVal.(float64); overridePriceFloat > 0 {
				dto.OverridePrice = &dtos.PriceItemDto{
					Value:        overridePriceFloat,
					CurrencyCode: branchCurrencyCode,
					UnitId:       nil, // Optional for override price
				}
			}
		}

		// Get override description - Excel service already validated this as optional string
		if overrideDescVal, exists := row["overrideDescription"]; exists && overrideDescVal != nil {
			if overrideDescStr := overrideDescVal.(string); overrideDescStr != "" {
				dto.OverrideDescription = &types.Localizable[string]{
					"en": overrideDescStr,
				}
			}
		}

		// Get out of stock - Excel service already validated this as optional bool
		if outOfStockVal, exists := row["outOfStock"]; exists && outOfStockVal != nil {
			outOfStockBool := outOfStockVal.(bool)
			dto.OutOfStock = outOfStockBool
		}

		// Get is active - Excel service already validated this as optional bool
		if isActiveVal, exists := row["isActive"]; exists && isActiveVal != nil {
			isActiveBool := isActiveVal.(bool)
			dto.IsActive = isActiveBool
		}

		branchProducts = append(branchProducts, dto)
	}

	return branchProducts, nil
}

// ConvertBranchProductsToExcelData converts branch products to Excel export data
func ConvertBranchProductsToExcelData(ctx context.Context, branchProducts []product_models.BranchProductResponse, products []product_models.ProductResponse) ([]map[string]interface{}, *helpers.ApiError) {
	if len(branchProducts) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Create a map for product ID to product lookup
	productIDToProduct := make(map[primitive.ObjectID]product_models.ProductResponse)
	for _, product := range products {
		productIDToProduct[product.Id] = product
	}

	excelData := make([]map[string]interface{}, 0)

	for _, branchProduct := range branchProducts {
		if branchProduct.ProductId == nil {
			continue
		}

		product, exists := productIDToProduct[*branchProduct.ProductId]
		if !exists {
			continue
		}

		row := map[string]interface{}{
			"productName": "",
		}

		// Set product name
		if product.Name != nil && product.Name["en"] != "" {
			row["productName"] = product.Name["en"]
		}

		// Set override price
		if branchProduct.OverridePrice != nil {
			row["overridePrice"] = branchProduct.OverridePrice.Value
		} else {
			row["overridePrice"] = ""
		}

		// Set override description
		if branchProduct.OverrideDescription != nil && (*branchProduct.OverrideDescription)["en"] != "" {
			row["overrideDescription"] = (*branchProduct.OverrideDescription)["en"]
		} else {
			row["overrideDescription"] = ""
		}

		// Set out of stock
		row["outOfStock"] = branchProduct.OutOfStock

		// Set is active
		row["isActive"] = branchProduct.IsActive

		excelData = append(excelData, row)
	}

	return excelData, nil
}

// GetBranchProductExampleData returns example data for the Excel template
func GetBranchProductExampleData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"productName":         "Sample Product 1",
			"overridePrice":       25.99,
			"overrideDescription": "Special branch description for product 1",
			"outOfStock":          false,
			"isActive":            true,
		},
		{
			"productName":         "Sample Product 2",
			"overridePrice":       "",
			"overrideDescription": "",
			"outOfStock":          true,
			"isActive":            false,
		},
	}
}

// ConvertBranchToUpdateDto converts branch data and new products to CreateBranchDto for updating
func ConvertBranchToUpdateDto(branch *product_models.BranchResponse, branchProducts []dtos.LinkBranchToProductDto) dtos.CreateBranchDto {
	// Update the branch with new products using the existing update functionality
	updateDto := dtos.CreateBranchDto{
		Name:                       branch.Name,
		Image:                      branch.Image,
		CountryIso3:                branch.CountryIso3,
		City:                       branch.City,
		Area:                       branch.Area,
		Coordinates:                branch.Coordinates,
		IsActive:                   branch.IsActive,
		PhoneNumber:                branch.PhoneNumber,
		MobileNumber:               branch.MobileNumber,
		BaseCurrencyCode:           branch.BaseCurrencyCode,
		SupportedCurrencyCodes:     branch.SupportedCurrencyCodes,
		RealTimeCurrencyConversion: branch.RealTimeCurrencyConversion,
		CanSetDeliveryOnPickup:     branch.CanSetDeliveryOnPickup,
		MinTimeBeforeOrder:         branch.MinTimeBeforeOrder,
		MaxTimeBeforeOrder:         branch.MaxTimeBeforeOrder,
		Products:                   &branchProducts,
		VendorId:                   branch.VendorId,
	}

	// Convert zones if they exist
	if branch.Zones != nil {
		zoneDtos := make([]dtos.CreateBranchZoneDto, len(*branch.Zones))
		for i, zone := range *branch.Zones {
			zoneDtos[i] = dtos.CreateBranchZoneDto{
				MinOrderAmount: zone.MinOrderAmount,
				DeliveryFee: dtos.PriceItemDto{
					Value:        zone.DeliveryFee.Value,
					CurrencyCode: zone.DeliveryFee.CurrencyCode,
					UnitId:       zone.DeliveryFee.UnitId,
				},
				EstimatedTime: zone.EstimatedTime,
				ZoneId:        zone.ZoneId,
			}
		}
		updateDto.Zones = &zoneDtos
	}

	// Convert opening days if they exist
	if branch.OpeningDays != nil {
		openingDaysDto := product_dtos.OpeningDaysDto{
			Days: make(map[time.Weekday]product_dtos.DayScheduleDto),
		}
		for day, daySchedule := range branch.OpeningDays.Days {
			dayRanges := make([]product_dtos.DateRangeDto, len(daySchedule.TimeRanges))
			for i, dateRange := range daySchedule.TimeRanges {
				dayRanges[i] = product_dtos.DateRangeDto{
					From: dateRange.From,
					To:   dateRange.To,
				}
			}
			openingDaysDto.Days[day] = product_dtos.DayScheduleDto{
				Label:      daySchedule.Label,
				TimeRanges: dayRanges,
			}
		}
		updateDto.OpeningDays = &openingDaysDto
	}

	return updateDto
}
