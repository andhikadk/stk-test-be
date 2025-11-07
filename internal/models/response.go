package models

// APIResponse is the standard API response wrapper
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse is the response wrapper for paginated data
type PaginatedResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
}

// LoginRequest is the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest is the request body for registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse is the response for successful login
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// CreateBookRequest is the request body for creating a book
type CreateBookRequest struct {
	Title  string `json:"title" binding:"required,min=2"`
	Author string `json:"author" binding:"required,min=2"`
	Year   int    `json:"year" binding:"required,min=1000,max=9999"`
	ISBN   string `json:"isbn" binding:"required"`
}

// UpdateBookRequest is the request body for updating a book
type UpdateBookRequest struct {
	Title  string `json:"title" binding:"omitempty,min=2"`
	Author string `json:"author" binding:"omitempty,min=2"`
	Year   int    `json:"year" binding:"omitempty,min=1000,max=9999"`
	ISBN   string `json:"isbn" binding:"omitempty"`
}
