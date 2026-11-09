# larsa-e-commerce-v2-microservice

E-commerce microservice for the Larsa platform.

**Stack:** Go, MongoDB, Redis, NATS

> Built by the Larsa engineering team and published with the company's permission. My contribution: 45 of 2527 commits. This is a clean snapshot: git history is removed and secrets are replaced with placeholders.

## Setup

```bash
cp .env.example .env   # then fill in the values
```

## Getting Started

### Initializing the Workspace

1. Clone the repository:

```bash
git clone https://git.larsa.io/larsa/larsa-e-commerce-v2-microservice.git
cd larsa-e-commerce-v2-microservice
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your development environment:

- Ensure you have Go 1.20+ installed
- Configure your IDE with Go tools
- Set up your environment variables

## Project Structure and Principles

### Directory Structure

```
module-name/
├── dtos/           # Data Transfer Objects
├── enums/          # Enumerations
├── filters/        # Filter definitions
├── handlers/       # HTTP handlers
├── models/         # Data models
├── repos/          # Repository implementations
├── services/       # Business logic
└── injector.go     # Dependency injection
```

### Naming Conventions

#### Package Naming

- Use the form: `moduleName_(models, dtos, handler, ...)`
- Example: `category_models`, `product_dtos`

#### File Naming

- Use hyphens between words for folder names (e.g., `plan-features`)
- Use hyphens between words for file names
- File names should reflect their purpose:
  - Models: `modelName.model.go`
  - Repositories: `repoName.repo.go`
  - DTOs: `dtoName.dto.go`
  - Handlers: `handlerName.handler.go`
  - Services: `serviceName.service.go`

#### Code Naming

- Use camelCase for BSON and JSON fields
- Example: `bson:"isFilterable" json:"isFilterable"`

### Core Implementation Patterns

#### Model Creation

1. Create the model file (e.g., `feature.model.go`)
2. Create corresponding DTOs:
   - Create DTO (`CreateFeatureDTO`)
   - Update DTO (`UpdateFeatureDTO`)
3. Implement converter functions in the service layer

Example converter function:

```go
func convertCreateDTOToModel(dto *feature_dto.CreateFeatureDTO, ctx context.Context) (feature_models.FeatureModel, *helpers.ApiError) {
    now := time.Now()
    cfg, err := util.GetReqAppCfg(ctx)
    if err != nil {
        return feature_models.FeatureModel{}, &helpers.ApiError{
            StatusCode: 500,
            Message:    "GetReqAppCfg went wrong",
        }
    }
    userId := cfg.User.Id
    // ... additional conversion logic
}
```

#### Repository Pattern

1. Create repository interface in the service package
2. Implement the repository in the repos package
3. Use dependency injection through the injector

#### Base Models and Repositories

- Always extend base models for common fields
- Implement base repository interfaces
- Follow consistent error handling patterns

### Base Models and Repositories

#### Base Models

Base models provide common fields and functionality that should be present in all models. Here's how to use them:

1. **Extending Base Model**

```go
type BaseModel struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    CreatedAt time.Time         `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time         `bson:"updatedAt" json:"updatedAt"`
    CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
    UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
    IsDeleted bool              `bson:"isDeleted" json:"isDeleted"`
}

type FeatureModel struct {
    BaseModel          `bson:",inline"`
    Name        string `bson:"name" json:"name"`
    Description string `bson:"description" json:"description"`
    // ... other fields
}
```

2. **Using Base Model Fields**

- Always set base model fields in your converter functions:

```go

func convertCreateDTOToModel(dto *feature_dto.CreateFeatureDTO, ctx context.Context) (feature_models.FeatureModel, *helpers.ApiError) {
    now := time.Now()
    cfg, err := util.GetReqAppCfg(ctx)
    if err != nil {
        return feature_models.FeatureModel{}, &helpers.ApiError{
            StatusCode: 500,
            Message:    "GetReqAppCfg went wrong",
        }
    }

    return feature_models.FeatureModel{
        BaseModel: BaseModel{
            CreatedAt: now,
            UpdatedAt: now,
            CreatedBy: cfg.User.Id,
            UpdatedBy: cfg.User.Id,
            IsDeleted: false,
        },
        Name:        dto.Name,
        Description: dto.Description,
    }, nil
}
```

