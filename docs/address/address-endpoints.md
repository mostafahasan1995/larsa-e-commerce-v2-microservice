### address endpoints (web)

- Base: `/addresses`

- create address: POST `/addresses`
  - Body: `AddressDTO`
  - Example:
```json
{
  "name": { "en": "Head Office", "ar": "المقر الرئيسي" },
  "type": "office",
  "country": "AE",
  "city": "Dubai",
  "street": "Harbor St",
  "addressLine1": "Building 10, Floor 3",
  "notes": "Leave with reception",
  "poBox": "12345",
  "postCode": "00000",
  "phone": "+9715000000",
  "location": { "type": "Point", "coordinates": [55.2708, 25.2048] },
  "isDefault": true
}
```

- get all: GET `/addresses/all`

- list (pagination + filters): GET `/addresses`
  - Query params:
    - `page` (default 1)
    - `perPage` (default 40)
    - `query` (JSON string of filters)

- get by id: GET `/addresses/{id}`

- update by id: PUT `/addresses/{id}`
  - Body: `AddressDTO` (same shape as create)

- delete by id: DELETE `/addresses/{id}`

- bulk delete: DELETE `/addresses`
  - Body:
```json
{ "ids": ["66a7a0a6d1f1a2b3c4d5e6f7", "66a7a0a6d1f1a2b3c4d5e6f8"] }
```

- headers required for all requests:
  - `x-access-token`: `{{auth_token}}`
  - `x-client`: `{{client_id}}`
