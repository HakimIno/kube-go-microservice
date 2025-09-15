package router

import (
	"kube/services/maps/internal/handler"
	"kube/services/maps/internal/service"
	"kube/shared/pkg/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes registers all maps service routes
func RegisterRoutes(r *server.Hertz, db *gorm.DB) {
	// Initialize maps service
	mapsService := service.NewMapsService(db)

	// Initialize maps handler
	mapsHandler := handler.NewMapsHandler(mapsService)

	// API v1 group
	v1 := r.Group("/api/v1")

	// Maps routes (location search - no authentication required for public API)
	maps := v1.Group("/maps")
	{
		maps.POST("/search", mapsHandler.SearchPlaces)             // Search places with POST
		maps.GET("/search/query", mapsHandler.SearchPlacesByQuery) // Search places with GET query params
		maps.GET("/places/:id", mapsHandler.GetPlaceDetails)       // Get place details
	}

	// Authenticated maps routes
	authMaps := v1.Group("/maps")
	authMaps.Use(middleware.AuthMiddleware("")) // JWT secret would be passed here
	{
		authMaps.POST("/favorites", mapsHandler.AddFavoritePlace)          // Add favorite place
		authMaps.GET("/favorites", mapsHandler.GetFavoritePlaces)          // Get favorite places
		authMaps.DELETE("/favorites/:id", mapsHandler.RemoveFavoritePlace) // Remove favorite place
		authMaps.GET("/history", mapsHandler.GetSearchHistory)             // Get search history
	}
}
