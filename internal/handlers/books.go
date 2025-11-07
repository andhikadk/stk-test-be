package handlers

import (
	"strconv"

	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetBooks retrieves all books
func GetBooks(c *fiber.Ctx) error {
	// Get pagination params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get books from service
	bookService := services.NewBookService()
	books, total, err := bookService.GetAllBooks(page, limit)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to fetch books")
	}

	return utils.PaginatedResponse(c, "Books retrieved successfully", books, page, limit, total)
}

// GetBook retrieves a specific book by ID
func GetBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid book ID")
	}

	bookService := services.NewBookService()
	book, err := bookService.GetBookByID(uint(id))
	if err != nil {
		return utils.NotFoundResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Book retrieved successfully", book)
}

// CreateBook creates a new book
func CreateBook(c *fiber.Ctx) error {
	var req models.CreateBookRequest

	// Parse and validate request
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Validate request data
	if req.Title == "" || req.Author == "" || req.ISBN == "" {
		return utils.BadRequestResponse(c, "title, author, and isbn are required")
	}

	if req.Year < 1000 || req.Year > 9999 {
		return utils.BadRequestResponse(c, "invalid year")
	}

	// Create book
	bookService := services.NewBookService()
	book, err := bookService.CreateBook(&req)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to create book")
	}

	return utils.CreatedResponse(c, "Book created successfully", book)
}

// UpdateBook updates an existing book
func UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid book ID")
	}

	var req models.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body")
	}

	// Validate year if provided
	if req.Year > 0 && (req.Year < 1000 || req.Year > 9999) {
		return utils.BadRequestResponse(c, "invalid year")
	}

	// Update book
	bookService := services.NewBookService()
	book, err := bookService.UpdateBook(uint(id), &req)
	if err != nil {
		return utils.NotFoundResponse(c, "book not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Book updated successfully", book)
}

// DeleteBook deletes a book
func DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.BadRequestResponse(c, "invalid book ID")
	}

	// Delete book
	bookService := services.NewBookService()
	if err := bookService.DeleteBook(uint(id)); err != nil {
		return utils.NotFoundResponse(c, "book not found")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Book deleted successfully", nil)
}

// SearchBooks searches for books
func SearchBooks(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return utils.BadRequestResponse(c, "search query is required")
	}

	bookService := services.NewBookService()
	books, err := bookService.SearchBooks(query)
	if err != nil {
		return utils.InternalErrorResponse(c, "search failed")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Search results", books)
}
