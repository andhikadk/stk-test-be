# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go Fiber Boilerplate is a production-ready REST API template using the Fiber framework (Express.js-like for Go), GORM ORM, and PostgreSQL/SQLite. It includes JWT authentication, database migrations, Docker support, and hot reload for development.

## Common Development Commands

### Building & Running
```bash
make build              # Build binary to ./bin/go-fiber-boilerplate
make run               # Run application immediately
make dev               # Run with hot reload (requires air installed)
```

### Testing
```bash
make test              # Run all tests with verbose output
make test-coverage     # Generate HTML coverage report (coverage.html)
```

### Code Quality
```bash
make fmt               # Format code with go fmt
make vet               # Run go vet analysis
make lint              # Run golangci-lint (install via: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
```

### Database
```bash
make migrate           # Run migrations (AutoMigrate in dev, SQL in prod)
make migrate-sql       # Run SQL migrations from embedded files
make migrate-status    # Show applied migrations and seeds
make seed              # Seed database with sample data
```

### Docker Development (Recommended)
```bash
make docker-dev        # Start with hot reload (best for development)
make docker-dev-logs   # View logs
make docker-dev-down   # Stop containers
make docker-dev-reset  # Reset containers and database
```

### Docker Production
```bash
make docker-up         # Start production containers
make docker-down       # Stop containers
make docker-logs       # View logs
make docker-reset      # Reset containers and database
```

### Utilities
```bash
make install-deps      # Download and tidy dependencies
make clean             # Remove build artifacts
make help              # Show all available commands
make all               # Clean, install, build, and test
```

## Project Architecture

### High-Level Structure

The application follows a **layered architecture** with clear separation of concerns:

```
main.go (entry point)
  ├── config/ (configuration management)
  ├── internal/ (core application logic)
  │   ├── handlers/ (HTTP request handlers)
  │   ├── services/ (business logic layer)
  │   ├── models/ (data structures)
  │   ├── middleware/ (request interceptors)
  │   ├── database/ (DB initialization & migrations)
  │   └── routes/ (route definitions)
  ├── pkg/ (reusable utilities)
  └── migrations/ (SQL files)
```

### Key Architectural Patterns

#### 1. **Handlers (HTTP Layer)**
- Located in `internal/handlers/`
- Convert HTTP requests to domain operations
- Use `fiber.Ctx` for request/response handling
- Call services for business logic
- Return standardized responses via utility functions

Example flow: `Register` handler → calls `AuthService.Register()` → returns JSON via `utils.CreatedResponse()`

#### 2. **Services (Business Logic Layer)**
- Located in `internal/services/`
- Contain all business logic and validation
- Access database through GORM
- Services are instantiated fresh per request (e.g., `NewAuthService()`)
- Services hold a reference to the database: `type AuthService struct { db *gorm.DB }`

#### 3. **Models (Domain & Request/Response)**
- Located in `internal/models/`
- Include both ORM models (e.g., `User`, `Book`) and DTOs (request/response structs)
- ORM models have GORM tags for database mapping
- Response models (e.g., `LoginResponse`) define API contract

#### 4. **Middleware Chain**
- Located in `internal/middleware/`
- Auth middleware extracts and validates JWT tokens
- Error handling middleware catches panics and formats errors
- Set in `main.go:setupMiddleware()` for global middleware
- Route-specific middleware applied via `group.Use(middleware.AuthMiddleware())`

#### 5. **Database Layer**
- Located in `internal/database/`
- `db.go`: Connection and migration orchestration
- `migrator.go`: SQL migration runner from embedded files
- `seeder.go`: Sample data seeding
- Uses `embed.FS` to bundle migration files into binary for production

#### 6. **Configuration Management**
- Located in `config/`
- `config.go`: Load from `.env` using `godotenv`, validate, expose via `config.AppConfig`
- Database connection created in `database.go` based on `DBDriver` (postgres/sqlite)
- JWT secret, timeouts, CORS settings all configured here

#### 7. **JWT Authentication**
- Located in `pkg/jwt/`
- Token generation and validation
- Access tokens (short-lived, 15min default) and refresh tokens (7d default)
- Claims include `UserID`, `Email`, `Role`
- Auth middleware validates token and extracts user info into context

### Request/Response Flow

1. **Request arrives** → Fiber routes to handler in `internal/handlers/`
2. **Handler parses** request body into DTO (e.g., `RegisterRequest`)
3. **Handler validates** input and calls service method
4. **Service executes** business logic: database queries, validations, token generation
5. **Service returns** domain objects or errors
6. **Handler formats** response using utility functions (`utils.SuccessResponse()`, `utils.CreatedResponse()`, etc.)
7. **Response sent** as JSON with standardized structure

