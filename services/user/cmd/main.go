package main

import (
	"kube/services/user/internal/model"
	"kube/services/user/internal/router"
	"kube/shared/pkg/config"
	"kube/shared/pkg/database"
	"kube/shared/pkg/middleware"
	"kube/shared/pkg/server"
	"time"
)

// @title User Service API
// @version 1.0
// @description This is a user management service API built with Hertz framework.

// @contact.name API Support
// @contact.url https://github.com/your-username/kube
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	cfg := config.Load()
	db := database.Init(cfg.Database)

	if err := db.AutoMigrate(&model.User{}, &model.QRLoginSession{}); err != nil {
		middleware.LogError("Failed to migrate database", err)
		panic("Database migration failed")
	}

	serverConfig := server.ServerConfig{
		Port:         "8081",
		ServiceName:  "user-service",
		SwaggerURL:   "http://localhost:8081",
		RateLimit:    100,
		RateDuration: time.Minute,
	}

	srv := server.NewServer(serverConfig)

	// Register routes using user service router
	router.RegisterRoutes(srv.Hertz, db, cfg.JWT.SecretKey)

	srv.Start()
}
