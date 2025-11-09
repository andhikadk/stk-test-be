package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/andhikadk/stk-test-be/docs"

	"github.com/andhikadk/stk-test-be/config"
	"github.com/andhikadk/stk-test-be/internal/database"
	"github.com/andhikadk/stk-test-be/internal/middleware"
	"github.com/andhikadk/stk-test-be/internal/routes"
	"github.com/andhikadk/stk-test-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

// @title           STK Test API - Menu Management
// @version         1.0
// @description     REST API for hierarchical menu management with reordering and moving capabilities
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4000
// @BasePath  /
// @schemes   http https

func main() {
	migrateCmd := flag.String("migrate", "", "Run migrations (use: -migrate or -migrate sql)")
	seedCmd := flag.Bool("seed", false, "Seed database with sample data")
	statusCmd := flag.Bool("status", false, "Show migration status")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := utils.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if *migrateCmd != "" {
		if *migrateCmd == "sql" || *migrateCmd == "true" {
			log.Println("Running SQL migrations from embedded files...")
			if err := database.MigrateFromFS(db, MigrationsFS); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		}
		return
	}

	if *seedCmd {
		log.Println("Seeding database...")
		if err := database.SeedFromFS(db, MigrationsFS); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		log.Println("Seeding completed successfully")
		return
	}

	if *statusCmd {
		showMigrationStatus(db)
		return
	}

	if err := database.Migrate(db, cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:           cfg.AppName,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		EnablePrintRoutes: cfg.IsDevelopment(),
	})

	setupMiddleware(app, cfg)

	routes.SetupRoutes(app)

	startServer(app, cfg)
}

func showMigrationStatus(db *gorm.DB) {
	fmt.Println("\n=== Migration Status ===")

	migrator := database.NewMigrator(db)
	migrations, err := migrator.GetAppliedMigrations()
	if err != nil {
		log.Fatalf("Failed to get migration status: %v", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No migrations applied yet")
	} else {
		fmt.Println("Applied migrations:")
		for _, m := range migrations {
			fmt.Printf("  ✓ %s\n", m)
		}
	}

	seeder := database.NewSeeder(db)
	seeds, err := seeder.GetAppliedSeeds()
	if err != nil {
		log.Fatalf("Failed to get seed status: %v", err)
	}

	fmt.Println("\nApplied seeds:")
	if len(seeds) == 0 {
		fmt.Println("  No seeds applied yet")
	} else {
		for _, s := range seeds {
			fmt.Printf("  ✓ %s\n", s)
		}
	}
	fmt.Println()
}

func setupMiddleware(app *fiber.App, cfg *config.Config) {
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(recover.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowMethods: cfg.CORSAllowedMethods,
		AllowHeaders: cfg.CORSAllowedHeaders,
	}))

	app.Use(helmet.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	app.Use(middleware.ErrorHandlingMiddleware())
}

func startServer(app *fiber.App, cfg *config.Config) {
	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting %s on %s [%s mode]", cfg.AppName, address, cfg.Env)

	if err := app.Listen(address); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