### Database Migrations

The project uses a **dual-migration system**:

- **Development**: `AutoMigrate` in `internal/database/db.go` for fast iteration
- **Production**: SQL files in `migrations/` directory embedded via `embed.go`

Run migrations via:
```bash
go run main.go -migrate        # AutoMigrate (dev)
go run main.go -migrate=sql    # SQL migrations (prod)
```

Migrations track applied status in `schema_migrations` table.

### Configuration & Environment

- Load configuration in `main.go` via `config.LoadConfig()`
- Read from `.env` file (or system env vars)
- Copy `.env.example` to `.env` and customize
- Critical settings: `PORT`, `DB_DRIVER`, `DB_HOST`, `JWT_SECRET`, `DB_USER`, `DB_PASSWORD`
- Different port mapping for Docker: `5432` internal → `6543` host-accessible

## Adding New Features

### Adding a New Entity (e.g., "Product")

1. **Create Model** in `internal/models/product.go`
   - Define GORM model struct with tags
   - Optionally add helper methods (e.g., `GetPublicProduct()`)

2. **Create Service** in `internal/services/product_service.go`
   - Implement business logic (CRUD operations)
   - Service receives database via `database.GetDB()`

3. **Create Handler** in `internal/handlers/product.go`
   - Implement HTTP endpoints
   - Parse requests, call service, format responses

4. **Add Routes** in `internal/routes/routes.go`
   - Define new route groups (e.g., `/api/products`)
   - Add auth middleware if needed

5. **Create Migration** in `migrations/` (optional for SQL migrations)
   - Or rely on `AutoMigrate` in development

6. **Write Tests** in `tests/` directory
   - Test service logic and handlers

### Modifying Existing Features

- Keep changes isolated to relevant layers (model → service → handler)
- Update migrations if schema changes
- Add tests for behavioral changes
- Run `make test-coverage` to ensure coverage

## Testing

The project uses `testify` and Go's standard `testing` package.

- Run all tests: `make test`
- Run specific test: `go test -v ./tests -run TestName`
- Generate coverage: `make test-coverage`
- Coverage report opens as HTML

## Docker Development Workflow

**Recommended setup for active development:**

```bash
make docker-dev              # Terminal 1: Start containers with hot reload
make docker-dev-logs         # Terminal 2: Monitor logs
# Edit code in your editor
# Air automatically rebuilds on file save
# Refresh browser to see changes
```

Benefits:
- Database runs in container (no local PostgreSQL needed)
- Hot reload: changes instantly visible (no rebuild)
- Production-like environment
- Easy to reset with `make docker-dev-reset`

**For production testing:**
```bash
make docker-up               # Build production image, run compiled binary
```

## Important Files & Their Purposes

| File | Purpose |
|------|---------|
| `main.go` | Entry point; sets up middleware, routes, migrations |
| `embed.go` | Embeds migration files into binary |
| `config/config.go` | Load and validate environment configuration |
| `config/database.go` | GORM connection setup based on driver |
| `internal/routes/routes.go` | Route definitions and grouping |
| `internal/middleware/auth.go` | JWT validation and user context extraction |
| `internal/middleware/error.go` | Global error handling |
| `pkg/jwt/manager.go` | JWT token creation and validation |
| `pkg/utils/responses.go` | Standard response formatting |
| `.env.example` | Template for environment variables |
| `.air.toml` | Hot reload configuration |

## Key Dependencies

- **Fiber v2**: Web framework
- **GORM**: ORM for database operations
- **golang-jwt**: JWT token handling
- **bcrypt**: Password hashing
- **validator/v10**: Input validation
- **godotenv**: .env file loading

## Code Organization Principles

1. **Single Responsibility**: Each file/function has one reason to change
2. **Dependency Injection**: Services receive database, don't create it
3. **Error Handling**: Return errors explicitly, use middleware for global handling
4. **Validation**: Validate input early in handlers; return 400 for bad requests
5. **No Sensitive Data in Logs**: Password hashes are excluded from JSON response
6. **Standard Response Format**: All endpoints return consistent JSON structure

## Concurrent Programming Patterns

This boilerplate includes educational examples of common Go concurrency patterns in `internal/services/concurrent_service.go` and `internal/handlers/concurrent.go`.

### Available Patterns

