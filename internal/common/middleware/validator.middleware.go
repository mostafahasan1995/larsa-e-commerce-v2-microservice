package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey string

const (
	dtoContextKey contextKey = "dto"
)

// ValidateDTO is a middleware function that validates incoming request data
// against a provided DTO struct.
func ValidateDTO(dto interface{}, source PayloadSource, injector *do.Injector) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxAppCfg, _ := util.AddCtxAppCfg(r)
			dtoInstance, err := validateRequest(w, r, dto, source, injector)

			if err != nil {
				log.Error().Err(err).Msg("Error validating DTO")
				http.Error(w, "err validate dto", http.StatusBadRequest)
				return
			}

			if dtoInstance != nil {
				db := do.MustInvoke[*mongo.Client](injector)
				cfg, _ := util.GetReqAppCfg(ctxAppCfg)
				ctx := context.WithValue(r.Context(), dtoContextKey, dtoInstance)

				// Validate the entire structure recursively
				customContexts, apiErr := validateNestedStructure(reflect.ValueOf(dtoInstance), CustomValidatorArgs{
					Context:  ctxAppCfg,
					DB:       db,
					Request:  r,
					Injector: injector,
					User:     cfg.User,
				})
				if apiErr != nil {
					helpers.HandleApiError(w, apiErr)
					return
				}

				// Add all custom contexts to the request context
				for _, customData := range customContexts {
					ctx = context.WithValue(ctx, customData.Key, customData.Value)
				}

				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// createInstance creates a new instance of the DTO type
func CreateInstance(dto interface{}) interface{} {
	dtoType := helpers.GetType(dto)
	return helpers.CreateInstance(dtoType)
}

func validateRequest(w http.ResponseWriter, r *http.Request, dto interface{}, source PayloadSource, injector *do.Injector) (interface{}, error) {
	switch source {
	case BODY:
		return readFromBody(w, r, dto, injector)
	case QUERY:
		return readFromQuery(w, r)
	case PARAM:
		return readFromParams(w, r, dto, injector)
	default:
		return nil, nil
	}
}

func validateJSON(dto interface{}, bodyBytes []byte, injector *do.Injector) error {
	// Get the validator instance from the injector
	validate := do.MustInvoke[*validator.Validate](injector)

	if err := json.Unmarshal(bodyBytes, dto); err != nil {
		return err
	}

	return validate.Struct(dto)

}

func ValidateInstance(dto interface{}, instance interface{}, injector *do.Injector) error {
	dtoInstance := CreateInstance(dto)

	// Get the validator instance from the injector
	validate := do.MustInvoke[*validator.Validate](injector)

	instanceBytes, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(instanceBytes, &dtoInstance); err != nil {
		return err
	}

	return validate.Struct(dtoInstance)
}

func handleInstanceError(w http.ResponseWriter, err error) {
	var syntaxErr *json.SyntaxError
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &syntaxErr) {
		handleJsonError(w, err)
	} else if errors.As(err, &unmarshalErr) {
		handleUnmarshalTypeError(w, err)
	} else {
		handleValidationError(w, err)
	}
}

// GetValidatedDTO retrieves a validated DTO from the request context.
func GetValidatedDTO[T any](r *http.Request) (*T, *helpers.ApiError) {
	dto := r.Context().Value(dtoContextKey)
	if dto == nil {
		return nil, nil
	}

	if typedDTO, ok := dto.(*T); ok {
		return typedDTO, nil
	}

	return nil, helpers.DefaultError(errors.New("invalid dto"))
}

// GetValidatedDTO retrieves a validated DTO from the request context.
func GetContext[T any](r *http.Request, contextKey string) (*T, *helpers.ApiError) {
	dto := r.Context().Value(contextKey)
	if dto == nil {
		return nil, nil
	}

	if typedDTO, ok := dto.(*T); ok {
		return typedDTO, nil
	}

	return nil, helpers.DefaultError(errors.New("something went wrong"))
}
