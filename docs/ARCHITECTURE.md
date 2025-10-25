# Architecture Documentation

## Overview

Darween ERP follows **Domain-Driven Design (DDD)** principles with a clean, layered architecture. This document explains the architectural decisions and structure of the project.

## DDD Layers

### 1. Domain Layer (`internal/domain/`)

The **core business logic** layer containing pure domain models with no external dependencies.

#### Purpose
- Define business entities and their behavior
- Establish repository contracts (interfaces)
- Maintain business rules and domain logic
- Be framework and infrastructure agnostic

#### Structure
```
internal/domain/
├── user/
│   ├── entity.go        # User entity with business methods
│   └── repository.go    # User repository interface
├── company/
│   ├── entity.go        # Company entity
│   └── repository.go    # Company repository interface
└── subscription/
    ├── entity.go        # Subscription entity with plan logic
    └── repository.go    # Subscription repository interface
```

#### Key Concepts

**Entities**: Core business objects with identity
- `User`: User accounts with authentication
- `Company`: Business organizations
- `Subscription`: Company subscription plans
- `UserCompanyRole`: Many-to-many relationship between users and companies

**Repository Interfaces**: Define data access contracts without implementation details

**Example**:
```go
type Repository interface {
    Create(user *User) error
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
}
```

### 2. Application Layer (`internal/application/`)

Orchestrates domain logic and implements **use cases**.

#### Purpose
- Coordinate domain objects to perform business operations
- Define DTOs (Data Transfer Objects) for API contracts
- Handle business workflows
- Transform data between layers

#### Structure
```
internal/application/
├── auth/
│   ├── dto.go          # Login/Register request/response DTOs
│   └── service.go      # Authentication use cases
├── user/
│   ├── dto.go          # User DTOs
│   └── service.go      # User management use cases
├── company/
│   ├── dto.go          # Company DTOs
│   └── service.go      # Company management use cases
└── subscription/
    ├── dto.go          # Subscription DTOs
    └── service.go      # Subscription management use cases
```

#### Key Concepts

**Services**: Implement business use cases
- Single Responsibility: Each service handles one domain concept
- No HTTP knowledge: Services don't know about HTTP requests/responses
- Return domain errors, not HTTP errors

**DTOs**: Data transfer between layers
- Request DTOs: Validation rules via struct tags
- Response DTOs: Simplified data structures for API responses

**Example**:
```go
type Service struct {
    userRepo   user.Repository
    jwtManager *security.JWTManager
}

func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
    // Business logic here
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

Implements **technical concerns** and external integrations.

#### Purpose
- Implement repository interfaces with actual storage
- Provide security utilities (JWT, encryption)
- Handle database migrations
- Integrate with external services

#### Structure
```
internal/infrastructure/
├── persistence/
│   ├── postgres/
│   │   ├── database.go              # DB connection setup
│   │   ├── user_repository.go       # User repo implementation
│   │   ├── company_repository.go    # Company repo implementation
│   │   └── subscription_repository.go
│   └── migrations/
│       └── migrate.go               # Auto-migration
└── security/
    └── jwt.go                       # JWT token management
```

#### Key Concepts

**Repository Implementations**: Concrete implementations using GORM
- Implement domain repository interfaces
- Handle database-specific logic
- Use GORM for ORM operations

**Database Connection**: Centralized database setup
- Connection pooling
- Logging configuration
- Error handling

**Security**: Authentication and authorization utilities
- JWT token generation and validation
- Password hashing (bcrypt)

**Example**:
```go
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(u *user.User) error {
    return r.db.Create(u).Error
}
```

### 4. Presentation Layer (`internal/presentation/`)

Handles **HTTP communication** and user interface concerns.

#### Purpose
- Receive and parse HTTP requests
- Validate input data
- Call application services
- Format and return HTTP responses
- Handle authentication middleware

#### Structure
```
internal/presentation/
├── http/
│   ├── handler/
│   │   ├── auth_handler.go         # Auth endpoints
│   │   ├── user_handler.go         # User endpoints
│   │   ├── company_handler.go      # Company endpoints
│   │   └── subscription_handler.go # Subscription endpoints
│   ├── middleware/
│   │   ├── auth.go                 # JWT authentication
│   │   └── cors.go                 # CORS configuration
│   └── router/
│       └── router.go               # Route definitions
└── response/
    └── response.go                 # Standard response format
```

#### Key Concepts

**Handlers**: HTTP request/response handling
- Parse request body
- Validate input (using Gin's binding)
- Call application services
- Return formatted responses

**Middleware**: Cross-cutting concerns
- Authentication: JWT token validation
- CORS: Cross-origin resource sharing
- Logging: Request/response logging (future)
- Rate limiting (future)

**Router**: Centralized route configuration
- Group routes by feature
- Apply middleware to route groups
- RESTful URL structure

**Response Format**: Consistent API responses
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
}
```

### 5. Shared Packages (`pkg/`)

**Reusable utilities** that can be used across layers.

#### Structure
```
pkg/
├── config/
│   └── config.go       # Configuration management
└── errors/
    └── errors.go       # Custom error types
```

#### Key Concepts

**Config**: Environment-based configuration
- Load from `.env` file
- Environment variable support
- Type-safe configuration structure
- Validation

**Errors**: Domain-specific error types
- Typed errors for different scenarios
- Error codes for API responses
- Error wrapping support

