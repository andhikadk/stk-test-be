# Database Migrations & Seeding Guide

This directory contains all database migrations and seed data for the Go Fiber Boilerplate.

## Structure

```
migrations/
├── 001_initial_schema.sql      # Create initial tables
├── 002_add_indexes.sql         # Add performance indexes
├── seeds/
│   ├── 001_admin_user.sql      # Create default admin user
│   └── 002_sample_books.sql    # Seed sample books
└── README.md                   # This file
```

## Migration System

The boilerplate uses a **hybrid migration approach**:

- **Development**: Uses GORM's `AutoMigrate` for fast iteration
- **Production**: Uses SQL migration files for controlled deployments

### Migration Files

- Numbered sequentially (001, 002, 003...)
- Each file contains complete SQL for that migration
- Migrations are tracked in `migration_versions` table
- Prevents duplicate migrations

### How Migrations Work

1. **Automatic Migration on Startup**
   ```bash
   go run cmd/main.go
   # In development: Uses AutoMigrate (from models)
   # In production: Uses SQL migrations (from files)
   ```

2. **Manual SQL Migrations**
   ```bash
   make migrate-sql
   # or
   go run cmd/main.go -migrate=sql
   ```

3. **Check Migration Status**
   ```bash
   make migrate-status
   # or
   go run cmd/main.go -status
   ```

## Creating New Migrations

### Step 1: Create Migration File

Create a new SQL file with the next number:

```bash
# Create file: migrations/003_add_new_table.sql
```

**Example migration file:**

```sql
-- Add new_table for new feature
-- Created at: 2024-10-25

CREATE TABLE IF NOT EXISTS new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_new_table_name ON new_table(name);
```

### Step 2: Also Update GORM Models

For development mode, also update your models in `internal/models/`:

```go
// internal/models/new_model.go
type NewModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    Description string  `gorm:"type:text" json:"description"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Step 3: Run Migration

```bash
make migrate-sql
# or
go run cmd/main.go -migrate=sql
```

## Seeding Database

### Available Seeds

1. **001_admin_user.sql** - Creates default admin user
   - Email: `admin@example.com`
   - Password: `admin123` (⚠️ Change in production!)

2. **002_sample_books.sql** - Creates sample books

### Running Seeds

```bash
# Seed the database
make seed
# or
go run cmd/main.go -seed

# Check what seeds have been applied
make migrate-status
```

### Creating New Seeds

Create a new seed file in `migrations/seeds/`:

```bash
# Create file: migrations/seeds/003_seed_new_data.sql
```

**Example seed file:**

```sql
-- Seed new_table with sample data
-- Created at: 2024-10-25

INSERT INTO new_table (name, description)
VALUES
    ('Sample 1', 'First sample record'),
    ('Sample 2', 'Second sample record')
ON CONFLICT (name) DO NOTHING;
```

## Database Schema

### Tables Created

#### `users`
- id (PK)
- name
- email (UNIQUE)
- password (hashed)
- role (admin/user)
- is_active
- created_at
- updated_at
- deleted_at (soft delete)

#### `books`
- id (PK)
- title
- author
- isbn (UNIQUE)
- year
- pages
- publisher
- description
- created_at
- updated_at
- deleted_at (soft delete)

#### `migration_versions` (system table)
- id (PK)
- version
- applied_at

#### `seed_versions` (system table)
- id (PK)
- seed_name
- applied_at

## Common Commands

```bash
# Run migrations (dev: AutoMigrate, prod: SQL files)
make migrate

# Run SQL migrations explicitly
make migrate-sql

# Seed database
make seed

# Check status
make migrate-status

# View migration table in psql
psql -U postgres -d fiber_boilerplate -c "SELECT * FROM migration_versions;"

# View seed table in psql
psql -U postgres -d fiber_boilerplate -c "SELECT * FROM seed_versions;"
```

## Best Practices

### ✅ Do

- Use numbered migration files (001, 002, 003...)
- Include timestamp comment in each file
- Make migrations idempotent (use `IF NOT EXISTS`)
- Include both SQL migrations and model updates
- Test migrations in development first
- Keep seed data relevant and small
- Use transactions for complex migrations

### ❌ Don't

- Manually edit migration history
- Skip numbering in migration files
- Mix migration and seed operations
- Include data cleanup in migrations
- Delete or modify committed migration files

## Troubleshooting

### Migration Won't Run

Check if migration is already applied:

```bash
make migrate-status
```

### Reset Database (Development Only!)

```bash
# Drop and recreate database
dropdb fiber_boilerplate
createdb fiber_boilerplate

# Re-run migrations
make migrate
make seed
```

### Migration Failed Midway

1. Check the error message
2. Fix the SQL in the migration file
3. Delete the failed migration from `migration_versions` table:
   ```sql
   DELETE FROM migration_versions WHERE version = 'migration_name.sql';
   ```
4. Re-run migration

## Environment Variables

Configure database connection in `.env`:

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_boilerplate
DB_SSL_MODE=disable
```

## References

- [GORM Migrations](https://gorm.io/docs/migration.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Best Practices](https://wiki.postgresql.org/wiki/Performance_Optimization)
