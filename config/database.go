package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	switch c.DBDriver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			c.DBHost,
			c.DBPort,
			c.DBUser,
			c.DBPassword,
			c.DBName,
			c.DBSSLMode,
		)
	case "sqlite":
		return c.DBName + ".db"
	default:
		log.Fatalf("Unsupported database driver: %s", c.DBDriver)
		return ""
	}
}

// GetDialector returns the appropriate GORM dialector
func (c *Config) GetDialector() gorm.Dialector {
	switch c.DBDriver {
	case "postgres":
		return postgres.Open(c.GetDatabaseURL())
	case "sqlite":
		return sqlite.Open(c.GetDatabaseURL())
	default:
		log.Fatalf("Unsupported database driver: %s", c.DBDriver)
		return nil
	}
}

// GetGormLogLevel returns the appropriate GORM log level
func (c *Config) GetGormLogLevel() logger.LogLevel {
	switch c.LogLevel {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}
