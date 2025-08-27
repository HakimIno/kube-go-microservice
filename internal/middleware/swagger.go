package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// SwaggerHandler creates a simple Swagger UI handler
func SwaggerHandler(url string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// Simple Swagger configuration
		handler := swagger.WrapHandler(
			swaggerFiles.Handler,
			swagger.URL(url+"/swagger/doc.json"),
			swagger.DocExpansion("list"),
			swagger.PersistAuthorization(true),
		)
		handler(ctx, c)
	}
}

// RegisterSwagger registers Swagger routes
func RegisterSwagger(h *server.Hertz, url string) {
	h.GET("/swagger/*any", SwaggerHandler(url))
}
