# Product Requirements Document (PRD) for E-Commerce Microservice

## Overview

This document outlines the requirements for an e-commerce microservice as part of a B2B platform for a manufacturing company. The system will be modular, scalable, and multi-tenant, using MongoDB and SaaS (Database per Tenant) architecture. Authentication and authorization will be handled externally via an existing user service with RBAC (Roles and Permissions). Each request will include an `x-client` header for tenant identification. The system will also support localization with the ability to define and manage languages in a picklist for different tenants. The system will cater to both B2B (Business-to-Business) and B2C (Business-to-Consumer) markets.

---

## System Components

### 1. **Vendor Dashboard & Super Admin Features**

#### **1.1 Inventory Management**

##### **1.1.1 Products**

- **Basic Information**

  - Name (StringLocalization)
  - Description (StringLocalization)
  - Tags (Array of StringLocalization)
  - Category (Reference to Category module)
  - SKU Code (String)
  - Availability (Boolean)
  - MOQ (Minimum Order Quantity - Integer)
  - Attributes (Dynamic Key-Value Pairs)

- **Pricing**

  - Unit Price (Decimal)
  - Bulk Pricing (Array of {Min Quantity, Max Quantity, Price})

- **Advanced Pricing**

  - Based on Conditions:
    - Customer Groups (Array of References)
    - Branches (Array of References)
    - Platforms (Array of References)
  - Price (Decimal)
  - Discount (Percentage or Fixed Amount)

- **Availability & Production**

  - Lead Time (Integer - Days)
  - Monthly Production Capacity (Integer)
  - Target Export Regions (Array of Strings)

- **Certifications**

  - Certifications (Array)
    - Certification Type (String)
    - Certification File (Attachment)

- **Shipping Details**

  - Production Time (Integer - Days)
  - Estimated Delivery Time (Integer - Days)
  - Shipping Methods (Enum: Air/Sea)

- **Attachments**

  - Product Images (Array of {URL, Type})
  - Product Videos (Array of {URL, Type})
  - Product Documents (Array of {URL, Type})

- **Customization Options**

  - Include Logo (Boolean)
  - Include Colors (Boolean)

- **Approval Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)

##### **1.1.2 Production Lines**

- **Basic Information**

  - Name (StringLocalization)
  - Description (StringLocalization)
  - Tags (Array of StringLocalization)
  - Category (Reference to Category module)
  - General Classification (Enum - TBD)
  - Country of Origin (String)
  - Establishment Year (Date)

- **Product Information**

  - Number of Sub-Products (Integer)
  - Names of Sub-Products (Array of References)
  - Customizable (Boolean)
  - Product Category (Reference to Category module)

- **Technology & Equipment**

  - Technology Used (String)
  - Advanced Techniques (String)

- **Availability & Production**

  - Lead Time (Integer - Days)
  - Monthly Production Capacity (Integer)
  - Target Export Markets (Array of Strings)
  - Shipping & Export Terms (String)

- **Attachments**

  - Product Images (Array of URLs)
  - Product Videos (Array of URLs)
  - Product Documents (Array of URLs)

- **Approval Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)
  - **Super Admin Actions:**
    - Approve/Reject submitted production lines
    - Modify production line details after approval
    - Delete production lines if necessary

##### **1.1.3 Categories**

- **Basic Information**

  - Name (StringLocalization)
  - Code (Array of StringLocalization)
  - Parent Category (Reference to Category module)
  - Visibility (Enum: Store, Website, Mobile App, …etc)

- **Media & Images**

  - Images (Array of URLs)

- **Categories Positions (Integer for Ordering)**

  - Order by position

- **SEO Optimization**

  - Meta Title (StringLocalization)
  - Meta Description (StringLocalization)
  - Meta Keywords (Array of StringLocalization)
  - Open Graph (OG) Title (StringLocalization)
  - Open Graph (OG) Description (StringLocalization)

- **Approval Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)
  - **Super Admin Actions:**
    - Approve/Reject category submissions
    - Edit pending categories
    - Modify or delete categories after approval

- **Hierarchy:** Categories will be structured in a hierarchical manner to allow for better organization and navigation.

##### **1.1.4 Attributes (Product-Specific)**

- **Definition**

  - Code (String) and unique
  - Label (StringLocalization)
  - Type (Enum: Select-One, Select-Multiple, Text, Number, Date, Boolean)
  - Configurable (Boolean - Determines if it can be modified per product)
  - Is Filterable (Boolean)
  - Is Comparable (Boolean)
  - Is Editable (Boolean - Case for Non-Editable)
  - Default Value (Depends on Type)
  - Required (Boolean)
  - Visibility (Enum: Store, Website, Mobile App, …etc)
  - Sortable (Boolean)
  - Selection Style (Enum: Dropdown, Checkbox, Radio Button)
  - Action (Enum: Add, Edit, Delete)
  - Status (Enum: Active, Pending, Inactive)
  - Tag (Enum: System, Custom - Indicates whether the attribute is system-generated or custom)

- **Approval Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)

##### **1.1.5 Units of Measurement**

