# Generic Excel Service

A powerful, generic Excel service for handling import, export, and example file generation across different entities in the application.

## Overview

The Excel service provides three main functions:
- **Import**: Read and validate Excel files against column definitions
- **Export**: Generate Excel files from data using column definitions  
- **GetExample**: Create example Excel files with headers and sample data

## Core Components

### 1. ExcelColumn Structure

Defines validation rules and metadata for each Excel column:

```go
type ExcelColumn struct {
    Key           string          // The key used in the resulting JSON data
    Header        string          // The header name displayed in Excel
    Type          ExcelColumnType // Data type of the column
    IsRequired    bool            // Whether this column is required
    DefaultValue  interface{}     // Default value if not provided
    MaxLength     *int            // Maximum length for string types
    MinLength     *int            // Minimum length for string types
    MinValue      *float64        // Minimum value for numeric types
    MaxValue      *float64        // Maximum value for numeric types
    AllowedValues []string        // List of allowed values (for enums)
    DateFormat    string          // Expected date format
    Regex         string          // Regular expression for validation
    Description   string          // Description/help text for the column
}
```

### 2. Column Types

Supported data types:
- `ExcelTypeString` - Text values
- `ExcelTypeInt` - Integer numbers
- `ExcelTypeFloat` - Decimal numbers
- `ExcelTypeBool` - Boolean values (true/false, 1/0, yes/no, y/n)
- `ExcelTypeDate` - Date values (default format: "2006-01-02")
- `ExcelTypeDateTime` - DateTime values (default format: "2006-01-02 15:04:05")
- `ExcelTypeEmail` - Email addresses with basic validation
- `ExcelTypeURL` - URLs starting with http:// or https://

### 3. Validation Results

```go
type ExcelValidationResult struct {
    IsValid bool                     // Whether all data is valid
    Errors  []ExcelValidationError   // List of validation errors
    Data    []map[string]interface{} // Validated and converted data
}

type ExcelValidationError struct {
    Row    int    // Row number where error occurred
    Column string // Column header name
    Key    string // Column key
    Value  string // The invalid value
    Error  string // Error description
}
```

## Usage

### Basic Service Usage

```go
// Create service instance
excelService := services.NewExcelService()

// Define columns
columns := []services.ExcelColumn{
    services.NewRequiredStringColumn("name", "Product Name", services.IntPtr(100)),
    services.NewFloatColumn("price", "Price", true, services.Float64Ptr(0), nil),
    services.NewEnumColumn("status", "Status", true, []string{"Active", "Inactive"}, "Active"),
}

// Import Excel file
result, err := excelService.Import(ctx, file, columns)
if err != nil {
    // Handle error
}

if !result.IsValid {
    // Handle validation errors
    for _, validationError := range result.Errors {
        fmt.Printf("Row %d, Column %s: %s\n", 
            validationError.Row, 
            validationError.Column, 
            validationError.Error)
    }
}

// Use validated data
for _, row := range result.Data {
    name := row["name"].(string)
    price := row["price"].(float64)
    // Process data...
}
```

### Column Builder Functions

Use the provided builder functions for common column types:

```go
// Required string with max length
services.NewRequiredStringColumn("name", "Product Name", services.IntPtr(100))

// Optional string with default value
services.NewOptionalStringColumn("description", "Description", "No description", services.IntPtr(500))

// Enum with allowed values
services.NewEnumColumn("status", "Status", true, []string{"Active", "Inactive", "Draft"}, "Active")

// Email validation
services.NewEmailColumn("email", "Email Address", true)

// Integer with min/max constraints
services.NewIntColumn("quantity", "Quantity", true, services.Float64Ptr(0), services.Float64Ptr(1000))

// Float with constraints
services.NewFloatColumn("price", "Price", true, services.Float64Ptr(0), nil)

// Boolean with default
services.NewBoolColumn("isActive", "Active", false, services.BoolPtr(true))

// Date with custom format
services.NewDateColumn("createdAt", "Created Date", false, "2006-01-02")

// URL validation
services.NewURLColumn("website", "Website URL", false)
```

### Helper Functions

```go
// Pointer helpers for constraints
services.IntPtr(100)        // *int
services.Float64Ptr(0.5)    // *float64
services.BoolPtr(true)      // *bool
```

### Predefined Column Sets

For common entities:

```go
// Basic product columns
columns := services.GetBasicProductColumns()

// Basic user columns  
columns := services.GetBasicUserColumns()

// Basic category columns
columns := services.GetBasicCategoryColumns()
```

## Module Integration Example

### 1. Create Module-Specific Excel Service

```go
// internal/modules/product/utils/excel-integration.go
package product_utils

import (
    "context"
    "mime/multipart"
    
    "larsa-e-commerce-microservice-v2/internal/common/services"
    product_models "larsa-e-commerce-microservice-v2/internal/modules/product/models"
)

type ProductExcelService struct {
    excelService services.ExcelService
}

func NewProductExcelService() *ProductExcelService {
    return &ProductExcelService{
        excelService: services.NewExcelService(),
    }
}

func (s *ProductExcelService) GetProductColumns() []services.ExcelColumn {
    return []services.ExcelColumn{
        services.NewRequiredStringColumn("nameEn", "Product Name (EN)", services.IntPtr(100)),
        services.NewRequiredStringColumn("nameAr", "Product Name (AR)", services.IntPtr(100)),
        services.NewOptionalStringColumn("descriptionEn", "Description (EN)", "", services.IntPtr(500)),
        services.NewOptionalStringColumn("descriptionAr", "Description (AR)", "", services.IntPtr(500)),
        services.NewEnumColumn("flag", "Flag", false, []string{"hot_pick", "trending", "innovative"}, ""),
    }
}

func (s *ProductExcelService) ImportProducts(ctx context.Context, file *multipart.FileHeader) (*services.ExcelValidationResult, *helpers.ApiError) {
    columns := s.GetProductColumns()
    return s.excelService.Import(ctx, file, columns)
}

func (s *ProductExcelService) ExportProducts(ctx context.Context, products []product_models.Product) ([]byte, *helpers.ApiError) {
    columns := s.GetProductColumns()
    
    // Convert models to map[string]interface{}
    data := make([]map[string]interface{}, len(products))
    for i, product := range products {
        data[i] = s.convertProductToMap(product)
    }
    
    return s.excelService.Export(ctx, data, columns, "Products")
}

func (s *ProductExcelService) GetProductsExampleFile(ctx context.Context) ([]byte, *helpers.ApiError) {
    columns := s.GetProductColumns()
    return s.excelService.GetExample(ctx, columns, "Products_Example")
}
```

