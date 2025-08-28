package auth

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

	api := h.Group("/api/v1/auth")
	{
		authLimiter := middleware.AuthRateLimitMiddleware(5, time.Minute)
		qrLimiter := middleware.AuthRateLimitMiddleware(10, time.Minute)

		api.POST("/login", authLimiter, func(ctx context.Context, c *app.RequestContext) { handler.Login(c) })

		qr := api.Group("/qr")
		{
			// QR Code Login Routes - เรียบง่ายและชัดเจน
			qr.POST("/generate", qrLimiter, func(ctx context.Context, c *app.RequestContext) { handler.GenerateQRCode(c) })
			qr.POST("/scan", qrLimiter, func(ctx context.Context, c *app.RequestContext) { handler.QRScan(c) }) // Mobile app สแกน QR code
			qr.GET("/status", func(ctx context.Context, c *app.RequestContext) { handler.GetQRLoginStatus(c) })
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		{
			protected.POST("/refresh", func(ctx context.Context, c *app.RequestContext) { handler.RefreshToken(c) })
			protected.POST("/change-password", func(ctx context.Context, c *app.RequestContext) { handler.ChangePassword(c) })
		}
	}
}
