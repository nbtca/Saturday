# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build and Run
```bash
# Build the application
go build -o saturday .

# Run the application directly
go run main.go

# Run with Docker
docker build -t saturday .
docker-compose up
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output and coverage
go test -v -cover ./...

# Run specific package tests
go test ./router
go test ./service
go test ./middleware

# Run benchmarks
go test -bench=. ./...
```

### Database Operations
- Database migrations run automatically on application startup
- Manual migration commands (if needed):
```bash
# Run migrations up
migrate -path migrations -database $DB_DATASOURCE up

# Rollback migrations
migrate -path migrations -database $DB_DATASOURCE down
```

### Code Quality
```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Clean up dependencies
go mod tidy
```

## Architecture Overview

This is a REST API server built in Go using the Gin web framework for the NBTCA repair service platform.

### Core Architecture
- **HTTP Layer**: `router/` - Gin-based HTTP handlers and middleware
- **Service Layer**: `service/` - Business logic and external service integrations  
- **Repository Layer**: `repo/` - Database access and queries using sqlx
- **Models**: `model/` - Data structures and DTOs
- **Utilities**: `util/` - Shared utilities, logging, validation, and helpers
- **Middleware**: `middleware/` - Authentication, logging, and request processing

### Key Dependencies
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: PostgreSQL with sqlx (github.com/jmoiron/sqlx)
- **Migrations**: golang-migrate/migrate/v4
- **Configuration**: Viper with Consul support
- **Testing**: Testify + Dockertest for integration tests
- **Authentication**: JWT with Logto integration

### Database Architecture
- **Primary Database**: PostgreSQL (production)
- **Test Database**: MySQL 8.0 (via Docker containers)
- **Migrations**: Located in `migrations/` directory, run automatically on startup
- **Connection**: Uses sqlx with connection pooling and logging hooks

### External Service Integrations
- **Logto**: Authentication and user management
- **WeChat**: Mini-program API integration
- **GitHub**: Issue tracking and webhook handling
- **Aliyun OSS**: File storage
- **Dify**: AI service integration
- **NSQ**: Message queue for events and logging

### Configuration
- Supports environment variables, .env files, and Consul
- Configuration loaded via Viper with automatic environment mapping
- Docker-friendly with override capabilities

### Testing Strategy
- Unit tests in same packages as source code (`*_test.go`)
- Integration tests use Docker containers for database testing
- Test data stored in `testdata/` directories as CSV files
- Comprehensive coverage of HTTP endpoints, services, and middleware