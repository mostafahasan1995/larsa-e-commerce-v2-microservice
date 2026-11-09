package middleware

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"net/http"
)

func AddQueryKeysToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryKeys := r.URL.Query()
		ctx := context.WithValue(r.Context(), enums.CtxQueryKeysKey, queryKeys)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
