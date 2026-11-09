package product_services

import (
	"context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"
	dtos "larsa-e-commerce-microservice-v2/internal/modules/product/dtos"
	filters "larsa-e-commerce-microservice-v2/internal/modules/product/filters"
	models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
	product_utils "larsa-e-commerce-microservice-v2/internal/modules/product/utils"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchRepo interface {
	Create(ctx context.Context, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, branchFilters filters.BranchFilters) (*models.BranchResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string, branchFilters filters.BranchFilters) ([]models.BranchResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string, branchFilters filters.BranchFilters) (*models.BranchesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	UpdateZones(ctx context.Context, id primitive.ObjectID, zones []models.BranchZone) (*models.BranchResponse, *helpers.ApiError)
	LinkAllProductsToTheFirstBranch(ctx context.Context) *helpers.ApiError
}

type BranchService interface {
	Create(ctx context.Context, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError)
	GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BranchResponse, *helpers.ApiError)
	GetAll(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError)
	GetAllV2(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError)
	GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BranchesWithPagination, *helpers.ApiError)
	UpdateById(ctx context.Context, id primitive.ObjectID, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError)
	DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError
	BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError
	ImportZones(ctx context.Context, branchId primitive.ObjectID, file *multipart.FileHeader) *helpers.ApiError
	GetZoneImportExample(ctx context.Context) ([]byte, *helpers.ApiError)
	ImportProducts(ctx context.Context, branchId primitive.ObjectID, file *multipart.FileHeader) *helpers.ApiError
	ExportProducts(ctx context.Context, branchId primitive.ObjectID) ([]byte, *helpers.ApiError)
	ExportZones(ctx context.Context, branchId primitive.ObjectID) ([]byte, *helpers.ApiError)
	GetProductImportExample(ctx context.Context) ([]byte, *helpers.ApiError)
	LinkAllProductsToTheFirstBranch(ctx context.Context) *helpers.ApiError
}

type branch struct {
	repo           BranchRepo
	zoneService    config_services.ZoneService
	productService ProductService
	excelService   excel_services.ExcelService
	rolesGw        *gateway.RolesGw
}

func NewBranchService(i *do.Injector) (BranchService, error) {
	return &branch{
		repo:           do.MustInvoke[BranchRepo](i),
		zoneService:    do.MustInvoke[config_services.ZoneService](i),
		productService: do.MustInvoke[ProductService](i),
		excelService:   do.MustInvoke[excel_services.ExcelService](i),
		rolesGw:        do.MustInvoke[*gateway.RolesGw](i),
	}, nil
}

func (p *branch) Create(ctx context.Context, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError) {
	createdBranch, apiErr := p.repo.Create(ctx, branch)
	if apiErr != nil {
		return nil, apiErr
	}

	// groupedOpeningDays := product_utils.GroupDatesSpinningAcrossMultipleDays(createdBranch.TimeZone.UtcOffset, createdBranch.OpeningDays)
	// createdBranch.OpeningDays = groupedOpeningDays

	return createdBranch, nil
}

func (p *branch) GetById(ctx context.Context, id primitive.ObjectID, query string) (*models.BranchResponse, *helpers.ApiError) {
	branchFilters, err := filters.NewBranchFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}

	storedBranch, apiErr := p.repo.GetById(ctx, id, *branchFilters)
	if apiErr != nil {
		return nil, apiErr
	}

	// groupedOpeningDays := product_utils.GroupDatesSpinningAcrossMultipleDays(storedBranch.TimeZone.UtcOffset, storedBranch.OpeningDays)
	// storedBranch.OpeningDays = groupedOpeningDays

	return storedBranch, nil
}

