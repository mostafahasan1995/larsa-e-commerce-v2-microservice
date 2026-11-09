// Package helpers provides utility functions for the application.
// This package includes validation helpers for generating user-friendly error messages.
package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// getErrorMessage returns a human-readable error message based on the validation tag
// that failed. It supports common validation tags like required, email, min, max, etc.
//
// Parameters:
//   - err: The validation error containing the tag and other details
//
// Returns:
//   - A string containing a user-friendly error message
func GetErrorMessage(err validator.FieldError) string {
	fieldName := err.Field()
	errorTag := err.Tag()

	// Handle nested field errors
	if err.Namespace() != "" {
		fieldName = err.Namespace()
	}

	switch errorTag {
	case "required":
		return fmt.Sprintf("The field '%s' is required", fieldName)
	case "email":
		return fmt.Sprintf("The field '%s' must be a valid email address", fieldName)
	case "min":
		switch err.Type().String() {
		case "string":
			return fmt.Sprintf("The field '%s' must be at least %s characters long", fieldName, err.Param())
		case "int", "int32", "int64", "uint", "uint32", "uint64":
			return fmt.Sprintf("The field '%s' must be at least %s", fieldName, err.Param())
		case "float32", "float64":
			return fmt.Sprintf("The field '%s' must be at least %s", fieldName, err.Param())
		case "slice", "array":
			// Handle slice/array validation more gracefully
			if fieldName == "addons" {
				return "Please select at least one addon option"
			}
			return fmt.Sprintf("Please select at least %s items for '%s'", err.Param(), fieldName)
		default:
			return fmt.Sprintf("The field '%s' must be at least %s", fieldName, err.Param())
		}
	case "max":
		switch err.Type().String() {
		case "string":
			return fmt.Sprintf("The field '%s' must not exceed %s characters", fieldName, err.Param())
		case "int", "int32", "int64", "uint", "uint32", "uint64":
			return fmt.Sprintf("The field '%s' must not exceed %s", fieldName, err.Param())
		case "float32", "float64":
			return fmt.Sprintf("The field '%s' must not exceed %s", fieldName, err.Param())
		default:
			return fmt.Sprintf("The field '%s' must not exceed %s", fieldName, err.Param())
		}
	case "len":
		return fmt.Sprintf("Value must be exactly %s characters long", err.Param())
	case "url":
		return "Invalid URL format"
	case "uri":
		return "Invalid URI format"
	case "uuid":
		return "Invalid UUID format"
	case "uuid4":
		return "Invalid UUID v4 format"
	case "uuid5":
		return "Invalid UUID v5 format"
	case "json":
		return "Invalid JSON format"
	case "jwt":
		return "Invalid JWT format"
	case "datetime":
		return "Invalid datetime format"
	case "date":
		return "Invalid date format"
	case "timezone":
		return "Invalid timezone"
	case "tcp_addr":
		return "Invalid TCP address"
	case "tcp4_addr":
		return "Invalid TCP4 address"
	case "tcp6_addr":
		return "Invalid TCP6 address"
	case "udp_addr":
		return "Invalid UDP address"
	case "udp4_addr":
		return "Invalid UDP4 address"
	case "udp6_addr":
		return "Invalid UDP6 address"
	case "ip":
		return "Invalid IP address"
	case "ipv4":
		return "Invalid IPv4 address"
	case "ipv6":
		return "Invalid IPv6 address"
	case "tcp":
		return "Invalid TCP address"
	case "udp":
		return "Invalid UDP address"
	case "oneof":
		return fmt.Sprintf("The field '%s' must be one of: %s", fieldName, err.Param())
	case "unique":
		return "Value must be unique"
	case "dive":
		return fmt.Sprintf("The field '%s' contains invalid values", fieldName)
	case "excludes":
		return fmt.Sprintf("Value must not contain: %s", err.Param())
	case "excludesall":
		return fmt.Sprintf("Value must not contain any of: %s", err.Param())
	case "excludesrune":
		return fmt.Sprintf("Value must not contain: %s", err.Param())
	case "contains":
		return fmt.Sprintf("Value must contain: %s", err.Param())
	case "containsany":
		return fmt.Sprintf("Value must contain at least one of: %s", err.Param())
	case "containsrune":
		return fmt.Sprintf("Value must contain: %s", err.Param())
	case "startswith":
		return fmt.Sprintf("Value must start with: %s", err.Param())
	case "endswith":
		return fmt.Sprintf("Value must end with: %s", err.Param())
	case "startsnotwith":
		return fmt.Sprintf("Value must not start with: %s", err.Param())
	case "endsnotwith":
		return fmt.Sprintf("Value must not end with: %s", err.Param())
	case "isbn":
		return "Invalid ISBN format"
	case "isbn10":
		return "Invalid ISBN-10 format"
	case "isbn13":
		return "Invalid ISBN-13 format"
	case "issn":
		return "Invalid ISSN format"
	case "ascii":
		return "Value must contain only ASCII characters"
	case "printascii":
		return "Value must contain only printable ASCII characters"
	case "multibyte":
		return "Value must contain multibyte characters"
	case "datauri":
		return "Invalid data URI format"
	case "latitude":
		return "Invalid latitude value"
	case "longitude":
		return "Invalid longitude value"
	case "ssn":
		return "Invalid SSN format"
	case "hostname":
		return "Invalid hostname format"
	case "hostname_rfc1123":
		return "Invalid hostname format (RFC 1123)"
	case "fqdn":
		return "Invalid FQDN format"
	case "dir":
		return "Invalid directory path"
	case "file":
		return "Invalid file path"
	case "image":
		return "Invalid image format"
	case "alpha":
		return "Value must contain only letters"
	case "alphanum":
		return "Value must contain only letters and numbers"
	case "alphaunicode":
		return "Value must contain only Unicode letters"
	case "alphanumunicode":
		return "Value must contain only Unicode letters and numbers"
	case "numeric":
		return "Value must be numeric"
	case "lowercase":
		return "Value must be lowercase"
	case "uppercase":
		return "Value must be uppercase"
	case "number":
		return "Value must be a number"
	case "e164":
		return "Invalid E.164 formatted phone number"
	case "hexcolor":
		return "Invalid hex color code"
	case "rgb":
		return "Invalid RGB color code"
	case "rgba":
		return "Invalid RGBA color code"
	case "hsl":
		return "Invalid HSL color code"
	case "hsla":
		return "Invalid HSLA color code"
	case "base64":
		return "Invalid base64 string"
	case "base64url":
		return "Invalid base64url string"
	case "mongodb":
		return fmt.Sprintf("The field '%s' must be a valid MongoDB ObjectID (24 characters)", fieldName)
	case "iso3166_1_alpha2":
		return "Invalid ISO 3166-1 alpha-2 country code"
	case "iso3166_1_alpha3":
		return "Invalid ISO 3166-1 alpha-3 country code"
	case "iso3166_1_alpha_numeric":
		return "Invalid ISO 3166-1 alpha-numeric country code"
	case "iso3166_2":
		return "Invalid ISO 3166-2 country code"
	case "iso4217":
		return "Invalid ISO 4217 currency code"
	case "iso8601":
		return "Invalid ISO 8601 datetime format"
	case "rfc3339":
		return "Invalid RFC 3339 datetime format"
	case "rfc3339_without_zone":
		return "Invalid RFC 3339 datetime format (without timezone)"
	case "rfc1123":
		return "Invalid RFC 1123 datetime format"
	case "rfc1123z":
		return "Invalid RFC 1123Z datetime format"
	case "rfc822":
		return "Invalid RFC 822 datetime format"
	case "rfc822z":
		return "Invalid RFC 822Z datetime format"
	case "rfc2822":
		return "Invalid RFC 2822 datetime format"
	case "ulid":
		return "Invalid ULID format"
	case "cron":
		return "Invalid cron expression"
	case "gte":
		return fmt.Sprintf("The field '%s' must be greater than or equal to %s", fieldName, err.Param())
	case "lte":
		return fmt.Sprintf("The field '%s' must be less than or equal to %s", fieldName, err.Param())
	case "gt":
		return fmt.Sprintf("The field '%s' must be greater than %s", fieldName, err.Param())
	case "lt":
		return fmt.Sprintf("The field '%s' must be less than %s", fieldName, err.Param())
	default:
		return fmt.Sprintf("The field '%s' has an invalid value", fieldName)
	}
}
