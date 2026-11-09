# Excel Service Examples

This document provides practical examples of using the generic Excel service with different modules and entities.

## Example 1: Simple User Import/Export

### Define User Columns

```go
func GetUserColumns() []services.ExcelColumn {
    return []services.ExcelColumn{
        services.NewRequiredStringColumn("firstName", "First Name", services.IntPtr(50)),
        services.NewRequiredStringColumn("lastName", "Last Name", services.IntPtr(50)),
        services.NewEmailColumn("email", "Email", true),
        services.NewOptionalStringColumn("phone", "Phone", "", services.IntPtr(20)),
        services.NewEnumColumn("role", "Role", true, []string{"admin", "user", "manager"}, "user"),
        services.NewBoolColumn("isActive", "Active", false, services.BoolPtr(true)),
        services.NewDateColumn("birthDate", "Birth Date", false, "2006-01-02"),
    }
}
```

### Import Users

```go
func ImportUsers(ctx context.Context, file *multipart.FileHeader) (*services.ExcelValidationResult, error) {
    excelService := services.NewExcelService()
    columns := GetUserColumns()
    
    result, err := excelService.Import(ctx, file, columns)
    if err != nil {
        return nil, err
    }
    
    // Process validated data
    if result.IsValid {
        for _, row := range result.Data {
            user := User{
                FirstName: row["firstName"].(string),
                LastName:  row["lastName"].(string),
                Email:     row["email"].(string),
                Role:      row["role"].(string),
                IsActive:  row["isActive"].(bool),
            }
            
            if phone, ok := row["phone"].(string); ok && phone != "" {
                user.Phone = &phone
            }
            
            if birthDate, ok := row["birthDate"].(time.Time); ok {
                user.BirthDate = &birthDate
            }
            
            // Save user to database
            // userService.Create(ctx, user)
        }
    }
    
    return result, nil
}
```

### Export Users

```go
func ExportUsers(ctx context.Context, users []User) ([]byte, error) {
    excelService := services.NewExcelService()
    columns := GetUserColumns()
    
    // Convert users to map format
    data := make([]map[string]interface{}, len(users))
    for i, user := range users {
        data[i] = map[string]interface{}{
            "firstName": user.FirstName,
            "lastName":  user.LastName,
            "email":     user.Email,
            "role":      user.Role,
            "isActive":  user.IsActive,
        }
        
        if user.Phone != nil {
            data[i]["phone"] = *user.Phone
        }
        
        if user.BirthDate != nil {
            data[i]["birthDate"] = user.BirthDate.Format("2006-01-02")
        }
    }
    
    return excelService.Export(ctx, data, columns, "Users")
}
```

## Example 2: Category Management with Localization

### Define Category Columns

```go
func GetCategoryColumns() []services.ExcelColumn {
    return []services.ExcelColumn{
        services.NewRequiredStringColumn("nameEn", "Name (English)", services.IntPtr(100)),
        services.NewRequiredStringColumn("nameAr", "Name (Arabic)", services.IntPtr(100)),
        services.NewOptionalStringColumn("descriptionEn", "Description (English)", "", services.IntPtr(500)),
        services.NewOptionalStringColumn("descriptionAr", "Description (Arabic)", "", services.IntPtr(500)),
        services.NewOptionalStringColumn("parentId", "Parent Category ID", "", services.IntPtr(24)),
        services.NewIntColumn("sortOrder", "Sort Order", false, services.Float64Ptr(0), services.Float64Ptr(999)),
        services.NewBoolColumn("isActive", "Active", false, services.BoolPtr(true)),
        services.NewEnumColumn("visibility", "Visibility", true, []string{"public", "private", "hidden"}, "public"),
    }
}
```

### Category Service

