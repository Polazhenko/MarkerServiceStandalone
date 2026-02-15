package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/marker-service-go/models"
	"github.com/marker-service-go/service"
)

// MarkerHandler handles all marker-related HTTP requests
type MarkerHandler struct {
	service service.IMarkerService
}

// NewMarkerHandler creates a new marker handler
func NewMarkerHandler(svc service.IMarkerService) *MarkerHandler {
	return &MarkerHandler{
		service: svc,
	}
}

// CreateMarker handles POST /api/v1/markers
func (h *MarkerHandler) CreateMarker(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var request models.CreateMarkerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate fields
	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if request.Latitude < -90 || request.Latitude > 90 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude must be between -90 and 90"})
		return
	}

	if request.Longitude < -180 || request.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "longitude must be between -180 and 180"})
		return
	}

	if !models.IsValidCategory(request.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	// Create marker via service
	marker, err := h.service.CreateMarkerAsync(c.Request.Context(), userId, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return 201 Created with Location header
	c.Header("Location", fmt.Sprintf("/api/v1/markers/%s", marker.ID))
	c.JSON(http.StatusCreated, marker)
}

// GetMarker handles GET /api/v1/markers/{markerId}
func (h *MarkerHandler) GetMarker(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	markerId := c.Param("markerId")
	if markerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "markerId is required"})
		return
	}

	marker, err := h.service.GetMarkerAsync(c.Request.Context(), userId, markerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if marker == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Marker not found"})
		return
	}

	c.JSON(http.StatusOK, marker)
}

// GetAllMarkers handles GET /api/v1/markers
func (h *MarkerHandler) GetAllMarkers(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Parse optional category filter
	var category *models.MarkerCategory
	if categoryStr := c.Query("category"); categoryStr != "" {
		categoryVal, err := strconv.Atoi(categoryStr)
		if err != nil || categoryVal < 1 || categoryVal > 4 {
			// Invalid category - return empty list
			c.JSON(http.StatusOK, []models.MarkerResponse{})
			return
		}
		cat := models.MarkerCategory(categoryVal)
		category = &cat
	}

	markers, err := h.service.GetAllMarkersAsync(c.Request.Context(), userId, category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if markers == nil {
		markers = []*models.MarkerResponse{}
	}

	c.JSON(http.StatusOK, markers)
}

// UpdateMarker handles PATCH /api/v1/markers/{markerId}
func (h *MarkerHandler) UpdateMarker(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	markerId := c.Param("markerId")
	if markerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "markerId is required"})
		return
	}

	var request models.UpdateMarkerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate coordinates if provided
	if request.Latitude != nil && (*request.Latitude < -90 || *request.Latitude > 90) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude must be between -90 and 90"})
		return
	}

	if request.Longitude != nil && (*request.Longitude < -180 || *request.Longitude > 180) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "longitude must be between -180 and 180"})
		return
	}

	marker, err := h.service.UpdateMarkerAsync(c.Request.Context(), userId, markerId, &request)
	if err != nil {
		if err.Error() == "marker not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Marker not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, marker)
}

// DeleteMarker handles DELETE /api/v1/markers/{markerId}
func (h *MarkerHandler) DeleteMarker(c *gin.Context) {
	userId, err := GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	markerId := c.Param("markerId")
	if markerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "markerId is required"})
		return
	}

	success, err := h.service.DeleteMarkerAsync(c.Request.Context(), userId, markerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Marker not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
