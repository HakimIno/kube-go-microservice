package server

import (
	"context"
	"log"
	"time"

	"kube/internal/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// ServerConfig holds configuration for the server
type ServerConfig struct {
	Port         string
	ServiceName  string
	SwaggerURL   string
	RateLimit    int
	RateDuration time.Duration
}

// CommonServer provides common server setup and utilities
type CommonServer struct {
	*server.Hertz
	config ServerConfig
}

// NewServer creates a new server with common middleware and setup
func NewServer(config ServerConfig) *CommonServer {
	h := server.Default(server.WithHostPorts(":" + config.Port))

	// Add common middleware
	h.Use(middleware.StrictCORSConfig()) // CORS middleware should be first
	h.Use(LoggingMiddleware())
	h.Use(ErrorHandlerMiddleware())
	h.Use(RequestIDMiddleware())

	// Add rate limiting
	rateLimiter := NewRateLimiter(config.RateLimit, config.RateDuration)
	h.Use(rateLimiter.RateLimitMiddleware())

	// Add common health endpoints
	h.GET("/ping", PingHandler)
	h.GET("/health", HealthHandler(config.ServiceName))

	// Add swagger if URL is provided
	if config.SwaggerURL != "" {
		middleware.RegisterSwagger(h, config.SwaggerURL)
	}

	return &CommonServer{
		Hertz:  h,
		config: config,
	}
}

// Start starts the server with logging
func (s *CommonServer) Start() {
	log.Printf("%s starting on port %s", s.config.ServiceName, s.config.Port)
	if s.config.SwaggerURL != "" {
		log.Printf("Swagger UI available at: %s/swagger/index.html", s.config.SwaggerURL)
	}
	s.Spin()
}

// PingHandler handles ping requests
func PingHandler(ctx context.Context, c *app.RequestContext) {
	c.JSON(200, map[string]interface{}{
		"message": "pong",
	})
}

// HealthHandler creates a health check handler
func HealthHandler(serviceName string) func(context.Context, *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"service": serviceName,
			"status":  "healthy",
			"time":    time.Now().UTC(),
		})
	}
}
