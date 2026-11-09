package middleware

import (
	"context"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"net/http"
	"strings"
)

func AddLangToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		langHeader := r.Header.Get("accept-language")
		languages := strings.Split(langHeader, ",")
		if len(languages) > 0 {
			primaryLang := languages[0]
			langParts := strings.Split(primaryLang, "-")
			lang := langParts[0] // Extract the primary language code
			ctx := context.WithValue(r.Context(), enums.CtxLangKey, lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Handle case where no language is provided
			next.ServeHTTP(w, r)
		}
	})
}
