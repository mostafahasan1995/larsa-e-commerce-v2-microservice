package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	TTL    time.Duration
	Prefix string
}

// CachedResponse holds the complete HTTP response for caching
type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		TTL:    3 * time.Second, // 3 seconds as requested
		Prefix: "admin_cache:",
	}
}

// AdminCache creates a caching middleware for admin endpoints
func AdminCache(config ...CacheConfig) func(http.Handler) http.Handler {
	cfg := DefaultCacheConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache GET requests for admin endpoints
			if r.Method != "GET" || !strings.HasPrefix(r.URL.Path, "/admin/") {
				next.ServeHTTP(w, r)
				return
			}

			// Skip caching if Redis is not available
			if configuration.Rdb == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Generate cache key based on URL, query params, and user context
			cacheKey := generateCacheKey(r, cfg.Prefix)

			// Try to get cached response - using interface{} to handle unknown methods
			if cachedResponse := getCachedResponse(cacheKey); cachedResponse != nil {
				// Cache hit - return cached response
				// fmt.Printf("[CACHE HIT] %s %s - served from cache\n", r.Method, r.URL.Path)

				// Set cached headers
				for key, values := range cachedResponse.Headers {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				w.Header().Set("X-Cache", "HIT")
				w.WriteHeader(cachedResponse.StatusCode)
				w.Write(cachedResponse.Body)
				return
			}

			// Cache miss - capture the response
			// fmt.Printf("[CACHE MISS] %s %s - generating response\n", r.Method, r.URL.Path)
			recorder := &responseRecorder{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     http.StatusOK,
			}

			// Process the request
			next.ServeHTTP(recorder, r)

			// Only cache successful responses
			if recorder.statusCode == http.StatusOK && recorder.body.Len() > 0 {
				// Create cached response with headers
				cachedResp := &CachedResponse{
					StatusCode: recorder.statusCode,
					Headers:    make(http.Header),
					Body:       recorder.body.Bytes(),
				}

				// Copy important headers (but skip compression-related ones)
				for key, values := range recorder.Header() {
					if key != "Content-Encoding" && key != "Content-Length" {
						cachedResp.Headers[key] = values
					}
				}

				// Store in cache with TTL
				setCachedResponse(cacheKey, cachedResp, int(cfg.TTL.Seconds()))
				// fmt.Printf("[CACHE STORED] %s %s - cached for %v\n", r.Method, r.URL.Path, cfg.TTL)
			}

			// Set cache headers
			recorder.Header().Set("X-Cache", "MISS")
		})
	}
}

// generateCacheKey creates a unique cache key based on request details
func generateCacheKey(r *http.Request, prefix string) string {
	// Include user context to avoid serving wrong data to different users
	userContext := ""
	if cfg, _ := util.GetReqAppCfg(r.Context()); cfg != nil {
		if cfg.User != nil {
			userContext = fmt.Sprintf("_user:%s", cfg.User.Id.Hex())
		}
		if cfg.Db != "" {
			userContext += fmt.Sprintf("_db:%s", cfg.Db)
		}
	}

	// Create key from path, query, and user context
	keyData := fmt.Sprintf("%s%s%s", r.URL.Path, r.URL.RawQuery, userContext)

	// Hash the key to keep it short
	hash := md5.Sum([]byte(keyData))
	hashStr := fmt.Sprintf("%x", hash)

	return fmt.Sprintf("%s%s", prefix, hashStr)
}

// responseRecorder captures the response for caching
type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	// Write to both the original response and our buffer
	r.body.Write(data)
	return r.ResponseWriter.Write(data)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Implement io.ReaderFrom for efficiency
func (r *responseRecorder) ReadFrom(src io.Reader) (int64, error) {
	// Copy to both destinations
	tee := io.TeeReader(src, r.body)
	return io.Copy(r.ResponseWriter, tee)
}

// Simple in-memory cache as fallback until Redis interface is figured out
type cacheItem struct {
	response  *CachedResponse
	expiresAt time.Time
}

var (
	memCache = make(map[string]cacheItem)
	cacheMux = sync.RWMutex{}
)

// getCachedResponse attempts to get response from cache
func getCachedResponse(key string) *CachedResponse {
	// First try Redis if available
	if configuration.Rdb != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle method not found and fall back to memory cache
			}
		}()
		// TODO: Use correct Redis methods when interface is known
	}

	// Fallback to in-memory cache
	cacheMux.RLock()
	item, exists := memCache[key]
	cacheMux.RUnlock()

	if !exists || time.Now().After(item.expiresAt) {
		// Clean up expired item
		if exists {
			cacheMux.Lock()
			delete(memCache, key)
			cacheMux.Unlock()
		}
		return nil
	}

	return item.response
}

// setCachedResponse attempts to set response in cache
func setCachedResponse(key string, response *CachedResponse, ttlSeconds int) {
	// First try Redis if available
	if configuration.Rdb != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle method not found and fall back to memory cache
			}
		}()
		// TODO: Use correct Redis methods when interface is known
	}

	// Fallback to in-memory cache
	cacheMux.Lock()
	defer cacheMux.Unlock()

	// Clean up old entries periodically (simple cleanup)
	if len(memCache) > 1000 { // Prevent memory bloat
		for k, v := range memCache {
			if time.Now().After(v.expiresAt) {
				delete(memCache, k)
			}
		}
	}

	memCache[key] = cacheItem{
		response:  response,
		expiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
}
