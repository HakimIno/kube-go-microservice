package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// CustomSwaggerHandler creates a custom Swagger UI handler
func CustomSwaggerHandler(url string, options ...func(*swagger.Config)) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Fallback to default Swagger handler
		handler := swagger.WrapHandler(swaggerFiles.Handler, append([]func(*swagger.Config){swagger.URL(url)}, options...)...)
		handler(ctx, c)
	}
}

// RegisterCustomSwaggerRoutes registers custom Swagger routes
func RegisterCustomSwaggerRoutes(h *server.Hertz, url string, options ...func(*swagger.Config)) {
	handler := CustomSwaggerHandler(url, options...)
	h.GET("/swagger/*any", handler)
}
