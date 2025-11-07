package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ==============================
// PATTERN 1: Basic Goroutines with WaitGroup
// ==============================

// ProcessBooksParallel demonstrates parallel processing
// @Summary Process multiple books in parallel
// @Description Fetches and processes multiple books simultaneously using goroutines and WaitGroup
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated book IDs" example(1,2,3,4,5)
// @Success 200 {object} map[string]interface{} "Books processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/concurrent/parallel [get]
func ProcessBooksParallel(c *fiber.Ctx) error {
	// Parse book IDs from query parameter
	idsParam := c.Query("ids")
	if idsParam == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ids parameter is required (e.g., ?ids=1,2,3)")
	}

	idStrings := strings.Split(idsParam, ",")
	bookIDs := make([]uint, 0, len(idStrings))
	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid book ID: "+idStr)
		}
		bookIDs = append(bookIDs, uint(id))
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	start := time.Now()
	books, err := service.ProcessBooksParallel(bookIDs)
	duration := time.Since(start)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":         "Books processed successfully using parallel goroutines",
		"pattern":         "Basic Goroutines with WaitGroup",
		"books_count":     len(books),
		"books":           books,
		"processing_time": duration.String(),
		"note":            "Each book was fetched in a separate goroutine",
	})
}

// ==============================
// PATTERN 2: Worker Pool Pattern
// ==============================

// ProcessBooksWorkerPool demonstrates worker pool pattern
// @Summary Process books using worker pool
// @Description Processes books using a fixed number of workers for better resource control
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated book IDs" example(1,2,3,4,5)
// @Param workers query int false "Number of workers" default(3)
// @Success 200 {object} map[string]interface{} "Books processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/concurrent/worker-pool [get]
func ProcessBooksWorkerPool(c *fiber.Ctx) error {
	// Parse book IDs
	idsParam := c.Query("ids")
	if idsParam == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ids parameter is required (e.g., ?ids=1,2,3)")
	}

	idStrings := strings.Split(idsParam, ",")
	bookIDs := make([]uint, 0, len(idStrings))
	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid book ID: "+idStr)
		}
		bookIDs = append(bookIDs, uint(id))
	}

	// Parse number of workers
	numWorkers := 3 // default
	if workersParam := c.Query("workers"); workersParam != "" {
		w, err := strconv.Atoi(workersParam)
		if err != nil || w < 1 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid workers parameter")
		}
		numWorkers = w
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	books, err := service.ProcessBooksWithWorkerPool(ctx, bookIDs, numWorkers)
	duration := time.Since(start)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":         "Books processed successfully using worker pool",
		"pattern":         "Worker Pool Pattern",
		"books_count":     len(books),
		"books":           books,
		"workers":         numWorkers,
		"processing_time": duration.String(),
		"note":            "Limited number of workers process jobs from a queue",
	})
}

// ==============================
// PATTERN 3: Fan-Out/Fan-In Pattern
// ==============================

// SearchBooksMultipleSources demonstrates fan-out/fan-in
// @Summary Search books across multiple fields
// @Description Searches title, author, and description simultaneously and merges results
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {object} map[string]interface{} "Search completed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/concurrent/fan-out-fan-in [get]
func SearchBooksMultipleSources(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "query parameter is required")
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	start := time.Now()
	books, err := service.SearchBooksMultipleSources(query)
	duration := time.Since(start)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":         "Search completed successfully",
		"pattern":         "Fan-Out/Fan-In Pattern",
		"query":           query,
		"books_count":     len(books),
		"books":           books,
		"processing_time": duration.String(),
		"note":            "Searched title, author, and description in parallel, then merged results",
	})
}

// ==============================
// PATTERN 4: Pipeline Pattern
// ==============================

