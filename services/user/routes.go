package user

import (
	"context"
	"time"

	"kube/internal/config"
	"kube/internal/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, service *Service) {
	handler := NewHandler(service)
	cfg := config.Load()

	// User routes
	api := h.Group("/api/v1/users")
	{
		// Public routes with rate limiting for security
		authLimiter := middleware.AuthRateLimitMiddleware(5, time.Minute) // 5 requests per minute for auth endpoints
		api.POST("/register", authLimiter, func(ctx context.Context, c *app.RequestContext) { handler.Register(c) })
		api.POST("/login", authLimiter, func(ctx context.Context, c *app.RequestContext) { handler.Login(c) })

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		{
			protected.POST("/refresh", func(ctx context.Context, c *app.RequestContext) { handler.RefreshToken(c) })
			protected.GET("/me", func(ctx context.Context, c *app.RequestContext) { handler.GetCurrentUser(c) })
			protected.POST("/change-password", func(ctx context.Context, c *app.RequestContext) { handler.ChangePassword(c) })
			protected.GET("/:id", func(ctx context.Context, c *app.RequestContext) { handler.GetUser(c) })
			protected.PUT("/:id", func(ctx context.Context, c *app.RequestContext) { handler.UpdateUser(c) })
			protected.DELETE("/:id", func(ctx context.Context, c *app.RequestContext) { handler.DeleteUser(c) })
		}

		// Admin routes (admin role required)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Add admin-specific routes here in the future
			// Example: admin.GET("/users", handler.GetAllUsers)
			// Example: admin.PUT("/users/:id/role", handler.UpdateUserRole)
		}
	}
}
