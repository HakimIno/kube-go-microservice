package router

import (
	"kube/services/user/internal/handler"
	"kube/services/user/internal/service"
	"kube/shared/pkg/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes registers all user service routes
func RegisterRoutes(r *server.Hertz, db *gorm.DB, jwtSecret string) {
	// Initialize services
	authService := service.NewAuthService(db, jwtSecret)
	userService := service.NewUserService(db)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// API v1 group
	v1 := r.Group("/api/v1")

	// Auth routes (authentication & account management)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)    // User registration
		auth.POST("/login", authHandler.Login)          // User login
		auth.POST("/refresh", authHandler.RefreshToken) // Refresh token
		auth.POST("/change-password", middleware.AuthMiddleware(jwtSecret), authHandler.ChangePassword)
		auth.POST("/logout", middleware.AuthMiddleware(jwtSecret), authHandler.Logout)
		auth.POST("/qr/generate", authHandler.GenerateQRCode)
		auth.POST("/qr/confirm", authHandler.QRConfirm)
		auth.POST("/qr/reject", authHandler.QRReject)
		auth.GET("/qr/status", authHandler.GetQRLoginStatus)
	}

	// User routes (user profile management - requires authentication)
	users := v1.Group("/users")
	{
		users.GET("/me", middleware.AuthMiddleware(jwtSecret), userHandler.GetCurrentUser)
		users.GET("/:id", middleware.AuthMiddleware(jwtSecret), userHandler.GetUser)
		users.PUT("/:id", middleware.AuthMiddleware(jwtSecret), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.AuthMiddleware(jwtSecret), userHandler.DeleteUser)
		users.GET("/", middleware.AuthMiddleware(jwtSecret), userHandler.ListUsers)
	}
}
