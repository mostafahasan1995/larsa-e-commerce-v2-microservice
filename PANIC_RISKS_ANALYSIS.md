# Potential Panic Risks Analysis

## Summary
Based on the logs and code inspection, I've identified several potential panic scenarios in the order module. While the middleware.Recoverer has been added to catch panics, these issues should be fixed to prevent runtime errors.

---

## Critical Issues (High Risk of Panic)

### 1. **Unchecked Type Assertions**

#### Location: `internal/modules/order/repos/order-reports.repo.go:850`
```go
item := order_models.ProductSales{
    Name: raw["name"].(string),  // ⚠️ PANIC RISK: No check if "name" exists or is string
}
```
**Risk**: If `raw["name"]` is nil or not a string, this will panic.

**Fix**:
```go
var name string
if n, ok := raw["name"].(string); ok {
    name = n
}
item := order_models.ProductSales{
    Name: name,
}
```

---

#### Location: `internal/modules/order/repos/order-reports.repo.go:1038-1048`
```go
doc := r.(bson.M)  // ⚠️ PANIC RISK: No check
report = append(report, order_models.ProductSalesReport{
    ProductName:     doc["productName"].(string),      // ⚠️ PANIC RISK
    CustomerName:    doc["customerName"].(string),     // ⚠️ PANIC RISK
    PaymentMethod:   order_enums.PaymentMethod(doc["paymentMethod"].(string)), // ⚠️ PANIC RISK
})
```
**Risk**: Multiple unchecked type assertions that will panic if fields are missing or wrong type.

**Fix**: Use helper functions or check each assertion:
```go
doc, ok := r.(bson.M)
if !ok {
    continue // or log error
}
productName, _ := doc["productName"].(string)
customerName, _ := doc["customerName"].(string)
paymentMethod, _ := doc["paymentMethod"].(string)
```

---

#### Location: `internal/modules/order/repos/order-reports.repo.go:1054-1060`
```go
doc := statDocs[0].(bson.M)  // ⚠️ PANIC RISK: No check
stat = order_models.ProductSalesStat{
    BestSellingProduct: doc["bestSellingProduct"].(string), // ⚠️ PANIC RISK
}
```
**Risk**: Double panic risk - array access without bounds check AND unchecked type assertion.

**Fix**:
```go
if len(statDocs) > 0 {
    if doc, ok := statDocs[0].(bson.M); ok {
        bestSelling, _ := doc["bestSellingProduct"].(string)
        stat = order_models.ProductSalesStat{
            BestSellingProduct: bestSelling,
        }
    }
}
```

---

#### Location: `internal/modules/order/repos/order.repo.go:540`
```go
result := res.result.([]bson.M)  // ⚠️ PANIC RISK: No check
```
**Risk**: If res.result is not `[]bson.M`, this will panic.

**Fix**:
```go
result, ok := res.result.([]bson.M)
if !ok {
    log.Printf("Unexpected result type for %s", res.name)
    continue
}
```

---

#### Location: `internal/modules/order/service/order-update.service.go:206`
```go
if paymentsUpdate, exists := updates["paymentInfo.payments"]; exists {
    updatedPayments = paymentsUpdate.([]order_models.Payments)  // ⚠️ PANIC RISK
}
```
**Risk**: Just checking if key exists doesn't guarantee the type is correct.

**Fix**:
```go
if paymentsUpdate, exists := updates["paymentInfo.payments"]; exists {
    if payments, ok := paymentsUpdate.([]order_models.Payments); ok {
        updatedPayments = payments
    } else {
        // Log error or use default
        updatedPayments = currentOrder.PaymentInfo.Payments
    }
}
```

---

### 2. **Channel Operations Issues**

#### Location: `internal/modules/order/repos/order.repo.go:438-519`
```go
resultChan := make(chan activityResult, 5)

// 5 goroutines launched
go func() { ... }()
go func() { ... }()
go func() { ... }()
go func() { ... }()
go func() { ... }()

// Collect results
for i := 0; i < 5; i++ {
    res := <-resultChan
}
```

