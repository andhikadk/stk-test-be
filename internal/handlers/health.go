package handlers

import (
	"github.com/andhikadk/stk-test-be/config"
	"github.com/andhikadk/stk-test-be/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Check API health status
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, "API is running", fiber.Map{
		"app":     config.AppConfig.AppName,
		"status":  "healthy",
		"version": "1.0.0",
		"env":     config.AppConfig.Env,
	})
}