// ProcessBooksPipeline demonstrates pipeline pattern
// @Summary Process books through pipeline stages
// @Description Processes books through multiple stages: fetch -> filter -> enrich
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Pipeline completed successfully"
// @Failure 500 {object} map[string]interface{} "Processing failed"
// @Router /api/concurrent/pipeline [get]
func ProcessBooksPipeline(c *fiber.Ctx) error {
	db := database.GetDB()
	service := services.NewConcurrentService(db)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	books, err := service.ProcessBooksPipeline(ctx)
	duration := time.Since(start)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":         "Pipeline processing completed successfully",
		"pattern":         "Pipeline Pattern",
		"books_count":     len(books),
		"books":           books,
		"processing_time": duration.String(),
		"note":            "Books processed through stages: fetch -> filter -> enrich",
	})
}

// ==============================
// PATTERN 5: Semaphore Pattern (Rate Limiting)
// ==============================

// BulkCreateBooksRequest represents bulk create request
type BulkCreateBooksRequest struct {
	Books         []models.CreateBookRequest `json:"books" validate:"required,min=1,dive"`
	MaxConcurrent int                        `json:"max_concurrent" validate:"omitempty,min=1,max=10"`
}

// BulkCreateBooksWithRateLimit demonstrates semaphore pattern
// @Summary Bulk create books with rate limiting
// @Description Creates multiple books with controlled concurrency using semaphore
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param request body BulkCreateBooksRequest true "Bulk create request"
// @Success 201 {object} map[string]interface{} "Books created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/concurrent/bulk-create [post]
func BulkCreateBooksWithRateLimit(c *fiber.Ctx) error {
	var req BulkCreateBooksRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Set default max concurrent if not provided
	if req.MaxConcurrent == 0 {
		req.MaxConcurrent = 3
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	books, err := service.BulkCreateBooksWithRateLimit(ctx, req.Books, req.MaxConcurrent)
	duration := time.Since(start)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.CreatedResponse(c, fiber.Map{
		"message":         "Books created successfully with rate limiting",
		"pattern":         "Semaphore Pattern (Rate Limiting)",
		"books_count":     len(books),
		"books":           books,
		"max_concurrent":  req.MaxConcurrent,
		"processing_time": duration.String(),
		"note":            "Maximum " + strconv.Itoa(req.MaxConcurrent) + " concurrent operations allowed",
	})
}

// ==============================
// PATTERN 6: Timeout Pattern
// ==============================

// FetchBookWithTimeout demonstrates timeout pattern
// @Summary Fetch book with timeout
// @Description Fetches a book with a configurable timeout
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param timeout query int false "Timeout in seconds" default(5)
// @Success 200 {object} map[string]interface{} "Book fetched successfully"
// @Failure 408 {object} map[string]interface{} "Request timeout"
// @Router /api/concurrent/timeout/{id} [get]
func FetchBookWithTimeout(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	// Parse timeout
	timeoutSec := 5 // default
	if timeoutParam := c.Query("timeout"); timeoutParam != "" {
		t, err := strconv.Atoi(timeoutParam)
		if err != nil || t < 1 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid timeout parameter")
		}
		timeoutSec = t
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	start := time.Now()
	book, err := service.FetchBookWithTimeout(uint(id), time.Duration(timeoutSec)*time.Second)
	duration := time.Since(start)

	if err != nil {
		if strings.Contains(err.Error(), "timed out") {
			return utils.ErrorResponse(c, fiber.StatusRequestTimeout, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":         "Book fetched successfully",
		"pattern":         "Timeout Pattern",
		"book":            book,
		"timeout":         timeoutSec,
		"processing_time": duration.String(),
		"note":            "Operation will timeout if it exceeds " + strconv.Itoa(timeoutSec) + " seconds",
	})
}

// ==============================
// PATTERN 7: Select with Multiple Channels
// ==============================

// MonitorBookUpdates demonstrates select with multiple channels
// @Summary Monitor book updates (SSE-like)
// @Description Monitors a book for updates at regular intervals
// @Tags concurrent-examples
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param interval query int false "Interval in seconds" default(2)
// @Param duration query int false "Monitoring duration in seconds" default(10)
// @Success 200 {object} map[string]interface{} "Monitoring completed"
// @Router /api/concurrent/monitor/{id} [get]
func MonitorBookUpdates(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	// Parse interval
	intervalSec := 2 // default
	if intervalParam := c.Query("interval"); intervalParam != "" {
		i, err := strconv.Atoi(intervalParam)
		if err != nil || i < 1 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid interval parameter")
		}
		intervalSec = i
	}

	// Parse duration
	durationSec := 10 // default
	if durationParam := c.Query("duration"); durationParam != "" {
		d, err := strconv.Atoi(durationParam)
		if err != nil || d < 1 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid duration parameter")
		}
		durationSec = d
	}

	db := database.GetDB()
	service := services.NewConcurrentService(db)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(durationSec)*time.Second)
	defer cancel()

	updatesChan := service.MonitorBookUpdates(ctx, uint(id), time.Duration(intervalSec)*time.Second)

	// Collect updates
	updates := make([]models.Book, 0)
	for update := range updatesChan {
		updates = append(updates, update)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":       "Monitoring completed",
		"pattern":       "Select with Multiple Channels",
		"book_id":       id,
		"interval":      intervalSec,
		"duration":      durationSec,
		"updates_count": len(updates),
		"updates":       updates,
		"note":          "Book was monitored for " + strconv.Itoa(durationSec) + " seconds with " + strconv.Itoa(intervalSec) + " second intervals",
	})
}

