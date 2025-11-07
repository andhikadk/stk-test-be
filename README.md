# Go Fiber Boilerplate

A production-ready boilerplate for building REST APIs with **Fiber**, a fast and lightweight Go web framework inspired by Express.js.

## 🚀 Features

- **Fiber Web Framework** - Fast, minimalist web framework
- **JWT Authentication** - Secure token-based authentication
- **GORM ORM** - Database abstraction layer
- **PostgreSQL & SQLite** - Multiple database support
- **Concurrent Programming Examples** - 7 production-ready concurrency patterns
- **Middleware Stack** - CORS, Logger, Recovery, Helmet
- **Request Validation** - Struct-based validation
- **Error Handling** - Centralized error management
- **Database Migrations** - Schema versioning
- **Unit Tests** - Testing setup ready
- **Docker Support** - Containerized deployment
- **Hot Reload** - Development mode with air
- **Environment Management** - .env configuration

## 📋 Project Structure

```
go-fiber-boilerplate/
├── main.go                        # Application entry point
├── embed.go                       # Embedded migrations
├── config/
│   ├── config.go                  # Configuration management
│   └── database.go                # Database setup
├── internal/
│   ├── handlers/                  # HTTP handlers
│   │   ├── auth.go                # Authentication handlers
│   │   ├── books.go               # Books CRUD handlers
│   │   ├── concurrent.go          # Concurrent patterns demo handlers
│   │   └── health.go              # Health check handlers
│   ├── models/                    # Data structures
│   ├── services/                  # Business logic
│   │   ├── auth_service.go        # Authentication service
│   │   ├── book_service.go        # Books service
│   │   └── concurrent_service.go  # Concurrent patterns service
│   ├── middleware/                # Custom middlewares
│   ├── database/                  # Database layer
│   └── routes/                    # Route definitions
├── pkg/
│   ├── utils/                     # Utility functions
│   └── jwt/                       # JWT utilities
├── migrations/                    # Database migrations (SQL files)
├── tests/                         # Test files
├── .env.example                   # Environment template
├── go.mod & go.sum                # Dependencies
├── Dockerfile                     # Multi-stage Docker build (production)
├── Dockerfile.dev                 # Development Docker with hot reload
├── docker-compose.yml             # Production Docker Compose configuration
├── docker-compose.dev.yml         # Development Docker Compose with hot reload
├── .air.toml                      # Air configuration for hot reload
├── .dockerignore                  # Docker build ignore rules
├── Makefile                       # Build and development commands
├── README.md                      # This file
└── CLAUDE.md                      # Claude Code instructions
```

## 🛠️ Tech Stack

- **Framework:** Fiber v2
- **Database:** GORM, PostgreSQL, SQLite
- **Authentication:** JWT (golang-jwt)
- **Security:** bcrypt (golang.org/x/crypto)
- **Validation:** go-playground/validator
- **Testing:** testify, standard library
- **Environment:** godotenv
- **Middleware:** Fiber built-in + custom

## 📦 Dependencies

### Core Dependencies
```
github.com/gofiber/fiber/v2 v2.52.5
github.com/gofiber/contrib/jwt v1.0.10
github.com/golang-jwt/jwt/v5 v5.2.1
gorm.io/gorm v1.31.0
gorm.io/driver/postgres v1.6.0
gorm.io/driver/sqlite v1.6.0
golang.org/x/crypto v0.43.0
github.com/go-playground/validator/v10 v10.28.0
github.com/joho/godotenv v1.5.1
```

## ⚡ Quick Start

### Prerequisites (Choose ONE based on your preferred setup)

**Option A: Docker + Hot Reload (Recommended)** ⭐
- Docker Desktop or Docker Engine
- Docker Compose v1.29+
- (No need to install Go, PostgreSQL, or Make locally!)

**Option B: Production-like (Docker)**
- Docker Desktop or Docker Engine
- Docker Compose v1.29+
- (No need to install Go, PostgreSQL, or Make locally!)

**Option C: Local Development (No Docker)**
- Go 1.25 or higher
- PostgreSQL 12+ (or SQLite)
- Make (optional, for using Makefile)
- git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/go-fiber-boilerplate.git
cd go-fiber-boilerplate
```

2. **Install dependencies**
```bash
make install-deps
# or
go mod download && go mod tidy
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Setup database and run application**

