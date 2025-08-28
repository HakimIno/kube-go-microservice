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

	api := h.Group("/api/v1/users")
	{
		authLimiter := middleware.AuthRateLimitMiddleware(5, time.Minute)

		api.POST("/register", authLimiter, func(ctx context.Context, c *app.RequestContext) { handler.Register(c) })

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		{
			protected.GET("/me", func(ctx context.Context, c *app.RequestContext) { handler.GetCurrentUser(c) })
			protected.GET("/:id", func(ctx context.Context, c *app.RequestContext) { handler.GetUser(c) })
			protected.PUT("/:id", func(ctx context.Context, c *app.RequestContext) { handler.UpdateUser(c) })
			protected.DELETE("/:id", func(ctx context.Context, c *app.RequestContext) { handler.DeleteUser(c) })
		}

		// admin := api.Group("/admin")
		// admin.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		// admin.Use(middleware.RoleMiddleware("admin"))
		// {
		// 	// Add admin-specific routes here in the future
		// 	// Example: admin.GET("/users", handler.GetAllUsers)
		// 	// Example: admin.PUT("/users/:id/role", handler.UpdateUserRole)
		// }
	}
}
