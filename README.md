# Darween ERP Platform

A modern ERP platform built with Domain-Driven Design (DDD) architecture using Go, Gin, GORM, and PostgreSQL.

## Features

- **User Management**: Complete user authentication and authorization
- **Multi-Company Support**: Users can belong to multiple companies with different roles
- **Subscription Management**: Flexible subscription plans (Free, Basic, Premium, Enterprise)
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Owner, Admin, Manager, and Employee roles

## Architecture

This project follows Domain-Driven Design (DDD) principles with a clean, layered architecture:

```
darween/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain layer (entities, interfaces)
│   ├── application/            # Application layer (use cases, DTOs)
│   ├── infrastructure/         # Infrastructure layer (database, security)
│   └── presentation/           # Presentation layer (HTTP handlers, middleware)
└── pkg/                        # Shared packages (config, errors)
```

## Technology Stack

- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator

## Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 12 or higher

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd darween
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   createdb erp_db
   ```

4. **Configure environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=erp_db

   # JWT Configuration
   JWT_SECRET=your-secret-key-change-in-production
   JWT_EXPIRATION=24

   # Server Configuration
   SERVER_PORT=8080
   GIN_MODE=debug
   ```

5. **Run the application**
   ```bash
   # Using go run
   go run cmd/api/main.go
   
   # Or using Makefile
   make run
   
   # Or build first
   make build
   ./darween
   ```

   The server will start on `http://localhost:8080`

## Makefile Commands

The project includes a Makefile for common tasks:

```bash
make build      # Build the application
make run        # Run the application
make clean      # Clean build artifacts
make test       # Run tests
make tidy       # Tidy Go modules
make deps       # Download dependencies
make fmt        # Format code
make db-create  # Create PostgreSQL database
make db-drop    # Drop PostgreSQL database
make db-reset   # Reset database (drop and recreate)
make help       # Show all available commands
```

## API Endpoints

### Authentication (Public)

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

### Users (Protected)

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `GET /api/v1/users?company_id=<id>` - List users in a company

### Companies (Protected)

- `POST /api/v1/companies` - Create a new company
- `GET /api/v1/companies` - List user's companies
- `GET /api/v1/companies/:id` - Get company details
- `PUT /api/v1/companies/:id` - Update company
- `POST /api/v1/companies/:id/users` - Add user to company
- `DELETE /api/v1/companies/:id/users/:userId` - Remove user from company

### Subscriptions (Protected)

- `GET /api/v1/companies/:id/subscription` - Get company subscription
- `PUT /api/v1/companies/:id/subscription` - Update subscription

### Health Check

- `GET /api/v1/health` - API health check

## API Usage Examples

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    }
  }
}
```

### 3. Create a Company

```bash
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "My Company",
    "code": "MYCO",
    "description": "My company description"
  }'
```

### 4. Get User's Companies

```bash
curl -X GET http://localhost:8080/api/v1/companies \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Add User to Company

```bash
curl -X POST http://localhost:8080/api/v1/companies/1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "email": "newuser@example.com",
    "role": "employee"
  }'
```

## User Roles

- **Owner**: Full control over company (cannot be removed)
- **Admin**: Can manage company settings and users
- **Manager**: Can manage operational tasks
- **Employee**: Basic access to company resources

## Subscription Plans

- **Free**: Up to 5 users
- **Basic**: Up to 20 users
- **Premium**: Up to 100 users
- **Enterprise**: Up to 1000 users

## Database Schema

### Users
- User profiles and credentials
- Many-to-many relationship with companies

### Companies
- Company information
- One-to-many relationship with subscriptions

### Subscriptions
- Plan type, status, dates, and user limits
- Belongs to one company

### UserCompanyRoles
- Junction table for user-company relationships
- Stores user roles within companies

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o darween cmd/api/main.go
```

### Running the Binary

```bash
./darween
```

## Project Structure Details

### Domain Layer
Contains pure business logic with no external dependencies:
- Entities: Core business objects
- Repository Interfaces: Data access contracts

### Application Layer
Orchestrates domain logic and implements use cases:
- Services: Business logic coordination
- DTOs: Data transfer objects for API contracts

### Infrastructure Layer
Implements technical concerns:
- Database: GORM repositories
- Security: JWT token management
- Migrations: Database schema management

### Presentation Layer
Handles HTTP communication:
- Handlers: HTTP request/response handling
- Middleware: Authentication, CORS
- Router: Route configuration

## Error Handling

The application uses custom error types with appropriate HTTP status codes:

- `NOT_FOUND` (404): Resource not found
- `UNAUTHORIZED` (401): Authentication required
- `FORBIDDEN` (403): Insufficient permissions
- `VALIDATION_ERROR` (400): Invalid input
- `CONFLICT` (409): Resource conflict
- `INTERNAL_ERROR` (500): Server error

## Security

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours (configurable)
- CORS enabled for cross-origin requests
- SQL injection prevention via GORM parameterization

## License

[Your License Here]

## Contributing

[Your Contributing Guidelines Here]

