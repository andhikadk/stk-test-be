package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port         string
	Env          string
	AppName      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret        string
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration

	// CORS
	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Logging
	LogLevel string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		// Server
		Port:         getEnv("PORT", "3000"),
		Env:          getEnv("ENV", "development"),
		AppName:      getEnv("APP_NAME", "Fiber Boilerplate API"),
		ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "10s")),
		WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
		IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s")),

		// Database
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "fiber_boilerplate"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// JWT
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiry:        parseDuration(getEnv("JWT_EXPIRY", "15m")),
		JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "7d")),

		// CORS
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Validate critical config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	AppConfig = config
	return config, nil
}

// Validate validates critical configuration
func (c *Config) Validate() error {
	if c.JWTSecret == "your-secret-key-change-this" && c.Env == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	if c.DBDriver != "postgres" && c.DBDriver != "sqlite" {
		return fmt.Errorf("DB_DRIVER must be either 'postgres' or 'sqlite'")
	}

	return nil
}

// IsDevelopment checks if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction checks if environment is production
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// getEnv gets environment variable with fallback default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// parseDuration parses duration string with fallback
func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Warning: Invalid duration '%s', using default 10s", s)
		return 10 * time.Second
	}
	return duration
}
