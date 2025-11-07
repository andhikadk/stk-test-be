package handlers

import (
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/middleware"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "name, email, and password are required")
	}

	if !utils.IsPasswordValid(req.Password) {
		return utils.BadRequestResponse(c, "password must be at least 6 characters")
	}

	// Register user
	authService := services.NewAuthService()
	user, err := authService.Register(&req)
	if err != nil {
		return utils.ConflictResponse(c, err.Error())
	}

	return utils.CreatedResponse(c, "User registered successfully", user.GetPublicUser())
}

// Login handles user login
func Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return utils.BadRequestResponse(c, "email and password are required")
	}

	// Authenticate user
	authService := services.NewAuthService()
	loginResp, err := authService.Login(&req)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", loginResp)
}

// RefreshToken refreshes the access token
func RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Refresh token
	authService := services.NewAuthService()
	newAccessToken, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"token": newAccessToken,
	})
}

// GetProfile retrieves current user profile
func GetProfile(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "invalid user")
	}

	// Get user
	authService := services.NewAuthService()
	user, err := authService.GetUserByID(userID)
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", user.GetPublicUser())
}

// UpdateProfile updates user profile
func UpdateProfile(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "invalid user")
	}

	var req struct {
		Name string `json:"name" binding:"required,min=2"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Update user
	authService := services.NewAuthService()
	user, err := authService.UpdateUser(userID, req.Name)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to update profile")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", user.GetPublicUser())
}

// ChangePassword changes user password
func ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "invalid user")
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	if req.OldPassword == req.NewPassword {
		return utils.BadRequestResponse(c, "new password must be different from old password")
	}

	// Change password
	authService := services.NewAuthService()
	if err := authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Password changed successfully", nil)
}
