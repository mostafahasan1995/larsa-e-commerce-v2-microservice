### agent endpoints (admin)

- Base: `/admin/agents`

- create agent: POST `/admin/agents`
  - Body: `CreateAgentDto`
  - Example:
```json
{
  "responsibleForLogistics": true,
  "active": true,
  "basicInformation": {
    "fullName": "Jane Doe",
    "companyName": "Acme Brokers",
    "licenseNumber": "LIC-12345",
    "yearsOfExperience": 5
  },
  "businessInfo": {
    "commercialRegistrationNumber": "CRN-987654",
    "taxIdentificationNumber": "TIN-456789",
    "phoneNumber": "+14155550123",
    "headOfficeAddress": "123 Market St, SF, CA",
    "officialEmail": "agent@example.com",
    "authorizedCountryCodes": ["US", "AE"],
    "websiteLink": "https://acme-brokers.example"
  },
  "financialInfo": {
    "shippingProfitPercentage": 7.5,
    "subscriptionProfitPercentage": 12.5
  },
  "contractInfo": {
    "startDate": "2025-01-01T00:00:00Z",
    "endDate": "2026-01-01T00:00:00Z",
    "attachments": [
      { "_id": "60d21b4667d0d8992e610c85", "path": "https://cdn.example.com/contracts/contract1.pdf" }
    ]
  },
  "bankAccountInfo": {
    "bankName": "Example Bank",
    "accountNumber": "1234567890",
    "accountHolderName": "Jane Doe",
    "iban": "AE070331234567890123456",
    "swiftCode": "EXAMPLEDX"
  },
  "attachments": [
    { "_id": "60d21b4667d0d8992e610c86", "path": "https://cdn.example.com/ids/id1.png" }
  ]
}
```

- get by id: GET `/admin/agents/{id}`
  - Query: `query` (JSON string of filters, optional)

- list (pagination + filters): GET `/admin/agents`
  - Query params:
    - `page` (default 1)
    - `perPage` (default 40)
    - `query` (JSON string of filters)
  - Filters in `query` (all optional):
    - `fullName`, `companyName`, `licenseNumber`, `yearsOfExperience`
    - `commercialRegistrationNumber`, `taxIdentificationNumber`, `phoneNumber`, `headOfficeAddress`, `officialEmail`, `authorizedCountryCodes` (array), `websiteLink`
    - `bankName`, `accountNumber`, `accountHolderName`, `iban`, `swiftCode`
  - Example: `?page=1&perPage=10&query={"fullName":"Jane Doe"}`

- update by id: PUT `/admin/agents/{id}`
  - Body: `CreateAgentDto` (same as create)

- delete by id: DELETE `/admin/agents/{id}`

- bulk delete: DELETE `/admin/agents`
  - Body:
```json
{ "ids": ["66a7a0a6d1f1a2b3c4d5e6f7", "66a7a0a6d1f1a2b3c4d5e6f8"] }
```

headers required for all requests:
- `x-access-token`: `{{auth_token}}`
- `x-client`: `{{client_id}}`

Authorization: Admin routes require bearer token and appropriate capabilities.
