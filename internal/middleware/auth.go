package middleware

import (
	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/pkg/jwt"
	"go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT token and extracts user info
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "missing authorization header")
		}

		// Extract token
		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		// Validate token
		tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
		claims, err := tm.ValidateAccessToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		// Store user info in context for next handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// OptionalAuthMiddleware validates JWT token if present but doesn't require it
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Token not provided, continue without authentication
			return c.Next()
		}

		// Try to extract and validate token
		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Invalid header format, continue without authentication
			return c.Next()
		}

		// Validate token
		tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)
		claims, err := tm.ValidateAccessToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			return c.Next()
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminMiddleware checks if user has admin role (must be used after AuthMiddleware)
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil || role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "admin access required")
		}
		return c.Next()
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uint, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0, fiber.ErrUnauthorized
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, fiber.ErrUnauthorized
	}

	return id, nil
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(c *fiber.Ctx) string {
	email := c.Locals("email")
	if email == nil {
		return ""
	}
	return email.(string)
}
