package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"

	"github.com/samber/do"
	"github.com/xuri/excelize/v2"
)

// ExcelService interface defines the contract for Excel operations
type ExcelService interface {
	Import(ctx context.Context, file *multipart.FileHeader, columns []ExcelColumn) (*ExcelValidationResult, *helpers.ApiError)
	Export(ctx context.Context, data []map[string]interface{}, columns []ExcelColumn, sheetName string) ([]byte, *helpers.ApiError)
	GetExample(ctx context.Context, columns []ExcelColumn, sheetName string) ([]byte, *helpers.ApiError)
}

// excelService implements the ExcelService interface
type excelService struct{}

// NewExcelService creates a new instance of ExcelService
func NewExcelService(i *do.Injector) (ExcelService, error) {
	return &excelService{}, nil
}

// Import reads an Excel file and validates it against the provided column definitions
func (s *excelService) Import(ctx context.Context, file *multipart.FileHeader, columns []ExcelColumn) (*ExcelValidationResult, *helpers.ApiError) {
	if file.Size == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Uploaded file is empty",
		}
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Failed to open uploaded file",
		}
	}
	defer src.Close()

	// Read file content into memory
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Failed to read file content",
		}
	}

	// Open Excel file
	f, err := excelize.OpenReader(bytes.NewReader(fileContent))
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Invalid Excel file format",
		}
	}

	// Get the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "No sheets found in Excel file",
		}
	}

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Failed to get rows from Excel file",
		}
	}
	if len(rows) == 0 {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    "Excel file is empty",
		}
	}

	// Validate headers
	headerRow := rows[0]
	columnMap, apiErr := s.validateHeaders(headerRow, columns)
	if apiErr != nil {
		return nil, apiErr
	}

	// Process data rows
	result := &ExcelValidationResult{
		IsValid: true,
		Errors:  []ExcelValidationError{},
		Data:    []map[string]interface{}{},
	}

	// Process each data row (skip header row)
	for rowIndex, row := range rows[1:] {
		actualRowIndex := rowIndex + 2 // +2 because we skip header and Excel is 1-indexed

		rowData := make(map[string]interface{})
		hasData := false

		// Process each column
		for colIndex, column := range columns {
			var cellValue interface{}

			// Get cell value if column exists in this row
			if colIndex < len(row) && columnMap[colIndex] {
				cellValue = row[colIndex]
				if fmt.Sprintf("%v", cellValue) != "" {
					hasData = true
				}
			}

			// Validate the cell value
			if validationErr := column.Validate(cellValue, actualRowIndex); validationErr != nil {
				result.Errors = append(result.Errors, *validationErr)
				result.IsValid = false
				continue
			}

			// Convert the value to appropriate type
			convertedValue, err := column.ConvertValue(cellValue)
			if err != nil {
				result.Errors = append(result.Errors, ExcelValidationError{
					Row:    actualRowIndex,
					Column: column.Header,
					Key:    column.Key,
					Value:  fmt.Sprintf("%v", cellValue),
					Error:  fmt.Sprintf("Failed to convert value: %v", err),
				})
				result.IsValid = false
				continue
			}

			rowData[column.Key] = convertedValue
		}

		// Only add row if it has at least some data
		if hasData {
			result.Data = append(result.Data, rowData)
		}
	}

	return result, nil
}

// Export creates an Excel file from the provided data using column definitions
func (s *excelService) Export(ctx context.Context, data []map[string]interface{}, columns []ExcelColumn, sheetName string) ([]byte, *helpers.ApiError) {
	f := excelize.NewFile()

	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// Create sheet if it doesn't exist (v1.x API)
	if sheetName != "Sheet1" {
		f.NewSheet(sheetName)
		f.DeleteSheet("Sheet1")
	}

	// Set headers
	for colIndex, column := range columns {
		cell := fmt.Sprintf("%s1", s.columnIndexToLetter(colIndex))
		f.SetCellValue(sheetName, cell, column.Header)
	}

	// Set data
	for rowIndex, rowData := range data {
		excelRowIndex := rowIndex + 2 // +2 because Excel is 1-indexed and we have header row

		for colIndex, column := range columns {
			cell := fmt.Sprintf("%s%d", s.columnIndexToLetter(colIndex), excelRowIndex)

			value := rowData[column.Key]
			if value != nil {
				f.SetCellValue(sheetName, cell, value)
			}
		}
	}

	// Generate Excel file as bytes (v1.x API)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to generate Excel file",
		}
	}

	return buffer.Bytes(), nil
}