**Risk**: 
- Channel is never closed
- If any goroutine panics before sending to channel, the receiver will deadlock
- No timeout on channel reads - could hang indefinitely if MongoDB is slow

**Fix**:
```go
resultChan := make(chan activityResult, 5)
defer close(resultChan)

// Use sync.WaitGroup for better control
var wg sync.WaitGroup
wg.Add(5)

go func() {
    defer wg.Done()
    defer func() {
        if r := recover(); r != nil {
            resultChan <- activityResult{"total", nil, fmt.Errorf("panic: %v", r)}
        }
    }()
    // ... actual work
}()

// Similar for other goroutines

// Use context with timeout
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

// Collect with timeout
for i := 0; i < 5; i++ {
    select {
    case res := <-resultChan:
        // Handle result
    case <-ctx.Done():
        return nil, &helpers.ApiError{
            StatusCode: 408,
            Message:    "Timeout waiting for activity results",
        }
    }
}
```

---

### 3. **Array/Slice Potential Index Issues**

#### Location: `internal/modules/order/repos/order-reports.repo.go:354-360`
```go
var results []order_models.PaginatedCustomersWithStats
if err := cursor.All(ctx, &results); err != nil || len(results) == 0 {
    return nil, &helpers.ApiError{...}
}

if len(results[0].Stats) == 0 {  // ⚠️ results[0] was checked above, but still good to be careful
    results[0].Stats = []order_models.OrderStat{{
        CustomersPurchasing: 0,
        TotalPurchases:      0,
        TotalQuantity:       0,
        TotalAmount:         0,
    }}
}
```
**Risk**: While checked above, accessing `results[0]` after the check could still fail if modified in between.

**Status**: Currently safe, but could be improved for clarity.

---

## Medium Risk Issues

### 4. **Type Assertion in GetRecentActivity**

#### Location: `internal/modules/order/repos/order.repo.go:137`
```go
if matchStage, exists := pipeline[0]["$match"]; exists {
    additionalFilters = matchStage.(bson.M)  // ⚠️ Could panic if not bson.M
}
```

**Fix**:
```go
if matchStage, exists := pipeline[0]["$match"]; exists {
    if filters, ok := matchStage.(bson.M); ok {
        additionalFilters = filters
    }
}
```

---

## Performance Issues (Related to Slow Queries)

### 5. **Aggregation Pipeline Performance**

The logs show consistent slow queries (173-774ms) on the orders collection. Issues:

1. **Missing Indexes**: The slow aggregations likely need compound indexes
2. **Population of Multiple Fields**: `popFields` contains many fields being populated
3. **No Query Result Caching**: Frequently repeated queries aren't cached

**Recommendations**:
```go
// Add indexes for commonly queried fields
// - branchIds
// - createdAt
// - status
// - type
// - buyerInfo.id

// Consider implementing query result caching for:
// - Recent orders by branch
// - Orders by status
// - Today's orders
```

---

## Prevention Checklist

✅ **Middleware.Recoverer** - DONE (Added in api.go)

🔧 **TODO**:
- [ ] Fix all unchecked type assertions with proper checks
- [ ] Add defer/recover to goroutines
- [ ] Close channels properly
- [ ] Add timeouts to channel operations
- [ ] Add database indexes for slow queries
- [ ] Implement query result caching
- [ ] Add nil pointer checks before dereferences
- [ ] Add bounds checking for array/slice access

---

## Monitoring Recommendations

1. **Add structured logging** for type assertion failures
2. **Add metrics** for slow query tracking
3. **Add alerts** for panic recovery events
4. **Add tracing** for aggregate query performance
5. **Monitor** goroutine counts and channel operations

---

## Priority

**CRITICAL (Fix Immediately)**:
- Lines 850, 1040, 1044, 1047, 1059 in order-reports.repo.go
- Line 540 in order.repo.go
- Line 206 in order-update.service.go

**HIGH (Fix Soon)**:
- Channel operations and goroutine management
- Add timeouts to aggregation queries

**MEDIUM (Fix When Possible)**:
- Performance optimization
- Add indexes
- Implement caching

---

Generated: 2025-12-01