## Dependency Flow

The dependency rule: **Dependencies point inward**

```
Presentation → Application → Domain
      ↓             ↓
Infrastructure ←────┘
```

### Rules:
1. **Domain** has no dependencies on other layers
2. **Application** depends only on Domain
3. **Infrastructure** implements Domain interfaces
4. **Presentation** depends on Application and uses Infrastructure

## Data Flow

### Request Flow (Inbound)
```
HTTP Request
    ↓
Router
    ↓
Middleware (Auth)
    ↓
Handler
    ↓
Application Service
    ↓
Domain Entity/Repository
    ↓
Infrastructure Repository (DB)
```

### Response Flow (Outbound)
```
Database Result
    ↓
Infrastructure Repository
    ↓
Domain Entity
    ↓
Application Service (DTO)
    ↓
Handler
    ↓
Response Formatter
    ↓
HTTP Response
```

## Design Patterns

### Repository Pattern
- **Interface**: Defined in Domain layer
- **Implementation**: In Infrastructure layer
- **Purpose**: Abstract data access, enable testing

### Dependency Injection
- **Manual DI**: Constructor injection in `main.go`
- **Benefits**: Loose coupling, easier testing, clear dependencies

### DTO Pattern
- **Purpose**: Decouple API contracts from domain models
- **Location**: Application layer
- **Benefits**: API versioning, validation, transformation

### Middleware Pattern
- **Purpose**: Cross-cutting concerns (auth, logging, CORS)
- **Location**: Presentation layer
- **Benefits**: Reusable, composable, separation of concerns

## Security Architecture

### Authentication Flow
1. User provides credentials (email/password)
2. Service validates credentials
3. JWT token generated with user claims
4. Token returned to client
5. Client includes token in subsequent requests
6. Middleware validates token on protected routes

### Authorization
- **Role-based**: Owner, Admin, Manager, Employee
- **Company-scoped**: Users have different roles in different companies
- **Enforcement**: In application services, not at database level

## Database Architecture

### Multi-Tenancy
- **Strategy**: Shared database, row-level isolation
- **Implementation**: `company_id` in relevant tables
- **Benefits**: Simpler architecture, cost-effective

### Key Relationships
```
Users ←→ UserCompanyRoles ←→ Companies
                              ↓
                         Subscriptions
```

### Indexes
- Primary keys: Automatically indexed
- Foreign keys: Indexed for join performance
- Unique constraints: Email, company code
- Composite indexes: (user_id, company_id) for fast lookups

## Testing Strategy (Future)

### Unit Tests
- **Domain**: Test business logic in entities
- **Application**: Test use cases with mocked repositories
- **Infrastructure**: Test repository implementations with test database

### Integration Tests
- **API**: Test HTTP endpoints end-to-end
- **Database**: Test repository operations with real database

### Test Structure
```
internal/domain/user/entity_test.go
internal/application/auth/service_test.go
internal/infrastructure/postgres/user_repository_test.go
internal/presentation/http/handler/auth_handler_test.go
```

## Scalability Considerations

### Horizontal Scaling
- Stateless design (JWT tokens)
- Database connection pooling
- No server-side sessions

### Performance
- Indexed queries
- Connection pooling (10 idle, 100 max open)
- Lazy loading with GORM

### Future Enhancements
- Caching layer (Redis)
- Event-driven architecture (message queue)
- CQRS pattern for complex queries
- Microservices split (if needed)

## Code Organization Principles

1. **Separation of Concerns**: Each layer has distinct responsibilities
2. **Single Responsibility**: Each file/struct does one thing well
3. **Dependency Inversion**: Depend on interfaces, not implementations
4. **Open/Closed**: Open for extension, closed for modification
5. **Don't Repeat Yourself**: Shared utilities in `pkg/`

## Adding New Features

### Example: Adding "Projects" Feature

1. **Domain Layer** (`internal/domain/project/`)
   ```go
   // entity.go - Define Project entity
   // repository.go - Define repository interface
   ```

2. **Infrastructure Layer** (`internal/infrastructure/persistence/postgres/`)
   ```go
   // project_repository.go - Implement repository
   ```

3. **Application Layer** (`internal/application/project/`)
   ```go
   // dto.go - Define request/response DTOs
   // service.go - Implement project use cases
   ```

4. **Presentation Layer** (`internal/presentation/http/handler/`)
   ```go
   // project_handler.go - Implement HTTP handlers
   ```

5. **Router** (`internal/presentation/http/router/router.go`)
   ```go
   // Add project routes
   ```

6. **Main** (`cmd/api/main.go`)
   ```go
   // Wire up dependencies
   projectRepo := postgres.NewProjectRepository(db)
   projectService := project.NewService(projectRepo)
   projectHandler := handler.NewProjectHandler(projectService)
   ```

## Best Practices

1. **Keep Domain Pure**: No external dependencies in domain layer
2. **Use Interfaces**: Define contracts, not implementations
3. **Return Errors**: Don't panic, return errors
4. **Validate Early**: Use DTOs with validation tags
5. **Log Strategically**: Log at boundaries (HTTP, DB)
6. **Test Thoroughly**: Unit tests for logic, integration tests for flows
7. **Document APIs**: Use comments and API documentation tools

## References

- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Web Framework](https://gin-gonic.com/docs/)