func (p *branch) GetAll(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	branchFilters, err := filters.NewBranchFilters(query)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}

	// userCapabilities := cfg.User.UserData.Capabilities
	//hasGetMyBranches := slices.Contains(userCapabilities, "ecommerceV2GetMyBranches")

	isBranchManager, isBranchManagerError := p.rolesGw.IsUserHaveRole(ctx, "branchManager")
	if isBranchManagerError != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to check if user is branch manager",
		}
	}

	if isBranchManager {
		branchFilters.ManagerId = &cfg.User.Id
	}

	if cfg != nil && cfg.User != nil && cfg.User.Type == enums.Vendor {
		branchFilters.VendorId = &cfg.User.Id
	}

	storedBranches, apiErr := p.repo.GetAll(ctx, query, *branchFilters)
	if apiErr != nil {
		return nil, apiErr
	}

	// for _, storedBranch := range storedBranches {
	// 	groupedOpeningDays := product_utils.GroupDatesSpinningAcrossMultipleDays(storedBranch.TimeZone.UtcOffset, storedBranch.OpeningDays)
	// 	storedBranch.OpeningDays = groupedOpeningDays
	// }

	return storedBranches, nil
}
func (p *branch) GetAllV2(ctx context.Context, query string) ([]models.BranchResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	branchFilters, err := filters.NewBranchFilters(query)

	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}

	if cfg != nil && cfg.User != nil && cfg.User.Type == enums.Vendor {
		branchFilters.VendorId = &cfg.User.Id
	}

	storedBranches, apiErr := p.repo.GetAll(ctx, query, *branchFilters)
	if apiErr != nil {
		return nil, apiErr
	}

	// for _, storedBranch := range storedBranches {
	// 	groupedOpeningDays := product_utils.GroupDatesSpinningAcrossMultipleDays(storedBranch.TimeZone.UtcOffset, storedBranch.OpeningDays)
	// 	storedBranch.OpeningDays = groupedOpeningDays
	// }

	return storedBranches, nil
}

func (p *branch) GetWithPagination(ctx context.Context, skip, limit int64, query string) (*models.BranchesWithPagination, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)
	branchFilters, err := filters.NewBranchFilters(query)

	isBranchManager, isBranchManagerError := p.rolesGw.IsUserHaveRole(ctx, "branchManager")
	fmt.Println("isBranchManagerError", isBranchManagerError)
	if isBranchManager {
		branchFilters.ManagerId = &cfg.User.Id
	}

	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query parameters",
		}
	}

	if cfg != nil && cfg.User != nil && cfg.User.Type == enums.Vendor {
		branchFilters.VendorId = &cfg.User.Id
	}

	return p.repo.GetWithPagination(ctx, skip, limit, query, *branchFilters)
}

func (p *branch) UpdateById(ctx context.Context, id primitive.ObjectID, branch dtos.CreateBranchDto) (*models.BranchResponse, *helpers.ApiError) {
	updatedBranch, apiErr := p.repo.UpdateById(ctx, id, branch)
	if apiErr != nil {
		return nil, apiErr
	}

	// groupedOpeningDays := product_utils.GroupDatesSpinningAcrossMultipleDays(updatedBranch.TimeZone.UtcOffset, updatedBranch.OpeningDays)
	// updatedBranch.OpeningDays = groupedOpeningDays

	return updatedBranch, nil
}

func (p *branch) DeleteById(ctx context.Context, id primitive.ObjectID) *helpers.ApiError {
	return p.repo.DeleteById(ctx, id)
}

func (p *branch) BulkDelete(ctx context.Context, ids []primitive.ObjectID) *helpers.ApiError {
	return p.repo.BulkDelete(ctx, ids)
}

func (p *branch) ImportZones(ctx context.Context, branchId primitive.ObjectID, file *multipart.FileHeader) *helpers.ApiError {
	// Check file format
	if strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "CSV files are not supported. Please upload an Excel (.xlsx) file or use the provided Excel template.",
		}
	}

	// Get the branch to validate it exists and get its currency
	_, apiErr := p.GetById(ctx, branchId, "")
	if apiErr != nil {
		return apiErr
	}

	// Create Excel service
	excelService := p.excelService

	// Define columns for import
	columns := product_utils.GetBranchZoneImportColumns()

	// Import data from Excel
	importedData, apiErr := excelService.Import(ctx, file, columns)
	if apiErr != nil {
		return apiErr
	}

	// Check if import validation was successful
	if !importedData.IsValid {
		errors := make([]string, len(importedData.Errors))
		for i, err := range importedData.Errors {
			errors[i] = err.Error
		}
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid data",
		}
	}

	// Convert Excel data to BranchZone objects
	branchZones, apiErr := product_utils.ConvertExcelDataToBranchZones(ctx, importedData.Data, p.zoneService)
	if apiErr != nil {
		return apiErr
	}

	// Update the branch with new zones
	_, apiErr = p.repo.UpdateZones(ctx, branchId, branchZones)
	if apiErr != nil {
		return apiErr
	}

	return nil
}

func (p *branch) GetZoneImportExample(ctx context.Context) ([]byte, *helpers.ApiError) {
	// Create Excel service
	excelService := p.excelService

	// Define columns for import
	columns := product_utils.GetBranchZoneImportColumns()

	// Generate example Excel file
	excelData, apiErr := excelService.GetExample(ctx, columns, "Zone Import Template")
	if apiErr != nil {
		return nil, apiErr
	}

	return excelData, nil
}

