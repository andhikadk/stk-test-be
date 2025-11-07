package utils

import (
	"go-fiber-boilerplate/internal/models"

	"github.com/gofiber/fiber/v2"
)

// SuccessResponse sends a success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := models.APIResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	response := models.APIResponse{
		Status:  statusCode,
		Message: message,
		Error:   message,
	}
	return c.Status(statusCode).JSON(response)
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *fiber.Ctx, message string, data interface{}, page, limit int, total int64) error {
	response := models.PaginatedResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    data,
		Page:    page,
		Limit:   limit,
		Total:   total,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// CreatedResponse sends a 201 created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return SuccessResponse(c, fiber.StatusCreated, message, data)
}

// BadRequestResponse sends a 400 bad request response
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

// UnauthorizedResponse sends a 401 unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

// ForbiddenResponse sends a 403 forbidden response
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

// NotFoundResponse sends a 404 not found response
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

// ConflictResponse sends a 409 conflict response
func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, message)
}

// InternalErrorResponse sends a 500 internal server error response
func InternalErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}
