package router

import (
	"kube/biz/handler"
	"kube/biz/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterMapsRoutes registers only maps-related routes
func RegisterMapsRoutes(r *server.Hertz, db *gorm.DB) {
	// Initialize maps service only
	mapsService := service.NewMapsService(db)

	// Initialize maps handler only
	mapsHandler := handler.NewMapsHandler(mapsService)

	// API v1 group
	v1 := r.Group("/api/v1")

	// Maps routes (location search - no authentication required for public API)
	maps := v1.Group("/maps")
	{
		maps.POST("/search", mapsHandler.SearchPlaces)             // Search places with POST
		maps.GET("/search/query", mapsHandler.SearchPlacesByQuery) // Search places with GET query params
	}
}
