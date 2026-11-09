# Delete Guest Profiles API Response

## Endpoint
```http
DELETE /profile/guests
```

## Headers
```
x-access-token: {{auth_token}}
x-client: {{client_id}}
Content-Type: application/json
```

## Response Examples

### Success Response (HTTP 200)
```json
{
  "deletedCount": 15,
  "message": "Guest profiles deleted successfully"
}
```

### When No Guest Profiles Found (HTTP 200)
```json
{
  "deletedCount": 0,
  "message": "Guest profiles deleted successfully"
}
```

### Error Response (HTTP 500)
```json
{
  "error": "Database connection failed",
  "statusCode": 500
}
```

## Notes
- The `deletedCount` field returns the exact number of guest profiles that were deleted
- Guest profiles are identified by having `firstName` equal to "Guest"
- The operation requires SuperAdmin privileges (currently commented out in middleware)
- Returns HTTP 200 with count information instead of HTTP 204 (No Content)
