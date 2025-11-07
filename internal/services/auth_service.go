package services

import (
	"errors"
	"time"

	"go-fiber-boilerplate/config"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/pkg/jwt"
	"go-fiber-boilerplate/pkg/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	db *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

// Register registers a new user
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := utils.VerifyPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)

	accessToken, err := tm.GenerateAccessToken(user.ID, user.Email, user.Role, config.AppConfig.JWTExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := tm.GenerateRefreshToken(user.ID, user.Email, config.AppConfig.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(config.AppConfig.JWTExpiry.Seconds()),
	}, nil
}

// RefreshToken generates a new access token from refresh token
func (s *AuthService) RefreshToken(refreshTokenString string) (string, error) {
	tm := jwt.NewTokenManager(config.AppConfig.JWTSecret)

	// Validate refresh token
	claims, err := tm.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return "", err
	}

	// Generate new access token
	accessToken, err := tm.GenerateAccessToken(user.ID, user.Email, user.Role, config.AppConfig.JWTExpiry)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *AuthService) UpdateUser(id uint, name string) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := utils.VerifyPassword(oldPassword, user.Password); err != nil {
		return errors.New("invalid password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := s.db.Model(user).Update("password", hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