// ==============================
// Overview Endpoint
// ==============================

// GetConcurrentPatterns returns overview of all concurrent patterns
// @Summary Get concurrent patterns overview
// @Description Returns information about all available concurrent programming patterns
// @Tags concurrent-examples
// @Produce json
// @Success 200 {object} map[string]interface{} "Patterns overview"
// @Router /api/concurrent [get]
func GetConcurrentPatterns(c *fiber.Ctx) error {
	patterns := []fiber.Map{
		{
			"pattern":     "Basic Goroutines with WaitGroup",
			"endpoint":    "GET /api/concurrent/parallel?ids=1,2,3",
			"description": "Process multiple items in parallel using goroutines and sync.WaitGroup",
			"use_case":    "Parallel data fetching, batch processing",
		},
		{
			"pattern":     "Worker Pool Pattern",
			"endpoint":    "GET /api/concurrent/worker-pool?ids=1,2,3&workers=3",
			"description": "Limit concurrent operations using a fixed number of workers",
			"use_case":    "Rate limiting, resource control, job queues",
		},
		{
			"pattern":     "Fan-Out/Fan-In Pattern",
			"endpoint":    "GET /api/concurrent/fan-out-fan-in?query=go",
			"description": "Split work across multiple goroutines, then merge results",
			"use_case":    "Multi-source data aggregation, parallel searches",
		},
		{
			"pattern":     "Pipeline Pattern",
			"endpoint":    "GET /api/concurrent/pipeline",
			"description": "Process data through multiple stages using channels",
			"use_case":    "Multi-stage data processing, ETL pipelines",
		},
		{
			"pattern":     "Semaphore Pattern (Rate Limiting)",
			"endpoint":    "POST /api/concurrent/bulk-create",
			"description": "Control concurrency using semaphore (buffered channel)",
			"use_case":    "API rate limiting, database connection pooling",
		},
		{
			"pattern":     "Timeout Pattern",
			"endpoint":    "GET /api/concurrent/timeout/1?timeout=5",
			"description": "Cancel operations that exceed time limit",
			"use_case":    "External API calls, slow queries, user-facing operations",
		},
		{
			"pattern":     "Select with Multiple Channels",
			"endpoint":    "GET /api/concurrent/monitor/1?interval=2&duration=10",
			"description": "Handle multiple channel operations simultaneously",
			"use_case":    "Event handling, monitoring, pub/sub systems",
		},
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":  "Concurrent programming patterns available in this boilerplate",
		"total":    len(patterns),
		"patterns": patterns,
		"note":     "These patterns demonstrate common Go concurrency patterns for production use",
	})
}
