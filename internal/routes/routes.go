package routes

import (
	"go-fiber-boilerplate/internal/handlers"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	apiGroup := app.Group("/api")
	{
		menusGroup := apiGroup.Group("/menus")
		{
			menusGroup.Get("/", handlers.GetMenus)
			menusGroup.Get("/:id", handlers.GetMenu)
			menusGroup.Post("/", handlers.CreateMenu)
			menusGroup.Put("/:id", handlers.UpdateMenu)
			menusGroup.Delete("/:id", handlers.DeleteMenu)
			menusGroup.Patch("/:id/move", handlers.MoveMenu)
			menusGroup.Patch("/:id/reorder", handlers.ReorderMenu)
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "endpoint not found",
		})
	})
}