**Option A: Using Docker Compose with Hot Reload (Recommended for Development)** ⭐
```bash
make docker-dev
```
This will:
- Start PostgreSQL database (accessible at `localhost:6543`)
- Run application with hot reload (instant reload on code changes)
- Run migrations automatically
- Start the API on `http://localhost:4000`

> 💡 **Tip:** This is perfect for development! Edit your code and see changes instantly without rebuilding.

View logs:
```bash
make docker-dev-logs
```

**Option B: Using Docker Compose (Production-like)**
```bash
make docker-up
```
This will:
- Start PostgreSQL database (accessible at `localhost:6543`)
- Run migrations automatically
- Start the API on `http://localhost:4000`

> ⚠️ **Important:** The application is already running inside Docker! You do NOT need to run `make run` after this. The API is ready at `http://localhost:4000`

**Option C: Local Development (without Docker)**
```bash
# Create PostgreSQL database manually
createdb fiber_boilerplate

# Run migrations
make migrate

# Run application
make run
```
The API will be available at `http://localhost:4000`

## 🚀 Usage

### Run Application
```bash
make run
# or
go run main.go
```

### Build
```bash
make build
```
Binary will be created at `./bin/go-fiber-boilerplate`

### Development with Hot Reload
```bash
make dev
# Requires: go install github.com/cosmtrek/air@latest
```

### Testing
```bash
make test                # Run all tests
make test-coverage       # Run tests with coverage report
```

### Database
```bash
make migrate            # Run migrations
make seed               # Seed sample data
```

### Code Quality
```bash
make fmt                # Format code
make lint               # Run linter
make vet                # Run go vet
```

### Docker
```bash
make docker-build       # Build Docker image
make docker-up          # Start containers
make docker-down        # Stop containers
make docker-logs        # View logs
```

## 🔐 Authentication

This boilerplate uses JWT (JSON Web Tokens) for authentication:

1. **Register** - POST `/auth/register`
2. **Login** - POST `/auth/login` (returns JWT token)
3. **Protected Routes** - Add `Authorization: Bearer <token>` header

Tokens expire after 15 minutes by default. Adjust in `.env` with `JWT_EXPIRY`.

## 📚 API Endpoints (Example)

### Health Check
```
GET /health
```

### Authentication
```
POST /auth/register      # Register new user
POST /auth/login         # Login and get JWT token
POST /auth/refresh       # Refresh JWT token
```

### Books (Protected)
```
GET    /api/books        # List all books (requires auth)
GET    /api/books/:id    # Get book by ID (requires auth)
POST   /api/books        # Create book (requires auth)
PUT    /api/books/:id    # Update book (requires auth)
DELETE /api/books/:id    # Delete book (requires auth)
```

### Concurrent Patterns (Protected) 🆕
```
GET    /api/concurrent                    # Overview of all patterns
GET    /api/concurrent/parallel           # Parallel processing with goroutines
GET    /api/concurrent/worker-pool        # Worker pool pattern
GET    /api/concurrent/fan-out-fan-in     # Fan-out/fan-in pattern
GET    /api/concurrent/pipeline           # Pipeline pattern
POST   /api/concurrent/bulk-create        # Semaphore (rate limiting)
GET    /api/concurrent/timeout/:id        # Timeout pattern
GET    /api/concurrent/monitor/:id        # Select with multiple channels
```

## ⚡ Concurrent Programming Patterns

This boilerplate includes **7 production-ready concurrent programming patterns** to help developers understand and implement Go's concurrency features.

### 🎯 Why Learn Concurrency?

Go's concurrency model (goroutines and channels) is one of its most powerful features. Understanding these patterns will help you:

- **Build faster applications** - Process multiple tasks simultaneously
- **Handle high traffic** - Scale your API to serve thousands of requests
- **Implement background jobs** - Run tasks asynchronously without blocking
- **Control resource usage** - Prevent overwhelming your database or external APIs

### 📚 Available Patterns

| Pattern | Use Case | Example |
|---------|----------|---------|
| **Basic Goroutines + WaitGroup** | Parallel data fetching | Fetch multiple books simultaneously |
| **Worker Pool** | Rate limiting, job queues | Process tasks with limited workers |
| **Fan-Out/Fan-In** | Multi-source aggregation | Search across multiple fields in parallel |
| **Pipeline** | Multi-stage processing | ETL operations with stages |
| **Semaphore (Rate Limiting)** | API rate limiting | Limit concurrent database writes |
| **Timeout** | External API calls | Cancel slow operations |
| **Select with Multiple Channels** | Event handling | Monitor changes in real-time |