```go
type CategoryExcelService struct {
    excelService services.ExcelService
}

func NewCategoryExcelService() *CategoryExcelService {
    return &CategoryExcelService{
        excelService: services.NewExcelService(),
    }
}

func (s *CategoryExcelService) ImportCategories(ctx context.Context, file *multipart.FileHeader) (*services.ExcelValidationResult, error) {
    columns := GetCategoryColumns()
    return s.excelService.Import(ctx, file, columns)
}

func (s *CategoryExcelService) ConvertToCategories(data []map[string]interface{}) []CreateCategoryDTO {
    categories := make([]CreateCategoryDTO, len(data))
    
    for i, row := range data {
        category := CreateCategoryDTO{}
        
        // Handle localized names
        if nameEn, ok := row["nameEn"].(string); ok {
            if nameAr, ok := row["nameAr"].(string); ok {
                category.Name = types.Localizable[string]{
                    enums.EN: nameEn,
                    enums.AR: nameAr,
                }
            }
        }
        
        // Handle localized descriptions
        if descEn, ok := row["descriptionEn"].(string); ok {
            if descAr, ok := row["descriptionAr"].(string); ok {
                category.Description = types.Localizable[string]{
                    enums.EN: descEn,
                    enums.AR: descAr,
                }
            }
        }
        
        // Handle parent ID
        if parentId, ok := row["parentId"].(string); ok && parentId != "" {
            if objId, err := primitive.ObjectIDFromHex(parentId); err == nil {
                category.ParentId = &objId
            }
        }
        
        // Handle other fields
        if sortOrder, ok := row["sortOrder"].(int); ok {
            category.SortOrder = &sortOrder
        }
        
        if isActive, ok := row["isActive"].(bool); ok {
            category.IsActive = &isActive
        }
        
        if visibility, ok := row["visibility"].(string); ok {
            category.Visibility = CategoryVisibility(visibility)
        }
        
        categories[i] = category
    }
    
    return categories
}
```

## Example 3: Inventory Management

### Define Inventory Columns

```go
func GetInventoryColumns() []services.ExcelColumn {
    return []services.ExcelColumn{
        services.NewRequiredStringColumn("sku", "SKU", services.IntPtr(50)),
        services.NewRequiredStringColumn("productName", "Product Name", services.IntPtr(100)),
        services.NewIntColumn("quantity", "Quantity", true, services.Float64Ptr(0), nil),
        services.NewFloatColumn("unitPrice", "Unit Price", true, services.Float64Ptr(0), nil),
        services.NewOptionalStringColumn("location", "Warehouse Location", "", services.IntPtr(50)),
        services.NewDateColumn("lastUpdated", "Last Updated", false, "2006-01-02"),
        services.NewEnumColumn("status", "Status", true, []string{"in_stock", "low_stock", "out_of_stock", "discontinued"}, "in_stock"),
        services.NewOptionalStringColumn("supplier", "Supplier", "", services.IntPtr(100)),
        services.NewFloatColumn("reorderLevel", "Reorder Level", false, services.Float64Ptr(0), nil),
    }
}
```

### Inventory Import with Validation

```go
func ImportInventory(ctx context.Context, file *multipart.FileHeader) (*services.ExcelValidationResult, error) {
    excelService := services.NewExcelService()
    columns := GetInventoryColumns()
    
    result, err := excelService.Import(ctx, file, columns)
    if err != nil {
        return nil, err
    }
    
    // Additional business logic validation
    if result.IsValid {
        for i, row := range result.Data {
            // Check if SKU already exists
            sku := row["sku"].(string)
            if existingProduct := findProductBySKU(sku); existingProduct != nil {
                result.Errors = append(result.Errors, services.ExcelValidationError{
                    Row:    i + 2, // +2 for header and 1-indexed
                    Column: "SKU",
                    Key:    "sku",
                    Value:  sku,
                    Error:  "SKU already exists in the system",
                })
                result.IsValid = false
            }
            
            // Validate reorder level vs quantity
            if quantity, ok := row["quantity"].(int); ok {
                if reorderLevel, ok := row["reorderLevel"].(float64); ok {
                    if float64(quantity) < reorderLevel {
                        // This is a warning, not an error
                        log.Printf("Warning: Product %s quantity (%d) is below reorder level (%.0f)", 
                            sku, quantity, reorderLevel)
                    }
                }
            }
        }
    }
    
    return result, nil
}
```

## Example 4: Financial Data Import

