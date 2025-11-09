package testutil

import (
	"io"
	"log"
	"testing"

	"github.com/andhikadk/stk-test-be/internal/utils"

	"github.com/andhikadk/stk-test-be/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "file::memory:?cache=shared",
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect test database: %v", err)
	}

	if err := db.AutoMigrate(&models.Menu{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TeardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil && sqlDB != nil {
		sqlDB.Close()
	}
}

func InitTestLogger() {
	utils.InfoLogger = log.New(io.Discard, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	utils.ErrorLogger = log.New(io.Discard, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