### 2. Add to Handler

```go
// internal/modules/product/handlers/product.handler.go

func (h *ProductHandler) ImportProducts(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("file")
    if err != nil {
        // Handle error
        return
    }
    
    excelService := product_utils.NewProductExcelService()
    result, apiErr := excelService.ImportProducts(r.Context(), file)
    if apiErr != nil {
        // Handle error
        return
    }
    
    if !result.IsValid {
        // Return validation errors
        helpers.WriteResponse(w, http.StatusBadRequest, result.Errors)
        return
    }
    
    // Convert to DTOs and process
    products := excelService.ConvertImportDataToProducts(result.Data)
    // Save products...
    
    helpers.WriteResponse(w, http.StatusOK, map[string]interface{}{
        "message": "Products imported successfully",
        "count":   len(products),
    })
}

func (h *ProductHandler) ExportProducts(w http.ResponseWriter, r *http.Request) {
    // Get products from database
    products, err := h.service.GetAll(r.Context())
    if err != nil {
        // Handle error
        return
    }
    
    excelService := product_utils.NewProductExcelService()
    excelData, apiErr := excelService.ExportProducts(r.Context(), products)
    if apiErr != nil {
        // Handle error
        return
    }
    
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=products.xlsx")
    w.Write(excelData)
}

func (h *ProductHandler) GetProductsExampleFile(w http.ResponseWriter, r *http.Request) {
    excelService := product_utils.NewProductExcelService()
    excelData, apiErr := excelService.GetProductsExampleFile(r.Context())
    if apiErr != nil {
        // Handle error
        return
    }
    
    w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    w.Header().Set("Content-Disposition", "attachment; filename=products_example.xlsx")
    w.Write(excelData)
}
```

### 3. Add Routes

```go
// Add to your router
r.Post("/products/import", h.ImportProducts)
r.Get("/products/export", h.ExportProducts)
r.Get("/products/example", h.GetProductsExampleFile)
```

## Advanced Features

### Custom Validation

You can create custom column types with specific validation:

```go
column := services.ExcelColumn{
    Key:        "customField",
    Header:     "Custom Field",
    Type:       services.ExcelTypeString,
    IsRequired: true,
    Regex:      `^[A-Z]{2,3}-\d{4}$`, // Custom pattern
    Description: "Format: XX-1234 or XXX-1234",
}
```

### Complex Data Handling

For complex fields like arrays or nested objects:

```go
// In your conversion function
func (s *ProductExcelService) ConvertImportDataToProducts(data []map[string]interface{}) []product_dtos.CreateProductDto {
    products := make([]product_dtos.CreateProductDto, len(data))
    
    for i, row := range data {
        product := product_dtos.CreateProductDto{}
        
        // Handle comma-separated values
        if categoryIds, ok := row["categoryIds"].(string); ok && categoryIds != "" {
            product.CategoryIds = s.parseObjectIds(categoryIds)
        }
        
        // Handle localized fields
        if nameEn, ok := row["nameEn"].(string); ok {
            if nameAr, ok := row["nameAr"].(string); ok {
                product.Name = types.Localizable[string]{
                    enums.EN: nameEn,
                    enums.AR: nameAr,
                }
            }
        }
        
        products[i] = product
    }
    
    return products
}
```

## Error Handling

The service provides detailed error information:

```go
result, err := excelService.Import(ctx, file, columns)
if err != nil {
    // System error (file read, parsing, etc.)
    return err
}

if !result.IsValid {
    // Validation errors
    for _, validationError := range result.Errors {
        log.Printf("Validation error at row %d, column %s: %s", 
            validationError.Row, 
            validationError.Column, 
            validationError.Error)
    }
}
```

## Best Practices

1. **Define Clear Column Headers**: Use descriptive headers that users will understand
2. **Provide Examples**: Always offer example files to guide users
3. **Validate Early**: Use the validation results to provide immediate feedback
4. **Handle Localization**: For multi-language apps, consider separate columns for each language
5. **Document Constraints**: Use the Description field to explain validation rules
6. **Test Edge Cases**: Test with empty files, missing columns, invalid data types
7. **Provide Feedback**: Return detailed error messages for validation failures

## Dependencies

- `github.com/360EntSecGroup-Skylar/excelize` - Excel file processing (qax-os/excelize fork)
- Internal helpers and types from the common package

**Note**: This implementation uses the qax-os/excelize library (v1.x) instead of the newer xuri/excelize/v2. Some advanced styling features are not available in the v1.x API, but all core functionality (import, export, validation) works as expected.

## File Structure

```
internal/common/services/
├── excel-column.go      # Column definitions and validation
├── excel.service.go     # Main service implementation  
├── excel-utils.go       # Helper functions and builders
└── README.md           # This documentation
``` 