### Define Financial Columns

```go
func GetTransactionColumns() []services.ExcelColumn {
    return []services.ExcelColumn{
        services.NewRequiredStringColumn("transactionId", "Transaction ID", services.IntPtr(50)),
        services.NewDateTimeColumn("timestamp", "Timestamp", true, "2006-01-02 15:04:05"),
        services.NewFloatColumn("amount", "Amount", true, nil, nil), // No min/max for amounts
        services.NewEnumColumn("currency", "Currency", true, []string{"USD", "EUR", "SAR", "AED"}, "USD"),
        services.NewEnumColumn("type", "Type", true, []string{"credit", "debit", "transfer"}, ""),
        services.NewRequiredStringColumn("description", "Description", services.IntPtr(200)),
        services.NewOptionalStringColumn("reference", "Reference", "", services.IntPtr(100)),
        services.NewEnumColumn("status", "Status", true, []string{"pending", "completed", "failed", "cancelled"}, "pending"),
        services.NewOptionalStringColumn("notes", "Notes", "", services.IntPtr(500)),
    }
}
```

### Custom Validation for Financial Data

```go
func ImportTransactions(ctx context.Context, file *multipart.FileHeader) (*services.ExcelValidationResult, error) {
    excelService := services.NewExcelService()
    columns := GetTransactionColumns()
    
    result, err := excelService.Import(ctx, file, columns)
    if err != nil {
        return nil, err
    }
    
    // Financial-specific validation
    if result.IsValid {
        for i, row := range result.Data {
            // Validate amount precision (max 2 decimal places)
            if amount, ok := row["amount"].(float64); ok {
                if math.Mod(amount*100, 1) != 0 {
                    result.Errors = append(result.Errors, services.ExcelValidationError{
                        Row:    i + 2,
                        Column: "Amount",
                        Key:    "amount",
                        Value:  fmt.Sprintf("%.4f", amount),
                        Error:  "Amount cannot have more than 2 decimal places",
                    })
                    result.IsValid = false
                }
            }
            
            // Validate transaction ID format
            if transactionId, ok := row["transactionId"].(string); ok {
                if !regexp.MustCompile(`^TXN-\d{8}-[A-Z0-9]{6}$`).MatchString(transactionId) {
                    result.Errors = append(result.Errors, services.ExcelValidationError{
                        Row:    i + 2,
                        Column: "Transaction ID",
                        Key:    "transactionId",
                        Value:  transactionId,
                        Error:  "Transaction ID must follow format: TXN-YYYYMMDD-XXXXXX",
                    })
                    result.IsValid = false
                }
            }
        }
    }
    
    return result, nil
}
```

## Example 5: Bulk Operations with Progress Tracking

### Import with Progress Callback

```go
type ImportProgress struct {
    Total     int `json:"total"`
    Processed int `json:"processed"`
    Errors    int `json:"errors"`
    Status    string `json:"status"`
}

func ImportWithProgress(ctx context.Context, file *multipart.FileHeader, progressChan chan<- ImportProgress) error {
    excelService := services.NewExcelService()
    columns := GetUserColumns()
    
    // Send initial progress
    progressChan <- ImportProgress{Status: "validating"}
    
    result, err := excelService.Import(ctx, file, columns)
    if err != nil {
        progressChan <- ImportProgress{Status: "error"}
        return err
    }
    
    if !result.IsValid {
        progressChan <- ImportProgress{
            Status: "validation_failed",
            Errors: len(result.Errors),
        }
        return fmt.Errorf("validation failed with %d errors", len(result.Errors))
    }
    
    // Process data with progress updates
    total := len(result.Data)
    progressChan <- ImportProgress{
        Total:  total,
        Status: "processing",
    }
    
    for i, row := range result.Data {
        // Process each row
        err := processUserRow(ctx, row)
        if err != nil {
            progressChan <- ImportProgress{
                Total:     total,
                Processed: i,
                Errors:    1,
                Status:    "processing",
            }
            continue
        }
        
        // Send progress update every 10 records
        if i%10 == 0 || i == total-1 {
            progressChan <- ImportProgress{
                Total:     total,
                Processed: i + 1,
                Status:    "processing",
            }
        }
    }
    
    progressChan <- ImportProgress{
        Total:     total,
        Processed: total,
        Status:    "completed",
    }
    
    return nil
}
```

