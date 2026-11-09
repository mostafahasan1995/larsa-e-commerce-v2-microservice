package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// RequestTiming middleware measures and logs HTTP request duration
func RequestTiming(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process the request
		next.ServeHTTP(wrapper, r)

		// Calculate duration
		duration := time.Since(start)

		// Enhanced logging with query parameters for slow requests
		logMessage := fmt.Sprintf("[PERFORMANCE] %s %s - %d - %v",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
		)

		// Add query parameters for slow requests (>1 second) to help identify patterns
		if duration > time.Second {
			if r.URL.RawQuery != "" {
				logMessage += fmt.Sprintf("\n  Query: %s", r.URL.RawQuery)
			}
			logMessage += " [SLOW REQUEST]"
		}

		fmt.Println(logMessage)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
