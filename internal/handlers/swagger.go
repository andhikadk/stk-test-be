package handlers

import "github.com/gofiber/fiber/v2"

func SwaggerJSON(c *fiber.Ctx) error {
	paths := []string{
		"./docs/swagger.json",
		"docs/swagger.json",
		"/app/docs/swagger.json",
	}

	for _, path := range paths {
		if err := c.SendFile(path); err == nil {
			return nil
		}
	}

	return c.JSON(SwaggerSpec)
}

func SwaggerUI(c *fiber.Ctx) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.css">
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; padding: 0; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.10.5/swagger-ui-bundle.min.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "BaseLayout"
      })
    }
  </script>
</body>
</html>`
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}

var SwaggerSpec = map[string]interface{}{
	"swagger": "2.0",
	"info": map[string]interface{}{
		"title":       "Go Fiber Boilerplate API",
		"description": "A production-ready REST API boilerplate built with Fiber, GORM, and PostgreSQL",
		"version":     "1.0.0",
	},
	"host":     "localhost:4000",
	"basePath": "/",
	"paths": map[string]interface{}{
		"/health": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Health"},
				"summary":     "Health Check",
				"description": "Check API health status",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
					},
				},
			},
		},
		"/api/books": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Get all books",
				"description": "Get all books with optional search and pagination",
				"parameters": []map[string]interface{}{
					{
						"name":        "search",
						"in":          "query",
						"type":        "string",
						"description": "Search keyword",
					},
					{
						"name":        "page",
						"in":          "query",
						"type":        "integer",
						"default":     1,
						"description": "Page number",
					},
					{
						"name":        "limit",
						"in":          "query",
						"type":        "integer",
						"default":     10,
						"description": "Page size",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
			"post": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Create a book",
				"description": "Create a new book",
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Created",
					},
					"400": map[string]interface{}{
						"description": "Invalid request",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
		},
		"/api/books/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Get a book",
				"description": "Get a book by ID",
				"parameters": []map[string]interface{}{
					{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "integer",
						"description": "Book ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
					},
					"404": map[string]interface{}{
						"description": "Book not found",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
			"put": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Update a book",
				"description": "Update a book",
				"parameters": []map[string]interface{}{
					{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "integer",
						"description": "Book ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
					},
					"400": map[string]interface{}{
						"description": "Invalid request",
					},
					"404": map[string]interface{}{
						"description": "Book not found",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Delete a book",
				"description": "Delete a book",
				"parameters": []map[string]interface{}{
					{
						"name":        "id",
						"in":          "path",
						"required":    true,
						"type":        "integer",
						"description": "Book ID",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Book deleted successfully",
					},
					"404": map[string]interface{}{
						"description": "Book not found",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
		},
		"/api/books/search": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Books"},
				"summary":     "Search books",
				"description": "Search books by title, author, or ISBN",
				"parameters": []map[string]interface{}{
					{
						"name":        "q",
						"in":          "query",
						"required":    true,
						"type":        "string",
						"description": "Search keyword",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "OK",
					},
					"400": map[string]interface{}{
						"description": "Missing search query",
					},
					"500": map[string]interface{}{
						"description": "Internal server error",
					},
				},
			},
		},
	},
}