#### Base Repository

The base repository provides common CRUD operations and should be implemented for all repositories:

1. **Implementing Base Repository**

```go

type FeatureRepository struct {
	*repo.BaseRepo[feature_models.FeatureModel, feature_models.FeatureModel]
}

func NewFeatureRepository(i *do.Injector) (feature_service.FeatureRepo, error) {
	db := do.MustInvoke[*mongo.Client](i)
	return &FeatureRepository{
		BaseRepo: repo.NewBaseRepo[feature_models.FeatureModel, feature_models.FeatureModel](db, string(enums.FEATURE)),
	}, nil
}


```

4. **Best Practices for Base Repository Usage**

- Always implement all base repository methods
- Handle errors consistently
- Use proper context for database operations
- Implement proper validation before database operations
- Use transactions when needed for multiple operations
- Implement proper indexing for better performance
- Use proper pagination for list operations

5. **Common Base Repository Operations**

- Soft delete instead of hard delete
- Track creation and update timestamps
- Track user who created/updated records
- Implement proper filtering and sorting
- Use proper error handling patterns

### Adding New Modules to Enums

When creating a new module, you need to add its configuration to two enum files:

1. **Add to `model.enums.go`**
   Add your new collection name to the `CollectionName` enum:

```go
const (
    // ... existing collections ...
    NEW_MODULE CollectionName = "NEW_MODULE"
)

func (i CollectionName) String() string {
    switch i {
    // ... existing cases ...
    case NEW_MODULE:
        return "ecommerceV2NewModules"
    default:
        return "ecommerceV2Default"
    }
}
```

2. **Add to `capabilities.go`**
   Add your new module to the capabilities configuration:

```go
const (
    // ... existing modules ...
    NewModuleModule ModuleType = "newModule"
)

var capabilityPrefixes = map[ModuleType]string{
    // ... existing prefixes ...
    NewModuleModule: "ecommerceV2",
}

type CapabilityRegistry struct {
    // ... existing capabilities ...
    NewModule ModuleCapability
}

func InitCapabilities() {
    CapabilitiesRegistry = CapabilityRegistry{
        // ... existing registry entries ...
        NewModule: initModuleCapability(NewModuleModule),
    }
}
```

3. **Usage in Repository**
   When creating your repository, use the enum for the collection name:

```go
func NewNewModuleRepository(i *do.Injector) (newmodule_service.NewModuleRepo, error) {
    db := do.MustInvoke[*mongo.Client](i)
    return &NewModuleRepository{
        BaseRepo: repo.NewBaseRepo[newmodule_models.NewModuleModel, newmodule_models.NewModuleModel](
            db,
            string(enums.NEW_MODULE),
        ),
    }, nil
}
```

4. **Best Practices for Enum Updates**

- Always add new modules to both enum files
- Use consistent naming conventions
- Follow the existing pattern for collection names
- Ensure the module name matches across all files
- Update the capabilities initialization
- Test the new module's capabilities

### Development Guidelines

#### Creating a New Module

1. Create the following folders:
   - `dtos/`
   - `enums/`
   - `filters/`
   - `handlers/`
   - `models/`
   - `repos/`
   - `services/`
2. Create an `injector.go` file
3. Follow the established patterns for each component

#### DTO Validation

- Always validate DTOs for create and update operations
- Use appropriate validation tags
- Implement custom validation when needed

#### Error Handling

- Use consistent error handling patterns
- Return appropriate HTTP status codes
- Provide meaningful error messages

#### Code Organization

- Maintain separation of concerns
- Keep business logic in services
- Handle data access in repositories
- Manage HTTP concerns in handlers

### Best Practices

1. **Validation**

   - Always validate input data
   - Use appropriate validation tags
   - Implement custom validation when needed

2. **Error Handling**

   - Use consistent error patterns
   - Provide meaningful error messages
   - Log errors appropriately

3. **Code Structure**

   - Follow the established directory structure
   - Maintain clean separation of concerns
   - Keep functions focused and small

4. **Documentation**
   - Document complex logic
   - Add comments for non-obvious code
   - Keep documentation up to date
