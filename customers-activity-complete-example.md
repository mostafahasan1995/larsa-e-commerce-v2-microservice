# Customers Activity Endpoint - Complete Filter Example

## GET Request with All Available Filters

```bash
curl -X GET "http://localhost:8080/api/stats/customers-activity?startDate=2024-01-01T00:00:00.000Z&endDate=2024-12-31T23:59:59.999Z&minOrders=5&maxOrders=100&branchId=507f1f77bcf86cd799439011&customerId=507f1f77bcf86cd799439012&customerName=Ahmed&companyName=Tech%20Solutions&email=%40company.com&phoneNumber=555&minTotalQuantity=10&maxTotalQuantity=1000&minTotalSalesAmount=200.50&maxTotalSalesAmount=10000.00&minAverageOrderValue=50.00&maxAverageOrderValue=500.00&minUsedCouponsCount=1&maxUsedCouponsCount=20&customerType=returning&lastBranchOrderedFrom=Main%20Branch&minLastOrderDate=2024-06-01T00:00:00.000Z&maxLastOrderDate=2024-12-31T23:59:59.999Z&sortBy=totalMonthSalesAmount&sortDir=desc&page=1&limit=25" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>"
```

## Alternative: URL-Encoded Format

```
GET /api/stats/customers-activity?startDate=2024-01-01T00:00:00.000Z&endDate=2024-12-31T23:59:59.999Z&minOrders=5&maxOrders=100&branchId=507f1f77bcf86cd799439011&customerId=507f1f77bcf86cd799439012&customerName=Ahmed&companyName=Tech%20Solutions&email=%40company.com&phoneNumber=555&minTotalQuantity=10&maxTotalQuantity=1000&minTotalSalesAmount=200.50&maxTotalSalesAmount=10000.00&minAverageOrderValue=50.00&maxAverageOrderValue=500.00&minUsedCouponsCount=1&maxUsedCouponsCount=20&customerType=returning&lastBranchOrderedFrom=Main%20Branch&minLastOrderDate=2024-06-01T00:00:00.000Z&maxLastOrderDate=2024-12-31T23:59:59.999Z&sortBy=totalMonthSalesAmount&sortDir=desc&page=1&limit=25
```

## Filter Parameters Breakdown

### Date Range Filters

- `startDate=2024-01-01T00:00:00.000Z` - Start date for order period (ISO 8601 format)
- `endDate=2024-12-31T23:59:59.999Z` - End date for order period (ISO 8601 format)
- `minLastOrderDate=2024-06-01T00:00:00.000Z` - Minimum last order date (ISO 8601 format)
- `maxLastOrderDate=2024-12-31T23:59:59.999Z` - Maximum last order date (ISO 8601 format)

### Order Count Filters

- `minOrders=5` - Minimum number of orders (default: 5)
- `maxOrders=100` - Maximum number of orders

### Customer Information Filters

- `customerId=507f1f77bcf86cd799439012` - Specific customer ID (ObjectID)
- `customerName=Ahmed` - Partial match on customer name (case-insensitive)
- `companyName=Tech%20Solutions` - Partial match on company name (case-insensitive, URL encoded)
- `email=%40company.com` - Partial match on email (case-insensitive, URL encoded)
- `phoneNumber=555` - Partial match on phone number (case-insensitive)
- `customerType=returning` - Customer type: "new" or "returning"

### Quantity & Sales Filters

- `minTotalQuantity=10` - Minimum total quantity ordered
- `maxTotalQuantity=1000` - Maximum total quantity ordered
- `minTotalSalesAmount=200.50` - Minimum total sales amount
- `maxTotalSalesAmount=10000.00` - Maximum total sales amount
- `minAverageOrderValue=50.00` - Minimum average order value
- `maxAverageOrderValue=500.00` - Maximum average order value

### Coupon Usage Filters

- `minUsedCouponsCount=1` - Minimum number of coupons used
- `maxUsedCouponsCount=20` - Maximum number of coupons used

### Branch & Location Filters

- `branchId=507f1f77bcf86cd799439011` - Specific branch ID (ObjectID)
- `lastBranchOrderedFrom=Main%20Branch` - Partial match on last branch name (case-insensitive, URL encoded)

### Sorting & Pagination

- `sortBy=totalMonthSalesAmount` - Sort field (totalOrdersCount, totalMonthSalesAmount, totalQuantityOrdered, customerName, lastOrderDate, usedCouponsCount, averageOrderValue)
- `sortDir=desc` - Sort direction: "asc" or "desc"
- `page=1` - Page number for pagination
- `limit=25` - Number of results per page

## Expected Response Format

```json
{
  "data": [
    {
      "customerId": "507f1f77bcf86cd799439012",
      "customerName": "Ahmed Al-Rashid",
      "companyName": "Tech Solutions Inc",
      "email": "ahmed@techsolutions.com",
      "phoneNumber": "+000000000000",
      "totalOrdersCount": 15,
      "totalQuantityOrdered": 45,
      "totalMonthSalesAmount": 2500.75,
      "lastOrderDate": "2024-12-15T10:30:00Z",
      "usedCouponsCount": 3,
      "couponDetails": [
        {
          "code": "SAVE10",
          "usageDate": "2024-12-10T14:20:00Z",
          "discountValue": 25.5
        }
      ],
      "averageOrderValue": 166.72,
      "customerType": "returning",
      "lastBranchOrderedFrom": "Main Branch",
      "isTopCustomer": false
    }
  ],
  "pagination": {
    "totalPages": 1,
    "currentPage": 1,
    "totalCount": 1,
    "limit": 25
  }
}
```

## Usage Notes

1. **Replace placeholders:**

   - `YOUR_AUTH_TOKEN_HERE` with your actual authentication token
   - `localhost:8080` with your actual server URL
   - ObjectIDs with real customer/branch IDs from your database

2. **Date format support:**

   - Primary: ISO 8601 format (`2024-01-01T00:00:00.000Z`)
   - Fallback: Simple date format (`2024-01-01`) for backward compatibility

3. **Filter combinations:**

   - All filters are optional and can be combined as needed
   - Filters work independently and can be mixed and matched
   - Empty or invalid filter values are ignored

4. **URL encoding:**

   - Spaces in company names and branch names are encoded as `%20`
   - Special characters in emails are properly encoded
   - Use proper URL encoding when constructing requests programmatically

5. **Performance considerations:**
   - Use specific filters to reduce result set size
   - Combine date ranges with other filters for better performance
   - Consider pagination for large result sets
