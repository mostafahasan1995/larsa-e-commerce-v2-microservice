package middleware

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"net/http"

	"github.com/samber/do"
)

func AddInjectorToContext(injector *do.Injector) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), enums.CtxInjectorKey, injector)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
