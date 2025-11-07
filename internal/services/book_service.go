package services

import (
	"errors"

	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/models"

	"gorm.io/gorm"
)

// BookService handles book business logic
type BookService struct {
	db *gorm.DB
}

// NewBookService creates a new book service
func NewBookService() *BookService {
	return &BookService{
		db: database.GetDB(),
	}
}

// GetAllBooks retrieves all books with pagination
func (s *BookService) GetAllBooks(page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Get total count
	if err := s.db.Model(&models.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// GetBookByID retrieves a book by ID
func (s *BookService) GetBookByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// CreateBook creates a new book
func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.Book, error) {
	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		Year:        req.Year,
		ISBN:        req.ISBN,
	}

	if err := s.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

// UpdateBook updates an existing book
func (s *BookService) UpdateBook(id uint, req *models.UpdateBookRequest) (*models.Book, error) {
	book, err := s.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	updateData := map[string]interface{}{}
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Author != "" {
		updateData["author"] = req.Author
	}
	if req.Year > 0 {
		updateData["year"] = req.Year
	}
	if req.ISBN != "" {
		updateData["isbn"] = req.ISBN
	}

	if err := s.db.Model(book).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return book, nil
}

// DeleteBook deletes a book (soft delete)
func (s *BookService) DeleteBook(id uint) error {
	if err := s.db.Delete(&models.Book{}, id).Error; err != nil {
		return err
	}
	return nil
}

// SearchBooks searches for books
func (s *BookService) SearchBooks(query string) ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Where("title ILIKE ? OR author ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
