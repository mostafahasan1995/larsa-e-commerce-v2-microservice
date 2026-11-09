# Logging Locations and Reduction Guide

This document lists all logging locations in your codebase and suggests which ones you can safely disable or reduce.

## 📍 Main Logging Locations

### 1. **Cache Logs** (High Volume)
**File:** `internal/common/middleware/cache.middleware.go`

**Lines:**
- Line 64: `[CACHE HIT]` - When serving from cache
- Line 79: `[CACHE MISS]` - When generating new response
- Line 107: `[CACHE STORED]` - When caching response

**Impact:** Very high - logs every GET request to admin endpoints
**Recommendation:** ✅ **DISABLE or reduce to debug level only**

```go
// To disable, comment out lines 64, 79, and 107:
// fmt.Printf("[CACHE HIT] %s %s - served from cache\n", r.Method, r.URL.Path)
// fmt.Printf("[CACHE MISS] %s %s - generating response\n", r.Method, r.URL.Path)
// fmt.Printf("[CACHE STORED] %s %s - cached for %v\n", r.Method, r.URL.Path, cfg.TTL)
```

---

### 2. **Performance Logs** (Very High Volume)
**File:** `internal/common/middleware/performance.middleware.go`

**Lines:**
- Line 24-39: `[PERFORMANCE]` - Logs every single HTTP request with duration
- Line 36: `[SLOW REQUEST]` - Additional log for requests > 1 second

**Impact:** Very high - logs EVERY HTTP request
**Recommendation:** ✅ **Keep only SLOW REQUEST logs, disable regular performance logs**

```go
// Modify to log only slow requests:
if duration > time.Second {
    logMessage := fmt.Sprintf("[PERFORMANCE] %s %s - %d - %v",
        r.Method, r.URL.Path, wrapper.statusCode, duration)
    if r.URL.RawQuery != "" {
        logMessage += fmt.Sprintf("\n  Query: %s", r.URL.RawQuery)
    }
    logMessage += " [SLOW REQUEST]"
    fmt.Println(logMessage)
}
// Comment out the normal performance log at line 39
```

---

### 3. **Database Performance Logs** (Medium Volume)
**File:** `internal/common/configuration/monitering-hook.go`

**Lines:**
- Line 197-203: `[DB PERFORMANCE] SLOW QUERY` - Logs queries > 100ms

**Impact:** Medium - only logs slow queries
**Recommendation:** ⚠️ **Keep but consider increasing threshold to 200-500ms**

```go
// Increase threshold on line 151:
if args.Duration.Milliseconds() > 200 { // Changed from 100ms to 200ms
```

---

### 4. **Settings Upsert Logs** (Startup Only)
**File:** `internal/modules/config/repos/setting.repo.go`

**Lines:**
- Line 85: `"Starting settings upsert process..."`
- Line 113: `"Indexes upsert completed successfully"`
- Line 146: `"Upserting Setting for Db: %v, collections: %v"`

**Impact:** Low - only runs at startup
**Recommendation:** ✅ **Keep or reduce to single startup message**

---

### 5. **Startup Logs** (Startup Only)
**Files:**
- `internal/api/api.go`
  - Line 59: `"start server on port:"`
  - Line 165: `"start server"`
  
- `internal/common/configuration/db.go`
  - Line 66: `"MongoDB successfully connected and pinged."`
  - Line 68: `"Config Mode: %v"`
  
- `internal/common/configuration/nats.go`
  - Line 39, 42, 61, 64: `"connected to nats server"` (multiple times)
  
- `internal/common/nats/listeners.go`
  - Line 56: `"starting nats jet listeners to:"`

**Impact:** Very low - only at startup
**Recommendation:** ✅ **Keep these - they're useful for startup verification**

---

### 6. **Other Scattered Logs** (Low to Medium Impact)
Found in 57 files using `fmt.Printf`, `fmt.Println`, `log.Printf`:

**High-impact files to check:**
- `internal/modules/order/service/order-web.service.go`
- `internal/modules/product/services/branch.service.go`
- `internal/modules/carts/services/checkout.service.go`
- `internal/gateway/*.go` files

**Recommendation:** Review these individually based on your needs

---

## 🎯 Quick Action Plan to Reduce Logs

### Priority 1: Immediate Impact (95% reduction)
1. **Disable cache logs** (lines 64, 79, 107 in `cache.middleware.go`)
2. **Disable regular performance logs** (keep only slow requests in `performance.middleware.go`)

### Priority 2: Medium Impact
3. **Increase DB slow query threshold** from 100ms to 200-500ms in `monitering-hook.go`
4. **Consolidate settings logs** to single startup message in `setting.repo.go`

### Priority 3: Optional
5. Review individual service logs in the 57 files listed above

---

## 🔧 Recommended Configuration

### Option A: Production (Minimal Logging)
```go
// Cache: Disabled
// Performance: Only requests > 2 seconds
// DB Performance: Only queries > 500ms
// Startup: Enabled (keep as-is)
```

### Option B: Development (Balanced)
```go
// Cache: Disabled or only MISS/STORED
// Performance: Only requests > 1 second
// DB Performance: Only queries > 200ms
// Startup: Enabled
```

### Option C: Debug (Verbose)
```go
// Cache: All logs enabled
// Performance: All requests
// DB Performance: Queries > 100ms
// Startup: Enabled
```

---

## 📊 Estimated Log Volume Reduction

| Change | Current Volume | After Change | Reduction |
|--------|---------------|--------------|-----------|
| Disable cache logs | ~1000/min | 0 | 100% |
| Only slow performance logs | ~500/min | ~10/min | 98% |
| Increase DB threshold | ~50/min | ~10/min | 80% |
| **Total** | **~1550/min** | **~20/min** | **~98.7%** |

---

## 🛠️ Implementation Files to Edit

1. `internal/common/middleware/cache.middleware.go` (Lines: 64, 79, 107)
2. `internal/common/middleware/performance.middleware.go` (Lines: 24-39)
3. `internal/common/configuration/monitering-hook.go` (Line: 151)
4. `internal/modules/config/repos/setting.repo.go` (Lines: 85, 113, 146)

---

## 📝 Next Steps

1. Decide which logging level you want (Production/Development/Debug)
2. Edit the files listed above according to your choice
3. Test to ensure critical errors are still logged
4. Monitor production to verify log reduction

---

**Note:** The HTTP request logs you see (like `"GET http://services-laffah.larsa.io/admin/orders HTTP/1.1"`) come from the chi router middleware, not your code. These are standard access logs and are usually good to keep.

