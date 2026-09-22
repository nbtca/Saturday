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

This is a REST API server built in Go using the Huma v2 framework for the NBTCA repair service platform.

### Core Architecture
- **HTTP Layer**: `router/` - Huma v2-based HTTP handlers with Chi router and OpenAPI generation
- **Service Layer**: `service/` - Business logic and external service integrations  
- **Repository Layer**: `repo/` - Database access and queries using sqlx
- **Models**: `model/` - Data structures and DTOs
- **Utilities**: `util/` - Shared utilities, logging, validation, and helpers
- **Middleware**: `middleware/` - Authentication, logging, and request processing

### Key Dependencies
- **Web Framework**: Huma v2 (github.com/danielgtaylor/huma/v2) with Chi router
- **Database**: PostgreSQL with sqlx (github.com/jmoiron/sqlx) and Squirrel query builder
- **Migrations**: golang-migrate/migrate/v4
- **Configuration**: Viper with Consul support and hot-reloading
- **Authentication**: JWT with Logto integration and JWKS validation

### Database Architecture
- **Primary Database**: PostgreSQL (production)
- **Migrations**: Located in `migrations/` directory, run automatically on startup
- **Connection**: Uses sqlx with connection pooling, logging hooks, and Squirrel query builder
- **Test Data**: CSV files in `testdata/` directories for reproducible test scenarios

### External Service Integrations
- **Logto**: Authentication and user management
- **WeChat**: Mini-program API integration
- **GitHub**: Issue tracking and webhook handling
- **Aliyun OSS**: File storage
- **Dify**: AI service integration
- **NSQ**: Message queue for events and logging

### Configuration
- Supports environment variables, .env files, and Consul
- Configuration loaded via Viper with automatic environment mapping (`CONSUL_KEY` → `consul.key`)
- Hot-reloading capability for configuration changes from Consul
- Docker-friendly with override capabilities
- Example configuration provided in `config.example.json`

### Testing Strategy
- Unit tests in same packages as source code (`*_test.go`)
- Test data stored in `testdata/` directories as CSV files for reproducible scenarios

### API Design Patterns
- **OpenAPI-First**: Huma v2 automatically generates OpenAPI documentation from Go structs
- **Type-Safe**: Strong typing with validation at API boundaries using struct tags
- **Input Types**: Centralized input validation types in `router/input_types.go`
- **Middleware Pipeline**: Authentication, logging, CORS, and custom middleware integration
- **Error Handling**: Structured error responses with proper HTTP status codes