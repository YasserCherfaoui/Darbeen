<!-- a4575041-5921-4563-b061-326c2f5f6534 7f9eb968-c0c7-4501-a50a-bc2589270203 -->
# ERP Platform with DDD Architecture

## Project Structure

Following Domain-Driven Design principles, the project will be organized as:

```
darween/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── domain/                        # Domain layer (entities, interfaces)
│   │   ├── user/
│   │   │   ├── entity.go             # User entity
│   │   │   └── repository.go         # User repository interface
│   │   ├── company/
│   │   │   ├── entity.go             # Company entity
│   │   │   └── repository.go         # Company repository interface
│   │   └── subscription/
│   │       ├── entity.go             # Subscription entity
│   │       └── repository.go         # Subscription repository interface
│   ├── application/                   # Application layer (use cases)
│   │   ├── user/
│   │   │   ├── dto.go                # DTOs for user operations
│   │   │   └── service.go            # User application service
│   │   ├── company/
│   │   │   ├── dto.go
│   │   │   └── service.go
│   │   ├── subscription/
│   │   │   ├── dto.go
│   │   │   └── service.go
│   │   └── auth/
│   │       ├── dto.go
│   │       └── service.go            # Authentication service (JWT)
│   ├── infrastructure/                # Infrastructure layer
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── database.go       # DB connection
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── company_repository.go
│   │   │   │   └── subscription_repository.go
│   │   │   └── migrations/
│   │   │       └── migrate.go        # Auto-migration setup
│   │   └── security/
│   │       └── jwt.go                # JWT token generation/validation
│   └── presentation/                  # Presentation layer (HTTP)
│       ├── http/
│       │   ├── handler/
│       │   │   ├── user_handler.go
│       │   │   ├── company_handler.go
│       │   │   ├── subscription_handler.go
│       │   │   └── auth_handler.go
│       │   ├── middleware/
│       │   │   ├── auth.go           # JWT authentication middleware
│       │   │   └── cors.go
│       │   └── router/
│       │       └── router.go         # Gin router setup
│       └── response/
│           └── response.go           # Standard HTTP response format
├── pkg/                               # Shared packages
│   ├── config/
│   │   └── config.go                 # Configuration management
│   └── errors/
│       └── errors.go                 # Custom error types
├── go.mod
├── go.sum
└── .env.example                      # Environment variables template
```

## Core Domain Models

### User Entity

- Fields: ID, Email, Password (hashed), FirstName, LastName, IsActive, CreatedAt, UpdatedAt
- Relationships: Many-to-many with Company through UserCompanyRole

### Company Entity

- Fields: ID, Name, Code (unique identifier), Description, IsActive, CreatedAt, UpdatedAt
- Relationships: Many users, one subscription

### Subscription Entity

- Fields: ID, CompanyID, PlanType (enum: free, basic, premium, enterprise), Status (enum: active, inactive, expired), StartDate, EndDate, MaxUsers, CreatedAt, UpdatedAt
- Relationships: Belongs to Company

### UserCompanyRole Entity (Junction Table)

- Fields: ID, UserID, CompanyID, Role (enum: owner, admin, manager, employee), IsActive, CreatedAt
- Purpose: Enables users to belong to multiple companies with different roles

## Key Implementation Details

### 1. Database Setup

- Use GORM with PostgreSQL driver
- Implement auto-migration for development
- Add unique constraints: User.Email, Company.Code
- Add indexes on foreign keys and frequently queried fields

### 2. JWT Authentication

- Token generation on login with user ID and active company roles
- Middleware to validate JWT and extract user context
- Token expiration: 24 hours (configurable)
- Refresh token mechanism (optional for future)

### 3. Repository Pattern

- Define interfaces in domain layer
- Implement concrete repositories in infrastructure/persistence
- Handle multi-company data isolation with company_id filtering where applicable

### 4. HTTP Endpoints

**Auth Routes** (public):

- POST `/api/v1/auth/register` - User registration
- POST `/api/v1/auth/login` - User login (returns JWT)

**User Routes** (protected):

- GET `/api/v1/users/me` - Get current user profile
- PUT `/api/v1/users/me` - Update current user profile
- GET `/api/v1/users` - List users (filtered by company context)

**Company Routes** (protected):

- POST `/api/v1/companies` - Create company (user becomes owner)
- GET `/api/v1/companies` - List user's companies
- GET `/api/v1/companies/:id` - Get company details
- PUT `/api/v1/companies/:id` - Update company (admin/owner only)
- POST `/api/v1/companies/:id/users` - Add user to company
- DELETE `/api/v1/companies/:id/users/:userId` - Remove user from company

**Subscription Routes** (protected):

- GET `/api/v1/companies/:id/subscription` - Get company subscription
- PUT `/api/v1/companies/:id/subscription` - Update subscription (owner only)

### 5. Configuration

Environment variables:

- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`, `JWT_EXPIRATION`
- `SERVER_PORT`
- `GIN_MODE` (debug/release)

### 6. Error Handling

- Custom error types for domain errors (NotFound, Unauthorized, ValidationError)
- Centralized error response formatting
- Proper HTTP status codes

### 7. Dependencies

Required packages:

- `github.com/gin-gonic/gin` - HTTP framework
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/joho/godotenv` - Environment variables
- `github.com/go-playground/validator/v10` - Request validation

### To-dos

- [ ] Create DDD folder structure and initialize Go modules with dependencies
- [ ] Implement configuration management and environment loading
- [ ] Create domain entities (User, Company, Subscription, UserCompanyRole) and repository interfaces
- [ ] Implement database connection, GORM setup, and auto-migration
- [ ] Implement repository pattern with PostgreSQL for all entities
- [ ] Implement JWT token generation and validation utilities
- [ ] Create authentication application service (register, login)
- [ ] Implement application services for User, Company, and Subscription
- [ ] Create HTTP handlers and middleware (auth, CORS)
- [ ] Set up Gin router with routes and create main.go entry point
- [ ] Create .env.example file with all required environment variables