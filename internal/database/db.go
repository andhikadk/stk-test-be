package database

import (
	"embed"
	"log"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize initializes the database connection
func Initialize(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		cfg.GetDialector(),
		&gorm.Config{
			Logger: logger.Default.LogMode(cfg.GetGormLogLevel()),
		},
	)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully")

	DB = db
	return db, nil
}

// Migrate runs database migrations
// Uses AutoMigrate for development, SQL migrations for production
func Migrate(db *gorm.DB, cfg *config.Config) error {
	log.Println("Running database migrations...")

	if cfg.IsDevelopment() {
		// Use AutoMigrate for fast development iteration
		log.Println("Using AutoMigrate for development mode")
		if err := db.AutoMigrate(
			&models.User{},
			&models.Book{},
		); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
			return err
		}
	} else {
		// Use SQL migrations for production
		log.Println("Using SQL migrations for production mode")
		// Note: Pass the migration files when calling this function
		// This will be handled in main.go
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// MigrateFromFS runs migrations from embedded filesystem
func MigrateFromFS(db *gorm.DB, migrations embed.FS) error {
	migrator := NewMigrator(db)
	return migrator.RunMigrationsFromFS(migrations)
}

// SeedFromFS seeds the database from embedded filesystem
func SeedFromFS(db *gorm.DB, seeds embed.FS) error {
	seeder := NewSeeder(db)
	return seeder.SeedFromFS(seeds)
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