Access the overview endpoint to see all available patterns:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:3000/api/concurrent
```

#### 1. **Basic Goroutines with WaitGroup**
**Endpoint:** `GET /api/concurrent/parallel?ids=1,2,3`

Process multiple items simultaneously using goroutines and `sync.WaitGroup`.

**Use Cases:**
- Parallel data fetching from database
- Batch processing of independent tasks
- Concurrent API calls

**Key Concepts:**
- `go func()` to launch goroutines
- `sync.WaitGroup` to wait for completion
- `sync.Mutex` to protect shared data

#### 2. **Worker Pool Pattern**
**Endpoint:** `GET /api/concurrent/worker-pool?ids=1,2,3&workers=3`

Limit concurrent operations using a fixed number of workers processing jobs from a queue.

**Use Cases:**
- Rate limiting external API calls
- Database connection pooling
- Background job processing

**Key Concepts:**
- Job queue (buffered channel)
- Fixed number of worker goroutines
- Results collection channel

#### 3. **Fan-Out/Fan-In Pattern**
**Endpoint:** `GET /api/concurrent/fan-out-fan-in?query=golang`

Split work across multiple goroutines, then merge results.

**Use Cases:**
- Multi-source data aggregation
- Parallel searches across different fields
- Distributed computation

**Key Concepts:**
- Multiple goroutines produce results (fan-out)
- Single goroutine collects results (fan-in)
- Result deduplication

#### 4. **Pipeline Pattern**
**Endpoint:** `GET /api/concurrent/pipeline`

Process data through multiple stages using channels.

**Use Cases:**
- Multi-stage data processing
- ETL (Extract, Transform, Load) pipelines
- Stream processing

**Key Concepts:**
- Chained channels between stages
- Each stage is a separate goroutine
- Data flows through pipeline

#### 5. **Semaphore Pattern (Rate Limiting)**
**Endpoint:** `POST /api/concurrent/bulk-create`

Control concurrency using a buffered channel as a semaphore.

**Use Cases:**
- API rate limiting
- Resource pooling (e.g., database connections)
- Controlled parallel writes

**Key Concepts:**
- Buffered channel as semaphore
- Acquire/release pattern
- Concurrent operation limiting

**Example Request:**
```json
{
  "books": [
    {"title": "Book 1", "author": "Author 1", "isbn": "ISBN1"},
    {"title": "Book 2", "author": "Author 2", "isbn": "ISBN2"}
  ],
  "max_concurrent": 3
}
```

#### 6. **Timeout Pattern**
**Endpoint:** `GET /api/concurrent/timeout/1?timeout=5`

Cancel operations that exceed a time limit using context.

**Use Cases:**
- External API calls with timeout
- Slow database queries
- User-facing operations requiring responsiveness

**Key Concepts:**
- `context.WithTimeout()`
- `select` statement with `<-ctx.Done()`
- Graceful timeout handling

#### 7. **Select with Multiple Channels**
**Endpoint:** `GET /api/concurrent/monitor/1?interval=2&duration=10`

Handle multiple channel operations simultaneously.

**Use Cases:**
- Event handling systems
- Real-time monitoring
- Pub/sub implementations

**Key Concepts:**
- `select` statement for channel multiplexing
- Ticker for periodic operations
- Context for cancellation

### Testing Concurrent Patterns

To test these patterns:

1. **Start the application:**
   ```bash
   make docker-dev
   ```

2. **Login to get JWT token:**
   ```bash
   curl -X POST http://localhost:3000/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"admin123"}'
   ```

3. **Test a pattern (example: worker pool):**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://localhost:3000/api/concurrent/worker-pool?ids=1,2,3,4,5&workers=3"
   ```

### Best Practices for Concurrent Code

1. **Always close channels** when done producing to avoid goroutine leaks
2. **Use context for cancellation** to handle timeouts and cleanup
3. **Protect shared data** with mutexes or use channels for communication
4. **Avoid goroutine leaks** by ensuring all goroutines can exit
5. **Handle errors properly** in goroutines (use error channels)
6. **Use `defer` for cleanup** (e.g., `defer wg.Done()`)
7. **Test with race detector** (`go test -race`)

### Common Pitfalls to Avoid

- **Not waiting for goroutines to finish** → incomplete work
- **Accessing shared memory without synchronization** → race conditions
- **Forgetting to close channels** → goroutine leaks
- **Deadlocks from improper channel usage** → program hangs
- **Creating too many goroutines** → resource exhaustion

### Additional Resources

- Source code: `internal/services/concurrent_service.go`
- Handler implementations: `internal/handlers/concurrent.go`
- Route definitions: `internal/routes/routes.go:47-61`

## Development Checklist for New Code

- [ ] Test locally: `make test`
- [ ] Format code: `make fmt`
- [ ] Check for issues: `make vet` and `make lint`
- [ ] Verify migrations: `make migrate-status`
- [ ] Run in Docker: `make docker-dev`
- [ ] Add/update tests for new functionality
- [ ] Update `.env.example` if new env vars are needed
