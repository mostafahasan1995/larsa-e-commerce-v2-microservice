package middleware

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"net/http"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

// PayloadSource represents where the data should be read from in the HTTP request.
type PayloadSource string

const (
	// BODY indicates that data should be read from the request body
	BODY PayloadSource = "body"
	// QUERY indicates that data should be read from query parameters
	QUERY PayloadSource = "query"
	// PARAM indicates that data should be read from URL parameters
	PARAM PayloadSource = "param"
)

type CustomValidatorArgs struct {
	Context  context.Context
	DB       *mongo.Client
	Request  *http.Request
	Injector *do.Injector
	User     *types.User
}

type CustomValidatorContext struct {
	Key   string
	Value interface{}
}

// CustomValidator is an interface that DTOs can implement to provide custom validation
type CustomValidator interface {
	CustomValidator(args CustomValidatorArgs) ([]CustomValidatorContext, *helpers.ApiError)
}
