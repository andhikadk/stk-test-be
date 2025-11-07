package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-fiber-boilerplate/internal/models"

	"gorm.io/gorm"
)

// ConcurrentService demonstrates various concurrent programming patterns in Go
type ConcurrentService struct {
	db *gorm.DB
}

// NewConcurrentService creates a new concurrent service instance
func NewConcurrentService(db *gorm.DB) *ConcurrentService {
	return &ConcurrentService{db: db}
}

// ==============================
// PATTERN 1: Basic Goroutines with WaitGroup
// ==============================

// ProcessBooksParallel demonstrates parallel processing of books
// Use case: Fetch and process multiple books simultaneously
func (s *ConcurrentService) ProcessBooksParallel(bookIDs []uint) ([]models.Book, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex // Protects shared slice
	books := make([]models.Book, 0, len(bookIDs))
	errChan := make(chan error, len(bookIDs))

	for _, id := range bookIDs {
		wg.Add(1)
		go func(bookID uint) {
			defer wg.Done()

			var book models.Book
			if err := s.db.First(&book, bookID).Error; err != nil {
				errChan <- fmt.Errorf("failed to fetch book %d: %w", bookID, err)
				return
			}

			// Simulate some processing
			time.Sleep(100 * time.Millisecond)

			// Safely append to shared slice
			mu.Lock()
			books = append(books, book)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	close(errChan)

	// Collect any errors
	if len(errChan) > 0 {
		return books, <-errChan
	}

	return books, nil
}

// ==============================
// PATTERN 2: Worker Pool Pattern
// ==============================

// BookJob represents a job to process a book
type BookJob struct {
	ID     uint
	Action string
}

// BookResult represents the result of processing a book
type BookResult struct {
	Book  models.Book
	Error error
}

// ProcessBooksWithWorkerPool demonstrates worker pool pattern
// Use case: Process large number of tasks with limited workers
func (s *ConcurrentService) ProcessBooksWithWorkerPool(ctx context.Context, bookIDs []uint, numWorkers int) ([]models.Book, error) {
	jobs := make(chan BookJob, len(bookIDs))
	results := make(chan BookResult, len(bookIDs))

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go s.worker(ctx, w, jobs, results, &wg)
	}

	// Send jobs
	go func() {
		for _, id := range bookIDs {
			select {
			case jobs <- BookJob{ID: id, Action: "process"}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Close results after all workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	books := make([]models.Book, 0, len(bookIDs))
	for result := range results {
		if result.Error != nil {
			return books, result.Error
		}
		books = append(books, result.Book)
	}

	return books, nil
}

// worker processes jobs from the jobs channel
func (s *ConcurrentService) worker(ctx context.Context, id int, jobs <-chan BookJob, results chan<- BookResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-ctx.Done():
			results <- BookResult{Error: ctx.Err()}
			return
		default:
			var book models.Book
			if err := s.db.First(&book, job.ID).Error; err != nil {
				results <- BookResult{Error: fmt.Errorf("worker %d: failed to fetch book %d: %w", id, job.ID, err)}
				continue
			}

			// Simulate processing
			time.Sleep(200 * time.Millisecond)

			results <- BookResult{Book: book, Error: nil}
		}
	}
}

// ==============================
// PATTERN 3: Fan-Out/Fan-In Pattern
// ==============================

// SearchBooksMultipleSources demonstrates fan-out/fan-in pattern
// Use case: Query multiple data sources simultaneously and merge results
func (s *ConcurrentService) SearchBooksMultipleSources(query string) ([]models.Book, error) {
	// Fan-out: Start multiple goroutines
	ch1 := s.searchByTitle(query)
	ch2 := s.searchByAuthor(query)
	ch3 := s.searchByDescription(query)

	// Fan-in: Merge results from multiple channels
	books := make([]models.Book, 0)
	seen := make(map[uint]bool) // Deduplicate

	for i := 0; i < 3; i++ {
		select {
		case result := <-ch1:
			for _, book := range result {
				if !seen[book.ID] {
					books = append(books, book)
					seen[book.ID] = true
				}
			}
		case result := <-ch2:
			for _, book := range result {
				if !seen[book.ID] {
					books = append(books, book)
					seen[book.ID] = true
				}
			}
		case result := <-ch3:
			for _, book := range result {
				if !seen[book.ID] {
					books = append(books, book)
					seen[book.ID] = true
				}
			}
		}
	}

	return books, nil
}

func (s *ConcurrentService) searchByTitle(query string) <-chan []models.Book {
	ch := make(chan []models.Book, 1)
	go func() {
		defer close(ch)
		var books []models.Book
		s.db.Where("title ILIKE ?", "%"+query+"%").Find(&books)
		ch <- books
	}()
	return ch
}

func (s *ConcurrentService) searchByAuthor(query string) <-chan []models.Book {
	ch := make(chan []models.Book, 1)
	go func() {
		defer close(ch)
		var books []models.Book
		s.db.Where("author ILIKE ?", "%"+query+"%").Find(&books)
		ch <- books
	}()
	return ch
}

func (s *ConcurrentService) searchByDescription(query string) <-chan []models.Book {
	ch := make(chan []models.Book, 1)
	go func() {
		defer close(ch)
		var books []models.Book
		s.db.Where("description ILIKE ?", "%"+query+"%").Find(&books)
		ch <- books
	}()
	return ch
}

// ==============================
// PATTERN 4: Pipeline Pattern
// ==============================

// ProcessBooksPipeline demonstrates pipeline pattern
// Use case: Multi-stage data processing
func (s *ConcurrentService) ProcessBooksPipeline(ctx context.Context) ([]models.Book, error) {
	// Stage 1: Fetch books
	booksChan := s.fetchAllBooks(ctx)

	// Stage 2: Filter books (e.g., published only)
	filteredChan := s.filterBooks(ctx, booksChan)

	// Stage 3: Enrich books (e.g., add additional data)
	enrichedChan := s.enrichBooks(ctx, filteredChan)

	// Collect final results
	books := make([]models.Book, 0)
	for book := range enrichedChan {
		books = append(books, book)
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return books, nil
}

func (s *ConcurrentService) fetchAllBooks(ctx context.Context) <-chan models.Book {
	out := make(chan models.Book)
	go func() {
		defer close(out)
		var books []models.Book
		if err := s.db.Find(&books).Error; err != nil {
			return
		}

		for _, book := range books {
			select {
			case out <- book:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func (s *ConcurrentService) filterBooks(ctx context.Context, in <-chan models.Book) <-chan models.Book {
	out := make(chan models.Book)
	go func() {
		defer close(out)
		for book := range in {
			// Filter logic (example: only books with non-empty title)
			if book.Title != "" {
				select {
				case out <- book:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func (s *ConcurrentService) enrichBooks(ctx context.Context, in <-chan models.Book) <-chan models.Book {
	out := make(chan models.Book)
	go func() {
		defer close(out)
		for book := range in {
			// Enrich logic (example: simulate adding metadata)
			time.Sleep(50 * time.Millisecond)

			select {
			case out <- book:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// ==============================
// PATTERN 5: Semaphore Pattern (Rate Limiting)
// ==============================

// BulkCreateBooksWithRateLimit demonstrates semaphore pattern for rate limiting
// Use case: Limit concurrent operations (e.g., API calls, DB writes)
func (s *ConcurrentService) BulkCreateBooksWithRateLimit(ctx context.Context, books []models.CreateBookRequest, maxConcurrent int) ([]models.Book, error) {
	sem := make(chan struct{}, maxConcurrent) // Semaphore
	var wg sync.WaitGroup
	var mu sync.Mutex
	createdBooks := make([]models.Book, 0, len(books))
	errChan := make(chan error, 1)

	for _, bookReq := range books {
		select {
		case <-ctx.Done():
			return createdBooks, ctx.Err()
		default:
		}

		wg.Add(1)
		go func(req models.CreateBookRequest) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			book := models.Book{
				Title:       req.Title,
				Author:      req.Author,
				ISBN:        req.ISBN,
				Description: req.Description,
			}

			if err := s.db.Create(&book).Error; err != nil {
				select {
				case errChan <- fmt.Errorf("failed to create book: %w", err):
				default:
				}
				return
			}

			mu.Lock()
			createdBooks = append(createdBooks, book)
			mu.Unlock()
		}(bookReq)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return createdBooks, err
	}

	return createdBooks, nil
}

// ==============================
// PATTERN 6: Timeout Pattern
// ==============================

// FetchBookWithTimeout demonstrates timeout pattern
// Use case: Operations that shouldn't take too long
func (s *ConcurrentService) FetchBookWithTimeout(bookID uint, timeout time.Duration) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan *models.Book, 1)
	errChan := make(chan error, 1)

	go func() {
		var book models.Book
		if err := s.db.First(&book, bookID).Error; err != nil {
			errChan <- err
			return
		}

		// Simulate slow operation
		time.Sleep(2 * time.Second)

		resultChan <- &book
	}()

	select {
	case book := <-resultChan:
		return book, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("operation timed out after %v", timeout)
	}
}

// ==============================
// PATTERN 7: Select with Multiple Channels
// ==============================

// MonitorBookUpdates demonstrates select with multiple channels
// Use case: Background monitoring, event handling
func (s *ConcurrentService) MonitorBookUpdates(ctx context.Context, bookID uint, interval time.Duration) <-chan models.Book {
	updates := make(chan models.Book)

	go func() {
		defer close(updates)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var book models.Book
				if err := s.db.First(&book, bookID).Error; err == nil {
					select {
					case updates <- book:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return updates
}
