.PHONY: build run clean test tidy migrate help

# Build the application
build:
	@echo "Building darween..."
	@go build -o darween cmd/api/main.go
	@echo "Build complete!"

# Run the application
run:
	@echo "Running darween..."
	@go run cmd/api/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f darween
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Tidy complete!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Dependencies downloaded!"

# Run with hot reload (requires air)
dev:
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

# Lint code (requires golangci-lint)
lint:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	@golangci-lint run

# Create database
db-create:
	@echo "Creating database..."
	@createdb erp_db || echo "Database may already exist"

# Drop database
db-drop:
	@echo "Dropping database..."
	@dropdb erp_db

# Reset database
db-reset: db-drop db-create
	@echo "Database reset complete!"

# Help
help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make test      - Run tests"
	@echo "  make tidy      - Tidy Go modules"
	@echo "  make deps      - Download dependencies"
	@echo "  make dev       - Run with hot reload (requires air)"
	@echo "  make fmt       - Format code"
	@echo "  make lint      - Lint code (requires golangci-lint)"
	@echo "  make db-create - Create database"
	@echo "  make db-drop   - Drop database"
	@echo "  make db-reset  - Reset database"
	@echo "  make help      - Show this help message"


