package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"

	"errors"

	"github.com/go-playground/validator/v10"
)

// sendValidationError sends a validation error response with the specified status,
// message, and list of validation errors.
func sendValidationError(w http.ResponseWriter, apiError *helpers.ApiError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(apiError)
}

// handleValidationError processes validation errors and sends appropriate response.
func handleValidationError(w http.ResponseWriter, err error) {
	validationErrors := make([]types.ValidationError, 0)

	var validErrs validator.ValidationErrors
	var timeErr *time.ParseError
	if errors.As(err, &validErrs) {
		// Handle validation errors
		for _, fieldError := range validErrs {
			// Get the field path including array indices
			fieldPath := fieldError.Field()

			// Handle nested field errors
			if fieldError.Namespace() != "" {
				fieldPath = fieldError.Namespace()
			}

			// Get the actual value that caused the error
			value := fieldError.Value()
			if value == nil {
				value = "null"
			}

			validationErrors = append(validationErrors, types.ValidationError{
				Field:   fieldPath,
				Value:   value,
				Tag:     fieldError.Tag(),
				Message: helpers.GetErrorMessage(fieldError),
			})
		}
	} else if errors.As(err, &timeErr) {
		validationErrors = append(validationErrors, types.ValidationError{
			Field:   "time",
			Value:   timeErr.Value,
			Tag:     "time_format",
			Message: fmt.Sprintf("Invalid time format: %s", timeErr.Error()),
		})
	} else {
		// Handle other error types
		validationErrors = append(validationErrors, types.ValidationError{
			Field:   "unknown",
			Value:   nil,
			Tag:     "error",
			Message: err.Error(),
		})
	}

	sendValidationError(w, &helpers.ApiError{
		StatusCode: http.StatusBadRequest,
		Message:    "Validation failed",
		Details:    &validationErrors,
	})
}

func handleJsonError(w http.ResponseWriter, err error) {
	validationErrors := make([]types.ValidationError, 0)

	var jsonErr *json.SyntaxError
	if errors.As(err, &jsonErr) {
		validationErrors = append(validationErrors, types.ValidationError{
			Field:   "body",
			Value:   nil,
			Tag:     "json",
			Message: fmt.Sprintf("Invalid JSON format at position %d: %s", jsonErr.Offset, jsonErr.Error()),
		})
	}

	sendValidationError(w, &helpers.ApiError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request data",
		Details:    &validationErrors,
	})
}

func handleUnmarshalTypeError(w http.ResponseWriter, err error) {
	validationErrors := make([]types.ValidationError, 0)

	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		validationErrors = append(validationErrors, types.ValidationError{
			Field: unmarshalErr.Field,
			Value: nil,
			Tag:   "json",
			Message: fmt.Sprintf("Invalid type for field '%s': expected %s, got %s",
				unmarshalErr.Field, unmarshalErr.Type, unmarshalErr.Value),
		})
	}

	sendValidationError(w, &helpers.ApiError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid request data",
		Details:    &validationErrors,
	})
}