### 🚀 Quick Start

1. **Start the application:**
   ```bash
   make docker-dev
   ```

2. **Login to get JWT token:**
   ```bash
   curl -X POST http://localhost:4000/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"admin123"}'
   ```

3. **View all available patterns:**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:4000/api/concurrent
   ```

4. **Test a pattern (Worker Pool example):**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://localhost:4000/api/concurrent/worker-pool?ids=1,2,3,4,5&workers=3"
   ```

### 📖 Pattern Details

Each pattern includes:
- **Production-ready code** with comprehensive error handling
- **Context-based cancellation** for graceful shutdown
- **Detailed comments** explaining each step
- **Real-world use cases** demonstrated via API endpoints

All patterns are fully functional and can be tested immediately via the API endpoints above.

### 💡 Key Concepts

**Goroutines** - Lightweight threads managed by Go runtime
```go
go func() {
    // This runs concurrently
}()
```

**Channels** - Communication between goroutines
```go
ch := make(chan int)
ch <- 42        // Send
value := <-ch   // Receive
```

**Select** - Handle multiple channels
```go
select {
case msg := <-ch1:
    // Handle ch1
case msg := <-ch2:
    // Handle ch2
case <-time.After(5*time.Second):
    // Timeout
}
```

**Context** - Cancellation and timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 🎓 Learning Path

1. **Start with Pattern 1** (Basic Goroutines) - Understand fundamentals
2. **Try Worker Pool** - Learn resource control
3. **Explore Fan-Out/Fan-In** - Master result aggregation
4. **Practice with examples** - Test all 7 patterns via API
5. **Read source code** - Study implementation details
6. **Apply to your project** - Use patterns in real scenarios

### 📁 Source Files

- **Service Layer:** `internal/services/concurrent_service.go` - All pattern implementations
- **Handler Layer:** `internal/handlers/concurrent.go` - API endpoints for each pattern
- **Routes:** `internal/routes/routes.go` - Route definitions

## 📝 Configuration

All configuration is managed through `.env` file. See `.env.example` for all available options.

### Key Configuration
- `PORT` - Server port (default: 4000)
- `ENV` - Environment (development/production)
- `DB_HOST` - Database host (localhost for local, postgres for Docker)
- `DB_PORT` - Database port (5432 for local/Docker internal, 6543 for Docker host access)
- `DB_DRIVER` - Database driver (postgres/sqlite)
- `JWT_SECRET` - Secret key for JWT signing
- `CORS_ALLOWED_ORIGINS` - Allowed origins for CORS
- `LOG_LEVEL` - Logging level (info/debug/error)

## 🧪 Testing

Run all tests:
```bash
go test -v ./...
```

Run specific test:
```bash
go test -v ./tests -run TestName
```

With coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐳 Docker Deployment

### Development Setup with Hot Reload (Recommended)
```bash
make docker-dev
# or
docker-compose -f docker-compose.dev.yml up -d
```

This will:
- Start PostgreSQL database (port 6543 from host, 5432 internal)
- Run application with hot reload enabled (air watches for code changes)
- Run database migrations automatically
- Expose API on port 4000 (`http://localhost:4000`)

> 💡 **Benefits:** Changes to Go files are automatically detected and compiled. Just save your file and refresh the browser - no rebuild needed!

### Production Setup (No Hot Reload)
```bash
make docker-up
# or
docker-compose up -d
```

This will:
- Build the Fiber application into a compiled binary
- Start PostgreSQL database (port 6543 from host, 5432 internal)
- Run database migrations automatically
- Expose API on port 4000 (`http://localhost:4000`)

### View Logs
```bash
# For development setup
make docker-dev-logs
# or
docker-compose -f docker-compose.dev.yml logs -f

# For production setup
make docker-logs
# or
docker-compose logs -f

# View specific service logs
docker-compose logs -f fiber_app    # Just app logs
docker-compose logs -f postgres     # Just database logs
```

### Check Container Status
```bash
docker-compose ps
```

Expected output (both containers should be "Up"):
```
NAME                    STATUS
fiber_boilerplate_app   Up (healthy)
fiber_boilerplate_db    Up (healthy)
```

### Verify Application is Running
```bash
# Check if app is responding
curl http://localhost:4000/health

# Expected response: {"status":"ok"}
```

### Stop

