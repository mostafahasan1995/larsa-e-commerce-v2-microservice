package services

// Common column builders for frequent use cases

// Helper functions for creating pointer types
func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func BoolPtr(b bool) *bool {
	return &b
}

// Column builder functions for common types

// NewRequiredStringColumn creates a required string column with optional max length
func NewRequiredStringColumn(key, header string, maxLength *int) ExcelColumn {
	return ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeString,
		IsRequired: true,
		MaxLength:  maxLength,
	}
}

// NewOptionalStringColumn creates an optional string column with default value and optional max length
func NewOptionalStringColumn(key, header, defaultValue string, maxLength *int) ExcelColumn {
	return ExcelColumn{
		Key:          key,
		Header:       header,
		Type:         ExcelTypeString,
		IsRequired:   false,
		DefaultValue: defaultValue,
		MaxLength:    maxLength,
	}
}

// NewEnumColumn creates a column with allowed values (enum-like)
func NewEnumColumn(key, header string, required bool, allowedValues []string, defaultValue string) ExcelColumn {
	column := ExcelColumn{
		Key:           key,
		Header:        header,
		Type:          ExcelTypeString,
		IsRequired:    required,
		AllowedValues: allowedValues,
	}
	if defaultValue != "" {
		column.DefaultValue = defaultValue
	}
	return column
}

// NewEmailColumn creates an email validation column
func NewEmailColumn(key, header string, required bool) ExcelColumn {
	return ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeEmail,
		IsRequired: required,
	}
}

// NewIntColumn creates an integer column with optional min/max constraints
func NewIntColumn(key, header string, required bool, minValue, maxValue *float64) ExcelColumn {
	return ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeInt,
		IsRequired: required,
		MinValue:   minValue,
		MaxValue:   maxValue,
	}
}

// NewFloatColumn creates a float column with optional min/max constraints
func NewFloatColumn(key, header string, required bool, minValue, maxValue *float64) ExcelColumn {
	return ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeFloat,
		IsRequired: required,
		MinValue:   minValue,
		MaxValue:   maxValue,
	}
}

// NewBoolColumn creates a boolean column with optional default value
func NewBoolColumn(key, header string, required bool, defaultValue *bool) ExcelColumn {
	column := ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeBool,
		IsRequired: required,
	}
	if defaultValue != nil {
		column.DefaultValue = *defaultValue
	}
	return column
}

// NewDateColumn creates a date column with optional custom format
func NewDateColumn(key, header string, required bool, dateFormat string) ExcelColumn {
	column := ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeDate,
		IsRequired: required,
	}
	if dateFormat != "" {
		column.DateFormat = dateFormat
	}
	return column
}

// NewDateTimeColumn creates a datetime column with optional custom format
func NewDateTimeColumn(key, header string, required bool, dateTimeFormat string) ExcelColumn {
	column := ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeDateTime,
		IsRequired: required,
	}
	if dateTimeFormat != "" {
		column.DateFormat = dateTimeFormat
	}
	return column
}

// NewURLColumn creates a URL validation column
func NewURLColumn(key, header string, required bool) ExcelColumn {
	return ExcelColumn{
		Key:        key,
		Header:     header,
		Type:       ExcelTypeURL,
		IsRequired: required,
	}
}

// Predefined column sets for common entities

// GetBasicOrderColumns returns standard columns for order export
func GetBasicOrderColumns() []ExcelColumn {
	return []ExcelColumn{
		NewRequiredStringColumn("orderId", "Order ID", IntPtr(50)),
		NewRequiredStringColumn("status", "Status", IntPtr(50)),
		NewFloatColumn("totalPrice", "Total Price", true, Float64Ptr(0), nil),
		NewRequiredStringColumn("buyerName", "Buyer Name", IntPtr(100)),
		NewOptionalStringColumn("vendorName", "Vendor Name", "", IntPtr(100)),
		NewRequiredStringColumn("paymentStatus", "Payment Status", IntPtr(50)),
		NewDateTimeColumn("createdAt", "Created Date", true, "2006-01-02 15:04:05"),
	}
}

// GetBasicProductColumns returns standard columns for product import/export
func GetBasicProductColumns() []ExcelColumn {
	return []ExcelColumn{
		NewRequiredStringColumn("name", "Product Name", IntPtr(100)),
		NewOptionalStringColumn("description", "Description", "", IntPtr(500)),
		NewFloatColumn("price", "Price", true, Float64Ptr(0), nil),
		NewIntColumn("quantity", "Quantity", false, Float64Ptr(0), nil),
		NewEnumColumn("status", "Status", true, []string{"Active", "Inactive", "Draft"}, "Active"),
	}
}

// GetBasicUserColumns returns standard columns for user import/export
func GetBasicUserColumns() []ExcelColumn {
	return []ExcelColumn{
		NewRequiredStringColumn("fullName", "Full Name", IntPtr(100)),
		NewEmailColumn("email", "Email", true),
		NewOptionalStringColumn("phone", "Phone", "", IntPtr(20)),
		NewEnumColumn("type", "User Type", true, []string{"Admin", "Vendor", "Buyer"}, "Buyer"),
		NewBoolColumn("isActive", "Active", false, BoolPtr(true)),
	}
}

// GetBasicCategoryColumns returns standard columns for category import/export
func GetBasicCategoryColumns() []ExcelColumn {
	return []ExcelColumn{
		NewRequiredStringColumn("name", "Category Name", IntPtr(100)),
		NewOptionalStringColumn("description", "Description", "", IntPtr(500)),
		NewOptionalStringColumn("parentCategory", "Parent Category", "", IntPtr(100)),
		NewBoolColumn("isActive", "Active", false, BoolPtr(true)),
	}
}
