package database

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// MigrationFile represents a single migration file
type MigrationFile struct {
	Version string
	SQL     string
}

// Migrator handles SQL migrations
type Migrator struct {
	db    *gorm.DB
	files embed.FS
	path  string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:   db,
		path: "migrations",
	}
}

// RunMigrationsFromFS runs migrations from embedded filesystem
func (m *Migrator) RunMigrationsFromFS(files embed.FS) error {
	m.files = files

	// Ensure migration_versions table exists
	if err := m.ensureMigrationTable(); err != nil {
		return err
	}

	// Read migration files
	entries, err := files.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Get SQL migration files (numbered .sql files)
	var migrations []MigrationFile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Check if migration is already applied
		if m.isMigrationApplied(entry.Name()) {
			log.Printf("Migration %s already applied, skipping", entry.Name())
			continue
		}

		// Read migration file
		content, err := files.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, MigrationFile{
			Version: entry.Name(),
			SQL:     string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Execute migrations in order
	for _, migration := range migrations {
		if err := m.executeMigration(&migration); err != nil {
			return err
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// executeMigration executes a single migration
func (m *Migrator) executeMigration(migration *MigrationFile) error {
	log.Printf("Running migration: %s", migration.Version)

	// Execute SQL
	if err := m.db.Exec(migration.SQL).Error; err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
	}

	// Record migration as applied
	if err := m.recordMigration(migration.Version); err != nil {
		return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
	}

	log.Printf("Migration %s completed successfully", migration.Version)
	return nil
}

// ensureMigrationTable ensures the migration versions table exists
func (m *Migrator) ensureMigrationTable() error {
	return m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migration_versions (
			id SERIAL PRIMARY KEY,
			version VARCHAR(50) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

// recordMigration records a migration as applied
func (m *Migrator) recordMigration(version string) error {
	return m.db.Exec(
		"INSERT INTO migration_versions (version) VALUES (?)",
		version,
	).Error
}

// isMigrationApplied checks if a migration has been applied
func (m *Migrator) isMigrationApplied(version string) bool {
	var count int64
	m.db.Table("migration_versions").
		Where("version = ?", version).
		Count(&count)
	return count > 0
}

// GetAppliedMigrations returns all applied migrations
func (m *Migrator) GetAppliedMigrations() ([]string, error) {
	var versions []string
	err := m.db.Table("migration_versions").
		Order("applied_at ASC").
		Pluck("version", &versions).Error
	return versions, err
}

// RollbackLastMigration rolls back the last applied migration
func (m *Migrator) RollbackLastMigration() error {
	// Note: This is a simplified implementation
	// For proper rollback, you would need down migrations
	log.Println("Rollback functionality requires down migrations to be implemented")
	return fmt.Errorf("rollback not fully implemented")
}
