package middleware

import (
	"go-fiber-boilerplate/internal/models"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Execute next handler
		err := c.Next()

		// Handle error if exists
		if err != nil {
			return handleError(c, err)
		}

		return nil
	}
}

// handleError processes different types of errors
func handleError(c *fiber.Ctx, err error) error {
	var code int
	var message string

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		// Generic error
		code = fiber.StatusInternalServerError
		message = "Internal Server Error"
	}

	response := models.APIResponse{
		Status:  code,
		Message: message,
		Error:   err.Error(),
	}

	return c.Status(code).JSON(response)
}