- **Definition**

  - Name (StringLocalization)
  - Type (Enum: Reference then Ratio is 1, Smaller, Bigger)
  - Category (Enum: Weight, Length, Volume, Working Time, Surface Area)
  - Ratio (Decimal - Conversion Factor)

- **Approval Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)

### **1.2 Picklist** Super Admin only

##### **1.2.1 Picklist**

##### **1.2.2 Picklist Pricing**

- **Pricing Currency**

##### **1.2.3 approval Process**

- **Category Process**
  - Submit for Approval (Boolean, Configurable)
  - Approved by Super Admin (Configurable)

##### **1.2.2 Supported Languages**
- **Language Code**: (e.g., "en" for English, "ar" for Arabic, etc.)
- **Language Name**: (e.g., "English", "Arabic")
- **Default Language**: (Boolean, indicates the default language for the system)

## Notes

- **Multi-Tenancy Support**: Each tenant will have its own database instance.
- **Microservices Architecture**: The e-commerce microservice will integrate with external authentication, user, and email services.
- **Scalability & Flexibility**: The platform is designed to support multiple business models and use cases within the manufacturing industry.
- **Localization Support**: The system supports multi-language options, allowing for localization of product information, categories, pricing, and more.

## Project Structure

my-project/
├── Dockerfile
├── .env
├── .env.example
├── .gitignore
│── cmd/                         
│   ├── main.go                  
│── config/                      
│   ├── config.go                
│── internal/                    
│   ├── api/                     # API Utilities
│   │   ├── api.go
│   ├── caching/                 # Caching (Redis)
│   │   ├── caching.go
│   ├── database/                # Database connection logic
│   │   ├── mongo.go             
│   │   ├── db.go                
│   ├── enums/                   # Enumerations and Constants
│   │   ├── enums.go
│   │   ├── msg.go
│   ├── helpers/                 # Utility Functions
│   │   ├── errors.helper.go
│   │   ├── numeric_filter.go
│   │   ├── pagination.go
│   │   ├── response.helper.go
│   │   ├── validation.helper.go
│   ├── middleware/              # Middlewares
│   │   ├── auth.go
│   ├── nats/                    # NATS Messaging
│   │   ├── nats.go
│   ├── modules/                 # Feature Modules
│   │   ├── inventory/
│   │   │   ├── di/
│   │   │   │   ├── product.di.go
│   │   │   ├── enums/
│   │   │   │   ├── product-status.go
│   │   │   ├── handler/
│   │   │   │   ├── product.handler.go
│   │   │   ├── models/
│   │   │   │   ├── product.model.go
│   │   │   ├── dtos/
│   │   │   │   ├── product.dto.go
│   │   │   ├── responses/
│   │   │   │   ├── product.response.go
│   │   │   ├── repo/
│   │   │   │   ├── product.repo.go
│   │   │   ├── services/
│   │   │   │   ├── product.service.go
│   │   │   ├── product.go
│   │   ├── db/
│   │   │   ├── di/
│   │   │   │   ├── db.di.go
│   │   │   ├── models/
│   │   │   │   ├── counter.model.go
│   │   │   ├── repo/
│   │   │   │   ├── main.repo.go
│   │   │   │   ├── sorting.repo.go
│   │   │   ├── sorting.go
│   ├── types/                   # Shared Types
│   │   ├── appcgf.go
│   │   ├── models.go
│   ├── util/                    # Utility functions
│   │   ├── context.go
│   │   ├── helper.go
│── api/                         # API documentation (Swagger, OpenAPI)
│── scripts/                     # Database migrations, seeders                  
│── go.mod                        
│── go.sum                        

## Localization

Use "github.com/goccy/go-json" for `MarshalJSON`

```golang

type LocalizableString map[string]string

func (ls LocalizableString) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string(ls))
}

func (ls *LocalizableString) UnmarshalJSON(data []byte) error {
	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}
	*ls = LocalizableString(translations)
	return nil
}

func (ls LocalizableString) Localize(languageKey string) string {
	if translation, exists := ls[languageKey]; exists {
		return strings.TrimSpace(translation)
	}
	// Default to English if the requested language is not available
	if translation, exists := ls["en"]; exists {
		return strings.TrimSpace(translation)
	}
	return ""
}

//example of use it

type Challenge struct {
	Title      LocalizableString  `bson:"title" json:"title"`
	Body       LocalizableString  `bson:"body" json:"body"`
	BaseDbModel
}

```

```golang

type BaseDbModel struct {
	CreatedAt time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time           `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CreatedBy *primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedBy *primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	DeletedAt time.Time           `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	DeletedBy *primitive.ObjectID `bson:"deletedBy,omitempty" json:"deletedBy,omitempty"`
}
```

### Tech Stack

- Go
- MongoDB
- Redis
- NATS
- Swagger
- github.com/goccy/go-json
- chi

### Constrains

- handle errors of 401 and 403 in auth middleware
- handle errors of 404 in not found middleware
- handle errors of 422 in validation middleware



### API Design

CRUD for <item>

- GET /<item>
- GET /<item>/all
- GET /<item>/:id
- POST /<item>
- PUT /<item>/:id
- DELETE /<item>/:id

use json.MarshalContext