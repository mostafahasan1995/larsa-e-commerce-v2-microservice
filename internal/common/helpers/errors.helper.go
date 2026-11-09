package helpers

import (
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"net/http"
	"strings"
	"time"
)

type ApiError struct {
	StatusCode int                      `json:"statusCode"`
	Message    string                   `json:"message"`
	Details    *[]types.ValidationError `json:"details,omitempty"`
	Timestamp  *string                  `json:"timestamp,omitempty"`
}

func DefaultError(err error) *ApiError {
	apiError := &ApiError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Something went wrong",
	}

	if err != nil {
		// Only log non-duplicate key errors (E11000 errors are expected during settings upsert)
		if !strings.Contains(err.Error(), "E11000 duplicate key error") {
			fmt.Printf("Error: %v\n", err.Error())
		}

		if configuration.EnvironmentVariables.Mode == "debug" {
			apiError.Details = &[]types.ValidationError{
				{
					Message: err.Error(),
				},
			}
		}
	}

	return apiError
}

func (e *ApiError) Error() string {
	if e == nil {
		return "nil ApiError"
	}
	return fmt.Sprintf("api error: %d - %s", e.StatusCode, e.Message)
}

func NewApiError(statusCode int, err error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}

func InvalidJSON() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		fmt.Errorf("invalid json request data"),
	)
}

func BadRequest(message string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		fmt.Errorf(message),
	)
}

func InternalServerError(message string) ApiError {
	return NewApiError(
		http.StatusInternalServerError,
		fmt.Errorf(message),
	)
}

func InvalidJsonCustomMessage(message string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		fmt.Errorf(message),
	)
}

func DataBaseInsertError() ApiError {
	return NewApiError(
		http.StatusInternalServerError,
		fmt.Errorf("database insert failed"),
	)
}

func DataBaseUpdateError() ApiError {
	return NewApiError(
		http.StatusInternalServerError,
		fmt.Errorf("database update failed"),
	)
}

func DataBaseDeleteError() ApiError {
	return NewApiError(
		http.StatusInternalServerError,
		fmt.Errorf("database delete failed"),
	)
}

func NotFoundError(message string) ApiError {
	return NewApiError(
		http.StatusNotFound,
		fmt.Errorf(message),
	)
}

func InvalidObjectId() ApiError {
	return NewApiError(
		http.StatusBadRequest,
		fmt.Errorf("invalid Object Id format"),
	)
}

func InvalidObjectIdCustomMessage(message string) ApiError {
	return NewApiError(
		http.StatusBadRequest,
		fmt.Errorf(message),
	)
}

func InvalidPagination() ApiError {
	return NewApiError(http.StatusBadRequest, fmt.Errorf("invalid Object Id format"))
}

func ValidationErrors(errors []types.ValidationError) *ApiError {
	return &ApiError{
		StatusCode: http.StatusBadRequest,
		Message:    "Validation errors",
		Details:    &errors,
	}
}

// HandleApiError is a generic error handler for HTTP responses
// It writes a standardized error response to the HTTP writer
func HandleApiError(w http.ResponseWriter, apiError *ApiError) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	apiError.Timestamp = &timestamp
	WriteJson(w, apiError.StatusCode, apiError)
}
