package services

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ExcelColumnType represents the data type of an Excel column
type ExcelColumnType string

const (
	ExcelTypeString   ExcelColumnType = "string"
	ExcelTypeInt      ExcelColumnType = "int"
	ExcelTypeFloat    ExcelColumnType = "float"
	ExcelTypeBool     ExcelColumnType = "bool"
	ExcelTypeDate     ExcelColumnType = "date"
	ExcelTypeDateTime ExcelColumnType = "datetime"
	ExcelTypeEmail    ExcelColumnType = "email"
	ExcelTypeURL      ExcelColumnType = "url"
)

// ExcelColumn defines the structure and validation rules for an Excel column
type ExcelColumn struct {
	Key           string          `json:"key"`           // The key used in the resulting JSON data
	Header        string          `json:"header"`        // The header name displayed in Excel
	Type          ExcelColumnType `json:"type"`          // Data type of the column
	IsRequired    bool            `json:"isRequired"`    // Whether this column is required
	DefaultValue  interface{}     `json:"defaultValue"`  // Default value if not provided
	MaxLength     *int            `json:"maxLength"`     // Maximum length for string types
	MinLength     *int            `json:"minLength"`     // Minimum length for string types
	MinValue      *float64        `json:"minValue"`      // Minimum value for numeric types
	MaxValue      *float64        `json:"maxValue"`      // Maximum value for numeric types
	AllowedValues []string        `json:"allowedValues"` // List of allowed values (for enums)
	DateFormat    string          `json:"dateFormat"`    // Expected date format
	Regex         string          `json:"regex"`         // Regular expression for validation
	Description   string          `json:"description"`   // Description/help text for the column
}

// ExcelValidationError represents a validation error in Excel data
type ExcelValidationError struct {
	Row    int    `json:"row"`    // Row number where error occurred (1-indexed)
	Column string `json:"column"` // Column header name
	Key    string `json:"key"`    // Column key
	Value  string `json:"value"`  // The invalid value
	Error  string `json:"error"`  // Error description
}

// ExcelValidationResult contains the result of Excel validation
type ExcelValidationResult struct {
	IsValid bool                     `json:"isValid"` // Whether all data is valid
	Errors  []ExcelValidationError   `json:"errors"`  // List of validation errors
	Data    []map[string]interface{} `json:"data"`    // Validated and converted data
}

// Validate validates a cell value against the column definition
func (c *ExcelColumn) Validate(value interface{}, row int) *ExcelValidationError {
	strValue := fmt.Sprintf("%v", value)

	// Check if required field is empty
	if c.IsRequired && (value == nil || strValue == "" || strValue == "<nil>") {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  strValue,
			Error:  "Field is required",
		}
	}

	// If value is empty and not required, it's valid
	if value == nil || strValue == "" || strValue == "<nil>" {
		return nil
	}

	// Type-specific validation
	switch c.Type {
	case ExcelTypeString:
		return c.validateString(strValue, row)
	case ExcelTypeInt:
		return c.validateInt(strValue, row)
	case ExcelTypeFloat:
		return c.validateFloat(strValue, row)
	case ExcelTypeBool:
		return c.validateBool(strValue, row)
	case ExcelTypeDate:
		return c.validateDate(strValue, row)
	case ExcelTypeDateTime:
		return c.validateDateTime(strValue, row)
	case ExcelTypeEmail:
		return c.validateEmail(strValue, row)
	case ExcelTypeURL:
		return c.validateURL(strValue, row)
	}

	return nil
}

// validateString validates string values
func (c *ExcelColumn) validateString(value string, row int) *ExcelValidationError {
	// Check length constraints
	if c.MinLength != nil && len(value) < *c.MinLength {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Minimum length is %d characters", *c.MinLength),
		}
	}

	if c.MaxLength != nil && len(value) > *c.MaxLength {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Maximum length is %d characters", *c.MaxLength),
		}
	}

	// Check allowed values
	if len(c.AllowedValues) > 0 {
		found := false
		for _, allowed := range c.AllowedValues {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			return &ExcelValidationError{
				Row:    row,
				Column: c.Header,
				Key:    c.Key,
				Value:  value,
				Error:  fmt.Sprintf("Value must be one of: %v", c.AllowedValues),
			}
		}
	}

	// Check regex pattern
	if c.Regex != "" {
		matched, err := regexp.MatchString(c.Regex, value)
		if err != nil || !matched {
			return &ExcelValidationError{
				Row:    row,
				Column: c.Header,
				Key:    c.Key,
				Value:  value,
				Error:  "Value does not match required pattern",
			}
		}
	}

	return nil
}

