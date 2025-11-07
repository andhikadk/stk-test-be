package database

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

// SeedFromFS seeds database from embedded filesystem
func (s *Seeder) SeedFromFS(files embed.FS) error {
	// Create seed tracking table if not exists
	if err := s.ensureSeedTable(); err != nil {
		return err
	}

	// Read seed files
	entries, err := files.ReadDir("migrations/seeds")
	if err != nil {
		log.Println("No seeds directory found, skipping seeding")
		return nil
	}

	var seedFiles []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// Check if seed is already applied
		if s.isSeedApplied(entry.Name()) {
			log.Printf("Seed %s already applied, skipping", entry.Name())
			continue
		}

		seedFiles = append(seedFiles, entry.Name())
	}

	// Execute seeds in order
	for _, seedFile := range seedFiles {
		if err := s.executeSeed(files, seedFile); err != nil {
			log.Printf("Warning: Failed to execute seed %s: %v", seedFile, err)
			// Don't fail completely if a seed fails
			continue
		}
	}

	log.Println("Seeding completed")
	return nil
}

// executeSeed executes a single seed file
func (s *Seeder) executeSeed(files embed.FS, seedFile string) error {
	log.Printf("Running seed: %s", seedFile)

	// Read seed file
	content, err := files.ReadFile(filepath.Join("migrations/seeds", seedFile))
	if err != nil {
		return fmt.Errorf("failed to read seed file %s: %w", seedFile, err)
	}

	// Execute SQL
	if err := s.db.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("failed to execute seed %s: %w", seedFile, err)
	}

	// Record seed as applied
	if err := s.recordSeed(seedFile); err != nil {
		return fmt.Errorf("failed to record seed %s: %w", seedFile, err)
	}

	log.Printf("Seed %s completed successfully", seedFile)
	return nil
}

// ensureSeedTable ensures the seed tracking table exists
func (s *Seeder) ensureSeedTable() error {
	return s.db.Exec(`
		CREATE TABLE IF NOT EXISTS seed_versions (
			id SERIAL PRIMARY KEY,
			seed_name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

// recordSeed records a seed as applied
func (s *Seeder) recordSeed(seedName string) error {
	return s.db.Exec(
		"INSERT INTO seed_versions (seed_name) VALUES (?)",
		seedName,
	).Error
}

// isSeedApplied checks if a seed has been applied
func (s *Seeder) isSeedApplied(seedName string) bool {
	var count int64
	s.db.Table("seed_versions").
		Where("seed_name = ?", seedName).
		Count(&count)
	return count > 0
}

// ClearSeeds clears all applied seed records (development only!)
func (s *Seeder) ClearSeeds() error {
	return s.db.Exec("DELETE FROM seed_versions").Error
}

// GetAppliedSeeds returns all applied seeds
func (s *Seeder) GetAppliedSeeds() ([]string, error) {
	var seeds []string
	err := s.db.Table("seed_versions").
		Order("applied_at ASC").
		Pluck("seed_name", &seeds).Error
	return seeds, err
}
