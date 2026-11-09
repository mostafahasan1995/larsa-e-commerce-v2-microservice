package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func readFromBody(w http.ResponseWriter, r *http.Request, dto interface{}, injector *do.Injector) (interface{}, error) {
	dtoInstance := CreateInstance(dto)
	var bodyBytes []byte
	var err error
	bodyBytes, err = io.ReadAll(r.Body)
	if err != nil {
		handleJsonError(w, err)
		return nil, err
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var data interface{}

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))

	if err := decoder.Decode(&data); err != nil {
		handleJsonError(w, err)
		return nil, err
	}

	if err := validateJSON(dtoInstance, bodyBytes, injector); err != nil {
		handleInstanceError(w, err)

		return nil, err
	}

	return dtoInstance, nil
}

// readFromQuery reads and decodes JSON data from query parameters.
func readFromQuery(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	queryStr := r.URL.Query().Get("query")

	if queryStr != "" {
		var data interface{}

		if err := json.Unmarshal([]byte(queryStr), &data); err != nil {
			handleJsonError(w, err)
			return nil, err
		}
	}

	return nil, nil
}

// readFromParams reads and decodes data from URL parameters.
func readFromParams(w http.ResponseWriter, r *http.Request, dto interface{}, injector *do.Injector) (interface{}, error) {
	id := chi.URLParam(r, "id")
	validationErrors := make([]types.ValidationError, 0)

	_, err := primitive.ObjectIDFromHex(id)
	apiErr := &types.ValidationError{
		Field:   "id",
		Value:   id,
		Tag:     "objectId",
		Message: "the provided hex string is not a valid ObjectID",
	}

	validationErrors = append(validationErrors, types.ValidationError{
		Field:   "id",
		Value:   id,
		Tag:     "json",
		Message: apiErr.Message,
	})

	if err != nil {
		sendValidationError(w, &helpers.ApiError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request data",
			Details:    &validationErrors,
		})
		return nil, apiErr
	}

	return nil, nil
}