// validateInt validates integer values
func (c *ExcelColumn) validateInt(value string, row int) *ExcelValidationError {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  "Value must be a valid integer",
		}
	}

	floatValue := float64(intValue)
	if c.MinValue != nil && floatValue < *c.MinValue {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Minimum value is %.0f", *c.MinValue),
		}
	}

	if c.MaxValue != nil && floatValue > *c.MaxValue {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Maximum value is %.0f", *c.MaxValue),
		}
	}

	return nil
}

// validateFloat validates float values
func (c *ExcelColumn) validateFloat(value string, row int) *ExcelValidationError {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  "Value must be a valid number",
		}
	}

	if c.MinValue != nil && floatValue < *c.MinValue {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Minimum value is %.2f", *c.MinValue),
		}
	}

	if c.MaxValue != nil && floatValue > *c.MaxValue {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Maximum value is %.2f", *c.MaxValue),
		}
	}

	return nil
}

// validateBool validates boolean values
func (c *ExcelColumn) validateBool(value string, row int) *ExcelValidationError {
	value = strings.ToLower(strings.TrimSpace(value))
	validValues := []string{"true", "false", "1", "0", "yes", "no", "y", "n"}

	for _, valid := range validValues {
		if value == valid {
			return nil
		}
	}

	return &ExcelValidationError{
		Row:    row,
		Column: c.Header,
		Key:    c.Key,
		Value:  value,
		Error:  "Value must be true/false, 1/0, yes/no, or y/n",
	}
}

// validateDate validates date values
func (c *ExcelColumn) validateDate(value string, row int) *ExcelValidationError {
	format := c.DateFormat
	if format == "" {
		format = "2006-01-02" // Default date format
	}

	_, err := time.Parse(format, value)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Invalid date format. Expected: %s", format),
		}
	}

	return nil
}

// validateDateTime validates datetime values
func (c *ExcelColumn) validateDateTime(value string, row int) *ExcelValidationError {
	format := c.DateFormat
	if format == "" {
		format = "2006-01-02 15:04:05" // Default datetime format
	}

	_, err := time.Parse(format, value)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  fmt.Sprintf("Invalid datetime format. Expected: %s", format),
		}
	}

	return nil
}

// validateEmail validates email addresses
func (c *ExcelColumn) validateEmail(value string, row int) *ExcelValidationError {
	_, err := mail.ParseAddress(value)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  "Invalid email address format",
		}
	}

	return nil
}

// validateURL validates URL values
func (c *ExcelColumn) validateURL(value string, row int) *ExcelValidationError {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  "Invalid URL format",
		}
	}

	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return &ExcelValidationError{
			Row:    row,
			Column: c.Header,
			Key:    c.Key,
			Value:  value,
			Error:  "URL must start with http:// or https://",
		}
	}

	return nil
}

// ConvertValue converts a value to the appropriate type based on column definition
func (c *ExcelColumn) ConvertValue(value interface{}) (interface{}, error) {
	if value == nil {
		return c.DefaultValue, nil
	}

	strValue := fmt.Sprintf("%v", value)
	if strValue == "" || strValue == "<nil>" {
		return c.DefaultValue, nil
	}

	switch c.Type {
	case ExcelTypeString, ExcelTypeEmail, ExcelTypeURL:
		return strValue, nil
	case ExcelTypeInt:
		return strconv.Atoi(strValue)
	case ExcelTypeFloat:
		return strconv.ParseFloat(strValue, 64)
	case ExcelTypeBool:
		return c.convertBool(strValue), nil
	case ExcelTypeDate:
		format := c.DateFormat
		if format == "" {
			format = "2006-01-02"
		}
		return time.Parse(format, strValue)
	case ExcelTypeDateTime:
		format := c.DateFormat
		if format == "" {
			format = "2006-01-02 15:04:05"
		}
		return time.Parse(format, strValue)
	default:
		return strValue, nil
	}
}

// convertBool converts string to boolean
func (c *ExcelColumn) convertBool(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1" || value == "yes" || value == "y"
}
