# Quick Start Guide

Get your Darween ERP platform up and running in minutes!

## Prerequisites

Before you begin, ensure you have:

- Go 1.24.2 or higher installed
- PostgreSQL 12 or higher installed and running
- Git (for version control)

## Step 1: Set Up Database

Create a PostgreSQL database for the application:

```bash
# Using createdb command
createdb erp_db

# Or using psql
psql -U postgres
CREATE DATABASE erp_db;
\q
```

Or use the Makefile:

```bash
make db-create
```

## Step 2: Configure Environment

Copy the example environment file and update with your settings:

```bash
cp env.example .env
```

Edit `.env` with your database credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=erp_db

JWT_SECRET=change-this-to-a-secure-random-string
JWT_EXPIRATION=24

SERVER_PORT=8080
GIN_MODE=debug
```

> **Important**: Change `JWT_SECRET` to a secure random string in production!

## Step 3: Install Dependencies

```bash
go mod download
```

Or using Makefile:

```bash
make deps
```

## Step 4: Run the Application

### Option 1: Using go run

```bash
go run cmd/api/main.go
```

### Option 2: Using Makefile

```bash
make run
```

### Option 3: Build and run

```bash
make build
./darween
```

The server will start on `http://localhost:8080` and automatically create the database tables on first run.

## Step 5: Test the API

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

You should receive a response with a JWT token:

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "first_name": "Admin",
      "last_name": "User"
    }
  }
}
```

### Create a company

Use the token from the registration response:

```bash
export TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "My First Company",
    "code": "MFC",
    "description": "This is my first company"
  }'
```

### Check your companies

```bash
curl -X GET http://localhost:8080/api/v1/companies \
  -H "Authorization: Bearer $TOKEN"
```

### Check company subscription

```bash
curl -X GET http://localhost:8080/api/v1/companies/1/subscription \
  -H "Authorization: Bearer $TOKEN"
```

## Makefile Commands

The project includes a Makefile with useful commands:

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make fmt           # Format code
make db-create     # Create database
make db-drop       # Drop database
make db-reset      # Reset database (drop and create)
```

## Troubleshooting

### Database Connection Error

If you get a database connection error:

1. Verify PostgreSQL is running: `pg_isready`
2. Check your database credentials in `.env`
3. Ensure the database exists: `psql -l | grep erp_db`

### Port Already in Use

If port 8080 is already in use:

1. Change `SERVER_PORT` in `.env` to another port (e.g., 8081)
2. Or stop the process using port 8080

### Migration Issues

If you encounter migration issues:

1. Drop and recreate the database:
   ```bash
   make db-reset
   ```
2. Restart the application

## Next Steps

Now that your ERP platform is running:

1. **Explore the API**: Check the full API documentation in [README.md](README.md)
2. **Add Users**: Invite team members to your company
3. **Manage Roles**: Assign different roles (Owner, Admin, Manager, Employee)
4. **Upgrade Subscription**: Test different subscription plans
5. **Build Features**: Start adding your custom ERP modules

## Development Tips

### Hot Reload (Optional)

For a better development experience with auto-reload on file changes:

```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Code Formatting

Always format your code before committing:

```bash
make fmt
```

### Testing

Run tests to ensure everything works:

```bash
make test
```

## Production Deployment

Before deploying to production:

1. ✅ Change `GIN_MODE` to `release` in `.env`
2. ✅ Use a strong, unique `JWT_SECRET`
3. ✅ Enable PostgreSQL SSL: Update DSN in `config.go`
4. ✅ Set up proper database backups
5. ✅ Use environment variables instead of `.env` file
6. ✅ Set up reverse proxy (nginx, traefik)
7. ✅ Enable HTTPS/TLS
8. ✅ Configure logging and monitoring
9. ✅ Review security best practices

## Support

For issues, questions, or contributions, please refer to the main [README.md](README.md) file.

---

**Happy coding! 🚀**