## Example 6: Template Generation

### Generate Template with Instructions

```go
func GenerateUserImportTemplate(ctx context.Context) ([]byte, error) {
    excelService := services.NewExcelService()
    
    // Define columns with detailed descriptions
    columns := []services.ExcelColumn{
        {
            Key:         "firstName",
            Header:      "First Name",
            Type:        services.ExcelTypeString,
            IsRequired:  true,
            MaxLength:   services.IntPtr(50),
            Description: "User's first name (required, max 50 characters)",
        },
        {
            Key:         "lastName",
            Header:      "Last Name", 
            Type:        services.ExcelTypeString,
            IsRequired:  true,
            MaxLength:   services.IntPtr(50),
            Description: "User's last name (required, max 50 characters)",
        },
        {
            Key:         "email",
            Header:      "Email Address",
            Type:        services.ExcelTypeEmail,
            IsRequired:  true,
            Description: "Valid email address (required, must contain @ and .)",
        },
        {
            Key:           "role",
            Header:        "Role",
            Type:          services.ExcelTypeString,
            IsRequired:    true,
            AllowedValues: []string{"admin", "user", "manager"},
            DefaultValue:  "user",
            Description:   "User role (required, must be: admin, user, or manager)",
        },
        {
            Key:         "birthDate",
            Header:      "Birth Date",
            Type:        services.ExcelTypeDate,
            IsRequired:  false,
            DateFormat:  "2006-01-02",
            Description: "Birth date in YYYY-MM-DD format (optional)",
        },
    }
    
    return excelService.GetExample(ctx, columns, "User_Import_Template")
}
```

## Best Practices from Examples

1. **Validation Layers**: Use both Excel service validation and business logic validation
2. **Error Handling**: Provide detailed error messages with row and column information
3. **Progress Tracking**: For large imports, implement progress tracking
4. **Template Generation**: Always provide example files with clear instructions
5. **Data Conversion**: Handle type conversions carefully, especially for dates and numbers
6. **Localization**: Use separate columns for multi-language content
7. **Business Rules**: Implement domain-specific validation after basic Excel validation
8. **Rollback Strategy**: Consider transaction management for bulk operations

## Testing Examples

### Unit Test for Column Validation

```go
func TestUserColumnValidation(t *testing.T) {
    columns := GetUserColumns()
    
    // Test valid data
    validData := map[string]interface{}{
        "firstName": "John",
        "lastName":  "Doe",
        "email":     "john.doe@example.com",
        "role":      "user",
        "isActive":  true,
    }
    
    for _, column := range columns {
        if value, exists := validData[column.Key]; exists {
            err := column.Validate(value, 1)
            assert.Nil(t, err, "Valid data should not produce validation errors")
        }
    }
    
    // Test invalid email
    emailColumn := findColumnByKey(columns, "email")
    err := emailColumn.Validate("invalid-email", 1)
    assert.NotNil(t, err, "Invalid email should produce validation error")
    assert.Contains(t, err.Error, "valid email address")
}
```

### Integration Test

```go
func TestUserImportIntegration(t *testing.T) {
    // Create test Excel file
    excelService := services.NewExcelService()
    columns := GetUserColumns()
    
    testData := []map[string]interface{}{
        {
            "firstName": "John",
            "lastName":  "Doe", 
            "email":     "john@example.com",
            "role":      "user",
            "isActive":  true,
        },
    }
    
    excelBytes, err := excelService.Export(context.Background(), testData, columns, "TestUsers")
    assert.NoError(t, err)
    
    // Create multipart file from bytes
    file := createMultipartFileFromBytes(excelBytes, "test.xlsx")
    
    // Test import
    result, err := ImportUsers(context.Background(), file)
    assert.NoError(t, err)
    assert.True(t, result.IsValid)
    assert.Len(t, result.Data, 1)
}
``` 