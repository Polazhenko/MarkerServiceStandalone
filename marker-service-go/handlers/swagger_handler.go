package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SwaggerHandler handles Swagger UI and documentation
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// HealthCheck handles GET / for basic health check
func (h *SwaggerHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "MarkerService API",
	})
}
