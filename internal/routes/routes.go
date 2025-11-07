package routes

import (
	"go-fiber-boilerplate/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)

	app.Get("/swagger/doc.json", handlers.SwaggerJSON)
	app.Get("/swagger/index.html", handlers.SwaggerUI)

	apiGroup := app.Group("/api")
	{
		booksGroup := apiGroup.Group("/books")
		{
			booksGroup.Get("/", handlers.GetBooks)
			booksGroup.Get("/search", handlers.SearchBooks)
			booksGroup.Get("/:id", handlers.GetBook)
			booksGroup.Post("/", handlers.CreateBook)
			booksGroup.Put("/:id", handlers.UpdateBook)
			booksGroup.Delete("/:id", handlers.DeleteBook)
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "endpoint not found",
		})
	})
}