**Development setup:**
```bash
make docker-dev-down
# or
docker-compose -f docker-compose.dev.yml down
```

**Production setup:**
```bash
make docker-down
# or
docker-compose down
```

### Reset Database (Remove all data and volumes)

**Development setup:**
```bash
make docker-dev-reset
# This removes development containers, networks, AND database volumes
# ⚠️ Warning: All data will be deleted!
```

**Production setup:**
```bash
make docker-reset
# This removes containers, networks, AND database volumes
# ⚠️ Warning: All data will be deleted!
```

## 📖 Project Structure Details

### `main.go`
Application entry point. Initializes config, database, and starts the server.

### `embed.go`
Embedded file system for database migrations. Uses Go's `embed` package to bundle migration files into the binary.

### `Dockerfile`
Production-ready multi-stage Docker build. Compiles Go code into an optimized binary with minimal dependencies.

### `Dockerfile.dev`
Development Docker image with hot reload support using air. Includes Go compiler and air tool for watching code changes.

### `.air.toml`
Configuration file for air hot reload tool. Specifies which files to watch, build commands, and reload behavior.

### `docker-compose.yml`
Production Docker Compose configuration. Runs compiled binary in isolated containers with PostgreSQL.

### `docker-compose.dev.yml`
Development Docker Compose configuration. Runs application with air hot reload and mounts source code as volume for instant updates.

### `config/`
Configuration management and database setup.

### `internal/handlers/`
HTTP request handlers for different routes.

### `internal/models/`
Data structures for the application.

### `internal/services/`
Business logic layer.

### `internal/middleware/`
Custom middleware for authentication, error handling, etc.

### `internal/database/`
Database connection and initialization.

### `internal/routes/`
Route definitions and grouping.

### `pkg/utils/`
Utility functions (response formatting, validation, etc).

### `pkg/jwt/`
JWT token generation and validation.

## 🔄 Development Workflow

> ⚠️ **Important:** Choose ONE path below. Do NOT run multiple paths at the same time - they will conflict on port 4000.

### With Docker + Hot Reload ⭐ (Recommended)
```bash
make docker-dev         # Start database and app with hot reload
make docker-dev-logs    # View logs (in another terminal)
make test              # Run tests in another terminal
make docker-dev-down   # Stop containers when done
```

**Features:**
- ✅ Instant reload on code changes (no rebuild needed)
- ✅ Database in Docker (simple setup)
- ✅ Production-like environment
- 💡 Perfect for rapid development

**When code changes:**
1. Edit file (e.g., `internal/handlers/books.go`)
2. Save file
3. Air automatically detects changes and rebuilds
4. Refresh `http://localhost:4000` → see your changes instantly!

### With Docker (Production-like, No Hot Reload)
```bash
make docker-up         # Start database and production-like app
make docker-logs       # View logs (in another terminal)
make test             # Run tests in another terminal
make docker-down      # Stop containers when done
```

**Features:**
- ✅ Production-ready build
- ✅ Uses compiled binary
- ❌ No hot reload (need rebuild on changes)

### Local Development (without Docker)
```bash
make run               # Start application with compiled binary
make test              # Run tests in another terminal
make dev               # OR run with hot reload (requires air installed locally)
make migrate           # Run migrations
make seed              # Seed sample data
```

**Features:**
- ✅ Hot reload with local air (`make dev`)
- ❌ Need to setup PostgreSQL locally
- 💡 Lightest setup, full control

**Note:** Make sure PostgreSQL is running locally before `make run` or `make dev`.

### Development Steps
1. **Create models** in `internal/models/`
2. **Create handlers** in `internal/handlers/`
3. **Add business logic** in `internal/services/`
4. **Define routes** in `internal/routes/`
5. **Write tests** in `tests/`
6. **Run and test** with `make run` and `make test`

## 📚 Learning Resources

### Framework & Libraries
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Guide](https://gorm.io/docs/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [Go Best Practices](https://golang.org/doc/effective_go)

### Concurrency (included in this boilerplate)
- **Source Code**: `internal/services/concurrent_service.go` - 7 production-ready patterns
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Concurrency Patterns (Video)](https://www.youtube.com/watch?v=f6kdp27TYZs)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the MIT License.

## 🙏 Acknowledgments

- Fiber team for the amazing framework
- GORM team for the powerful ORM
- Go community for best practices

---

**Happy coding! 🚀**

For issues and questions, please open an issue on GitHub.
