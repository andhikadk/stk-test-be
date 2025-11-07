package handlers

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, "API is running", fiber.Map{
		"app":     config.AppConfig.AppName,
		"status":  "healthy",
		"version": "1.0.0",
		"env":     config.AppConfig.Env,
	})
}
