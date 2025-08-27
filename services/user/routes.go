package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, service *Service) {
	handler := NewHandler(service)

	// User routes
	api := h.Group("/api/v1/users")
	{
		api.POST("/register", func(ctx context.Context, c *app.RequestContext) { handler.Register(c) })
		api.POST("/login", func(ctx context.Context, c *app.RequestContext) { handler.Login(c) })
		api.GET("/:id", func(ctx context.Context, c *app.RequestContext) { handler.GetUser(c) })
		api.PUT("/:id", func(ctx context.Context, c *app.RequestContext) { handler.UpdateUser(c) })
		api.DELETE("/:id", func(ctx context.Context, c *app.RequestContext) { handler.DeleteUser(c) })
	}
}
