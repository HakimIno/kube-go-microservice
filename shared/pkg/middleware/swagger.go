package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// SwaggerHandler creates a simple Swagger UI handler
func SwaggerHandler(url string, serviceName string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Path())

		// Handle different swagger JSON requests
		if strings.HasSuffix(path, "/doc.json") || strings.HasSuffix(path, "/maps-swagger.json") || strings.HasSuffix(path, "/user-swagger.json") {
			serveSwaggerDoc(c, serviceName)
			return
		}

		// For other swagger requests, use the standard handler
		handler := swagger.WrapHandler(
			swaggerFiles.Handler,
			swagger.URL(url+"/swagger/doc.json"),
			swagger.DocExpansion("list"),
			swagger.PersistAuthorization(true),
		)
		handler(ctx, c)
	}
}

// serveSwaggerDoc serves the appropriate swagger.json file based on service name
func serveSwaggerDoc(c *app.RequestContext, serviceName string) {
	var swaggerFile string

	switch serviceName {
	case "maps-service":
		swaggerFile = "docs/maps-swagger.json"
	case "user-service":
		swaggerFile = "docs/user-swagger.json"
	default:
		// Fallback to combined swagger file
		swaggerFile = "docs/swagger.json"
	}

	// Try multiple possible paths for the swagger file
	possiblePaths := []string{
		filepath.Join(".", swaggerFile),              // Current directory (for services running from root)
		filepath.Join("..", "..", swaggerFile),       // From services/*/cmd/ directory
		filepath.Join("..", "..", "..", swaggerFile), // From deeper nested directories
	}

	var file *os.File
	var err error
	var filePath string

	for _, path := range possiblePaths {
		file, err = os.Open(path)
		if err == nil {
			filePath = path
			break
		}
	}

	if err != nil {
		LogError(fmt.Sprintf("Failed to open swagger file: %s (tried paths: %v)", swaggerFile, possiblePaths), err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Swagger documentation not found",
		})
		return
	}
	defer file.Close()

	// Read and parse the JSON
	content, err := io.ReadAll(file)
	if err != nil {
		LogError("Failed to read swagger file", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to read swagger documentation",
		})
		return
	}

	// Validate JSON
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal(content, &swaggerDoc); err != nil {
		LogError("Invalid JSON in swagger file", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Invalid swagger documentation format",
		})
		return
	}

	LogInfo(fmt.Sprintf("Successfully loaded swagger file: %s", filePath))

	// Set content type and serve the JSON
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", content)
}

// RegisterSwagger registers Swagger routes
func RegisterSwagger(h *server.Hertz, url string, serviceName string) {
	h.GET("/swagger/*any", SwaggerHandler(url, serviceName))

	// Add specific routes for different swagger JSON files
	h.GET("/swagger/maps-swagger.json", func(ctx context.Context, c *app.RequestContext) {
		if serviceName == "maps-service" {
			serveSwaggerDoc(c, serviceName)
		} else {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Service not found",
			})
		}
	})

	h.GET("/swagger/user-swagger.json", func(ctx context.Context, c *app.RequestContext) {
		if serviceName == "user-service" {
			serveSwaggerDoc(c, serviceName)
		} else {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Service not found",
			})
		}
	})
}
