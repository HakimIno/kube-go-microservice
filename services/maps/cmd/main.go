package main

import (
	"kube/services/maps/internal/model"
	"kube/services/maps/internal/router"
	"kube/shared/pkg/config"
	"kube/shared/pkg/database"
	"kube/shared/pkg/middleware"
	"kube/shared/pkg/server"
	"time"
)

// @title Maps Service API
// @version 1.0
// @description This is a maps service API for location search using HERE API with Redis caching to reduce costs.

// @contact.name API Support
// @contact.url https://github.com/your-username/kube
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
// @schemes http https

func main() {
	cfg := config.Load()
	db := database.Init(cfg.Database)

	// Maps service doesn't need user tables, but we can keep it for future features
	if err := db.AutoMigrate(&model.User{}, &model.QRLoginSession{}); err != nil {
		middleware.LogError("Failed to migrate database", err)
		panic("Database migration failed")
	}

	serverConfig := server.ServerConfig{
		Port:         "8082",
		ServiceName:  "maps-service",
		SwaggerURL:   "http://localhost:8082",
		RateLimit:    100,
		RateDuration: time.Minute,
	}

	srv := server.NewServer(serverConfig)

	// Register routes using maps service router
	router.RegisterRoutes(srv.Hertz, db)

	srv.Start()
}
