package main

import (
	"github.com/gin-gonic/gin"

	"github.com/marker-service-go/handlers"
	"github.com/marker-service-go/repository"
	"github.com/marker-service-go/service"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Initialize repository and service
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)

	// Initialize handlers
	markerHandler := handlers.NewMarkerHandler(svc)
	swaggerHandler := handlers.NewSwaggerHandler()

	// Apply middleware to all routes
	router.Use(handlers.ExtractUserID())

	// Health check
	router.GET("/", swaggerHandler.HealthCheck)

	// Marker endpoints
	api := router.Group("/api/v1")
	{
		// Create marker
		api.POST("/markers", markerHandler.CreateMarker)

		// Get all markers
		api.GET("/markers", markerHandler.GetAllMarkers)

		// Get specific marker
		api.GET("/markers/:markerId", markerHandler.GetMarker)

		// Update marker
		api.PATCH("/markers/:markerId", markerHandler.UpdateMarker)

		// Delete marker
		api.DELETE("/markers/:markerId", markerHandler.DeleteMarker)
	}

	// Start server
	router.Run(":8080")
}
