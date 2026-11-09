### shipping endpoints (admin)

- Base: `/admin/shipping-companies`

- create shipping company: POST `/admin/shipping-companies`
  - Body: `CreateShippingDto`
  - Example:
```json
{
  "basicInformation": {
    "name": { "en": "Acme Shipping", "ar": "اكمي للشحن" },
    "type": { "en": "International", "ar": "دولي" },
    "commercialRegistrationNumber": 1234567890,
    "taxIdentificationNumber": 987654321
  },
  "service": {
    "landFreight": true,
    "airFreight": false,
    "seaFreight": true,
    "warehousingAndDistribution": true
  },
  "contract": {
    "email": "ops@acme-shipping.example",
    "headOfficeAddress": "456 Harbor Rd, Dubai",
    "phone": "+000000000000",
    "countryCodes": ["AE", "US"],
    "website": "https://acme-shipping.example"
  },
  "financial": {
    "startDate": "2025-01-01T00:00:00Z",
    "endDate": "2026-01-01T00:00:00Z",
    "bankInformation": "Example Bank - IBAN AE070331234567890123456",
    "attachments": [
      { "_id": "60d21b4667d0d8992e610c85", "path": "https://cdn.example.com/contracts/contract1.pdf" }
    ]
  }
}
```

- get by id: GET `/admin/shipping-companies/{id}`
  - Query: `query` (JSON string of filters, optional)

- list (pagination + filters): GET `/admin/shipping-companies`
  - Query params:
    - `page` (default 1)
    - `perPage` (default 40)
    - `query` (JSON string of filters)
  - Filters in `query` (all optional):
    - `name` (localizable), `code`, `country`
  - Example: `?page=1&perPage=10&query={"country":"AE"}`

- update by id: PUT `/admin/shipping-companies/{id}`
  - Body: `CreateShippingDto` (same shape as create)

- delete by id: DELETE `/admin/shipping-companies/{id}`

- bulk delete: DELETE `/admin/shipping-companies`
  - Body:
```json
{ "ids": ["66a7a0a6d1f1a2b3c4d5e6f7", "66a7a0a6d1f1a2b3c4d5e6f8"] }
```

- headers required for all requests:
  - `x-access-token`: `{{auth_token}}`
  - `x-client`: `{{client_id}}`