func (p *branch) ImportProducts(ctx context.Context, branchId primitive.ObjectID, file *multipart.FileHeader) *helpers.ApiError {
	if strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "CSV files are not supported. Please upload an Excel (.xlsx) file or use the provided Excel template.",
		}
	}

	branch, apiErr := p.GetById(ctx, branchId, "")
	if apiErr != nil {
		return apiErr
	}

	// Create Excel service
	excelService := p.excelService

	// Define columns for import
	columns := product_utils.GetBranchProductImportColumns()

	// Import data from Excel
	importedData, apiErr := excelService.Import(ctx, file, columns)
	if apiErr != nil {
		return apiErr
	}

	// Check if import validation was successful
	if !importedData.IsValid {
		errors := make([]string, len(importedData.Errors))
		for i, err := range importedData.Errors {
			errors[i] = err.Error
		}
		return &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid data",
		}
	}

	// Get all products for lookup
	products, apiErr := p.productService.GetAll(ctx, []bson.M{})
	if apiErr != nil {
		return apiErr
	}

	// Convert Excel data to LinkBranchToProductDto objects
	branchProducts, apiErr := product_utils.ConvertExcelDataToBranchProducts(ctx, importedData.Data, products, branch.BaseCurrencyCode)
	if apiErr != nil {
		return apiErr
	}

	// Use utility function to convert branch data to update DTO
	updateDto := product_utils.ConvertBranchToUpdateDto(branch, branchProducts)

	_, apiErr = p.UpdateById(ctx, branchId, updateDto)
	if apiErr != nil {
		return apiErr
	}

	return nil
}

func (p *branch) ExportProducts(ctx context.Context, branchId primitive.ObjectID) ([]byte, *helpers.ApiError) {
	branch, apiErr := p.GetById(ctx, branchId, `{"popFields":["products"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	if len(branch.Products) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "No products found for this branch",
		}
	}

	// Get product IDs for lookup
	productIds := make([]primitive.ObjectID, 0)
	for _, bp := range branch.Products {
		if bp.ProductId != nil {
			productIds = append(productIds, *bp.ProductId)
		}
	}

	// Get all products for export using proper MongoDB aggregation pipeline
	products, apiErr := p.productService.GetAll(ctx, []bson.M{{"$match": bson.M{"_id": bson.M{"$in": productIds}}}})
	if apiErr != nil {
		return nil, apiErr
	}

	// Convert branch products to Excel data
	excelData, apiErr := product_utils.ConvertBranchProductsToExcelData(ctx, branch.Products, products)
	if apiErr != nil {
		return nil, apiErr
	}

	// Create Excel service
	excelService := p.excelService

	// Define columns for export
	columns := product_utils.GetBranchProductExportColumns()

	// Generate Excel file
	excelBytes, apiErr := excelService.Export(ctx, excelData, columns, "Branch Products Export")
	if apiErr != nil {
		return nil, apiErr
	}

	return excelBytes, nil
}

func (p *branch) ExportZones(ctx context.Context, branchId primitive.ObjectID) ([]byte, *helpers.ApiError) {
	branch, apiErr := p.GetById(ctx, branchId, `{"popFields":["zones"]}`)
	if apiErr != nil {
		return nil, apiErr
	}

	if branch.Zones == nil || len(*branch.Zones) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: http.StatusNotFound,
			Message:    "No zones found for this branch",
		}
	}

	// Convert branch zones to Excel data
	excelData, apiErr := product_utils.ConvertBranchZonesToExcelData(ctx, *branch.Zones, p.zoneService)
	if apiErr != nil {
		return nil, apiErr
	}

	// Create Excel service
	excelService := p.excelService

	// Define columns for export
	columns := product_utils.GetBranchZoneExportColumns()

	// Generate Excel file
	excelBytes, apiErr := excelService.Export(ctx, excelData, columns, "Branch Zones Export")
	if apiErr != nil {
		return nil, apiErr
	}

	return excelBytes, nil
}

func (p *branch) GetProductImportExample(ctx context.Context) ([]byte, *helpers.ApiError) {
	// Create Excel service
	excelService := p.excelService

	// Define columns for import
	columns := product_utils.GetBranchProductImportColumns()

	// Generate example Excel file
	excelData, apiErr := excelService.GetExample(ctx, columns, "Product Import Template")
	if apiErr != nil {
		return nil, apiErr
	}

	return excelData, nil
}

func (p *branch) LinkAllProductsToTheFirstBranch(ctx context.Context) *helpers.ApiError {
	return p.repo.LinkAllProductsToTheFirstBranch(ctx)
}