// GetExample creates an example Excel file with headers and sample data
func (s *excelService) GetExample(ctx context.Context, columns []ExcelColumn, sheetName string) ([]byte, *helpers.ApiError) {
	f := excelize.NewFile()

	if sheetName == "" {
		sheetName = "Example"
	}

	// Create sheet if it doesn't exist (v1.x API)
	if sheetName != "Sheet1" {
		f.NewSheet(sheetName)
		f.DeleteSheet("Sheet1")
	}

	// Set headers
	for colIndex, column := range columns {
		cell := fmt.Sprintf("%s1", s.columnIndexToLetter(colIndex))
		f.SetCellValue(sheetName, cell, column.Header)
	}

	// Add example data and validation info
	for colIndex, column := range columns {
		exampleCell := fmt.Sprintf("%s2", s.columnIndexToLetter(colIndex))
		infoCell := fmt.Sprintf("%s3", s.columnIndexToLetter(colIndex))

		// Add example value
		exampleValue := s.getExampleValue(column)
		if exampleValue != nil {
			f.SetCellValue(sheetName, exampleCell, exampleValue)
		}

		// Add validation info as comment
		validationInfo := s.buildValidationInfo(column)
		if validationInfo != "" {
			f.SetCellValue(sheetName, infoCell, validationInfo)
		}
	}

	// Set column widths (v1.x API)
	for colIndex, column := range columns {
		columnLetter := s.columnIndexToLetter(colIndex)
		width := float64(len(column.Header) + 5)
		if width < 15 {
			width = 15
		}
		if width > 50 {
			width = 50
		}
		f.SetColWidth(sheetName, columnLetter, columnLetter, width)
	}

	// Generate Excel file as bytes (v1.x API)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, &helpers.ApiError{
			StatusCode: 500,
			Message:    "Failed to generate example Excel file",
		}
	}

	return buffer.Bytes(), nil
}

// validateHeaders validates that the Excel headers match the expected columns
func (s *excelService) validateHeaders(headerRow []string, columns []ExcelColumn) (map[int]bool, *helpers.ApiError) {
	columnMap := make(map[int]bool)

	// Check if we have the expected number of columns
	if len(headerRow) < len(columns) {
		return nil, &helpers.ApiError{
			StatusCode: 400,
			Message:    fmt.Sprintf("Excel file has %d columns but expected %d", len(headerRow), len(columns)),
		}
	}

	// Validate each column header
	for colIndex, expectedColumn := range columns {
		if colIndex >= len(headerRow) {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    fmt.Sprintf("Missing column: %s", expectedColumn.Header),
			}
		}

		if headerRow[colIndex] != expectedColumn.Header {
			return nil, &helpers.ApiError{
				StatusCode: 400,
				Message:    fmt.Sprintf("Column %d: expected '%s' but got '%s'", colIndex+1, expectedColumn.Header, headerRow[colIndex]),
			}
		}

		columnMap[colIndex] = true
	}

	return columnMap, nil
}

// columnIndexToLetter converts a column index to Excel column letter (A, B, C, ...)
func (s *excelService) columnIndexToLetter(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
}

// getExampleValue generates an example value for a column based on its type
func (s *excelService) getExampleValue(column ExcelColumn) interface{} {
	if column.DefaultValue != nil {
		return column.DefaultValue
	}

	switch column.Type {
	case ExcelTypeString:
		if len(column.AllowedValues) > 0 {
			return column.AllowedValues[0]
		}
		return "Example " + column.Key
	case ExcelTypeEmail:
		return "example@email.com"
	case ExcelTypeURL:
		return "https://example.com"
	case ExcelTypeInt:
		if column.MinValue != nil {
			return int(*column.MinValue)
		}
		return 1
	case ExcelTypeFloat:
		if column.MinValue != nil {
			return *column.MinValue
		}
		return 1.0
	case ExcelTypeBool:
		return true
	case ExcelTypeDate:
		return "2024-01-01"
	case ExcelTypeDateTime:
		return "2024-01-01 10:00:00"
	default:
		return "Example"
	}
}

// buildValidationInfo creates validation information text for a column
func (s *excelService) buildValidationInfo(column ExcelColumn) string {
	info := []string{}

	if column.IsRequired {
		info = append(info, "Required")
	}

	switch column.Type {
	case ExcelTypeString:
		if column.MinLength != nil || column.MaxLength != nil {
			lengthInfo := "Length: "
			if column.MinLength != nil && column.MaxLength != nil {
				lengthInfo += fmt.Sprintf("%d-%d chars", *column.MinLength, *column.MaxLength)
			} else if column.MinLength != nil {
				lengthInfo += fmt.Sprintf("min %d chars", *column.MinLength)
			} else if column.MaxLength != nil {
				lengthInfo += fmt.Sprintf("max %d chars", *column.MaxLength)
			}
			info = append(info, lengthInfo)
		}
	case ExcelTypeInt, ExcelTypeFloat:
		if column.MinValue != nil || column.MaxValue != nil {
			rangeInfo := "Range: "
			if column.MinValue != nil && column.MaxValue != nil {
				rangeInfo += fmt.Sprintf("%.0f-%.0f", *column.MinValue, *column.MaxValue)
			} else if column.MinValue != nil {
				rangeInfo += fmt.Sprintf("min %.0f", *column.MinValue)
			} else if column.MaxValue != nil {
				rangeInfo += fmt.Sprintf("max %.0f", *column.MaxValue)
			}
			info = append(info, rangeInfo)
		}
	}

	if len(column.AllowedValues) > 0 {
		allowedInfo := "Allowed: " + fmt.Sprintf("%v", column.AllowedValues)
		if len(allowedInfo) > 50 {
			allowedInfo = allowedInfo[:47] + "..."
		}
		info = append(info, allowedInfo)
	}

	if column.DateFormat != "" {
		info = append(info, "Format: "+column.DateFormat)
	}

	if column.Description != "" {
		info = append(info, column.Description)
	}

	if len(info) == 0 {
		return ""
	}

	return fmt.Sprintf("(%s)", fmt.Sprintf("%v", info))
}
