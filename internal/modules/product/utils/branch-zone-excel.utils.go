package product_utils

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetBranchZoneImportColumns returns the column definitions for branch zone import
func GetBranchZoneImportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("zoneName", "Zone Name", excel_services.IntPtr(100)),
		excel_services.NewFloatColumn("minOrderAmount", "Min Order Amount", true, excel_services.Float64Ptr(0), nil),
		excel_services.NewFloatColumn("deliveryFee", "Delivery Fee", true, excel_services.Float64Ptr(0), nil),
		excel_services.NewRequiredStringColumn("currencyCode", "Currency Code", excel_services.IntPtr(10)),
		excel_services.NewFloatColumn("estimatedTime", "Estimated Delivery Time (mins)", true, excel_services.Float64Ptr(0), nil),
	}
}

// ConvertExcelDataToBranchZones converts Excel import data to BranchZone objects
func ConvertExcelDataToBranchZones(ctx context.Context, excelData []map[string]interface{}, zoneService config_services.ZoneService) ([]product_models.BranchZone, *helpers.ApiError) {
	branchZones := make([]product_models.BranchZone, 0, len(excelData))

	// Get all zones for name lookup
	zones, apiErr := zoneService.GetAll(ctx, `{"isActive":true}`)
	if apiErr != nil {
		return nil, apiErr
	}

	// Create a map for zone name to zone ID lookup
	zoneNameToID := make(map[string]primitive.ObjectID)
	for _, zone := range zones {
		zoneNameToID[zone.Name] = zone.Id
	}

	for i, row := range excelData {
		// Get zone name - Excel service already validated this as required string
		zoneName := row["zoneName"].(string)

		// Business logic validation: check if zone exists
		zoneID, exists := zoneNameToID[zoneName]
		if !exists {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    fmt.Sprintf("Row %d: Zone '%s' not found", i+1, zoneName),
			}
		}

		// Get values - Excel service already validated and converted these
		minOrderAmount := row["minOrderAmount"].(float64)
		deliveryFee := row["deliveryFee"].(float64)
		currencyCode := row["currencyCode"].(string)
		estimatedTime := row["estimatedTime"].(float64)

		// Create BranchZone
		branchZone := product_models.BranchZone{
			ZoneId:         zoneID,
			MinOrderAmount: minOrderAmount,
			DeliveryFee: product_models.PriceItem{
				Value:        deliveryFee,
				CurrencyCode: currencyCode,
				UnitId:       nil, // Optional for delivery fee
			},
			EstimatedTime: uint64(estimatedTime),
		}

		branchZones = append(branchZones, branchZone)
	}

	return branchZones, nil
}

// GetBranchZoneExampleData returns example data for the Excel template
func GetBranchZoneExampleData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"zoneName":       "Zone A",
			"minOrderAmount": 50.0,
			"deliveryFee":    10.0,
			"currencyCode":   "USD",
			"estimatedTime":  30.0,
		},
		{
			"zoneName":       "Zone B",
			"minOrderAmount": 100.0,
			"deliveryFee":    15.0,
			"currencyCode":   "AED",
			"estimatedTime":  45.0,
		},
	}
}

// GetBranchZoneExportColumns returns the column definitions for branch zone export
func GetBranchZoneExportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("zoneName", "Zone Name", excel_services.IntPtr(100)),
		excel_services.NewFloatColumn("minOrderAmount", "Min Order Amount", true, excel_services.Float64Ptr(0), nil),
		excel_services.NewFloatColumn("deliveryFee", "Delivery Fee", true, excel_services.Float64Ptr(0), nil),
		excel_services.NewRequiredStringColumn("currencyCode", "Currency Code", excel_services.IntPtr(10)),
		excel_services.NewFloatColumn("estimatedTime", "Estimated Delivery Time (mins)", true, excel_services.Float64Ptr(0), nil),
	}
}

// ConvertBranchZonesToExcelData converts branch zones to Excel export data
func ConvertBranchZonesToExcelData(ctx context.Context, branchZones []product_models.BranchZone, zoneService config_services.ZoneService) ([]map[string]interface{}, *helpers.ApiError) {
	if len(branchZones) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Get all zones for name lookup
	zones, apiErr := zoneService.GetAll(ctx, `{"isActive":true}`)
	if apiErr != nil {
		return nil, apiErr
	}

	// Create a map for zone ID to zone name lookup
	zoneIDToName := make(map[primitive.ObjectID]string)
	for _, zone := range zones {
		zoneIDToName[zone.Id] = zone.Name
	}

	excelData := make([]map[string]interface{}, 0)

	for _, branchZone := range branchZones {
		zoneName, exists := zoneIDToName[branchZone.ZoneId]
		if !exists {
			continue
		}

		row := map[string]interface{}{
			"zoneName":       zoneName,
			"minOrderAmount": branchZone.MinOrderAmount,
			"deliveryFee":    branchZone.DeliveryFee.Value,
			"currencyCode":   branchZone.DeliveryFee.CurrencyCode,
			"estimatedTime":  float64(branchZone.EstimatedTime),
		}

		excelData = append(excelData, row)
	}

	return excelData, nil
}
